package components

import (
	"fmt"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneJalousie struct {
	*Component
	*service.WindowCovering
	up   bool
	down bool
}

func NewJalousie(config ComponentConfig, control *loxone.Control, lox loxone.WebsocketInterface) *Component {
	component := &LoxoneJalousie{
		Component: newComponent(config, control, lox),
	}

	component.WindowCovering = service.NewWindowCovering()
	component.AddService(component.Service)

	component.up = false
	component.down = false

	component.TargetPosition.OnValueRemoteUpdate(component.onTargetPositionStateRemoteUpdate)

	component.PositionState.OnValueGet(component.onPositionStateGet)

	component.addHook("up", component.upHook)
	component.addHook("down", component.downHook)
	component.addHook("position", component.positionHook)

	return component.Component
}

func (l *LoxoneJalousie) getState() int {
	switch {
	case l.up:
		return characteristic.PositionStateIncreasing
	case l.down:
		return characteristic.PositionStateDecreasing
	default:
		return characteristic.PositionStateStopped
	}
}

func (l *LoxoneJalousie) setState(state int) {
	l.PositionState.SetValue(state)
}

func (l *LoxoneJalousie) setPosition(position int) {
	pos := int32(100 - position)

	l.Logger.Infof("Asking for loxone position %d", pos)

	command := fmt.Sprintf("ManualPosition/%d", pos)
	if pos == 0 {
		command = "FullUp"
	} else if pos == 100 {
		command = "FullDown"
	}

	l.command(command)
}

func (l *LoxoneJalousie) onTargetPositionStateRemoteUpdate(target int) {
	l.Logger.Infof("Remote target %d, current position %d, current target %d, current state %d",
		target, l.PositionState.Value, l.TargetPosition.Value, l.getState())

	if l.getState() != characteristic.PositionStateStopped {
		l.command("stop")
		l.setState(characteristic.PositionStateStopped)
		l.TargetPosition.SetValue(l.PositionState.Value.(int))
		return
	}
	l.setPosition(target)
}

func (l *LoxoneJalousie) onPositionStateGet() interface{} {
	return l.getState()
}

func (l *LoxoneJalousie) upHook(event *events.Event) {
	l.up = event.Value == 1
	l.Logger.Infof("Updating state up to %t", l.up)
	l.setState(l.getState())
}

func (l *LoxoneJalousie) downHook(event *events.Event) {
	l.down = event.Value == 1
	l.Logger.Infof("Updating state down to %t", l.down)
	l.setState(l.getState())
}

func (l *LoxoneJalousie) positionHook(event *events.Event) {
	l.CurrentPosition.SetValue(int(100 - (event.Value * 100)))

	state := l.getState()

	l.setState(state)

	if state == characteristic.PositionStateStopped {
		l.TargetPosition.SetValue(l.CurrentPosition.Value.(int))
	}

	l.Logger.Infof("Websocket event: %d, Target position %d, state %d",
		l.CurrentPosition.Value, l.TargetPosition.Value, l.getState())
}
