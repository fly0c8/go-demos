package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_FindTariffModelsAndDurations_CorrecTariffmodelsAreUsed(t *testing.T) {

	tp := Tariffplan{
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

			{Weekday: 6, OffsetInMinutes: 0, TariffModel: "fullday"},
		},
	}

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
	fmt.Println("************************************************************************")
	fmt.Println("From:", sFrom)
	fmt.Println("To:", sTo)
	fmt.Println("************************************************************************")
	tariffModelsUsed := FindTariffModelsAndDurations(&tp, from.Unix(), to.Unix())

	got := len(tariffModelsUsed)
	want := 15
	if got != want {
		prettyPrint(tariffModelsUsed)
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
			DurationInSecs: 43200,
			TariffModel:    "vormittag",
		},
		{
			Day:            "thu",
			DurationInSecs: 43200,
			TariffModel:    "nachmittag",
		},
		{
			Day:            "fri",
			DurationInSecs: 43200,
			TariffModel:    "vormittag",
		},
		{
			Day:            "fri",
			DurationInSecs: 43200,
			TariffModel:    "nachmittag",
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
			TariffModel:    "frueh",
		},
	}
	for i := 0; i<len(tariffModelsUsed); i ++ {
		if tariffModelsUsed[i].Day != wantedValues[i].Day {
			t.Errorf("wrong day at index: %v. %v != %v",i, tariffModelsUsed[i].Day, wantedValues[i].Day)
		}
		if tariffModelsUsed[i].DurationInSecs != wantedValues[i].DurationInSecs {
			t.Errorf("wrong DurationInSecs at index: %v. %v != %v",i, tariffModelsUsed[i].DurationInSecs, wantedValues[i].DurationInSecs)
		}
		if tariffModelsUsed[i].TariffModel != wantedValues[i].TariffModel {
			t.Errorf("wrong TariffModel at index: %v. %v != %v",i, tariffModelsUsed[i].TariffModel, wantedValues[i].TariffModel)
		}
	}

}
func prettyPrint(tariffModelsUsed []TariffModelsUsed) {
	lastWeekday := ""
	for _, tmUsed := range tariffModelsUsed {
		if lastWeekday != tmUsed.Day {
			fmt.Println("------------------------------------------------------------------------------------------------")
			lastWeekday = tmUsed.Day
		}
		fmt.Printf("%+v\n", tmUsed)
	}
}