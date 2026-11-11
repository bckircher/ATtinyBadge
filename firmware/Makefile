GOFLAGS=-target attiny1616 -opt=s

all:
	@tinygo build ${GOFLAGS} -o ATtinyBadge.elf

clean:
	@tinygo clean

flash:
	@tinygo flash ${GOFLAGS}

size:
	@tinygo build ${GOFLAGS} -size=full -o ATtinyBadge.elf

allocs:
	@tinygo build ${GOFLAGS} -print-allocs=. -o ATtinyBadge.elf

sim:
	@go run .
