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
// Constants for the clock select field of the MCLKCTRLA register; the values
// provided by the TinyGo board support package are not as friendly.
//
//*****************************************************************************
const (
	CLKCTRL_MCLKCTRLA_CLKSEL_OSC20M    = 0x00
	CLKCTRL_MCLKCTRLA_CLKSEL_OSCULP32K = 0x01
	CLKCTRL_MCLKCTRLA_CLKSEL_XOSC32K   = 0x02
	CLKCTRL_MCLKCTRLA_CLKSEL_EXTCLK    = 0x03
)

//*****************************************************************************
//
// Constants for the clock divider field of the MCLKCTRLB register; the values
// provided by the TinyGo board support package are not as friendly.
//
//*****************************************************************************
const (
	CLKCTRL_MCLKCTRLB_PDIV_2  = 0x00
	CLKCTRL_MCLKCTRLB_PDIV_4  = 0x02
	CLKCTRL_MCLKCTRLB_PDIV_8  = 0x04
	CLKCTRL_MCLKCTRLB_PDIV_16 = 0x06
	CLKCTRL_MCLKCTRLB_PDIV_32 = 0x08
	CLKCTRL_MCLKCTRLB_PDIV_64 = 0x0a
	CLKCTRL_MCLKCTRLB_PDIV_6  = 0x10
	CLKCTRL_MCLKCTRLB_PDIV_10 = 0x12
	CLKCTRL_MCLKCTRLB_PDIV_12 = 0x14
	CLKCTRL_MCLKCTRLB_PDIV_24 = 0x16
	CLKCTRL_MCLKCTRLB_PDIV_48 = 0x18
)

//*****************************************************************************
//
// Constants for the input sense configuration field of the PIN2CTRL register;
// the values provided by the TinyGo board support package are not as friendly.
//
//*****************************************************************************
const (
	PORT_PIN2CTRL_ISC_INTDISABLE    = 0x00
	PORT_PIN2CTRL_ISC_BOTHEDGES     = 0x01
	PORT_PIN2CTRL_ISC_RISING        = 0x02
	PORT_PIN2CTRL_ISC_FALLING       = 0x03
	PORT_PIN2CTRL_ISC_INPUT_DISABLE = 0x04
	PORT_PIN2CTRL_ISC_LEVEL         = 0x05
)

//*****************************************************************************
//
// Constants for the sleep mode field of the SLPCTRL register; the values
// provided by the TinyGo board support package are not as friendly.
//
//*****************************************************************************
const (
	SLPCTRL_CTRLA_SMODE_IDLE    = 0x00
	SLPCTRL_CTRLA_SMODE_STANDBY = 0x02
	SLPCTRL_CTRLA_SMODE_PDOWN   = 0x04
)
