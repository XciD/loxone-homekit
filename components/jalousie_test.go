package components

import (
	"testing"

	"github.com/brutella/hc/characteristic"

	"github.com/magiconair/properties/assert"
)

func TestJalousie(t *testing.T) {
	fixture := NewFixture("test-name", "test-uuid", "Jalousie", map[string]interface{}{
		"up":       "uuid-up",
		"down":     "uuid-down",
		"position": "uuid-position",
	})

	dimmer := NewJalousie(*fixture.ComponentConfig, fixture.Control, fixture.FakeWebsocket)

	characteristics := dimmer.GetServices()[1].GetCharacteristics()
	currentPosition := characteristics[0]
	targetPosition := characteristics[1]
	positionState := characteristics[2]

	// Trigger default values from WS
	fixture.TriggerEvent("uuid-up", 0)
	fixture.TriggerEvent("uuid-down", 0)
	fixture.TriggerEvent("uuid-position", 0)

	assert.Equal(t, currentPosition.GetValue(), 100)
	assert.Equal(t, targetPosition.GetValue(), 100)
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateStopped)

	// Trigger up from WS
	fixture.TriggerEvent("uuid-up", 1)
	fixture.TriggerEvent("uuid-position", 0.1)

	assert.Equal(t, currentPosition.GetValue(), 90)
	assert.Equal(t, targetPosition.GetValue(), 100)
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateIncreasing)

	// Trigger stop from WS
	fixture.TriggerEvent("uuid-up", 0)
	fixture.TriggerEvent("uuid-down", 0)
	fixture.TriggerEvent("uuid-position", 0.2)

	assert.Equal(t, currentPosition.GetValue(), 80)
	assert.Equal(t, targetPosition.GetValue(), 80)
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateStopped)

	// Down from Client
	targetPosition.UpdateValueFromConnection(0, TestConn)
	// callback
	fixture.TriggerEvent("uuid-up", 0)
	fixture.TriggerEvent("uuid-down", 1)
	fixture.TriggerEvent("uuid-position", 0.3)

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/test-uuid/FullDown")
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateDecreasing)
	assert.Equal(t, currentPosition.GetValue(), 70)
	assert.Equal(t, targetPosition.GetValue(), 0)

	// Up from Client should stop the blind
	targetPosition.UpdateValueFromConnection(100, TestConn)
	// callback
	fixture.TriggerEvent("uuid-up", 0)
	fixture.TriggerEvent("uuid-down", 0)
	fixture.TriggerEvent("uuid-position", 0.3)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/test-uuid/stop")
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateStopped)
	assert.Equal(t, currentPosition.GetValue(), 70)
	assert.Equal(t, targetPosition.GetValue(), 70)

	// Up from Client should go up
	targetPosition.UpdateValueFromConnection(10, TestConn)
	// callback
	fixture.TriggerEvent("uuid-up", 1)
	fixture.TriggerEvent("uuid-down", 0)
	fixture.TriggerEvent("uuid-position", 0.7)

	assert.Equal(t, fixture.GetCommands()[2], "jdev/sps/io/test-uuid/ManualPosition/90")
	assert.Equal(t, positionState.GetValue(), characteristic.PositionStateIncreasing)
	assert.Equal(t, currentPosition.GetValue(), 30)
	assert.Equal(t, targetPosition.GetValue(), 10)

}
