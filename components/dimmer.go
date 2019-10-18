package components

import (
	"strconv"

	"github.com/XciD/loxone-ws/events"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneDimmer struct {
	*Component
	*service.Service
	*characteristic.On
	*characteristic.Brightness
}

func NewLoxoneDimmer(config ComponentConfig, control *loxone.Control, lox loxone.WebsocketInterface) *Component {
	component := &LoxoneDimmer{
		Component: newComponent(config, control, lox),
	}

	component.Service = service.New(service.TypeLightbulb)
	component.AddService(component.Service)

	component.On = characteristic.NewOn()
	component.AddCharacteristic(component.On.Characteristic)

	component.Brightness = characteristic.NewBrightness()
	component.AddCharacteristic(component.Brightness.Characteristic)

	component.On.OnValueRemoteUpdate(component.onOnRemoteUpdate)
	component.On.OnValueGet(component.onOnGet)

	component.Brightness.OnValueRemoteUpdate(component.onBrightnessUpdate)

	// Add status updates
	component.addHook("position", component.brightnessHook)
	return component.Component
}

func (l *LoxoneDimmer) brightnessHook(event *events.Event) {
	l.Brightness.SetValue(int(event.Value))
	if l.Brightness.GetValue() > 0 {
		l.On.SetValue(true)
	} else {
		l.On.SetValue(false)
	}
}

func (l *LoxoneDimmer) onOnRemoteUpdate(on bool) {
	if on {
		l.command("on")
	} else {
		l.command("off")
	}
}

func (l *LoxoneDimmer) onBrightnessUpdate(value int) {
	l.command(strconv.Itoa(value))
}

func (l *LoxoneDimmer) onOnGet() interface{} {
	return l.Brightness.GetValue() > 0
}
