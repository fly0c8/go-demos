package main

import (
	"fmt"
	"testing"
	"time"
)



//func Test1(t *testing.T) {
//	tp := makeTariffPlan()
//
//	from, _ := time.Parse(time.RFC3339, "2021-01-21T13:00:00+01:00")
//	to, _ := time.Parse(time.RFC3339, "2021-01-25T00:01:01+01:00")
//	tms := GetTariffmodelsForCalculation(tp, from.Unix(), to.Unix())
//	prettyPrint(tms)
//}

func Test_CalcTariffModel_CorrectAmountReturned(t *testing.T) {
	tm := makeTariffModel()

	got := tm.Calculate(3)
	expected := 7
	if got != expected {
		t.Errorf("Expected: %d, Got: %d", expected, got)
	}

	got = tm.Calculate(20)
	expected = 57
	if got != expected {
		t.Errorf("Expected: %d, Got: %d", expected, got)
	}

	got = tm.Calculate(60)
	expected = 257
	if got != expected {
		t.Errorf("Expected: %d, Got: %d", expected, got)
	}
}

func Test_FindTariffModelsAndDurations_CorrecTariffmodelsAreUsed(t *testing.T) {

	tp := makeTariffPlan()

	sFrom := "2021-01-18T01:00:00+01:00"
	sTo := "2021-01-25T00:02:01+01:00"

	from, err := time.Parse(time.RFC3339, sFrom)
	if err != nil {
		t.Fatal(err)
	}
	to, err := time.Parse(time.RFC3339, sTo)
	if err != nil {
		t.Fatal(err)
	}
	//tariffmodelsForCalculation := GetTariffmodelsForCalculation(tp, from.Unix(), to.Unix())
	tariffmodelsForCalculation := tp.FindTariffmodelsForCalculation(from.Unix(), to.Unix())

	got := len(tariffmodelsForCalculation.Tariffmodels)
	want := 13
	if got != want {
		prettyPrint(tariffmodelsForCalculation)
		t.Errorf("lenTariffModelsUsed: got: %v, want: %v", got, want)
	}

	wantedValues := []struct {
		Day            string
		DurationInSecs int64
		TariffModel    string
	}{
		{
			Day:            "mon",
			DurationInSecs: 18000,
			TariffModel:    "frueh",
		},
		{
			Day:            "mon",
			DurationInSecs: 21600,
			TariffModel:    "vormittag",
		},
		{
			Day:            "mon",
			DurationInSecs: 21600,
			TariffModel:    "nachmittag",
		},
		{
			Day:            "mon",
			DurationInSecs: 21600,
			TariffModel:    "abend",
		},
		{
			Day:            "tue",
			DurationInSecs: 43200,
			TariffModel:    "vormittag",
		},
		{
			Day:            "tue",
			DurationInSecs: 43200,
			TariffModel:    "nachmittag",
		},
		{
			Day:            "wed",
			DurationInSecs: 43200,
			TariffModel:    "vormittag",
		},
		{
			Day:            "wed",
			DurationInSecs: 43200,
			TariffModel:    "nachmittag",
		},
		{
			Day:            "thu",
			DurationInSecs: 86400,
			TariffModel:    "feiertag",
		},
		{
			Day:            "fri",
			DurationInSecs: 86400,
			TariffModel:    "feiertag",
		},
		{
			Day:            "sat",
			DurationInSecs: 86400,
			TariffModel:    "fullday",
		},
		{
			Day:            "sun",
			DurationInSecs: 86400,
			TariffModel:    "fullday",
		},
		{
			Day:            "mon",
			DurationInSecs: 121,
			TariffModel:    "feiertag",
		},
	}

	for i := 0; i<len(tariffmodelsForCalculation.Tariffmodels); i ++ {
		if tariffmodelsForCalculation.Tariffmodels[i].Day != wantedValues[i].Day {
			t.Errorf("wrong day at index: %v. %v != %v",i, tariffmodelsForCalculation.Tariffmodels[i].Day, wantedValues[i].Day)
		}
		if tariffmodelsForCalculation.Tariffmodels[i].DurationInSecs != wantedValues[i].DurationInSecs {
			t.Errorf("wrong DurationInSecs at index: %v. %v != %v",i, tariffmodelsForCalculation.Tariffmodels[i].DurationInSecs, wantedValues[i].DurationInSecs)
		}
		if tariffmodelsForCalculation.Tariffmodels[i].TariffModel != wantedValues[i].TariffModel {
			t.Errorf("wrong TariffModel at index: %v. %v != %v",i, tariffmodelsForCalculation.Tariffmodels[i].TariffModel, wantedValues[i].TariffModel)
		}
	}


}
func prettyPrint(tariffModelsUsed *TariffmodelsForCalc) {
	lastWeekday := ""
	for _, tmUsed := range tariffModelsUsed.Tariffmodels {
		if lastWeekday != tmUsed.Day {
			fmt.Println("------------------------------------------------------------------------------------------------")
			lastWeekday = tmUsed.Day
		}
		fmt.Printf("%+v\n", tmUsed)
	}
	fmt.Println("=== Summary ===")
	fmt.Printf("Timerange: From: %v, To: %v\n", time.Unix(tariffModelsUsed.FromEpoch,0), time.Unix(tariffModelsUsed.ToEpoch,0))
	for k, v := range tariffModelsUsed.DurationPerModelSummary {
		h,m,s := offsetToHMS(v)
		fmt.Printf("%v => %v \t(%vh:%02vm:%02vs)\n", k,v, h,m,s)
	}


	fmt.Println("=========================================================================================================")

}
