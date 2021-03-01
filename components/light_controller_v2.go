package components

import (
	"fmt"

	"github.com/brutella/hc/service"
)

type LoxoneLightController struct {
	*Component
	*service.Service
}

func NewLoxoneLightController(f *Factory, config ComponentConfig) []*Component {
	components := make([]*Component, 0)
	for id, control := range config.Control.SubControls {
		subConfig := ComponentConfig{
			ID:           id,
			Name:         fmt.Sprintf("%s %s", config.Name, control.Name),
			CategoryType: config.CategoryType,
			ControlType:  control.Type,
			Control:      control,
		}

		components = append(components, f.CreateComponents(subConfig)...)
	}

	return components
}
