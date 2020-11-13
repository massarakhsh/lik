package likapi

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
)

type DataDriver interface {
	Clear()
	//Clone() DataDriver
	GetPage() DataPager
	BuildPath() string
	ToString() string
	Top() string
	GetPath() []string
	GetExt() string
	GetIP() string
	GetLogin() string
	GetPassword() string
	Shift() string
	ShiftAll() string
	IsShift(part string) bool
	ItShift(part string) string
	IsShiftInt() (int, bool)
	IsEmpty() bool
	IsMethod(method string) bool
	IsAtMethod(method string) bool
	InitializePage(version string) likdom.Domer
	StoreItem(item ...likdom.Domer)
	SetGoPart(part string)
	SetOnPart(part string)
	BuildUrl(part string) string
	SeekPageSize() bool
	LoadRequest(r *http.Request)
	LoadContext(data string, isget bool)
	GetContext(key string) string
	SetContext(val interface{}, key string)
	GetBuffers() map[string][]byte
	SetResponse(val interface{}, key string)
	GetCookie(key string) string
	SetCookie(val string, key string)
	GetAllContext() lik.Seter
	GetAllResponse() lik.Seter
	GetAllCookies() lik.Seter
}

type DataDrive struct {
	Page     DataPager
	RootUri  string
	RootUrl  string
	site     string
	port     int
	stack    []string
	path     []string
	buffers  map[string][]byte
	ext      string
	method   string
	ip       string
	login    string
	password string
	isformat bool

	mutex    sync.Mutex
	context  lik.Seter
	responce lik.Seter
	cookies  lik.Seter
}

func GetHttpAccount(r *http.Request) (string, string) {
	user, pass := "", ""
	if auth := r.Header.Get("Authorization"); auth != "" {
		if match := lik.RegExParse(auth, "Basic\\s+(\\S+)"); match != nil {
			if deco, err := base64.StdEncoding.DecodeString(match[1]); err == nil {
				if match := lik.RegExParse(string(deco), "([^:]*):(.*)"); match != nil {
					user = match[1]
					pass = match[2]
				}
			}
		}
	}
	return user, pass
}

func (drive *DataDrive) Clear() {
	drive.mutex.Lock()
	drive.port = 0
	drive.site = ""
	drive.port = 0
	drive.ip = ""
	drive.stack = []string{}
	drive.path = []string{}
	drive.ext = ""
	drive.method = ""
	drive.context = nil
	drive.responce = nil
	//drive.cookies = nil
	drive.mutex.Unlock()
}

func (drive *DataDrive) GetPage() DataPager {
	return drive.Page
}

func (drive *DataDrive) LoadRequest(r *http.Request) {
	drive.Clear()
	drive.ip = r.Header.Get("X-Real-IP")
	if drive.ip == "" {
		drive.ip = r.RemoteAddr
	}
	if drive.login == "" {
		drive.login, drive.password = GetHttpAccount(r)
	}
	if sess := drive.Page.GetSession().GetSelf(); sess != nil && sess.IP == "" {
		sess.IP = drive.ip
	}
	drive.RootUri = r.RequestURI
	if host, _ := os.Hostname(); host == "Shaman" {
		drive.RootUrl = "http"
	} else {
		drive.RootUrl = "https"
	}
	drive.RootUrl += "://" + r.Host
	site := r.Host
	if pos := strings.Index(site, ":"); pos >= 0 {
		if port, ok := lik.StrToIntIf(site[pos+1:]); ok && port > 0 {
			drive.port = int(port)
		}
		site = site[:pos]
	}
	drive.site = site
	if drive.port == 0 {
		drive.port = 80
	}
	for _, cookie := range r.Cookies() {
		drive.SetCookie(cookie.Value, cookie.Name)
	}
	var path = r.RequestURI
	if pos := strings.Index(path, "?"); pos >= 0 {
		parm := path[pos+1 : len(path)]
		path = path[0:pos]
		drive.LoadContext(parm, true)
	}
	if drive.GetContext("_mf") == "" {
		if body, err := ioutil.ReadAll(r.Body); err == nil && body != nil {
			drive.LoadContext(string(body), false)
		}
	} else {
		r.ParseMultipartForm(10 << 20)
		if file, handler, err := r.FormFile("file"); err == nil {
			//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			//fmt.Printf("File Size: %+v\n", handler.Size)
			//fmt.Printf("MIME Header: %+v\n", handler.Header)
			if fileBytes, err := ioutil.ReadAll(file); err == nil {
				//fmt.Printf("Loaded: %d\n", len(fileBytes))
				drive.buffers = make(map[string][]byte)
				drive.buffers[handler.Filename] = fileBytes
			}
			file.Close()
		}
	}

	if mtch := regexp.MustCompile("\\.(.*)$").FindStringSubmatch(path); mtch != nil {
		drive.ext = mtch[1]
	}
	drive.path = []string{}
	for _, part := range strings.Split(path, "/") {
		if len(part) > 0 {
			drive.path = append(drive.path, part)
		}
	}
	drive.method = r.Method
}

