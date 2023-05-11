package device

import (
	"fmt"
	"machine"
	"strings"
	"time"
)

const (
	// Commands
	DISPLAY_CLEAR        = 0x01
	CURSOR_HOME          = 0x02
	ENTRY_MODE           = 0x04
	DISPLAY_ON_OFF       = 0x08
	CURSOR_DISPLAY_SHIFT = 0x10
	FUNCTION_MODE        = 0x20
	CGRAM_SET            = 0x40
	DDRAM_SET            = 0x80

	// flags for display entry mode
	CURSOR_INCREASE = 0x02
	CURSOR_DECREASE = 0x00
	DISPLAY_SHIFT   = 0x01

	// flags for display on/off control
	CURSOR_BLINK_OFF = 0x00
	CURSOR_BLINK_ON  = 0x01
	CURSOR_OFF       = 0x00
	CURSOR_ON        = 0x02
	DISPLAY_ON       = 0x04
	DISPLAY_OFF      = 0x00

	// flags for display/cursor shift
	CURSOR_MOVE_LEFT  = 0x00
	CURSOR_MOVE_RIGHT = 0x04
	CURSOR_MOVE       = 0x08

	// flags for function set
	DATA_LENGTH_4BIT = 0x00
	DATA_LENGTH_8BIT = 0x10

	ONE_LINE = 0x00
	TWO_LINE = 0x08

	FONT_5X8  = 0x00
	FONT_5X10 = 0x04

	// flags for Backlight control
	BACKLIGHT_ON  = 0x08
	BACKLIGHT_OFF = 0x00
)

const (
	Rs byte = 0x01 // Register select bit
	Rw byte = 0x02 // Read/Write bit
	En byte = 0x04 // Enable bit
)

type lcdI2CType int

const (
	LCD_UNKNOWN lcdI2CType = iota
	LCD_16x2
	LCD_20x4
)

type ShowOptions int

const (
	SHOW_NO_OPTIONS ShowOptions = 0
	SHOW_LINE_1                 = 1 << iota
	SHOW_LINE_2
	SHOW_LINE_3
	SHOW_LINE_4
	SHOW_ELIPSE_IF_NOT_FIT
	SHOW_BLANK_PADDING
)

type LcdI2C struct {
	i2c            machine.I2C
	addr           uint16
	lcdI2CType     lcdI2CType
	setup          byte
	displaycontrol byte
	displaymode    byte
	backlightval   byte
}

func NewLcdI2C(i2c machine.I2C, addr uint16, lcdI2CType lcdI2CType) (*LcdI2C, error) {
	this := &LcdI2C{
		i2c:            i2c,
		addr:           addr,
		lcdI2CType:     lcdI2CType,
		setup:          TWO_LINE | FONT_5X8 | DATA_LENGTH_4BIT,
		displaycontrol: DISPLAY_ON | CURSOR_OFF | CURSOR_BLINK_OFF,
		displaymode:    CURSOR_INCREASE | CURSOR_DECREASE,
		backlightval:   BACKLIGHT_OFF,
	}

	initByteSeq := []byte{
		0x03, 0x03, 0x03, // base initialization
		0x02, // setting up 4-bit transfer mode
		FUNCTION_MODE | this.setup,
		DISPLAY_ON_OFF | this.displaycontrol,
		ENTRY_MODE | this.displaymode,
	}

	for _, b := range initByteSeq {
		if err := this.writeByte(b, 0); err != nil {
			return nil, err
		}
	}

	if err := this.Clear(); err != nil {
		return nil, err
	}

	if err := this.Home(); err != nil {
		return nil, err
	}

	return this, nil
}

type rawData struct {
	Data  byte
	Delay time.Duration
}

func (this *LcdI2C) writeRawDataSeq(seq []rawData) error {
	r := make([]byte, 1)

	for _, item := range seq {
		if err := this.i2c.Tx(this.addr, []byte{item.Data}, r); err != nil {
			return err
		}

		time.Sleep(item.Delay)
	}

	return nil
}

