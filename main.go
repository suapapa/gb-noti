package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/receipt"
)

const (
	topic = "homin-dev/gb"
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

	if flagTestRunPrinter {
		jsonStr := `{"from":"","msg":"우리의 소원은 통일\r\n꿈에도 소원은 통일\r\n이 정성 다해서 통일\r\n통일을 이루자\r\n\r\n이 겨레 살리는 통일\r\n이 나라 살리는 통일\r\n통일이여 어서오라\r\n통일이여 오라","remoteAddr":"10.128.0.7:42213","timestamp":"2022-09-09T13:40:12Z"}`
		var c chat
		if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
			log.Fatal(errors.Wrap(err, "fail to print"))
		}
		printToReceipt(&c)
		return
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
			if err := printToReceipt(&c); err != nil {
				log.Fatal(errors.Wrap(err, "fail to print"))
			}
		},
	)

	tk := time.NewTicker(20 * time.Second)
	defer tk.Stop()
	for t := range tk.C {
		if !subC.IsConnected() {
			log.Fatalf("lost connect with mqtt: %b", t)
		}
	}
}
