package main

import (
    //"fmt"
    "log"
    "os"
    "path/filepath"
    "net/http"

    "github.com/miralgj/si/pkg/config"
    "github.com/miralgj/si/pkg/router"

    "github.com/urfave/cli/v2"
)

func main() {
    app := cli.NewApp()
    app.Name = "Si"
    app.Version = "v0.0.1"
    app.Usage = "Expose commands as an API"
    app.Action = cliAction
    app.Before = func(c *cli.Context) error {
        if len(c.StringSlice("command")) > 0 {
            // StringSliceFlag doesn't support Destination
            for _, path := range c.StringSlice("command") {
                name := filepath.Base(path)
                config.Config.Commands[name] = path
            }
        }
        return nil
    }
    app.Flags = config.GetFlags()
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

func cliAction(c *cli.Context) error {
    r := router.NewRouter()
    http.ListenAndServe(config.Config.Listen+":"+config.Config.Port, r)
    return nil
}
