package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/receipt"
	"github.com/suapapa/site-gb/msg"
	"golang.org/x/sync/errgroup"
)

const (
	topicGB = "homin-dev/gb" // guest book
	topicHB = "homin-dev/hb" // heart beat
)

var (
	programName                   = "gb-noti"
	buildStamp, gitHash, buildTag string

	flagRpType    string
	flagRpDevPath string
	flagFontPath  string
	flagHQ        bool

	lastPork    = time.Now()
	maxWaitPork = 35 * time.Minute

	rp *receipt.Printer
)

func main() {
	log.Printf("%s-%s-%s(%s)", programName, buildTag, gitHash, buildStamp)
	defer log.Printf("%s finished", programName)

	flag.StringVar(&flagRpType, "t", "serial", "receipt printer type [serial|usb]")
	flag.StringVar(&flagRpDevPath, "d", "/dev/ttyACM0", "receipt printer dev path")
	flag.StringVar(&flagFontPath, "f", "", "external font path to use")
	flag.BoolVar(&flagHQ, "q", false, "better quality printing")
	flag.Parse()

	switch flagRpType {
	case "serial":
		rp = receipt.NewSerialPrinter(flagRpDevPath, 9600)
	case "usb":
		rp = receipt.NewUSBPrinter(flagRpDevPath)
	default:
		log.Println("select serial or usb")
		os.Exit(-1)

	}

	conf := Config{
		Host:     "homin.dev",
		Port:     9001,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
	}
	confSub, confPub := conf, conf
	confSub.ClientID = "suapapa-gb-noti-sub"
	confPub.ClientID = "suapapa-gb-noti-pub"

	subF := func() error {
		mqttC, err := connectBrokerByWSS(&confSub)
		if err != nil {
			return errors.Wrap(err, "fail to sub")
		}
		defer mqttC.Disconnect(1000)
		log.Println("sub: connected with MQTT broker")
		mqttC.Subscribe(topicGB, 1,
			func(mqttClient mqtt.Client, mqttMsg mqtt.Message) {
				topic, payload := mqttMsg.Topic(), mqttMsg.Payload()
				log.Printf("got %v from %s", string(payload), topic)

				switch topic {
				case "homin-dev/gb":
					var m msg.Message
					if err := json.Unmarshal(mqttMsg.Payload(), &m); err != nil {
						log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
					}

					// dont print if it is just a pork
					if m.Type == msg.MTGuestBook {
						if gb, ok := m.Data.(msg.GuestBook); ok {
							if err := printToReceipt(&gb); err != nil {
								log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
							}
						} else {
							log.Printf("err: fail to convert msg.Data to GuestBook")
						}
					}
					lastPork = time.Now()

				default:
					log.Printf("err: unknown topic %s", topic)
				}
			},
		)
		tk := time.NewTicker(10 * time.Second)
		defer tk.Stop()
		for range tk.C {
			if !mqttC.IsConnected() || !mqttC.IsConnectionOpen() {
				return errors.Wrap(err, "mqtt sub conn lost")
			}
		}
		return nil
	}

	/*
		pubF := func() error {
			mqttC, err := connectBrokerByWSS(&confPub)
			if err != nil {
				return errors.Wrap(err, "fail to pub")
			}
			defer mqttC.Disconnect(1000)
			log.Println("pub: connected with MQTT broker")
			tk := time.NewTicker(60 * time.Second)
			defer tk.Stop()
			for range tk.C {
				mqttC.Publish(topicHB, 0, false, "gb-noti")
			}
			return nil
		}
	*/

	eg, _ := errgroup.WithContext(context.Background())
	// eg.Go(pubF)
	eg.Go(subF)
	eg.Go(checkPork)
	err := eg.Wait()
	log.Printf("eg stop: %v", err)
}

func checkPork() error {
	tk := time.NewTicker(5 * time.Second)
	defer tk.Stop()
	for range tk.C {
		if time.Since(lastPork) > maxWaitPork {
			return fmt.Errorf("no porking over %v", maxWaitPork)
		}
	}
	return nil
}
