package main

import (
	"machine"
	"time"
)

const (
	// commands
	LCD_CLEARDISPLAY   = 0x01
	LCD_RETURNHOME     = 0x02
	LCD_ENTRYMODESET   = 0x04
	LCD_DISPLAYCONTROL = 0x08
	LCD_CURSORSHIFT    = 0x10
	LCD_FUNCTIONSET    = 0x20
	LCD_SETCGRAMADDR   = 0x40
	LCD_SETDDRAMADDR   = 0x80

	// flags for display entry mode
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00

	// flags for display on/off control
	LCD_DISPLAYON  = 0x04
	LCD_DISPLAYOFF = 0x00
	LCD_CURSORON   = 0x02
	LCD_CURSOROFF  = 0x00
	LCD_BLINKON    = 0x01
	LCD_BLINKOFF   = 0x00

	// flags for display/cursor shift
	LCD_DISPLAYMOVE = 0x08
	LCD_CURSORMOVE  = 0x00
	LCD_MOVERIGHT   = 0x04
	LCD_MOVELEFT    = 0x00

	// flags for function set
	LCD_8BITMODE = 0x10
	LCD_4BITMODE = 0x00
	LCD_2LINE    = 0x08
	LCD_1LINE    = 0x00
	LCD_5x10DOTS = 0x04
	LCD_5x8DOTS  = 0x00

	// flags for backlight control
	LCD_BACKLIGHT   = 0x08
	LCD_NOBACKLIGHT = 0x00

	En = 0b00000100 // Enable bit
	Rw = 0b00000010 // Read/Write bit
	Rs = 0b00000001 // Register select bit
)

type LCD struct {
	Frequency       uint32
	Address         uint16
	Columns         byte
	Rows            byte
	SCL             machine.Pin
	SDA             machine.Pin
	dotSize         byte
	i2c             *machine.I2C
	err             error
	ready           bool
	displayfunction byte
	displaymode     byte
	displaycontrol  byte
	numlines        byte
	backlightval    byte
	col             byte
}

func (lcd *LCD) Init() error {
	lcd.displayfunction = LCD_4BITMODE | LCD_1LINE | LCD_5x8DOTS
	return lcd.begin()
}

func (lcd *LCD) begin() error {
	lcd.i2c = machine.I2C0

	lcd.err = lcd.i2c.Configure(machine.I2CConfig{
		Frequency: lcd.Frequency,
		SCL:       lcd.SCL,
		SDA:       lcd.SDA,
	})

	lcd.ready = lcd.err != nil

	if !lcd.ready {
		return lcd.err
	}

	if lcd.Rows > 1 {
		lcd.displayfunction |= LCD_2LINE
	}

	// for some 1 line displays you can select a 10 pixel high font
	if lcd.dotSize != 0 && lcd.Rows == 1 {
		lcd.displayfunction |= LCD_5x10DOTS
	}

	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way before 4.5V so we'll wait 50
	time.Sleep(time.Millisecond * 50)

	// Now we pull both RS and R/W low to begin commands
	lcd.expanderWrite(lcd.backlightval) // reset expander and turn backlight off (Bit 8 =1)
	time.Sleep(time.Millisecond * 1000)

	//put the LCD into 4 bit mode
	// this is according to the hitachi HD44780 datasheet
	// figure 24, pg 46

	// we start in 8bit mode, try to set 4 bit mode
	lcd.write4bits(0x03 << 4)
	time.Sleep(time.Microsecond * 4500) // wait min 4.1ms

	// second try
	lcd.write4bits(0x03 << 4)
	time.Sleep(time.Microsecond * 4500) // wait min 4.1ms

	// third go!
	lcd.write4bits(0x03 << 4)
	time.Sleep(time.Microsecond * 150)

	// finally, set to 4-bit interface
	lcd.write4bits(0x02 << 4)

	// set # lcd.Rows, font size, etc.
	lcd.command(LCD_FUNCTIONSET | lcd.displayfunction)

	// turn the display on with no cursor or blinking default
	lcd.displaycontrol = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	lcd.Display()

	// clear it off
	lcd.Clear()

	// Initialize to default text direction (for roman languages)
	lcd.displaymode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT

	// set the entry mode
	lcd.command(LCD_ENTRYMODESET | lcd.displaymode)

	lcd.Home()

	return nil
}

func (lcd *LCD) Clear() {
	lcd.command(LCD_CLEARDISPLAY)       // clear display, set cursor position to zero
	time.Sleep(time.Microsecond * 2000) // this command takes a long time!
}

func (lcd *LCD) Home() {
	lcd.command(LCD_RETURNHOME)         // set cursor position to zero
	time.Sleep(time.Microsecond * 2000) // this command takes a long time!
}

func (lcd *LCD) SetCursor(column byte, row byte) {
	row_offsets := [4]byte{0x00, 0x40, 0x14, 0x54}

	if row > lcd.numlines {
		row = lcd.numlines - 1 // we count rows starting w/0
	}

	lcd.command(LCD_SETDDRAMADDR | (lcd.col + row_offsets[row]))
}

