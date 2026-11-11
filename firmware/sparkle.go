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
// Displays random "sparkles" on the three LED arrays.  Up to four LEDs are on
// at any time in each array (on occasion, the random numbers may result in the
// same LED for more than one entry in an array, which results in it being on
// for longer for that step, making it a bit brighter).
//
//*****************************************************************************
func sparkle(idx uint16) bool {
	//
	// End the animation after 312 steps (~5 second).
	//
	if idx >= 312 {
		return true
	}

	//
	// Only update the animation every other time.
	//
	if (idx & 1) != 0 {
		return false
	}

	//
	// Divide the step by two since this is run every other time.
	//
	idx /= 2

	//
	// Select a random red LED for the index based on the step number.
	//
	redLEDs.Set(uint8(idx & 3), uint8(randomLessThan8(11)))

	//
	// Select a random white LED for the index based on the step number.
	//
	whiteLEDs.Set(uint8(idx & 3), uint8(randomLessThan8(11)))

	//
	// Select a random blue LED for the index based on the step number.
	//
	blueLEDs.Set(uint8(idx & 3), uint8(randomLessThan8(11)))

	//
	// The animation is not done.
	//
	return false
}
