package logging

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var once sync.Once

func ConfigureZeroLogger() {
	once.Do(configureZeroLogger)
}

func configureZeroLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
}
// NewLogWriter returns a writer for the given logDir and logFile
// If the log file already exists it will be removed and a fresh file will be created
func NewLogWriter(logDir, logFile string) *os.File {
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
