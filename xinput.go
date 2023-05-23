package xinput

import (
	"errors"
	"math"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dll = windows.MustLoadDLL("xinput1_4.dll")

	// XInputGetState              = dll.MustFindProc("XInputGetState")
	XInputGetState              = dll.MustFindProcByOrdinal(100) // XInputGetStateEx
	XInputSetState              = dll.MustFindProc("XInputSetState")
	XInputGetCapabilities       = dll.MustFindProc("XInputGetCapabilities")
	XInputGetAudioDeviceIds     = dll.MustFindProc("XInputGetAudioDeviceIds")
	XInputGetBatteryInformation = dll.MustFindProc("XInputGetBatteryInformation")
	XInputGetKeystroke          = dll.MustFindProc("XInputGetKeystroke")

	deadzoneLeftThumbInner = XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE
	deadzoneLeftThumbOuter = math.MaxInt16

	deadzoneRightThumbInner = XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE
	deadzoneRightThumbOuter = math.MaxInt16

	// deadzoneT = XINPUT_GAMEPAD_TRIGGER_THRESHOLD
)

// default is (7849, 32767, 8689, 32767)
func SetThumbstickDeadzones(leftInner, leftOuter, rightInner, rightOuter int16) {
	if leftInner < 0 {
		panic(leftInner)
	}
	if leftOuter < 0 {
		panic(leftOuter)
	}
	if rightInner < 0 {
		panic(rightInner)
	}
	if rightOuter < 0 {
		panic(rightOuter)
	}

	deadzoneLeftThumbInner = int(leftInner)
	deadzoneLeftThumbOuter = int(leftOuter)
	deadzoneRightThumbInner = int(rightInner)
	deadzoneRightThumbOuter = int(rightOuter)
}

func SetTriggerDeadzones(leftInner, leftOuter, rightInner, rightOuter int8) {
	// TODO
}

func GetCapabilities(userIndex int) (*XINPUT_CAPABILITIES, error) {
	caps := XINPUT_CAPABILITIES{}

	if userIndex < 0 || userIndex > 3 {
		return &caps, errors.New("Invalid user index (0-3)")
	}

	caps_p := uintptr(unsafe.Pointer(&caps))

	result, _, _ := XInputGetCapabilities.Call(uintptr(userIndex), XINPUT_DEVTYPE_GAMEPAD, caps_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &caps, syscall.Errno(result)
	}

	return &caps, nil
}

func GetBatteryInformation(userIndex int, deviceSubtype int) (*XINPUT_BATTERY_INFORMATION, error) {
	batt := XINPUT_BATTERY_INFORMATION{}

	if userIndex < 0 || userIndex > 3 {
		return &batt, errors.New("Invalid user index (0-3)")
	}
	if deviceSubtype != BATTERY_DEVTYPE_GAMEPAD && deviceSubtype != BATTERY_DEVTYPE_HEADSET {
		return &batt, errors.New("Invalid device subtype (0, 1)")
	}

	batt_p := uintptr(unsafe.Pointer(&batt))

	result, _, _ := XInputGetBatteryInformation.Call(uintptr(userIndex), uintptr(deviceSubtype), batt_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &batt, syscall.Errno(result)
	}

	return &batt, nil
}

func GetState(userIndex int) (*XINPUT_STATE, error) {
	state := XINPUT_STATE{}

	if userIndex < 0 || userIndex > 3 {
		return &state, errors.New("Invalid user index (0-3)")
	}

	state_p := uintptr(unsafe.Pointer(&state))

	result, _, _ := XInputGetState.Call(uintptr(userIndex), state_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &state, syscall.Errno(result)
	}

	applyDeadzones(&state)

	return &state, nil
}

func applyDeadzones(state *XINPUT_STATE) {
	// left thumbstick X
	if uint8(state.Gamepad.ThumbLX) < uint8(deadzoneLeftThumbInner) {
		state.Gamepad.ThumbLX = 0
	} else if uint8(state.Gamepad.ThumbLX) >= uint8(deadzoneLeftThumbOuter) {
		state.Gamepad.ThumbLX = math.MaxInt16
	}

	// left thumbstick Y
	if uint8(state.Gamepad.ThumbLY) < uint8(deadzoneLeftThumbInner) {
		state.Gamepad.ThumbLY = 0
	} else if uint8(state.Gamepad.ThumbLY) >= uint8(deadzoneLeftThumbOuter) {
		state.Gamepad.ThumbLY = math.MaxInt16
	}

	// right thumbstick X
	if uint8(state.Gamepad.ThumbRX) < uint8(deadzoneRightThumbInner) {
		state.Gamepad.ThumbRX = 0
	} else if uint8(state.Gamepad.ThumbRX) >= uint8(deadzoneRightThumbOuter) {
		state.Gamepad.ThumbRX = math.MaxInt16
	}

	// right thumbstick Y
	if uint8(state.Gamepad.ThumbRY) < uint8(deadzoneRightThumbInner) {
		state.Gamepad.ThumbRY = 0
	} else if uint8(state.Gamepad.ThumbRY) >= uint8(deadzoneRightThumbOuter) {
		state.Gamepad.ThumbRY = math.MaxInt16
	}

	// TODO triggers
}
