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
	"runtime/interrupt"
)

//*****************************************************************************
//
// Constants that control the timing of the animation loop:
//
// * animationDelay - the amount of time between successive animations, in
//                    milliseconds; if this is 1000 milliseconds, there is
//                    instead a 1 second pause between the end of one animation
//                    and the start of the next one (in other words, the delay
//                    adapts to the animation length).
//
// * animationStep - the amount of time between successive steps of an
//                   animation, in milliseconds.
//
// * sleepDelay - the amount of time before the animations halt and the
//                processor is put into standby sleep mode (conserving battery
//                power) until the button is pressed (in milliseconds).
//
//*****************************************************************************
var animationDelay uint32 = 60 * 1000
const animationStep uint32 = 16
const sleepDelay uint32 = 8 * 60 * 60 * 1000

//*****************************************************************************
//
// Set this to true, for testing purposes, to disable putting the processor
// into standby sleep mode (and therefore starting animations as soon as the
// application starts).
//
//*****************************************************************************
const noSleep = false

//*****************************************************************************
//
// Set this to true, for testing purposes, to show the active time of the
// processor on the on-board LED of the ATtiny1616 Breakout with seesaw board
// from Adafruit.  Note that this interferes with the blue LED array on the
// actual board!
//
//*****************************************************************************
const showActive = false

//*****************************************************************************
//
// Set this to true, for testing purposes, to display each animation in order.
// This allows all the animations to be seen quickly and easily, without having
// to wait for each of them to be randomly selected.
//
//*****************************************************************************
const sequential = false

