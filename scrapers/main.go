package main

import (
	"fmt"
	"github.com/gookit/goutil/dump"
	"github.com/jessevdk/go-flags"
	"github.com/mfojtik/euclid/scrapers/acond"
	solar "github.com/mfojtik/euclid/scrapers/solar"
	"github.com/mfojtik/euclid/scrapers/types"
	"github.com/mfojtik/euclid/scrapers/weather"
	"log"
	"os"
	"time"
)

var opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	DisplayFile string `short:"f" long:"file"`

	AcondURL      string `long:"acond-url" default:"https://localhost:4443"`
	AcondUser     string `long:"acond-user" default:"acond"`
	AcondPassword string `long:"acond-password" default:"acond"`

	SofarPort   string `long:"sofar-port" default:"10.0.0.30:8899"`
	SofarSerial int64  `long:"sofar-serial"`

	WeatherAPI string `long:"weather-api-key"`

	MaxKeepValues int `long:"max-values" default:"5"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	if len(opts.DisplayFile) == 0 {
		fmt.Println("ERROR: No display file specified")
		os.Exit(1)
	}
	//dump.P(opts)

	heatPump := acond.New(opts.AcondURL, opts.AcondUser, opts.AcondPassword)
	heatPumpVal, err := heatPump.Scrape()
	heatPumpStatus := "on"
	if err != nil {
		heatPumpStatus = "off"
	}

	solarScraper := solar.New(opts.SofarPort, opts.SofarSerial)
	solarVal, err := solarScraper.Scrape()
	solarStatus := "on"
	if err != nil {
		solarStatus = "off"
	}

	weatherScraper := weather.New(opts.WeatherAPI)
	weatherData, err := weatherScraper.Scrape()
	if err != nil {
		log.Printf("No weather data: %v", err)
	}

	display, err := types.ReadDisplayFromFile(opts.DisplayFile)
	if os.IsNotExist(err) {
		fmt.Println("ERR: creating new file")
		display = &types.Display{
			HeatPumpState: heatPumpStatus,
			Temperatures: []types.Temperature{
				{Name: "Outside", Values: []types.Value{}, MaxValue: 30.0, MinValue: -5.0},
				{Name: "Living Room", Values: []types.Value{}, MaxValue: 24.0, MinValue: 17.0},
				{Name: "Upstairs", Values: []types.Value{}, MaxValue: 24.0, MinValue: 17.0},
				{Name: "Cellar", Values: []types.Value{}, MaxValue: 17.0, MinValue: 5.0},
			},
			Weather: types.Weather{
				Temperature:   0,
				ConditionIcon: "",
				ConditionText: "Not Available",
				Precipitation: 0,
			},
			SolarPower: types.Solar{
				Timestamp: time.Now().Unix(),
				Status:    "off",
			},
		}
	}

	if weatherData != nil {
		display.Weather = types.Weather{
			Temperature:   weatherData.Temperature,
			ConditionIcon: weatherData.ConditionIcon,
			ConditionText: weatherData.ConditionText,
			Precipitation: weatherData.Precipitation,
		}
	}

	if heatPumpStatus == "on" {
		display.HeatPumpState = heatPumpStatus
		for i := range display.Temperatures {
			switch display.Temperatures[i].Name {
			case "Outside":
				display.Temperatures[i].Values = types.RecordValue(heatPumpVal.Outdoor, display.Temperatures[i].Values, opts.MaxKeepValues)
			case "Living Room":
				display.Temperatures[i].Values = types.RecordValue(heatPumpVal.DownstairsCurrent, display.Temperatures[i].Values, opts.MaxKeepValues)
			case "Upstairs":
				display.Temperatures[i].Values = types.RecordValue(heatPumpVal.UpstairsCurrent, display.Temperatures[i].Values, opts.MaxKeepValues)
			}
			display.Temperatures[i].Trend = types.SetTrend(display.Temperatures[i].Values)
		}
	}

	if solarStatus == "on" {
		dump.P(solarVal)
		display.SolarPower.Status = solarVal.Status
		display.SolarPower.GenerationNow = solarVal.GenerationNow
		display.SolarPower.ConsumptionToday = solarVal.ConsumptionToday
		display.SolarPower.GenerationToday = solarVal.GenerationToday
		display.SolarPower.GenerationTotal = solarVal.GenerationTotal
		display.SolarPower.Timestamp = solarVal.Timestamp
	} else {
		display.SolarPower.Status = "off"
		display.SolarPower.GenerationNow = 0
		display.SolarPower.ConsumptionToday = 0
		display.SolarPower.Timestamp = time.Now().Unix()
	}

	if err := types.WriteDisplayToFile(opts.DisplayFile, *display); err != nil {
		panic(err)
	}
	//dump.P(display)
}
