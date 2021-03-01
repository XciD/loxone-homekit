package components

import (
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc/accessory"

	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneSwitch struct {
	*Component
	*service.Service
	*characteristic.On
}

func NewLoxoneSwitch(f *Factory, config ComponentConfig) []*Component {
	if config.CategoryType != "lights" {
		// Only handles light for now.
		return nil
	}

	component := &LoxoneSwitch{
		Component: f.newComponent(config, accessory.TypeLightbulb),
	}

	component.Service = service.New(service.TypeLightbulb)
	component.AddService(component.Service)

	component.On = characteristic.NewOn()
	component.AddCharacteristic(component.On.Characteristic)

	component.On.OnValueRemoteUpdate(component.remoteUpdate)

	// Add status updates
	component.addHook("active", component.activeHook)
	return []*Component{component.Component}
}

func (l *LoxoneSwitch) activeHook(event events.Event) {
	l.On.SetValue(event.Value == 1)
}

func (l *LoxoneSwitch) remoteUpdate(on bool) {
	if on {
		l.command("on")
	} else {
		l.command("off")
	}
}
