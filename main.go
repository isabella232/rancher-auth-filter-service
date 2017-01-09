package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
	"github.com/rancher/rancher-auth-filter-service/service"
	"github.com/urfave/cli"
)

var VERSION = "v0.1.0-dev"

func main() {
	logrus.Infof("Starting authantication filtering Service")
	//init parsing command line
	app := cli.NewApp()
	app.Name = "rancher-auth-filter-service"
	app.Version = "v0.1.0-dev"
	app.Usage = "Rancher authantication Filter Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rancherUrl",
			Value: "http://54.255.182.226:8080/",
			Usage: "Rancher server url",
		},
		cli.StringFlag{
			Name:  "localport",
			Value: "8080",
			Usage: "Local server port ",
		},
	}

	app.Action = func(c *cli.Context) error {
		manager.URL = c.String("rancherUrl")
		manager.Port = c.String("localport")
		logrus.Infof("Rancher server URL:" + manager.URL + " The validation filter server running on local port:" + manager.Port)
		//create mux router
		router := service.NewRouter()
		http.Handle("/", router)
		serverString := ":" + manager.Port
		//start local server
		logrus.Fatal(http.ListenAndServe(serverString, nil))
		return nil
	}

	app.Run(os.Args)

}
