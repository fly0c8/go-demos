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

var (
	tariffPlan1 = Tariffplan{
		TariffMax: 20,
		AssignedTariffModels: []AssignedTariffModel{
			{Weekday: 0, OffsetInMinutes: 0, TariffModel: "fullday"},

			{Weekday: 1, OffsetInMinutes: 0, TariffModel: "frueh"},          // 0-6h
			{Weekday: 1, OffsetInMinutes: 21600, TariffModel: "vormittag"},  // 6-12h
			{Weekday: 1, OffsetInMinutes: 43200, TariffModel: "nachmittag"}, // 12-18h
			{Weekday: 1, OffsetInMinutes: 64800, TariffModel: "abend"},      // 18-24h

			{Weekday: 2, OffsetInMinutes: 0, TariffModel: "frueh"},          // 0-6h
			{Weekday: 2, OffsetInMinutes: 21600, TariffModel: "vormittag"},  // 6-12h
			{Weekday: 2, OffsetInMinutes: 43200, TariffModel: "nachmittag"}, // 12-18h
			{Weekday: 2, OffsetInMinutes: 64800, TariffModel: "abend"},      // 18-24h

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

func offsetToHMS(offsetInSeconds int64) (int64, int64, int64) {
	hours := offsetInSeconds / 3600
	remainingSecs := offsetInSeconds % 3600
	minutes := remainingSecs / 60
	seconds := remainingSecs % 60
	return hours, minutes, seconds
}
func calcDay(day int, fromOffset, toOffset int64) {
	fmt.Printf("calcDay: day=%v, fromOffset=%v, toOffset=%v\n", day, fromOffset, toOffset)
	h1, m1, s1 := offsetToHMS(fromOffset)
	h2, m2, s2 := offsetToHMS(toOffset)
	fmt.Printf("fromOffset: %v = %v:%v:%v, toOffset: %v = %v:%v:%v\n", fromOffset, h1, m1, s1, toOffset, h2, m2, s2)

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
			fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
		} else if i == firstIndexToTake {
			duration = tariffModels[i+1].OffsetInMinutes - fromOffset
			h, m, s := offsetToHMS(duration)
			fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
		} else if i == lastIndexToTake {
			duration = toOffset - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
		} else {
			duration = tariffModels[i+1].OffsetInMinutes - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fmt.Printf("TariffModel: %v, Duration: %v:%v:%v\n", tariffModels[i], h, m, s)
		}
	}

}
func PrintDays(fromEpoch int64, toEpoch int64) {

	startEpoch := fromEpoch

	for {
		from := time.Unix(startEpoch, 0)
		fmt.Printf("Weekday=%v, %v\n", from.Weekday(), int(from.Weekday()))
		fromDay := int(from.Weekday())
		to := time.Unix(toEpoch, 0)
		toDay := int(to.Weekday())
		t0 := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local).Unix()
		fromOffset := startEpoch - t0
		var toOffset int64
		if fromDay == toDay {
			toOffset = toEpoch - t0
		} else {
			toOffset = 86400
		}
		calcDay(fromDay, fromOffset, toOffset)
		startEpoch += 86400
		if startEpoch >= toEpoch {
			break
		}
	}
}

func main() {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day()+14, 0, 0, 11, 0, time.Local)
	PrintDays(from.Unix(), to.Unix())
}
