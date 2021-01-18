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
		it.Close()
		return nil
	}
	return it
}

func (it *LikTel) open(server string, user string, password string) bool {
	it.ItClient, _ = telnet.DialTo(server)
	if it.ItClient == nil {
		return false
	}
	if _,ok := it.scanTo(":"); !ok {
		return false
	}
	it.ItClient.Write([]byte(user + "\r\n"))
	if _,ok := it.scanTo(":"); !ok {
		return false
	}
	it.ItClient.Write([]byte(password + "\r\n"))
	if _,ok := it.scanTo("#"); !ok {
		return false
	}
	return true
}

func (it *LikTel) scanTo(dlm string) (string,bool) {
	wait := []byte(dlm)
	answer := []byte{}
	word := []byte{}
	get := []byte{0}
	for len(word) < len(wait) {
		if n,err := it.ItClient.Read(get); err != nil {
			answer = append(answer, word...)
			return string(answer),false
		} else if n == 1 && get[0] == wait[len(word)] {
			word = append(word, get...)
		} else if n > 0 {
			answer = append(answer, word...)
			if get[0] == wait[0] {
				word = []byte{ get[0] }
			} else {
				answer = append(answer, get[0])
				word = []byte{}
			}
		}
	}
	return string(answer), true
}

func (it *LikTel) Close() {
	if it.ItClient != nil {
		it.ItClient.Close()
		it.ItClient = nil
	}
}

func (it *LikTel) Execute(cmd string) (string,bool) {
	if _,err := it.ItClient.Write([]byte(cmd + "\r\n")); err != nil {
		return "", false
	}
	return it.scanTo("#")
}

