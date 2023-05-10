package main

import (
	"machine"
	"time"
)

/*
  LiquidCrystal Library - Autoscroll

  Demonstrates the use a 16x2 LCD display.  The LiquidCrystal
  library works with all LCD displays that are compatible with the
  Hitachi HD44780 driver. There are many of them out there, and you
  can usually tell them by the 16-pin interface.

  This sketch demonstrates the use of the autoscroll()
  and noAutoscroll() functions to make new text scroll or not.

 The circuit:
 * LCD RS pin to digital pin 12
 * LCD Enable pin to digital pin 11
 * LCD D4 pin to digital pin 5
 * LCD D5 pin to digital pin 4
 * LCD D6 pin to digital pin 3
 * LCD D7 pin to digital pin 2
 * LCD R/W pin to ground
 * 10K resistor:
 * ends to +5V and ground
 * wiper to LCD VO pin (pin 3)

 Library originally added 18 Apr 2008
 by David A. Mellis
 library modified 5 Jul 2009
 by Limor Fried (http://www.ladyada.net)
 example added 9 Jul 2009
 by Tom Igoe
 modified 22 Nov 2010
 by Tom Igoe

 This example code is in the public domain.

 http://www.arduino.cc/en/Tutorial/LiquidCrystal
*/

func main() {
	address := uint16(0x3F)
	columns := byte(16)
	rows := byte(2)
	scl := machine.GP16
	sda := machine.GP17

	// initialize the library with the numbers of the interface pins
	lcd := LCD{
		Address: address,
		Columns: columns,
		Rows:    rows,
		SCL:     scl,
		SDA:     sda,
	}

	// set up the LCD's number of columns and rows:
	lcd.Init()

	for {
		// set the cursor to (0,0):
		lcd.SetCursor(0, 0)

		// print from 0 to 9:
		for thisChar := 0; thisChar < 10; thisChar++ {
			lcd.Print([]byte{byte(thisChar + 0x30)})
			time.Sleep(time.Millisecond * 500)
		}

		time.Sleep(1000)

		// set the cursor to (16,1):
		lcd.SetCursor(16, 1)

		// set the display to automatically scroll:
		lcd.Autoscroll()

		// print from 0 to 9:
		for thisChar := 0; thisChar < 10; thisChar++ {
			lcd.Print([]byte{byte(thisChar + 0x30)})
			time.Sleep(time.Millisecond * 500)
		}

		// turn off automatic scrolling
		lcd.NoAutoscroll()

		time.Sleep(time.Millisecond * 500)

		// clear screen for the next loop:
		lcd.Clear()
	}
}