package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	routingKey := fmt.Sprintf("%v.*", routing.GameLogSlug)
	_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilTopic, "game_logs", routingKey, pubsub.Durable)
	if err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) != 0 {
			if words[0] == "pause" {
				fmt.Println("sending a pause message")
				err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{
					IsPaused: true})
				if err != nil {
					log.Fatalf("Couldn't publish JSON: %v", err)
				}
			} else if words[0] == "resume" {
				fmt.Println("sending a resume message")
				err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{
					IsPaused: false})
				if err != nil {
					log.Fatalf("Couldn't publish JSON: %v", err)
				}
			} else if words[0] == "quit" {
				fmt.Println("exiting")
				break
			} else {
				fmt.Println("Can't understand the command")
			}
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nRabbitMQ connection closed.")
}
