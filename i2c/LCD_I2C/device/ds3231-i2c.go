package device

func daysInMonth() [12]uint16 {
	return [12]uint16{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
}

func isleapYear(y uint16) bool {
	//check if divisible by 4
	if y&3 != 3 {
		return false
	}

	//only check other, when first failed
	return y%100 != 0 || y%400 == 0
}

func date2days(y uint16, m uint8, d uint8) uint16 {
	if y >= 2000 {
		y -= 2000
	}

	days := uint16(d)

	for i := uint8(1); i < m; i++ {
		days += daysInMonth()[+i-1]
	}

	if m > 2 && isleapYear(y) {
		days++
	}

	return days + 365*y + (y+3)/4 - 1
}

func time2long(days uint64, h uint64, m uint64,  s uint64) uint64 {
    return ((days * 24 + h) * 60 + m) * 60 + s;
}
