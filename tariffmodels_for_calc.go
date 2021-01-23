package main

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
