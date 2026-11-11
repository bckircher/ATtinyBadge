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
// A test animation for the red LED array.  Lights the LEDs in order forever.
//
//*****************************************************************************
func testRed(idx uint16) bool {
	//
	// Only update the red LED array every 16 steps.
	//
	if (idx & 15) != 0 {
		return false
	}

	//
	// Convert the step into the LED of the array to light.
	//
	idx = (idx / 16) % 16
	if idx > 11 {
		idx = 11
	}

	//
	// Turn on this LED.
	//
	redLEDs.Set(0, uint8(idx))

	//
	// The animation is never done.
	//
	return false
}

//*****************************************************************************
//
// A test animation for the white LED array.  Lights the LEDs in order forever.
//
//*****************************************************************************
func testWhite(idx uint16) bool {
	//
	// Only update the white LED array every 16 steps.
	//
	if (idx & 15) != 0 {
		return false
	}

	//
	// Convert the step into the LED of the array to light.
	//
	idx = (idx / 16) % 16
	if idx > 11 {
		idx = 11
	}

	//
	// Turn on this LED.
	//
	whiteLEDs.Set(0, uint8(idx))

	//
	// The animation is never done.
	//
	return false
}

//*****************************************************************************
//
// A test animation for the blue LED array.  Lights the LEDs in order forever.
//
//*****************************************************************************
func testBlue(idx uint16) bool {
	//
	// Only update the blue LED array every 16 steps.
	//
	if (idx & 15) != 0 {
		return false
	}

	//
	// Convert the step into the LED of the array to light.
	//
	idx = (idx / 16) % 16
	if idx > 11 {
		idx = 11
	}

	//
	// Turn on this LED.
	//
	blueLEDs.Set(0, uint8(idx))

	//
	// The animation is never done.
	//
	return false
}
