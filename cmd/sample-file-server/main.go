package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	PORT        = "port"
	FILE_URL    = "fileurl"
	FILE_LENGTH = "filelength"
)

func main() {
	app := &cli.App{
		Name:  "sample-file-server",
		Usage: "Runs a simple http server that serves a single file.",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    PORT,
				Usage:   "Specifies the port to listen on.",
				Value:   8088,
				Aliases: []string{"p"},
			},

			&cli.UintFlag{
				Name:    FILE_LENGTH,
				Usage:   "Specifies the length of the file to serve.",
				Value:   100,
				Aliases: []string{"l"},
			},
		},
		Action: func(c *cli.Context) error {
			const (
				fileName = "test.txt"
			)

			fileContent := generateFileData(c.Int(FILE_LENGTH))
			filePath, cleanup := setupFile(fileName, fileContent)
			defer cleanup()

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Set the Content-Disposition header to suggest a filename
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

				// Add CORS headers to allow all origins (*).
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Expose-Headers", "*")

				http.ServeFile(w, r, filePath)
			})

			fmt.Printf("Serve listening on http://localhost:%d%s\n", c.Uint(PORT), c.String(FILE_URL))
			err := http.ListenAndServe(fmt.Sprintf(":%d", c.Uint(PORT)), nil)
			if err != nil {
				fmt.Printf("Error starting server: %s", err)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	waitForKillSignal()
}

// waitForKillSignal blocks until we receive a kill or interrupt signal
func waitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}

// setupFile creates a file with the given name and content, and returns a cleanup function
func setupFile(fileName string, fileContent string) (string, func()) {
	dataFolder, err := os.MkdirTemp("", "sample-file-server-*")
	if err != nil {
		panic(err)
	}
	filePath := fmt.Sprintf("%s/%s", dataFolder, fileName)
	// Open the file for writing (create or truncate)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		os.Remove(filePath)
		panic(err)
	}
	return filePath, func() {
		err := os.Remove(fileName)
		if err != nil {
			panic(err)
		}
	}
}

// generateFileData generates a string of the given length composed of random words
func generateFileData(length int) (fileData string) {
	if length < 10 {
		panic("file length must be at least 10")
	}
	wordSelection := []string{
		"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel",
		"India", "Juliet", "Kilo", "Lima", "Mike", "November", "Oscar", "Papa",
		"Quebec", "Romeo", "Sierra", "Tango", "Uniform", "Victor", "Whiskey",
		"X-ray", "Yankee", "Zulu",
	}
	fileData = "START"
	// Continue adding words until we reach the desired length or beyond
	for len(fileData) < length {
		randomIndex := rand.Intn(len(wordSelection))
		fileData = fileData + " " + wordSelection[randomIndex]
	}

	return fileData[:length-3] + "END"
}
