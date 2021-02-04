package main





/*
import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

//Применение выбранной морфологической операции
func choiceMorphologyOperation(img MyImage, choice string, param int) MyImage {
	mask1 := [][] int {{0, 1, 0}, {1, 1, 1}, {0, 1, 0}}
	mask2 := [][] int {{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	mask3 := [][] int {{0, 0, 1, 1, 1, 0, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 0, 1, 1, 1, 0, 0}}

	var mask [][]int
	var size int
	if param == 0 {
		mask = mask1
		size = 3
	} else {
		if param == 1 {
			mask = mask2
			size = 3
		} else {
			if param == 2 {
				mask = mask3
				size = 7
			}
		}
	}

	if choice == "dilation" {
		return dilation(img, mask, size)
	}
	if choice == "erosion" {
		return erosion(img, mask, size)
	}
	if choice == "opening" {
		return opening(img, mask, size)
	}
	if choice == "closing" {
		return closing(img, mask, size)
	}

	return img
}

//Применение выбранной бинаризации
func choiceBinarization(img MyImage, choice string, param int) MyImage {
	if choice == "bradley" {
		return BradleyThreshold(img)
	}
	if choice == "otsu" {
		return otsuBinarization(img)
	}
	if choice == "mean" {
		return sredBinarization(img, param)
	}

	return img
}

//Выбор фильтра
func choiceFiltr(img MyImage, choice string, param int) MyImage {
	if choice == "median" {
		return medianFilth(img, param)
	} else {
		if choice == "opening" || choice == "closing" || choice == "erosion" || choice == "dilation" {
			return choiceMorphologyOperation(img, choice, param)
		}
	}

	return img
}

//Распечатка статистики
func printStat(patterns int, timeOpen, timeRead, timeBinarization, timeFilter, timeDetection, timeAll time.Duration) {
	fmt.Println("count_detect_pattern", patterns, "time_all", timeAll, "time_open", timeOpen, "time_read", timeRead,
		"time_binarization", timeBinarization, "time_filter", timeFilter, "time_detection", timeDetection)
}

//Запуск анализа изображений
func run(path string, binar string, paramBinar int, filter string, paramFilter int, skipPows, countImage int) {

	for i := 1; i < countImage; i++ {
		fmt.Println(i)
		start1 := time.Now()
		// Код для измерения
		pathRes := ""
		if i < 10 {
			pathRes = path + "/000" + strconv.Itoa(i) + ".jpg"
		} else {
			if i < 100 {
				pathRes = path + "/00" + strconv.Itoa(i) + ".jpg"
			} else {
				if i < 1000 {
					pathRes = path + "/0" + strconv.Itoa(i) + ".jpg"
				} else {
					pathRes = path + "/" + strconv.Itoa(i) + ".jpg"
				}
			}
		}
		fmt.Println("begin1")
		f, err := os.Open(pathRes)
		if err != nil {
			// Handle error
			fmt.Println("Не возможно открыть изображение " + pathRes)
			return
		}
		defer f.Close()
		fmt.Println("begin2")
		img, _, err := image.Decode(f)
		if err != nil {
			// Handle error
			fmt.Println("Ошибка декодирования " + pathRes)
		}
		timeOpen := time.Since(start1)
		start := time.Now()
		fmt.Println("begin3")
		im := readImage(img)
		timeRead := time.Since(start)
		start = time.Now()
		fmt.Println("begin4")
		im = choiceBinarization(im, binar, paramBinar)
		timeBinarization := time.Since(start)
		start = time.Now()
		createGrayImage(pathRes, im, "binar")
		fmt.Println("begin5")
		im = choiceFiltr(im, filter, paramFilter)
		timeFilter := time.Since(start)
		start = time.Now()
		fmt.Println("begin6")
		createGrayImage(pathRes, im, "filter")

		var result int
		result, im = detector(im, skipPows)
		timeDetection := time.Since(start)
		timeAll := time.Since(start1)
		fmt.Println("begin7")
		createResultImage(pathRes, img, im)
		printStat(result, timeOpen, timeRead, timeBinarization, timeFilter, timeDetection, timeAll)
	}
}

func main(){

	pathImagePtr := flag.String("path_image", "./Image/TestSet1/", "Путь к изображениям для обработки")
	countImagePtr := flag.Int("count_image", 47, "Число изображений")
	typeBinarizationPtr := flag.String("type_binarization", "bradley", "Тип бинаризации: bradley, otsu, mean")
	windowMeanPtr := flag.Int("window_mean", 20, "Размер окна при бинаризации среднее в окне")
	typeFilterPtr := flag.String("type_filter", "median", "Тип фильтра: none, median, dilation, erosion, opening, closing")
	windowMedianPtr := flag.Int("window_median", 5, "Размер окна при медианной фильтрации")
	numberMaskPtr := flag.Int("number_mask", 0, "Номер маски при морфологических операциях: 0, 1, 2")
	skipRowsPtr := flag.Int("skip_row", 3, "Число пропущенных строк при детекции паттернов")

	flag.Parse()

	if *typeFilterPtr == "median" {
		fmt.Println("begin1")
		run(*pathImagePtr, *typeBinarizationPtr, *windowMeanPtr, *typeFilterPtr, *windowMedianPtr, *skipRowsPtr, *countImagePtr)
	} else {
		fmt.Println("begin2")
		run(*pathImagePtr, *typeBinarizationPtr, *windowMeanPtr, *typeFilterPtr, *numberMaskPtr, *skipRowsPtr, *countImagePtr)
	}
}


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


//Функция детекции паттернов
func detector(img MyImage, skipRows int) (int, MyImage){
	centers, sizes := find(img, skipRows)
	found := len(centers) > 0
	var result MyImage
	result.height = img.height
	result.width = img.width

	if found == true {
		var imgRes MyImage
		imgRes.height = img.height
		imgRes.width = img.width
		result = drawFinderPatterns(centers, sizes, imgRes)
	}

	return len(centers), result
}

//Функция нахождения паттернов
func find(img MyImage, skipRows int) ([]point, []int){
	var centers []point
	var moduleSize []int
	currentState := 0

	for row := 0 + skipRows - 1; row < img.height; row += skipRows {
		var counterState [5]int
		currentState = 0
		ptr := img.imagePix[row][:]

		for col := 0; col < img.width; col++ {
			if ptr[col] < 128 { // черный пиксель
				if (currentState & 1) == 1 {
					currentState += 1
				}
				counterState[currentState] += 1
			} else { //белый пиксель
				if (currentState & 1) == 1 { //если 1 или 3 положение
					counterState[currentState] += 1
				} else {

					if currentState == 4 {
						if checkRatio(counterState, img.height, img.width) {
							//this is where we do some more checks
							_, centers, moduleSize = handlePossibleCenter(img, counterState, row, col, centers, moduleSize)
						} else {
							currentState = 3
							counterState[0] = counterState[2]
							counterState[1] = counterState[3]
							counterState[2] = counterState[4]
							counterState[3] = 1
							counterState[4] = 0
							continue
						}

						currentState = 0
						counterState[0] = 0
						counterState[1] = 0
						counterState[2] = 0
						counterState[3] = 0
						counterState[4] = 0
					} else {
						currentState += 1
						counterState[currentState] += 1
					}
				}
			}
		}
	}
	return centers, moduleSize
}

//Функция отрисовки паттернов
func drawFinderPatterns(centers []point, sizes []int, img MyImage) MyImage {
	if len(centers) == 0 {
		return img
	}

	for i := range centers {
		pt := centers[i]
		diff := float64(sizes[i]) * 3.5
		point1 := point{pt.x - int(diff), pt.y - int(diff)}
		if point1.x < 0 {
			point1.x = 0
		}

		if point1.y < 0{
			point1.y = 0
		}

		point2 := point{pt.x + int(diff), pt.y + int(diff)}
		if math.Abs(float64(point1.x - point2.x)) > float64(img.height / 100) &&
			math.Abs(float64(point1.y - point2.y)) > float64(img.width / 100) {
			img = drawRectangle(img, point1, point2, false)
		}
	}

	return img
}

//Проверка пропорциональности найденной линии
func checkRatio(stateCount [5]int, height, width int) bool {
	totalWidth := 0
	for i := 0; i < 5; i++ {
		if stateCount[i] == 0 {
			return false
		}
		totalWidth += stateCount[i]
	}

	if totalWidth < 7 && (totalWidth > min(height / 100, width / 100)) {
		return false
	}

	widthPattern := math.Round(float64(totalWidth) / 7.0)
	dispersion := widthPattern / 2

	result := (math.Abs(widthPattern- float64(stateCount[0])) < dispersion) &&
		(math.Abs(widthPattern- float64(stateCount[1])) < dispersion) &&
		(math.Abs(3 *widthPattern- float64(stateCount[2])) < 3 * dispersion) &&
		(math.Abs(widthPattern- float64(stateCount[3])) < dispersion) &&
		(math.Abs(widthPattern- float64(stateCount[4])) < dispersion)
	return result
}

//Функция добавления нового центра
func addCenter(centerCol, centerRow, totalState int, centers []point, sizes []int) ([]point, []int){
	newCenter := point{centerCol, centerRow}

	newModuleSize := totalState / 7.0
	found := false

	for i := range centers {
		diff := point{0, 0}
		diff.x = centers[i].x - newCenter.x
		diff.y = centers[i].y - newCenter.y
		distance := math.Sqrt(float64(diff.x * diff.x + diff.y * diff.y))

		if distance < 10 {
			centers[i].x = centers[i].x + newCenter.x
			centers[i].y = centers[i].y + newCenter.y
			centers[i].x /= 2.0
			centers[i].y /= 2.0
			sizes[i] = (sizes[i] + newModuleSize) / 2.0
			found = false
			break
		}
	}

	if !found && !isValueInArray(centers, newCenter) {
		centers = append(centers, newCenter)
		sizes = append(sizes, newModuleSize)
	}

	return centers, sizes
}

//Проверка и уточнение центра паттерна
func handlePossibleCenter(img MyImage, stateCount [5]int, row, col int, centers []point, sizes []int) (bool, []point, []int) {
	totalState := 0
	for i := 0; i < 5; i++ {
		totalState += stateCount[i]
	}

	if totalState > img.height / 100 {
		centerCol := getCenter(stateCount, col)
		centerRow := checkVertical(img, row, centerCol, stateCount[2], totalState)
		if centerRow == -1.0 {
			return false, centers, sizes
		}

		newCenter := point{centerCol, centerRow}
		if isValueInArray(centers, newCenter) != true {
			centerCol = checkHorizontal(img, centerRow, centerCol, stateCount[2], totalState)
			if centerCol == -1.0 {
				return false, centers, sizes
			}
			newCenter := point{centerCol, centerRow}
			if isValueInArray(centers, newCenter) != true {
				if !checkDiagonal(img, centerRow, centerCol, stateCount[2], totalState) {
					return false, centers, sizes
				}
				centers, sizes = addCenter(centerCol, centerRow, totalState, centers, sizes)
			}
		}
	}

	return false, centers, sizes
}

//Проверка диагонали
func checkDiagonal(img MyImage, centerRow, centerCol, maxCount, stateCountTotal int) bool {
	var stateCount [5]int
	i := 0
	centerRow = int(centerRow) + 1
	centerCol = int(centerCol) + 1

	for centerRow >= i && centerCol >= 1 && img.imagePix[centerRow- i][centerCol- i] < 128 {
		stateCount[2] += 1
		i += 1
		if centerRow < i || centerCol < i {
			return false
		}
	}

	for centerRow >= i && centerCol >= i && img.imagePix[centerRow- i][centerCol- i] >= 128 && stateCount[1] <= maxCount {
		stateCount[1] += 1
		i += 1
		if centerRow < i || centerCol < i || stateCount[1] > maxCount {
			return false
		}
	}

	for centerRow >= i && centerCol >= i && img.imagePix[centerRow- i][centerCol- i] < 128 && stateCount[0] <= maxCount {
		stateCount[0] += 1
		i += 1
		if stateCount[0] > maxCount {
			return false
		}
	}

	i = 1
	for centerRow+ i < img.height && centerCol+ i < img.width && img.imagePix[centerRow+ i][centerCol+ i] < 128 {
		stateCount[2] += 1
		i += 1
		if centerRow+ i >= img.height || centerCol+ i >= img.width {
			return false
		}
	}

	for centerRow+ i < img.height && centerCol+ i < img.width && img.imagePix[centerRow+ i][centerCol+ i] >= 128 && stateCount[3] < maxCount {
		stateCount[3] += 1
		i += 1
		if centerRow+ i >= img.height || centerCol+ i >= img.width || stateCount[3] > maxCount {
			return false
		}
	}

	for centerRow+ i < img.height && centerCol+ i < img.width && img.imagePix[centerRow+ i][centerCol+ i] < 128 && stateCount[4] < maxCount {
		stateCount[4] += 1
		i += 1
		if stateCount[4] > maxCount {
			return false
		}
	}

	newStateCountTotal := 0
	for j := 0; j < 5; j++ {
		newStateCountTotal += stateCount[j]
	}
	res1 := math.Abs(float64(stateCountTotal-newStateCountTotal)) < float64(2 *stateCountTotal)
	res2 := checkRatio(stateCount, img.height, img.width)

	return res1 && res2
}

//Проверка вертикали
func checkVertical(img MyImage, startRow, centerCol, centralCount, stateCountTotal int) int {
	var counterState [5]int
	row := startRow
	centerCol = int(centerCol)

	for row >= 0 && img.imagePix[row][centerCol] < 128 {
		counterState[2] += 1
		row -= 1
	}
	if row < 0 {
		return -1
	}

	for row >= 0 && img.imagePix[row][centerCol] >= 128 && counterState[1] < centralCount {
		counterState[1] += 1
		row -= 1
	}
	if row < 0 || counterState[1] > centralCount + centralCount / 5 {
		return  -1
	}

	for row >= 0 && img.imagePix[row][centerCol] < 128 && counterState[0] < centralCount {
		counterState[0] += 1
		row -= 1
	}
	if counterState[0] > centralCount + centralCount / 5 {
		return -1
	}

	row = startRow + 1
	for row < img.height && img.imagePix[row][centerCol] < 128 {
		counterState[2] += 1
		row += 1
	}
	if row == img.height {
		return -1
	}

	for row < img.height && img.imagePix[row][centerCol] >= 128 && counterState[3] < centralCount {
		counterState[3] += 1
		row += 1
	}
	if row == img.height || counterState[3] > centralCount + centralCount / 5 {
		return -1
	}

	for row < img.height && img.imagePix[row][centerCol] < 128 && counterState[4] < centralCount {
		counterState[4] += 1
		row += 1
	}
	if counterState[4] > centralCount + centralCount / 5 {
		return  -1
	}

	counterStateTotal := 0
	for i := 0; i < 5; i++ {
		counterStateTotal += counterState[i]
	}
	if 5 * math.Abs(float64(counterStateTotal-stateCountTotal)) >= float64(2 *stateCountTotal) {
		return -1
	}

	center := getCenter(counterState, row)
	if checkRatio(counterState, img.height, img.width) {
		return center
	} else {
		return -1
	}
}

//Проверка горизонтали
func checkHorizontal(img MyImage, centerRow, startCol, centerCount, stateCountTotal int) int {
	var counterState [5]int
	col := int(startCol)
	centerRow = int(centerRow)
	ptr := img.imagePix[centerRow][:]

	for col >= 0 && ptr[col] < 128 {
		counterState[2] += 1
		col -= 1
	}
	if col < 0 {
		return -1
	}

	for col >= 0 && ptr[col] >= 128 && counterState[1] < centerCount {
		counterState[1] += 1
		col -= 1
	}
	if col < 0 || counterState[1] == centerCount {
		return -1
	}

	for col >= 0 && ptr[col] < 128 && counterState[0] < centerCount {
		counterState[0] += 1
		col -= 1
	}
	if counterState[0] == centerCount {
		return -1
	}

	col = int(startCol) + 1
	for col < img.width && ptr[col] < 128 {
		counterState[2] += 1
		col += 1
	}
	if col == img.width {
		return -1
	}

	for col < img.width && ptr[col] >= 128 && counterState[3] < centerCount {
		counterState[3] += 1
		col += 1
	}
	if col == img.width || counterState[3] == centerCount {
		return -1
	}

	for col < img.width && ptr[col] < 128 && counterState[4] < centerCount {
		counterState[4] += 1
		col += 1
	}
	if counterState[4] >= centerCount {
		return -1
	}

	counterStateTotal := 0
	for i := 0; i < 5; i++ {
		counterStateTotal += counterState[i]
	}

	if 5 * math.Abs(float64(counterStateTotal-stateCountTotal)) >= float64(stateCountTotal) {
		return -1
	}

	center := getCenter(counterState, col)
	if checkRatio(counterState, img.height, img.width) {
		return center
	} else {
		return -1
	}
}

//Вычеслить центр паттерна
func getCenter(stateCount [5]int, end int) int {
	return end - stateCount[4] - stateCount[3] - stateCount[2] / 2.0
}


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
func readImage(mas image.Image) MyImage {
	var Mas [4000][4000]uint8

	for i := 0; i < mas.Bounds().Max.X; i++ {
		for j := 0; j < mas.Bounds().Max.Y; j++{
			r, g, b, _  := mas.At(i, j).RGBA()
			Mas[i][j] = uint8((0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) * (255.0 / 65535))
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

	pathResultImage := "result//"
	pathResultImage += path[len(path) - 8:len(path) - 4] + "_result.png"
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

	pathResultImage := typeSave + "/"
	pathResultImage += path[len(path) - 8:len(path) - 4] + "_" + typeSave + ".png"
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





func sravn(img MyImage, otv image.Image) bool {

	rasst := 0.0

	for i := 0; i < img.height; i++ {
		for j := 0; j < img.width; j++ {
			r, g, b, _ := otv.At(i, j).RGBA()
			if math.Abs(float64(uint32(img.imagePix[i][j]) - (r + g + b) * 255.0 / (65535 * 3))) > 10 {
				rasst = math.Abs(float64(uint32(img.imagePix[i][j]) - (r + g + b) * 255.0 / (65535 * 3)))
			}
		}
	}
	if rasst > 100 {
		return true
	} else {
		return false
	}
}

func absdf() {
	pathRes := "./Test/0001.jpg"
	f, err := os.Open(pathRes)
	if err != nil {
		// Handle error
		fmt.Println("Не возможно открыть изображение " + pathRes)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		// Handle error
		fmt.Println("Ошибка декодирования " + pathRes)
	}
	result := readImage(img)

	pathRes_otv := "./Test/0001_otv.jpg"
	f1, err1 := os.Open(pathRes_otv)
	if err1 != nil {
		// Handle error
		fmt.Println("Не возможно открыть изображение " + pathRes_otv)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		fmt.Println("Ошибка декодирования " + pathRes_otv)
	}

	if sravn(result, img1) {
		//t.Error("Неверный перевод в серый")
		fmt.Println("ajshahdhshd")
	}
}

func main()  {
	absdf()
}
*/