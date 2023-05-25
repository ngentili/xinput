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

	api *XInputApi
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
	if proc, err := dll.FindProc("XInputEnable"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetAudioDeviceIds"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetDSoundAudioDeviceGuids"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetBatteryInformation"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetCapabilities"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetKeystroke"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputGetState"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProc("XInputSetState"); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(100); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(101); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(102); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(103); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(104); err == nil {
		api.XInputEnable = proc
	}
	if proc, err := dll.FindProcByOrdinal(108); err == nil {
		api.XInputEnable = proc
	}

	fmt.Println("Available procedures:")
	fmt.Println("XInputEnable: ", api.XInputEnable != nil)
	fmt.Println("XInputGetAudioDeviceIds: ", api.XInputGetAudioDeviceIds != nil)
	fmt.Println("XInputGetDSoundAudioDeviceGuids: ", api.XInputGetDSoundAudioDeviceGuids != nil)
	fmt.Println("XInputGetBatteryInformation: ", api.XInputGetBatteryInformation != nil)
	fmt.Println("XInputGetCapabilities: ", api.XInputGetCapabilities != nil)
	fmt.Println("XInputGetKeystroke: ", api.XInputGetKeystroke != nil)
	fmt.Println("XInputGetState: ", api.XInputGetState != nil)
	fmt.Println("XInputSetState: ", api.XInputSetState != nil)
	fmt.Println("100: ", api.XInputOrdinal100 != nil)
	fmt.Println("101: ", api.XInputOrdinal101 != nil)
	fmt.Println("102: ", api.XInputOrdinal102 != nil)
	fmt.Println("103: ", api.XInputOrdinal103 != nil)
	fmt.Println("104: ", api.XInputOrdinal104 != nil)
	fmt.Println("108: ", api.XInputOrdinal108 != nil)

	connected := 0
	for i := 0; i < 4; i++ {
		if _, err := GetCapabilities(i); err == nil {
			connected += 1
		}
	}

	fmt.Println("Devices connected: ", connected)
}

func GetCapabilities(userIndex int) (*XINPUT_CAPABILITIES, error) {
	if api.XInputGetCapabilities != nil {
		panic("Procedure not suppported in current XInput version")
	}

	caps := XINPUT_CAPABILITIES{}

	if userIndex < 0 || userIndex > 3 {
		return &caps, errors.New("Invalid user index (0-3)")
	}

	caps_p := uintptr(unsafe.Pointer(&caps))

	result, _, _ := api.XInputGetCapabilities.Call(uintptr(userIndex), XINPUT_DEVTYPE_GAMEPAD, caps_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &caps, syscall.Errno(result)
	}

	return &caps, nil
}

func GetBatteryInformation(userIndex int, deviceSubtype int) (*XINPUT_BATTERY_INFORMATION, error) {
	if api.XInputGetBatteryInformation != nil {
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

	result, _, _ := api.XInputGetBatteryInformation.Call(uintptr(userIndex), uintptr(deviceSubtype), batt_p)

	if result != uintptr(windows.ERROR_SUCCESS) {
		return &batt, syscall.Errno(result)
	}

	return &batt, nil
}

func GetState(userIndex int) (bool, *XINPUT_STATE) {
	if api.XInputGetState != nil {
		panic("Procedure not suppported in current XInput version")
	}

	state := XINPUT_STATE{}

	result, _, _ := api.XInputGetState.Call(uintptr(userIndex), uintptr(unsafe.Pointer(&state)))

	if result == uintptr(windows.ERROR_SUCCESS) {
		return true, &state

	} else if result != uintptr(windows.ERROR_DEVICE_NOT_CONNECTED) {
		return false, &state

	} else {
		panic(syscall.Errno(result))
	}
}

func adjustForDeadzone(x, y, dz float64) (float64, float64) {
	// https://learn.microsoft.com/en-us/windows/win32/xinput/getting-started-with-xinput#dead-zone

	magnitude := math.Sqrt(x*x + y*y)

	normalizedX := x / magnitude
	normalizedY := y / magnitude

	// normalizedMagnitude := float64(0)

	// check if the controller is outside a circular dead zone
	if magnitude > dz {

		// clip the magnitude at its expected maximum value
		if magnitude > 32767 {
			magnitude = 32767
		}

		// adjust magnitude relative to the end of the dead zone
		magnitude -= dz

		// optionally normalize the magnitude with respect to its expected range
		// giving a magnitude value of 0.0 to 1.0
		// normalizedMagnitude = magnitude / (32767 - dz)

	} else {
		magnitude = 0
		// normalizedMagnitude = 0
	}

	adjustedX := normalizedX * magnitude
	adjustedY := normalizedY * magnitude

	return adjustedX, adjustedY
}
