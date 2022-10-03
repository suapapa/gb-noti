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
					lastPork = time.Now()
					if gb, err := getGBFromMsg(mqttMsg.Payload()); err != nil {
						log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
					} else {
						if err := printToReceipt(gb); err != nil {
							log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
						}
					}

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

func getGBFromMsg(msgBytes []byte) (*msg.GuestBook, error) {
	m := msg.Message{
		Data: &msg.GuestBook{}, // it is needed. if not data will be map[string]any
	}
	if err := json.Unmarshal(msgBytes, &m); err != nil {
		return nil, errors.Wrap(err, "fail to get gb from msg")
	}

	return m.GetGuestBook()
}
