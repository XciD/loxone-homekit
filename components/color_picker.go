package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brutella/hc/accessory"

	"github.com/XciD/loxone-ws/events"

	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneColorPicker struct {
	*Component
	*service.Service
	*characteristic.On
	*characteristic.Brightness
	*characteristic.Saturation
	*characteristic.Hue
}

func NewLoxoneColorPicker(f *Factory, config ComponentConfig) []*Component {
	component := &LoxoneColorPicker{
		Component: f.newComponent(config, accessory.TypeLightbulb),
	}

	component.Service = service.New(service.TypeLightbulb)
	component.AddService(component.Service)

	component.On = characteristic.NewOn()
	component.AddCharacteristic(component.On.Characteristic)

	component.Brightness = characteristic.NewBrightness()
	component.AddCharacteristic(component.Brightness.Characteristic)

	component.Hue = characteristic.NewHue()
	component.AddCharacteristic(component.Hue.Characteristic)

	component.Saturation = characteristic.NewSaturation()
	component.AddCharacteristic(component.Saturation.Characteristic)

	component.On.OnValueRemoteUpdate(component.onOnRemoteUpdate)
	component.On.OnValueGet(component.onOnGet)

	component.Brightness.OnValueRemoteUpdate(func(i int) {
		component.onUpdate()
	})
	component.Hue.OnValueRemoteUpdate(func(f float64) {
		component.onUpdate()
	})
	component.Saturation.OnValueRemoteUpdate(func(f float64) {
		component.onUpdate()
	})

	component.addHook("color", component.colorHook)

	return []*Component{component.Component}
}

func (l *LoxoneColorPicker) colorHook(event events.Event) {
	if !strings.HasPrefix(event.Text, "hsv") {
		return
	}
	HSV := event.Text[4 : len(event.Text)-1]
	HSVSlice := strings.Split(HSV, ",")
	if len(HSVSlice) != 3 {
		return
	}
	H, err := strconv.Atoi(HSVSlice[0])
	if err != nil {
		l.Logger.WithError(err).Error("Invalid parse")
		return
	}
	S, err := strconv.Atoi(HSVSlice[1])
	if err != nil {
		l.Logger.WithError(err).Error("Invalid parse")
		return
	}
	V, err := strconv.Atoi(HSVSlice[2])
	if err != nil {
		l.Logger.WithError(err).Error("Invalid parse")
		return
	}
	if V > 0 {
		l.On.SetValue(true)
	} else {
		l.On.SetValue(false)
	}
	l.Hue.SetValue(float64(H))
	l.Saturation.SetValue(float64(S))
	l.Brightness.SetValue(V)
}

func (l *LoxoneColorPicker) onOnRemoteUpdate(on bool) {
	if on {
		l.command("on")
	} else {
		l.command("off")
	}
}

func (l *LoxoneColorPicker) onUpdate() {
	l.command(fmt.Sprintf("hsv(%d,%d,%d)", int64(l.Hue.GetValue()), int64(l.Saturation.GetValue()), l.Brightness.Value))
}

func (l *LoxoneColorPicker) onOnGet() interface{} {
	return l.Brightness.GetValue() > 0
}
