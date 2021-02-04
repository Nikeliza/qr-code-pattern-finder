package main

import "math"

//Функция детекции паттернов
func detector(img MyImage, skipRows int, printFillRectangle bool) (int, MyImage){
	centers, sizes := find(img, skipRows)
	found := len(centers) > 0
	var result MyImage
	result.height = img.height
	result.width = img.width

	if found == true {
		var imgRes MyImage
		imgRes.height = img.height
		imgRes.width = img.width
		result = drawFinderPatterns(centers, sizes, imgRes, printFillRectangle)
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
func drawFinderPatterns(centers []point, sizes []int, img MyImage, printFillRectangle bool) MyImage {
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
			img = drawRectangle(img, point1, point2, printFillRectangle)
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