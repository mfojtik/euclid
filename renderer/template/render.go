package template

import (
	"bytes"
	"fmt"
	"github.com/mfojtik/euclid/scrapers/types"
	"html/template"
	"sort"
	"time"
)

func GetLastUpdate(display types.Display) string {
	v := []time.Time{}
	for i := range display.Temperatures {
		for d := range display.Temperatures[i].Values {
			v = append(v, display.Temperatures[i].Values[d].Timestamp)
		}
	}
	sort.Slice(v, func(i, j int) bool {
		return v[j].Before(v[i])
	})
	if len(v) == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%s", v[0].Format("02 January 2006 15:04"))
}

func GetStatus(status string) template.HTML {
	if status == "on" {
		return `<td style="width:30%" class="text-md-end bg-success">ONLINE</td>`
	}
	return `<td style="width:30%" class="text-md-end bg-danger"><b>OFFLINE</b></td>`
}

func RenderTemplate(display types.Display) ([]byte, error) {
	type Template struct {
		Date       string
		Outside    template.HTML
		LivingRoom template.HTML
		Upstairs   template.HTML
		LastUpdate string

		SolarStatus   template.HTML
		HeatPumpState template.HTML

		Weather template.HTML

		SolarToday  string
		SolarTotal  string
		SolarNow    string
		EnergyToday string
	}
	t, err := template.New("test").Parse(HTML)
	if err != nil {
		return nil, err
	}
	val := Template{
		Date:          time.Now().Format("02 January 2006"),
		Outside:       template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Outside", display.Temperatures))),
		LivingRoom:    template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Living Room", display.Temperatures))),
		Upstairs:      template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Upstairs", display.Temperatures))),
		LastUpdate:    GetLastUpdate(display),
		SolarStatus:   GetStatus(display.SolarPower.Status),
		HeatPumpState: GetStatus(display.HeatPumpState),
		SolarNow:      fmt.Sprintf("%.2f kW", display.SolarPower.GenerationNow),
		EnergyToday:   fmt.Sprintf("%.2f kWh", display.SolarPower.ConsumptionToday),
		SolarTotal:    fmt.Sprintf("%.2f kWh", display.SolarPower.GenerationTotal),
		SolarToday:    fmt.Sprintf("%.2f kWh", display.SolarPower.GenerationToday),
		Weather: template.HTML(fmt.Sprintf("<img src=\"%s\" style=\"width:24px\"/> %s, %.1f °C, %d%% chance of rain",
			display.Weather.ConditionIcon, display.Weather.ConditionText, display.Weather.Temperature, display.Weather.Precipitation)),
	}
	w := bytes.NewBuffer([]byte{})
	err = t.Execute(w, val)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
