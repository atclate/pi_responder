package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const (
	PIN_LOCATION = "%v"
	//PIN_LOCATION = "/sys/class/gpio/gpio%v/value"
)

var (
	gpioMap    = make(map[string]Pin)
	buttonPush = make(chan Pin)
)

type Pin struct {
	Name    string
	Value   string
	Color   string
	Updated time.Time
}

func initGpioPort(num, color string) {
	b, err := ioutil.ReadFile(fmt.Sprintf(PIN_LOCATION, num))
	if err != nil {
		log.Fatal(err)
	}
	value := strings.Trim(string(b), "\n")
	gpioMap[num] = Pin{Name: num, Value: value, Color: color, Updated: time.Now()}
	fmt.Printf("Initializing GPIO %v: %v\n", num, value)
}

func InitGpioPoll() {
	initGpioPort("23", "yellow")
	initGpioPort("24", "blue")
	initGpioPort("25", "red")
	go func() {
		for {
			for k, p := range gpioMap {
				// loop through all gpioMap and watch for differences
				b, err := ioutil.ReadFile(fmt.Sprintf(PIN_LOCATION, k))
				if err != nil {
					log.Fatal(err)
				}
				value := strings.Trim(string(b), "\n")
				if p.Value != value && (value == "0" || value == "1") {
					fmt.Printf("pin %v changed from %v to %v\n", k, p.Value, value)
					p.Value = value
					gpioMap[k] = p
					buttonPush <- p
				}
			}
		}
	}()
}
