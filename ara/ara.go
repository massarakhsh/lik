package ara

import (
	"fmt"
	"time"

	"go.uber.org/atomic"
)

// Дескриптор  соединения с БД ARA
type ItAra struct {
	peers map[string]*Peer // Список пиров

	part       string    //	Раздел базы данных сервиса
	seconds    int       //	Keep-live сервиса в секундах
	token      string    //	Токен для регистрации в ARA
	key        string    //	Ключ экземпляра сервиса
	lastTry    time.Time //	Время последней попытки регистрации
	lastOnline time.Time //	Время последней удачной регистрации

	isStarted *atomic.Bool // Запущена автоматическая перерегистрация
	isStopped *atomic.Bool // Останавливается автоматическая перерегистрация
}

// Построить дескриптор сервиса ARA
// address - список адресов пиров
// part - раздел базы данных
// seconds - keep-live сервиса в секундах
// token - токен для регистрациив ARA
func Build(address []string, part string, seconds int, token string) *ItAra {
	it := &ItAra{
		peers:     map[string]*Peer{},
		part:      part,
		seconds:   seconds,
		token:     token,
		isStarted: atomic.NewBool(false),
		isStopped: atomic.NewBool(false),
	}

	for _, addr := range address {
		peer := newPeer(addr)
		it.peers[addr] = peer
	}

	return it
}

// Запуск автоматической перерегистрации
func (it *ItAra) StartAutoRegister(seconds int, jsonData []byte) bool {
	if !it.isStarted.Load() {
		it.isStarted.Store(true)

		go func() {
			for !it.isStopped.Load() {
				if time.Since(it.lastTry) >= time.Second*time.Duration(seconds) {
					it.Register(jsonData)
				}
				time.Sleep(time.Millisecond * 100)
			}
			it.isStarted.Store(false)
		}()

		return true
	}

	return false
}

// Останов автоматической перерегистрации
func (it *ItAra) StopAutoRegister() bool {
	if it.isStarted.Load() {
		it.isStopped.Store(true)
		return true
	}
	return false
}

// Регистрация сервиса в ARA
// data - конфигурация сервиса (разгруженная структура JSON)
func (it *ItAra) Register(jsonData []byte) bool {
	for httpAddr, peer := range it.peers {
		uri := fmt.Sprintf("%s/register?part=%s&duration=%d", httpAddr, it.part, it.seconds)

		if it.key != "" {
			uri += "&key=" + it.key
		}

		it.lastTry = time.Now()

		if regs := peer.HttpRequest(uri, it.token, jsonData); regs != nil {
			if key := regs.GetString("key"); key != "" {
				it.key = key
				it.lastOnline = time.Now()
				return true
			}
		}
	}

	return false
}
