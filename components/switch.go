package components

import (
	"github.com/XciD/loxone-ws/events"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneSwitch struct {
	*Component
	*service.Service
	*characteristic.On
}

func NewLoxoneSwitch(config ComponentConfig, control *loxone.Control, lox loxone.WebsocketInterface) *Component {
	component := &LoxoneSwitch{
		Component: newComponent(config, control, lox),
	}

	component.Service = service.New(service.TypeLightbulb)
	component.AddService(component.Service)

	component.On = characteristic.NewOn()
	component.AddCharacteristic(component.On.Characteristic)

	component.On.OnValueRemoteUpdate(component.remoteUpdate)

	// Add status updates
	component.addHook("active", component.activeHook)
	return component.Component
}

func (l *LoxoneSwitch) activeHook(event *events.Event) {
	l.On.SetValue(event.Value == 1)
}

func (l *LoxoneSwitch) remoteUpdate(on bool) {
	if on {
		l.command("on")
	} else {
		l.command("off")
	}
}
