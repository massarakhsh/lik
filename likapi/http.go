package likapi

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/massarakhsh/lik"
)

type EnterHttp struct {
	Port  int
	Tls   bool
	Enter func(drive *DataDrive)
}

var (
	EnterList []EnterHttp = []EnterHttp{}
)

func RegisterHttp(port int, tls bool, enter func(*DataDrive)) {
	entelm := EnterHttp{
		Port:  port,
		Tls:   tls,
		Enter: enter,
	}
	EnterList = append(EnterList, entelm)
}

func GetParm(r *http.Request, key string) string {
	val := ""
	if match := lik.RegExParse(r.RequestURI, "[?&]"+key+"=([^&]+)"); match != nil {
		val = match[1]
	}
	return val
}

func GetHeader(r *http.Request, key string) string {
	val := ""
	if vals, ok := r.Header[key]; ok && len(vals) > 0 {
		val = vals[len(vals)-1]
	}
	return val
}

func ProbeRouteFile(w http.ResponseWriter, r *http.Request, path string) bool {
	if match := lik.RegExParse(path, "^/?([^?]*)"); match != nil {
		path = match[1]
	}
	return RouteFile(w, r, path)
}

func RouteFile(w http.ResponseWriter, r *http.Request, name string) bool {
	ok := false
	if _, err := os.Stat(name); err == nil {
		w.Header().Set("Cache-control", "private,no-cache,no-store,must-revalidate")
		http.ServeFile(w, r, name)
		ok = true
	} else if match := lik.RegExParse(name, "^js/(.*)$"); match != nil {
		w.Header().Set("Cache-control", "private,no-cache,no-store,must-revalidate")
		http.ServeFile(w, r, "../lik/"+match[1])
		ok = true
	}
	return ok
}

func RouteHtml(w http.ResponseWriter, rc int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	w.WriteHeader(rc)
	fmt.Fprint(w, html)
}

func Route401(w http.ResponseWriter, rc int, msg string) {
	w.Header().Set("Www-Authenticate", "Basic "+msg)
	w.WriteHeader(rc)
}

func RouteRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(302)
}

func RouteRast(w http.ResponseWriter, rc int, image *image.RGBA) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	w.WriteHeader(rc)
	if image != nil {
		png.Encode(w, image)
	}
}

func RouteCsv(w http.ResponseWriter, rc int, content lik.Lister, dlm string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	w.WriteHeader(rc)
	if content == nil {
		fmt.Fprint(w, "")
	} else {
		fmt.Fprint(w, content.ToCsv(dlm))
	}
}

func RouteJson(w http.ResponseWriter, rc int, content lik.Seter, format bool) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(rc)
	if content == nil {
		fmt.Fprint(w, "{}")
	} else if format {
		serial := content.Format("")
		//fmt.Printf("JSON format length: %d\n", len(serial))
		fmt.Fprint(w, serial)
	} else {
		serial := content.Serialize()
		//fmt.Printf("JSON length: %d\n", len(serial))
		fmt.Fprint(w, serial)
	}
}

func RouteXml(w http.ResponseWriter, rc int, content lik.Lister) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	w.WriteHeader(rc)
	if content == nil {
		fmt.Fprint(w, "[]")
	} else {
		fmt.Fprint(w, content.ToXml())
	}
}

func RouteCookies(w http.ResponseWriter, cookies lik.Seter) {
	if cookies != nil {
		for _, set := range cookies.Values() {
			if val := set.Val.ToString(); val != "" {
				key := set.Key
				expiration := time.Now().Add(30 * 24 * time.Hour)
				cookie := http.Cookie{Path: "/", Name: key, Value: val, SameSite: http.SameSiteLaxMode, Expires: expiration}
				http.SetCookie(w, &cookie)
			}
		}
	}
}

func GetApiJson(url string, sets lik.Seter) (int, lik.Seter) {
	var rc int
	var info lik.Seter
	uri := buildUri(url, sets)
	if resp, err := http.Get(uri); err != nil {
		fmt.Print(err)
	} else {
		rc = resp.StatusCode
		if body, err := ioutil.ReadAll(resp.Body); err == nil && body != nil {
			info = lik.SetFromRequest(string(body))
		}
	}
	<-time.After(time.Millisecond * 1000)
	//fmt.Println("Try: ", trgo + 1)
	return rc, info
}

func buildUri(requrl string, sets lik.Seter) string {
	requri := requrl
	if sets != nil {
		for _, set := range sets.Values() {
			if key := set.Key; key != "" {
				if strings.Index(requri, "?") >= 0 {
					requri += "&"
				} else {
					requri += "?"
				}
				requri += key + "="
				if set.Val != nil {
					str := set.Val.ToString()
					requri += url.QueryEscape(str)
				}
			}
		}
	}
	return requri
}

func SetToPart(drive DataDriver, to DataDriver, part string) int {
	id := drive.GetPage().ContinueToPage(to.GetPage())
	url := BuildUrl(id, part)
	drive.SetResponse(url, "_topart")
	return id
}

func BuildUrl(id int, part string) string {
	url := part
	if url == "" {
		url = "/"
	}
	if strings.Index(url, "?_sp=") < 0 && strings.Index(url, "&_sp=") < 0 {
		if strings.Index(url, "?") < 0 {
			url += "?"
		} else {
			url += "&"
		}
		url += "_sp=" + fmt.Sprint(id)
	}
	return url
}
