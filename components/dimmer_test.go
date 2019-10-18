package components

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestDimmer(t *testing.T) {
	fixture := NewFixture("test-name", "test-uuid", "EIBDimmer", map[string]interface{}{"position": "uuid-position"})

	dimmer := NewLoxoneDimmer(*fixture.ComponentConfig, fixture.Control, fixture.LoxoneInterface)

	characteristics := dimmer.GetServices()[1].GetCharacteristics()
	on := characteristics[0]
	brigtness := characteristics[1]

	// Trigger down from WS
	fixture.TriggerEvent("uuid-position", 0)

	assert.Equal(t, on.GetValue(), false)
	assert.Equal(t, brigtness.GetValue(), 0)

	// Trigger 100 from WS
	fixture.TriggerEvent("uuid-position", 100)
	assert.Equal(t, on.GetValue(), true)
	assert.Equal(t, brigtness.GetValue(), 100)

	// Shutdown from client
	on.UpdateValueFromConnection(false, TestConn)

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/test-uuid/off")

	// LightOn from client
	on.UpdateValueFromConnection(true, TestConn)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/test-uuid/on")

	// Brigtness to 50
	brigtness.UpdateValueFromConnection(50, TestConn)

	assert.Equal(t, fixture.GetCommands()[2], "jdev/sps/io/test-uuid/50")

}
