package scrape

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Temperature struct {
	Timestamp            time.Time
	Outdoor              float64
	BivalenceActive      bool
	WaterCurrent         float64
	WaterDesired         float64
	DownstairsCurrent    float64
	DownstairsDesired    float64
	DownstairsPumpActive bool
	UpstairsCurrent      float64
	UpstairsDesired      float64
	UpstairsPumpActive   bool
}

type Scraper struct {
	BaseURL        *url.URL
	user, password string
}

func New(baseUrl, user, pass string) *Scraper {
	parsedUrl, _ := url.Parse(baseUrl)
	return &Scraper{
		BaseURL:  parsedUrl,
		user:     user,
		password: pass,
	}
}

func (s *Scraper) httpGet(u *url.URL) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.user+":"+s.password)))
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

type XMLResult struct {
	XMLName xml.Name   `xml:"PAGE"`
	Inputs  []XMLInput `xml:"INPUT"`
}

type XMLInput struct {
	XMLName xml.Name `xml:"INPUT"`
	Name    string   `xml:"NAME,attr"`
	Value   string   `xml:"VALUE,attr"`
}

func getValueFor(name string, inputs []XMLInput) string {
	for i := range inputs {
		if inputs[i].Name == name {
			return strings.TrimSpace(inputs[i].Value)
		}
	}
	return ""
}

func (s *Scraper) getLatestTemperature(tempDataURL *url.URL) (*Temperature, error) {
	tempDataBytes, err := s.httpGet(tempDataURL)
	if err != nil {
		return nil, err
	}

	var result XMLResult

	reader := bytes.NewReader(tempDataBytes)

	dec := xml.NewDecoder(reader)
	dec.CharsetReader = charset.NewReaderLabel

	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	t := Temperature{
		Timestamp:         time.Now(),
		Outdoor:           stringToNum(getValueFor("__T033A2538_REAL_.1f", result.Inputs)),
		WaterCurrent:      stringToNum(getValueFor("__T50A32455_REAL_.1f", result.Inputs)),
		WaterDesired:      stringToNum(getValueFor("__T20B5E623_REAL_.1f", result.Inputs)),
		DownstairsCurrent: stringToNum(getValueFor("__T46AA2571_REAL_.1f", result.Inputs)),
		DownstairsDesired: stringToNum(getValueFor("__T05D9E707_REAL_.1f", result.Inputs)),
		UpstairsCurrent:   stringToNum(getValueFor("__T7CB1261D_REAL_.1f", result.Inputs)),
		UpstairsDesired:   stringToNum(getValueFor("__T6A6DE46B_REAL_.1f", result.Inputs)),
	}

	if v := getValueFor("__T6EE9225F_BOOL_i", result.Inputs); v != "0" {
		t.DownstairsPumpActive = true
	}
	if v := getValueFor("__T911354AE_BOOL_i", result.Inputs); v != "0" {
		t.UpstairsPumpActive = true
	}

	return &t, nil
}

func stringToNum(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (s *Scraper) Scrape() (*Temperature, error) {
	u, _ := url.Parse(s.BaseURL.String() + "/PAGE115.XML")
	return s.getLatestTemperature(u)
}
