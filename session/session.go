package session

// 会话结构

type Session struct {
	Sid string
	Expire int64 //失效时间，指定多长时间后失效
	ClientAddr string
	CreateTime string
	User *UserInfo
	SVer string // 信道版本
	Other struct{} //存储其它信息
}

type UserInfo struct {
	Status bool //是否登录
	Uid string
	
}

type CertInfo struct {
	PublicKey []byte
	Privatekey []byte
	PreMaster []byte
}


var sessions map[string]*Session

func Init() {
	sessions = make(map[string]*Session)
}


func NewSession() string {
	return ""
}
