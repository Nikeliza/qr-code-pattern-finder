package main

import (
	"flag"
	"fmt"
	"image"
	"os"
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
func printStat(numberImage, patterns int, timeOpen, timeRead, timeBinarization, timeFilter, timeDetection, timeAll time.Duration) {
	fmt.Println(numberImage, ": count_detect_pattern", patterns, "time_all", timeAll, "time_open", timeOpen, "time_read", timeRead,
		"time_binarization", timeBinarization, "time_filter", timeFilter, "time_detection", timeDetection)
}

//Запуск анализа изображений
func run(path string, binar string, paramBinar int, filter string, paramFilter int, skipPows, countImage int, printBinar, printFilter, printFillRectangle bool) int {
	summ_res := 0
	for i := 1; i <= countImage; i++ {
		start1 := time.Now()
		// Код для измерения
		pathRes := ""
		if i < 10 {
			pathRes = path + "/000" + strconv.Itoa(i) + ".png"
		} else {
			if i < 100 {
				pathRes = path + "/00" + strconv.Itoa(i) + ".png"
			} else {
				if i < 1000 {
					pathRes = path + "/0" + strconv.Itoa(i) + ".png"
				} else {
					pathRes = path + "/" + strconv.Itoa(i) + ".png"
				}
			}
		}

		f, err := os.Open(pathRes)
		if err != nil {
			// Handle error
			fmt.Println("Не возможно открыть изображение " + pathRes)
			return summ_res
		}


		img, _, err := image.Decode(f)
		if err != nil {
			// Handle error
			fmt.Println("Ошибка декодирования " + pathRes)
		}
		timeOpen := time.Since(start1)
		start := time.Now()

		im := readImage(img, false)
		timeRead := time.Since(start)
		start = time.Now()

		im = choiceBinarization(im, binar, paramBinar)
		timeBinarization := time.Since(start)
		start = time.Now()
		if printBinar {
			createGrayImage(pathRes, im, "binar" + binar + strconv.Itoa(paramBinar))
		}

		im = choiceFiltr(im, filter, paramFilter)
		timeFilter := time.Since(start)
		start = time.Now()
		if printFilter {
			createGrayImage(pathRes, im, "filter" + filter + strconv.Itoa(paramFilter))
		}

		var result int
		result, im = detector(im, skipPows, printFillRectangle)
		timeDetection := time.Since(start)
		timeAll := time.Since(start1)

		summ_res += result
		createResultImage(pathRes, img, im)
		printStat(i, result, timeOpen, timeRead, timeBinarization, timeFilter, timeDetection, timeAll)
		f.Close()
	}
	return summ_res
}

func main(){

	pathImagePtr := flag.String("path_image", "./Image/TestSet1/", "Путь к изображениям для обработки")
	countImagePtr := flag.Int("count_image", 47, "Число изображений")
	typeBinarizationPtr := flag.String("type_binarization", "bradley", "Тип бинаризации: bradley, otsu, mean")
	windowMeanPtr := flag.Int("window_mean", 20, "Размер окна при бинаризации среднее в окне")
	typeFilterPtr := flag.String("type_filter", "none", "Тип фильтра: none, median, dilation, erosion, opening, closing")
	windowMedianPtr := flag.Int("window_median", 5, "Размер окна при медианной фильтрации")
	numberMaskPtr := flag.Int("number_mask", 0, "Номер маски при морфологических операциях: 0, 1, 2")
	skipRowsPtr := flag.Int("skip_row", 3, "Число пропущенных строк при детекции паттернов")
	printBinarPtr := flag.Bool("print_binar", true, "Печать после бинаризации")
	printFilterPtr := flag.Bool("print_filter", true, "Печать после фильтрации")
	printFillRectanglePtr := flag.Bool("print_fill_rectangle", false, "Печать заполненного прямоугольника")

	flag.Parse()

	err := false
	if *countImagePtr < 1 {
		fmt.Println("Число изображений не может быть меньше 0")
		err = true
	}
	if *typeBinarizationPtr != "bradley" && *typeBinarizationPtr != "otsu" && *typeBinarizationPtr != "mean" {
		fmt.Println("Неверная бинаризация")
		err = true
	}
	if *windowMeanPtr < 10 {
		fmt.Println("Слишком маленькое окно")
		err = true
	}
	if *typeFilterPtr != "none" && *typeFilterPtr != "median" && *typeFilterPtr != "dilation" &&
		*typeFilterPtr != "erosion" && *typeFilterPtr != "opening" && *typeFilterPtr != "closing" {
		fmt.Println("Неверный фильтр")
		err = true
	}
	if *windowMedianPtr < 3 {
		fmt.Println("CСлищком маленькое окно при медианной фильтрации")
		err = true
	}
	if *numberMaskPtr < 0 || *numberMaskPtr > 2 {
		fmt.Println("Неверный номер маски")
		err = true
	}
	if *skipRowsPtr < 1 {
		fmt.Println("Слишком маленькое число пропущенных строк")
		err = true
	}

	if err {
		return
	}

	fmt.Println("path_image", *pathImagePtr)
	fmt.Println("count_images", *countImagePtr)
	fmt.Println("type_binarization", *typeBinarizationPtr)
	if *typeBinarizationPtr == "mean" {
		fmt.Println("window_mean", *windowMeanPtr)
	}
	fmt.Println("type_filter", *typeFilterPtr)
	if *typeFilterPtr == "median" {
		fmt.Println("window_median", *windowMedianPtr)
	} else {
		if *typeFilterPtr != "none" {
			fmt.Println("number_mask", *numberMaskPtr)
		}
	}
	fmt.Println("skip_rows", *skipRowsPtr)

	if *typeFilterPtr == "median" {
		run(*pathImagePtr, *typeBinarizationPtr, *windowMeanPtr, *typeFilterPtr, *windowMedianPtr,
			*skipRowsPtr, *countImagePtr, *printBinarPtr, *printFilterPtr, *printFillRectanglePtr)
	} else {
		run(*pathImagePtr, *typeBinarizationPtr, *windowMeanPtr, *typeFilterPtr, *numberMaskPtr,
			*skipRowsPtr, *countImagePtr, *printBinarPtr, *printFilterPtr, *printFillRectanglePtr)
	}
}