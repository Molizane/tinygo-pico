package main

import (
	"machine"
	"time"
)

var period uint64 = 1e9 / 500

func main() {
	pin := machine.LED
	pwm := machine.PWM4

	pwm.Configure(machine.PWMConfig{Period: period})

	ch, err := pwm.Channel(pin)

	if err != nil {
		println(err.Error())
		return
	}

	for i := 1; i < 255; i++ {
		pwm.Set(ch, pwm.Top()/uint32(i))
		time.Sleep(time.Millisecond * 5)
	}
}
