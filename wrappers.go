package xinput

import "golang.org/x/sys/windows"

type XInputApi struct {
	XInputEnable                    *windows.Proc
	XInputGetAudioDeviceIds         *windows.Proc
	XInputGetDSoundAudioDeviceGuids *windows.Proc
	XInputGetBatteryInformation     *windows.Proc
	XInputGetCapabilities           *windows.Proc
	XInputGetKeystroke              *windows.Proc
	XInputGetState                  *windows.Proc
	XInputSetState                  *windows.Proc
	XInputOrdinal100                *windows.Proc
	XInputOrdinal101                *windows.Proc
	XInputOrdinal102                *windows.Proc
	XInputOrdinal103                *windows.Proc
	XInputOrdinal104                *windows.Proc
	XInputOrdinal108                *windows.Proc
}

type Deadzones struct {
	LeftThumb  int16
	RightThumb int16
	Triggers   int8
}

type XboxController struct {
	UserIndex int
	Deadzones Deadzones
}

func NewXboxController(userIndex int) *XboxController {
	return &XboxController{
		UserIndex: userIndex,
		Deadzones: Deadzones{
			LeftThumb:  XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE,
			RightThumb: XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE,
			Triggers:   XINPUT_GAMEPAD_TRIGGER_THRESHOLD,
		},
	}
}

func (c *XboxController) SetThumbstickDeadzones(left, right int16) {
	if left < 0 {
		panic(left)
	}
	if right < 0 {
		panic(right)
	}
}

func (c *XboxController) GetCapabilities() (*XINPUT_CAPABILITIES, error) {
	return GetCapabilities(c.UserIndex)
}

func (c *XboxController) GetGamepadBatteryInformation() (*XINPUT_BATTERY_INFORMATION, error) {
	return GetBatteryInformation(c.UserIndex, BATTERY_DEVTYPE_GAMEPAD)
}

func (c *XboxController) GetHeadsetBatteryInformation() (*XINPUT_BATTERY_INFORMATION, error) {
	return GetBatteryInformation(c.UserIndex, BATTERY_DEVTYPE_HEADSET)
}

func (c *XboxController) GetState() (bool, *XINPUT_STATE) {
	return GetState(c.UserIndex)
}

func (c *XboxController) GetStateEx() (bool, *XINPUT_STATE) {
	return GetStateEx(c.UserIndex)
}
