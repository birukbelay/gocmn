package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime"
)

func GetDef(opt *Opt, keyColor, valColor Color, ctx context.Context) (num int, conxt context.Context, keyClr Color, valClr Color) {
	num = 2
	if ctx == nil {
		ctx = context.Background()
	}

	if opt != nil {
		if opt.Num != 0 {
			num = opt.Num
		}
		if opt.Ctx != nil {
			ctx = opt.Ctx
		}
		if opt.Color != "" {
			keyColor = opt.Color
		}
		if opt.ValCol != "" {
			valColor = opt.ValCol
		}
	}
	return num, ctx, keyColor, valColor
}

// InfoCtx logs info with context
func InfoCtx(ctx context.Context, name string, output any, opt *Opt) {
	num, _, keyColor, valColor := GetDef(opt, ColorCyan, ColorBlue, ctx)

	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("output", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.InfoContext(ctx, coloredName, slogArgs...)
	log.Printf("%s:%d\n", frame.File, frame.Line)
}

// ErrorCtx logs error with context
func ErrorCtx(ctx context.Context, name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorMagenta, ColorYellow, ctx)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("output", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.ErrorContext(ctx, coloredName, slogArgs...)
	log.Printf("%s:%d\n", frame.File, frame.Line)
}

// WarnCtx logs warning with context
func WarnCtx(ctx context.Context, name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorYellow, ColorMagenta, ctx)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("output", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.WarnContext(ctx, coloredName, slogArgs...)
	log.Printf("%s:%d\n", frame.File, frame.Line)
}

// DebugCtx logs debug with context
func DebugCtx(ctx context.Context, name string, output any, opt *Opt) {
	num, ctx, keyColor, valColor := GetDef(opt, ColorBlue, ColorMagenta, ctx)
	pc := make([]uintptr, 15)
	n := runtime.Callers(num, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	coloredName := keyColor.S() + name + string(ColorReset)
	coloredOutput := valColor.S() + fmt.Sprintf("=| %v", output) + string(ColorReset)

	var slogArgs []any
	slogArgs = append(slogArgs, slog.String("output", coloredOutput), slog.String("file", fmt.Sprintf("%s:%d", frame.File, frame.Line)))
	if opt != nil && len(opt.Args) > 0 {
		slogArgs = append(slogArgs, opt.Args)
	}

	slog.DebugContext(ctx, coloredName, slogArgs...)
	log.Printf("%s:%d\n", frame.File, frame.Line)
}