func (drive *DataDrive) LoadContext(data string, isget bool) {
	if lik.RegExCompare(data, "^{") {
		if set := lik.SetFromRequest(data); set != nil {
			drive.SetContext(set, "_json")
		}
	} else if len(data) > 0 {
		//if strings.Contains(data, "data=") {
		//	lik.SayInfo("Query: " + data)
		//}
		data = strings.Replace(data, "+", "#", -1)
		data, _ := url.QueryUnescape(data)
		data = strings.Replace(data, "#", "+", -1)
		//if strings.Contains(data, "data=") {
		//	lik.SayInfo("Unescape: " + data)
		//}
		words := strings.Split(data, "&")
		for _, word := range words {
			if peq := strings.Index(word, "="); peq > 0 {
				key := word[:peq]
				val := word[peq+1:]
				//if strings.Contains(data, "data=") {
				//	lik.SayInfo("key: " + key + ", val: " + val)
				//}
				drive.SetContext(val, key)
			}
		}
	}
}

func (drive *DataDrive) GetBuffers() map[string][]byte {
	return drive.buffers
}

func (drive *DataDrive) GetIP() string {
	return drive.ip
}

func (drive *DataDrive) GetLogin() string {
	return drive.login
}

func (drive *DataDrive) GetPassword() string {
	return drive.password
}

func (drive *DataDrive) GetPath() []string {
	return drive.path
}

func (drive *DataDrive) BuildPath() string {
	return strings.Join(drive.stack, "/") + "/" + strings.Join(drive.path, "/")
}

func (drive *DataDrive) GetExt() string {
	return drive.ext
}

func (drive *DataDrive) ToString() string {
	var text = "Drive:\n"
	site := drive.site
	if drive.port > 0 {
		site += ":" + fmt.Sprint(drive.port)
	}
	text += "Site: " + site + "\n"
	if len(drive.stack) > 0 {
		text += "Stack: /" + strings.Join(drive.stack, "/") + "\n"
	}
	if len(drive.path) > 0 {
		text += "Path: /" + strings.Join(drive.path, "/") + "\n"
	}
	if len(drive.method) > 0 {
		text += "Method: " + drive.method + "\n"
	}
	return text
}

func (drive *DataDrive) Top() string {
	part := ""
	if len(drive.path) > 0 {
		part = drive.path[0]
	}
	return part
}

func (drive *DataDrive) Shift() string {
	part := ""
	if len(drive.path) > 0 {
		part = drive.path[0]
		drive.path = drive.path[1:]
		drive.stack = append(drive.stack, part)
	}
	return part
}

func (drive *DataDrive) ShiftAll() string {
	part := drive.Shift()
	for len(drive.path) > 0 {
		part += "/" + drive.Shift()
	}
	return part
}

func (drive *DataDrive) IsShift(part string) bool {
	if strings.ToLower(part) == strings.ToLower(drive.Top()) {
		drive.Shift()
		return true
	}
	return false
}

func (drive *DataDrive) ItShift(part string) string {
	if cmd := strings.ToLower(drive.Top()); strings.HasPrefix(cmd, strings.ToLower(part)) {
		drive.Shift()
		return cmd
	}
	return ""
}

func (drive *DataDrive) IsShiftInt() (int, bool) {
	var pi int = 0
	var ok bool = false
	if pi, ok = lik.StrToIntIf(drive.Top()); ok {
		drive.Shift()
	}
	return pi, ok
}

func (drive *DataDrive) IsEmpty() bool {
	return len(drive.path) == 0
}

func (drive *DataDrive) IsMethod(method string) bool {
	return strings.ToLower(drive.method) == strings.ToLower(method)
}

func (drive *DataDrive) IsAtMethod(method string) bool {
	return drive.IsEmpty() && drive.IsMethod(method)
}

