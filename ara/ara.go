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

// Дескриптор  соединения с БД ARA
type ItAra struct {
	Addresses  []string  //	Список url - адресов пиров
	Part       string    //	Раздел базы данных сервиса
	Seconds    int       //	Keep-live сервиса в секундах
	Token      string    //	Токен для регистрации в ARA
	Key        string    //	Ключ экземпляра сервиса
	LastTry    time.Time //	Время последней попытки регистрации
	LastOnline time.Time //	Время последней удачной регистрации
	isStarted  bool      // Запущена автоматическая перерегистрация
	isStoping  bool      // Останавливается автоматическая перерегистрация
}

// Построить дескриптор сервиса ARA
// address - список адресов пиров
// part - раздел базы данных
// seconds - keep-live сервиса в секундах
// token - токен для регистрациив ARA
func Build(address []string, part string, seconds int, token string) *ItAra {
	it := &ItAra{Addresses: address, Part: part, Seconds: seconds, Token: token}
	return it
}

// Запуск автоматической перерегистрации
func (it *ItAra) StartAutoRegister(seconds int, data string) bool {
	if !it.isStarted {
		it.isStarted = true
		go func() {
			for !it.isStoping {
				if time.Since(it.LastTry) >= time.Second * time.Duration(seconds) {
					it.Register(data)
				}
				time.Sleep(time.Millisecond * 100)
			}
			it.isStarted = false
		}()
		return true
	}
	return false
}

// Останов автоматической перерегистрации
func (it *ItAra) StopAutoRegister() bool {
	if it.isStarted {
		it.isStoping = true
		return true
	}
	return false
}

// Регистрация сервиса в ARA
// data - конфигурация сервиса (разгруженная структура JSON)
func (it *ItAra) Register(data string) bool {
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
