package components

import (
	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneGate struct {
	*Component
	*service.GarageDoorOpener
	*characteristic.CurrentPosition
	active int
}

func NewGate(config ComponentConfig, control *loxone.Control, lox loxone.WebsocketInterface) *Component {
	component := &LoxoneGate{
		Component: newComponent(config, control, lox),
	}

	component.GarageDoorOpener = service.NewGarageDoorOpener()
	component.AddService(component.Service)

	component.CurrentPosition = characteristic.NewCurrentPosition()
	component.AddCharacteristic(component.CurrentPosition.Characteristic)

	component.active = 0
	// Not initialized value
	component.TargetDoorState.Value = -1

	component.CurrentDoorState.OnValueGet(component.onCurrentDoorStateGet)

	component.TargetDoorState.OnValueRemoteUpdate(component.onTargetDoorStateRemoteUpdate)

	component.addHook("position", component.positionHook)
	component.addHook("active", component.activeHook)

	component.addDebugHook("preventOpen")
	component.addDebugHook("preventClose")

	return component.Component
}

func (l *LoxoneGate) getState() int {
	switch l.active {
	case -1:
		l.Logger.Info("Door is closing")
		return characteristic.CurrentDoorStateClosing
	case 1:
		l.Logger.Info("Door is Opening")
		return characteristic.CurrentDoorStateOpening
	default:
		switch l.CurrentPosition.Value {
		case 100:
			l.Logger.Info("Door is Closed")
			return characteristic.CurrentDoorStateClosed
		case 0:
			l.Logger.Info("Door is Open")
			return characteristic.CurrentDoorStateOpen
		default:
			l.Logger.Info("Door is Stopped")
			return characteristic.CurrentDoorStateStopped
		}
	}
}

func (l *LoxoneGate) setState(state int) {
	l.CurrentDoorState.SetValue(state)
}

func (l *LoxoneGate) onCurrentDoorStateGet() interface{} {
	return l.getState()
}

func (l *LoxoneGate) onTargetDoorStateRemoteUpdate(i int) {
	if l.active != 0 {
		l.Logger.Info("Door already moving, Stopping")
		l.command("stop")
		return
	}
	if i == characteristic.TargetDoorStateOpen {
		l.Logger.Info("Asking to open the door")
		l.command("open")
	} else {
		l.Logger.Info("Asking to close the door")
		l.command("close")
	}
}
func (l *LoxoneGate) positionHook(event *events.Event) {
	position := int(100 - (event.Value * 100))
	l.CurrentPosition.SetValue(position)

	state := l.getState()
	l.checkTarget(state)
	l.setState(state)
}

func (l *LoxoneGate) activeHook(event *events.Event) {
	l.active = int(event.Value)
	state := l.getState()
	l.checkTarget(state)
	l.setState(state)
}

func (l *LoxoneGate) checkTarget(state int) {
	target := l.TargetDoorState.Value

	if (state == characteristic.CurrentDoorStateOpening || state == characteristic.CurrentDoorStateOpen) && target != characteristic.TargetDoorStateOpen {
		l.TargetDoorState.SetValue(characteristic.TargetDoorStateOpen)
	}

	if (state == characteristic.CurrentDoorStateClosing || state == characteristic.CurrentDoorStateClosed) && target != characteristic.TargetDoorStateClosed {
		l.TargetDoorState.SetValue(characteristic.TargetDoorStateClosed)
	}
}
