package components

import (
	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/accessory"
)

type LoxoneIntercom struct {
	*accessory.Accessory
	loxone *loxone.Loxone
	uuid   string
}

// TODO
func NewIntercom(component Component, control *loxone.Control, lox *loxone.Loxone) *LoxoneIntercom {
	acc := LoxoneIntercom{}
	info := accessory.Info{
		Name:         control.Name,
		Manufacturer: "Loxone",
		SerialNumber: control.UUIDAction,
	}
	acc.Accessory = accessory.New(info, accessory.AccessoryType(component.Type))

	acc.uuid = control.UUIDAction
	acc.loxone = lox

	return &acc
}
