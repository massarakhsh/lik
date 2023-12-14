package stack

import (
	"math/rand"
	"sync"
	"time"
)

const SessionDuration = time.Minute * 30

type itSession struct {
	id      int
	lastIn  time.Time
	session ItSession
}

type ItSession interface {
}

var sessionGate sync.Mutex
var sessionMap map[int]*itSession

func GetSession(id int) ItSession {
	purgeSessions()
	sessionGate.Lock()
	var session ItSession
	if isess, ok := sessionMap[id]; ok {
		session = isess.session
		isess.lastIn = time.Now()
	}
	sessionGate.Unlock()
	return session
}

func CreateSession(session ItSession) int {
	purgeSessions()
	sessionGate.Lock()
	id := 100000000 + rand.Intn(900000000)
	isess := &itSession{id: id, session: session}
	isess.lastIn = time.Now()
	sessionMap[id] = isess
	sessionGate.Unlock()
	return id
}

func purgeSessions() {
	sessionGate.Lock()
	if sessionMap == nil {
		sessionMap = make(map[int]*itSession)
	}
	for id, session := range sessionMap {
		if time.Since(session.lastIn) > SessionDuration {
			delete(sessionMap, id)
		}
	}
	sessionGate.Unlock()
}
