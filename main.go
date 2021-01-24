package main

import (
	"fmt"
	"log"
	"time"
)

func main() {


	call_tariffplan()
	//call_tariffmodel()


}
func call_tariffmodel() {
	tariffModel := makeSimpleTariffModel()
	amount := tariffModel.CalculateAmountInCents(60*24)
	fmt.Println(amount)
}
func call_tariffplan() {
	tariffPlan := makeTariffPlan()
	m := make(map[string]*TariffModel)
	m["morgen"] = makeSimpleTariffModel()
	m["vormittag"] = makeSimpleTariffModel()
	m["nachmittag"] = makeSimpleTariffModel()
	m["abend"] = makeSimpleTariffModel()
	m["feiertag"] = makeSimpleTariffModel()
	m["fullday"] = makeSimpleTariffModel()
	tariffPlan.SetTariffModelMap(m)
	sFrom := "2021-01-18T01:00:00+01:00"
	sTo := "2021-01-25T00:02:01+01:00"
	//sFrom := 	"2021-01-18T00:00:00+01:00"
	//sTo := 		"2021-01-20T00:00:00+01:00"
	from, _ := time.Parse(time.RFC3339, sFrom)
	to, _ := time.Parse(time.RFC3339, sTo)
	amount, err := tariffPlan.CalculateAmountInCents(from.Unix(), to.Unix())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Calculated Amount:", amount)
}

