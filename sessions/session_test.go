package sessions

import (
	"testing"
	"log"
)

func TestSessionNew(t *testing.T) {
	session := NewSession()
	if session == "" {
		t.Error("session create failed!")
	}
	session1 := NewSession()
	log.Print("session:", session, ", session1:", session1)
	if session1 == session {
		t.Errorf("s1:%s , s2: %s", session, session1)
	}

}

// 会话创建函数压力测试
func BenchmarkNewSession(b *testing.B) {
	for i := 0; i < b.N ;i++  {
		NewSession()
	}
}