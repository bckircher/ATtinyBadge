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
	"runtime/volatile"
)

//*****************************************************************************
//
// The definition of the structure used to track the number of milliseconds.
//
//*****************************************************************************
type Millis struct {
	millis uint32
}

//*****************************************************************************
//
// Increments the count of milliseconds and returns the new value.  Must be
// called every millisecond (or, twice every two milliseconds, and so on).
//
//*****************************************************************************
func (m *Millis) Increment() uint32 {
	//
	// Increment the counter.
	//
	m.millis++

	//
	// Return the new value.
	//
	return m.millis
}

//*****************************************************************************
//
// Gets the current value of the millisecond counter.
//
//*****************************************************************************
func (m *Millis) Get() uint32 {
	//
	// Perform a volatile read of the millisecond counter and return the value.
	// The volatile read is needed to make sure the value is re-read every time
	// when this is called in a loop waiting for time to pass.
	//
	return volatile.LoadUint32(&(m.millis))
}
