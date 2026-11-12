package bufio

import (
	"io"
	"net"
	"sync"
	"testing"

	"github.com/sagernet/sing/common/buf"
)

func TestCachedConnCloseWhileReading(t *testing.T) {
	for i := 0; i < 2000; i++ {
		client, server := net.Pipe()
		buffer := buf.NewSize(64)
		buffer.WriteString("cached")
		conn := NewCachedConn(client, buffer)
		buffer.Release()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := make([]byte, 64)
			for {
				if _, err := conn.Read(p); err != nil {
					return
				}
			}
		}()
		conn.Close()
		server.Close()
		wg.Wait()
	}
}

func TestCachedConnCloseWhileWriteTo(t *testing.T) {
	for i := 0; i < 2000; i++ {
		client, server := net.Pipe()
		buffer := buf.NewSize(64)
		buffer.WriteString("cached")
		conn := NewCachedConn(client, buffer)
		buffer.Release()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn.WriteTo(io.Discard)
		}()
		conn.Close()
		server.Close()
		wg.Wait()
	}
}
