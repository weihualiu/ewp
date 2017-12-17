package sessions

// 会话结构
import (
	"time"
)

type Session struct {
	Sid string
	Expire int64 //失效时间，指定多长时间后失效
	ClientAddr string
	CreateTime time.Time
	UpdateTime time.Time
	User *UserInfo
	SVer string // 信道版本
	Other struct{} //存储其它信息
	Conn *ConnectionState
	Client *ClientInfo
	CipherS *CipherState
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
	id := "1111111111111111111111111111111"
	sessions[id] = &Session{
		Sid:id,
		CreateTime: time.Now(),
		UpdateTime: time.Now()}
	return id
}


// 将连接状态写入会话中
func InsertConnState(sessionid string, conn *ConnectionState) {
	s, ok := sessions[sessionid]
	if ok {
		s.Conn = conn
	}
	
}

// 将客户端信息写入会话中
func InsertClientInfo(sessionid string, info *ClientInfo) {
	s, ok := sessions[sessionid]
	if ok {
		s.Client = info
	}
}

func GetSession(sessionid string) *Session {
	var s *Session
	for k,v := range sessions {
		if k == sessionid {
			s = v
			break
		}
	}
	return s
	
}