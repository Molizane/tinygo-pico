package device

import (
	"fmt"
	"machine"
	"strings"
	"time"
)

const (
	// Commands
	CMD_Clear_Display   = 0x01
	CMD_Return_Home     = 0x02
	CMD_Display_Entry   = 0x04
	CMD_Display_Control = 0x08
	CMD_Cursor_Shift    = 0x10
	CMD_Setup           = 0x20
	CMD_CGRAM_Set       = 0x40
	CMD_DDRAM_Set       = 0x80

	// flags for display entry mode
	FLG_Entry_Left            = 0x02
	FLG_Entry_Shift_Decrement = 0x00
	FLG_Entry_Shift_Increment = 0x01

	// flags for display on/off control
	FLG_Blink_Off  = 0x00
	FLG_Blink_On   = 0x01
	FLG_Cursor_Off = 0x00
	FLG_Cursor_On  = 0x02
	FLG_Display_On = 0x04

	// flags for display/cursor shift
	FLG_Cursor_Move  = 0x00
	FLG_Move_Left    = 0x00
	FLG_Move_Right   = 0x04
	FLG_Display_Move = 0x08

	// flags for function set
	FLG_4Bit_Mode = 0x00
	FLG_8Bit_Mode = 0x10

	FLG_1_Line  = 0x00
	FLG_2_Lines = 0x08

	FLG_5x8_Dots  = 0x00
	FLG_5x10_Dots = 0x04

	// flags for Backlight control
	FLG_Backlight_On  = 0x08
	FLG_Backlight_Off = 0x00
)

const (
	PIN_RS byte = 0x01 // Register select bit
	PIN_RW byte = 0x02 // Read/Write bit
	PIN_EN byte = 0x04 // Enable bit
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
		setup:          FLG_2_Lines | FLG_5x8_Dots | FLG_4Bit_Mode,
		displaycontrol: FLG_Display_On | FLG_Cursor_Off | FLG_Blink_Off,
		displaymode:    FLG_Entry_Left | FLG_Entry_Shift_Decrement,
		backlightval:   FLG_Backlight_Off,
	}

	initByteSeq := []byte{
		0x03, 0x03, 0x03, // base initialization
		0x02, // setting up 4-bit transfer mode
		CMD_Setup | this.setup,
		CMD_Display_Control | this.displaycontrol,
		CMD_Display_Entry | this.displaymode,
	}

	for _, b := range initByteSeq {
		err := this.writeByte(b, 0)

		if err != nil {
			return nil, err
		}
	}

	err := this.Clear()

	if err != nil {
		return nil, err
	}

	err = this.Home()

	if err != nil {
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
		err := this.i2c.Tx(this.addr, []byte{item.Data}, r)

		if err != nil {
			return err
		}

		time.Sleep(item.Delay)
	}

	return nil
}

func (this *LcdI2C) writeDataWithStrobe(data byte) error {
	data |= this.backlightval

	seq := []rawData{
		{data, 0},                               // send data
		{data | PIN_EN, 200 * time.Microsecond}, // set strobe
		{data, 30 * time.Microsecond},           // reset strobe
	}

	return this.writeRawDataSeq(seq)
}

