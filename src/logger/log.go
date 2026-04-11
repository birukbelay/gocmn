package logger

import (
	"fmt"
	"log/slog"
	"runtime"
)

// InfoCtx logs info with context
func Info(name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorCyan, ColorBlue, nil)

	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("val", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.InfoContext(ctx, coloredName, slogArgs...)
}

// ErrorCtx logs error with context
func Error(name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorMagenta, ColorYellow, nil)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("val", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.ErrorContext(ctx, coloredName, slogArgs...)

}

// WarnCtx logs warning with context
func Warn(name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorYellow, ColorMagenta, nil)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("val", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.WarnContext(ctx, coloredName, slogArgs...)

}

// DebugCtx logs debug with context
func Debug(name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorBlue, ColorMagenta, nil)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("val", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.DebugContext(ctx, coloredName, slogArgs...)

}
