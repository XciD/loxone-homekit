package components

import (
	"testing"

	"github.com/brutella/hc/characteristic"
	"github.com/magiconair/properties/assert"
)

func TestGate(t *testing.T) {
	fixture := NewFixture("test-name", "test-uuid", "Gate", map[string]interface{}{
		"position":     "uuid-position",
		"preventOpen":  "uuid-preventOpen",
		"preventClose": "uuid-preventClose",
		"active":       "uuid-active",
	})

	dimmer := NewGate(*fixture.ComponentConfig, fixture.Control, fixture.LoxoneInterface)

	characteristics := dimmer.GetServices()[1].GetCharacteristics()
	currentDoorState := characteristics[0]
	targetDoorState := characteristics[1]
	currentPosition := characteristics[3]

	// Trigger event from WS
	fixture.TriggerEvent("uuid-active", 0)
	fixture.TriggerEvent("uuid-position", 0)

	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateClosed, "Door should be closed")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateClosed, "Target should be set same as door")
	assert.Equal(t, currentPosition.GetValue(), 100)

	// Trigger event from WS
	fixture.TriggerEvent("uuid-active", 1)
	fixture.TriggerEvent("uuid-position", 0.1)

	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateOpening, "Door should opening")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateOpen, "Target should be updated")
	assert.Equal(t, currentPosition.GetValue(), 90)

	// Trigger event from WS
	fixture.TriggerEvent("uuid-active", -1)
	fixture.TriggerEvent("uuid-position", 0.2)

	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateClosing, "Door should closing")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateClosed, "Target should be updated")
	assert.Equal(t, currentPosition.GetValue(), 80)

	// Trigger event from Client
	targetDoorState.UpdateValueFromConnection(characteristic.TargetDoorStateOpen, TestConn)

	// Send result from WS
	fixture.TriggerEvent("uuid-active", 0)
	fixture.TriggerEvent("uuid-position", 0.3)

	assert.Equal(t, fixture.GetCommands()[0], "jdev/sps/io/test-uuid/stop")
	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateStopped, "Door should stop")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateOpen, "Target should equal open")
	assert.Equal(t, currentPosition.GetValue(), 70)

	// Trigger event from Client, close a stopped door
	targetDoorState.UpdateValueFromConnection(characteristic.TargetDoorStateClosed, TestConn)

	// Send result from WS
	fixture.TriggerEvent("uuid-active", -1)
	fixture.TriggerEvent("uuid-position", 0.3)

	assert.Equal(t, fixture.GetCommands()[1], "jdev/sps/io/test-uuid/close")
	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateClosing, "Door should be closing")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateClosed, "Target should equal closed")
	assert.Equal(t, currentPosition.GetValue(), 70)

	// Trigger event from Client, open the door
	fixture.TriggerEvent("uuid-active", 0)
	targetDoorState.UpdateValueFromConnection(characteristic.TargetDoorStateOpen, TestConn)

	// Send result from WS
	fixture.TriggerEvent("uuid-active", 1)
	fixture.TriggerEvent("uuid-position", 0.3)

	assert.Equal(t, fixture.GetCommands()[2], "jdev/sps/io/test-uuid/open")
	assert.Equal(t, currentDoorState.GetValue(), characteristic.CurrentDoorStateOpening, "Door should be opening")
	assert.Equal(t, targetDoorState.GetValue(), characteristic.TargetDoorStateOpen, "Target should equal to open")
	assert.Equal(t, currentPosition.GetValue(), 70)
}