//*****************************************************************************
//
// The array of animation functions.  Adding a new function to this array (or
// removing one!) allows it to automatically selected (at random).
//
// Note that having an odd number of animations makes the random select much
// better!
//
//*****************************************************************************
var animations = []func(idx uint16) bool {
	//
	// Random "sparkles" on each symbol.  This is only here to be uncommented
	// for testing purposes (commenting out all the others); normally, it is
	// occasionally selected as a special, longer animation.
	//
	//sparkle,

	//
	// Outlines the logo with chasing LEDs.
	//
	outline,

	//
	// Simultaneous chasing LEDs around each symbol.
	//
	chase,

	//
	// Sweep from left to right and then back from right to left.
	//
	sweepH,

	//
	// Sweep from top to bottom and then back from bottom to top.
	//
	sweepV,

    //
	// Sweep from the top left to the bottom right and then back from bottom
	// right to top left.
	//
	sweep45,

    //
	// Sweep from the top right to the bottom left and then back from bottom
	// left to top right.
	//
	sweep135,

	//
	// Highlights each symbol in order.
	//
	logo,

	//
	// Test patterns for checking the LED arrays.  Uncomment one of these and
	// comment out all other animations.
	//
	//testRed,
	//testWhite,
	//testBlue,
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
// The pin "address" of the red LEDs on the board (the upper nibble being the
// pin to drive high and the lower nibble being the pin to drive low in order
// to light the given LED of the Charlieplexed array).
//
//*****************************************************************************
var ledsRed = []uint8 {
	0x47, 0x45, 0x54, 0x56, 0x65, 0x67, 0x76, 0x46, 0x64, 0x57, 0x75, 0x74,
}

//*****************************************************************************
//
// The LEDs object for the red LED Charlieplexed array.
//
//*****************************************************************************
var redLEDs = LEDs {
	port:    avr.PORTA,
	mask:    0xf0,
	mapping: ledsRed,
}

//*****************************************************************************
//
// The pin "address" of the white LEDs on the board (the upper nibble being the
// pin to drive high and the lower nibble being the pin to drive low in order
// to light the given LED of the Charlieplexed array).
//
//*****************************************************************************
var ledsWhite = []uint8 {
	0x24, 0x42, 0x35, 0x53, 0x25, 0x23, 0x32, 0x34, 0x43, 0x45, 0x54, 0x52,
}

//*****************************************************************************
//
// The LEDs object for the white LED Charlieplexed array.
//
//*****************************************************************************
var whiteLEDs = LEDs {
	port:    avr.PORTB,
	mask:    0x3c,
	mapping: ledsWhite,
}

//*****************************************************************************
//
// The pin "address" of the blue LEDs on the board (the upper nibble being the
// pin to drive high and the lower nibble being the pin to drive low in order
// to light the given LED of the Charlieplexed array).
//
//*****************************************************************************
var ledsBlue = []uint8 {
	0x01, 0x10, 0x12, 0x21, 0x23, 0x32, 0x02, 0x20, 0x13, 0x31, 0x03, 0x30,
}

//*****************************************************************************
//
// The LEDs object for the blue LED Charlieplexed array.
//
//*****************************************************************************
var blueLEDs = LEDs {
	port:    avr.PORTC,
	mask:    0x0f,
	mapping: ledsBlue,
}

//*****************************************************************************
//
// Tracks the number of milliseconds that have ellapsed while the processor is
// not in standby sleep mode; used to control timing of the animations, the LED
// refresh, and the sampling of the button.
//
//*****************************************************************************
var millis Millis

//*****************************************************************************
//
// Tracks when the animation is active; the LEDs are only updated by the timer
// interrupt handler when the animation is active (to save power).
//
//*****************************************************************************
var active = false

//*****************************************************************************
//
// Called when the button is pressed while in sleep mode (which is when the
// interrupt is enabled).  Only clears the interrupt, since the entire purpose
// is to bring the microcontroller out of sleep mode.
//
//*****************************************************************************
func handlePORTAInterrupt(intr interrupt.Interrupt) {
	//
	// Clear the PORTA interrupt.
	//
	avr.PORTA.INTFLAGS.Set(0x04)
}

//*****************************************************************************
//
// Called when
//
//*****************************************************************************
func handleTCB0Interrupt(intr interrupt.Interrupt) {
	//
	// Increment the millisecond counter.
	//
	elapsed := millis.Increment()

	//
	// Determine the index of the LEDs to display, and refresh the three LED
	// arrays.
	//
	if active {
		idx := int(elapsed & 3)
		redLEDs.Refresh(idx)
		whiteLEDs.Refresh(idx)
		blueLEDs.Refresh(idx)
	}

	//
	// Every 16ms, read the state of the button and perform debouncing.
	//
	if (elapsed & 15) == 0 {
		debounce((avr.PORTA.IN.Get() & 0x04) >> 2)
	}

	//
	// Clear the timer interrupt.
	//
	avr.TCB0.INTFLAGS.Set(avr.TCB_INTFLAGS_CAPT)
}

//*****************************************************************************
//
// Performs on-time setup of the ATtiny.
//
//*****************************************************************************
func setup() {
	//
	// Configure the button pin as an input with the on-chip pull-up enabled.
	//
	avr.PORTA.DIRCLR.Set(0x04)
	avr.PORTA.PIN2CTRL.Set(avr.PORT_PIN2CTRL_PULLUPEN)

	//
	// Disable the RTC that is setup by the base runtime; it is not needed (and
	// gets in the way of staying in sleep mode for extended periods of time).
	//
	avr.RTC.INTCTRL.Set(0)
	avr.RTC.CTRLA.Set(0)

	//
	// If the processor active time should be shown on the on-board LED, enable
	// it now.
	//
	if showActive {
		avr.PORTA.DIRSET.Set(0x04)
		avr.PORTA.OUTCLR.Set(0x04)
	}

    //
	// Seed the random number generators from the unique serial number of this
	// ATtiny1616.
	//
	srandom8(avr.SIGROW.SERNUM0.Get() ^ avr.SIGROW.SERNUM1.Get() ^
		avr.SIGROW.SERNUM2.Get() ^ avr.SIGROW.SERNUM3.Get() ^
		avr.SIGROW.SERNUM4.Get() ^ avr.SIGROW.SERNUM5.Get() ^
		avr.SIGROW.SERNUM6.Get() ^ avr.SIGROW.SERNUM7.Get() ^
		avr.SIGROW.SERNUM8.Get() ^ avr.SIGROW.SERNUM9.Get())
	srandom32((uint32(avr.SIGROW.SERNUM0.Get()) |
		(uint32(avr.SIGROW.SERNUM1.Get()) << 8) |
		(uint32(avr.SIGROW.SERNUM2.Get()) << 16) |
		(uint32(avr.SIGROW.SERNUM3.Get()) << 24)) |
		(uint32(avr.SIGROW.SERNUM4.Get()) |
		(uint32(avr.SIGROW.SERNUM5.Get()) << 8) |
		(uint32(avr.SIGROW.SERNUM6.Get()) << 16) |
		(uint32(avr.SIGROW.SERNUM7.Get()) << 24)) |
		(uint32(avr.SIGROW.SERNUM8.Get()) |
		(uint32(avr.SIGROW.SERNUM9.Get()) << 8)))
}

//*****************************************************************************
//
// Configures the processor to run slowly (32 kHz).
//
//*****************************************************************************
func slow() {
	//
	// Disable the timer.
	//
	avr.TCB0.INTCTRL.Set(0)
	avr.TCB0.CTRLA.Set(0)

	//
	// Set the CPU speed to 32 kHz, being careful to not set it to run at
	// 20 MHz at an intermediate step (which is not supported since VDD is
	// 3 V).
	//
	avr.CPU.CCP.Set(0xd8)
	avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_64 |
		avr.CLKCTRL_MCLKCTRLB_PEN)
	avr.CPU.CCP.Set(0xd8)
	avr.CLKCTRL.MCLKCTRLA.Set(CLKCTRL_MCLKCTRLA_CLKSEL_OSCULP32K)
	avr.CPU.CCP.Set(0xd8)
	avr.CLKCTRL.MCLKCTRLB.Set(0)
}

