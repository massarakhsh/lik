package likssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type LikSSH struct {
	ItClient	*ssh.Client
}

func Open(server string, user string, password string, keyfile string) *LikSSH {
	it := &LikSSH{}
	if !it.open(server, user, password, keyfile) {
		return nil
	}
	return it
}

func (it *LikSSH) open(server string, user string, password string, keyfile string) bool {
	var auth []ssh.AuthMethod
	if password != "" {
		auth = append(auth, ssh.Password(password))
	}
	if keyfile != "" {
		key, err := ioutil.ReadFile(keyfile)
		if err != nil {
			return false
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return false
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return false
	}
	it.ItClient = client
	return true
}

func (it *LikSSH) Close() {
	if it.ItClient != nil {
		it.ItClient.Close()
		it.ItClient = nil
	}
}

func (it *LikSSH) Execute(cmd string) string {
	answer := ""
	if session, err := it.ItClient.NewSession(); err == nil {
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run(cmd); err == nil {
			answer = b.String()
		}
		session.Close()
	}
	return answer
}

