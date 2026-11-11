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
// Constants to help build the state table.
//
//*****************************************************************************
const RED   uint8 = 0 << 6
const WHITE uint8 = 1 << 6
const BLUE  uint8 = 2 << 6
const SLOT0 uint8 = 0 << 4
const SLOT1 uint8 = 1 << 4
const SLOT2 uint8 = 2 << 4
const SLOT3 uint8 = 3 << 4
const NOP   uint8 = 0xff
const OFF   uint8 = 15

//*****************************************************************************
//
// The states, representing the changes to the LED arrays to be made at each
// step of the horizontal sweep.  This lights up the LEDs in order from left to
// right, then back the other way.  The horizontal positions of the LEDs are as
// follows:
//
//     R0  -> 139.6875    W0  -> 154.9300    B0  -> 155.7500
//     R1  -> 126.3625    W1  -> 147.3000    B1  -> 150.1500
//     R2  -> 119.7000    W2  -> 139.7700    B2  -> 155.5500
//     R3  -> 121.1625    W3  -> 133.3300    B3  -> 160.9500
//     R4  -> 122.6250    W4  -> 130.0400    B4  -> 166.3500
//     R5  -> 124.0875    W5  -> 138.0000    B5  -> 171.9700
//     R6  -> 125.5500    W6  -> 146.2000    B6  -> 177.5800
//     R7  -> 130.7500    W7  -> 154.6000    B7  -> 183.2000
//     R8  -> 135.9500    W8  -> 165.4800    B8  -> 177.7800
//     R9  -> 141.1500    W9  -> 166.2700    B9  -> 172.3700
//     R10 -> 146.3500    W10 -> 161.8200    B10 -> 166.9500
//
//*****************************************************************************
var sweepStatesH = []uint8 {
	RED   | SLOT0 | 2,   NOP,                 // 0... .... ....
	RED   | SLOT1 | 3,   NOP,                 // 10.. .... ....
	RED   | SLOT2 | 4,   NOP,                 // 210. .... ....
	RED   | SLOT3 | 5,   NOP,                 // 3210 .... ....
	RED   | SLOT0 | 6,   NOP,                 // 0321 .... ....
	RED   | SLOT1 | 1,   NOP,                 // 1032 .... ....
	RED   | SLOT2 | OFF, WHITE | SLOT0 | 4,   // 21.3 0... ....
	RED   | SLOT3 | 7,   NOP,                 // 32.0 1... ....
	RED   | SLOT0 | OFF, WHITE | SLOT1 | 3,   // .3.1 20.. ....
	RED   | SLOT1 | 8,   NOP,                 // .0.2 31.. ....
	RED   | SLOT3 | OFF, WHITE | SLOT0 | 5,   // .1.3 02.. ....
	RED   | SLOT3 | 0,   NOP,                 // .2.0 13.. ....
	WHITE | SLOT1 | 2,   NOP,                 // .3.1 20.. ....
	RED   | SLOT1 | 9,   NOP,                 // .0.2 31.. ....
	WHITE | SLOT0 | 6,   NOP,                 // .1.3 02.. ....
	RED   | SLOT3 | 10,  NOP,                 // .2.0 13.. ....
	WHITE | SLOT1 | 1,   NOP,                 // .3.1 20.. ....
	RED   | SLOT1 | OFF, BLUE  | SLOT0 | 1,   // ...2 31.. 0...
	WHITE | SLOT0 | 7,   NOP,                 // ...3 02.. 1...
	RED   | SLOT3 | OFF, WHITE | SLOT2 | 0,   // .... 130. 2...
	WHITE | SLOT1 | OFF, BLUE  | SLOT1 | 2,   // .... 1..2 30..
	BLUE  | SLOT0 | 0,   NOP,                 // .... 2..3 01..
	WHITE | SLOT3 | OFF, BLUE  | SLOT2 | 3,   // .... 3... 120.
	WHITE | SLOT0 | 10,  NOP,                 // .... 0... 231.
	WHITE | SLOT1 | 8,   BLUE  | SLOT1 | OFF, // .... 10.. 3.2.
	WHITE | SLOT2 | 9,   BLUE  | SLOT0 | OFF, // .... 210. ..3.
	BLUE  | SLOT2 | 4,   NOP,                 // .... 321. ..0.
	WHITE | SLOT0 | OFF, BLUE  | SLOT0 | 10,  // .... .32. 0.1.
	WHITE | SLOT1 | OFF, BLUE  | SLOT1 | 5,   // .... ..3. 102.
	WHITE | SLOT2 | OFF, BLUE  | SLOT3 | 9,   // .... .... 2130
	BLUE  | SLOT2 | 6,   NOP,                 // .... .... 3201
	BLUE  | SLOT0 | 8,   NOP,                 // .... .... 0312
	BLUE  | SLOT1 | 7,   NOP,                 // .... .... 1023
	BLUE  | SLOT3 | OFF, NOP,                 // .... .... 213.
	BLUE  | SLOT2 | OFF, NOP,                 // .... .... 32..
	BLUE  | SLOT0 | OFF, NOP,                 // .... .... .3..
	BLUE  | SLOT1 | OFF, NOP,                 // .... .... ....
	BLUE  | SLOT0 | 7,   NOP,                 // .... .... 0...
	BLUE  | SLOT1 | 8,   NOP,                 // .... .... 10..
	BLUE  | SLOT2 | 6,   NOP,                 // .... .... 210.
	BLUE  | SLOT3 | 9,   NOP,                 // .... .... 3210
	BLUE  | SLOT0 | 5,   NOP,                 // .... .... 0321
	BLUE  | SLOT1 | 10,  NOP,                 // .... .... 1032
	BLUE  | SLOT2 | 4,   NOP,                 // .... .... 2103
	WHITE | SLOT0 | 9,   BLUE  | SLOT3 | OFF, // .... 0... 321.
	WHITE | SLOT1 | 8,   BLUE  | SLOT0 | OFF, // .... 10.. .32.
	WHITE | SLOT2 | 10,  BLUE  | SLOT1 | OFF, // .... 210. ..3.
	BLUE  | SLOT2 | 3,   NOP,                 // .... 321. ..0.
	WHITE | SLOT0 | OFF, BLUE  | SLOT0 | 0,   // .... .32. 0.1.
	WHITE | SLOT1 | OFF, BLUE  | SLOT1 | 2,   // .... ..3. 102.
	WHITE | SLOT2 | 0,   NOP,                 // .... ..0. 213.
	WHITE | SLOT0 | 7,   BLUE  | SLOT2 | OFF, // .... 0.1. 32..
	BLUE  | SLOT0 | 1,   NOP,                 // .... 1.2. 03..
	WHITE | SLOT1 | 1,   BLUE  | SLOT1 | OFF, // .... 203. 1...
	RED   | SLOT0 | 10,  WHITE | SLOT2 | OFF, // 0... 31.. 2...
	WHITE | SLOT0 | 6,   NOP,                 // 1... 02.. 3...
	RED   | SLOT1 | 9,   BLUE  | SLOT0 | OFF, // 20.. 13.. ....
	WHITE | SLOT1 | 2,   NOP,                 // 31.. 20.. ....
	RED   | SLOT0 | 0,   NOP,                 // 02.. 31.. ....
	WHITE | SLOT0 | 5,   NOP,                 // 13.. 02.. ....
	RED   | SLOT1 | 8,   NOP,                 // 20.. 13.. ....
	WHITE | SLOT1 | 3,   NOP,                 // 31.. 20.. ....
	RED   | SLOT0 | 7,   NOP,                 // 02.. 31.. ....
	WHITE | SLOT0 | 4,   NOP,                 // 13.. 02.. ....
	RED   | SLOT1 | 1,   NOP,                 // 20.. 13.. ....
	RED   | SLOT2 | 6,   WHITE | SLOT1 | OFF, // 310. 2... ....
	RED   | SLOT0 | 5,   NOP,                 // 021. 3... ....
	RED   | SLOT3 | 4,   WHITE | SLOT0 | OFF, // 1320 .... ....
	RED   | SLOT1 | 3,   NOP,                 // 2031 .... ....
	RED   | SLOT2 | 2,   NOP,                 // 3102 .... ....
	RED   | SLOT0 | OFF, NOP,                 // .213 .... ....
	RED   | SLOT3 | OFF, NOP,                 // .32. .... ....
	RED   | SLOT1 | OFF, NOP,                 // ..3. .... ....
	RED   | SLOT2 | OFF, NOP,                 // .... .... ....
}

