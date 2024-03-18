package weather

import (
	resty "github.com/go-resty/resty/v2"
	"github.com/mfojtik/euclid/scrapers/types"
)

type Scraper struct {
	weatherAPI string
}

func New(weatherAPI string) *Scraper {
	return &Scraper{weatherAPI: weatherAPI}
}

func (s *Scraper) Scrape() (*types.Weather, error) {
	weather := &WeatherData{}
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"key":    s.weatherAPI,
			"days":   "1",
			"aqi":    "no",
			"alerts": "no",
			"q":      "Kamenec Pod Vtacnikom",
		}).
		SetHeader("Accept", "application/json").
		SetResult(weather).
		Get("http://api.weatherapi.com/v1/forecast.json")
	if err != nil || resp.IsError() {
		return nil, err
	}
	return &types.Weather{
		Temperature:   weather.Forecast.Forecastday[0].Day.AvgTempC,
		ConditionIcon: weather.Forecast.Forecastday[0].Day.Condition.Icon,
		ConditionText: weather.Forecast.Forecastday[0].Day.Condition.Text,
		Precipitation: weather.Forecast.Forecastday[0].Day.DailyChanceOfRain,
	}, nil
}
