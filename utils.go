package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

//Структура для хранения серого изображения
type MyImage struct {
	imagePix [4000][4000]uint8
	height   int
	width    int
}

//Структура для хранения изображения в интегральном формате
type MyImageFloat struct {
	imagePix [4000][4000]float64
	height   int
	width    int
}

//Структура для хранения точки
type point struct {
	x, y int
}

//Функция для чтения изображения в структуру
func readImage(mas image.Image, fl bool) MyImage {
	var Mas [4000][4000]uint8

	for i := 0; i < mas.Bounds().Max.X; i++ {
		for j := 0; j < mas.Bounds().Max.Y; j++{
			r, g, b, _  := mas.At(i, j).RGBA()
			if fl {
				Mas[i][j] = uint8((float64(r) + float64(g) + float64(b) / 3) * (255.0 / 65535))
			} else {
				Mas[i][j] = uint8((0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) * (255.0 / 65535))

			}

		}
	}

	return MyImage{Mas, int(mas.Bounds().Max.X), int(mas.Bounds().Max.Y)}
}

//Функция для создания файла
func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

//Функция для печати цветного изображения
func createResultImage(path string, im image.Image, rect MyImage) {
	width := int(rect.height)
	height := int(rect.width)

	upLeft := image.Point{}
	lowRight := image.Point{X: int(width), Y: int(height)}

	img1 := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rect.imagePix[x][y] == 255 {
				img1.Set(x, y, color.RGBA{R: 255, A: 255})
			} else {
				img1.Set(x, y, im.At(x, y))
			}
		}
	}

	pathResultImage := path[:len(path) - 8] + "result//" + path[len(path) - 8:len(path) - 4] + "_result.png"
	// Encode as PNG.
	f1, _ := create(pathResultImage)
	png.Encode(f1, img1)
}

//Функция для печати черно-белого изображения
func createGrayImage(path string, im MyImage, typeSave string) {
	width := int(im.height)
	height := int(im.width)

	upLeft := image.Point{}
	lowRight := image.Point{X: int(width), Y: int(height)}

	img1 := image.NewGray(image.Rectangle{Min: upLeft, Max: lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img1.Set(x, y, color.Gray{Y: im.imagePix[x][y]})
		}
	}

	pathResultImage := path[:len(path) - 8] + typeSave + "/" + path[len(path) - 8:len(path) - 4] + "_" + typeSave + ".png"
	// Encode as PNG.
	f1, _ := create(pathResultImage)
	png.Encode(f1, img1)
}

//Функция для рисования прямоугльника
func drawRectangle(img MyImage, x, y point, fill bool) MyImage {
	if fill {
		return drawRectangleFill(img, x, y)
	} else {
		return drawRectangleContour(img, x, y)
	}
}

//Функция для рисования заполненного прямоугольника
func drawRectangleFill(img MyImage, x, y point) MyImage {
	for i := x.y; i < y.y; i++ {
		for j := x.x; j < y.x; j++ {
			img.imagePix[i][j] = 255
		}
	}

	return img
}

//Функция для рисования контура изображения
func drawRectangleContour(img MyImage, x, y point) MyImage {
	for i := x.y - 1; i < x.y + 1; i++ {
		for j := x.x - 1; j < y.x + 1; j++ {
			img.imagePix[i][j] = 255
		}
	}
	for i := y.y - 1; i < y.y + 1; i++ {
		for j := x.x - 1; j < y.x + 1; j++ {
			img.imagePix[i][j] = 255
		}
	}
	for i := x.y - 1; i < y.y + 1; i++ {
		for j := x.x - 1; j < x.x + 1; j++ {
			img.imagePix[i][j] = 255
		}
	}
	for i := x.y - 1; i < y.y + 1; i++ {
		for j := y.x - 1; j < y.x + 1; j++ {
			img.imagePix[i][j] = 255
		}
	}
	return img
}

//Функция нахождения минимального
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Входит ли значение в массив
func isValueInArray(array []point, value point) bool {
	for i := range array {
		if array[i].x == value.x && array[i].y == value.y {
			return true
		}
	}
	return false
}

// создать гистограмму
func createHistogram(img MyImage) [256]int {
	var histogram [256]int
	for x := 0; x < img.height; x++ {
		for y := 0; y < img.width; y++ {
			instensity := int(img.imagePix[x][y])
			histogram[instensity] += 1
		}
	}

	return histogram
}