package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/nerd2/solaredgego"
)

func main() {
	username := os.Args[1]
	password := os.Args[2]

	nh := solaredgego.NewSolarEdge(&solaredgego.Options{Username: username, Password: password})
	sites, err := nh.Login()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("%d sites\n", sites.Count)

	p, b, err := nh.GetData(sites.Sites[0].Id)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("IN: %f OUT: %f BAT: Pwr:%f Charge:%d Status:%s\n", p.Grid.CurrentPower, p.Consumption.CurrentPower, p.Storage.CurrentPower, p.Storage.ChargeLevel, p.Storage.Status, b.DevicesByType.BATTERY[0].ChargeEnergy)

	log.Printf("POWER:\n%s\n", pp(p))

	log.Printf("BATTERY:\n%s\n", pp(b))
}

func pp(x interface{}) string {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.SetIndent("", "    ")
	e.Encode(x)
	return b.String()
}
