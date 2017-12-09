package basic

import (
	"bytes"
	
	mstring "github.com/weihualiu/toolkit/string"
	"github.com/weihualiu/ewp/session"
	log "github.com/Sirupsen/logrus"
)

func enc_msg(typeval byte, data []byte) []byte {
	m := bytes.NewBuffer([]byte{typeval})
	m.Write(mstring.IntToBytes(len(data)))
	m.Write(data)
	return m.Bytes()
}

type clientHello struct {
	Version session.SecurityVersion
	Random []byte
	SessionId string
	SeqGroup uint8
	CipherS []session.CipherSuite
	SerialNumber []byte
}

func clientHelloNew() *clientHello {
	return new(clientHello)
}

func (this *clientHello)Parse(msg []byte) {
	this.Version = session.SecurityVersion{uint8(msg[0]), uint8(msg[1])}
	this.Random = msg[2:34]
	if this.Version.ToInt() >= 104 {
		this.SeqGroup = uint8(msg[34])
		cs_len := mstring.BytesToUInt16(msg[35:37])
		// many ciphersuite
		this.CipherS = session.ParseCipherSuites(msg[37: int(cs_len) + 37])
		//this.CipherS = msg[37: int(cs_len) + 37]
		sn_len := int(msg[int(cs_len) + 37])
		this.SerialNumber = msg[int(cs_len) + 38:int(cs_len) + 38 + sn_len]
	}else{
		sid_len := int(msg[34])
		this.SessionId = mstring.BytesToString(msg[35:sid_len+35])
		cs_len := mstring.BytesToUInt16(msg[sid_len+35:sid_len+37])
		this.CipherS = session.ParseCipherSuites(msg[sid_len+37: int(cs_len) + sid_len + 37])
		sn_len := int(msg[int(cs_len) + sid_len + 37])
		this.SerialNumber = msg[int(cs_len) + sid_len + 38:int(cs_len) + sid_len + 38 + sn_len]
	}
	
	log.Println("client hello:", this)
}

