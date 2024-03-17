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
	solarTime := time.Unix(display.SolarPower.Timestamp, 0).Format("02 January 2006 15:04:05")
	return fmt.Sprintf("solar: %s | temp:%s", solarTime, v[0].Format("02 January 2006 15:04:05"))
}

func GetSolarStatus(status string) template.HTML {
	if status == "on" {
		return `<td class="w-50 text-md-end bg-success">ONLINE</td>`
	}
	return `<td class="w-50 text-md-end bg-danger">OFFLINE</td>`
}

func RenderTemplate(display types.Display) ([]byte, error) {
	type Template struct {
		Date       string
		Outside    template.HTML
		LivingRoom template.HTML
		Upstairs   template.HTML
		LastUpdate string

		SolarStatus template.HTML
		SolarToday  string
		SolarTotal  string
		SolarNow    string
	}
	t, err := template.New("test").Parse(HTML)
	if err != nil {
		return nil, err
	}
	val := Template{
		Date:        time.Now().Format("02 January 2006"),
		Outside:     template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Outside", display.Temperatures))),
		LivingRoom:  template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Living Room", display.Temperatures))),
		Upstairs:    template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Upstairs", display.Temperatures))),
		LastUpdate:  GetLastUpdate(display),
		SolarStatus: GetSolarStatus(display.SolarPower.Status),
		SolarNow:    fmt.Sprintf("%.2f kW", display.SolarPower.GenerationNow),
		SolarTotal:  fmt.Sprintf("%.2f kWh", display.SolarPower.GenerationTotal),
		SolarToday:  fmt.Sprintf("%.2f kWh", display.SolarPower.GenerationToday),
	}
	w := bytes.NewBuffer([]byte{})
	err = t.Execute(w, val)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