//*****************************************************************************
//
// Configures the processor to run quickly (312.5 kHz).
//
//*****************************************************************************
func fast() {
	//
	// Set the CPU speed to 312.5 kHz, being careful to not set it to run at
	// 20 MHz at an intermediate step (which is not supported since VDD is
	// 3 V).
	//
	avr.CPU.CCP.Set(0xd8)
	//avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_16 |
	//	avr.CLKCTRL_MCLKCTRLB_PEN)
	//avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_24 |
	//	avr.CLKCTRL_MCLKCTRLB_PEN)
	//avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_32 |
	//	avr.CLKCTRL_MCLKCTRLB_PEN)
	//avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_48 |
	//	avr.CLKCTRL_MCLKCTRLB_PEN)
	avr.CLKCTRL.MCLKCTRLB.Set(CLKCTRL_MCLKCTRLB_PDIV_64 |
		avr.CLKCTRL_MCLKCTRLB_PEN)
	avr.CPU.CCP.Set(0xd8)
	avr.CLKCTRL.MCLKCTRLA.Set(CLKCTRL_MCLKCTRLA_CLKSEL_OSC20M)

	//
	// Set the timer to generate an interrupt at 1 kHz.
	//
	//avr.TCB0.CCMPL.Set(1249 & 0xff)
	//avr.TCB0.CCMPH.Set(1249 >> 8)
	//avr.TCB0.CCMPL.Set(833 & 0xff)
	//avr.TCB0.CCMPH.Set(833 >> 8)
	//avr.TCB0.CCMPL.Set(624 & 0xff)
	//avr.TCB0.CCMPH.Set(624 >> 8)
	//avr.TCB0.CCMPL.Set(416 & 0xff)
	//avr.TCB0.CCMPH.Set(416 >> 8)
	avr.TCB0.CCMPL.Set(312  & 0xff)
	avr.TCB0.CCMPH.Set(312 >> 8)
	avr.TCB0.CTRLA.Set(avr.TCB_CTRLA_ENABLE)
	avr.TCB0.INTFLAGS.Set(avr.TCB_INTFLAGS_CAPT)
	avr.TCB0.INTCTRL.Set(avr.TCB_INTCTRL_CAPT)
}

