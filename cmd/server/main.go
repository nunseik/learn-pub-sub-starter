package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"	
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"


	
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionAddress := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionAddress)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Couldn't create channel: %v", err)
	}

	err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect),string(routing.PauseKey), routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("Couldn't publish JSON: %v", err)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nRabbitMQ connection closed.")
}
