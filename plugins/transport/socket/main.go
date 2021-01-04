package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/infrawatch/sg-core-refactor/pkg/config"
	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

const maxBufferSize = 4096

var msgBuffer []byte

func init() {
	msgBuffer = make([]byte, maxBufferSize)
}

type configuration struct {
	Address string `validate:"required"`
}

//Socket basic struct
type Socket struct {
	conf configuration
}

//Run implements type Transport
func (s *Socket) Run(ctx context.Context, wg *sync.WaitGroup, w transport.WriteFn) {
	defer wg.Done()

	// var laddr net.UnixAddr

	// laddr.Name = s.conf.Address
	// laddr.Net = "unixgram"

	// os.Remove(s.conf.Address)

	// pc, err := net.ListenUnixgram("unixgram", &laddr)
	// if err != nil {
	// 	return err
	// }
	// defer os.Remove(s.conf.Address)
	// defer pc.Close()

	// for {
	// 	n, err := pc.Read(msgBuffer[:])
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if n < 1 {
	// 		return nil
	// 	}
	// 	t <- msgBuffer
	// }

	for {
		select {
		case <-ctx.Done():
			goto done
		default:
			time.Sleep(time.Second * 1)
			ret := fmt.Sprintf("message from socket")
			w([]byte(ret))
		}
	}

done:
}

//Config load configurations
func (s *Socket) Config(c []byte) error {
	s.conf = configuration{}
	err := config.ParseConfig(bytes.NewReader(c), &s.conf)
	if err != nil {
		return err
	}
	return nil
}

//New create new socket transport
func New() transport.Transport {
	return &Socket{}
}
