package likcom

import (
	"fmt"
	"github.com/tarm/serial"
	//"os"
	"time"
)

type ComPort struct {
	NCom     int
	Baud     int
	DataOut  []byte
	DataIn   []byte
	ser      *serial.Port
	receiver ComReceiver
	stop     bool
}

type ComPorter interface {
	Close() bool
	GetNCom() int
	GetBaud() int
	SetComReceived(receiver ComReceiver)
	SendStream(raw []byte) bool
	goSend(data []byte, done chan bool)
	goReading()
}

type ComReceiver interface {
	SignReceived(data *[]byte) int
}

var (
	Debug bool
)

func GetListCom() []int {
	var list []int
	for nc := 1; nc < 99; nc++ {
		/*if file,err :=os.Open(fmt.Sprintf("COM%d",nc)); err == nil {
			list = append(list, nc)
			file.Close()
		}*/
		cnf := &serial.Config{Name: "COM" + fmt.Sprint(nc), Baud: 9600, ReadTimeout: time.Millisecond * 250}
		if ser, err := serial.OpenPort(cnf); err == nil {
			//fmt.Printf("Probe COM%d\n", nc)
			ser.Close()
			list = append(list, nc)
		}
	}
	return list
}

func OpenCom(ncom int, baud int) (ComPorter, bool) {
	cnf := &serial.Config{Name: "COM" + fmt.Sprint(ncom), Baud: baud, ReadTimeout: time.Millisecond * 250}
	ser, err := serial.OpenPort(cnf)
	if err != nil {
		return nil, false
	}
	port := &ComPort{NCom: ncom, Baud: baud, ser: ser, stop: false}
	port.DataOut = []byte{}
	port.DataIn = []byte{}
	go port.goReading()
	return port, true
}

func (port *ComPort) Close() bool {
	port.stop = true
	if port.ser == nil {
		return false
	}
	return true
}

func (port *ComPort) GetNCom() int {
	return port.NCom
}

func (port *ComPort) GetBaud() int {
	return port.Baud
}

func (port *ComPort) SetComReceived(receiver ComReceiver) {
	port.receiver = receiver
}

func (port *ComPort) SendStream(raw []byte) bool {
	if Debug {
		fmt.Print("Out:")
		for _, bt := range raw {
			fmt.Printf(" %02X", bt)
		}
		fmt.Print("\n")
	}
	res := false
	done := make(chan bool, 1)
	go port.goSend(raw, done)
	for wt := 0; wt < 300; wt++ {
		if len(done) > 0 {
			res = true
			break
		}
		<-time.After(time.Microsecond * 10)
	}
	//if _,err := port.ser.Write(raw); err != nil {
	return res
}

func (port *ComPort) goSend(data []byte, done chan bool) {
	if _, err := port.ser.Write(data); err != nil {
		done <- false
	} else {
		done <- true
	}
}

func (port *ComPort) goReading() {
	buf := make([]byte, 256, 256)
	active := false
	for !port.stop {
		nb := 0
		if port.ser != nil {
			nb, _ = port.ser.Read(buf)
		}
		if nb > 0 {
			active = true
			for b := 0; b < nb; b++ {
				bt := buf[b]
				if Debug {
					fmt.Printf(" %02X", bt)
				}
				port.DataIn = append(port.DataIn, bt)
			}
		}
		for port.receiver != nil && active && len(port.DataIn) > 0 {
			cut := port.receiver.SignReceived(&port.DataIn)
			if cut > 0 {
				port.DataIn = port.DataIn[cut:]
			} else {
				active = false
			}
		}
		if nb == 0 && !active {
			time.Sleep(time.Millisecond * 10)
		}
	}
	if port.ser != nil {
		_ = port.ser.Close()
		port.ser = nil
	}
}
