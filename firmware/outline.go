//*****************************************************************************
//
// Copyright (c) 2026 Brian Kircher
//
// Open Source Software: you can modify and/or share it under the terms of the
// BSD license file in the root directory of this project.
//
//*****************************************************************************

package main

//*****************************************************************************
//
// Outlines the logo by chasing LEDs around each symbol in order.
//
//*****************************************************************************
func outline(idx uint16) bool {
	//
	// End the animation after 148 steps (~2.37 seconds).
	//
	if idx >= 148 {
		return true
	}

	//
	// No updates are required on the odd steps.
	//
	if (idx & 1) != 0 {
		return false
	}

	//
	// Divide the step by 2 since updates are only made on even steps.
	//
	idx /= 2

	//
	// The second half of the animation matches the first half.
	//
	if idx >= 37 {
		idx -= 37
	}

	//
	// Turn on a red LED when and where appropriate.
	//
	if idx < 11 {
		redLEDs.Set(uint8(idx & 3), uint8(idx))
	} else if idx == 11 {
		redLEDs.Set(3, 0xff)
	} else if idx == 12 {
		redLEDs.Set(0, 0xff)
	} else if idx == 13 {
		redLEDs.Set(1, 0xff)
	} else if idx == 14 {
		redLEDs.Set(2, 0xff)
	}

	//
	// Turn on a white LED when and where appropriate.
	//
	if (idx >= 11) && (idx < 22) {
		whiteLEDs.Set(uint8((idx - 11) & 3), uint8(idx - 11))
	} else if idx == 22 {
		whiteLEDs.Set(3, 0xff)
	} else if idx == 23 {
		whiteLEDs.Set(0, 0xff)
	} else if idx == 24 {
		whiteLEDs.Set(1, 0xff)
	} else if idx == 25 {
		whiteLEDs.Set(2, 0xff)
	}

	//
	// Turn on a blue LED when and where appropriate.
	//
	if (idx >= 22) && (idx < 33) {
		blueLEDs.Set(uint8((idx - 22) & 3), uint8(idx - 22))
	} else if idx == 33 {
		blueLEDs.Set(3, 0xff)
	} else if idx == 34 {
		blueLEDs.Set(0, 0xff)
	} else if idx == 35 {
		blueLEDs.Set(1, 0xff)
	} else if idx == 36 {
		blueLEDs.Set(2, 0xff)
	}

	//
	// The animation is not done.
	//
	return false
}
