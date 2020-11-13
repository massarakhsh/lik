package likapi

import (
	_ "fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/massarakhsh/lik"
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

type DataPage struct {
	Self     DataPager
	id       int
	atLast   time.Time
	isTrust  bool
	session  DataSessioner
	marshal  lik.Seter
	width    int
	height   int
	sizeNeed bool
}

type DataPager interface {
	GetSelf() DataPager
	ContinueToPage(to DataPager) int
	GetSession() DataSessioner
	GetPageId() int
	GetSessionId() int
	GetTrust() bool
	GetSize() (int, int)
	GetSizeFix() (int, int, bool)
	SetSize(width int, height int, fix bool)
	GetAt() time.Time
	StoreMarshal(set lik.Seter)
	GetMarshal() lik.Seter
	setAt(at time.Time)
	getSession() DataSessioner
	bindSession(id int, session DataSessioner)
	freeSession()
	setTrust(trust bool)
	markTime()
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

func (page *DataPage) GetSelf() DataPager {
	return page.Self
}

func (page *DataPage) ContinueToPage(to DataPager) int {
	semOn()
	page.session.appendPage(page, to)
	semOff()
	return to.GetPageId()
}

func (page *DataPage) GetPageId() int {
	return page.id
}

func (page *DataPage) GetSessionId() int {
	return page.getSession().GetSelf().IdSession
}

func (page *DataPage) GetSize() (int, int) {
	return page.width, page.height
}

func (page *DataPage) GetSizeFix() (int, int, bool) {
	need := page.sizeNeed
	page.sizeNeed = false
	return page.width, page.height, need
}

func (page *DataPage) SetSize(width int, height int, fix bool) {
	page.sizeNeed = (width != page.width || height != page.height) && !fix
	page.width = width
	page.height = height
}

func (page *DataPage) GetTrust() bool {
	return page.isTrust
}

func (page *DataPage) StoreMarshal(set lik.Seter) {
	if set != nil {
		if page.marshal == nil {
			page.marshal = lik.BuildSet()
		}
		for _, set := range set.Values() {
			page.marshal.SetItem(set.Val, set.Key)
		}
	}
}

func (page *DataPage) GetMarshal() lik.Seter {
	set := page.marshal
	page.marshal = nil
	return set
}

func (page *DataPage) getSession() DataSessioner {
	return page.session
}

func (page *DataPage) bindSession(id int, session DataSessioner) {
	page.id = id
	page.session = session
}

func (page *DataPage) freeSession() {
	page.id = 0
	page.session = nil
}

func (page *DataPage) setTrust(trust bool) {
	page.isTrust = trust
}

func (page *DataPage) GetAt() time.Time {
	return page.atLast
}

func (page *DataPage) setAt(at time.Time) {
	page.atLast = at
}

func (page *DataPage) GetSession() DataSessioner {
	return page.getSession()
}

func (page *DataPage) markTime() {
	at := time.Now()
	page.setAt(at)
	page.session.setAt(at)
}
