package config

import (
	"io"
	"log/slog"
	"time"

	"gabe565.com/utils/slogx"
	"gabe565.com/utils/termx"
	"github.com/lmittmann/tint"
)

// InitLog initializes the logger with the config's log level and format.
func (c *Config) InitLog(w io.Writer) {
	InitLog(w, c.logLevel, c.logFormat)
}

// InitLog initializes the global logger with the specified level and format.
func InitLog(w io.Writer, level slogx.Level, format slogx.Format) {
	switch format {
	case slogx.FormatJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.Level(level),
		})))
	default:
		var color bool
		switch format {
		case slogx.FormatAuto:
			color = termx.IsColor(w)
		case slogx.FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      slog.Level(level),
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
