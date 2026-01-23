package logger

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed    Color = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"

	ColorMagenta Color = "\u001b[35m"
	ColorCyan    Color = "\u001b[36m"
	ColorWhite   Color = "\u001b[37m"
	ColorReset         = "\u001b[0m"
)

// BRIGHT COLORS
const (
	ColorBrightBlack   Color = "\u001b[90m"
	ColorBrightRed     Color = "\u001b[91m"
	ColorBrightGreen   Color = "\u001b[92m"
	ColorBrightYellow  Color = "\u001b[93m"
	ColorBrightBlue    Color = "\u001b[94m"
	ColorBrightMagenta Color = "\u001b[95m"
	ColorBrightCyan    Color = "\u001b[96m"
	ColorBrightWhite   Color = "\u001b[97m"
)

// Back Ground Colors
const (
	BgBlack   Color = "\u001b[40m"
	BgRed     Color = "\u001b[41m"
	BgGreen   Color = "\u001b[42m"
	BgYellow  Color = "\u001b[43m"
	BgBlue    Color = "\u001b[44m"
	BgMagenta Color = "\u001b[45m"
	BgCyan    Color = "\u001b[46m"
	BgWhite   Color = "\u001b[47m"
)
