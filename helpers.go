package main

import "time"
func makeTariffModel() *TariffModel{
	tm := &TariffModel{
		Uuid:        "u1",
		Name:        "tm1",
		Description: "feiertagsmodell",
		Tariffsteps: []Tariffstep{
			{
				Index:                 0,
				StepDurationInMinutes: 10,
				ValueInCents:          7,
				TimeUnitInMinutes:     0,
				CalcMethod:            CALCMETHOD_CONSTANT,
			},
			{
				Index:                 1,
				StepDurationInMinutes: 20,
				ValueInCents:          300,
				TimeUnitInMinutes:     60,
				CalcMethod:            CALCMETHOD_VALUE_PER_TIMEUNIT,
			},
			// every 10 minutes, raise value by 50 cents, starting with 50 cents
			{
				Index:                 2,
				StepDurationInMinutes: 0,
				ValueInCents:          50,
				TimeUnitInMinutes:     10,
				CalcMethod:            CALCMETHOD_STEPWISE_STARTING_WITH_VALUE,
			},

		},
	}
	return tm
}
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