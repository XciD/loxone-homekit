package components

import (
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneAlarm struct {
	*Component
	*service.SecuritySystem
	active bool
}

func NewAlarm(f *Factory, config ComponentConfig) []*Component {
	component := &LoxoneAlarm{
		Component: f.newComponent(config, accessory.TypeSecuritySystem),
	}

	component.SecuritySystem = service.NewSecuritySystem()
	component.AddService(component.Service)

	component.addHook("armed", component.armedHook)
	component.addHook("events", component.intrusion)

	component.SecuritySystemCurrentState.Value = characteristic.SecuritySystemCurrentStateDisarmed
	component.SecuritySystemTargetState.Value = characteristic.SecuritySystemCurrentStateDisarmed

	component.SecuritySystemCurrentState.OnValueRemoteGet(component.getState)
	component.SecuritySystemTargetState.OnValueRemoteUpdate(component.onTargetSet)

	return []*Component{component.Component}
}

func (l *LoxoneAlarm) armedHook(event events.Event) {
	l.active = event.Value == 1
	l.SecuritySystemCurrentState.SetValue(l.getState())
	l.SecuritySystemTargetState.SetValue(l.getState())
}

func (l *LoxoneAlarm) intrusion(event events.Event) {
	if event.Value > 0 {
		l.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateAlarmTriggered)
	} else {
		l.SecuritySystemCurrentState.SetValue(l.getState())
	}
}

func (l *LoxoneAlarm) getState() int {
	if l.active {
		return characteristic.SecuritySystemCurrentStateStayArm
	}

	return characteristic.SecuritySystemCurrentStateDisarmed
}

func (l *LoxoneAlarm) onTargetSet(i int) {
	if i == characteristic.SecuritySystemCurrentStateDisarmed {
		l.command("off")
	} else {
		l.command("on/1")
	}
}
