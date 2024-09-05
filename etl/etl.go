package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type ETL struct {
	r        *Repository
	producer sarama.SyncProducer
}

// ReadFromKafkaAndProcess is lifted from sarama's example code
func (e *ETL) ReadFromKafkaAndProcess() {
	keepRunning := true
	fmt.Println("Starting a new Sarama consumer")

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		fmt.Printf("Error parsing Kafka version: %v", err)
		panic(err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready:    make(chan bool),
		r:        e.r,
		producer: e.producer,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		fmt.Printf("Error creating consumer group client: %v", err)
		panic(err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				fmt.Printf("Error from consumer: %v", err)
				panic(err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	fmt.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			fmt.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			fmt.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		fmt.Printf("Error closing client: %v", err)
		panic(err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		fmt.Println("Resuming consumption")
	} else {
		client.PauseAll()
		fmt.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
