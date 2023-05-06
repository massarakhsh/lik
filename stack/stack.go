package stack

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
)

type ItStack struct {
	Method         string
	Request        string
	IP             string
	Host           string
	Port           int
	Path           []string
	Parms          lik.Seter
	Info           lik.Seter
	InPath         string
	IdSession      int
	Authentication string

	Dom        likdom.Domer
	Json       lik.Seter
	CookiesOut map[string]string
}

func BuildRequest(r *http.Request) *ItStack {
	it := &ItStack{}
	it.Parms = lik.BuildSet()
	it.loadHost(r)
	it.loadPath(r.RequestURI)
	it.loadCookies(r)
	it.loadAuth(r)
	it.loadInfo(r)
	return it
}

func BuildMethodPath(method, path string) *ItStack {
	it := &ItStack{}
	it.Parms = lik.BuildSet()
	it.Method = method
	it.loadPath(path)
	return it
}

func (it *ItStack) loadHost(r *http.Request) {
	it.Method = r.Method
	it.IP = r.RemoteAddr
	if match := lik.RegExParse(r.Host, "^(.+):(\\d+)$"); match != nil {
		it.Host = match[1]
		it.Port = lik.StrToInt(match[2])
	} else {
		it.Host = r.Host
		it.Port = 80
	}
}

func (it *ItStack) loadPath(path string) {
	it.Request = path
	var parts []string
	if pos := strings.Index(path, "?"); pos >= 0 {
		parts = lik.PathToNames(path[:pos])
		it.loadContext(path[pos+1:])
	} else {
		parts = lik.PathToNames(it.Request)
	}
	for _, name := range parts {
		if match := lik.RegExParse(name, "^([^=]+)=(.*)"); match != nil {
			it.Parms.SetValue(match[1], match[2])
		} else if name != "" {
			it.Path = append(it.Path, name)
		}
	}
	it.InPath = "/" + strings.Join(it.Path, "/")
}

func (it *ItStack) AddCookie(name string, value string) {
	it.CookiesOut[name] = value
}

func (it *ItStack) loadCookies(r *http.Request) {
	it.CookiesOut = make(map[string]string)
	for _, cookie := range r.Cookies() {
		if cookie.Name == "token" {
			it.Authentication = cookie.Value
		} else if it.IdSession == 0 && cookie.Name == "id" {
			it.IdSession = lik.StrToInt(cookie.Value)
		}
	}
}

func (it *ItStack) loadAuth(r *http.Request) {
	if auth := r.Header.Get("Authentication"); auth != "" {
		it.Authentication = auth
	} else if auth := it.Parms.GetString("token"); auth != "" {
		it.Authentication = auth
	}
}

func (it *ItStack) loadInfo(r *http.Request) {
	if body, err := io.ReadAll(r.Body); err == nil && body != nil {
		//		metrics.CounterInBytes.Add(float64(len(body)))
		it.Info = lik.SetFromRequest(string(body))
	}
}

func (it *ItStack) Top() string {
	if len(it.Path) == 0 {
		return ""
	}
	return it.Path[0]
}

func (it *ItStack) Pop() string {
	if len(it.Path) == 0 {
		return ""
	}
	cmd := it.Path[0]
	it.Path = it.Path[1:]
	return cmd
}

func (it *ItStack) Probe(cmd string) bool {
	if len(it.Path) == 0 {
		return false
	}
	if it.Path[0] != cmd {
		return false
	}
	it.Path = it.Path[1:]
	return true
}

func (it *ItStack) loadContext(data string) {
	if lik.RegExCompare(data, "^{") {
		if set := lik.SetFromRequest(data); set != nil {
			for _, elm := range set.Values() {
				it.Parms.SetValue(elm.Key, elm.Val)
			}
		}
	} else if len(data) > 0 {
		data = strings.Replace(data, "+", "#", -1)
		data, _ := url.QueryUnescape(data)
		data = strings.Replace(data, "#", "+", -1)
		words := strings.Split(data, "&")
		for _, word := range words {
			if peq := strings.Index(word, "="); peq > 0 {
				id := word[:peq]
				val := word[peq+1:]
				it.Parms.SetValue(id, val)
			}
		}
	}
}
