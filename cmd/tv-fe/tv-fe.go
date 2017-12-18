package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/keitax/tv-fe"
)

func main() {
	c, err := tvfe.Parse("./config.toml")
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	if c.RunLevel == tvfe.DevelopmentRunLevel {
		logrus.SetLevel(logrus.DebugLevel)
	}

	app, err := tvfe.NewApplication(c)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	logrus.Info("Listen on 8080.")
	if err := http.ListenAndServe(":8080", app); err != nil {
		logrus.Fatal(err)
	}
}
