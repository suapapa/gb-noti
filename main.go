package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic = "homin-dev/gb"
)

// TODO
// type GBMsg struct {

// }

func main() {
	log.Println("start gb-noti")
	defer log.Println("gb-noti finished")

	conf := Config{
		Host:     "homin.dev",
		Port:     9001,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
	}

	subC, err := connectBrokerByWSS(&conf)
	if err != nil {
		log.Fatal(err)
	}
	defer subC.Disconnect(1000)
	subC.Subscribe(topic, 0,
		func(subClient mqtt.Client, msg mqtt.Message) {
			log.Printf("got %v from %s", string(msg.Payload()), msg.Topic())
			sendMsgToTelegram(msg.Payload())
		},
	)
	ch := make(chan any)
	<-ch
}
