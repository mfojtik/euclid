package template

import (
	"bytes"
	"github.com/mfojtik/euclid/scraper/types"
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
	return v[0].Format("02 January 2006 15:04:05")
}

func RenderTemplate(display types.Display) ([]byte, error) {
	type Template struct {
		Date       string
		Outside    template.HTML
		LivingRoom template.HTML
		Upstairs   template.HTML
		LastUpdate string
	}
	t, err := template.New("test").Parse(HTML)
	if err != nil {
		return nil, err
	}
	val := Template{
		Date:       time.Now().Format("02 January 2006"),
		Outside:    template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Outside", display.Temperatures))),
		LivingRoom: template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Living Room", display.Temperatures))),
		Upstairs:   template.HTML(types.GetLastTemperatureHTMLValue(types.GetTemperature("Upstairs", display.Temperatures))),
		LastUpdate: GetLastUpdate(display),
	}
	w := bytes.NewBuffer([]byte{})
	err = t.Execute(w, val)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
