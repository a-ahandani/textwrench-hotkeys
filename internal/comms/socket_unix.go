//go:build darwin || linux

package comms

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
)

type socketCommunicator struct {
	listener net.Listener
	conn     net.Conn
	mu       sync.Mutex
	path     string
}

func newPlatformSpecificCommunicator() Communicator {
	path := filepath.Join(os.TempDir(), "textwrench.sock")
	return &socketCommunicator{path: path}
}

func (s *socketCommunicator) Start(ctx context.Context, handler MessageHandler) error {
	_ = os.Remove(s.path) // Remove stale socket file if it exists

	ln, err := net.Listen("unix", s.path)
	if err != nil {
		return fmt.Errorf("failed to start unix socket: %w", err)
	}
	s.listener = ln

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
				select {
				case <-ctx.Done():
					return
				default:
				}
				continue
			}

			// Close any old connection before overwriting
			s.mu.Lock()
			if s.conn != nil {
				_ = s.conn.Close()
			}
			s.conn = conn
			s.mu.Unlock()

			// Read incoming messages
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				handler(scanner.Text())
			}

			// Clean up this connection
			s.mu.Lock()
			_ = conn.Close()
			if s.conn == conn {
				s.conn = nil
			}
			s.mu.Unlock()
		}
	}()

	return nil
}

func (s *socketCommunicator) Send(message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return fmt.Errorf("no connection established")
	}
	_, err := fmt.Fprintln(s.conn, message)
	return err
}

func (s *socketCommunicator) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	_ = os.Remove(s.path)
	return nil
}
