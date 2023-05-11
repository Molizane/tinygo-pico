package main

import (
	"lcd_i2c/device"
	"log"
	"machine"
	"time"
)

func main() {
	i2c := machine.I2C0

	if err := i2c.Configure(machine.I2CConfig{
		SDA: machine.GP16,
		SCL: machine.GP17,
	}); err != nil {
		log.Fatal(err)
	}

	lcd, err := device.NewLcdI2C(*i2c, 0x3F, device.LCD_16x2)

	if err != nil {
		log.Fatal(err)
	}

	if err = lcd.ShowMessage("Backlight ON", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.SetBacklight(true); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	if err = lcd.ShowMessage("Backlight OFF", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.SetBacklight(false); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	if err = lcd.SetBacklight(true); err != nil {
		log.Fatal(err)
	}

	lcd.Clear()

	if err = lcd.ShowMessage("Cursor ON", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.CursorOn(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	lcd.Clear()

	if err = lcd.ShowMessage("Cursor OFF", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.CursorOff(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	lcd.Clear()

	if err = lcd.ShowMessage("Blink ON", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.BlinkOn(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	lcd.Clear()

	if err = lcd.ShowMessage("Cursor r 2, c 7", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.SetPosition(1, 7); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	lcd.Clear()

	if err = lcd.ShowMessage("Blink OFF", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.BlinkOff(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	bell := []byte{0x04, 0x0e, 0x0e, 0x0e, 0x1f, 0x00, 0x04}
	note := []byte{0x02, 0x03, 0x02, 0x0e, 0x1e, 0x0c, 0x00}
	clock := []byte{0x00, 0x0e, 0x15, 0x17, 0x11, 0x0e, 0x00}
	heart := []byte{0x00, 0x0a, 0x1f, 0x1f, 0x0e, 0x04, 0x00}
	duck := []byte{0x00, 0x0c, 0x1d, 0x0f, 0x0f, 0x06, 0x00}
	check := []byte{0x00, 0x01, 0x03, 0x16, 0x1c, 0x08, 0x00}
	cross := []byte{0x00, 0x1b, 0x0e, 0x04, 0x0e, 0x1b, 0x00}
	retarrow := []byte{0x01, 0x01, 0x05, 0x09, 0x1f, 0x08, 0x04}
	// degree := []byte{0x06, 0x09, 0x09, 0x06}
	// heart := []byte{0x00, 0x0a, 0x1f, 0x1f, 0x0e, 0x04}
	// dot := []byte{0x00, 0x04, 0x0e, 0x1f, 0x1f, 0x0e, 0x04}
	// triang := []byte{0x00, 0x04, 0x0e, 0x1f}
	// ovrscr := []byte{0xff}
	// unchecked := []byte{0x00, 0x0e, 0x11, 0x11, 0x11, 0x0e}
	// checked := []byte{0x00, 0x0e, 0x1f, 0x1f, 0x1f, 0x0e}
	// uparr := []byte{0x04, 0x0e, 0x1f, 0x04, 0x04, 0x04, 0x04, 0x04}
	// alien := []byte{0b11111, 0b10101, 0b11111, 0b11111, 0b01110, 0b01010, 0b11011, 0b00000}
	// speaker := []byte{0b00001, 0b00011, 0b01111, 0b01111, 0b01111, 0b00011, 0b00001, 0b00000}
	// sound := []byte{0b00001, 0b00011, 0b00101, 0b01001, 0b01001, 0b01011, 0b11011, 0b11000}
	// skull := []byte{0b00000, 0b01110, 0b10101, 0b11011, 0b01110, 0b01110, 0b00000, 0b00000}
	// lock := []byte{0b01110, 0b10001, 0b10001, 0b11111, 0b11011, 0b11011, 0b11111, 0b00000}

	lcd.CreateChar(0, bell)
	lcd.CreateChar(1, note)
	lcd.CreateChar(2, clock)
	lcd.CreateChar(3, heart)
	lcd.CreateChar(4, duck)
	lcd.CreateChar(5, check)
	lcd.CreateChar(6, cross)
	lcd.CreateChar(7, retarrow)

	chars := string([]byte{0, 1, 2, 3, 4, 5, 6, 7})

	lcd.Clear()

	if err = lcd.ShowMessage("Custom chars", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}

	if err = lcd.ShowMessage(chars, device.SHOW_LINE_2); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	/*
		if err = lcd.RightToLeft(); err != nil {
			log.Fatal(err)
		}

		lcd.Clear()

		if err = lcd.ShowMessage("Right to left text", device.SHOW_LINE_1); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 5)

		if err = lcd.LeftToRight(); err != nil {
			log.Fatal(err)
		}

		lcd.Clear()

		if err = lcd.ShowMessage("Left to right text", device.SHOW_LINE_1); err != nil {
			log.Fatal(err)
		}
	*/

	lcd.Clear()

	if err = lcd.ShowMessage("Bye", device.SHOW_LINE_1); err != nil {
		log.Fatal(err)
	}
}
