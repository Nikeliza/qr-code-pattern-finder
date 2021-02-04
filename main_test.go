package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"strconv"
	"testing"
	"time"
)

func sravnImage(ish image.Image, res image.Image) int {
	count := 0

	for i := 0; i < ish.Bounds().Max.X; i++ {
		for j := 0; j < ish.Bounds().Max.Y; j++ {
			r1, g1, b1, _ := ish.At(i, j).RGBA()
			r2, g2, b2, _ := res.At(i, j).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 {
				count += 1
			}
		}
	}

	return count
}

//Сравниваем исходное и результируещее изображение, должны отличаться не сильно
func TestRun(t *testing.T)  {
	run("./Test", "mean", 20, "median", 5, 3, 1, false, false, false)

	pathRes := "./Test/0001.png"
	f, err := os.Open(pathRes)
	if err != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathRes, err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathRes)
		return
	}

	pathRes_res := "./Test/result/0001_result.png"
	f1, err1 := os.Open(pathRes_res)
	if err1 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathRes_res)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathRes_res)
		return
	}

	res := sravnImage(img, img1)
	if res > 1000  {
		t.Error("Неверный вывод изображения", res)
	}
}

//При одинаковой бинаризации и разных параметрах, разные бинаризованные изображения
func TestRun1 (t *testing.T) {
	run("./Test", "mean", 20, "median", 5, 3, 1, true, false, false)
	run("./Test", "mean", 40, "median", 5, 3, 1, true, false, false)

	pathresRes := "./Test/binarmean20/0001_binarmean20.png"
	f1, err1 := os.Open(pathresRes)
	if err1 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes)
		return
	}


	pathresRes1 := "./Test/binarmean40/0001_binarmean40.png"
	f2, err2 := os.Open(pathresRes1)
	if err2 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes1)
		return
	}
	defer f2.Close()

	img2, _, err2 := image.Decode(f2)
	if err2 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes1)
		return
	}

	result := sravnImage(img1, img2)
	if result < 100 {
		t.Error("Неверная бинаризация среднем", result)
	}
}

//При разной бинаризации разные изображения
func TestRun2 (t *testing.T) {
	run("./Test", "mean", 20, "median", 5, 3, 1, true, false, false)
	run("./Test", "otsu", 20, "median", 5, 3, 1, true, false, false)

	pathresRes := "./Test/binarmean20/0001_binarmean20.png"
	f1, err1 := os.Open(pathresRes)
	if err1 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes)
		return
	}


	pathresRes1 := "./Test/binarotsu20/0001_binarotsu20.png"
	f2, err2 := os.Open(pathresRes1)
	if err2 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes1)
		return
	}
	defer f2.Close()

	img2, _, err2 := image.Decode(f2)
	if err2 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes1)
		return
	}

	result := sravnImage(img1, img2)
	if result < 100 {
		t.Error("Неверно работает бинаризация", result)
	}
}

//При медианной фильтрации и разных окнах разный изображения после фильтра
func TestRun3 (t *testing.T) {
	run("./Test", "bradley", 20, "median", 5, 3, 1, false, true, false)
	run("./Test", "bradley", 20, "median", 7, 3, 1, false, true, false)

	pathresRes := "./Test/filtermedian5/0001_filtermedian5.png"
	f1, err1 := os.Open(pathresRes)
	if err1 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes)
		return
	}


	pathresRes1 := "./Test/filtermedian7/0001_filtermedian7.png"
	f2, err2 := os.Open(pathresRes1)
	if err2 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes1)
		return
	}
	defer f2.Close()

	img2, _, err2 := image.Decode(f2)
	if err2 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes1)
		return
	}

	result := sravnImage(img1, img2)
	if result < 100 {
		t.Error("Неверная медианная фильтрация", result)
	}
}

//При морфологических операциях при разных масках разные результирующие изображения
func TestRun4 (t *testing.T) {
	run("./Test", "bradley", 20, "erosion", 0, 3, 1, false, true, false)
	run("./Test", "bradley", 20, "erosion", 1, 3, 1, false, true, false)
	run("./Test", "bradley", 20, "erosion", 2, 3, 1, false, true, false)

	pathresRes := "./Test/filtererosion0/0001_filtererosion0.png"
	f1, err1 := os.Open(pathresRes)
	if err1 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes)
		return
	}
	defer f1.Close()

	img1, _, err1 := image.Decode(f1)
	if err1 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes)
		return
	}


	pathresRes1 := "./Test/filtererosion1/0001_filtererosion1.png"
	f2, err2 := os.Open(pathresRes1)
	if err2 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes1)
		return
	}
	defer f2.Close()

	img2, _, err2 := image.Decode(f2)
	if err2 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes1)
		return
	}


	pathresRes2 := "./Test/filtererosion2/0001_filtererosion2.png"
	f3, err3 := os.Open(pathresRes2)
	if err3 != nil {
		// Handle error
		t.Error("Не возможно открыть изображение " + pathresRes2)
		return
	}
	defer f3.Close()

	img3, _, err3 := image.Decode(f3)
	if err3 != nil {
		// Handle error
		t.Error("Ошибка декодирования " + pathresRes2)
		return
	}

	result1 := sravnImage(img1, img2)
	result2 := sravnImage(img1, img3)
	result3 := sravnImage(img3, img2)
	if result1 < 100 || result2 < 100 || result3 < 100 {
		t.Error("Неверная работа фильтра морфологических операций", result1, result2, result3)
	}
}

//При увелечении пропущенных строк уменьшается время работы
func TestRun5 (t *testing.T) {
	start1 := time.Now()
	run("./Test", "bradley", 20, "median", 5, 1, 1, false, false, false)
	time1 := time.Since(start1)
	start1 = time.Now()
	run("./Test", "bradley", 20, "median", 5, 5, 1, false, false, false)
	time5 := time.Since(start1)

	if (time1 - time5).Microseconds() <= 0  {
		t.Error("Нет увелечения скорости, при увелечении пропущенных строк")
	}
}

//Проверка числа определенных паттернов
func TestRun6 (t *testing.T) {
	res := run("./Test", "bradley", 20, "median", 5, 5, 5, false, false, false)
	if res < 12  {
		t.Error("Плохая работа детектора ", res, "определенных pattern qr кода")
	}
}

//Проверка точности определения паттернов
func TestRun7 (t *testing.T) {
	run("./Test", "bradley", 20, "median", 5, 5, 5, false, false, true)
	summSrav := 0
	for i := 1; i <= 5; i++ {

		path := "./Test/result/000" + strconv.Itoa(i) + "_result.png"
		f1, err1 := os.Open(path)
		if err1 != nil {
			// Handle error
			t.Error("Не возможно открыть изображение " + path)
			return
		}
		defer f1.Close()

		img1, _, err1 := image.Decode(f1)
		if err1 != nil {
			// Handle error
			t.Error("Ошибка декодирования " + path)
			return
		}

		pathRect := "./Test/000" + strconv.Itoa(i) + "_mask.png"
		f2, err2 := os.Open(pathRect)
		if err2 != nil {
			// Handle error
			t.Error("Не возможно открыть изображение " + pathRect)
			return
		}
		defer f2.Close()

		img2, _, err2 := image.Decode(f2)
		if err2 != nil {
			// Handle error
			t.Error("Ошибка декодирования " + pathRect)
			return
		}
		summSrav += sravnImage(img2, img1)
	}

	if summSrav / 18 > 1200  {
		t.Error("Плохая точность детекции", summSrav / 18, "на 1 pattern qr кода")
	}
}
