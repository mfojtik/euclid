package solar

import (
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/comms/tcpip"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/devices/sofar"
	"github.com/mfojtik/euclid/scrapers/types"
	"log"
	"time"
)

type Scraper struct {
	device *sofar.Logger
}

func New(baseUrl string, serial int64) *Scraper {
	return &Scraper{device: sofar.NewSofarLogger(uint(serial), tcpip.New(baseUrl), []string{
		"PV_Generation_Today",
		"PV_Generation_Total",
		"Temperature_HeatSink1",
		"ActivePower_Output_Total",
		"Load_Consumption_Today",
	}, []string{})}
}

func (s *Scraper) Scrape() (*types.Solar, error) {
	data, err := s.device.Query()
	if err != nil {
		log.Printf("ERR: %v", err)
		return &types.Solar{
			Timestamp: time.Now().Unix(),
			Status:    "off",
		}, nil
	}
	return &types.Solar{
		Timestamp:        time.Now().Unix(),
		Status:           "on",
		GenerationNow:    (float32(data["ActivePower_Output_Total"].(int16)) * 10) / 1000,
		ConsumptionToday: (float32(data["Load_Consumption_Today"].(uint32)) * 10) / 1000,
		GenerationTotal:  (float32(data["PV_Generation_Total"].(uint32)) * 100) / 1000, // 10*W to kWh
		GenerationToday:  (float32(data["PV_Generation_Today"].(uint32)) * 10) / 1000,
	}, nil
}