//*****************************************************************************
//
// The states, representing the changes to the LED arrays to be made at each
// step of the vertical sweep.  This lights up the LEDS in order from top to
// bottom, then back the other way.  The vertical positions of the LEDs are as
// follows:
//
//     R0  -> 106.1000    W0  -> 111.9000    B0  -> 104.3700
//     R1  -> 112.2000    W1  -> 112.7700    B1  ->  99.0000
//     R2  -> 115.2500    W2  -> 111.2300    B2  ->  93.3700
//     R3  -> 107.3000    W3  -> 107.0600    B3  ->  87.7300
//     R4  ->  99.3500    W4  -> 100.1100    B4  ->  82.1000
//     R5  ->  91.4000    W5  ->  88.2000    B5  ->  87.4700
//     R6  ->  83.4500    W6  ->  86.0000    B6  ->  92.8300
//     R7  ->  88.3500    W7  ->  86.7000    B7  ->  98.2000
//     R8  ->  93.2500    W8  ->  94.5600    B8  -> 103.8300
//     R9  ->  98.1500    W9  -> 102.2000    B9  -> 109.4700
//     R10 -> 103.0500    W10 -> 108.4700    B10 -> 115.1000
//
//*****************************************************************************
var sweepStatesV = []uint8 {
	BLUE  | SLOT0 | 4,   NOP,                 // .... .... 0...
	RED   | SLOT0 | 6,   NOP,                 // 0... .... 1...
	WHITE | SLOT0 | 6,   NOP,                 // 1... 0... 2...
	WHITE | SLOT1 | 7,   NOP,                 // 2... 10.. 3...
	BLUE  | SLOT0 | 5,   NOP,                 // 3... 21.. 0...
	BLUE  | SLOT1 | 3,   RED   | SLOT0 | OFF, // .... 32.. 10..
	WHITE | SLOT0 | 5,   NOP,                 // .... 03.. 21..
	RED   | SLOT0 | 7,   WHITE | SLOT1 | OFF, // 0... 1... 32..
	RED   | SLOT1 | 5,   BLUE  | SLOT0 | OFF, // 10.. 2... .3..
	BLUE  | SLOT1 | 6,   NOP,                 // 21.. 3... .0..
	RED   | SLOT2 | 8,   WHITE | SLOT0 | OFF, // 320. .... .1..
	BLUE  | SLOT0 | 2,   RED   | SLOT0 | OFF, // .31. .... 02..
	WHITE | SLOT0 | 8,   RED   | SLOT1 | OFF, // ..2. 0... 13..
	RED   | SLOT0 | 9,   BLUE  | SLOT1 | OFF, // 0.3. 1... 2...
	BLUE  | SLOT1 | 7,   RED   | SLOT2 | OFF, // 1... 2... 30..
	BLUE  | SLOT0 | 1,   NOP,                 // 2... 3... 01..
	RED   | SLOT1 | 4,   WHITE | SLOT0 | OFF, // 30.. .... 12..
	WHITE | SLOT0 | 4,   RED   | SLOT0 | OFF, // .1.. 0... 23..
	WHITE | SLOT1 | 9,   BLUE  | SLOT1 | OFF, // .2.. 10.. 3...
	RED   | SLOT0 | 10,  BLUE  | SLOT0 | OFF, // 03.. 21.. ....
	BLUE  | SLOT0 | 8,   RED   | SLOT1 | OFF, // 1... 32.. 0...
	BLUE  | SLOT1 | 0,   WHITE | SLOT0 | OFF, // 2... .3.. 10..
	RED   | SLOT1 | 0,   WHITE | SLOT1 | OFF, // 30.. .... 21..
	WHITE | SLOT0 | 3,   RED   | SLOT0 | OFF, // .1.. 0... 32..
	RED   | SLOT0 | 3,   BLUE  | SLOT0 | OFF, // 02.. 1... .3..
	WHITE | SLOT1 | 10,  BLUE  | SLOT1 | OFF, // 13.. 20.. ....
	BLUE  | SLOT0 | 9,   RED   | SLOT1 | OFF, // 2... 31.. 0...
	WHITE | SLOT0 | 2,   NOP,                 // 3... 02.. 1...
	WHITE | SLOT2 | 0,   RED   | SLOT0 | OFF, // .... 130. 2...
	RED   | SLOT0 | 1,   WHITE | SLOT1 | OFF, // 0... 2.1. 3...
	WHITE | SLOT1 | 1,   BLUE  | SLOT0 | OFF, // 1... 302. ....
	BLUE  | SLOT0 | 10,  WHITE | SLOT0 | OFF, // 2... .13. 0...
	RED   | SLOT1 | 2,   WHITE | SLOT2 | OFF, // 30.. .2.. 1...
	RED   | SLOT0 | OFF, NOP,                 // .1.. .3.. 2...
	WHITE | SLOT1 | OFF, NOP,                 // .2.. .... 3...
	BLUE  | SLOT0 | OFF, NOP,                 // .3.. .... ....
	RED   | SLOT1 | OFF, NOP,                 // .... .... ....
	RED   | SLOT0 | 2,   NOP,                 // 0... .... ....
	BLUE  | SLOT0 | 10,  NOP,                 // 1... .... 0...
	WHITE | SLOT0 | 1,   NOP,                 // 2... 0... 1...
	RED   | SLOT1 | 1,   NOP,                 // 30.. 1... 2...
	WHITE | SLOT1 | 0,   RED   | SLOT0 | OFF, // .1.. 20.. 3...
	WHITE | SLOT2 | 2,   BLUE  | SLOT0 | OFF, // .2.. 310. ....
	BLUE  | SLOT0 | 9,   WHITE | SLOT0 | OFF, // .3.. .21. 0...
	WHITE | SLOT0 | 10,  RED   | SLOT1 | OFF, // .... 032. 1...
	RED   | SLOT0 | 3,   WHITE | SLOT1 | OFF, // 0... 1.3. 2...
	WHITE | SLOT2 | 3,   NOP,                 // 1... 2.0. 3...
	RED   | SLOT1 | 0,   BLUE  | SLOT0 | OFF, // 20.. 3.1. ....
	BLUE  | SLOT0 | 0,   WHITE | SLOT0 | OFF, // 31.. ..2. 0...
	BLUE  | SLOT1 | 8,   RED   | SLOT0 | OFF, // .2.. ..3. 10..
	RED   | SLOT0 | 10,  WHITE | SLOT2 | OFF, // 03.. .... 21..
	WHITE | SLOT0 | 9,   RED   | SLOT1 | OFF, // 1... 0... 32..
	WHITE | SLOT1 | 4,   BLUE  | SLOT0 | OFF, // 2... 10.. .3..
	RED   | SLOT1 | 4,   BLUE  | SLOT1 | OFF, // 30.. 21.. ....
	BLUE  | SLOT0 | 1,   RED   | SLOT0 | OFF, // .1.. 32.. 0...
	BLUE  | SLOT1 | 7,   WHITE | SLOT0 | OFF, // .2.. .3.. 10..
	RED   | SLOT0 | 9,   WHITE | SLOT1 | OFF, // 03.. .... 21..
	WHITE | SLOT0 | 8,   RED   | SLOT1 | OFF, // 1... 0... 32..
	BLUE  | SLOT0 | 2,   NOP,                 // 2... 1... 03..
	RED   | SLOT1 | 8,   BLUE  | SLOT1 | OFF, // 30.. 2... 1...
	BLUE  | SLOT1 | 6,   RED   | SLOT0 | OFF, // .1.. 3... 20..
	RED   | SLOT0 | 5,   WHITE | SLOT0 | OFF, // 02.. .... 31..
	RED   | SLOT2 | 7,   BLUE  | SLOT0 | OFF, // 130. .... .2..
	WHITE | SLOT0 | 5,   RED   | SLOT1 | OFF, // 2.1. 0... .3..
	BLUE  | SLOT1 | 3,   NOP,                 // 3.2. 1... .0..
	BLUE  | SLOT0 | 5,   RED   | SLOT0 | OFF, // ..3. 2... 01..
	WHITE | SLOT1 | 7,   RED   | SLOT2 | OFF, // .... 30.. 12..
	WHITE | SLOT0 | 6,   NOP,                 // .... 01.. 23..
	RED   | SLOT0 | 6,   BLUE  | SLOT1 | OFF, // 0... 12.. 3...
	BLUE  | SLOT0 | 4,   NOP,                 // 1... 23.. 0...
	WHITE | SLOT1 | OFF, NOP,                 // 2... 3... 1...
	WHITE | SLOT0 | OFF, NOP,                 // 3... .... 2...
	RED   | SLOT0 | OFF, NOP,                 // .... .... 3...
	BLUE  | SLOT0 | OFF, NOP,                 // .... .... ....
}

