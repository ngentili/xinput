package xinput

import "math"

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
