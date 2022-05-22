package main

import (
    "fmt"
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
    app.Action = cliAction
    app.Before = func(c *cli.Context) error {
        conf := config.GetConfig()
        if len(c.StringSlice("command")) > 0 {
            // StringSliceFlag doesn't support Destination
            conf.Commands = c.StringSlice("command")
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
    conf := config.GetConfig()
    for _, command := range conf.Commands {
        fmt.Println("Command: "+command)
    }
    r := router.NewRouter()
    http.ListenAndServe(conf.Listen+":"+conf.Port, r)
    return nil
}
