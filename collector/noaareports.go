package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

func getNOAAStormReports() {
	// calculate and store yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1)
	baseURL := "https://www.spc.noaa.gov/climo/reports/"

	// create kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Printf("error creating Kafka producer: %v", err)
		return
	}
	defer producer.Close()

	// get tornado report from noaa
	tornadoURL := baseURL + yesterday.Format("060102") + "_rpts_filtered_torn.csv"
	tornadoesRaw, err := retrieveAndParseCSV(tornadoURL)
	if err != nil {
		fmt.Printf("Error getting tornado reports:%s from %s", err, tornadoURL)
	} else if tornadoesRaw != nil {
		// parse each record and push into kafka topic
		for _, record := range tornadoesRaw[1:] {
			tornado := parseTornadoReport(record, yesterday)
			err = putStormReportInKafka(producer, tornado)
			if err != nil {
				fmt.Printf("error putting tornado report in kafka: %s", err)
			}
		}
	}

	// get hail report from noaa
	hailURL := baseURL + yesterday.Format("060102") + "_rpts_filtered_hail.csv"
	hailRaw, err := retrieveAndParseCSV(hailURL)
	if err != nil {
		fmt.Printf("Error getting hail reports:%s from %s", err, hailURL)
	} else if hailRaw != nil {
		// parse each record and push into kafka topic
		for _, record := range hailRaw[1:] {
			hail := parseHailReport(record, yesterday)
			err = putStormReportInKafka(producer, hail)
			if err != nil {
				fmt.Printf("error putting tornado report in kafka: %s", err)
			}
		}
	}

	// get wind report from noaa
	windURL := baseURL + yesterday.Format("060102") + "_rpts_filtered_wind.csv"
	windRaw, err := retrieveAndParseCSV(windURL)
	if err != nil {
		fmt.Printf("Error getting wind reports:%s from %s", err, windRaw)
	} else if windRaw != nil {
		// parse each record and push into kafka topic
		for _, record := range windRaw[1:] {
			wind := parseWindReport(record, yesterday)
			err = putStormReportInKafka(producer, wind)
			if err != nil {
				fmt.Printf("error putting tornado report in kafka: %s", err)
			}
		}
	}
}

// StormReport is a super set of all the fields from different reports with
// additional data to track the type of storm and original report
type StormReport struct {
	Source     string `json:"source,omitempty"`
	ReportDate string `json:"reportDate,omitempty"`
	StormType  string `json:"stormType,omitempty"`
	Time       string `json:"time,omitempty"`
	Size       string `json:"size,omitempty"`
	Speed      string `json:"speed,omitempty"`
	FScale     string `json:"FScale,omitempty"`
	Location   string `json:"location,omitempty"`
	County     string `json:"county,omitempty"`
	State      string `json:"state,omitempty"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
	Comments   string `json:"comments,omitempty"`
}

func parseTornadoReport(record []string, reportDate time.Time) *StormReport {
	report := &StormReport{
		ReportDate: reportDate.Format("2006-01-02"),
		StormType:  "Tornado",
		Time:       record[0],
		FScale:     record[1],
		Location:   record[2],
		County:     record[3],
		State:      record[4],
		Latitude:   record[5],
		Longitude:  record[6],
		Comments:   record[7],
	}
	return report
}

func parseHailReport(record []string, reportDate time.Time) *StormReport {
	report := &StormReport{
		ReportDate: reportDate.Format("2006-01-02"),
		StormType:  "Hail",
		Time:       record[0],
		Size:       record[1],
		Location:   record[2],
		County:     record[3],
		State:      record[4],
		Latitude:   record[5],
		Longitude:  record[6],
		Comments:   record[7],
	}
	return report
}

func parseWindReport(record []string, reportDate time.Time) *StormReport {
	report := &StormReport{
		ReportDate: reportDate.Format("2006-01-02"),
		StormType:  "Wind",
		Time:       record[0],
		Speed:      record[1],
		Location:   record[2],
		County:     record[3],
		State:      record[4],
		Latitude:   record[5],
		Longitude:  record[6],
		Comments:   record[7],
	}
	return report
}

func putStormReportInKafka(producer sarama.SyncProducer, report *StormReport) error {
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("error marshaling report to JSON: %v", err)
	}

	message := &sarama.ProducerMessage{
		Topic: "raw-weather-reports",
		Value: sarama.ByteEncoder(reportBytes),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("error sending message to Kafka: %v", err)
	}

	return nil
}
