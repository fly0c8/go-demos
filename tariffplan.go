package main

import (
	"fmt"
	"time"
)

type AssignedTariffModel struct {
	OffsetInMinutes int64
	TariffModelUuid string
}
type WeekdayModel struct {
	Weekday int
	AssignedTariffModels []AssignedTariffModel
}
type ExceptiondayModel struct {
	LocalDate time.Time
	Name string
	AssignedTariffModels []AssignedTariffModel
}
type Tariffplan struct {
	Uuid string
	Name string
	Description string
	MaxTariff int
	ValidFromEpoch int64
	WeekdayModels []WeekdayModel
	ExceptiondayModels []ExceptiondayModel

}

func findTariffModelByUuid(uuid string) (*TariffModel, error) {
	return makeSimpleTariffModel(), nil
}
func (tp *Tariffplan) CalculateAmountInCents(fromEpoch, toEpoch int64) (int64, error) {
	tms := tp.findTariffmodelsForCalculation(fromEpoch, toEpoch)
	var amount int64
	for tmUuid, duration := range tms.DurationPerModelSummary {
		tm, err := findTariffModelByUuid(tmUuid)
		if err != nil {
			return 0, err
		}
		amount += tm.CalculateAmountInCents(duration/60)
	}
	return amount, nil
}

func (tp *Tariffplan) findTariffmodelsForCalculation(fromEpoch int64, toEpoch int64) *TariffmodelsForCalc {

	tariffModelsUsed := []TariffmodelForCalc{}
	startEpoch := fromEpoch

	for {
		start_lt := time.Unix(startEpoch, 0)
		t0 := time.Date(start_lt.Year(), start_lt.Month(), start_lt.Day(), 0,0,0, 0, time.Local)
		fromOffset :=  startEpoch - t0.Unix()

		toOffset := SECS_PER_DAY
		if SECS_PER_DAY > toEpoch - t0.Unix() {
			toOffset = toEpoch - t0.Unix()
		}
		tariffModelsUsed = append(tariffModelsUsed, tp.calcDay(start_lt, int(start_lt.Weekday()), fromOffset, toOffset)...)
		startEpoch = startEpoch + SECS_PER_DAY-fromOffset
		if startEpoch >= toEpoch {
			break
		}
	}
	res := TariffmodelsForCalc{
		FromEpoch: fromEpoch,
		ToEpoch: toEpoch,
		Tariffmodels:            tariffModelsUsed,
		DurationPerModelSummary: make(map[string]int64),
	}
	for _, tm := range res.Tariffmodels {
		res.DurationPerModelSummary[tm.TariffModel] += tm.DurationInSecs
	}

	return &res
}
func (tp* Tariffplan) calcDay(from_lt time.Time, day int, fromOffset, toOffset int64) []TariffmodelForCalc {

	//fromH, fromM, fromS := offsetToHMS(fromOffset)
	//toH, toM, toS := offsetToHMS(toOffset)
	//fmt.Printf("calcDay: day=%v, fromOffset=%v:%v:%v, toOffset=%v:%v:%v\n", day, fromH, fromM, fromS, toH,toM,toS)

	dayYear, dayMonth, dayDay := from_lt.Date()
	tariffmodelForCalcs := []TariffmodelForCalc{}

	tariffModels := []AssignedTariffModel{}

	// Check if current day is exception day
	// If yes, use Tariffmodels from Exception day
	// else use Tariffmodels from WeekdayModels
	exceptiondayModelIndex := -1
	for i, em := range tp.ExceptiondayModels {
		emYear, emMonth, emDay := em.LocalDate.Date()
		if emYear == dayYear && emMonth == dayMonth && emDay == dayDay {
			exceptiondayModelIndex = i
		}
	}
	if exceptiondayModelIndex != -1 {
		for _, tm := range tp.ExceptiondayModels[exceptiondayModelIndex].AssignedTariffModels {
			tariffModels = append(tariffModels, tm)
		}
	} else {
		// just a normal weekday
		for _, wm := range tp.WeekdayModels {
			if wm.Weekday == day {
				for _, tm := range wm.AssignedTariffModels {
					tariffModels = append(tariffModels, tm)
				}
			}
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

			tariffmodelForCalcs = append(tariffmodelForCalcs, TariffmodelForCalc{
				Date:           fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:            fmt.Sprintf("%v", day2str(day)),
				From:           fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:             fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:       fmt.Sprintf("%v:%v:%v", h, m, s),
				DurationInSecs: duration,
				TariffModel:    fmt.Sprintf("%v", tariffModels[i].TariffModelUuid),
			})

		} else if i == firstIndexToTake {
			duration = tariffModels[i+1].OffsetInMinutes - fromOffset
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(fromOffset)
			toHH, toMM, toSS := offsetToHMS(tariffModels[i+1].OffsetInMinutes)

			//fmt.Printf("TariffModel: %v, StepDurationInMinutes: %v:%v:%v\n", tariffModels[i], h, m, s)
			tariffmodelForCalcs = append(tariffmodelForCalcs, TariffmodelForCalc{
				Date:           fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:            fmt.Sprintf("%v", day2str(day)),
				From:           fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:             fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:       fmt.Sprintf("%v:%v:%v", h, m, s),
				DurationInSecs: duration,
				TariffModel:    fmt.Sprintf("%v", tariffModels[i].TariffModelUuid),
			})
		} else if i == lastIndexToTake {
			duration = toOffset - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(tariffModels[i].OffsetInMinutes)
			toHH, toMM, toSS := offsetToHMS(toOffset)

			//fmt.Printf("TariffModel: %v, StepDurationInMinutes: %v:%v:%v\n", tariffModels[i], h, m, s)
			tariffmodelForCalcs = append(tariffmodelForCalcs, TariffmodelForCalc{
				Date:           fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:            fmt.Sprintf("%v", day2str(day)),
				From:           fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:             fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:       fmt.Sprintf("%v:%v:%v", h, m, s),
				DurationInSecs: duration,
				TariffModel:    fmt.Sprintf("%v", tariffModels[i].TariffModelUuid),
			})
		} else {
			duration = tariffModels[i+1].OffsetInMinutes - tariffModels[i].OffsetInMinutes
			h, m, s := offsetToHMS(duration)
			fromHH, fromMM, fromSS := offsetToHMS(tariffModels[i].OffsetInMinutes)
			toHH, toMM, toSS := offsetToHMS(tariffModels[i+1].OffsetInMinutes)

			//fmt.Printf("TariffModel: %v, StepDurationInMinutes: %v:%v:%v\n", tariffModels[i], h, m, s)
			tariffmodelForCalcs = append(tariffmodelForCalcs, TariffmodelForCalc{
				Date:           fmt.Sprintf("%v/%v/%v", dayYear, dayMonth, dayDay),
				Day:            fmt.Sprintf("%v", day2str(day)),
				From:           fmt.Sprintf("%v:%v:%v", fromHH, fromMM, fromSS),
				To:             fmt.Sprintf("%v:%v:%v", toHH, toMM, toSS),
				Duration:       fmt.Sprintf("%v:%v:%v", h, m, s),
				DurationInSecs: duration,
				TariffModel:    fmt.Sprintf("%v", tariffModels[i].TariffModelUuid),
			})
		}
	}

	return tariffmodelForCalcs
}
