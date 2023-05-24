package xinput

import (
	"errors"
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	versions = [...]string{
		"XInput1_4.dll",
		"XInput9_1_0.dll",
		"xinput1_3.dll",
		"xinput1_2.dll",
		"xinput1_1.dll",
		"XInputUap.dll",
	}

	dll *windows.DLL

	// procedures
	xInputEnable                    *windows.Proc
	xInputGetAudioDeviceIds         *windows.Proc
	xInputGetDSoundAudioDeviceGuids *windows.Proc
	xInputGetBatteryInformation     *windows.Proc
	xInputGetCapabilities           *windows.Proc
	xInputGetKeystroke              *windows.Proc
	xInputGetState                  *windows.Proc
	xInputSetState                  *windows.Proc
	xInputOrdinal100                *windows.Proc
	xInputOrdinal101                *windows.Proc
	xInputOrdinal102                *windows.Proc
	xInputOrdinal103                *windows.Proc
	xInputOrdinal104                *windows.Proc
	xInputOrdinal108                *windows.Proc

	// deadzones
	deadzoneLeftThumbInner  int
	deadzoneLeftThumbOuter  int
	deadzoneRightThumbInner int
	deadzoneRightThumbOuter int

	// version compatibility
	support_XInputEnable                    bool
	support_XInputGetAudioDeviceIds         bool
	support_XInputGetDSoundAudioDeviceGuids bool
	support_XInputGetBatteryInformation     bool
	support_XInputGetCapabilities           bool
	support_XInputGetKeystroke              bool
	support_XInputGetState                  bool
	support_XInputSetState                  bool
	support_XInputOrdinal100                bool
	support_XInputOrdinal101                bool
	support_XInputOrdinal102                bool
	support_XInputOrdinal103                bool
	support_XInputOrdinal104                bool
	support_XInputOrdinal108                bool
)

