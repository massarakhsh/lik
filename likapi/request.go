package likapi

import (
	"bytes"
	"github.com/massarakhsh/lik"
	"io"
	"io/ioutil"
	"net/http"
)

func GetHttpRequest(url string, headers lik.Seter) lik.Seter {
	return SendHttpRequest("GET", url, headers, nil)
}

func PostHttpRequest(url string, headers lik.Seter, data lik.Seter) lik.Seter {
	return SendHttpRequest("POST", url, headers, data)
}

func SendHttpRequest(method string, url string, headers lik.Seter, data lik.Seter) lik.Seter {
	answer := lik.BuildSet()
	client := http.Client{}
	var body io.Reader
	if data != nil {
		bin := []byte(data.ToString())
		body = bytes.NewReader(bin)
	}
	if request, err := http.NewRequest(method, url, body); err == nil {
		if headers != nil {
			for _, set := range headers.Values() {
				request.Header.Add(set.Key, set.Val.ToString())
			}
		}
		if resp, err := client.Do(request); err == nil {
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				str := string(body)
				answer = lik.SetFromRequest(str)
			}
		}
	}
	return answer
}
