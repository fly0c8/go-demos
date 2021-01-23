package main

import "math"

//1) const : "die erste halbe stunde 500 cent" : f(duration) = 500
//2) value pro zeiteinheit: "300 cent pro stunde" : f(duration) = duration300/60
//3) getaktete preiserhoehung: "alle angefangenen 10 minuten erhoehe um 5 cent"
//StartwithZero: amount = Math.floor(duration/interval)pricestep
//StartwithValue: amount = Math.ceil(duration/interval)*pricestep
const (
	CALCMETHOD_CONSTANT = 0							// amount = value
	CALCMETHOD_VALUE_PER_TIMEUNIT = 1 				// amount = duration * value/timeunit
	CALCMETHOD_STEPWISE_STARTING_WITH_ZERO = 2    	// amount = Math.floor(duration/timeunit)*value
	CALCMETHOD_STEPWISE_STARTING_WITH_VALUE = 3   	// amount = Math.ceil(duration/timeunit)*value
)

type Tariffstep struct {
	Index                 int
	StepDurationInMinutes int
	ValueInCents          int
	TimeIntervalInMinutes int
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
			result += stepDurationForCalc *step.ValueInCents/step.TimeIntervalInMinutes
		case CALCMETHOD_STEPWISE_STARTING_WITH_ZERO:
			result += int(math.Floor(float64(stepDurationForCalc)/float64(step.TimeIntervalInMinutes)))*step.ValueInCents
		case CALCMETHOD_STEPWISE_STARTING_WITH_VALUE:
			result += int(math.Ceil(float64(stepDurationForCalc/step.TimeIntervalInMinutes)))*step.ValueInCents
		}
		if totalCalculatedDuration == duration {
			break
		}
	}
	return result
}

