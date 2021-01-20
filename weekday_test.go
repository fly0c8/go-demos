package main

import (
	"fmt"
	"testing"
	"time"
)


func Test_FindTariffModelsAndDurations(t *testing.T) {
	from, err := time.Parse(time.RFC3339, "2021-01-20T01:00:00+01:00")
	if err != nil {
		t.Fatal(err)
	}
	to, err := time.Parse(time.RFC3339, "2021-01-21T01:00:02+01:00")
	if err != nil {
		t.Fatal(err)
	}

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

			{Weekday: 6, OffsetInMinutes: 0, TariffModel: "vormittag"},
			{Weekday: 6, OffsetInMinutes: 43200, TariffModel: "nachmittag"},
		},
	}

	tariffModelsUsed := FindTariffModelsAndDurations(&tp, from.Unix(), to.Unix())

	lastWeekday := ""
	for _, tmUsed := range tariffModelsUsed {
		if lastWeekday != tmUsed.Day {
			fmt.Println("--------------------------------")
			lastWeekday = tmUsed.Day
		}
		fmt.Printf("%+v\n", tmUsed)
	}
}

//func main() {
//
//
//	//now := time.Now()
//	//from := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.Local)
//	//to := time.Date(now.Year(), now.Month(), now.Day()+5, 12, 0, 0, 1, time.Local)
//
//
//}
