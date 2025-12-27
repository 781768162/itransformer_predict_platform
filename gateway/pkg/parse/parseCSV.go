package parse

import (
	"os"
	"time"
	"sort"
	"strconv"
	"encoding/csv"

	"gateway/internal/code"
)

// TODO: need patch

type parsedRow struct {
	ts      time.Time
	feature [13]float64
}

func ParseCSV(dst, date string) ([13][72]float64, [12][24]float64, error) {
	f, err := os.Open(dst)
	if err != nil {
		return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
	}

	targetDay, err := parseDateParam(date) // 解析给定的预测日期
	if err != nil {
		return [13][72]float64{}, [12][24]float64{}, code.ErrInvalidParam
	}

	var rows []parsedRow
	for idx, record := range records {
		if len(record) < 15 {
			if idx == 0 {// 跳过第一行
				continue
			}
			return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
		}

		ts, err := parseFlexibleTime(record[0]) // 解析记录时间戳
		if err != nil {
			if idx == 0 {
				continue
			}
			return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
		}

		var feature [13]float64
		for i := 0; i < 13; i++ { // 取出数值
			val, err := strconv.ParseFloat(record[i+2], 64)
			if err != nil {
				return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
			}
			feature[i] = val
		}

		rows = append(rows, parsedRow{ts: ts, feature: feature})
	}

	if len(rows) == 0 {
		return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].ts.Before(rows[j].ts) })

	pivot := -1
	for i, row := range rows {
		if row.ts.Truncate(24*time.Hour).Equal(targetDay) && row.ts.Hour() == 0 { // 找到目标日期的0点
			pivot = i
			break
		}
	}

	if pivot == -1 || pivot < 72 || pivot+24 > len(rows) {
		return [13][72]float64{}, [12][24]float64{}, code.ErrParseFile
	}

	var passData [13][72]float64
	for t := 0; t < 72; t++ { // 取前[72, 13]
		src := rows[pivot-72+t]
		for fIdx := 0; fIdx < 13; fIdx++ {
			passData[fIdx][t] = src.feature[fIdx]
		}
	}

	var futureData [12][24]float64
	for t := 0; t < 24; t++ { // 取后[24, 12]
		src := rows[pivot+t]
		for fIdx := 0; fIdx < 12; fIdx++ {
			futureData[fIdx][t] = src.feature[fIdx]
		}
	}

	return passData, futureData, nil
}

func parseFlexibleTime(val string) (time.Time, error) {
	// CSV 首列格式固定为 YYMMDDHH，例如 "12040101" -> 2012-04-01 01:00
	if len(val) != 8 {
		return time.Time{}, code.ErrParseFile
	}

	ts, err := time.ParseInLocation("06010215", val, time.Local)
	if err != nil {
		return time.Time{}, code.ErrParseFile
	}

	return ts, nil
}

func parseDateParam(val string) (time.Time, error) {
	// date格式固定为 YYMMDDHH，例如 "120401" -> 2012-04-01
	if len(val) != 6 {
		return time.Time{}, code.ErrInvalidParam
	}

	ts, err := time.ParseInLocation("060102", val, time.Local)
	if err != nil {
		return time.Time{}, code.ErrInvalidParam
	}

	return ts.Truncate(24 * time.Hour), nil
}
