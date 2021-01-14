package liktel

import (
	"github.com/reiver/go-telnet"
)

type LikTel struct {
	ItClient	*telnet.Conn
}

func Open(server string, user string, password string) *LikTel {
	it := &LikTel{}
	if !it.open(server, user, password) {
		return nil
	}
	return it
}

func (it *LikTel) open(server string, user string, password string) bool {
	it.ItClient, _ = telnet.DialTo(server)
	if it.ItClient == nil {
		return false
	}
	return true
}

func SetTest() {
	conn, _ := telnet.DialTo("localhost:5555")
	conn.Write([]byte("hello world"))
	conn.Write([]byte("\n"))
}

func (it *LikTel) Close() {
	if it.ItClient != nil {
		it.ItClient = nil
	}
}

func (it *LikTel) Execute(cmd string) string {
	answer := ""
	return answer
}

