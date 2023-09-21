package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
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
				fileName = "test.png"
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
func setupFile(fileName string, fileContent *image.RGBA) (string, func()) {
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

	if err := png.Encode(file, fileContent); err != nil {
		fmt.Println("Failed to encode image:", err)
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
func generateFileData(length int) (img *image.RGBA) {
	// Define image dimensions
	width, height := length, length

	// Create an empty RGBA image
	img = image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with a gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Gradient: Horizontal red, Vertical blue
			c := color.RGBA{
				R: uint8(x * 255 / width),
				B: uint8(y * 255 / height),
				G: 0,
				A: 255, // Fully opaque
			}
			img.Set(x, y, c)
		}
	}

	return img
}
