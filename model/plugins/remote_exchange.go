package channels

// define three security communication
import (
	"crypto/rsa"
	"net/http"
	"net/url"
	"encoding/base64"
	"strings"
	"io/ioutil"
	"crypto/rand"
	
	"github.com/weihualiu/ewp/sessions"
	"github.com/weihualiu/ewp/utils"
	log "github.com/Sirupsen/logrus"
)

// 信道握手随机数远程协商
func ExchangeCallback(sessionid string, ciphersuite, cipher, clientrandom []byte) ([]byte, []byte, []byte, error) {
	session := sessions.GetSession(sessionid)
	// define client attr
	headers := http.Header{}
	headers.Add("Content-Type","application/x-www-form-urlencoded;charset=utf-8")
	
	params := url.Values{}
	params.Add("protocolVersion", "1.0")
	params.Add("mobKey", base64.StdEncoding.EncodeToString(cipher))
	params.Add("cipherSpec", base64.StdEncoding.EncodeToString(ciphersuite))
	params.Add("RNC", base64.StdEncoding.EncodeToString(clientrandom))
	clientvers := strings.Split(session.Client.OtaVer, "-")
	params.Add("clientVersion", clientvers[0]+"-"+clientvers[2]+"-"+clientvers[3])
	preMasterSecret := utils.BytesAppend([]byte{1}, []byte{0}, utils.Random2(46))
	publicKey,_ := sessions.GetSecOptions().ServerCert.PublicKey.(*rsa.PublicKey)
	log.Println("premastersecret length:", len(preMasterSecret))
	pms, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey, preMasterSecret)
	params.Add("PMS", base64.StdEncoding.EncodeToString(pms))
	log.Println("parameters:", params.Encode())
	
	resp, err := HttpDo("http://182.207.176.72:8080/ebws/MobileBank?tranCode=HKMB000000&userLocale=zh_HK&locale=zh_HK", headers, params)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(string(resp))
		return nil,nil,nil, err
	}
	log.Println("remote response body: ", string(resp))
	
	
	
	return nil, nil, nil, nil
	
}

func HttpDo(url string, headers http.Header, params url.Values) ([]byte, error){
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = headers
	
	resp, err := client.Do(req)
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return body, nil
}

