package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/receipt"
	"golang.org/x/sync/errgroup"
)

const (
	topicGB = "homin-dev/gb" // guest book
	topicHB = "homin-dev/hb" // heart beat
)

var (
	programName                   = "gb-noti"
	buildStamp, gitHash, buildTag string

	flagRpType         string
	flagRpDevPath      string
	flagTestRunPrinter bool

	rp *receipt.Printer
)

func main() {
	log.Printf("%s-%s-%s(%s)", programName, buildTag, gitHash, buildStamp)
	defer log.Printf("%s finished", programName)

	flag.StringVar(&flagRpType, "t", "serial", "receipt printer type [serial|usb]")
	flag.StringVar(&flagRpDevPath, "d", "/dev/ttyACM0", "receipt printer dev path")
	flag.BoolVar(&flagTestRunPrinter, "tp", false, "test run printer")
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

	if flagTestRunPrinter {
		jsonStr := `{"from":"","msg":"우리의 소원은 통일\r\n꿈에도 소원은 통일\r\n이 정성 다해서 통일\r\n통일을 이루자\r\n\r\n이 겨레 살리는 통일\r\n이 나라 살리는 통일\r\n통일이여 어서오라\r\n통일이여 오라","remoteAddr":"10.128.0.7:42213","timestamp":"2022-09-09T13:40:12Z"}`
		var c chat
		if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
			log.Fatal(errors.Wrap(err, "fail to print"))
		}
		printToReceipt(&c)
		return
	}

	subF := func() error {
		mqttC, err := connectBrokerByWSS(&confSub)
		if err != nil {
			return errors.Wrap(err, "fail to sub")
		}
		log.Println("sub: connected with MQTT broker")
		mqttC.Subscribe(topicGB, 1,
			func(mqttClient mqtt.Client, msg mqtt.Message) {
				log.Printf("got %v from %s", string(msg.Payload()), msg.Topic())

				var c chat
				if err := json.Unmarshal(msg.Payload(), &c); err != nil {
					log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
				}
				if err := printToReceipt(&c); err != nil {
					log.Printf("err: %v", errors.Wrap(err, "fail in sub"))
				}
			},
		)
		tk := time.NewTicker(60 * time.Second)
		defer tk.Stop()
		for range tk.C {
			if !mqttC.IsConnectionOpen() {
				return errors.Wrap(err, "mqtt sub conn lost")
			}
		}
		return nil
	}

	pubF := func() error {
		mqttC, err := connectBrokerByWSS(&confPub)
		if err != nil {
			return errors.Wrap(err, "fail to pub")
		}
		log.Println("pub: connected with MQTT broker")
		tk := time.NewTicker(60 * time.Second)
		defer tk.Stop()
		for range tk.C {
			mqttC.Publish(topicHB, 0, false, "gb-noti")
		}
		return nil
	}

	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(pubF)
	eg.Go(subF)
	err := eg.Wait()
	log.Printf("eg stop: %v", err)
}
