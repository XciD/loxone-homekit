package components

import (
	"fmt"

	"github.com/XciD/loxone-ws/events"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/accessory"
	log "github.com/sirupsen/logrus"
)

type ComponentConfig struct {
	ID         string
	Name       string
	Type       int
	LoxoneType string `mapstructure:"loxone_type"`
}

type Component struct {
	*accessory.Accessory
	Logger *log.Entry
	loxone.WebsocketInterface
	*loxone.Control
	uuid string
}

func newComponent(config ComponentConfig, control *loxone.Control, loxone loxone.WebsocketInterface) *Component {
	c := &Component{}
	c.WebsocketInterface = loxone
	c.Control = control

	c.Logger = log.WithFields(log.Fields{"type": config.LoxoneType, "id": config.ID, "name": config.Name})
	c.uuid = config.ID

	info := accessory.Info{
		Name:         control.Name,
		Manufacturer: "Loxone",
		SerialNumber: control.UUIDAction,
	}
	c.Accessory = accessory.New(info, accessory.AccessoryType(config.Type))
	return c
}

func (c *Component) command(command string) {
	c.Logger.Infof("command %s", command)
	status, err := c.WebsocketInterface.SendCommand(fmt.Sprintf("jdev/sps/io/%s/%s", c.uuid, command), nil)
	if err != nil {
		log.Error(err)
	}
	c.Logger.Infof("Result %d", status.Code)
}

func (c *Component) addHook(stateName string, callback func(*events.Event)) {
	// TODO Check if state exist
	c.WebsocketInterface.AddHook(c.Control.States[stateName].(string), callback)
}

func (c *Component) addDebugHook(stateName string) {
	// TODO Check if state exist
	c.WebsocketInterface.AddHook(c.Control.States[stateName].(string), c.debugHook(stateName))
}

func (c *Component) debugHook(name string) func(event *events.Event) {
	return func(event *events.Event) {
		c.Logger.Infof("Received event from %s with value %.2f", name, event.Value)
	}
}

func CreateComponent(config ComponentConfig, control *loxone.Control, lox loxone.WebsocketInterface) *Component {
	switch config.LoxoneType {
	case "Switch":
		return NewLoxoneSwitch(config, control, lox)
	case "Jalousie":
		return NewJalousie(config, control, lox)
	case "Gate":
		return NewGate(config, control, lox)
	case "EIBDimmer":
		return NewLoxoneDimmer(config, control, lox)
	default:
		log.Warnf("Unknown component %s", config.LoxoneType)
		return nil
	}
}
