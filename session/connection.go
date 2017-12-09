package session

// connection
import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"hash"
	"crypto/x509"
	"crypto/hmac"
	"io/ioutil"
	"encoding/pem"
	
	c "github.com/weihualiu/ewp/model/constant"
	g "github.com/weihualiu/ewp/conf"
	mstring "github.com/weihualiu/toolkit/string"
)

// 连接状态相关数据
type ConnectionState struct {
	Version SecurityVersion
	Verify bool
	ClientRandom []byte
	ServerRandom []byte
	SessionId string
	CipherS CipherSuite
	ClientPMS []byte
	ServerPMS []byte
	MasterSecret1 []byte
	MasterSecret2 []byte
	ClientCertificate []byte
	VerifyD *VerifyData
	InitContent []byte
}

func ConnectionStateInit() *ConnectionState {
	this := new(ConnectionState)
	this.VerifyD = VerifyDataNew()
	return this
}

// 用于验证签名
type VerifyData struct {
	Md5 hash.Hash
	Sha hash.Hash
	Buffer []byte
}

// sha,md5 init()
func VerifyDataNew() *VerifyData {
	return &VerifyData{Md5:md5.New(), Sha:sha1.New()}
}

// md5,sha update()
func (this *VerifyData)Update(content []byte) {
	this.Md5.Write(content)
	this.Sha.Write(content)
}

// md5,sha final()
func (this *VerifyData)Finish(secret []byte) {
	md5 := this.Md5.Sum(nil)
	sha := this.Sha.Sum(nil)
	//prf 伪随机数
	buf := bytes.NewBuffer(md5)
	buf.Write(sha)
	this.Buffer = prf(secret, []byte("server finished"), buf.Bytes(), 12)
}

func (this *VerifyData)FinishClient(client, secret []byte) bool {
	md5 := this.Md5.Sum(nil)
	sha := this.Sha.Sum(nil)
	buf := bytes.NewBuffer(md5)
	buf.Write(sha)
	
	this.Buffer = prf(secret, []byte("client finished"), buf.Bytes(), 12)
	return bytes.Equal(this.Buffer, client)
	
}

// 伪随机数算法
func prf(secret, label, seed []byte, wantedLength int) []byte {
	s1 := secret[0 : (len(secret)+1)/2]
	s2 := secret[len(secret)/2:]

	labelAndSeed := make([]byte, len(label)+len(seed))
	copy(labelAndSeed, label)
	copy(labelAndSeed[len(label):], seed)
	
	hashMD5 := md5.New
	hashSHA1 := sha1.New
	
	result := make([]byte, wantedLength)
	pHash(result, s1, labelAndSeed, wantedLength, hashMD5)
	result2 := make([]byte, len(result))
	pHash(result2, s2, labelAndSeed, wantedLength, hashSHA1)

	for i, b := range result2 {
		result[i] ^= b
	}
	
	return result
}

func pHash(result, secret, seed []byte, wantedLength int, hash func() hash.Hash) {
	// 定义的HMAC算法函数
	hm := func(a, b []byte) []byte {
		h := hmac.New(hash, a)
		h.Write(b)
		return h.Sum(nil)
	}
	// 定义的字节合并函数
	join := func(a,b []byte) []byte {
		tmp := make([]byte, len(a)+len(b))
		copy(tmp, a)
		copy(tmp[len(a):], b)
		return tmp
	}
	
	// 定义一个存储结果值的变量
	rbuf := make([]byte, 0)
	
	// 每次循环使用的变量
	buf := make([]byte, len(seed))
	copy(buf, seed)
	
	for len(rbuf) < wantedLength {
		a := hm(secret, buf)
		b := hm(secret, join(a, seed))
		copy(rbuf[len(rbuf):], b)
		buf = make([]byte, len(a))
		copy(buf, a)
	}
	
	// 截取长度，返回指定长度
	copy(result[0:wantedLength],rbuf[0:wantedLength])
	
}

type CipherSuite struct {
	Id uint16
	// the lengths, in bytes, of the key material needed for each component.
	//keyLen int
	//macLen int
	//ivLen  int
	//ka     func(version uint16) keyAgreement
	// flags is a bitmask of the suite* values, above.
	//flags  int
	//cipher func(key, iv []byte, isRead bool) interface{}
	//mac    func(version uint16, macKey []byte) macFunction
	//aead   func(key, fixedNonce []byte) cipher.AEAD
}

func (this CipherSuite)Bytes() []byte {
	return mstring.UInt16ToBytes(this.Id)
}

// 解析加密套件
func ParseCipherSuites(data []byte) []CipherSuite {
	return []CipherSuite{}
}

// 从加密套件组中选取一个
func CipherSuiteSelect(client []CipherSuite) CipherSuite {
	// 定义的已有的cipherSuites
	return CipherSuite{c.TLS_RSA_WITH_AES_256_CBC_SHA}
}

var cipherSuites = []*CipherSuite{
	{c.TLS_RSA_WITH_AES_256_CBC_MD5},
	{c.TLS_RSA_WITH_AES_256_CBC_SHA},
	{c.TLS_SM2_WITH_SM4_128_CBC_SM3}}


// 证书相关
type SecOptions struct {
	OldestVer []SecurityVersion  //支持的旧版本
	Verify bool  // 是否开启双向验证
	ServerCert *x509.Certificate // 用户服务证书
	//ServerKey []byte
	//CaCerts []byte
	//IssuerKey []byte
	//IssuerId []byte
}

func SecOptionsInit() *SecOptions {
	this := new(SecOptions)
	this.Verify = g.Config().Security.Verify
	this.OldestVer = []SecurityVersion{} // g.Config().Security.OldestVer
	serverCertBuf, _ := ioutil.ReadFile(g.Config().Security.ServerCertPath)
	block, _ := pem.Decode(serverCertBuf)
	if block == nil {
	    panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
	    panic("failed to parse certificate: " + err.Error())
	}
	this.ServerCert = cert
	return this
}

// 信道版本结构
// 例如 1.4
type SecurityVersion struct {
	Major uint8 //主版本
	Minor uint8 //次版本
}

func (this SecurityVersion)ToInt() int {
	return int(this.Major) * 100 + int(this.Minor)
}

// 选择版本，如果客户端版本低于服务端版本使用客户端版本；否则使用服务端版本
func SelectVersion(clientver SecurityVersion, oldver []SecurityVersion) SecurityVersion {
	return clientver
}

