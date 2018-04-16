package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

// var fileName = "test.png"

var fileName = "image.png"

type Pixel struct {
	R, G, B, A int
}

func main() {
	// TESTING STRING TO BINARY
	newString := fmt.Sprintf("%08b", []byte("is fun security")) // Pad with leading 0s
	replacer := strings.NewReplacer("[", "", "]", "", " ", "")  // Stripping [, ], and whitespace
	newString = replacer.Replace(newString)
	messageLength := len(newString)
	temp := fmt.Sprintf("%b", messageLength)
	fmt.Println([]byte(temp))
	fmt.Println("\noutput: ", newString, len(newString))

	fmt.Println(i2b(len(newString)))
	bitLength := i2b(len(newString))
	fmt.Print(bitLength)
	pixelArray, width, height := getImage()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	fmt.Println("\nTesting...")
	// Have to add cases to make sure it is doing the right boolean logic for each portion... atm it's only ORing 1.
	//		This is only helpful when you are writing a 1, not when you need to write a 0 (& 0)
	for i := 0; i < 11; i++ {
		fmt.Print(" r", i, ":")
		if int(bitLength[i*3]) == 1 {
			pixelArray[len(pixelArray)-1-i].R = pixelArray[len(pixelArray)-1-i].R | int(bitLength[i*3])
			fmt.Print(1)
		} else {
			pixelArray[len(pixelArray)-1-i].R = pixelArray[len(pixelArray)-1-i].R & 254
			fmt.Print(0)
		}
		fmt.Print(" g", i, ":")
		if int(bitLength[(i*3)+1]) == 1 {
			pixelArray[len(pixelArray)-1-i].G = pixelArray[len(pixelArray)-1-i].G | int(bitLength[(i*3)+1])
			fmt.Print(1)
		} else {
			pixelArray[len(pixelArray)-1-i].G = pixelArray[len(pixelArray)-1-i].G & 254
			fmt.Print(0)
		}
		if i == 10 {
			continue
		}
		fmt.Print(" b", i, ":")

		if int(bitLength[(i*3)+2]) == 1 {
			fmt.Print(1)
			pixelArray[len(pixelArray)-1-i].B = pixelArray[len(pixelArray)-1-i].B | int(bitLength[(i*3)+2])
		} else {
			fmt.Print(0)
			pixelArray[len(pixelArray)-1-i].B = pixelArray[len(pixelArray)-1-i].B & 254
		}
		// fmt.Print(" ", pixelArray[len(pixelArray)-1-i].B&1, ":", pixelArray[len(pixelArray)-1-i].B&1, ",")
	}
	fmt.Print("\n", pixelArray[len(pixelArray)-11:])

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

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	extract()
}

func extract() {
	// Open the file and get the pixel array
	pixelArray, _, _ := getImage()

	// Getting stored message length, including extra unused bit
	var leng []byte
	for i := range pixelArray[len(pixelArray)-11:] {
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].R&1))
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].G&1))
		leng = append(leng, byte(pixelArray[len(pixelArray)-i-1].B&1))
	}

	messageLength := b2d(leng[:len(leng)-1])

	// Debugging-- Getting len
	fmt.Println("Message length: ", messageLength/8, "Characters.")

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
			// char = append(char, byte(tmpMessage[(i*8)+j]))
			char = (char << 1) + tmpMessage[(i*8)+j]
		}
		message = append(message, char)
	}
	// Converting
	fmt.Print(string(message))

}

func getImage() ([]Pixel, int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return getPixelArray(file)

}

func getPixelArray(file io.Reader) ([]Pixel, int, int) {

	loadedImage, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("woops")
		log.Fatal(err)
	}

	height := loadedImage.Bounds().Dy()
	width := loadedImage.Bounds().Dx()
	var pixels []Pixel
	// fmt.Println(loadedImage.At(width-1, height-1).RGBA())
	// fmt.Println(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, toPixel(loadedImage.At(x, y).RGBA()))
		}
	}

	return pixels, width, height
}

func toPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	// Takes in a golang pixel in rgba format
	// and converts all the values from uint32s to normal integers
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func b2d(input []byte) float64 {
	result := 0.0
	for i, val := range input {
		power := len(input) - 1 - i
		if val == 1 {
			result += math.Pow(2, float64(power))
		}
		fmt.Println(power, result)
	}
	return result
}

func i2b(input int) [32]byte {
	var result [32]byte
	for i := 0; i < 32; i++ {
		if (input & 1) == 1 {
			result[31-i] = 1
		}
		input = input >> 1
	}
	return result
}
