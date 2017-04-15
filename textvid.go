package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/keitax/textvid/application"
	"github.com/keitax/textvid/config"
)

func main() {
	c, err := config.Parse("./config.yml")
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	app := application.New(c)

	logrus.Info("Listen on 8080.")
	if err := http.ListenAndServe(":8080", app); err != nil {
		logrus.Fatal(err)
	}
}
