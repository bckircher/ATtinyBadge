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
// A quick animation to indicate the badge powering off.
//
//*****************************************************************************
func off(idx uint16) bool {
	//
	// End the animation after 64 steps (~1 second).
	//
	if idx >= 64 {
		return true
	}

	//
	// Update the animation every 16 steps.
	//
	if (idx & 15) != 0 {
		return false
	}

	//
	// Convert the step number into a value between 0 and 3, inclusive.
	//
	idx = idx >> 4

	//
	// See if this is the first step.
	//
	if idx == 0 {
		//
		// Turn on four LEDs in each array.
		//
		redLEDs.Set(0, 2)
		redLEDs.Set(1, 3)
		redLEDs.Set(2, 4)
		redLEDs.Set(3, 5)
		whiteLEDs.Set(0, 1)
		whiteLEDs.Set(1, 0)
		whiteLEDs.Set(2, 2)
		whiteLEDs.Set(3, 10)
		blueLEDs.Set(0, 10)
		blueLEDs.Set(1, 9)
		blueLEDs.Set(2, 8)
		blueLEDs.Set(3, 7)
	} else {
		//
		// Turn off the next LED in each array.
		//
		redLEDs.Set(uint8(4 - idx), 11)
		whiteLEDs.Set(uint8(4 - idx), 11)
		blueLEDs.Set(uint8(4 - idx), 11)
	}

	//
	// The animation is not done.
	//
	return false
}