func init() {
	// load dll
	for _, filename := range versions {
		d, err := windows.LoadDLL(filename)
		if err != nil {
			continue
		} else {
			dll = d
			fmt.Println("Loaded XInput DLL version: ", filename)
			break
		}
	}
	if dll == nil {
		panic("Could not load XInput DLL")
	}

	// load procedures and determine compatibility
	var err error

	xInputEnable, err = dll.FindProc("XInputEnable")
	support_XInputEnable = err == nil

	xInputGetAudioDeviceIds, err = dll.FindProc("XInputGetAudioDeviceIds")
	support_XInputGetAudioDeviceIds = err == nil

	xInputGetDSoundAudioDeviceGuids, err = dll.FindProc("XInputGetDSoundAudioDeviceGuids")
	support_XInputGetDSoundAudioDeviceGuids = err == nil

	xInputGetBatteryInformation, err = dll.FindProc("XInputGetBatteryInformation")
	support_XInputGetBatteryInformation = err == nil

	xInputGetCapabilities, err = dll.FindProc("XInputGetCapabilities")
	support_XInputGetCapabilities = err == nil

	xInputGetKeystroke, err = dll.FindProc("XInputGetKeystroke")
	support_XInputGetKeystroke = err == nil

	xInputGetState, err = dll.FindProc("XInputGetState")
	support_XInputGetState = err == nil

	xInputSetState, err = dll.FindProc("XInputSetState")
	support_XInputSetState = err == nil

	xInputOrdinal100, err = dll.FindProcByOrdinal(100)
	support_XInputOrdinal100 = err == nil

	xInputOrdinal101, err = dll.FindProcByOrdinal(101)
	support_XInputOrdinal101 = err == nil

	xInputOrdinal102, err = dll.FindProcByOrdinal(102)
	support_XInputOrdinal102 = err == nil

	xInputOrdinal103, err = dll.FindProcByOrdinal(103)
	support_XInputOrdinal103 = err == nil

	xInputOrdinal104, err = dll.FindProcByOrdinal(104)
	support_XInputOrdinal104 = err == nil

	xInputOrdinal108, err = dll.FindProcByOrdinal(108)
	support_XInputOrdinal108 = err == nil

	fmt.Println("Available procedures:")
	fmt.Println("XInputEnable: ", support_XInputEnable)
	fmt.Println("XInputGetAudioDeviceIds: ", support_XInputGetAudioDeviceIds)
	fmt.Println("XInputGetDSoundAudioDeviceGuids: ", support_XInputGetDSoundAudioDeviceGuids)
	fmt.Println("XInputGetBatteryInformation: ", support_XInputGetBatteryInformation)
	fmt.Println("XInputGetCapabilities: ", support_XInputGetCapabilities)
	fmt.Println("XInputGetKeystroke: ", support_XInputGetKeystroke)
	fmt.Println("XInputGetState: ", support_XInputGetState)
	fmt.Println("XInputSetState: ", support_XInputSetState)
	fmt.Println("100: ", support_XInputOrdinal100)
	fmt.Println("101: ", support_XInputOrdinal101)
	fmt.Println("102: ", support_XInputOrdinal102)
	fmt.Println("103: ", support_XInputOrdinal103)
	fmt.Println("104: ", support_XInputOrdinal104)
	fmt.Println("108: ", support_XInputOrdinal108)

	// set deadzones
	deadzoneLeftThumbInner = XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE
	deadzoneLeftThumbOuter = math.MaxInt16
	deadzoneRightThumbInner = XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE
	deadzoneRightThumbOuter = math.MaxInt16

	// deadzoneT = XINPUT_GAMEPAD_TRIGGER_THRESHOLD

	fmt.Println("Deadzones:")
	fmt.Println("Left thumbstick inner:  ", deadzoneLeftThumbInner)
	fmt.Println("Left thumbstick outer:  ", deadzoneLeftThumbOuter)
	fmt.Println("Right thumbstick inner: ", deadzoneRightThumbInner)
	fmt.Println("Right thumbstick outer: ", deadzoneRightThumbOuter)

	connected := 0
	for i := 0; i < 4; i++ {
		if _, err := GetCapabilities(i); err == nil {
			connected += 1
		}
	}

	fmt.Println("Devices connected: ", connected)
}

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
	if !support_XInputGetCapabilities {
		panic("Procedure not suppported in current XInput version")
	}

	caps := XINPUT_CAPABILITIES{}

	if userIndex < 0 || userIndex > 3 {
		return &caps, errors.New("Invalid user index (0-3)")
	}

	caps_p := uintptr(unsafe.Pointer(&caps))

	result, _, _ := xInputGetCapabilities.Call(uintptr(userIndex), XINPUT_DEVTYPE_GAMEPAD, caps_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &caps, syscall.Errno(result)
	}

	return &caps, nil
}

func GetBatteryInformation(userIndex int, deviceSubtype int) (*XINPUT_BATTERY_INFORMATION, error) {
	if !support_XInputGetBatteryInformation {
		panic("Procedure not suppported in current XInput version")
	}

	batt := XINPUT_BATTERY_INFORMATION{}

	if userIndex < 0 || userIndex > 3 {
		return &batt, errors.New("Invalid user index (0-3)")
	}
	if deviceSubtype != BATTERY_DEVTYPE_GAMEPAD && deviceSubtype != BATTERY_DEVTYPE_HEADSET {
		return &batt, errors.New("Invalid device subtype (0, 1)")
	}

	batt_p := uintptr(unsafe.Pointer(&batt))

	result, _, _ := xInputGetBatteryInformation.Call(uintptr(userIndex), uintptr(deviceSubtype), batt_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &batt, syscall.Errno(result)
	}

	return &batt, nil
}

func GetState(userIndex int) (*XINPUT_STATE, error) {
	if !support_XInputGetState {
		panic("Procedure not suppported in current XInput version")
	}

	state := XINPUT_STATE{}

	if userIndex < 0 || userIndex > 3 {
		return &state, errors.New("Invalid user index (0-3)")
	}

	state_p := uintptr(unsafe.Pointer(&state))

	result, _, _ := xInputGetState.Call(uintptr(userIndex), state_p)

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
