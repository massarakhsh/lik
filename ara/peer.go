package ara

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/massarakhsh/lik"
)

type Peer struct {
	httpClient *http.Client
}

func newPeer(httpAddr string) *Peer {
	p := &Peer{}

	httpTransport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	if strings.Contains(httpAddr, "https") {
		caCert, err := os.ReadFile("")
		if err != nil {
			log.Println(err)
			return nil
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
		}
		httpTransport.TLSClientConfig = tlsConfig
	}

	p.httpClient = &http.Client{Transport: httpTransport}

	return p
}

func (p *Peer) HttpRequest(uri string, token string, jsonData []byte) lik.Seter {
	method, reqBody := "GET", bytes.NewBuffer(nil)
	if len(jsonData) > 0 {
		method = "POST"
		reqBody = bytes.NewBuffer(jsonData)
	}

	request, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		return nil
	}

	if token != "" {
		request.Header.Add("Authentication", token)
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := p.httpClient.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	var result lik.Seter
	if body, err := io.ReadAll(response.Body); err == nil {
		result = lik.SetFromRequest(string(body))
	}
	return result
}
