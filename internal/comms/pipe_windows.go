//go:build windows

package comms

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Microsoft/go-winio"
)

type pipeCommunicator struct {
	listener net.Listener
	conn     net.Conn
	mu       sync.Mutex
}

func newPlatformSpecificCommunicator() Communicator {
	return &pipeCommunicator{}
}

func (p *pipeCommunicator) Start(ctx context.Context, handler MessageHandler) error {
	ln, err := winio.ListenPipe(`\\.\\pipe\\textwrench-pipe`, nil)
	if err != nil {
		return fmt.Errorf("failed to start named pipe: %w", err)
	}
	p.listener = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			p.mu.Lock()
			p.conn = conn
			p.mu.Unlock()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				handler(scanner.Text())
			}
		}
	}()

	return nil
}

func (p *pipeCommunicator) Send(message string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn == nil {
		return fmt.Errorf("no connection established")
	}
	_, err := fmt.Fprintln(p.conn, message)
	return err
}

func (p *pipeCommunicator) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		p.conn.Close()
	}
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}
