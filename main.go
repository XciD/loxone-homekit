package main

import (
	"loxone-homekit/components"
	"os"
	"strings"

	"github.com/XciD/loxone-ws"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	"github.com/mdp/qrterminal/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		logrus.Fatal(err)
	}

	level, err := logrus.ParseLevel(viper.GetString("log"))
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})

	log.Debug.Enable()
}

func main() {
	// Open socket
	lox, err := loxone.New(
		viper.GetString("loxone.host"),
		viper.GetInt("loxone.port"),
		viper.GetString("loxone.user"),
		viper.GetString("loxone.password"),
	)

	if err != nil {
		logrus.Fatal(err)
	}

	defer lox.Close()

	// Get config
	loxoneConfig, err := lox.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Get Config OK")

	// Register events
	if err := lox.RegisterEvents(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Register Events OK")

	accessories := components.GetAccessories(lox, loxoneConfig)

	logrus.Infof("Found %d accessories", len(accessories))

	info := accessory.Info{
		Name:         loxoneConfig.MsInfo["projectName"].(string),
		SerialNumber: loxoneConfig.MsInfo["serialNr"].(string),
		Manufacturer: "Loxone",
		Model:        loxoneConfig.MsInfo["msName"].(string),
	}

	bridge := accessory.NewBridge(info)

	t, err := hc.NewIPTransport(hc.Config{
		Pin:         viper.GetString("homekit.pin"),
		Port:        viper.GetString("homekit.port"),
		StoragePath: viper.GetString("homekit.storagePath"),
	}, bridge.Accessory, accessories...)

	if err != nil {
		logrus.Fatal(err)
	}

	uri, _ := t.XHMURI()
	qrterminal.Generate(uri, qrterminal.L, os.Stdout)

	stop := make(chan bool)

	hc.OnTermination(func() {
		logrus.Info("Stopping")
		close(stop)
	})
	go t.Start()
	logrus.Info("Start reading events")
	go lox.PumpEvents(stop)

	<-stop
}
