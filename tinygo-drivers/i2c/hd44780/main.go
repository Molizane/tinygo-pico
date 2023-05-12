package main

import (
	"fmt"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
)

func main() {
	// Note: most HD44780 LCD modules requires 5V power, however some variations
	// use 3.3V (and may be damaged by 5V).

	if err := machine.I2C0.Configure(machine.I2CConfig{
		SDA:       machine.GP16,
		SCL:       machine.GP17,
		Frequency: machine.TWI_FREQ_400KHZ,
	}); err != nil {
		log.Fatal(err)
	}

	// lcd := hd44780i2c.New(machine.I2C0, 0x27)
	lcd := hd44780i2c.New(machine.I2C0, 0x3F)

	if err := lcd.Configure(hd44780i2c.Config{
		Width:       16, // required
		Height:      2,  // required
		CursorOn:    false,
		CursorBlink: false,
	}); err != nil {
		log.Fatal(err)
	}

	lcd.ClearDisplay()
	lcd.Print([]byte("Backlight ON"))
	lcd.BacklightOn(true)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Backlight OFF"))
	lcd.BacklightOn(false)
	time.Sleep(time.Second * 3)

	lcd.BacklightOn(true)
	lcd.ClearDisplay()
	lcd.Print([]byte("Cursor ON "))
	lcd.CursorOn(true)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Cursor OFF"))
	lcd.CursorOn(false)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Blink ON "))
	lcd.CursorBlink(true)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()

	for i := uint8(0); i <= 15; i++ {
		lcd.SetCursor(0, 0)
		lcd.Print([]byte(fmt.Sprintf("Cursor at 1,%d", i)))
		lcd.SetCursor(i, 1)
		time.Sleep(time.Second * 1)
	}

	lcd.ClearDisplay()
	lcd.Print([]byte("Blink OFF"))
	lcd.CursorBlink(false)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Cursor/Blink ON"))
	lcd.CursorOn(true)
	lcd.CursorBlink(true)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Cursor/Blink OFF"))
	lcd.CursorOn(false)
	lcd.CursorBlink(false)
	time.Sleep(time.Second * 3)

	bell := []byte{0x04, 0x0e, 0x0e, 0x0e, 0x1f, 0x00, 0x04}
	note := []byte{0x02, 0x03, 0x02, 0x0e, 0x1e, 0x0c, 0x00}
	clock := []byte{0x00, 0x0e, 0x15, 0x17, 0x11, 0x0e, 0x00}
	heart := []byte{0x00, 0x0a, 0x1f, 0x1f, 0x0e, 0x04, 0x00}
	duck := []byte{0x00, 0x0c, 0x1d, 0x0f, 0x0f, 0x06, 0x00}
	check := []byte{0x00, 0x01, 0x03, 0x16, 0x1c, 0x08, 0x00}
	cross := []byte{0x00, 0x1b, 0x0e, 0x04, 0x0e, 0x1b, 0x00}
	retarrow := []byte{0x01, 0x01, 0x05, 0x09, 0x1f, 0x08, 0x04}
	degree := []byte{0x06, 0x09, 0x09, 0x06}
	dot := []byte{0x00, 0x04, 0x0e, 0x1f, 0x1f, 0x0e, 0x04}
	triang := []byte{0x00, 0x04, 0x0e, 0x1f}
	ovrscr := []byte{0xff}
	unchecked := []byte{0x00, 0x0e, 0x11, 0x11, 0x11, 0x0e}
	checked := []byte{0x00, 0x0e, 0x1f, 0x1f, 0x1f, 0x0e}
	uparr := []byte{0x04, 0x0e, 0x1f, 0x04, 0x04, 0x04, 0x04, 0x04}
	alien := []byte{0b11111, 0b10101, 0b11111, 0b11111, 0b01110, 0b01010, 0b11011, 0b00000}
	speaker := []byte{0b00001, 0b00011, 0b01111, 0b01111, 0b01111, 0b00011, 0b00001, 0b00000}
	sound := []byte{0b00001, 0b00011, 0b00101, 0b01001, 0b01001, 0b01011, 0b11011, 0b11000}
	skull := []byte{0b00000, 0b01110, 0b10101, 0b11011, 0b01110, 0b01110, 0b00000, 0b00000}
	lock := []byte{0b01110, 0b10001, 0b10001, 0b11111, 0b11011, 0b11011, 0b11111, 0b00000}
	frame := []byte{0b11111, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b11111}
	window := []byte{0b11111, 0b10101, 0b10101, 0b10111, 0b11101, 0b10101, 0b10101, 0b11111}
	crux := []byte{0b00100, 0b00100, 0b11111, 0b00100, 0b00100, 0b00100, 0b00100, 0b01110}
	kylix := []byte{0b11111, 0b01110, 0b01110, 0b00100, 0b01110, 0b00100, 0b01110, 0b11111}

	lcd.CreateCharacter(0, bell)
	lcd.CreateCharacter(1, note)
	lcd.CreateCharacter(2, clock)
	lcd.CreateCharacter(3, heart)
	lcd.CreateCharacter(4, duck)
	lcd.CreateCharacter(5, check)
	lcd.CreateCharacter(6, cross)
	lcd.CreateCharacter(7, retarrow)

	chars := []byte{0, 1, 2, 3, 4, 5, 6, 7}

	lcd.ClearDisplay()
	lcd.Print([]byte("Custom chars 1"))
	lcd.SetCursor(0, 1)
	lcd.Print(chars)
	time.Sleep(time.Second * 3)

	lcd.CreateCharacter(0, degree)
	lcd.CreateCharacter(1, dot)
	lcd.CreateCharacter(2, triang)
	lcd.CreateCharacter(3, ovrscr)
	lcd.CreateCharacter(4, unchecked)
	lcd.CreateCharacter(5, checked)
	lcd.CreateCharacter(6, uparr)
	lcd.CreateCharacter(7, alien)

	lcd.ClearDisplay()
	lcd.Print([]byte("Custom chars 3"))
	lcd.SetCursor(0, 1)
	lcd.Print(chars)
	time.Sleep(time.Second * 3)

	lcd.CreateCharacter(0, speaker)
	lcd.CreateCharacter(1, sound)
	lcd.CreateCharacter(2, skull)
	lcd.CreateCharacter(3, lock)
	lcd.CreateCharacter(4, frame)
	lcd.CreateCharacter(5, window)
	lcd.CreateCharacter(6, crux)
	lcd.CreateCharacter(7, kylix)
	time.Sleep(time.Second * 3)

	lcd.ClearDisplay()
	lcd.Print([]byte("Bye"))
}
