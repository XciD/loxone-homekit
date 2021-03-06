package components

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestDimmer(t *testing.T) {
	fixture := NewFixture("test-name", "10516e85-0239-2320-ffffb7fe2005e936", "EIBDimmer", map[string]interface{}{"position": "uuid-position"})

	dimmer := NewLoxoneDimmer(fixture.Factory, *fixture.ComponentConfig)

	characteristics := dimmer[0].GetServices()[1].GetCharacteristics()
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

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/10516e85-0239-2320-ffffb7fe2005e936/off")

	// LightOn from client
	on.UpdateValueFromConnection(true, TestConn)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/10516e85-0239-2320-ffffb7fe2005e936/on")

	// Brigtness to 50
	brigtness.UpdateValueFromConnection(50, TestConn)

	assert.Equal(t, fixture.GetCommands()[2], "jdev/sps/io/10516e85-0239-2320-ffffb7fe2005e936/50")

}
