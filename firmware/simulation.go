//*****************************************************************************
//
// Copyright (c) 2026 Brian Kircher
//
// Open Source Software: you can modify and/or share it under the terms of the
// BSD license file in the root directory of this project.
//
//*****************************************************************************

//go:build !baremetal

package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//*****************************************************************************
//
// The ANSI control sequence for setting the simulated LED into the various
// states/colors.  Note that when LEDs are off, they are displayed as a dot so
// that the position of all the LEDs can be seen at all times.
//
//*****************************************************************************
const colorRedOn = "\033[41m \033[40m"
const colorRedOff = "\033[31mx\033[37m"
const colorWhiteOn = "\033[47m \033[40m"
const colorWhiteOff = "x"
const colorBlueOn = "\033[44m \033[40m"
const colorBlueOff = "\033[34mx\033[37m"

//*****************************************************************************
//
// The current location of the cursor on the screen.
//
//*****************************************************************************
var cursorX = 0
var cursorY = 18

//*****************************************************************************
//
// Sets an LED at the given location on the screen with the given ANSI control
// sequence.
//
//*****************************************************************************
func SetLED(x int, y int, c string) {
	//
	// See if the X column needs to be adjusted.
	//
	if x != cursorX {
		//
		// Move the cursor left or right by the required number of columns.
		//
		if x < cursorX {
			fmt.Printf("\033[%dD", cursorX - x)
		} else {
			fmt.Printf("\033[%dC", x - cursorX)
		}

		//
		// Update the X position of the cursor.
		//
		cursorX = x
	}

	//
	// See if the Y row needs to be adjusted.
	//
	if y != cursorY {
		//
		// Move the cursor up or down by the required number of rows.
		//
		if y < cursorY {
			fmt.Printf("\033[%dA", cursorY - y)
		} else {
			fmt.Printf("\033[%dB", y - cursorY)
		}

		//
		// Update the Y position of the cursor.
		//
		cursorY = y
	}

	//
	// Display the LED in its given state.
	//
	fmt.Printf("%s", c)

	//
	// The cursor has moved one column to the right.
	//
	cursorX++
}

//*****************************************************************************
//
// A structure to provide the position of an LED on the screen.
//
//*****************************************************************************
type Point struct {
	//
	// The X column of this LED.
	//
	x int

	//
	// The Y row of this LED.
	//
	y int
}

//*****************************************************************************
//
// A substitute LEDs structure for use by the simulator.
//
//*****************************************************************************
type LEDs struct {
	//
	// Tne ANSI control string to turn on LEDs in this array.
	//
	colorOn  string

	//
	// The ANSI control string to turn off LEDs in this array.
	//
	colorOff string

	//
	// The screen positions of this set of LEDs.
	//
	position [11]Point

	//
	// The four slots for this LED array, indicating the LEDs that are turned
	// on at any given time.
	//
	slots    [4]uint8
}

//*****************************************************************************
//
// Sets the state of an LED in the given slot.
//
//*****************************************************************************
func (l *LEDs) Set(idx uint8, led uint8) {
	//
	// See if the LED in this slot is currently turned on.
	//
	if l.slots[idx] != 0xff {
		//
		// Turn off the LED in this slot.
		//
		SetLED(l.position[l.slots[idx]].x, l.position[l.slots[idx]].y,
			l.colorOff)

		//
		// Indicate that th LED in this slot is not on.
		//
		l.slots[idx] = 0xff
	}

	//
	// Do nothing if this is either an LED that doesn't exist or is the special
	// indicator for turning off an LED.
	//
	if (led == 11) || (led == 0xff) {
		return
	}

	//
	// Turn on the LED at this position.
	//
	SetLED(l.position[led].x, l.position[led].y, l.colorOn)

	//
	// Save the state of the LED in this slot.
	//
	l.slots[idx] = led
}

//*****************************************************************************
//
// Turn off all the LEDs in the array.
//
//*****************************************************************************
func (l *LEDs) ForceOff() {
	//
	// Loop through all the LEDs in the array.
	//
	for idx := 0; idx < 11; idx++ {
		//
		// Turn off the LED at this position.
		//
		SetLED(l.position[idx].x, l.position[idx].y, l.colorOff)
	}

	//
	// Loop through all the slots in this array.
	//
	for idx := 0; idx < 4; idx++ {
		//
		// Indicate that this slot is not in use.
		//
		l.slots[idx] = 0xff
	}
}

//*****************************************************************************
//
// The substitute LEDs definition for the red LED array.
//
//*****************************************************************************
var redLEDs = LEDs {
	colorOn:  colorRedOn,
	colorOff: colorRedOff,
	position: [11]Point {
		{ 24, 12 },
		{ 11, 15 },
		{ 4, 17 },
		{ 5, 13 },
		{ 7, 9 },
		{ 8, 5 },
		{ 10, 1 },
		{ 15, 3 },
		{ 20, 6 },
		{ 25, 8 },
		{ 31, 10 },
	},
}

//*****************************************************************************
//
// The substitute LEDs definition for the white LED array.
//
//*****************************************************************************
var whiteLEDs = LEDs {
	colorOn:  colorWhiteOn,
	colorOff: colorWhiteOff,
	position: [11]Point {
		{ 39, 15 },
		{ 32, 15 },
		{ 24, 15 },
		{ 18, 12 },
		{ 14, 9 },
		{ 22, 3 },
		{ 31, 2 },
		{ 39, 2 },
		{ 50, 6 },
		{ 51, 10 },
		{ 46, 13 },
	},
}

