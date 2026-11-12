package canceler

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/bufio"
	M "github.com/sagernet/sing/common/metadata"
)

// A packet written while a read waits counts as activity, and the session ends one timeout after
// it, not one timeout after the read's deadline: re-arming a whole timeout at that deadline let
// an idle session live up to twice its timeout.
func TestH_TimeoutPacketConnEndsOneTimeoutAfterTheLastWrite(t *testing.T) {
	const timeout = 400 * time.Millisecond
	server, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	peer, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	ctx, conn := NewTimeoutPacketConn(context.Background(), bufio.NewPacketConn(server), timeout)
	defer conn.Close()

	readDone := make(chan time.Time, 1)
	go func() {
		for {
			buffer := buf.NewPacket()
			_, err := conn.ReadPacket(buffer)
			buffer.Release()
			if err != nil {
				readDone <- time.Now()
				return
			}
		}
	}()

	// the read is waiting (its deadline one timeout from now); a reply goes out halfway
	time.Sleep(timeout / 2)
	reply := buf.NewPacket()
	reply.WriteString("reply")
	lastActive := time.Now()
	if err := conn.WritePacket(reply, M.SocksaddrFromNet(peer.LocalAddr())); err != nil {
		t.Fatal(err)
	}

	select {
	case ended := <-readDone:
		idle := ended.Sub(lastActive)
		if idle < timeout {
			t.Fatalf("ended %v after the last write, before the %v timeout", idle, timeout)
		}
		if idle > timeout+timeout/2 {
			t.Fatalf("ended %v after the last write, a %v timeout re-armed in full", idle, timeout)
		}
	case <-time.After(3 * timeout):
		t.Fatal("the idle session never ended")
	}
	if context.Cause(ctx) == nil {
		t.Fatal("the context was not canceled")
	}
}
