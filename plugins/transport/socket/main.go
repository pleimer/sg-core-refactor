package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/pkg/config"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

const maxBufferSize = 4096

type configT struct {
	Address string `validate:"required"`
}

//Socket basic struct
type Socket struct {
	conf   configT
	logger *logging.Logger
}

//Run implements type Transport
func (s *Socket) Run(ctx context.Context, wg *sync.WaitGroup, w transport.WriteFn) {
	defer wg.Done()

	msgBuffer := make([]byte, maxBufferSize)
	var laddr net.UnixAddr

	laddr.Name = s.conf.Address
	laddr.Net = "unixgram"

	os.Remove(s.conf.Address)

	pc, err := net.ListenUnixgram("unixgram", &laddr)
	if err != nil {
		s.logger.Metadata(logging.Metadata{"plugin": "socket", "error": err})
		s.logger.Error("failed to listen on unix soc")
		return
	}
	defer os.Remove(s.conf.Address)
	defer pc.Close()

	go func() {
		for {
			n, err := pc.Read(msgBuffer)

			if err != nil {
				s.logger.Metadata(logging.Metadata{"plugin": "socket", "error": err})
				s.logger.Error("failed reading data")
				return //done
			}
			if n < 1 {
				return
			}
			w(msgBuffer[:n])
		}
	}()

	<-ctx.Done()
}

//Listen ...
func (s *Socket) Listen(e data.Event) {
	fmt.Printf("Recieved event: %v\n", e)
}

//Config load configurations
func (s *Socket) Config(c []byte) error {
	s.conf = configT{}
	err := config.ParseConfig(bytes.NewReader(c), &s.conf)
	if err != nil {
		return err
	}
	return nil
}

//New create new socket transport
func New(l *logging.Logger) transport.Transport {
	return &Socket{
		logger: l,
	}
}
