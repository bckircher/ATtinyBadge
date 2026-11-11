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
// The number of consecutive ticks the button must in the opposite state in
// order to be considered to be the opposite state (in other words, it has been
// debounced).
//
//*****************************************************************************
const DEBOUNCE_TIME uint8 = 3

//*****************************************************************************
//
// The number of consecutive ticks the debounced button must be pressed in
// order to consider it a hold.
//
//*****************************************************************************
const HOLD_TIME uint8 = 125

//*****************************************************************************
//
// The current raw debounced state of the button.
//
//*****************************************************************************
var btnState1 uint8 = 1

//*****************************************************************************
//
// The current press/hold state of the button.
//
//*****************************************************************************
var btnState2 uint8 = 1

//*****************************************************************************
//
// The count of consecutive ticks the button has been in the opposite state.
//
//*****************************************************************************
var btnDebounceCount uint8 = 0

//*****************************************************************************
//
// The count of consecutive ticks the debounced button has been pressed.
//
//*****************************************************************************
var btnHoldCount uint8 = 0

//*****************************************************************************
//
// The flags that indicate detected button events.
//
//*****************************************************************************
var btnFlags uint8 = 0

//*****************************************************************************
//
// The bit in btnFlags that indicates that a button press has been detected.
//
//*****************************************************************************
const BTN_PRESS uint8 = 1

//*****************************************************************************
//
// The bit in btnFlags that indicates that a button hold has been detected.
//
//*****************************************************************************
const BTN_HOLD uint8 = 2

//*****************************************************************************
//
// Debounces the state of the button, detecting presses versus holds while
// rejecting spurious noise and button bouncing.
//
//*****************************************************************************
func debounce(newState uint8) {
	//
	// See if the button state has changed.
	//
	if newState == btnState1 {
		//
		// It has not changed, so reset the press/release counter.
		//
		btnDebounceCount = 0
	} else {
		//
		// It has changed, so increment the press/release counter.
		//
		btnDebounceCount++

		//
		// See if the debounce time has been reached.
		//
		if btnDebounceCount >= DEBOUNCE_TIME {
			//
			// Accept the button state as the debounced button state.
			//
			btnState1 = newState

			//
			// Reset the press/release counter.
			//
			btnDebounceCount = 0
		}
	}

	//
	// There is nothing to do if the debounced pin states are both not pressed.
	//
	if (btnState1 != 0) && (btnState2 != 0) {
		return
	}

	//
	// See if there has been rising edge on the debounced button state.
	//
	if (btnState1 != 0) && (btnState2 == 0) {
		//
		// If the button was held for less than the hold time, a button press
		// has been detected.
		//
		if btnHoldCount < HOLD_TIME {
			btnFlags |= BTN_PRESS
		}

		//
		// Update the debounced button state tracker.
		//
		btnState2 = btnState1
		btnHoldCount = 0

		//
		// There is nothing further to do.
		//
		return
	}

	//
	// See if there has been a falling edge on the debounced button state.
	//
	if (btnState1 == 0) && (btnState2 != 0) {
		//
		// Update the debounced button state tracker.
		//
		btnState2 = btnState1
		btnHoldCount = DEBOUNCE_TIME

		//
		// There is nothing further to do.
		//
		return
	}

	//
	// The button is being held, so increment the hold counter, if it is not
	// at its maximum value.
	//
	if btnHoldCount < 255 {
		btnHoldCount++
	}

	//
	// If the button hold time has been reached, indicate a button hold.
	//
	if btnHoldCount == HOLD_TIME {
	   btnFlags |= BTN_HOLD
	}
}

//*****************************************************************************
//
// Determines if a button press has been detected.  This returns true exactly
// one time per button press; a subsequent button release and press must occur
// before it returns true again.
//
//*****************************************************************************
func wasPressed() bool {
	//
	// Get the state of the button press flag.
	//
	pressed := (btnFlags & BTN_PRESS) != 0

	//
	// Clear the button press flag if it was set.
	//
	if pressed {
		btnFlags &= ^BTN_PRESS
	}

	//
	// Return the state of the button press flag.
	//
	return pressed
}

//*****************************************************************************
//
// Determines if a button hold has been detected.  This returns true exactly
// one time per button hold; a subsequent button release and hold must occur
// before it returns true again.
//
//*****************************************************************************
func wasHeld() bool {
	//
	// Get the state of the button hold flag.
	//
	held := (btnFlags & BTN_HOLD) != 0

	//
	// Clear the button hold flag if it was set.
	//
	if held {
		btnFlags &= ^BTN_HOLD
	}

	//
	// Return the state of the button hold flag.
	//
	return held
}
