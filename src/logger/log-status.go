package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// LogTrace logs only value with color
func LogInfo(name string, output any) {
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(ColorCyan, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogTrace logs only value with color
func LogError(name string, output interface{}) {
	// PrintUi("=", 40, false)
	colorizeSameLine(BgGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(ColorBrightRed, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogTraceN logs only value with color
func LogTraceNOp(name string, output any, num int, opt Opt) {
	color := ColorRed
	if opt.Color != "" {
		color = opt.Color
	}
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(color, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	// PrintUi("=", 40, false)

}

// LogTraceR logs only value with color
func LogTraceRop(name string, output interface{}, num int, opt Opt) (file, function string) {
	color := ColorRed
	if opt.Color != "" {
		color = opt.Color
	}
	// PrintUi("=", 40, false)
	colorizeSameLine(ColorGreen, name)
	o := fmt.Sprintf("=| %v", output)
	colorize(color, o)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	fileName := filepath.Base(frame.File)
	functionName := filepath.Base(frame.Function)
	return fmt.Sprintf("%s:%d", fileName, frame.Line), functionName
	// PrintUi("=", 40, false)

}
