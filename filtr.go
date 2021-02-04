package main

import "sort"

//Функция для применения медианного фильтра для устанения шума
func medianFilth(img MyImage, window int) MyImage {
	var res MyImage
	res.width = img.width
	res.height = img.height

	for i := window / 2; i < img.height - window / 2; i++ {
		for j := window / 2; j < img.width - window / 2; j++ {
			var mas []int
			for k := 0; k < window; k++ {
				for q := 0; q < window; q++ {
					mas = append(mas, int(img.imagePix[i + k - window / 2][j + q - window / 2]))
				}
			}
			sort.Ints(mas)
			res.imagePix[i][j] = uint8(mas[len(mas) / 2])
		}
	}

	return res
}

//Функция морфологической операции делатации
func dilation(img MyImage, mask [][]int, size int) MyImage {
	var result MyImage
	result.width = img.width
	result.height = img.height

	for i := size / 2; i < img.height - size / 2; i++ {
		for j := size / 2; j < img.width - size / 2; j++ {
			maxZnach := uint8(0)
			for x := -size / 2; x <= size / 2; x++ {
				for y := -size / 2; y <= size / 2; y++ {
					if mask[size / 2 + x][size / 2 + y] == 1 && img.imagePix[i + x][j + y] > maxZnach {
						maxZnach = img.imagePix[i + x][j + y]
					}
				}
			}
			result.imagePix[i][j] = maxZnach
		}
	}

	return result
}

//Функция морфологической операции эрозии
func erosion(img MyImage, mask [][]int, size int) MyImage {
	var result MyImage
	result.width = img.width
	result.height = img.height

	for i := size / 2; i < img.height - size / 2; i++ {
		for j := size / 2; j < img.width - size / 2; j++ {
			minZnach := uint8(255)
			for x := -size / 2; x <= size / 2; x++ {
				for y := -size / 2; y <= size / 2; y++ {
					if mask[size / 2 + x][size / 2 + y] == 1 && img.imagePix[i + x][j + y] < minZnach {
						minZnach = img.imagePix[i + x][j + y]
					}
				}
			}
			result.imagePix[i][j] = minZnach
		}
	}

	return result
}

//Функция морфологической операции открытия
func opening(img MyImage, mask [][]int, size int) MyImage {
	result := erosion(img, mask, size)
	result = dilation(result, mask, size)

	return result
}

//Функция морфологической операции закрытия
func closing(img MyImage, mask [][]int, size int) MyImage {
	result := dilation(img, mask, size)
	result = erosion(result, mask, size)

	return result
}