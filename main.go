package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/SmartMeshFoundation/matrix-regservice/models"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/matrix-regservice/rest"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/matrix-regservice/internal/debug"

	"github.com/SmartMeshFoundation/matrix-regservice/params"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	fmt.Printf("args=%q\n", os.Args)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "astoken",
			Usage: "acess  token for matrix HomeServer",
			Value: params.ASToken,
		},
		cli.StringFlag{
			Name:  "hstoken",
			Usage: "homeserver access token for my service",
			Value: params.HSToken,
		},
		cli.StringFlag{
			Name:  "matrixurl",
			Usage: "creat user url of matrix",
			Value: params.MatrixRegisterUrl,
		},
		cli.StringFlag{
			Name:  "host",
			Usage: "listen host",
			Value: params.APIHost,
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "listen port",
			Value: params.APIPort,
		},
		cli.StringFlag{
			Name:  "datapath",
			Usage: "database file directory",
			Value: ".matrix",
		},
		cli.StringFlag{
			Name:  "matrixdomain",
			Usage: "matrix server's domain for this register service",
			Value: params.MatrixDomain,
		},
	}
	app.Flags = append(app.Flags, debug.Flags...)
	app.Action = mainCtx
	app.Name = "matrix registor service"
	app.Version = "0.1"
	app.Before = func(ctx *cli.Context) error {
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		return nil
	}
	app.Commands = []cli.Command{
		{
			Action:      genyaml,
			Name:        "genconfig",
			Usage:       "generate config yaml file for matrix and run.sh for service run",
			Description: "generate config yaml file for matrix and run.sh for service run",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Error(fmt.Sprintf("run err %s", err))
	}
}

func mainCtx(ctx *cli.Context) {
	log.Info(fmt.Sprintf("welcome matrix application service %s\n", ctx.App.Version))
	err := config(ctx)
	if err != nil {
		panic(err)
	}
	err = setupdb(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		quitSignal := make(chan os.Signal, 1)
		signal.Notify(quitSignal, os.Interrupt, os.Kill)
		<-quitSignal
		signal.Stop(quitSignal)
		models.CloseDB()
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()
	rest.Start()
}

func config(ctx *cli.Context) error {
	params.MatrixRegisterUrl = ctx.GlobalString("matrixurl")
	params.ASToken = ctx.GlobalString("astoken")
	params.HSToken = ctx.GlobalString("hstoken")
	params.APIHost = ctx.GlobalString("host")
	params.APIPort = ctx.GlobalInt("port")
	params.MatrixDomain = ctx.GlobalString("matrixdomain")
	return nil
}

func setupdb(ctx *cli.Context) error {
	datadir := params.DBPath
	if !utils.Exists(datadir) {
		err := os.MkdirAll(datadir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", datadir, err)
			return err
		}
	}
	userDbPath := filepath.Join(datadir, "register.db")
	models.SetUpDB(userDbPath)
	return nil
}

type hsConfig struct {
	/*
			id: regapp_transport01
		hs_token: 6350c08ee06eed113f84da06c2f85369dcf0777d79d9679efd3ad2abdadd01d9
		as_token: 8309e0e83ea8ddfb41ed860eb76be627256d4651ccd2c313597d6f286f41bf82
		namespaces:
		  users:
		    - exclusive: false
		      regex: '@.*'
		  aliases: []
		  rooms: []
		url: 'http://localhost:8009/regapp/1'
		sender_localpart: app.transport01.Photon.network
		rate_limited: true
		protocols:
		  - app.transport01.Photon.network
	*/
	ID              string   `yaml:"id"`
	HSToken         string   `yaml:"hs_token"`
	ASToken         string   `yaml:"as_token"`
	URL             string   `yaml:"url"`
	SenderLocalPart string   `yaml:"sender_localpart"`
	Protocols       []string `yaml:"protocols"`
	NameSpaces      struct {
		Users   []interface{} `yaml:"users"`
		Aliases []string
		Rooms   []string
	} `yaml:"namespaces"`
}
type exclusive struct {
	Exclusive bool   `yaml:"exclusive"`
	RegEx     string `yaml:"regex"`
}

func genyaml(ctx *cli.Context) error {
	err := config(ctx)
	if err != nil {
		return err
	}
	c := &hsConfig{
		ID:              fmt.Sprintf("%s-%s", utils.RandomString(10), params.MatrixDomain),
		HSToken:         utils.RandomString(40),
		ASToken:         utils.RandomString(40),
		URL:             fmt.Sprintf("http://%s:%d/regapp/1", params.APIHost, params.APIPort),
		SenderLocalPart: params.MatrixDomain,
		Protocols:       []string{fmt.Sprintf("regapp.%s", params.MatrixDomain)},
	}
	e := exclusive{Exclusive: false, RegEx: "@.*"}
	users := []interface{}{e}
	c.NameSpaces.Users = users
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	log.Trace(fmt.Sprintf("out=%s", string(out)))
	err = ioutil.WriteFile("registration.yaml", out, 0600)
	if err != nil {
		return err
	}
	runStr := fmt.Sprintf("#!/bin/sh\n%s --astoken %s --hstoken %s --matrixurl %s --host %s --port %d --datapath %s --matrixdomain %s --verbosity 5",
		os.Args[0], c.ASToken, c.HSToken, params.MatrixRegisterUrl, params.APIHost, params.APIPort, params.DBPath, params.MatrixDomain,
	)
	return ioutil.WriteFile("run.sh", []byte(runStr), os.ModePerm)
}
