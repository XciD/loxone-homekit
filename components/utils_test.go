package components

import (
	"net"
	"testing"
	"time"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
)

type LoxoneFake struct {
	hooks    *map[string]func(*events.Event)
	commands *[]string
}

func (l *LoxoneFake) AddHook(uuid string, callback func(*events.Event)) {
	i := *l.hooks
	i[uuid] = callback
}
func (l *LoxoneFake) SendCommand(command string, class interface{}) (*loxone.Body, error) {
	*l.commands = append(*l.commands, command)
	return &loxone.Body{Code: 200}, nil
}

func (l *LoxoneFake) Close() {

}

func (l *LoxoneFake) RegisterEvents() error {
	return nil
}

func (l *LoxoneFake) PumpEvents(stop <-chan bool) {

}

func (l *LoxoneFake) GetConfig() (*loxone.Config, error) {
	return nil, nil
}

type fixture struct {
	*testing.T
	*ComponentConfig
	*loxone.Control
	loxone.LoxoneInterface
	hooks    map[string]func(*events.Event)
	commands []string
}

func (l *fixture) TriggerEvent(uuid string, value float64) {
	if hook, ok := l.hooks[uuid]; ok {
		hook(&events.Event{Value: value})
	}
}

func (l fixture) GetCommands() []string {
	return l.commands
}
func NewFixture(name, id, loxoneType string, states map[string]interface{}) *fixture {
	fixture := &fixture{}

	fixture.ComponentConfig = &ComponentConfig{
		ID:         id,
		Name:       name,
		Type:       1,
		LoxoneType: loxoneType,
	}

	fixture.Control = &loxone.Control{
		Name:       name,
		Type:       loxoneType,
		UUIDAction: id,
		States:     states,
	}

	fixture.commands = make([]string, 0)
	fixture.hooks = make(map[string]func(*events.Event))
	fixture.LoxoneInterface = &LoxoneFake{
		hooks:    &fixture.hooks,
		commands: &fixture.commands,
	}

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
