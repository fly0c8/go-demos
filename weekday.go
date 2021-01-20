package main

import (
	"fmt"
	"time"
)

type AssignedTariffModel struct {
	Weekday         int
	OffsetInMinutes int64
	TariffModel     string
}

type Tariffplan struct {
	TariffMax            int
	AssignedTariffModels []AssignedTariffModel
}

type TariffModelsUsed struct {
	Date        string
	Day         string
	From        string
	To          string
	Duration    string
	TariffModel string
}

var (
	tariffPlan1 = Tariffplan{
		TariffMax: 20,
		AssignedTariffModels: []AssignedTariffModel{
			{Weekday: 0, OffsetInMinutes: 0, TariffModel: "fullday"},

			{Weekday: 1, OffsetInMinutes: 0, TariffModel: "frueh"},          // 0-6h
			{Weekday: 1, OffsetInMinutes: 21600, TariffModel: "vormittag"},  // 6-12h
			{Weekday: 1, OffsetInMinutes: 43200, TariffModel: "nachmittag"}, // 12-18h
			{Weekday: 1, OffsetInMinutes: 64800, TariffModel: "abend"},      // 18-24h

			{Weekday: 2, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 2, OffsetInMinutes: 43200, TariffModel: "nachmittag"},

			{Weekday: 3, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 3, OffsetInMinutes: 43200, TariffModel: "nachmittag"},

			{Weekday: 4, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 4, OffsetInMinutes: 43200, TariffModel: "nachmittag"},

			{Weekday: 5, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 5, OffsetInMinutes: 43200, TariffModel: "nachmittag"},

			{Weekday: 6, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 6, OffsetInMinutes: 43200, TariffModel: "nachmittag"},
		},
	}
)

