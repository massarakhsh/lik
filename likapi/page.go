package likapi

import (
	"github.com/massarakhsh/lik"
	"sync"
	"time"
)

type DataPage struct {
	Self     	DataPager
	Sync		sync.Mutex
	Stack   	[]Controller
	Collect   	map[string]Controller

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
	GetLevel() int
	SeekControl(index string) (int, Controller)
	SetControlStack(level int, controller Controller)
	setAt(at time.Time)
	getSession() DataSessioner
	bindSession(id int, session DataSessioner)
	freeSession()
	setTrust(trust bool)
	markTime()
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

func (page *DataPage) GetLevel() int {
	return len(page.Stack)
}

func (page *DataPage) SeekControl(index string) (int, Controller) {
	level := -1
	var ctrl Controller
	for nc,ctr := range page.Stack {
		if ctr.GetIndex() == index {
			level = nc
			ctrl = ctr
		}
	}
	if ctrl != nil {
		page.Sync.Lock()
		ctrl,_ = page.Collect[index]
		page.Sync.Unlock()
	}
	return level, ctrl
}

func (page *DataPage) SetControlStack(level int, controller Controller) {
	page.Sync.Lock()
	levold := len(page.Stack)
	if level <= levold {
		var ctrls []Controller
		for nc := 0; nc < level; nc++ {
			ctrls = append(ctrls, page.Stack[nc])
		}
		if controller != nil {
			ctrls = append(ctrls, controller)
			controller.Init("")
		}
		page.Stack = ctrls
	}
	if controller != nil {
		page.Collect[controller.GetIndex()] = controller
	}
	page.Sync.Unlock()
}

