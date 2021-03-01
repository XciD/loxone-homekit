package components

import (
	"net"
	"testing"
	"time"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/test"
)

type fixture struct {
	*testing.T
	*ComponentConfig
	*loxone.Control
	*test.FakeWebsocket
	*Factory
}

func NewFixture(name, id, loxoneType string, states map[string]interface{}) *fixture {
	fixture := &fixture{}

	fixture.ComponentConfig = &ComponentConfig{
		ID:          id,
		Name:        name,
		ControlType: loxoneType,
		Control: &loxone.Control{
			Name:       name,
			Type:       loxoneType,
			UUIDAction: id,
			States:     states,
		},
	}

	websocket := test.NewFakeWebsocket()
	fixture.FakeWebsocket = websocket.(*test.FakeWebsocket)
	fixture.Factory = &Factory{Loxone: fixture.FakeWebsocket}

	return fixture
}

var TestConn net.Conn = &fakeConn{}

type fakeConn struct {
}

func (f *fakeConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (f *fakeConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (f *fakeConn) Close() error {
	return nil
}

func (f *fakeConn) LocalAddr() net.Addr {
	return nil
}

func (f *fakeConn) RemoteAddr() net.Addr {
	return nil
}

func (f *fakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (f *fakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (f *fakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
