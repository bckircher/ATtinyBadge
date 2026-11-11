//*****************************************************************************
//
// Copyright (c) 2026 Brian Kircher
//
// Open Source Software: you can modify and/or share it under the terms of the
// BSD license file in the root directory of this project.
//
//*****************************************************************************

//go:build baremetal

package main

import (
	"device/avr"
)

//*****************************************************************************
//
// The number of LEDs that can be driven in the LED array.
//
//*****************************************************************************
const numLEDs = 4

//*****************************************************************************
//
// The definition for an array of LEDs.
//
//*****************************************************************************
type LEDs struct {
	//
	// The ATtiny port that drives these LEDs.
	//
	port    *avr.PORT_Type

	//
	// The bit mask for the pins in the port that drive these LEDs.
	//
	mask    uint8

	//
	// The mapping of LED index and the pin "address".
	//
	mapping []uint8

	//
	// The bits in the direction register to clear for each LED.
	//
	dirClr  [numLEDs]uint8

	//
	// The bits in the direction register to set for each LED.
	//
	dirSet  [numLEDs]uint8

	//
	// The bits in the output register to clear for each LED.
	//
	outClr  [numLEDs]uint8

	//
	// The bits in the output register to set for each LED.
	//
	outSet  [numLEDs]uint8
}

//*****************************************************************************
//
// Sets an LED to be driven.  The given slot in the drive array is immediately
// set, causing the LED to light up as soon as the givne slot is updated.  For
// execution speed, the slot index is assumed to be valid; bad things happen if
// it is not!
//
//*****************************************************************************
func (l *LEDs) Set(idx uint8, led uint8) {
	//
	// Map the LED to its address.
	//
	if led != 0xff {
		led = l.mapping[led]
	}

	//
	// Get the bit number to drive high and low from the LED number.
	//
	high := led / 16
	low := led % 16

	//
	// Compute the bits in the direction and output registers that should be
	// set to properly drive the requested LED.
	//
	dir := uint8((1 << high) | (1 << low))
	out := uint8(1 << high)

	//
	// Set this array entry to drive the appropriate pin, and tri-state the
	// others, to properly drive only the requested LED.
	//
	l.dirClr[idx] = dir ^ l.mask
	l.dirSet[idx] = dir
	l.outClr[idx] = out ^ l.mask
	l.outSet[idx] = out
}

//*****************************************************************************
//
// Drives the pins to light up the LED in the given index (or none if one was
// not enabled for this index).
//
//*****************************************************************************
func (l *LEDs) Refresh(idx int) {
	//
	// Set the direction and output registers as appropriate.  Note that bits
	// are cleared in both first, then other bits set, to avoid briefly driving
	// unexpected LEDs (which shows up as a very dim ghosting).
	//
	l.port.OUTCLR.Set(l.outClr[idx])
	l.port.DIRCLR.Set(l.dirClr[idx])
	l.port.DIRSET.Set(l.dirSet[idx])
	l.port.OUTSET.Set(l.outSet[idx])
}

//*****************************************************************************
//
// Immediately forces all of the LEDs to be off.
//
//*****************************************************************************
func (l *LEDs) ForceOff() {
	//
	// Loop through the register arrays for the LEDs, setting each value to
	// the value that results in the entire array being tri-stated.
	//
	for idx := 0; idx < numLEDs; idx++ {
		l.dirClr[idx] = l.mask
		l.dirSet[idx] = 0
		l.outClr[idx] = l.mask
		l.outSet[idx] = 0
	}

	//
	// Set the direction and output registers to turn off all the LEDs.
	//
	l.port.OUTCLR.Set(l.mask)
	l.port.DIRCLR.Set(l.mask)
}