func (this *LcdI2C) writeDataWithStrobe(data byte) error {
	data |= this.backlightval

	seq := []rawData{
		{data, 0},                           // send data
		{data | En, 200 * time.Microsecond}, // set strobe
		{data, 30 * time.Microsecond},       // reset strobe
	}

	return this.writeRawDataSeq(seq)
}

func (this *LcdI2C) writeByte(data byte, controlPins byte) error {
	if err := this.writeDataWithStrobe(data&0xF0 | controlPins); err != nil {
		return err
	}

	return this.writeDataWithStrobe((data<<4)&0xF0 | controlPins)
}

func (this *LcdI2C) getLineRange(options ShowOptions) (startLine, endLine int) {
	var lines [4]bool

	lines[0] = options&SHOW_LINE_1 != 0
	lines[1] = options&SHOW_LINE_2 != 0
	lines[2] = options&SHOW_LINE_3 != 0
	lines[3] = options&SHOW_LINE_4 != 0

	startLine = -1

	for i := 0; i < len(lines); i++ {
		if lines[i] {
			startLine = i
			break
		}
	}

	endLine = -1

	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] {
			endLine = i
			break
		}
	}

	return startLine, endLine
}

func (this *LcdI2C) splitText(text string, options ShowOptions) []string {
	var lines []string
	startLine, endLine := this.getLineRange(options)
	w, _ := this.getSize()

	if w != -1 && startLine != -1 && endLine != -1 {
		for i := 0; i <= endLine-startLine; i++ {
			if len(text) == 0 {
				break
			}

			j := w

			if j > len(text) {
				j = len(text)
			}

			lines = append(lines, text[:j])
			text = text[j:]
		}

		if len(text) > 0 {
			if options&SHOW_ELIPSE_IF_NOT_FIT != 0 {
				j := len(lines) - 1
				lines[j] = lines[j][:len(lines[j])-1] + "~"
			}
		} else {
			if options&SHOW_BLANK_PADDING != 0 {
				j := len(lines) - 1
				lines[j] = lines[j] + strings.Repeat(" ", w-len(lines[j]))

				for k := j + 1; k <= endLine-startLine; k++ {
					lines = append(lines, strings.Repeat(" ", w))
				}
			}
		}
	} else if len(text) > 0 {
		lines = append(lines, text)
	}

	return lines
}

func (this *LcdI2C) command(cmd byte) error {
	return this.writeByte(cmd, 0)
}

func (this *LcdI2C) ShowMessage(text string, options ShowOptions) error {
	lines := this.splitText(text, options)
	startLine, endLine := this.getLineRange(options)
	i := 0

	for {
		if startLine != -1 && endLine != -1 {
			if err := this.SetPosition(i+startLine, 0); err != nil {
				return err
			}
		}

		line := lines[i]

		for _, c := range line {
			if err := this.writeByte(byte(c), Rs); err != nil {
				return err
			}
		}

		if i == len(lines)-1 {
			break
		}

		i++
	}

	return nil
}

func (this *LcdI2C) TestWriteCGRam() error {
	if err := this.writeByte(CGRAM_SET, 0); err != nil {
		return err
	}

	var a byte = 0x55

	for i := 0; i < 80; i++ {
		if err := this.writeByte(a, Rs); err != nil {
			return err
		}

		a = a ^ 0xFF
	}

	return nil
}

func (this *LcdI2C) BacklightOn() error {
	this.backlightval = BACKLIGHT_ON
	return this.writeByte(0x00, 0)
}

func (this *LcdI2C) BacklightOff() error {
	this.backlightval = BACKLIGHT_OFF
	return this.writeByte(0x00, 0)
}

func (this *LcdI2C) Clear() error {
	err := this.writeByte(DISPLAY_CLEAR, 0)
	time.Sleep(3 * time.Millisecond)
	return err
}

func (this *LcdI2C) Home() error {
	err := this.writeByte(CURSOR_HOME, 0)
	time.Sleep(3 * time.Millisecond)
	return err
}

