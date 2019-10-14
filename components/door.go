package components

import (
	"fmt"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	log "github.com/sirupsen/logrus"
)

type LoxoneDoor struct {
	*accessory.Accessory
	service        *service.GarageDoorOpener
	loxone         *loxone.Loxone
	uuid           string
	position       float64
	targetPosition int
	active         int
}

func NewDoor(component Component, control *loxone.Control, lox *loxone.Loxone) *LoxoneDoor {
	acc := LoxoneDoor{}
	info := accessory.Info{
		Name:         component.Name,
		Manufacturer: "Loxone",
		SerialNumber: control.UUIDAction,
	}
	acc.Accessory = accessory.New(info, accessory.AccessoryType(component.Type))
	acc.service = service.NewGarageDoorOpener()

	acc.AddService(acc.service.Service)

	acc.uuid = control.UUIDAction
	acc.loxone = lox

	acc.position = 0
	acc.active = 0
	acc.targetPosition = 1

	acc.service.CurrentDoorState.OnValueGet(func() interface{} {
		return acc.getState()
	})

	acc.service.TargetDoorState.OnValueGet(func() interface{} {
		return acc.targetPosition
	})

	acc.service.TargetDoorState.OnValueRemoteUpdate(func(i int) {
		acc.targetPosition = i
		if i == characteristic.TargetDoorStateOpen {
			log.Infof("[%s] Asking to open the door", acc.Info.Name.Value)
			acc.command("open")
		} else {
			log.Infof("[%s] Asking to close the door", acc.Info.Name.Value)
			acc.command("close")
		}
	})

	lox.AddHook(control.States["position"].(string), func(event *events.Event) {
		acc.position = event.Value
		acc.setState(acc.getState())
	})
	lox.AddHook(control.States["active"].(string), func(event *events.Event) {
		acc.active = int(event.Value)
		acc.setState(acc.getState())
	})

	lox.AddHook(control.States["preventOpen"].(string), func(event *events.Event) {
		log.Infof("[%s] preventOpen %f", acc.Info.Name.Value, event.Value)
	})
	lox.AddHook(control.States["preventClose"].(string), func(event *events.Event) {
		log.Infof("[%s] preventClose %f", acc.Info.Name.Value, event.Value)
	})

	return &acc
}

func (l *LoxoneDoor) getState() int {
	switch l.active {
	case -1:
		log.Infof("[%s] Door is closing", l.Info.Name.Value)
		return characteristic.CurrentDoorStateClosing
	case 1:
		log.Infof("[%s] Door is Opening", l.Info.Name.Value)
		return characteristic.CurrentDoorStateOpening
	default:
		switch l.position {
		case 0:
			log.Infof("[%s] Door is Closed", l.Info.Name.Value)
			return characteristic.CurrentDoorStateClosed
		case 1:
			log.Infof("[%s] Door is Open", l.Info.Name.Value)
			return characteristic.CurrentDoorStateOpen
		default:
			log.Infof("[%s] Door is Stopped", l.Info.Name.Value)
			return characteristic.CurrentDoorStateStopped
		}
	}
}

func (l *LoxoneDoor) setState(state int) {
	l.service.CurrentDoorState.SetValue(state)
}

func (l *LoxoneDoor) command(command string) {
	log.Infof("[%s] Stopping", l.Info.Name.Value)
	status, err := l.loxone.SendCommand(fmt.Sprintf("jdev/sps/io/%s/%s", l.uuid, command), nil)
	if err != nil {
		log.Error(err)
	}
	log.Infof("[%s] Result %d", l.Info.Name.Value, status.Code)
}
