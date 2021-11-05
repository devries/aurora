package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/pflag"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricsAddress = pflag.StringP("metrics", "m", ":8080", "metrics server address in format ip:port")
var version = "DEV"

type SpaceWeatherConditions struct {
	DateStamp     string `json:"DateStamp"`
	TimeStamp     string `json:"TimeStamp"`
	RadioBlackout struct {
		Scale     *string `json:"Scale"`
		Text      *string `json:"Text"`
		MinorProb *string `json:"MinorProb"`
		MajorProb *string `json:"MajorProb"`
	} `json:"R"`
	SolarRadiation struct {
		Scale *string `json:"Scale"`
		Text  *string `json:"Text"`
		Prob  *string `json:"Prob"`
	} `json:"S"`
	Geomagnetic struct {
		Scale *string `json:"Scale"`
		Text  *string `json:"Text"`
	} `json:"G"`
}

type SpaceWeatherResponse map[string]SpaceWeatherConditions

func main() {
	pflag.Parse()
	collector := newGeomagneticCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Version: %s", version)
	log.Printf("Starting to serve on %s", *metricsAddress)
	log.Fatal(http.ListenAndServe(*metricsAddress, nil))
}

func getSpaceWeather() (SpaceWeatherResponse, error) {
	url := "https://services.swpc.noaa.gov/products/noaa-scales.json"

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "aurora-tracker")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var swr SpaceWeatherResponse

	if err := dec.Decode(&swr); err != nil {
		return nil, err
	}

	return swr, nil
}

type geomagneticCollector struct {
	geoMetric *prometheus.Desc
}

func newGeomagneticCollector() *geomagneticCollector {
	return &geomagneticCollector{
		geoMetric: prometheus.NewDesc("aurora_geomagnetic_storm",
			"Geomagnetic storm index.",
			[]string{"time"}, nil),
	}
}

func (collector *geomagneticCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.geoMetric
}

func (collector *geomagneticCollector) Collect(ch chan<- prometheus.Metric) {
	var scale float64

	swr, err := getSpaceWeather()
	if err != nil {
		log.Printf("Error getting space weather: %s", err)
		return
	}

	for _, entry := range []string{"0", "1", "2", "3"} {
		current := swr[entry]

		scaleText := current.Geomagnetic.Scale
		scale, err = strconv.ParseFloat(*scaleText, 64)
		if err != nil {
			log.Printf("Failed to convert to float: %s", err)
			continue
		}

		switch entry {
		case "0":
			ch <- prometheus.MustNewConstMetric(collector.geoMetric, prometheus.GaugeValue, scale, "current")
		case "1":
			ch <- prometheus.MustNewConstMetric(collector.geoMetric, prometheus.GaugeValue, scale, "predicted")
		case "2":
			ch <- prometheus.MustNewConstMetric(collector.geoMetric, prometheus.GaugeValue, scale, "two_day")
		case "3":
			ch <- prometheus.MustNewConstMetric(collector.geoMetric, prometheus.GaugeValue, scale, "three_day")
		}
	}
}