//*****************************************************************************
//
// The substitute LEDs definition for the blue LED array.
//
//*****************************************************************************
var blueLEDs = LEDs {
	colorOn:  colorBlueOn,
	colorOff: colorBlueOff,
	position: [11]Point {
		{ 40, 11 },
		{ 34, 8 },
		{ 40, 6 },
		{ 45, 3 },
		{ 51, 0 },
		{ 56, 3 },
		{ 62, 5 },
		{ 68, 8 },
		{ 62, 11 },
		{ 57, 14 },
		{ 51, 17 },
	},
}

//*****************************************************************************
//
// The mapping of the current index to the next index, when stepping through
// the LEDs by 3 with a modulo of 12.  Using this avoids using the modulo
// operator, which is very expensive.
//
//*****************************************************************************
var step3Mod12 [12]uint8 = [12]uint8{ 3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2, 0 }

//*****************************************************************************
//
// The minimum duration of the run of the animation.
//
//*****************************************************************************
var min time.Duration = 1000000000

//*****************************************************************************
//
// The accumulation of the run durations of the animation.
//
//*****************************************************************************
var sum time.Duration = 0

//*****************************************************************************
//
// The maximum duration of the run of the animation.
//
//*****************************************************************************
var max time.Duration = 0

//*****************************************************************************
//
// The number of runs of the animation.
//
//*****************************************************************************
var cnt int = 0

//*****************************************************************************
//
// Runs a simulation of an LED animation.
//
//*****************************************************************************
func main() {
	//
	// Parse the command line flags.
	//
	animPtr := flag.String("anim", "sparkle", "select the animation to run")
	ratePtr := flag.Uint("rate", 1, "the rate to run the animation")
	timePtr := flag.Bool("time", false, "display the run time statistics")
	flag.Parse()

	//
	// Select the animation based on the command line flag selection.
	//
	anim := sparkle
	switch *animPtr {
		case "chase":
			anim = chase
		case "logo":
			anim = logo
		case "off":
			anim = off
		case "outline":
			anim = outline
		case "sweep135":
			anim = sweep135
		case "sweep45":
			anim = sweep45
		case "sweepH":
			anim = sweepH
		case "sweepV":
			anim = sweepV
		case "testRed":
			anim = testRed
		case "testWhite":
			anim = testWhite
		case "testBlue":
			anim = testBlue
		default:
	}

	//
	// Confirm that the selected rate is valid.
	//
	if *ratePtr == 0 {
		*ratePtr = 1
	}

	//
	// Hide the cursor and clear out enough space for the display.
	//
	fmt.Printf("\033[?25l%s", strings.Repeat("\n", 18))

	//
	// Turn off all the LEDs in all three arrays (resulting in the placeholder
	// LEDs being drawn on the screen).
	//
	redLEDs.ForceOff()
	whiteLEDs.ForceOff()
	blueLEDs.ForceOff()

	//
	// Create channels for communication with the animation goroutine that it
	// needs to halt.
	//
	halt := make(chan int)
	halted := make(chan int)

	//
	// A goroutine for updating the display with the selected animation.
	//
	go func() {
		//
		// Create a ticker that fires every 16 ms (to match the update rate
		// on the ATtinyBadge).
		//
		ticker := time.NewTicker(16 * time.Millisecond)

		//
		// Loop forever, running through the loop every time the ticker fires.
		//
		var idx uint16 = 0
		done := 0
		for range ticker.C {
			//
			// See if a halt request has been received.
			//
			select {
				//
				// A halt request was received.
				//
				case <- halt:
					//
					// Indicate that the goroutine has halted.
					//
					halted <- 1

					//
					// Exit the goroutine.
					//
					return

				//
				// A halt request was not received, so fall out of the select
				// statement.
				//
				default:
			}

			//
			// See if the animation is done.
			//
			if done == 0 {
				//
				// Run the next step of the animation.  Also capture and record
				// the run time.
				//
				result := false
				if (idx % uint16(*ratePtr)) == 0 {
					start := time.Now()
					result = anim(idx / uint16(*ratePtr))
					end := time.Now()
					delta := end.Sub(start)
					if delta < min {
						min = delta
					}
					sum += delta
					if delta > max {
						max = delta
					}
					cnt++
				}

				//
				// See if the animation is done.
				//
				if result {
					//
					// The animation is done, so set the done counter to keep
					// the display quiescent for a second.
					//
					done = 63

					//
					// Reset the animation counter.
					//
					idx = 0

					//
					// Turn off all the LEDs.
					//
					redLEDs.ForceOff()
					whiteLEDs.ForceOff()
					blueLEDs.ForceOff()
				} else {
					//
					// The animation isn't done yet, so increment its counter.
					//
					idx++
				}
			} else {
				//
				// Since the animation is done, decrement the down counter by
				// one to reach the done timeout.
				//
				done--
			}
		}
	}()

	//
	// Create a signal catcher for SIGINT and SIGTERM.
	//
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT,
		syscall.SIGTERM)
	defer stop()

	//
	// Wait until a signal is caught.
	//
	<-ctx.Done()

	//
	// Request that the animation goroutine halt, and wait for it to confirm
	// it has halted.
	//
	halt <- 1
	<- halted

	//
	// Make the cursor visible again.  An X of -2 is used to account for the
	// Ctrl-C that gets printed to the display (moving the cursor two
	// characters further to the right than expected); requesting an X of -2
	// results in the cursur being on the far left column.
	//
	SetLED(-2, 18, "\033[?25h\033[0m")

	//
	// Display the run time statistics if requested.
	//
	if *timePtr {
		fmt.Printf("\nMin: %d\n", int(min))
		fmt.Printf("Avg: %d\n", int(sum) / cnt)
		fmt.Printf("Max: %d\n", int(max))
	}
}
