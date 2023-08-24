package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// newLogWriter returns a writer for the given logDir and logFile
// If the log file already exists it will be removed and a fresh file will be created
func newLogWriter(logDir, logFile string) *os.File {
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(logDir, logFile)
	// Clear the file
	os.Remove(filename)
	logDestination, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	return logDestination
}

const (
	LOG_DIR = "../artifacts"

	CHANNEL_ID_LOG_KEY   = "channel-id"
	OBJECTIVE_ID_LOG_KEY = "objective-id"
	ADDRESS_LOG_KEY      = "address"
)

// WithChannelIdAttribute returns a logging attribute for the given channel id
func WithChannelIdAttribute(c types.Destination) slog.Attr {
	return slog.String(CHANNEL_ID_LOG_KEY, c.String())
}

// WithObjectiveIdAttribute returns a logging attribute for the given objective id
func WithObjectiveIdAttribute(o protocols.ObjectiveId) slog.Attr {
	return slog.String(OBJECTIVE_ID_LOG_KEY, string(o))
}

// LoggerWithAddress returns a logger with the address attribute set to the given address
func LoggerWithAddress(logger *slog.Logger, a types.Address) *slog.Logger {
	return logger.With(slog.String(ADDRESS_LOG_KEY, a.String()))
}

// SetupDefaultFileLogger sets up a default logger that writes to the specified file
// The file will be created in the artifacts directory
func SetupDefaultFileLogger(filename string, level slog.Level) {
	logFile := newLogWriter(LOG_DIR, filename)
	SetupDefaultLogger(logFile, level)
}

// SetupDefaultLogger sets up a default logger that writes to the specified writer
func SetupDefaultLogger(w io.Writer, level slog.Level) {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(h))
}
