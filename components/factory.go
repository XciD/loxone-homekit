package components

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/viper"

	"github.com/XciD/loxone-ws/events"
	"github.com/XciD/loxone-ws/types"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc/accessory"
	log "github.com/sirupsen/logrus"
)

type Factory struct {
	Loxone loxone.Loxone
	Config FactoryConfig
}

type FactoryConfig struct {
	Switch            bool
	Jalousie          bool
	Alarm             bool
	LightControllerV2 bool
	ColorPickerV2     bool
	Gate              bool
	EIBDimmer         bool
	Dimmer            bool
}

type ComponentConfig struct {
	ID           string
	Name         string
	ControlType  string
	CategoryType string
	Control      *loxone.Control
}

type Component struct {
	*accessory.Accessory
	Logger        *log.Entry
	loxone        loxone.Loxone
	loxoneControl *loxone.Control
	uuid          string
}

func (f *Factory) newComponent(config ComponentConfig, typ accessory.AccessoryType) *Component {
	c := &Component{}
	c.loxone = f.Loxone
	c.loxoneControl = config.Control
	c.uuid = config.ID

	ID := getID(config.ID)

	c.Logger = log.WithFields(log.Fields{
		"type":  config.ControlType,
		"id":    config.ID,
		"name":  config.Name,
		"hc-id": ID,
	})

	info := accessory.Info{
		ID:           ID,
		Name:         config.Name,
		Manufacturer: "Loxone",
		SerialNumber: config.ID,
	}
	c.Accessory = accessory.New(info, typ)
	return c
}

func (c *Component) command(command string) {
	c.Logger.Infof("command %s", command)
	status, err := c.loxone.SendCommand(fmt.Sprintf("jdev/sps/io/%s/%s", c.uuid, command), nil)
	if err != nil {
		log.Error(err)
	}
	c.Logger.Infof("Result %d", status.Code)
}

func (c *Component) addHook(stateName string, callback func(events.Event)) {
	// TODO Check if state exist
	c.loxone.AddHook(c.loxoneControl.States[stateName].(string), callback)
}

func (c *Component) addDebugHook(stateName string) {
	// TODO Check if state exist
	c.loxone.AddHook(c.loxoneControl.States[stateName].(string), c.debugHook(stateName))
}

func (c *Component) debugHook(name string) func(event events.Event) {
	return func(event events.Event) {
		c.Logger.Infof("Received event from %s with value %+v", name, event)
	}
}

func NewFactory(loxone loxone.Loxone) *Factory {
	var config FactoryConfig

	if err := viper.UnmarshalKey("autodiscover", &config); err != nil {
		panic(err)
	}

	return &Factory{
		Loxone: loxone,
		Config: config,
	}
}

func (f *Factory) CreateComponents(config ComponentConfig) []*Component {
	log.Infof("Creating control %s", config.ID)
	switch types.Type(config.ControlType) {
	case types.Switch:
		if f.Config.Switch {
			return NewLoxoneSwitch(f, config)
		}
	case types.Jalousie:
		if f.Config.Jalousie {
			return NewJalousie(f, config)
		}
	case types.Gate:
		if f.Config.Gate {
			return NewGate(f, config)
		}
	case types.EIBDimmer:
		if f.Config.EIBDimmer {
			return NewLoxoneDimmer(f, config)
		}
	case types.Dimmer:
		if f.Config.Dimmer {
			return NewLoxoneDimmer(f, config)
		}
	case types.Alarm:
		if f.Config.Alarm {
			return NewAlarm(f, config)
		}
	case types.LightControllerV2:
		if f.Config.LightControllerV2 {
			return NewLoxoneLightController(f, config)
		}
	case types.ColorPickerV2:
		if f.Config.ColorPickerV2 {
			return NewLoxoneColorPicker(f, config)
		}
	default:
		log.Warnf("Unknown component %s, %s, %s", config.ControlType, config.Control.Name, config.ID)
	}
	return []*Component{}
}

func GetAccessories(loxone loxone.Loxone, config *loxone.Config) []*accessory.Accessory {
	factory := NewFactory(loxone)

	accessories := make([]*accessory.Accessory, 0)

	for id, control := range config.Controls {
		config := ComponentConfig{
			ID:           id,
			Name:         fmt.Sprintf("%s %s", config.RoomName(control.Room), control.Name),
			CategoryType: config.Cats[control.Cat].Type,
			ControlType:  control.Type,
			Control:      control,
		}

		components := factory.CreateComponents(config)
		for _, c := range components {
			accessories = append(accessories, c.Accessory)
		}
	}

	return accessories
}

func getID(uuid string) uint64 {
	uuid = strings.ReplaceAll(uuid, "-", "")
	removeSlash := strings.Split(uuid, "/")

	hexID, err := hex.DecodeString(removeSlash[0])

	if err != nil {
		panic(err)
	}

	firstPart := binary.BigEndian.Uint64(hexID)

	if len(removeSlash) > 1 {
		// Only we way to create a unique uint64...
		text := []byte(removeSlash[1])
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(text)))
		base64.StdEncoding.Encode(dst, text)
		var i big.Int
		i.SetBytes(dst)
		firstPart += i.Uint64()
	}

	return firstPart
}