func (drive *DataDrive) InitializePage(version string) likdom.Domer {
	html := likdom.BuildPageHtml()
	if item, _ := html.GetDataTag("head"); item != nil {
		item.BuildString("<script type='text/javascript' src='/js/jquery.js?"+version+"'></script>",
			"<script type='text/javascript' src='/js/lik.js?"+version+"'></script>",
			"<meta http-equiv=\"Content-Language\" content=\"ru\">",
			"<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">",
		)
	}
	if item, _ := html.GetDataTag("body"); item != nil {
		script := item.BuildItem("script")
		script.BuildString(fmt.Sprintf("lik_time=%d;", time.Now().Unix()))
		script.BuildString(fmt.Sprintf("lik_page=%d;", drive.Page.GetPageId()))
		if drive.Page.GetTrust() {
			script.BuildString("lik_trust=1;")
		}
		if width, height, _ := drive.Page.GetSizeFix(); width > 0 && height > 0 {
			script.BuildString(fmt.Sprintf("screen_width=%d;", width))
			script.BuildString(fmt.Sprintf("screen_height=%d;", height))
		}
		script.BuildString("jQuery(document).ready(function () { lik_start(); });")
	}
	return html
}

func (drive *DataDrive) StoreItem(items ...likdom.Domer) {
	for _, item := range items {
		if item != nil {
			code := item.ToString()
			if code != "" {
				id := item.GetAttr("id")
				if id == "" {
					id = "_self"
				}
				if id != "" {
					drive.SetResponse(code, id)
				}
			}
		}
	}
}

func (drive *DataDrive) SetGoPart(part string) {
	url := drive.BuildUrl(part)
	drive.SetResponse(url, "_topart")
}

func (drive *DataDrive) SetWindowPart(part string) {
	url := drive.BuildUrl(part) + "&_tp=1"
	drive.SetResponse(url, "_function_lik_window_part")
}

func (drive *DataDrive) SetOnPart(part string) {
	url := drive.BuildUrl(part)
	drive.SetResponse(url, "_url")
	drive.SetResponse(drive.Page.GetPageId(), "_sp")
}

func (drive *DataDrive) PushOnPart(part string) {
	url := drive.BuildUrl(part)
	drive.SetResponse(url, "_history")
	drive.SetResponse(drive.Page.GetPageId(), "_sp")
}

func (drive *DataDrive) SetTitle(title string) {
	drive.SetResponse(title, "_title")
}

func (drive *DataDrive) BuildUrl(part string) string {
	url := part
	if url == "" {
		url = "/"
	}
	if !strings.Contains(url, "?_sp=") && !strings.Contains(url, "&_sp=") {
		if !strings.Contains(url, "?") {
			url += "?"
		} else {
			url += "&"
		}
		url += "_sp=" + fmt.Sprint(int(drive.Page.GetPageId()))
	}
	return url
}

func (drive *DataDrive) SeekPageSize() bool {
	ok := true
	width, height := drive.Page.GetSize()
	sx, sy := width, height
	if val := lik.StrToInt(drive.GetContext("_sw")); val > 0 {
		sx = val
	}
	if val := lik.StrToInt(drive.GetContext("_sh")); val > 0 {
		sy = val
	}
	if sx != width || sy != height {
		ok = false
		drive.Page.SetSize(sx, sy, false)
	}
	return ok
}

func (drive *DataDrive) GetContext(key string) string {
	val := ""
	drive.mutex.Lock()
	if drive.context != nil {
		val = drive.context.GetString(key)
	}
	drive.mutex.Unlock()
	return val
}

func (drive *DataDrive) SetContext(val interface{}, key string) {
	drive.mutex.Lock()
	if drive.context == nil {
		drive.context = lik.BuildSet()
	}
	drive.context.SetItem(val, key)
	drive.mutex.Unlock()
}

func (drive *DataDrive) GetAllContext() lik.Seter {
	return drive.context
}

func (drive *DataDrive) SetCookie(val string, key string) {
	drive.mutex.Lock()
	if drive.cookies == nil {
		drive.cookies = lik.BuildSet()
	}
	drive.cookies.SetItem(val, key)
	drive.mutex.Unlock()
}

func (drive *DataDrive) GetCookie(key string) string {
	val := ""
	if drive.cookies != nil {
		val = drive.cookies.GetString(key)
	}
	return val
}

func (drive *DataDrive) SetResponse(val interface{}, key string) {
	drive.mutex.Lock()
	if drive.responce == nil {
		drive.responce = lik.BuildSet()
	}
	drive.responce.SetItem(val, key)
	drive.mutex.Unlock()
}

func (drive *DataDrive) GetAllResponse() lik.Seter {
	return drive.responce
}

func (drive *DataDrive) GetAllCookies() lik.Seter {
	return drive.cookies
}
