package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
)

var fileName = "image.png"

type Pixel struct {
	R, G, B, A int
}

func main() {
	if len(os.Args) < 3 || (os.Args[1] != "-w" && os.Args[1] != "-r" && os.Args[1] != "-f") || (os.Args[1] == "-w" && len(os.Args) != 4) || (os.Args[1] == "-f" && len(os.Args) != 4) {
		fmt.Println("Invalid Usage")
		fmt.Println("Usage: go run main.go -w fileName \"message here\" ")
		fmt.Println("Usage: go run main.go -f fileName \"file name\" ")
		fmt.Println("Usage: go run main.go -r fileName ")
		os.Exit(1)
	}
	mode := os.Args[1]
	fileName = os.Args[2]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("Image file doesn't exist")
		os.Exit(1)
	}

	if mode == "-w" {
		write(os.Args[3])
	} else if mode == "-r" {
		read()
	} else if mode == "-f" {
		f, err := ioutil.ReadFile(os.Args[3])
		if err != nil {
			fmt.Println("Input text file doesn't exist")
			log.Fatal(err)
		}
		write(string(f))
	}
}

func write(input string) {
	// TESTING STRING TO BINARY
	newString := convertString(input)
	bitLength := i2b(len(newString), 32)

	// Load fileName and create a placeholder with same dimensions
	pixelArray, width, height := getImage()
	if (height * width * 3) < (len(newString)*8)+11 {
		log.Fatal("Message is too long for given image file.\nTry with a bigger image.")
	}
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Write the length of the message
	writeBits(0, pixelArray, bitLength)
	writeBits(11, pixelArray, newString)

	fmt.Printf("Length of message: %d characters\n", len(input))
	fmt.Println("Number of pixels being used to write message: ", int(math.Ceil((float64(len(newString)) / 3))))

	// Copying pixelArray to placeholder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(pixelArray[(y*width)+(x)].R),
				G: uint8(pixelArray[(y*width)+(x)].G),
				B: uint8(pixelArray[(y*width)+(x)].B),
				A: 255,
			})
		}
	}

	fmt.Println("Writing to file...")
	// Writing the placeholder to a png file
	output := fileName[:len(fileName)-3] + "png"
	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	fmt.Println("Done writing to file")
}

func writeBits(offset int, pixelArray []Pixel, source []byte) []Pixel {
	// Takes input (source) and writes to pixelArray starting at a given offset
	numPix := int(math.Ceil(float64(len(source)) / 3))
	arrayLen := len(pixelArray)
	sourceLen := len(source)

	for i := 0; i < numPix; i++ {
		position := arrayLen - 1 - offset - i
		if int(source[i*3]) == 1 {
			pixelArray[position].R = pixelArray[position].R | int(source[i*3])
		} else {
			pixelArray[position].R = pixelArray[position].R & 254
		}

		if (i*3)+1 >= sourceLen {
			break
		}

		if int(source[(i*3)+1]) == 1 {
			pixelArray[position].G = pixelArray[position].G | int(source[(i*3)+1])
		} else {
			pixelArray[position].G = pixelArray[position].G & 254
		}

		if (i*3)+2 >= sourceLen {
			break
		}

		if int(source[(i*3)+2]) == 1 {
			pixelArray[position].B = pixelArray[position].B | int(source[(i*3)+2])
		} else {
			pixelArray[position].B = pixelArray[position].B & 254
		}
	}
	return pixelArray
}

func read() {
	// Open the file and get the pixel array
	pixelArray, _, _ := getImage()
	// Getting stored message length, including extra unused bit
	var leng []byte
	for i := range pixelArray[len(pixelArray)-11:] {
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].R&1))
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].G&1))
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].B&1))
	}
	// Convert from binary to decimal
	messageLength := b2d(leng[:len(leng)-1])
	// Retreive the message pixels
	var tmpMessage []byte
	for i := range pixelArray[(len(pixelArray) - (11 + int(math.Ceil(messageLength/3)))) : len(pixelArray)-11] {
		tmpMessage = append(tmpMessage, byte(pixelArray[len(pixelArray)-i-12].R&1))
		tmpMessage = append(tmpMessage, byte(pixelArray[len(pixelArray)-i-12].G&1))
		tmpMessage = append(tmpMessage, byte(pixelArray[len(pixelArray)-i-12].B&1))
	}
	// Retreive their bits
	var message []byte
	for i := 0; i < int(messageLength/8); i++ {
		var char byte = 0
		for j := 0; j < 8; j++ {
			char = (char << 1) + tmpMessage[(i*8)+j]
		}
		message = append(message, char)
	}
	fmt.Printf("Message size: %d characters\n", len(message))
	fmt.Println("Message:\n", string(message))
}

func getImage() ([]Pixel, int, int) {
	// Opens the image and return its pixels and dimensions
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileType := fileName[len(fileName)-3:] // Get the file extension
	var loadedImage image.Image
	if fileType == "jpg" { // Input is a jpg
		loadedImage, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	} else if fileType == "png" { // Input is a png
		loadedImage, _, err = image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unsuported image file type")
	}

	height := loadedImage.Bounds().Dy()
	width := loadedImage.Bounds().Dx()
	var pixels []Pixel
	// Gather all the pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, toPixel(loadedImage.At(x, y).RGBA()))
		}
	}
	return pixels, width, height
}

func toPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	// Takes in a golang pixel in rgba format and converts all the values from uint32s to normal integers
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func b2d(input []byte) float64 {
	// Takes in a binary representation and returns a decimal equivalent
	result := 0.0
	for i, val := range input {
		power := len(input) - 1 - i
		if val == 1 {
			result += math.Pow(2, float64(power))
		}
	}
	return result
}

func i2b(input int, numBytes int) []byte {
	// Converts an integer to binary
	var result []byte
	for i := 0; i < numBytes; i++ {
		result = append(result, 0)
	}
	for i := 0; i < numBytes; i++ {
		if (input & 1) == 1 {
			result[numBytes-1-i] = 1
		}
		input = input >> 1
	}
	return result
}

func convertString(input string) []byte {
	var temp []byte
	for _, char := range input {
		for _, bytes := range i2b(int(char), 8) {
			// Convert each character's ord value to bits
			temp = append(temp, bytes)
		}
	}
	// fmt.Println(temp)
	return temp
}