func (this *LcdI2C) writeByte(data byte, controlPins byte) error {
	err := this.writeDataWithStrobe(data&0xF0 | controlPins)

	if err != nil {
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
			err := this.SetPosition(i+startLine, 0)

			if err != nil {
				return err
			}
		}

		line := lines[i]

		for _, c := range line {
			err := this.writeByte(byte(c), PIN_RS)

			if err != nil {
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
	err := this.writeByte(CMD_CGRAM_Set, 0)

	if err != nil {
		return err
	}

	var a byte = 0x55

	for i := 0; i < 80; i++ {
		err := this.writeByte(a, PIN_RS)

		if err != nil {
			return err
		}

		a = a ^ 0xFF
	}

	return nil
}

func (this *LcdI2C) BacklightOn() error {
	this.backlightval = FLG_Backlight_On
	return this.writeByte(0x00, 0)
}

func (this *LcdI2C) BacklightOff() error {
	this.backlightval = FLG_Backlight_Off
	return this.writeByte(0x00, 0)
}

func (this *LcdI2C) Clear() error {
	return this.writeByte(CMD_Clear_Display, 0)
}

func (this *LcdI2C) Home() error {
	err := this.writeByte(CMD_Return_Home, 0)
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
	w, h := this.getSize()

	if w != -1 && (pos < 0 || pos > w-1) {
		return fmt.Errorf("Cursor position %d "+"must be within the range [0..%d]", pos, w-1)
	}

	if h != -1 && (line < 0 || line > h-1) {
		return fmt.Errorf("Cursor line %d "+"must be within the range [0..%d]", line, h-1)
	}

	lineOffset := []byte{0x00, 0x40, 0x14, 0x54}
	var b byte = CMD_DDRAM_Set + lineOffset[line] + byte(pos)
	return this.writeByte(b, 0)
}

func (this *LcdI2C) Write(buf []byte) (int, error) {
	for i, c := range buf {
		err := this.writeByte(c, PIN_RS)

		if err != nil {
			return i, err
		}
	}

	return len(buf), nil
}

func (this *LcdI2C) DisplayOn() error {
	this.displaycontrol |= FLG_Display_On
	return this.command(CMD_Display_Control | this.displaycontrol)
}

func (this *LcdI2C) DisplayOff() error {
	this.displaycontrol &= FLG_Display_On ^ 0xFF
	return this.command(CMD_Display_Control | this.displaycontrol)
}

func (this *LcdI2C) CursorOn() error {
	this.displaycontrol |= FLG_Cursor_On
	return this.command(CMD_Display_Control | this.displaycontrol)
}

func (this *LcdI2C) CursorOff() error {
	this.displaycontrol &= FLG_Cursor_On ^ 0xFF
	return this.command(CMD_Display_Control | this.displaycontrol)
}

func (this *LcdI2C) BlinkOn() error {
	this.displaycontrol |= FLG_Display_On
	return this.command(CMD_Display_Control | this.displaycontrol)
}

func (this *LcdI2C) BlinkOff() error {
	this.displaycontrol &= FLG_Display_On ^ 0xFF
	return this.command(CMD_Display_Control | this.displaycontrol)
}

// These commands scroll the display without changing the RAM
func (this *LcdI2C) ScrollDisplayLeft() error {
	return this.command(CMD_Display_Control | FLG_Display_Move | FLG_Move_Left)
}

func (this *LcdI2C) ScrollDisplayRight() error {
	return this.command(CMD_Display_Control | FLG_Display_Move | FLG_Move_Right)
}

// This is for text that flows Left to Right
func (this *LcdI2C) LeftToRight() error {
	this.displaymode |= FLG_Entry_Left
	return this.command(CMD_Display_Entry | this.displaymode)
}

// This is for text that flows Right to Left
func (this *LcdI2C) RightToLeft() error {
	this.displaymode &= FLG_Entry_Left ^ 0xFF
	return this.command(CMD_Display_Entry | this.displaymode)
}

// This will 'right justify' text from the cursor
func (this *LcdI2C) Autoscroll() error {
	this.displaymode |= FLG_Entry_Shift_Increment
	return this.command(CMD_Display_Entry | this.displaymode)
}

// This will 'left justify' text from the cursor
func (this *LcdI2C) NoAutoscroll() error {
	this.displaymode &= FLG_Entry_Shift_Increment ^ 0xFF
	return this.command(CMD_Display_Entry | this.displaymode)
}

func (this *LcdI2C) SetCursor(col, row int) error {
	row_offsets := []int{0x00, 0x40, 0x14, 0x54}
	rows, _ := this.getSize()

	if row > rows {
		row = rows - 1 // we count rows starting w/0
	}

	return this.command(CMD_DDRAM_Set | byte(col+row_offsets[row]))
}

// Allows us to fill the first 8 CGRAM locations
// with custom characters
func (this *LcdI2C) CreateChar(location byte, charmap []byte) error {
	location &= 0x7 // we only have 8 locations 0-7
	err := this.command(CMD_CGRAM_Set | (location << 3))

	if err != nil {
		return err
	}

	_, err = this.Write(charmap)
	return err
}

func (this *LcdI2C) SetBacklight(state bool) error {
	if state {
		return this.BacklightOn() // turn backlight on
	}

	return this.BacklightOff() // turn backlight off
}
