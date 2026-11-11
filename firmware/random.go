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
// The current state of the 32-bit random number generator.
//
//*****************************************************************************
var random32State uint32 = 1

//*****************************************************************************
//
// The current state of the 8-bit random number generator.
//
//*****************************************************************************
var random8State uint8 = 1

//*****************************************************************************
//
// Seeds the 32-bit random number generator.
//
//*****************************************************************************
func srandom32(seed uint32) {
	if seed != 0 {
		random32State = seed
	}
}

//*****************************************************************************
//
// Generates a new 32-bit random number.
//
//*****************************************************************************
func random32() uint32 {
	//
	// Get the current random number state.
	//
	temp := random32State

	//
	// Perform a (7, 1, 9) XOR-shift to generate a new random number state.
	//
	temp ^= temp << 7
	temp ^= temp >> 1
	temp ^= temp << 9

	//
	// Update the random number state.
	//
	random32State = temp

	//
	// Return the random number state as the new random number.
	//
	return temp
}

//*****************************************************************************
//
// Seeds the 8-bit random number generator.
//
//*****************************************************************************
func srandom8(seed uint8) {
	if seed != 0 {
		random8State = seed
	}
}

//*****************************************************************************
//
// Generates a new 8-bit random number.
//
//*****************************************************************************
func random8() uint8 {
	//
	// Get the current random number state.
	//
	temp := random8State

	//
	// Perform a (1, 1, 3) XOR-shift to generate a new random number state.
	//
	temp ^= temp << 1
	temp ^= temp >> 1
	temp ^= temp << 3

	//
	// Update the random number state.
	//
	random8State = temp

	//
	// Return the random number state as the new random number.
	//
	return temp
}

//*****************************************************************************
//
// Generates a new 32-bit random number less than a given value.
//
//*****************************************************************************
func randomLessThan32(max int) uint32 {
	//
	// Generate a random number, and discard the lower 16 bits.
	//
	value := random32() >> 16

	//
	// Multiply by the maximum value.  This effectively leaves a 16.16
	// fixed-point value that ranges from 0 to N minus a tiny fraction.
	//
	value *= uint32(max)

	//
	// Discard the lower 16 bits again, leaving only integer portion (a number
	// between 0 and N - 1).
	//
	return value >> 16
}

//*****************************************************************************
//
// Generates a new 8-bit random number less than a given value.
//
//*****************************************************************************
func randomLessThan8(max uint8) uint8 {
	//
	// Generate a new 8-bit random number, multiply it by the maximum value,
	// and divide it by 256 to generate value between 0 and max - 1.
	//
	return uint8((uint16(random8()) * uint16(max)) >> 8)
}
