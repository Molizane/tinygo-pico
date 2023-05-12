package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/ds3231"
	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		SDA:       machine.GP16,
		SCL:       machine.GP17,
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	// lcd := hd44780i2c.New(machine.I2C0, 0x27)
	lcd := hd44780i2c.New(machine.I2C0, 0x3F)

	if err := lcd.Configure(hd44780i2c.Config{
		Width:       16,
		Height:      2,
		CursorOn:    false,
		CursorBlink: false,
	}); err != nil {
		log.Fatal(err)
	}

	lcd.ClearDisplay()

	degree := []byte{0b00110, 0b01001, 0b01001, 0b00110, 0b00000, 0b00000, 0b00000, 0b00000}
	lcd.CreateCharacter(0, degree)

	display := ssd1306.NewI2C(machine.I2C0)

	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	display.ClearDisplay()

	rtc := ds3231.New(machine.I2C0)
	rtc.Configure()

	/*
		valid := rtc.IsTimeValid()

		if !valid {
			date := time.Date(2023, 05, 11, 23, 34, 50, 0, time.UTC)
			rtc.SetTime(date)
		}
	*/

	running := rtc.IsRunning()

	if !running {
		if err := rtc.SetRunning(true); err != nil {
			lcd.Print([]byte("Error configuring RTC"))
		}
	}

	var sec int = -1

	x := int16(0)
	y := int16(0)
	deltaX := int16(1)
	deltaY := int16(1)

	for {
		lcd.Home()

		dt, err := rtc.ReadTime()

		if sec != dt.Second() {
			sec = dt.Second()

			var s string

			if err != nil {
				lcd.Print([]byte(fmt.Sprintf("Error reading date: %s", err.Error())))
			} else {
				s = fmt.Sprintf("%02d/%02d/%04d", dt.Day(), dt.Month(), dt.Year())
				lcd.Print([]byte(s))

				s = fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
				lcd.SetCursor(0, 1)
				lcd.Print([]byte(s))
			}

			temp, _ := rtc.ReadTemperature()
			s = fmt.Sprintf("%.2f%cC", float32(temp)/1000, 0)

			lcd.SetCursor(9, 1)
			lcd.Print([]byte(s))

			// time.Sleep(time.Second * 1)
		}

		pixel := display.GetPixel(x, y)
		c := color.RGBA{255, 255, 255, 255}

		if pixel {
			c = color.RGBA{0, 0, 0, 255}
		}

		display.SetPixel(x, y, c)
		display.Display()

		x += deltaX
		y += deltaY

		if x == 0 || x == 127 {
			deltaX = -deltaX
		}

		if y == 0 || y == 63 {
			deltaY = -deltaY
		}

		time.Sleep(1 * time.Millisecond)
	}
}
