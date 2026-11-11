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
// The current logo index.
//
//*****************************************************************************
var logoIdx uint8

//*****************************************************************************
//
// Displays each symbol in the logo in sequence.
//
//*****************************************************************************
func logo(idx uint16) bool {
	//
	// End the animation after 132 steps (~2.1 seconds).
	//
	if idx >= 132 {
		return true
	}

	//
	// Get the index of the first LED to turn on.
	//
	if idx == 0 {
		logoIdx = 0
	} else {
		logoIdx = step3Mod12[logoIdx]
	}

	//
	// Turn on an LED based on the step.
	//
	if idx < 44 {
		//
		// For the first 44 steps (~0.7 seconds), turn on a new set of red
		// LEDs.
		//
		redLEDs.Set(0, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		redLEDs.Set(1, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		redLEDs.Set(2, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		redLEDs.Set(3, logoIdx)
	} else if idx < 88 {
		//
		// For the second 44 steps (~0.7 seconds), turn on a new white LED.
		//
		whiteLEDs.Set(0, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		whiteLEDs.Set(1, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		whiteLEDs.Set(2, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		whiteLEDs.Set(3, logoIdx)
	} else {
		//
		// For the final 44 steps (~0.7 seconds), turn on a new blue LED.
		//
		blueLEDs.Set(0, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		blueLEDs.Set(1, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		blueLEDs.Set(2, logoIdx)
		logoIdx = step3Mod12[logoIdx]
		blueLEDs.Set(3, logoIdx)
	}

	//
	// At the 44th step, when the first white LED is turned on, turn off all
	// the red LEDs.
	//
	if idx == 44 {
		redLEDs.ForceOff()
	}

	//
	// At the 88th step, when the first blue LED is turned on, turn off all the
	// white LEDs.
	//
	if idx == 88 {
		whiteLEDs.ForceOff()
	}

	//
	// The animation is not done.
	//
	return false
}
