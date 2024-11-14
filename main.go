// cmd/main.go
package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Create directories if they don't exist
	if err := os.MkdirAll("data/images", 0755); err != nil {
		log.Fatal("Failed to create images directory:", err)
	}
	if err := os.MkdirAll("data/text", 0755); err != nil {
		log.Fatal("Failed to create text directory:", err)
	}

	// Open webcam
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatal("Failed to open webcam:", err)
	}
	defer webcam.Close()

	// Create a window to display the image
	window := gocv.NewWindow("Webcam")
	defer window.Close()

	// Create a Mat to store the image
	img := gocv.NewMat()
	defer img.Close()

	fmt.Println("Press 'space' to capture an image or 'q' to quit")

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("Failed to read from webcam")
			break
		}
		if img.Empty() {
			continue
		}

		window.IMShow(img)
		key := window.WaitKey(1)

		// Press 'space' to capture
		if key == 32 { // 32 is spacebar
			timestamp := time.Now().Format("20060102-150405")
			imagePath := filepath.Join("data/images", fmt.Sprintf("capture_%s.jpg", timestamp))
			textPath := filepath.Join("data/text", fmt.Sprintf("text_%s.txt", timestamp))

			// Save the image
			if ok := gocv.IMWrite(imagePath, img); !ok {
				log.Println("Failed to save image")
				continue
			}
			fmt.Printf("Image saved to: %s\n", imagePath)

			// Perform OCR
			client := gosseract.NewClient()
			defer client.Close()

			client.SetImage(imagePath)
			text, err := client.Text()
			if err != nil {
				log.Printf("OCR failed: %v\n", err)
				continue
			}

			// Save the OCR text
			if err := os.WriteFile(textPath, []byte(text), 0644); err != nil {
				log.Printf("Failed to save text: %v\n", err)
				continue
			}
			fmt.Printf("OCR text saved to: %s\n", textPath)
			fmt.Printf("Extracted text:\n%s\n", text)
		}

		// Press 'q' to quit
		if key == 113 { // 113 is 'q'
			break
		}
	}
}
