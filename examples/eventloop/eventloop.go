package main

import (
	"math"
	"time"

	"github.com/ngentili/xinput"
)

func main() {
	PLAYER_1 := 0

	xinput.SetThumbstickDeadzones(7000, math.MaxInt16, 7000, math.MaxInt16)

	interval := time.Duration(math.Floor((1.0 / 60.0) * float64(time.Second)))

	events := xinput.GetEventLoop(PLAYER_1, xinput.EVENT_TYPE_BUTTON_PRESS, interval)

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
