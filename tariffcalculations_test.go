package main

import (
	"fmt"
	"testing"
	"time"
)


func makeTariffPlan() *Tariffplan {

	exceptiondayThursday, _ := time.Parse(time.RFC3339, "2021-01-21T01:00:00+01:00")
	exceptiondayFriday, _ := time.Parse(time.RFC3339, "2021-01-22T01:00:00+01:00")
	exceptiondayMonday, _ := time.Parse(time.RFC3339, "2021-01-25T01:00:00+01:00")

	return &Tariffplan{
		MaxTariff: 20,
		WeekdayModels: []WeekdayModel{
			{ Weekday: 0, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "fullday"},
			}},
			{ Weekday: 1, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "frueh"},
				{OffsetInMinutes: 21600, TariffModelUuid: "vormittag"},
				{OffsetInMinutes: 43200, TariffModelUuid: "nachmittag"},
				{OffsetInMinutes: 64800, TariffModelUuid: "abend"},
			}},
			{ Weekday: 2, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "vormittag"},
				{OffsetInMinutes: 43200, TariffModelUuid: "nachmittag"},
			}},
			{ Weekday: 3, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "vormittag"},
				{OffsetInMinutes: 43200, TariffModelUuid: "nachmittag"},
			}},
			{ Weekday: 4, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "vormittag"},
				{OffsetInMinutes: 43200, TariffModelUuid: "nachmittag"},
			}},
			{ Weekday: 5, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "vormittag"},
				{OffsetInMinutes: 43200, TariffModelUuid: "nachmittag"},
			}},
			{ Weekday: 6, AssignedTariffModels: []AssignedTariffModel{
				{OffsetInMinutes: 0, TariffModelUuid: "fullday"},
			}},
		},
		ExceptiondayModels: []ExceptiondayModel{
			{
				LocalDate: exceptiondayThursday,
				Name:                 "DonnerstagsFeiertag",
				AssignedTariffModels: []AssignedTariffModel{{
					OffsetInMinutes: 0,
					TariffModelUuid: "feiertag",
				}},
			},
			{
				LocalDate: exceptiondayFriday,
				Name:                 "FreitagsFeiertag",
				AssignedTariffModels: []AssignedTariffModel{{
					OffsetInMinutes: 0,
					TariffModelUuid: "feiertag",
				}},
			},
			{
				LocalDate: exceptiondayMonday,
				Name:                 "FreitagsMontag",
				AssignedTariffModels: []AssignedTariffModel{{
					OffsetInMinutes: 0,
					TariffModelUuid: "feiertag",
				}},
			},
		},
	}
}
func Test1(t *testing.T) {
	tp := makeTariffPlan()

	from, _ := time.Parse(time.RFC3339, "2021-01-21T13:00:00+01:00")
	to, _ := time.Parse(time.RFC3339, "2021-01-25T00:01:01+01:00")
	tms := GetTariffmodelsForCalculation(tp, from.Unix(), to.Unix())
	prettyPrint(tms)
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
	tariffmodelsForCalculation := GetTariffmodelsForCalculation(tp, from.Unix(), to.Unix())

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

//sparkApiServer=# select * from tariffstructure where id=198;
//id  |                 uuid                 | tariffstructure | description |          created           |          updated           | deleted
//-----+--------------------------------------+-----------------+-------------+----------------------------+----------------------------+---------
//198 | 8815fd6a-6e13-4199-9d01-72ad47b82603 | arnitest        | asdfadf     | 2021-01-20 19:04:57.509739 | 2021-01-20 19:05:38.029605 |
//(1 row)
//
//sparkApiServer=# select * from tariffstep where tariffstructure_id=198;
//id  | uuid | tariffstructure_id | index | duration | multiplier | value |          created           | updated | deleted
//-----+------+--------------------+-------+----------+------------+-------+----------------------------+---------+---------
//201 |      |                198 |     0 |      600 |          3 |   100 | 2021-01-20 19:05:38.029605 |         |
//202 |      |                198 |     1 |      600 |          1 |   200 | 2021-01-20 19:05:38.029605 |         |
