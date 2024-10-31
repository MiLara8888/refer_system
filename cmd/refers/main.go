package main

import (
	"errors"
	"log"
	"os"
	"refers_rest/pkg/settings"
	"runtime"

	"github.com/urfave/cli"
	"refers_rest/internal/refers_rest"
)

var (
	config *settings.Config
)

func init() {

	var err error
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	config, err = settings.InitEnv()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	app := &cli.App{
		Name:  "Logger-monitoring rest server",
		Usage: "",
		Commands: []cli.Command{
			{
				Name:  "serve",
				Usage: "Start test http",
				Action: func(*cli.Context) error {
					return Run()
				},
			},
		},
	}
	app.Run(os.Args)
}

func Run() error {
	r, err := refersrest.New(config)
	if err != nil {
		log.Fatal(err)
	}
	err = r.Start()
	if err != nil {
		log.Fatal(errors.Join(errors.New("err on start service"), err))
	}
	return nil
}
