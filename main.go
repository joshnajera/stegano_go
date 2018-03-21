package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"
)

// var file_name = "test.png"
var file_name = "test.png"

func main() {
	pixel_array := get_image()
	fmt.Println(len(pixel_array), len(pixel_array[0]))
	fmt.Println(pixel_array[399][734])
	// for _, val := range pixel_array[len(pixel_array)-1][len(pixel_array[0])-11:] {
	for _, val := range pixel_array[399][735-11:] {
		fmt.Println(" ", val.R&1, val.G&1, val.B&1)
	}
	new_string := fmt.Sprintf("%08b", []byte("ab ab"))         // Pad with leading 0s
	replacer := strings.NewReplacer("[", "", "]", "", " ", "") // Stripping [, ], and whitespace
	new_string = replacer.Replace(new_string)
	fmt.Println("\noutput: ", new_string)

}

func get_image() [][]Pixel {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pixel_array, err := get_pixel_array(file)
	fmt.Println(pixel_array[399][len(pixel_array[0])-11:])
	return pixel_array

}

type Pixel struct {
	R, G, B, A int
}

func convert_to_pixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	// Takes in a golang pixel in rgba format
	// and converts all the values from uint32s to normal integers
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func get_pixel_array(file io.Reader) ([][]Pixel, error) {
	loaded_image, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("woops")
		log.Fatal(err)
	}

	height := loaded_image.Bounds().Dy()
	width := loaded_image.Bounds().Dx()
	fmt.Println(loaded_image.At(width-1, height-1).RGBA())
	fmt.Println(width, height)
	var pixels [][]Pixel

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, convert_to_pixel(loaded_image.At(x, y).RGBA()))
			if x == width-1 && y == height-1 {
				fmt.Println(row[width-1])
			}
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}
