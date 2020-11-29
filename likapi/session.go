package likapi

import (
	_ "fmt"
	"math/rand"
	"sync"
	"time"
)

type DataSession struct {
	IdSession int
	AtLast    time.Time
	Uri       string
	IP        string
	CP        int
}

type DataSessioner interface {
	StartToPage(page DataPager) int
	GetSelf() *DataSession
	bindSession()
	getAt() time.Time
	setAt(at time.Time)
	appendPage(pred DataPager, next DataPager)
	inspectSession()
}

var (
	mutex        sync.Mutex
	initialized  bool
	ListSessions map[int]DataSessioner
	ListPages    map[int]DataPager
)

func FindPage(idp int) DataPager {
	semOn()
	page := findPage(idp)
	if page != nil {
		page.setTrust(true)
	}
	semOff()
	return page
}

func semOn() {
	mutex.Lock()
	if !initialized {
		rand.Seed(time.Now().UnixNano())
		ListSessions = make(map[int]DataSessioner)
		ListPages = make(map[int]DataPager)
		initialized = true
	}
}

func semOff() {
	mutex.Unlock()
}

func findPage(idp int) DataPager {
	page, ok := ListPages[idp]
	if !ok {
		page = nil
	} else if page != nil {
		page.markTime()
	}
	inspectSessions()
	return page
}

func inspectSessions() {
	for _, session := range ListSessions {
		session.inspectSession()
	}
}

func (session *DataSession) StartToPage(page DataPager) int {
	semOn()
	session.bindSession()
	session.appendPage(nil, page)
	semOff()
	return page.GetPageId()
}

func (session *DataSession) GetSelf() *DataSession {
	return session
}

func (session *DataSession) bindSession() {
	session.IdSession = 1 + int(rand.Int31n(999999999))
	ListSessions[session.IdSession] = session
}

func (session *DataSession) getAt() time.Time {
	return session.AtLast
}

func (session *DataSession) setAt(at time.Time) {
	session.AtLast = at
}

func (session *DataSession) appendPage(pred DataPager, next DataPager) {
	id := 1 + int(rand.Int31n(999999999))
	ListPages[id] = next.GetSelf()
	next.bindSession(id, session)
	sx, sy := 800, 600
	if pred != nil {
		sx, sy = pred.GetSize()
	}
	next.SetSize(sx, sy, false)
	session.CP++
	next.markTime()
}

func (session *DataSession) inspectSession() {
	now := time.Now()
	idsession := session.IdSession
	delsession := (now.Sub(session.getAt()) > 2*time.Hour)
	for id, page := range ListPages {
		if page.getSession().GetSelf().IdSession == idsession {
			if delsession || now.Sub(page.GetAt()) > 2*time.Hour {
				delete(ListPages, id)
				page.freeSession()
			}
		}
	}
	if delsession {
		delete(ListSessions, idsession)
		session.IdSession = 0
	}
}

