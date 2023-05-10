package main

import (
	"lcd_i2c/device"
	"log"
	"machine"
	"time"
)

func main() {
	i2c := machine.I2C0

	err := i2c.Configure(machine.I2CConfig{
		SDA: machine.GP16,
		SCL: machine.GP17,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Construct lcd-device connected via I2C connection
	lcd, err := device.NewLcdI2C(*i2c, 0x3F, device.LCD_16x2)

	if err != nil {
		log.Fatal(err)
	}

	// Turn on the backlight
	err = lcd.BacklightOn()

	if err != nil {
		log.Fatal(err)
	}

	// Put text on 1 line of lcd-display
	err = lcd.ShowMessage(" ! Let's rock ! ", device.SHOW_LINE_1)

	if err != nil {
		log.Fatal(err)
	}

	// Wait 5 secs
	time.Sleep(5 * time.Second)

	// Output text to 2 line of lcd-screen
	err = lcd.ShowMessage("Welcome to RPi dude!", device.SHOW_LINE_2)

	if err != nil {
		log.Fatal(err)
	}

	// Wait 5 secs
	time.Sleep(5 * time.Second)

	// Turn off the backlight and exit
	err = lcd.BacklightOff()

	if err != nil {
		log.Fatal(err)
	}

	// Wait 5 secs
	time.Sleep(5 * time.Second)

	// Turn off the backlight and exit
	err = lcd.BacklightOn()

	if err != nil {
		log.Fatal(err)
	}
}