func day2str(day int) string {
	switch day {
	case 0:
		return "sun"
	case 1:
		return "mon"
	case 2:
		return "tue"
	case 3:
		return "wed"
	case 4:
		return "thu"
	case 5:
		return "fri"
	case 6:
		return "sat"
	default:
		panic("day invalid")
	}
}
func offsetToHMS(offsetInSeconds int64) (int64, int64, int64) {
	hours := offsetInSeconds / 3600
	remainingSecs := offsetInSeconds % 3600
	minutes := remainingSecs / 60
	seconds := remainingSecs % 60
	return hours, minutes, seconds
}
func calcDay(from time.Time, day int, fromOffset, toOffset int64) []TariffModelsUsed {

	//fromH, fromM, fromS := offsetToHMS(fromOffset)
	//toH, toM, toS := offsetToHMS(toOffset)
	//fmt.Printf("calcDay: day=%v, fromOffset=%v:%v:%v, toOffset=%v:%v:%v\n", day, fromH, fromM, fromS, toH,toM,toS)

	dayYear, dayMonth, dayDay := from.Date()
	usedTariffModels := []TariffModelsUsed{}

	tariffModels := []AssignedTariffModel{}
	for _, tm := range tariffPlan1.AssignedTariffModels {
		if tm.Weekday == day {
			tariffModels = append(tariffModels, tm)
		}
	}

	var firstIndexToTake = 0

	for i := 0; i < len(tariffModels); i++ {
		if fromOffset == tariffModels[i].OffsetInMinutes {
			firstIndexToTake = i
			break
		} else {
			var nextOffset int64
			if i == len(tariffModels)-1 {
				nextOffset = 86400
			} else {
				nextOffset = tariffModels[i+1].OffsetInMinutes
			}
			if fromOffset > tariffModels[i].OffsetInMinutes && fromOffset < nextOffset {
				firstIndexToTake = i
				break
			}
		}
	}

	var lastIndexToTake = 0

	for i := len(tariffModels) - 1; i >= 0; i-- {

		tmFrom := tariffModels[i].OffsetInMinutes
		var tmTo int64
		if i == len(tariffModels)-1 {
			tmTo = 86400
		} else {
			tmTo = tariffModels[i+1].OffsetInMinutes
		}

		if toOffset == tmTo {
			lastIndexToTake = i
			break
		} else {
			if toOffset > tmFrom && toOffset < tmTo {
				lastIndexToTake = i
				break
			}
		}
	}

	for i := firstIndexToTake; i <= lastIndexToTake; i++ {
		var duration int64
		if firstIndexToTake == lastIndexToTake {
			// sonderfall: es gibt nur 1
			duration := toOffset - fromOffset
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(fromOffset)
			toHH, toMM, toSS := offsetToHMS(toOffset)

			//fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)

			usedTariffModels = append(usedTariffModels, TariffModelsUsed{
				Date:        fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:         fmt.Sprintf("%v", day2str(day)),
				From:        fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:          fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:    fmt.Sprintf("%v:%v:%v", h, m, s),
				TariffModel: fmt.Sprintf("%v", tariffModels[i].TariffModel),
			})

		} else if i == firstIndexToTake {
			duration = tariffModels[i+1].OffsetInMinutes - fromOffset
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(fromOffset)
			toHH, toMM, toSS := offsetToHMS(tariffModels[i+1].OffsetInMinutes)

			//fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
			usedTariffModels = append(usedTariffModels, TariffModelsUsed{
				Date:        fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:         fmt.Sprintf("%v", day2str(day)),
				From:        fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:          fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:    fmt.Sprintf("%v:%v:%v", h, m, s),
				TariffModel: fmt.Sprintf("%v", tariffModels[i].TariffModel),
			})
		} else if i == lastIndexToTake {
			duration = toOffset - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(tariffModels[i].OffsetInMinutes)
			toHH, toMM, toSS := offsetToHMS(toOffset)

			//fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
			usedTariffModels = append(usedTariffModels, TariffModelsUsed{
				Date:        fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:         fmt.Sprintf("%v", day2str(day)),
				From:        fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:          fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:    fmt.Sprintf("%v:%v:%v", h, m, s),
				TariffModel: fmt.Sprintf("%v", tariffModels[i].TariffModel),
			})
		} else {
			duration = tariffModels[i+1].OffsetInMinutes - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(tariffModels[i].OffsetInMinutes)
			toHH, toMM, toSS := offsetToHMS(tariffModels[i+1].OffsetInMinutes)

			//fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
			usedTariffModels = append(usedTariffModels, TariffModelsUsed{
				Date:        fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:         fmt.Sprintf("%v", day2str(day)),
				From:        fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:          fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:    fmt.Sprintf("%v:%v:%v", h, m, s),
				TariffModel: fmt.Sprintf("%v", tariffModels[i].TariffModel),
			})
		}
	}
	return usedTariffModels

}
const(
	SECS_PER_DAY  = int64(86400)
)
func PrintDays(fromEpoch int64, toEpoch int64) {

	tariffModelsUsed := []TariffModelsUsed{}

	startEpoch := fromEpoch

	for {
		//fmt.Printf("%v --> %v\n",time.Unix(startEpoch,0) , time.Unix(toEpoch,0))

		start := time.Unix(startEpoch, 0)
		t0 := time.Date(start.Year(), start.Month(), start.Day(), 0,0,0, 0, time.Local)
		fromOffset :=  startEpoch - t0.Unix()

		toOffset := SECS_PER_DAY
		if SECS_PER_DAY > toEpoch - t0.Unix() {
			toOffset = toEpoch - t0.Unix()
		}

		tariffModelsUsed = append(tariffModelsUsed, calcDay(start, int(start.Weekday()), fromOffset, toOffset)...)
		startEpoch = startEpoch + SECS_PER_DAY-fromOffset
		if startEpoch > toEpoch {
			break
		}
	}
	for _, tmUsed := range tariffModelsUsed {
		fmt.Printf("%+v\n", tmUsed)
	}

}

func main() {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day()+2, 11, 0, 0, 0, time.Local)
	PrintDays(from.Unix(), to.Unix())
}
