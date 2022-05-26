package main

import (
    "log"
    "os"
    "errors"
    "net/http"
    "path/filepath"

    "github.com/miralgj/si/pkg/config"
    "github.com/miralgj/si/pkg/router"

    "github.com/urfave/cli/v2"
)

func main() {
    app := cli.NewApp()
    app.Name = "Si"
    app.Version = "v0.1.0"
    app.Usage = "Expose commands as an API"
    app.Action = cliActionHandler
    app.Before = cliBeforeHandler
    app.Flags = config.GetFlags()
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

func cliActionHandler(c *cli.Context) error {
    r := router.NewRouter()
    http.ListenAndServe(config.Config.Listen+":"+config.Config.Port, r)
    return nil
}

func cliBeforeHandler(c *cli.Context) error {
    if len(c.StringSlice("command")) > 0 {
        // StringSliceFlag doesn't support Destination
        for _, path := range c.StringSlice("command") {
            // Verify file exists
            if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
                log.Fatal(err.Error())
            }
            name := filepath.Base(path)
            config.Config.Commands[name] = path
        }
    }
    return nil
}
