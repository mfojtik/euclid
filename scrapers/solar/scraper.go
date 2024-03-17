package solar

import (
	"encoding/json"
	"github.com/mfojtik/euclid/scrapers/types"
	"io"
	"net/http"
	"net/url"
)

type Scraper struct {
	BaseURL *url.URL
}

func New(baseUrl string) *Scraper {
	parsedUrl, _ := url.Parse(baseUrl)
	return &Scraper{
		BaseURL: parsedUrl,
	}
}

func (s *Scraper) httpGet(u *url.URL) ([]byte, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", u.String()+"/inverter.json", nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func (s *Scraper) Scrape() (*types.Solar, error) {
	jsonBytes, err := s.httpGet(s.BaseURL)
	if err != nil {
		return nil, err
	}
	data := types.Solar{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
