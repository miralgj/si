package main

import (
    "fmt"
    "log"
    "os"
    "github.com/miralgj/si/pkg/config"
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
    for _, command := range c.StringSlice("cmd") {
        fmt.Println(command)
    }
    return nil
}
