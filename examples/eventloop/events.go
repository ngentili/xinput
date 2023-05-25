package main

import (
	"time"

	"github.com/ngentili/xinput"
)

const (
	EVENT_LEFT_THUMBSTICK  = 0
	EVENT_RIGHT_THUMBSTICK = 1

	EVENT_LEFT_TRIGGER  = 0
	EVENT_RIGHT_TRIGGER = 1
)

const (
	EVENT_TYPE_BUTTON_PRESS   = 0x01
	EVENT_TYPE_BUTTON_RELEASE = 0x02
	EVENT_TYPE_THUMBSTICK     = 0x04
	EVENT_TYPE_TRIGGER        = 0x08
)

type ThumbstickEvent struct {
	Thumbstick int
	ValueX     int16 // (-32768 to 32767)
	ValueY     int16 // (-32768 to 32767)
}

type TriggerEvent struct {
	Trigger int
	Value   int8 // double check
}

type XInputEvent struct {
	EventType       int
	Button          int
	ThumbstickEvent *ThumbstickEvent
	TriggerEvent    *TriggerEvent
	State           *xinput.XINPUT_STATE
}

var buttons = [...]int{
	xinput.XINPUT_GAMEPAD_DPAD_UP,
	xinput.XINPUT_GAMEPAD_DPAD_DOWN,
	xinput.XINPUT_GAMEPAD_DPAD_LEFT,
	xinput.XINPUT_GAMEPAD_DPAD_RIGHT,
	xinput.XINPUT_GAMEPAD_START,
	xinput.XINPUT_GAMEPAD_BACK,
	xinput.XINPUT_GAMEPAD_LEFT_THUMB,
	xinput.XINPUT_GAMEPAD_RIGHT_THUMB,
	xinput.XINPUT_GAMEPAD_LEFT_SHOULDER,
	xinput.XINPUT_GAMEPAD_RIGHT_SHOULDER,
	xinput.XINPUT_GAMEPAD_A,
	xinput.XINPUT_GAMEPAD_B,
	xinput.XINPUT_GAMEPAD_X,
	xinput.XINPUT_GAMEPAD_Y,
	xinput.XINPUT_GAMEPAD_GUIDE,
}

func GetEventLoop(userIndex int, eventFilter int, tick time.Duration) <-chan *XInputEvent {

	ch := make(chan *XInputEvent)

	go func() {
		defer close(ch)

		prevState := &xinput.XINPUT_STATE{}

		for {

			connected, state := xinput.GetState(userIndex)
			_ = connected

			// if change in state
			if state.PacketNumber != prevState.PacketNumber {

				// buttons
				if eventFilter&(EVENT_TYPE_BUTTON_PRESS|EVENT_TYPE_BUTTON_RELEASE) != 0 {

					for _, button := range buttons {

						isPressed := state.Gamepad.Buttons&uint16(button) != 0
						wasPressed := prevState.Gamepad.Buttons&uint16(button) != 0

						if isPressed != wasPressed {

							e := XInputEvent{
								State:  state,
								Button: int(button),
							}

							if isPressed {
								e.EventType = EVENT_TYPE_BUTTON_PRESS

							} else {
								e.EventType = EVENT_TYPE_BUTTON_RELEASE
							}

							if e.EventType&eventFilter == 0 {
								continue
							}

							ch <- &e
						}
					}
				}

				// thumbsticks
				if eventFilter&EVENT_TYPE_THUMBSTICK != 0 {

					// left thumbstick
					if state.Gamepad.ThumbLX != prevState.Gamepad.ThumbLX || state.Gamepad.ThumbLY != prevState.Gamepad.ThumbLY {

						e := XInputEvent{
							State:     state,
							EventType: EVENT_TYPE_THUMBSTICK,
							ThumbstickEvent: &ThumbstickEvent{
								Thumbstick: EVENT_LEFT_THUMBSTICK,
								ValueX:     int16(state.Gamepad.ThumbLX),
								ValueY:     int16(state.Gamepad.ThumbLY),
							},
						}

						ch <- &e
					}

					// right thumbstick
					if state.Gamepad.ThumbRX != prevState.Gamepad.ThumbRX || state.Gamepad.ThumbRY != prevState.Gamepad.ThumbRY {

						e := XInputEvent{
							State:     state,
							EventType: EVENT_TYPE_THUMBSTICK,
							ThumbstickEvent: &ThumbstickEvent{
								Thumbstick: EVENT_RIGHT_THUMBSTICK,
								ValueX:     int16(state.Gamepad.ThumbRX),
								ValueY:     int16(state.Gamepad.ThumbRY),
							},
						}

						ch <- &e
					}
				}

				// triggers
				if eventFilter&EVENT_TYPE_TRIGGER != 0 {

					// left trigger
					if state.Gamepad.LeftTrigger != prevState.Gamepad.LeftTrigger {

						e := XInputEvent{
							State:     state,
							EventType: EVENT_TYPE_TRIGGER,
							TriggerEvent: &TriggerEvent{
								Trigger: EVENT_LEFT_TRIGGER,
								Value:   int8(state.Gamepad.LeftTrigger),
							},
						}

						ch <- &e
					}

					// right trigger
					if state.Gamepad.RightTrigger != prevState.Gamepad.RightTrigger {

						e := XInputEvent{
							State:     state,
							EventType: EVENT_TYPE_TRIGGER,
							TriggerEvent: &TriggerEvent{
								Trigger: EVENT_RIGHT_TRIGGER,
								Value:   int8(state.Gamepad.RightTrigger),
							},
						}

						ch <- &e
					}
				}

				// store current gamepad state
				prevState = state
			}

			// wait for next tick
			time.Sleep(tick)
		}

	}()

	return ch
}
