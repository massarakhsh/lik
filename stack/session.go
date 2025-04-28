package stack

import (
	"math/rand"
	"sync"
	"time"
)

const SessionDuration = time.Minute * 30

type itSession struct {
	id     int
	lastIn time.Time
	self   ItSession
}

type ItSession interface {
	Close()
}

var sessionGate sync.Mutex
var sessionMap map[int]*itSession

func GetSession(id int) ItSession {
	purgeSessions()
	sessionGate.Lock()
	var session ItSession
	if isess, ok := sessionMap[id]; ok {
		session = isess.self
		isess.lastIn = time.Now()
	}
	sessionGate.Unlock()
	return session
}

func GetAllSessions() []ItSession {
	purgeSessions()
	var list []ItSession
	sessionGate.Lock()
	for _, sess := range sessionMap {
		list = append(list, sess.self)
	}
	sessionGate.Unlock()
	return list
}

func CreateSession(session ItSession) int {
	purgeSessions()
	sessionGate.Lock()
	id := 100000000 + rand.Intn(900000000)
	isess := &itSession{id: id, self: session}
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
			session.self.Close()
			delete(sessionMap, id)
		}
	}
	sessionGate.Unlock()
}
