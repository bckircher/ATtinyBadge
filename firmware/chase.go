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
// The current chase index.
//
//*****************************************************************************
var chaseIdx uint8

//*****************************************************************************
//
// Chasing LEDs around each symbol.
//
//*****************************************************************************
func chase(idx uint16) bool {
	//
	// End the animation after 128 steps (~2 second).
	//
	if idx >= 128 {
		return true
	}

	//
	// Get the index of the first LED to turn on.
	//
	if idx == 0 {
		chaseIdx = 0
	} else {
		chaseIdx = step3Mod12[chaseIdx]
	}

	//
	// Turn on the first red, white, and blue LEDs.
	//
	redLEDs.Set(0, chaseIdx)
	whiteLEDs.Set(0, chaseIdx)
	blueLEDs.Set(0, chaseIdx)

	//
	// Get the index of the next LED to turn on.
	//
	chaseIdx = step3Mod12[chaseIdx]

	//
	// Turn on the second red, white, and blue LEDs.
	//
	redLEDs.Set(1, chaseIdx)
	whiteLEDs.Set(1, chaseIdx)
	blueLEDs.Set(1, chaseIdx)

	//
	// Get the index of the next LED to turn on.
	//
	chaseIdx = step3Mod12[chaseIdx]

	//
	// Turn on the third red, white, and blue LEDs.
	//
	redLEDs.Set(2, chaseIdx)
	whiteLEDs.Set(2, chaseIdx)
	blueLEDs.Set(2, chaseIdx)

	//
	// Get the index of the next LED to turn on.
	//
	chaseIdx = step3Mod12[chaseIdx]

	//
	// Turn on the fourth red, white, and blue LEDs.
	//
	redLEDs.Set(3, chaseIdx)
	whiteLEDs.Set(3, chaseIdx)
	blueLEDs.Set(3, chaseIdx)

	//
	// The animation is not done.
	//
	return false
}
