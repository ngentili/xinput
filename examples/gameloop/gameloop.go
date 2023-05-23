package main

import (
	"math"
	"syscall"
	"time"

	"github.com/ngentili/xinput"
	"golang.org/x/sys/windows"
)

func main() {
	PLAYER_1 := 0

	xinput.SetThumbstickDeadzones(7000, math.MaxInt16, 7000, math.MaxInt16)

	prevState := &xinput.XINPUT_STATE{}

	for {
		// get current gamepad state
		state, err := xinput.GetState(PLAYER_1)

		if err != nil {
			errno, isWinErr := err.(syscall.Errno)

			// ignore device not connected error
			if !isWinErr || errno != windows.ERROR_DEVICE_NOT_CONNECTED {
				panic(err)
			}
		}

		// gamepad buttons
		dpadUp := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_DPAD_UP != 0
		dpadDown := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_DPAD_DOWN != 0
		dpadLeft := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_DPAD_LEFT != 0
		dpadRight := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_DPAD_RIGHT != 0
		start := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_START != 0
		back := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_BACK != 0
		leftThumbClick := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_LEFT_THUMB != 0
		rightThumbClick := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_RIGHT_THUMB != 0
		leftShoulder := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_LEFT_SHOULDER != 0
		rightShoulder := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_RIGHT_SHOULDER != 0
		a := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_A != 0
		b := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_B != 0
		x := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_X != 0
		y := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_Y != 0
		guide := state.Gamepad.Buttons&xinput.XINPUT_GAMEPAD_GUIDE != 0

		// gamepad thumbstick
		leftThumbX, leftThumbY := state.Gamepad.ThumbLX, state.Gamepad.ThumbLY
		rightThumbX, rightThumbY := state.Gamepad.ThumbRX, state.Gamepad.ThumbRY

		// gamepad trigger
		leftTrigger := state.Gamepad.LeftTrigger
		rightTrigger := state.Gamepad.RightTrigger

		prevState = state

		time.Sleep(1 / 60)

		if false {
			_ = prevState
			_ = dpadUp
			_ = dpadDown
			_ = dpadLeft
			_ = dpadRight
			_ = start
			_ = back
			_ = leftThumbClick
			_ = rightThumbClick
			_ = leftShoulder
			_ = rightShoulder
			_ = a
			_ = b
			_ = x
			_ = y
			_ = guide
			_ = leftThumbX
			_ = leftThumbY
			_ = rightThumbX
			_ = rightThumbY
			_ = leftTrigger
			_ = rightTrigger
		}
	}
}
