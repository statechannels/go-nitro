package logging

import (
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
