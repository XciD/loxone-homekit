package components

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestSwitch(t *testing.T) {
	fixture := NewFixture("test-name", "10516e85-0239-2320-ffffb7fe2005e936", "Switch", map[string]interface{}{"active": "uuid-active"})

	fixture.ComponentConfig.CategoryType = "lights"

	loxoneSwitch := NewLoxoneSwitch(fixture.Factory, *fixture.ComponentConfig)[0]

	characteristics := loxoneSwitch.GetServices()[1].GetCharacteristics()
	on := characteristics[0]

	// Trigger down from WS
	fixture.TriggerEvent("uuid-active", 1)

	assert.Equal(t, on.GetValue(), true)

	// Trigger 100 from WS
	fixture.TriggerEvent("uuid-active", 0)
	assert.Equal(t, on.GetValue(), false)

	// Shutdown from client
	on.UpdateValueFromConnection(true, TestConn)

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/10516e85-0239-2320-ffffb7fe2005e936/on")
	assert.Equal(t, on.GetValue(), true)

	// LightOn from client
	on.UpdateValueFromConnection(false, TestConn)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/10516e85-0239-2320-ffffb7fe2005e936/off")
	assert.Equal(t, on.GetValue(), false)

}
