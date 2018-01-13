package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/keitax/tv-fe"
)

func main() {
	c := tvfe.GetFromEnv()
	if c.RunLevel == tvfe.DevelopmentRunLevel {
		logrus.SetLevel(logrus.DebugLevel)
	}

	app, err := tvfe.NewApplication(c)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	logrus.Infof("Listen on %s.\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), app); err != nil {
		logrus.Fatal(err)
	}
}
