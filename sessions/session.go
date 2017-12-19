package sessions

// 会话结构
import (
	"time"
	"sync"

	log "github.com/Sirupsen/logrus"
	"bytes"
	"encoding/binary"
	"fmt"
	"crypto/md5"
	"github.com/weihualiu/ewp/utils"
)

type Session struct {
	Sid        string
	Expire     int64 //失效时间，指定多长时间后失效
	ClientAddr string
	CreateTime time.Time
	UpdateTime time.Time
	User       *UserInfo
	SVer       string   // 信道版本
	Other      struct{} //存储其它信息
	Conn       *ConnectionState
	Client     *ClientInfo
	CipherS    *CipherState
}

type UserInfo struct {
	Status bool //是否登录
	Uid    string
}

type CertInfo struct {
	PublicKey  []byte
	Privatekey []byte
	PreMaster  []byte
}

var (
	sessions map[string]*Session
	lock sync.Mutex
)

func createSid() string {
	point := time.Now().Unix()
	b := bytes.NewBuffer([]byte{})
	binary.Write(b, binary.BigEndian, point)
	m := md5.New()
	m.Write(b.Bytes())
	m.Write(utils.Random2(20))
	tmp := bytes.NewBuffer(m.Sum(nil))
	tmp.Write(utils.Random2(12))
	sid := fmt.Sprintf("%x", tmp.Bytes())
	return sid
}

func Init() {
	if sessions == nil {
		sessions = make(map[string]*Session)
	}
}

func NewSession() string {
	Init()
	newflag := true
	var id string
	for newflag{
		id = createSid()
		_, newflag = sessions[id]
	}

	lock.Lock()
	sessions[id] = &Session{
		Sid:        id,
		CreateTime: time.Now(),
		UpdateTime: time.Now()}
	lock.Unlock()
	return id
}

// 将连接状态写入会话中
func InsertConnState(sessionid string, conn *ConnectionState) {
	s, ok := sessions[sessionid]
	lock.Lock()
	if ok {
		s.Conn = conn
	}else {
		log.Error("not found session:",sessionid)
	}
	lock.Unlock()
}

// 将客户端信息写入会话中
func InsertClientInfo(sessionid string, info *ClientInfo) {
	s, ok := sessions[sessionid]
	lock.Lock()
	if ok {
		s.Client = info
	}else{
		log.Error("not found session:",sessionid)
	}
	lock.Unlock()
}

func GetSession(sessionid string) *Session {
	var s *Session
	for k, v := range sessions {
		if k == sessionid {
			s = v
			break
		}
	}
	return s
}

func SessionSelectId(sessionid string) string {
	if GetSession(sessionid) == nil {
		return NewSession()
	}
	return sessionid
}