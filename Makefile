build:
	go build

reload: build run

run:
	./pi_responder /sys/class/gpio/gpio23/value
	disable_gpio

enable_gpio: 
	echo "23" > /sys/class/gpio/export
	echo "in" > /sys/class/gpio/gpio23/direction
	echo "24" > /sys/class/gpio/export
	echo "in" > /sys/class/gpio/gpio23/direction
	echo "25" > /sys/class/gpio/export
	echo "in" > /sys/class/gpio/gpio23/direction

disable_gpio:
	echo "23" > /sys/class/gpio/unexport
	echo "24" > /sys/class/gpio/unexport
	echo "25" > /sys/class/gpio/unexport
