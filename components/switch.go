package components

import (
	"fmt"

	"github.com/prometheus/common/log"

	"github.com/XciD/loxone-ws/events"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type LoxoneSwitch struct {
	*accessory.Accessory
	service *SwitchService
	loxone  *loxone.Loxone
	uuid    string
	state   bool
}

type SwitchService struct {
	*service.Service
	On *characteristic.On
}

func NewLoxoneSwitch(component Component, control *loxone.Control, lox *loxone.Loxone) *LoxoneSwitch {
	acc := LoxoneSwitch{}
	info := accessory.Info{
		Name:         control.Name,
		Manufacturer: "Loxone",
		SerialNumber: control.UUIDAction,
	}
	acc.Accessory = accessory.New(info, accessory.AccessoryType(component.Type))
	acc.service = newLightService()

	acc.AddService(acc.service.Service)

	acc.uuid = control.UUIDAction
	acc.loxone = lox
	acc.state = true

	acc.service.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			acc.command("on")
		} else {
			acc.command("off")
		}
	})

	acc.service.On.OnValueGet(func() interface{} {
		return acc.state
	})

	lox.AddHook(control.States["active"].(string), func(event *events.Event) {
		acc.state = event.Value == 1
		acc.service.On.SetValue(acc.state)
	})

	return &acc
}

func newLightService() *SwitchService {
	svc := SwitchService{}
	svc.Service = service.New(service.TypeLightbulb)

	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	return &svc
}

func (l *LoxoneSwitch) command(command string) {
	log.Infof("[%s] Stopping", l.Info.Name.Value)
	status, err := l.loxone.SendCommand(fmt.Sprintf("jdev/sps/io/%s/%s", l.uuid, command), nil)
	if err != nil {
		log.Error(err)
	}
	log.Infof("[%s] Result %d", l.Info.Name.Value, status.Code)
}