// Turn the display on/off (quickly)
func (lcd *LCD) NoDisplay() {
	lcd.displaycontrol &= LCD_DISPLAYON ^ 255
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

func (lcd *LCD) Display() {
	lcd.displaycontrol |= LCD_DISPLAYON
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

// Turns the underline cursor on/off
func (lcd *LCD) NoCursor() {
	lcd.displaycontrol &= LCD_CURSORON ^ 255
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

func (lcd *LCD) Cursor() {
	lcd.displaycontrol |= LCD_CURSORON
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

func (lcd *LCD) Print(v []byte) {
	for _, b := range v {
		lcd.write(b)
	}
}

// Turn on and off the blinking cursor
func (lcd *LCD) NoBlink() {
	lcd.displaycontrol &= LCD_BLINKON ^ 255
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

func (lcd *LCD) Blink() {
	lcd.displaycontrol |= LCD_BLINKON
	lcd.command(LCD_DISPLAYCONTROL | lcd.displaycontrol)
}

// These commands scroll the display without changing the RAM
func (lcd *LCD) ScrollDisplayLeft() {
	lcd.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT)
}

func (lcd *LCD) ScrollDisplayRight() {
	lcd.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT)
}

// This is for text that flows Left to Right
func (lcd *LCD) LeftToRight() {
	lcd.displaymode |= LCD_ENTRYLEFT
	lcd.command(LCD_ENTRYMODESET | lcd.displaymode)
}

// This is for text that flows Right to Left
func (lcd *LCD) RightToLeft() {
	lcd.displaymode &= LCD_ENTRYLEFT ^ 255
	lcd.command(LCD_ENTRYMODESET | lcd.displaymode)
}

// This will 'right justify' text from the cursor
func (lcd *LCD) Autoscroll() {
	lcd.displaymode |= LCD_ENTRYSHIFTINCREMENT
	lcd.command(LCD_ENTRYMODESET | lcd.displaymode)
}

// This will 'left justify' text from the cursor
func (lcd *LCD) NoAutoscroll() {
	lcd.displaymode &= LCD_ENTRYSHIFTINCREMENT ^ 255
	lcd.command(LCD_ENTRYMODESET | lcd.displaymode)
}

// Allows us to fill the first 8 CGRAM locations
// with custom characters
func (lcd *LCD) CreateChar(location byte, charmap [8]byte) {
	location &= 0x7 // we only have 8 locations 0-7
	lcd.command(LCD_SETCGRAMADDR | (location << 3))

	for i := 0; i < 8; i++ {
		lcd.write(charmap[i])
	}
}

func (lcd *LCD) write(w byte) error {
	if !lcd.ready || lcd.err != nil {
		return lcd.err
	}

	r := make([]byte, 1)
	return lcd.i2c.Tx(lcd.Address, []byte{w}, r)
}

// Turn the (optional) backlight off/on
func (lcd *LCD) NoBacklight() {
	lcd.backlightval = LCD_NOBACKLIGHT
	lcd.expanderWrite(0)
}

func (lcd *LCD) Backlight() {
	lcd.backlightval = LCD_BACKLIGHT
	lcd.expanderWrite(0)
}

/*********** mid level commands, for sending data/cmds */

func (lcd *LCD) command(value byte) {
	if lcd.ready && lcd.err == nil {
		lcd.send(value, 0)
	}
}

/************ low level data pushing commands **********/

// write either command or data
func (lcd *LCD) send(value byte, mode byte) {
	highnib := value & 0xf0
	lownib := (value << 4) & 0xf0
	lcd.write4bits(highnib | mode)
	lcd.write4bits(lownib | mode)
}

func (lcd *LCD) write4bits(value byte) {
	lcd.expanderWrite(value)
	lcd.pulseEnable(value)
}

func (lcd *LCD) expanderWrite(_data byte) {
	if !lcd.ready || lcd.err != nil {
		return
	}

	r := make([]byte, 1)
	lcd.i2c.Tx(lcd.Address<<1, []byte{_data | lcd.backlightval}, r)
}

func (lcd *LCD) pulseEnable(_data byte) {
	lcd.expanderWrite(_data | En)    // En high
	time.Sleep(time.Microsecond * 1) // enable pulse must be >450ns

	lcd.expanderWrite(_data & (En ^ 255)) // En low
	time.Sleep(time.Microsecond * 50)     // commands need > 37us to settle
}

// Alias functions

func (lcd *LCD) Cursor_on() {
	lcd.Cursor()
}

func (lcd *LCD) Cursor_off() {
	lcd.NoCursor()
}

func (lcd *LCD) blink_on() {
	lcd.Blink()
}

func (lcd *LCD) blink_off() {
	lcd.NoBlink()
}

// func (lcd *LCD) load_custom_character(char_num byte,  *rows byte) {
//     lcd.CreateChar(char_num, rows);
// }

func (lcd *LCD) SetBacklight(new_val byte) {
	if new_val != 0 {
		lcd.Backlight() // turn backlight on
	} else {
		lcd.NoBacklight() // turn backlight off
	}
}

func (lcd *LCD) printstr(s string) {
	//This function is not identical to the function used for "real" I2C displays
	//it's here so the user sketch doesn't have to be changed
	c := []byte(s)
	lcd.Print(c)
}

// unsupported API functions
func (lcd *LCD) Off() {
}

func (lcd *LCD) On() {
}

func (lcd *LCD) SetDelay(cmdDelay byte, charDelay byte) {
}

func (lcd *LCD) Status() byte {
	return 0
}

func (lcd *LCD) Keypad() byte {
	return 0
}

func (lcd *LCD) Init_bargraph(graphtype byte) byte {
	return 0
}

func (lcd *LCD) Draw_horizontal_graph(row byte, column byte, len byte, pixel_col_end byte) {
}

func (lcd *LCD) Draw_vertical_graph(row byte, column byte, len byte, pixel_row_end byte) {
}

func (lcd *LCD) SetContrast(new_val byte) {
}
