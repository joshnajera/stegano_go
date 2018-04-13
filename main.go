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
	"math"
)

// var file_name = "test.png"
var file_name = "testImage.png"

func b2d(input []byte) float64 {
	result := 0.0
	for i, val := range(input) {
		power := len(input) - 1 - i
		if val == 1 {
			result += math.Pow(2, float64(power))
		}
		fmt.Println(power, result)
		}
	return result
}

func main() {
	pixel_array := get_image()
	fmt.Println(len(pixel_array))
	// for _, val := range pixel_array[len(pixel_array)-1][len(pixel_array[0])-11:] {
	fmt.Println("Last values")
	var leng []byte
	for i := range pixel_array[len(pixel_array)-11:] {
		fmt.Println(i)
		leng = append(leng, byte(pixel_array[len(pixel_array)-i-1].R&1))
		leng = append(leng, byte(pixel_array[len(pixel_array)-i-1].G&1))
		leng = append(leng, byte(pixel_array[len(pixel_array)-i-1].B&1))
	}
	// output := binary.BigEndian.Uint32(leng[:len(leng)-1])
	// output := uint32(leng[:len(leng)-1])
	fmt.Println(leng[:len(leng)-1])
	fmt.Println(b2d(leng[:len(leng)-1]))
	new_string := fmt.Sprintf("%08b", []byte("ab ab"))         // Pad with leading 0s
	replacer := strings.NewReplacer("[", "", "]", "", " ", "") // Stripping [, ], and whitespace
	new_string = replacer.Replace(new_string)
	fmt.Println("\noutput: ", new_string)

}

func get_image() []Pixel {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pixel_array, err := get_pixel_array(file)
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

func get_pixel_array(file io.Reader) ([]Pixel, error) {
	loaded_image, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("woops")
		log.Fatal(err)
	}

	height := loaded_image.Bounds().Dy()
	width := loaded_image.Bounds().Dx()
	fmt.Println(loaded_image.At(width-1, height-1).RGBA())
	fmt.Println(width, height)
	var pixels []Pixel

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, convert_to_pixel(loaded_image.At(x, y).RGBA()))
			// if x == width-1 && y == height-1 {
			// 	fmt.Println(row[width-1])
			// }
		}
	}

	return pixels, nil
}
