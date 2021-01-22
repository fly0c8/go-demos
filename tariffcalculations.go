package main

import (
	"fmt"
	"math"
	"time"
)


const(
	SECS_PER_DAY  = int64(86400)
)

func main() {
	//tp := makeTariffPlan()
	//sFrom := "2021-01-18T01:00:00+01:00"
	//sTo := "2021-01-25T00:02:01+01:00"
	//from, _ := time.Parse(time.RFC3339, sFrom)
	//to, _ := time.Parse(time.RFC3339, sTo)
	//amount := CalculateTariff(tp, from.Unix(), to.Unix())
	//fmt.Println("Calculated Amount:", amount)

	tm := makeTariffModel()
	amount := tm.Calculate(60)
	fmt.Println(amount)

}

//1) const : "die erste halbe stunde 500 cent" : f(duration) = 500
//2) value pro zeiteinheit: "300 cent pro stunde" : f(duration) = duration300/60
//3) getaktete preiserhoehung: "alle angefangenen 10 minuten erhoehe um 5 cent"
//StartwithZero: amount = Math.floor(duration/interval)pricestep
//StartwithValue: amount = Math.ceil(duration/interval)*pricestep

const (
	CALCMETHOD_UNKNOWN = 0
	CALCMETHOD_CONSTANT = 1							// amount = value
	CALCMETHOD_VALUE_PER_TIMEUNIT = 2 				// amount = duration * value/timeunit
	CALCMETHOD_STEPWISE_STARTING_WITH_ZERO = 3    	// amount = Math.floor(duration/timeunit)*value
	CALCMETHOD_STEPWISE_STARTING_WITH_VALUE = 4   	// amount = Math.ceil(duration/timeunit)*value
)

type Tariffstep struct {
	Index                 int
	StepDurationInMinutes int
	ValueInCents          int
	TimeUnitInMinutes     int
	CalcMethod            int
}
type TariffModel struct {
	Uuid string
	Name string
	Description string
	Tariffsteps []Tariffstep
}
func(tm *TariffModel) Calculate(duration int) int {

	var result int
	var totalCalculatedDuration int

	for i, step := range tm.Tariffsteps {

		var stepDurationForCalc int
		if i < len(tm.Tariffsteps)-1 {
			stepDurationForCalc = step.StepDurationInMinutes
			if totalCalculatedDuration+ stepDurationForCalc > duration {
				diff := totalCalculatedDuration + stepDurationForCalc - duration
				stepDurationForCalc -= diff
			}
		} else {
			// last tariffstep used for the remaining time
			stepDurationForCalc = duration - totalCalculatedDuration
		}
		totalCalculatedDuration += stepDurationForCalc

		switch step.CalcMethod {
		case CALCMETHOD_CONSTANT:
			result += step.ValueInCents
		case CALCMETHOD_VALUE_PER_TIMEUNIT:
			result += stepDurationForCalc *step.ValueInCents/step.TimeUnitInMinutes
		case CALCMETHOD_STEPWISE_STARTING_WITH_ZERO:
			result += int(math.Floor(float64(stepDurationForCalc)/float64(step.TimeUnitInMinutes)))*step.ValueInCents
		case CALCMETHOD_STEPWISE_STARTING_WITH_VALUE:
			result += int(math.Ceil(float64(stepDurationForCalc/step.TimeUnitInMinutes)))*step.ValueInCents
		}
		if totalCalculatedDuration == duration {
			break
		}
	}
	return result
}


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

type TariffmodelsForCalc struct {
	FromEpoch int64
	ToEpoch int64
	Tariffmodels []TariffmodelForCalc
	DurationPerModelSummary map[string]int64
}
type TariffmodelForCalc struct {
	Date           string
	Day            string
	From           string
	To             string
	Duration       string
	DurationInSecs int64
	TariffModel    string
}

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
func calcDay(tp* Tariffplan, from_lt time.Time, day int, fromOffset, toOffset int64) []TariffmodelForCalc {

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
func CalculateTariff(tariffPlan *Tariffplan, fromEpoch int64, toEpoch int64) int {
	tariffmodelsToUse := GetTariffmodelsForCalculation(tariffPlan, fromEpoch, toEpoch)
	for model, duration := range tariffmodelsToUse.DurationPerModelSummary {
		fmt.Println(model, duration)
	}
	return 0
}
func GetTariffmodelsForCalculation(tariffPlan *Tariffplan, fromEpoch int64, toEpoch int64) *TariffmodelsForCalc {

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
		tariffModelsUsed = append(tariffModelsUsed, calcDay(tariffPlan, start_lt, int(start_lt.Weekday()), fromOffset, toOffset)...)
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

