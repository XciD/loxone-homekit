package main

import (
	"loxone-homekit/components"
	"os"
	"strings"

	"github.com/brutella/hc/accessory"

	"github.com/XciD/loxone-ws"
	"github.com/XciD/loxone-ws/events"
	"github.com/brutella/hc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(viper.GetString("log"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
}

func main() {
	// Parse arguments
	var configComponents []components.ComponentConfig
	err := viper.UnmarshalKey("components", &configComponents)
	if err != nil {
		log.Fatal(err)
	}

	// Open socket
	lox, err := loxone.New(
		viper.GetString("loxone.host"),
		viper.GetString("loxone.user"),
		viper.GetString("loxone.password"),
	)

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

	for _, config := range configComponents {
		if control, ok := loxoneConfig.Controls[config.ID]; ok {
			component := components.CreateComponent(config, control, lox)
			if component != nil {
				accessories = append(accessories, component.Accessory)
			}
		}
	}

	info := accessory.Info{
		Name:         "Loxone",
		Manufacturer: "Loxone",
	}

	bridge := accessory.NewBridge(info)

	t, err := hc.NewIPTransport(hc.Config{Pin: viper.GetString("pin")}, bridge.Accessory, accessories...)

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
