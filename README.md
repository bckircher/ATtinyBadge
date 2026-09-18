Copyright &copy; 2026 Brian Kircher

Open Source Software: you can modify and/or share it under the terms of the BSD
license file in the root directory of this project.

# ATtinyBadge

This project is a small (75.2 mm by 44.7 mm) circuit board in the shape of the
FIRST&reg; logo.  Each symbol in the logo is painted the appropriate color, and
has a set of 11 LEDs (of matching color).  It is powered by a single CR2032
coin cell battery, which provides a typical run time of N hours.

[animated gif image]

Once turned on, a random, short animation is played one a minute.  Once in a
while, a special, longer animation is played instead.  Pressing the button
causes another animation to be played without further delay (canceling the one
minute delay).  Holding the button for two seconds turns the badge off, placing
it into an extremely low power mode.  This also happens after the badge has
been running for eight hours (so it doesn't run continously until the battery
is fully depleted).  When off, pressing the button turns it back on and repeats
the process.

Full details of the design, development, and construction process can be found
[here](https://bckircher.github.io/ATtinyBadge).


## Hardware

The board is designed using [KiCad](https://www.kicad.org/), with the design
files and resulting gerbers provided in the `hardware` directory.  It uses 6/6
design rules, with a 12 mil minimum drill size, so it should be easily
manufactured by any number of prototype board shops (so long as they are able
to make the very odd, non-rectangular shape of the board).  My personal
preference is [PCBWay](https://www.pcbway.com/).

The components are populated by hand (since the board has components on both
sides, the hot plate technique only works for one side!), and the PCB is then
hand painted to match the FIRST logo and protected with a coating of a spray
shellac.  Safety pins like
[these](https://www.amazon.com/Locking-Safety-Brooch-Silver-Jewelry/dp/B0DNQFC7CN)
are soldered to the exposed pad at the top of the back and used to attach the
badge to a shirt/lapel/etc.


## Firmware

The firmware that runs on the ATtiny1616 is written using
[TinyGo](https://tinygo.org/), and is provided in the `firmware` directory.
One space saving optimization that could be done is a custom target that does
not enable the timing services that are a standard part of the TinyGo targets;
that timer is simply disabled after the application starts, making all the
timer setup code effectively "dead code" (even though it still exists and
runs).

The firmware can also be built and run natively (using normal Go) in order to
develop, debug, and test new animations.  It uses ANSI control sequences with
the command line terminal to generate an ASCII-art simulation of the LED array,
with the LEDs "lighting up" as defined by the animation (with the same timing).


## Trademarks

FIRST&reg; is a trademark of For Inspiration and Recognition of Science and
Technology (FIRST), which does not sponsor, authorize, or endorse this
application.

It is very much hoped that this does not ruffly any feathers at FIRST,
especially since the color palette could not be _exactly_ adhered to with
off-the-shelf acrylic paints.
