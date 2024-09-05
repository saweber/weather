package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/uptrace/bun"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready    chan bool
	r        *Repository
	producer sarama.SyncProducer
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			c.processRecord(message.Value)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

// Example processor function that will handle the data
func (c *Consumer) processRecord(data []byte) {
	fmt.Printf("Received message: %s", string(data))

	raw := RawStormReport{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		fmt.Printf("error unmarshalling raw storm report: %v", err)
		return
	}

	sr := c.transformRawStormReport(&raw)

	err = c.r.InsertRecordIntoPostgres(sr)
	if err != nil {
		fmt.Printf("error inserting record into postgres: %v", err)
		return
	}

	err = c.insertRecordIntoKafka(sr)
	if err != nil {
		fmt.Printf("error inserting record into kafka: %v", err)
		return
	}
}

type RawStormReport struct {
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

// func to accept a record and transform to processed record
func (c *Consumer) transformRawStormReport(raw *RawStormReport) StormReport {
	// there are other transformations that can be done here, as well as
	// validation, but omitted due to time constraints
	return StormReport{
		ReportDate: raw.ReportDate,
		StormType:  raw.StormType,
		Latitude:   raw.Latitude,
		Longitude:  raw.Longitude,
		Location:   raw.Location,
		County:     raw.County,
		State:      raw.State,
		Comments:   raw.Comments,
		Speed:      convertStrToInt(raw.Speed),
		Size:       convertStrToInt(raw.Size),
		FScale:     convertStrToInt(raw.FScale),
		Time:       c.getStormReportTime(raw.ReportDate, raw.Time),
	}
}

func convertStrToInt(raw string) int {
	if raw != "" {
		size, err := strconv.Atoi(raw)
		if err != nil {
			if raw != "UNK" {
				fmt.Printf("failed to convert %s to integer", raw)
			}
		} else {
			return size
		}
	}
	return 0
}

func (c *Consumer) getStormReportTime(rawDate, rawTime string) time.Time {
	format := "2006-01-02 1504"
	timeStr := rawDate + " " + rawTime
	stormTime, err := time.Parse(format, timeStr)
	if err != nil {
		// for now, just return default time if unable to parse
		return time.Time{}
	}

	// NOAA reports ase not midnight to midnight, but noon to noon
	if stormTime.Hour() < 12 {
		stormTime = stormTime.AddDate(0, 0, -1)
	}
	return stormTime
}

type StormReport struct {
	bun.BaseModel `bun:"table:storm_reports,alias:sr"`
	ReportDate    string    `json:"reportDate,omitempty" bun:"report_date"`
	StormType     string    `json:"stormType,omitempty" bun:"storm_type"`
	Latitude      string    `json:"latitude,omitempty"`
	Longitude     string    `json:"longitude,omitempty"`
	Location      string    `json:"location,omitempty"`
	County        string    `json:"county,omitempty"`
	State         string    `json:"state,omitempty"`
	Comments      string    `json:"comments,omitempty"`
	Speed         int       `json:"speed,omitempty"`
	Size          int       `json:"size,omitempty"`
	FScale        int       `json:"fScale,omitempty" bun:"f_scale"`
	Time          time.Time `json:"time,omitempty" bun:"time"`
}

// func to insert into processed kafka queue
func (c *Consumer) insertRecordIntoKafka(report StormReport) error {
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("error marshaling report to JSON: %v", err)
	}

	message := &sarama.ProducerMessage{
		Topic: "transformed-weather-data",
		Value: sarama.ByteEncoder(reportBytes),
	}

	_, _, err = c.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("error sending message to Kafka: %v", err)
	}

	return nil
}
