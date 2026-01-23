package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func colorizeSameLine(color Color, message string) {
	fmt.Print(string(color), message, string(ColorReset))
}
type Opt struct {
	Color Color
}
// LogTrace logs only value with color
func LogTrace(name string, output interface{}) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(ColorRed, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogTraceN  allows caller number to be specified
func LogTraceN(name string, output interface{}, num int) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(ColorRed, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogTraceR logs & return file & funciton Name
func LogTraceR(name string, output interface{}, num int) (file, function string) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(ColorRed, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	fileName := filepath.Base(frame.File)
	functionName := filepath.Base(frame.Function)
	return fmt.Sprintf("%s:%d", fileName, frame.Line), functionName

}

// LogColor allows color of logs to be specified
func LogColor(name string, output interface{}, color Color) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(color, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogVal logs only value with color, no lines
func LogVal(name string, output interface{}, color Color) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(color, o)

}
