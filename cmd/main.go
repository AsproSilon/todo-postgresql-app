package main

import (
	"aspro/internal/app"
	"log"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	log.Print("config initializing")
	if err := app.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	a, err := app.NewApp()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("Running Application")
	a.Run()
}
