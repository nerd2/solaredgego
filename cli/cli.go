package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nerd2/solaredgego"
)

func main() {
	username := os.Args[1]
	password := os.Args[2]
	cmd := ""
	if len(os.Args) > 3 {
		cmd = os.Args[3]
	}

	nh := solaredgego.NewSolarEdge(&solaredgego.Options{Username: username, Password: password})
	sites, err := nh.Login()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("%d sites\n", sites.Count)

	siteId := sites.Sites[0].Id
	switch cmd {
	case "":
		p, b, err := nh.GetData(siteId)
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("IN: %f OUT: %f BAT: Pwr:%f Charge:%d Status:%s\n", p.Grid.CurrentPower, p.Consumption.CurrentPower, p.Storage.CurrentPower, p.Storage.ChargeLevel, p.Storage.Status, b.DevicesByType.BATTERY[0].ChargeEnergy)

		log.Printf("POWER:\n%s\n", pp(p))

		log.Printf("BATTERY:\n%s\n", pp(b))
		nh.GetBatteryMode(siteId)
	case "charge":
		checkErr(nh.PutBatteryMode(siteId, solaredgego.BatteryModeCharge()))
	case "discharge":
		checkErr(nh.PutBatteryMode(siteId, solaredgego.BatteryModeDischarge()))
	case "msc":
		checkErr(nh.PutBatteryMode(siteId, solaredgego.BatteryModeMsc()))
	case "pause":
		checkErr(nh.PutBatteryMode(siteId, solaredgego.BatteryModeDisable()))
	}
}

func checkErr(err error) {
	if err == nil {
		fmt.Println("Success")
	} else {
		fmt.Println("Failed", err.Error())
	}
}

func pp(x interface{}) string {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.SetIndent("", "    ")
	e.Encode(x)
	return b.String()
}
