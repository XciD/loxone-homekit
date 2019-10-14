package main

import (
	"flag"
	"loxone-homekit/components"
	"os"

	"github.com/brutella/hc/accessory"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse arguments
	// TODO Config values in viper
	host := flag.String("host", "", "Loxone Host Name")
	user := flag.String("user", "", "Loxone User Name")
	password := flag.String("password", "", "Loxone Password")

	var configComponents []components.Component
	err := viper.UnmarshalKey("configComponents", &configComponents)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	// Open socket
	lox, err := loxone.New(*host, *user, *password)

	if err != nil {
		log.Error(err)
		return
	}

	defer lox.Close()

	// Get config
	loxoneConfig, err := lox.GetConfig()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Get Config OK")

	// Register events
	err = lox.RegisterEvents()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Register Events OK")

	_ = make(map[string]chan *events.Event)

	accessories := make([]*accessory.Accessory, 0)

	// TODO Refactor this
	for _, component := range configComponents {
		if control, ok := loxoneConfig.Controls[component.ID]; ok {
			switch component.LoxoneType {
			case "Switch":
				light := components.NewLoxoneSwitch(component, control, lox)
				accessories = append(accessories, light.Accessory)
			case "Jalousie":
				jalousie := components.NewJalousie(component, control, lox)
				accessories = append(accessories, jalousie.Accessory)
			case "Gate":
				door := components.NewDoor(component, control, lox)
				accessories = append(accessories, door.Accessory)
			case "Intercom":
				intercom := components.NewIntercom(component, control, lox)
				accessories = append(accessories, intercom.Accessory)
			}
		}
	}

	info := accessory.Info{
		Name:         "Loxone",
		Manufacturer: "Loxone",
	}
	bridge := accessory.NewBridge(info)

	// TODO Config pin
	t, err := hc.NewIPTransport(hc.Config{Pin: "32191123"}, bridge.Accessory, accessories...)

	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan bool)

	hc.OnTermination(func() {
		log.Info("Stopping")
		close(stop)
	})

	go t.Start()
	log.Info("Start reading events")
	go lox.PumpEvents(stop)

	<-stop
}
