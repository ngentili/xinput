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

	events := GetEventLoop(0, EVENT_TYPE_BUTTON_PRESS, tick)

	// block until a new event is received
	for event := range events {

		// if LB is held
		if event.State.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_LEFT_SHOULDER != 0 {

			// if D-pad Up is pressed
			if event.Button == xinput.XINPUT_GAMEPAD_DPAD_UP {
				println("LB + D-Up")

				// if D-pad Down is pressed
			} else if event.Button == xinput.XINPUT_GAMEPAD_DPAD_UP {
				println("LB + D-Down")

			}
		}
	}
}
