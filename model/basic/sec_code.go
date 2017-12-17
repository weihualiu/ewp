package basic

import (
	"bytes"
	
	mstring "github.com/weihualiu/toolkit/string"
	"github.com/weihualiu/ewp/sessions"
	log "github.com/Sirupsen/logrus"
	"github.com/weihualiu/ewp/utils"
)

type message struct {
	msgblock map[byte][]byte
}

func messageNew(data []byte) *message {
	m := new(message)
	m.msgblock = make(map[byte][]byte, 100)
	m.parse(data)
	
	return m
}

// 信道握手报文组包
func encodemsg(typeval byte, data []byte) []byte {
	m := bytes.NewBuffer([]byte{typeval})
	m.Write(mstring.IntToBytes(len(data)))
	m.Write(data)
	return m.Bytes()
}

// 信道握手报文解包
func (this message)decode(typeval byte) ([]byte, bool){
	flag := false
	var content []byte
	for k,v := range this.msgblock {
		if k == typeval {
			content = v
			flag = true
			break
		}
	}
	return content, flag
	
}

func (this *message)parse(data []byte) {
	end := 0
	for end < len(data) {
		if len(data) < end + 5 {
			break
		}
		tval := byte(data[end])
		lenBuff := int(mstring.BytesToUInt32(data[end+1: end+5]))
		if end + 4 + lenBuff > len(data) {
			break
		}
		this.msgblock[tval] = data[end+5:end+5+lenBuff]
		end += (lenBuff + 5)
	}
	
}

type clientHello struct {
	Version sessions.SecurityVersion
	Random []byte
	SessionId string
	SeqGroup uint8
	CipherS []sessions.CipherSuite
	SerialNumber []byte
}

func clientHelloNew() *clientHello {
	return new(clientHello)
}

func (this *clientHello)Parse(msg []byte) {
	this.Version = sessions.SecurityVersion{uint8(msg[0]), uint8(msg[1])}
	this.Random = msg[2:34]
	if this.Version.ToInt() >= 104 {
		this.SeqGroup = uint8(msg[34])
		cs_len := mstring.BytesToUInt16(msg[35:37])
		// many ciphersuite
		this.CipherS = sessions.ParseCipherSuites(msg[37: int(cs_len) + 37])
		//this.CipherS = msg[37: int(cs_len) + 37]
		sn_len := int(msg[int(cs_len) + 37])
		this.SerialNumber = msg[int(cs_len) + 38:int(cs_len) + 38 + sn_len]
	}else{
		sid_len := int(msg[34])
		this.SessionId = mstring.BytesToString(msg[35:sid_len+35])
		cs_len := mstring.BytesToUInt16(msg[sid_len+35:sid_len+37])
		this.CipherS = sessions.ParseCipherSuites(msg[sid_len+37: int(cs_len) + sid_len + 37])
		sn_len := int(msg[int(cs_len) + sid_len + 37])
		this.SerialNumber = msg[int(cs_len) + sid_len + 38:int(cs_len) + sid_len + 38 + sn_len]
	}
	
	log.Println("client hello:", this)
}

func dec_cke(data []byte) (premaster, serverRandom, extension []byte, reflag bool) {
	if len(data) < 81 {
		return nil,nil,nil, false
	}
	premaster = data[0:48]
	serverRandom = data[48:80]
	elen := int(data[80])
	extension = data[81:elen+81]
	reflag = true
	return
}

func enc_pms(version sessions.SecurityVersion, rand []byte) []byte {
	b := bytes.NewBuffer([]byte{})
	b.Write([]byte{version.Major})
	b.Write([]byte{version.Minor})
	b.Write(rand)
	
	return b.Bytes()
}

func enc_ske(random, pms, mac, masterSecret []byte) []byte {
	b := bytes.NewBuffer(random)
	b.Write(pms)
	b.Write(mac)
	
	key := masterSecret[0:32]
	iv := masterSecret[32:48]
	
	encrypt, err := utils.AesEncrypt(b.Bytes(), key, iv)
	if err != nil {
		return nil
	}
	return encrypt
}
