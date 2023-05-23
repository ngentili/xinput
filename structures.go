package xinput

// Xinput.h

// wintypes
// type (
// 	BYTE  uint8
// 	WORD  uint16
// 	DWORD uint32
// 	SHORT int16
// 	WCHAR int16
// )

//
// Structures used by XInput APIs
//
type XINPUT_GAMEPAD struct {
	Buttons      uint16 // WORD
	LeftTrigger  uint8  // BYTE
	RightTrigger uint8  // BYTE
	ThumbLX      int16  // SHORT
	ThumbLY      int16  // SHORT
	ThumbRX      int16  // SHORT
	ThumbRY      int16  // SHORT
}

type XINPUT_STATE struct {
	PacketNumber uint32 // DWORD
	Gamepad        XINPUT_GAMEPAD
}

type XINPUT_VIBRATION struct {
	LeftMotorSpeed  uint16 // WORD
	RightMotorSpeed uint16 // WORD
}

type XINPUT_CAPABILITIES struct {
	Type      uint8  // BYTE
	SubType   uint8  // BYTE
	Flags     uint16 // WORD
	Gamepad   XINPUT_GAMEPAD
	Vibration XINPUT_VIBRATION
}

type XINPUT_BATTERY_INFORMATION struct {
	BatteryType  uint8 // BYTE
	BatteryLevel uint8 // BYTE
}

type XINPUT_KEYSTROKE struct {
	VirtualKey uint16 // WORD
	Unicode    int16  // WCHAR
	Flags      uint16 // WORD
	UserIndex  uint8  // BYTE
	HidCode    uint8  // BYTE
}