//*****************************************************************************
//
// Puts the processor into idle sleep mode (leaving everything else running,
// until the next interrupt).
//
//*****************************************************************************
func idle() {
	//
	// If the processor active time should be shown on the on-board LED, turn
	// it off now.
	//
	if showActive {
		avr.PORTA.OUTCLR.Set(0x04)
	}

	//
	// Select the idle sleep mode.
	//
	avr.SLPCTRL.CTRLA.Set(SLPCTRL_CTRLA_SMODE_IDLE | avr.SLPCTRL_CTRLA_SEN)

	//
	// Put the processor to sleep.
	//
	avr.Asm("sleep")

	//
	// If the processor active time should be shown on the on-board LED, turn
	// it on now.
	//
	if showActive {
		avr.PORTA.OUTSET.Set(0x04)
	}
}

//*****************************************************************************
//
// Puts the processor into standby sleep mode (only to be awoken by the button
// pin interrupt).
//
//*****************************************************************************
func sleep() {
	//
	// If the processor active time should be shown on the on-board LED, turn
	// it off now.
	//
	if showActive {
		avr.PORTA.OUTCLR.Set(0x04)
	}

	//
	// Select the standby sleep mode.
	//
	avr.SLPCTRL.CTRLA.Set(SLPCTRL_CTRLA_SMODE_STANDBY | avr.SLPCTRL_CTRLA_SEN)

	//
	// Enable the button pin interrupt, and clear any previously detected pin
	// interrupt.
	//
	avr.PORTA.PIN2CTRL.Set(avr.PORT_PIN2CTRL_PULLUPEN |
		PORT_PIN2CTRL_ISC_FALLING)
	avr.PORTA.INTFLAGS.Set(0x04)

	//
	// Put the processor to sleep.
	//
	avr.Asm("sleep")

	//
	// Disable the button pin interrupt.
	//
	avr.PORTA.PIN2CTRL.Set(avr.PORT_PIN2CTRL_PULLUPEN)

	//
	// If the processor active time should be shown on the on-board LED, turn
	// it on now.
	//
	if showActive {
		avr.PORTA.OUTSET.Set(0x04)
	}
}

//*****************************************************************************
//
// Turns off all the LEDs.
//
//*****************************************************************************
func allOff() {
	//
	// Turn off all the LEDs.
	//
	redLEDs.ForceOff()
	whiteLEDs.ForceOff()
	blueLEDs.ForceOff()
}

