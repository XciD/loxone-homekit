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

type LoxoneJalousie struct {
	*accessory.Accessory
	service        *service.WindowCovering
	loxone         *loxone.Loxone
	uuid           string
	position       int
	targetPosition int
	up             bool
	down           bool
}

func NewJalousie(component Component, control *loxone.Control, lox *loxone.Loxone) *LoxoneJalousie {
	acc := LoxoneJalousie{}
	info := accessory.Info{
		Name:         component.Name,
		Manufacturer: "Loxone",
		SerialNumber: control.UUIDAction,
	}
	acc.Accessory = accessory.New(info, accessory.AccessoryType(component.Type))
	acc.service = service.NewWindowCovering()

	acc.AddService(acc.service.Service)

	acc.uuid = control.UUIDAction
	acc.loxone = lox

	acc.targetPosition = 0
	acc.position = 0
	acc.up = false
	acc.down = false

	acc.service.TargetPosition.OnValueRemoteUpdate(func(i int) {
		log.Infof("[%s] setPosition %d, Target position %d, state %d, asking position %d", acc.Info.Name.Value, acc.position, acc.targetPosition, acc.getState(), i)

		if acc.getState() != characteristic.PositionStateStopped {
			acc.stop()
			acc.setState(characteristic.PositionStateStopped)
			acc.service.TargetPosition.SetValue(acc.position)
			return
		}
		acc.targetPosition = i
		acc.setPosition(i)
	})

	acc.service.TargetPosition.OnValueGet(func() interface{} {
		return acc.targetPosition
	})

	acc.service.CurrentPosition.OnValueGet(func() interface{} {
		return acc.position
	})

	acc.service.PositionState.OnValueGet(func() interface{} {
		return acc.getState()
	})

	lox.AddHook(control.States["up"].(string), func(event *events.Event) {
		acc.up = event.Value == 1
		log.Infof("[%s] Updating state up to %t", acc.Info.Name.Value, acc.up)
		acc.setState(acc.getState())
	})
	lox.AddHook(control.States["down"].(string), func(event *events.Event) {
		acc.down = event.Value == 1
		log.Infof("[%s] Updating state down to %t", acc.Info.Name.Value, acc.down)
		acc.setState(acc.getState())
	})

	lox.AddHook(control.States["position"].(string), func(event *events.Event) {
		acc.position = int(100 - (event.Value * 100))
		log.Infof("[%s] setPosition from Websocket %d", acc.Info.Name.Value, acc.position)
		acc.service.CurrentPosition.SetValue(acc.position)
		state := acc.getState()

		acc.setState(state)

		if state == characteristic.PositionStateStopped {
			acc.targetPosition = acc.position
			acc.service.TargetPosition.SetValue(acc.position)
		}

		log.Infof("[%s] setPosition %d, Target position %d, state %d", acc.Info.Name.Value, acc.position, acc.targetPosition, acc.getState())
	})

	return &acc
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
	l.service.PositionState.SetValue(state)
}

func (l *LoxoneJalousie) setPosition(position int) {
	pos := int32(100 - position)
	log.Infof("[%s] Asking for loxone position %d", l.Info.Name.Value, pos)

	command := fmt.Sprintf("ManualPosition/%d", pos)
	if pos == 0 {
		command = "FullUp"
	} else if pos == 100 {
		command = "FullDown"
	}

	status, err := l.loxone.SendCommand(fmt.Sprintf("jdev/sps/io/%s/%s", l.uuid, command), nil)
	if err != nil {
		log.Error(err)
	}
	log.Infof("[%s] Result %d", l.Info.Name.Value, status.Code)
}

func (l *LoxoneJalousie) stop() {
	log.Infof("[%s] Stopping", l.Info.Name.Value)
	status, err := l.loxone.SendCommand(fmt.Sprintf("jdev/sps/io/%s/stop", l.uuid), nil)
	if err != nil {
		log.Error(err)
	}
	log.Infof("[%s] Result %d", l.Info.Name.Value, status.Code)
}
