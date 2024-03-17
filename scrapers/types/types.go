package types

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Display struct {
	Temperatures []Temperature `json:"temperatures,omitempty"`
	SolarPower   Solar         `json:"solar"`
}

type Value struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type Solar struct {
	Timestamp       int64   `json:"timestamp"`
	Status          string  `json:"status"`
	GenerationNow   float32 `json:"generation_now"`
	GenerationTotal float32 `json:"generation_total"`
	GenerationToday float32 `json:"generation_today"`
}

type Temperature struct {
	Name     string  `json:"name"`
	Values   []Value `json:"values"`
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
	Trend    string  `json:"trend"`
}

func ReadDisplayFromFile(file string) (*Display, error) {
	fb, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var display Display
	if err := json.Unmarshal(fb, &display); err != nil {
		return nil, err
	}
	return &display, nil
}

func WriteDisplayToFile(file string, display Display) error {
	fb, err := json.Marshal(display)
	if err != nil {
		return err
	}
	return os.WriteFile(file, fb, os.ModePerm)
}

func SetTrend(values []Value) string {
	if len(values) == 0 || len(values) == 1 {
		return "stale"
	}
	sum := float64(0)
	count := float64(0)
	for _, v := range values[1 : len(values)-1] {
		sum += v.Value
		count++
	}
	avg := sum / count
	switch {
	case values[0].Value > avg:
		return "up"
	case values[0].Value < avg:
		return "down"
	default:
		return "stale"
	}
}

var (
	staleHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-three-dots" viewBox="0 0 16 16"><path d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3"/></svg>`
	downHTML  = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-arrow-down" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M8 1a.5.5 0 0 1 .5.5v11.793l3.146-3.147a.5.5 0 0 1 .708.708l-4 4a.5.5 0 0 1-.708 0l-4-4a.5.5 0 0 1 .708-.708L7.5 13.293V1.5A.5.5 0 0 1 8 1"/></svg>`
	upHTML    = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-arrow-up" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M8 15a.5.5 0 0 0 .5-.5V2.707l3.146 3.147a.5.5 0 0 0 .708-.708l-4-4a.5.5 0 0 0-.708 0l-4 4a.5.5 0 1 0 .708.708L7.5 2.707V14.5a.5.5 0 0 0 .5.5"/></svg>`
)

func GetLastTemperatureHTMLValue(t *Temperature) string {
	if t == nil {
		return fmt.Sprintf(`<th scope="row">N/A</th><td style="width:27%%" class="text-md-end">0 ℃ %s</td>`, staleHTML)
	}
	sign := staleHTML
	trendColor := ""
	switch t.Trend {
	case "up":
		sign = upHTML
		if t.Values[0].Value >= t.MaxValue {
			trendColor = " bg-warning-subtle"
		}
	case "down":
		sign = downHTML
		if t.Values[0].Value <= t.MinValue {
			trendColor = " bg-info"
		}
	}
	o := fmt.Sprintf(`<th scope="row">%s</th><td style="width:27%%" class="text-md-end%s">%2.1f ℃ %s</td>`, t.Name, trendColor, t.Values[0].Value, sign)
	//dump.P(o)
	return o
}

func GetTemperature(name string, t []Temperature) *Temperature {
	for i := range t {
		if t[i].Name == name {
			return &t[i]
		}
	}
	return nil
}

func RecordValue(v float64, values []Value, maxValues int) []Value {
	n := []Value{
		{
			Timestamp: time.Now(),
			Value:     v,
		},
	}
	if len(values) == 0 {
		return n
	}
	sort.Slice(values, func(i, j int) bool {
		return values[j].Timestamp.Before(values[i].Timestamp)
	})
	count := 1
	for i := range values {
		n = append(n, values[i])
		count++
		if count >= maxValues {
			break
		}
	}
	return n
}
