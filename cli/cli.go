package main

import (
	"github.com/nerd2/solaredgego"
	"log"
	"os"
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

	log.Printf("IN: %f OUT: %f BAT: %f\n", p.Grid.CurrentPower, p.Consumption.CurrentPower, b.DevicesByType.BATTERY[0].ChargeEnergy)
}
