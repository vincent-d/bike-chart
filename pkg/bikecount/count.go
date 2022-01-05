package bikecount

import (
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/vincent-d/bike-count/pkg/ecovisio"
)

type BikeCountData struct {
	Name           string
	Latitude       float32
	Longitude      float32
	DisplayDates   []string
	Counts         []opts.BarData
	Start          int // starting point percentage
	End            int // ending point percentage
	OverallAverage int
	WindowAverage  int
	WindowMax      int
	MaxDate        time.Time
	WeekdaysMax    int
	WeekdaysAvg    int
	WeekendsAvg    int
	WeekendsMax    int
	DaysAverage    [7]int
	DaysMax        [7]int
	BestDay        time.Weekday
}

type Average struct {
	value float64
	nb    int
}

func (avg *Average) addValue(val float64) {
	avg.value += val
	avg.nb++
}

func (avg Average) compute() float64 {
	if avg.nb != 0 {
		return avg.value / float64(avg.nb)
	}
	return 0.0
}

func isSameDate(a time.Time, b time.Time) bool {
	if a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day() {
		return true
	}
	return false
}

func PrepareBikeCountData(totemInfo []ecovisio.TotemInfo, bikeCount [][]ecovisio.BikeCount, start time.Time, end time.Time) (bikeCountData []BikeCountData) {

	for i, info := range totemInfo {
		if len(bikeCount[i]) == 0 {
			continue
		}
		countData := BikeCountData{
			Name:           info.Name,
			OverallAverage: info.DailyAverage,
			Latitude:       info.Latitude,
			Longitude:      info.Longitude,
			Start:          0,
			End:            100,
		}
		windowAvg := Average{}
		weekDays := Average{}
		weekends := Average{}
		var days [7]Average
		for j, value := range bikeCount[i] {
			if isSameDate(start, value.Date) {
				countData.Start = int(float32(j) / float32(len(bikeCount[i])) * 100.0)
			}
			if isSameDate(end, value.Date) {
				countData.End = int(float32(j) / float32(len(bikeCount[i])) * 100.0)
			}
			countData.DisplayDates = append(countData.DisplayDates, value.Date.Format("02/01/2006"))
			countData.Counts = append(countData.Counts, opts.BarData{Value: value.Count})
			windowAvg.addValue(float64(value.Count))
			if value.Count > countData.WindowMax {
				countData.WindowMax = value.Count
				countData.MaxDate = value.Date
			}
			switch value.Date.Weekday() {
			case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
				weekDays.addValue(float64(value.Count))
				if value.Count > countData.WeekdaysMax {
					countData.WeekdaysMax = value.Count
				}
			case time.Saturday, time.Sunday:
				weekends.addValue(float64(value.Count))
				if value.Count > countData.WeekendsMax {
					countData.WeekendsMax = value.Count
				}
			}
			days[int(value.Date.Weekday())].addValue(float64(value.Count))
		}
		maxDay := 0
		for j := range days {
			countData.DaysAverage[j] = int(days[j].compute())
			if countData.DaysAverage[j] > maxDay {
				maxDay = countData.DaysAverage[j]
				countData.BestDay = time.Weekday(j)
			}
		}
		countData.WindowAverage = int(windowAvg.compute())
		countData.WeekdaysAvg = int(weekDays.compute())
		countData.WeekendsAvg = int(weekends.compute())
		bikeCountData = append(bikeCountData, countData)
	}

	return bikeCountData
}
