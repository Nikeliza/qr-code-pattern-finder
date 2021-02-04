package main

import (
	_ "image/jpeg"
)

//Вычисление интегрального представления изображения
func integralImage(img MyImage) MyImageFloat {
	var mas [4000][4000]float64

	for i := 0; i < img.height; i++ {
		for j := 0; j < img.width; j++ {
			if i > 0 && j > 0 {
				mas [i][j] = float64(img.imagePix[i][j]) + mas[i][j-1] + mas[i - 1][j] - mas[i - 1][j - 1]
			} else {
				if i == 0 && j > 0 {
					mas [i][j] = float64(img.imagePix[i][j]) + mas[i][j-1]
				} else {
					if j == 0 && i > 0 {
						mas [i][j] = float64(img.imagePix[i][j]) + mas[i - 1][j]
					} else {
						mas [i][j] = float64(img.imagePix[i][j])
					}
				}
			}

		}
	}

	return MyImageFloat{mas, img.height, img.width}
}

//Бинаризация Брейдли
func BradleyThreshold(src MyImage) MyImage {
	S := src.width / 8
	s2 := S / 2
	t := 0.15
	integralImg := integralImage(src)
	sum := 0
	count := 0

	var x1, y1, x2, y2 int
	var res MyImage
	res.height = src.height
	res.width = src.width
	//находим границы для локальные областей
	for i := 0; i < src.width; i++ {
		for j := 0; j < src.height; j++ {

			x1 = i - s2
			x2 = i + s2
			y1 = j - s2
			y2 = j + s2

			if x1 < 0 {
				x1 = 0
			}
			if x2 >= src.width {
				x2 = src.width - 1
			}
			if y1 < 0 {
				y1 = 0
			}
			if y2 >= src.height {
				y2 = src.height - 1
			}

			count = (x2 - x1) * (y2 - y1)

			sum = int(integralImg.imagePix[y2][x2] - integralImg.imagePix[y1][x2] - integralImg.imagePix[y2][x1] + integralImg.imagePix[y1][x1])
			if float64(int(src.imagePix[j][i]) * count) < float64(sum) * (1.0 - t){
				res.imagePix[j][i] = 0
			} else {
				res.imagePix[j][i] = 255
			}
		}
	}

	return res
}

// посчитать порог по Оцу
func calculateOtsuThreshold(img MyImage) int {
	// вычисляем гистограмму
	histogram := createHistogram(img)

	// аккумулятор суммы яркостей
	sumOfLuminance := 0

	// вычисляем сумму яркостей
	for x := 0; x < img.height; x++ {
		for y := 0; y < img.width; y++ {
			sumOfLuminance += int(img.imagePix[x][y])
		}
	}

	// общее количество пикселей
	allPixelCount := float64(img.width * img.height)

	// оптимальный порог
	bestThreshold := 0
	// количество полезных пикселей
	firstClassPixelCount := 0
	// суммарная яркость полезных пикселей
	firstClassLuminanceSum := 0

	// оптимальный разброс яркостей
	bestSigma := 0.0

	for threshold := 0; threshold < 255; threshold++ {
		firstClassPixelCount += histogram[threshold]
		firstClassLuminanceSum += threshold * histogram[threshold]
		// доля полезных пикселей double
		firstClassProbability := float64(firstClassPixelCount) / allPixelCount
		// доля фоновых пикселей
		secondClassProbability := 1.0 - firstClassProbability
		// средняя доля полезных пикселей
		firstClassMean := 0.0
		if firstClassPixelCount == 0 {
			firstClassMean = 0.0
		} else {
			firstClassMean = float64(firstClassLuminanceSum) / float64(firstClassPixelCount)
		}
		// средняя доля фоновых пикселей
		secondClassMean := float64(sumOfLuminance- firstClassLuminanceSum) / (allPixelCount - float64(firstClassPixelCount))
		// величина разброса
		meanDelta := firstClassMean - secondClassMean
		// общий разброс
		sigma := firstClassProbability * secondClassProbability * meanDelta * meanDelta
		// находим оптимальный разброс
		if sigma > bestSigma {
			bestSigma = sigma
			bestThreshold = threshold
		}
	}

	return bestThreshold
}

// бинаризация по Оцу
func otsuBinarization(img MyImage) MyImage {
	newImage := MyImage{}
	newImage.height = img.height
	newImage.width = img.width
	threshold := calculateOtsuThreshold(img)

	for x := 0; x < img.height; x++ {
		for y := 0; y < img.width; y++ {
			luminance := int(img.imagePix[x][y])

			if luminance > threshold {
				newImage.imagePix[x][y] = 255
			} else {
				newImage.imagePix[x][y] = 0
			}
		}
	}

	return newImage
}

//Бинаризация по среднему значению в области
func sredBinarization(img MyImage, window int) MyImage {
	result := MyImage{}
	result.height = img.height
	result.width = img.width

	integr := integralImage(img)

	for i := window / 2; i < img.height - window / 2; i++ {
		for j := window; j < img.width - window / 2; j++ {
			summ := integr.imagePix[i - window / 2][j - window / 2] + integr.imagePix[i + window / 2][j + window / 2] -
				integr.imagePix[i - window / 2][j + window / 2] - integr.imagePix[i + window / 2][j - window / 2]
			porog := summ / float64(window * window)

			if float64(img.imagePix[i][j]) < porog {
				result.imagePix[i][j] = 0
			} else {
				result.imagePix[i][j] = 255
			}
		}
	}

	return result
}