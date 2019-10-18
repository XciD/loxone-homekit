package components

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestSwitch(t *testing.T) {
	fixture := NewFixture("test-name", "test-uuid", "Switch", map[string]interface{}{"active": "uuid-active"})

	dimmer := NewLoxoneSwitch(*fixture.ComponentConfig, fixture.Control, fixture.FakeWebsocket)

	characteristics := dimmer.GetServices()[1].GetCharacteristics()
	on := characteristics[0]

	// Trigger down from WS
	fixture.TriggerEvent("uuid-active", 1)

	assert.Equal(t, on.GetValue(), true)

	// Trigger 100 from WS
	fixture.TriggerEvent("uuid-active", 0)
	assert.Equal(t, on.GetValue(), false)

	// Shutdown from client
	on.UpdateValueFromConnection(true, TestConn)

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/test-uuid/on")
	assert.Equal(t, on.GetValue(), true)

	// LightOn from client
	on.UpdateValueFromConnection(false, TestConn)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/test-uuid/off")
	assert.Equal(t, on.GetValue(), false)

}