//*****************************************************************************
//
// The main program.  Yes, there are labels and gotos; they greatly simplify
// the logic and are for exceptional circumstances (which it could be argued is
// an acceptable use for a goto).
//
//*****************************************************************************
func main() {
	//
	// Set the interrupt handlers.  In TinyGo, this is actually handled by the
	// compiler, not at runtime (as one would expect).
	//
	interrupt.New(avr.IRQ_PORTA_PORT, handlePORTAInterrupt)
	interrupt.New(avr.IRQ_TCB0_INT, handleTCB0Interrupt)

	//
	// Perform one-time setup.
	//
	setup()

	//
	// Check the button.  If it is being pressed when the battery is installed,
	// change the animation delay to 1 s (as a demo mode).
	//
	if (avr.PORTA.IN.Get() & 0x04) == 0 {
		animationDelay = 1000
	}

	//
	// Loop forever.
	//
	for {
		//
		// Turn off off the LEDs.
		//
		allOff();

		//
		// Wait for the button to be released.
		//
		for (avr.PORTA.IN.Get() & 0x04) == 0 {
		}

		//
		// Switch to the slow clock.
		//
		if !noSleep {
			slow()
		}

		//
		// Wait for the button to be pressed.
		//
		if !noSleep {
			sleep()
		}

		//
		// Switch to the fast clock.
		//
		fast()

		//
		// Get the time to start the next animation.
		//
		next := millis.Get() + animationDelay

		//
		// Get the time to go back to sleep.
		//
		end := millis.Get() + sleepDelay

		//
		// Clear any button press/hold events that resulted from the button
		// press to wake up.
		//
		wasPressed()
		wasHeld()

		//
		// The previous animation, to avoid repeating them.
		//
		var prevAnim uint8 = 0xff

		//
		//
		// Loop until the sleep time arrives.
		//
		for int(end - millis.Get()) > 0 {
newAnimation:
			//
			// Select an animation to play.  Every so often, play the special,
			// longer, sparkle animation.
			//
			var anim uint8 = 0
			var animFn func(idx uint16) bool = nil
			if sequential {
				if prevAnim == 0xff {
					animFn = sparkle
				} else {
					animFn = animations[prevAnim]
				}
				prevAnim++
				if prevAnim == uint8(len(animations)) {
					prevAnim = 0xff
				}
			} else {
				if random8() > 0xf8 {
					animFn = sparkle
					prevAnim = 0xff
				} else {
					for {
						anim = randomLessThan8(uint8(len(animations)))
						if (len(animations) == 1) || (anim != prevAnim) {
							break
						}
					}
					animFn = animations[anim]
					prevAnim = anim
				}
			}

			//
			// Get the next time to step the animation.
			//
			start := (millis.Get() & 0xfffffff0) + animationStep

			//
			// Set the initial animation step.
			//
			var idx uint16 = 0

			//
			// Indicate that the animation is active, so that the LEDs get
			// updated by the timer interrupt handler.
			//
			active = true

			//
			// Loop while the animation is not done.
			//
			for done := false; !done; {
				//
				// Wait until it is time for the next animation step.
				//
				for int(start - millis.Get()) > 0 {
					//
					// Save power by halting the processor.
					//
					idle()
				}

				//
				// Immediately go to sleep if the button was held.
				//
				if !noSleep && wasHeld() {
					goto goToSleep
				}

				//
				// Call the selected animation.
				//
				done = animFn(idx)

				//
				// Advance to the next animation step and time.
				//
				start += animationStep
				idx++
			}

			//
			// Indicate that the animation is no longer active, so that the
			// LEDs do not get updated by the timer interrupt handler (saving
			// power).
			//
			active = false

			//
			// Turn off all the LEDs (so the individual animations don't have
			// to each handle this at the end).
			//
			allOff()

			//
			// Discard any button presses that may have occurred during the
			// animation.
			//
			wasPressed()

			//
			// If the animation delay is one second, update the next animation
			// time to be one second after the previous animation completed.
			//
			if animationDelay == 1000 {
				next = millis.Get() + 1000
			}

			//
			// Wait until it is time for the next animiation.
			//
			for int(next - millis.Get()) > 0 {
				//
				// Save power by halting the processor.
				//
				idle()

				//
				// Advance to the next animation if the button was pressed.
				//
				if wasPressed() {
					goto newAnimation
				}

				//
				// Immediately go to sleep if the button was held.
				//
				if !noSleep && wasHeld() {
					goto goToSleep
				}
			}

			//
			// Set the time for the next animation to start.
			//
			next += animationDelay
		}

		//
		// Get the next time to step the animation.
		//
goToSleep:
		start := (millis.Get() & 0xfffffff0) + animationStep

		//
		// Set the initial animation step.
		//
		var idx uint16 = 0

		//
		// Indicate that the animation is active, so that the LEDs get updated
		// by the timer interrupt handler.
		//
		active = true

		//
		// Loop while the animation is not done.
		//
		for done := false; !done; {
			//
			// Wait until it is time for the next animation step.
			//
			for int(start - millis.Get()) > 0 {
				//
				// Save power by halting the processor.
				//
				idle()
			}

			//
			// Call the selected animation.
			//
			done = off(idx)

			//
			// Advance to the next animation step and time.
			//
			start += animationStep
			idx++
		}

		//
		// Indicate that the animation is no longer active, so that the LEDs do
		// not get updated by the timer interrupt handler (saving power).
		//
		active = false

		//
		// Turn off all the LEDs (so the individual animations don't have to
		// each handle this at the end).
		//
		allOff()
	}
}
