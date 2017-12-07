package basic

import (
	mstring "github.com/weihualiu/toolkit/string"
)

func enc_msg(typeval byte, data []byte) []byte {
	m := bytes.NewBuffer([]byte{typeval})
	m.Write(mstring.IntToBytes(len(data)))
	m.Write(data)
	return m
}

type clientHello struct {
	Version []int
	Random []byte
	SessionId []byte
	SeqGroup []byte
	CipherSuites []byte
	SerialNumber []byte
}

func clientHelloNew() *clientHello {
	return new(clientHello)
}

func (this *clientHello)Parse(msg []byte) {
	major := int(msg[0])
	minor := int(msg[1])
	this.Version = []int{major,minor}
	this.Random = msg[2:34]
	sid_len := int(msg[34])
	this.SessionId = msg[35:sid_len+35]
	cs_len := mstring.BytesToUInt32(msg[sid_len+35:sid_len+37])
	this.CipherSuites = msg[sid_len+37: int(cs_len) + sid_len + 37]
	sn_len := int(msg[int(cs_len) + sid_len + 37])
	this.SerialNumer = msg[int(cs_len) + sid_len + 38:int(cs_len) + sid_len + 37 + sn_len]
}

