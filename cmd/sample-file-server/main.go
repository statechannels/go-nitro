package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	PORT     = "port"
	FILE_URL = "fileurl"
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
			&cli.StringFlag{
				Name:    FILE_URL,
				Usage:   "Specifies the url to serve the file at.",
				Value:   "/test.txt",
				Aliases: []string{"f"},
			},
		},
		Action: func(c *cli.Context) error {
			const (
				fileName    = "test.txt"
				fileContent = "Hello world!"
			)

			cleanup := setupFile(fileName, fileContent)
			defer cleanup()

			http.HandleFunc(c.String(FILE_URL), func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, fileName)
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
func setupFile(testFileName string, testFileContent string) func() {
	// Open the file for writing (create or truncate)
	file, err := os.Create(testFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(testFileContent)
	if err != nil {
		os.Remove(testFileName)
		panic(err)
	}
	return func() {
		err := os.Remove(testFileName)
		if err != nil {
			panic(err)
		}
	}
}
