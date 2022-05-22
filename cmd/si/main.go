package main

import (
    //"fmt"
    "log"
    "os"
    "net/http"
    "github.com/miralgj/si/pkg/config"
    "github.com/miralgj/si/pkg/router"
    "github.com/urfave/cli/v2"
)

func main() {
    app := cli.NewApp()
    app.Name = "Si"
    app.Version = "v0.0.1"
    app.Usage = "Expose system commands as an API"
    app.Action = cliHandler
    app.Flags = config.Flags
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

func cliHandler(c *cli.Context) error {
    r := router.New()
    conf := config.New()
    http.ListenAndServe(conf.Listen+":"+conf.Port, r)
    return nil
}
