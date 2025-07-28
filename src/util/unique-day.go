package util

import (
	"fmt"
	"time"
)

var base62Digits = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', // 0-9
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', // 10-19
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', // 20-29
	'u', 'v', 'w', 'x', 'y', 'z', // 30-45
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', // 46-55
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', // 56-65
	'U', 'V', 'W', 'X', 'Y', 'Z', // 66-71
}

// toBase62 converts a number to the base-62 number system with a minimum number of digits
func toBase62(num int64) string {
	if num == 0 {
		return string(base62Digits[0])
	}
	digits := []rune{}
	for num > 0 {
		digits = append([]rune{base62Digits[num%62]}, digits...)
		num /= 62
	}
	// // Pad with leading zeros if needed to meet minDigits
	// for len(digits) < minDigits {
	// 	digits = append([]rune{'0'}, digits...)
	// }
	if len(digits) == 0 {
		digits = []rune{'0'}
	}
	return string(digits)
}
func DateTimeToBase62(dateStr string) (string, error) {
	// Parse the input date
	inputDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	// Get last two digits of year, month, day
	year := inputDate.Year() % 100  // e.g., 2025 -> 25
	month := int(inputDate.Month()) // 1-12
	day := inputDate.Day()          // 1-31

	// Validate date components
	if year < 0 || year > 99 || month < 1 || month > 12 || day < 1 || day > 31 {
		return "", fmt.Errorf("invalid date components: year=%d, month=%d, day=%d", year, month, day)
	}

	// Get current hour and minute in EAT (UTC+3)
	eat := time.FixedZone("EAT", 3*3600)
	now := time.Now().In(eat)
	hour := now.Hour()     // 0-23
	minute := now.Minute() // 0-59

	// Validate time components
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return "", fmt.Errorf("invalid time components: hour=%d, minute=%d", hour, minute)
	}

	// Combine year and month: year*12 + month
	yearMonth := year*12 + month // 0*12 + 1 to 99*12 + 12 = 1 to 1200

	// Combine day and hour: (day + hour) % 62
	dayHour := (day + hour) % 62 // 1+0 to 31+23 = 1 to 54, mod 62 ensures 0-61

	// Convert to base-62
	yearMonthStr := toBase62(int64(yearMonth)) // 2 digits
	dayHourStr := toBase62(int64(dayHour))     // 1 digit
	minuteStr := toBase62(int64(minute))

	// Combine into a single string
	result := yearMonthStr + dayHourStr + minuteStr

	return result, nil
}
func DateTimeUniqueToBase62(dateStr string) (string, error) {
	// Parse the input date
	inputDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	// Get last two digits of year, month, day
	year := inputDate.Year() % 100  // e.g., 2025 -> 25
	month := int(inputDate.Month()) // 1-12
	day := inputDate.Day()          // 1-31

	// Validate date components
	if year < 0 || year > 99 || month < 1 || month > 12 || day < 1 || day > 31 {
		return "", fmt.Errorf("invalid date components: year=%d, month=%d, day=%d", year, month, day)
	}

	// Get current hour and minute in EAT (UTC+3)
	eat := time.FixedZone("EAT", 3*3600)
	now := time.Now().In(eat)
	hour := now.Hour()     // 0-23
	minute := now.Minute() // 0-59

	// Validate time components
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return "", fmt.Errorf("invalid time components: hour=%d, minute=%d", hour, minute)
	}

	// Combine into a single number for unique mapping
	combined := int64((((year*12+(month-1))*31+(day-1))*24+hour)*60 + minute)

	// Convert to base-62, pad to 5 digits for consistency
	result := toBase62(combined)

	return result, nil
}
