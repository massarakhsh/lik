package ara

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/massarakhsh/lik"
)

type ItAra struct {
	Addresses  []string
	Part       string
	Seconds    int
	Token      string
	Key        string
	LastTry    time.Time
	LastOnline time.Time
}

func Build(address []string, part string, seconds int, token string) *ItAra {
	it := &ItAra{Addresses: address, Part: part, Seconds: seconds, Token: token}
	return it
}

func (it *ItAra) Connect(data string) bool {
	for _, peer := range it.Addresses {
		uri := fmt.Sprintf("%s/register?part=%s&duration=%d", peer, it.Part, it.Seconds)
		if it.Key != "" {
			uri += "&key=" + it.Key
		}
		it.LastTry = time.Now()
		if regs := it.callPost(uri, data); regs != nil {
			if key := regs.GetString("key"); key != "" {
				it.Key = key
				it.LastOnline = time.Now()
				return true
			}
		}
	}
	return false
}

func (it *ItAra) callPost(uri string, data string) lik.Seter {
	var transport *http.Transport
	if strings.HasPrefix(uri, "https") {
		var tlsConfig *tls.Config
		caCert, err := os.ReadFile("") //core.ParmFilePem)
		if err != nil {
			return nil
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
		}
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}
	} else {
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}
	}
	client := &http.Client{Transport: transport}

	method, body := "GET", io.Reader(nil)
	if data != "" {
		method = "POST"
		body = bytes.NewBuffer([]byte(data))
	}
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil
	}

	if it.Token != "" {
		req.Header.Add("Authentication", it.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	/*if data == nil {
		resp, err = http.Get(uri)
	} else {
		databytes := []byte(data.Serialize())
		CounterOutBytes.Add(float64(len(databytes)))
		resp, err = http.Post(uri, "application/json", bytes.NewBuffer(databytes))
	}*/
	var result lik.Seter
	defer resp.Body.Close()
	if body, err := io.ReadAll(resp.Body); err == nil {
		result = lik.SetFromRequest(string(body))
	}
	return result
}
