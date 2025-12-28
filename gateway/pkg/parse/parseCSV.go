package parse

import (
	"encoding/csv"
	"os"
	"strconv"

	"gateway/internal/code"
)

// ParseCSV 分别解析已切好的 passdata 和 futuredata CSV。
// passPath: 73 行 14 列，取前 72 行的数据。
// futurePath: 25 行 13 列，取前 24 行的数据。
// 首行若无法解析视为表头跳过，后续行解析失败则报错。
func ParseCSV(passPath, futurePath string) ([13][72]float64, [12][24]float64, error) {
	passData, err := parsePassData(passPath)
	if err != nil {
		return [13][72]float64{}, [12][24]float64{}, err
	}

	futureData, err := parseFutureData(futurePath)
	if err != nil {
		return [13][72]float64{}, [12][24]float64{}, err
	}

	return passData, futureData, nil
}

func parsePassData(path string) ([13][72]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return [13][72]float64{}, code.ErrParseFile
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [13][72]float64{}, code.ErrParseFile
	}

	var data [13][72]float64
	rowIdx := 0
	for i, r := range records {
		if len(r) < 14 {
			if i == 0 {
				continue
			}
			return [13][72]float64{}, code.ErrParseFile
		}

		if rowIdx >= 72 {
			break
		}

		for c := 0; c < 13; c++ {
			v, err := strconv.ParseFloat(r[c+1], 64) // 跳过首列
			if err != nil {
				if i == 0 {
					return [13][72]float64{}, code.ErrParseFile
				}
				return [13][72]float64{}, code.ErrParseFile
			}
			data[c][rowIdx] = v
		}

		rowIdx++
	}

	if rowIdx < 72 {
		return [13][72]float64{}, code.ErrParseFile
	}

	return data, nil
}

func parseFutureData(path string) ([12][24]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return [12][24]float64{}, code.ErrParseFile
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [12][24]float64{}, code.ErrParseFile
	}

	var data [12][24]float64
	rowIdx := 0
	for i, r := range records {
		if len(r) < 13 {
			if i == 0 {
				continue
			}
			return [12][24]float64{}, code.ErrParseFile
		}

		if rowIdx >= 24 {
			break
		}

		for c := 0; c < 12; c++ {
			v, err := strconv.ParseFloat(r[c+1], 64) // 跳过首列
			if err != nil {
				if i == 0 {
					return [12][24]float64{}, code.ErrParseFile
				}
				return [12][24]float64{}, code.ErrParseFile
			}
			data[c][rowIdx] = v
		}

		rowIdx++
	}

	if rowIdx < 24 {
		return [12][24]float64{}, code.ErrParseFile
	}

	return data, nil
}
