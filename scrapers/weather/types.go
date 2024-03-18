package weather

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	TimeZoneID     string  `json:"tz_id"`
	LocalTimeEpoch int64   `json:"localtime_epoch"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Day struct {
	MaxTempC          float64   `json:"maxtemp_c"`
	MaxTempF          float64   `json:"maxtemp_f"`
	MinTempC          float64   `json:"mintemp_c"`
	MinTempF          float64   `json:"mintemp_f"`
	AvgTempC          float64   `json:"avgtemp_c"`
	AvgTempF          float64   `json:"avgtemp_f"`
	MaxWindMPH        float64   `json:"maxwind_mph"`
	MaxWindKPH        float64   `json:"maxwind_kph"`
	TotalPrecipMM     float64   `json:"totalprecip_mm"`
	TotalPrecipIN     float64   `json:"totalprecip_in"`
	TotalSnowCM       float64   `json:"totalsnow_cm"`
	AvgVisKM          float64   `json:"avgvis_km"`
	AvgVisMiles       float64   `json:"avgvis_miles"`
	AvgHumidity       int       `json:"avghumidity"`
	DailyWillItRain   int       `json:"daily_will_it_rain"`
	DailyChanceOfRain int       `json:"daily_chance_of_rain"`
	DailyWillItSnow   int       `json:"daily_will_it_snow"`
	DailyChanceOfSnow int       `json:"daily_chance_of_snow"`
	Condition         Condition `json:"condition"`
	UV                float64   `json:"uv"`
}

type Astro struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination int    `json:"moon_illumination"`
	IsMoonUp         int    `json:"is_moon_up"`
	IsSunUp          int    `json:"is_sun_up"`
}

type Forecastday struct {
	Date      string `json:"date"`
	DateEpoch int64  `json:"date_epoch"`
	Day       Day    `json:"day"`
	Astro     Astro  `json:"astro"`
	Hour      []Hour `json:"hour"`
}

type Hour struct {
	// Define hour fields if needed
}

type Current struct {
	LastUpdatedEpoch int64     `json:"last_updated_epoch"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMPH          float64   `json:"wind_mph"`
	WindKPH          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMB       float64   `json:"pressure_mb"`
	PressureIN       float64   `json:"pressure_in"`
	PrecipMM         float64   `json:"precip_mm"`
	PrecipIN         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	VisKM            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	UV               float64   `json:"uv"`
	GustMPH          float64   `json:"gust_mph"`
	GustKPH          float64   `json:"gust_kph"`
}

type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

type WeatherData struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}
