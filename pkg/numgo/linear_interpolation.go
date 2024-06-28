package numgo

import "math"

func Linspace(start, stop float64, num int, endpoint bool) []float64 {
	res := make([]float64, num)
	if num == 0 {
		return res
	}
	var den float64
	if endpoint {
		den = float64(num - 1)
	} else {
		den = float64(num)
	}
	for i := 0; i < num; i++ {
		if i == 0 {
			res[i] = start
		} else {
			res[i] = start + float64(i)*(stop-start)/den
		}
	}
	return res
}

func Circspace(start, stop, radius float64, num int, endpoint bool) [][]float64 {
	res := make([][]float64, num)
	if num == 0 {
		return res
	}
	var den float64
	if endpoint {
		den = float64(num - 1)
	} else {
		den = float64(num)
	}
	for i := 0; i < num; i++ {
		t := start + float64(i)*(stop-start)/den
		res[i] = []float64{radius * math.Cos(t), radius * math.Sin(t), 1}

	}
	return res
}
