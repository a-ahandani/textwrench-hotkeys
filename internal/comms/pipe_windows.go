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
			// Exit if context is cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			conn, err := ln.Accept()
			if err != nil {
				// If listener was closed via Close(), stop looping
				select {
				case <-ctx.Done():
					return
				default:
				}
				continue
			}

			// Close any old connection before overwriting
			p.mu.Lock()
			if p.conn != nil {
				_ = p.conn.Close()
			}
			p.conn = conn
			p.mu.Unlock()

			// Read messages until EOF or error
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				handler(scanner.Text())
			}

			// Clean up this connection
			p.mu.Lock()
			_ = conn.Close()
			if p.conn == conn {
				p.conn = nil
			}
			p.mu.Unlock()
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
		_ = p.conn.Close()
		p.conn = nil
	}
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}
