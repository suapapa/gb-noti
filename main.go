package main

import (
	"encoding/json"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	topic = "homin-dev/gb"
)

var (
	programName                   = "gb-noti"
	buildStamp, gitHash, buildTag string
)

func main() {
	log.Printf("%s-%s-%s(%s)", programName, buildTag, gitHash, buildStamp)
	defer log.Println("program-name finished")

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
	subC.Subscribe(topic, 1,
		func(subClient mqtt.Client, msg mqtt.Message) {
			log.Printf("got %v from %s", string(msg.Payload()), msg.Topic())

			var c chat
			if err := json.Unmarshal(msg.Payload(), &c); err != nil {
				log.Fatal(errors.Wrap(err, "fail to print"))
			}
			// TODO: print from here!
			// log.Println(string(msg.Payload()))
			print(&c)
		},
	)

	ch := make(chan any)
	<-ch
}
