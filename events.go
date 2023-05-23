package xinput

import "time"

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
	State           *XINPUT_STATE
}

var buttons = [...]int{
	XINPUT_GAMEPAD_DPAD_UP,
	XINPUT_GAMEPAD_DPAD_DOWN,
	XINPUT_GAMEPAD_DPAD_LEFT,
	XINPUT_GAMEPAD_DPAD_RIGHT,
	XINPUT_GAMEPAD_START,
	XINPUT_GAMEPAD_BACK,
	XINPUT_GAMEPAD_LEFT_THUMB,
	XINPUT_GAMEPAD_RIGHT_THUMB,
	XINPUT_GAMEPAD_LEFT_SHOULDER,
	XINPUT_GAMEPAD_RIGHT_SHOULDER,
	XINPUT_GAMEPAD_A,
	XINPUT_GAMEPAD_B,
	XINPUT_GAMEPAD_X,
	XINPUT_GAMEPAD_Y,
	XINPUT_GAMEPAD_GUIDE,
}

func GetEventLoop(userIndex int, eventFilter int, pollInterval time.Duration) <-chan *XInputEvent {

	ch := make(chan *XInputEvent)

	go func() {
		defer close(ch)

		prev_state := &XINPUT_STATE{}

		for {

			state, err := GetState(userIndex)
			if err != nil {
				panic(err)
			}

			// buttons
			if eventFilter&(EVENT_TYPE_BUTTON_PRESS|EVENT_TYPE_BUTTON_RELEASE) != 0 {

				for _, button := range buttons {

					isPressed := state.Gamepad.Buttons&uint16(button) != 0
					wasPressed := prev_state.Gamepad.Buttons&uint16(button) != 0

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
				if state.Gamepad.ThumbLX != prev_state.Gamepad.ThumbLX || state.Gamepad.ThumbLY != prev_state.Gamepad.ThumbLY {

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
				if state.Gamepad.ThumbRX != prev_state.Gamepad.ThumbRX || state.Gamepad.ThumbRY != prev_state.Gamepad.ThumbRY {

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
				if state.Gamepad.LeftTrigger != prev_state.Gamepad.LeftTrigger {

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
				if state.Gamepad.RightTrigger != prev_state.Gamepad.RightTrigger {

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

			prev_state = state

			time.Sleep(pollInterval)
		}

	}()

	return ch
}