func (this *LcdI2C) getSize() (width, height int) {
	switch this.lcdI2CType {
	case LCD_16x2:
		return 16, 2
	case LCD_20x4:
		return 20, 4
	default:
		return -1, -1
	}
}

func (this *LcdI2C) SetPosition(line, pos int) error {
	cols, rows := this.getSize()

	if cols != -1 && (pos < 0 || pos > cols-1) {
		return fmt.Errorf("Cursor position %d must be within the range [0..%d]", pos, cols-1)
	}

	if rows != -1 && (line < 0 || line > rows-1) {
		return fmt.Errorf("Cursor line %d must be within the range [0..%d]", line, rows-1)
	}

	var lineOffset []byte

	switch cols {
	case 16:
		lineOffset = []byte{0x00, 0x40}
	case 20:
		lineOffset = []byte{0x00, 0x40, 0x14, 0x54}
	default:
		return nil
	}

	var b byte = DDRAM_SET + lineOffset[line] + byte(pos)
	return this.writeByte(b, 0)
}

func (this *LcdI2C) Write(buf []byte) (int, error) {
	for i, c := range buf {
		if err := this.writeByte(c, Rs); err != nil {
			return i, err
		}
	}

	return len(buf), nil
}

func (this *LcdI2C) DisplayOn() error {
	this.displaycontrol |= DISPLAY_ON
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

func (this *LcdI2C) DisplayOff() error {
	this.displaycontrol &= DISPLAY_ON ^ 0xFF
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

func (this *LcdI2C) CursorOn() error {
	this.displaycontrol |= CURSOR_ON
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

func (this *LcdI2C) CursorOff() error {
	this.displaycontrol &= CURSOR_ON ^ 0xFF
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

func (this *LcdI2C) BlinkOn() error {
	this.displaycontrol |= CURSOR_BLINK_ON
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

func (this *LcdI2C) BlinkOff() error {
	this.displaycontrol &= CURSOR_BLINK_ON ^ 0xFF
	return this.command(DISPLAY_ON_OFF | this.displaycontrol)
}

// These commands scroll the display without changing the RAM
func (this *LcdI2C) ScrollDisplayLeft() error {
	return this.command(DISPLAY_ON_OFF | CURSOR_MOVE | CURSOR_MOVE_LEFT)
}

func (this *LcdI2C) ScrollDisplayRight() error {
	return this.command(DISPLAY_ON_OFF | CURSOR_MOVE | CURSOR_MOVE_RIGHT)
}

// This is for text that flows Left to Right
func (this *LcdI2C) LeftToRight() error {
	this.displaymode |= CURSOR_INCREASE
	return this.command(ENTRY_MODE | this.displaymode)
}

// This is for text that flows Right to Left
func (this *LcdI2C) RightToLeft() error {
	this.displaymode &= CURSOR_INCREASE ^ 0xFF
	return this.command(ENTRY_MODE | this.displaymode)
}

// This will 'right justify' text from the cursor
func (this *LcdI2C) Autoscroll() error {
	this.displaymode |= DISPLAY_SHIFT
	return this.command(ENTRY_MODE | this.displaymode)
}

// This will 'left justify' text from the cursor
func (this *LcdI2C) NoAutoscroll() error {
	this.displaymode &= DISPLAY_SHIFT ^ 0xFF
	return this.command(ENTRY_MODE | this.displaymode)
}

// Allows us to fill the first 8 CGRAM locations
// with custom characters
func (this *LcdI2C) CreateChar(location byte, charmap []byte) error {
	location &= 0x7 // we only have 8 locations 0-7

	if err := this.command(CGRAM_SET | (location << 3)); err != nil {
		return err
	}

	_, err := this.Write(charmap)
	return err
}

func (this *LcdI2C) SetBacklight(state bool) error {
	if state {
		return this.BacklightOn() // turn backlight on
	}

	return this.BacklightOff() // turn backlight off
}
