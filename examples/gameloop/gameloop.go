package main

import (
	"time"

	"github.com/ngentili/xinput"
)

func main() {
	// tickrate
	tickHz := 60.0
	tick := time.Duration((1.0 / tickHz) * float64(time.Second))

	// initialize controller
	xbc := xinput.NewXboxController(0)
	xbc.SetThumbstickDeadzones(5000, 5000)

	// keep track of previous controller state
	prevState := &xinput.XINPUT_STATE{}

	for {
		// get current gamepad state
		connected, state := xbc.GetState()
		_ = connected

		// if change in state
		if state.PacketNumber != prevState.PacketNumber {

			// query gamepad state
			aButtonWasHeld := prevState.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_A != 0
			aButtonIsHeld := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_A != 0

			// game logic
			if aButtonWasHeld && !aButtonIsHeld {
				println("Released A button")

			} else if !aButtonWasHeld && aButtonIsHeld {
				println("Pressed A button")
			}

			// store current gamepad state
			prevState = state
		}

		// wait for next tick
		time.Sleep(tick)
	}
}
