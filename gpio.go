package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

const (
	PIN_LOCATION = "/sys/class/gpio/gpio%v/value"
)

func initGpioPort(num, color string) {
	b, err := ioutil.ReadFile(fmt.Sprintf(PIN_LOCATION, num))
	if err != nil {
		log.Fatal(err)
	}
	gpioMap[num] = Pin{name: num, value: string(b), color: color}
}

func InitGpioPoll() {
	initGpioPort("23", "yellow")
	initGpioPort("24", "blue")
	initGpioPort("25", "red")
	change := make(chan Pin)
	go func() {
		for {
			for k, p := range gpioMap {
				// loop through all gpioMap and watch for differences
				//b, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%v/value", k))
				b, err := ioutil.ReadFile(fmt.Sprintf(PIN_LOCATION, k))
				if err != nil {
					log.Fatal(err)
				}
				if p.value != string(b) {
					fmt.Printf("pin %v changed from %v to %v", k, p.value, string(b))
					p.value = string(b)
					gpioMap[k] = p
					change <- p
				}
			}
		}
	}()
	for {
		pin := <-change
		fmt.Printf("handle pin change here: Pin = %V\n", pin)
	}
}

type Pin struct {
	name  string
	value string
	color string
}