//*****************************************************************************
//
// Runs through the sweep encoded in the given table.
//
//*****************************************************************************
func sweep(idx uint16, states []uint8, count uint16) bool {
	//
	// End the animation after the entire state table has been processed.
	//
	if idx >= count {
		return true
	}

	//
	// No updates are required on the odd steps.
	//
	if (idx & 1) != 0 {
		return false
	}

	//
	// Loop through the two entries of the state table for this step.
	//
	var idx2 uint16
	for idx2 = 0; idx2 < 2; idx2++ {
		//
		// Ignore this entry if it is a NOP.
		//
		if states[idx + idx2] == NOP {
			continue
		}

		//
		// Extract the fields of this entry in the state table.
		//
		color := (states[idx + idx2] & 0xc0) >> 6
		slot := (states[idx + idx2] & 0x30) >> 4
		led := states[idx + idx2] & 0x0f
		if led == OFF {
			led = 0xff
		}

		//
		// Apply this change to the appropriate LED array.
		//
		if color == 0 {
			redLEDs.Set(slot, led)
		} else if color == 1 {
			whiteLEDs.Set(slot, led)
		} else if color == 2 {
			blueLEDs.Set(slot, led)
		}
	}

	//
	// The animation is not done.
	//
	return false
}

//*****************************************************************************
//
// Sweeps the LEDs from left to right, then back right to left.
//
//*****************************************************************************
func sweepH(idx uint16) bool {
	return sweep(idx, sweepStatesH, uint16(len(sweepStatesH)))
}

//*****************************************************************************
//
// Sweeps the LEDs from top to bottom, then back bottom to top.
//
//*****************************************************************************
func sweepV(idx uint16) bool {
	return sweep(idx, sweepStatesV, uint16(len(sweepStatesV)))
}
