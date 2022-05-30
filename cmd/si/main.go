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

// Compile with --ldflags="-w -X main.Version=$(VERSION)"
// to set version number
var Version string

func main() {
    app := cli.NewApp()
    app.Name = "Si"
    app.Version = Version
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
    if (config.Config.Tls) {
        err := http.ListenAndServeTLS(config.Config.Listen+":"+config.Config.Port, config.Config.TlsCert, config.Config.TlsKey, r)
        if err != nil {
            return err
        }
    } else {
        err := http.ListenAndServe(config.Config.Listen+":"+config.Config.Port, r)
        if err != nil {
            return err
        }
    }
    return nil
}

func cliBeforeHandler(c *cli.Context) error {
    // Verify both tls-cert and tls-key were used together
    if (c.IsSet("tls-cert") && !c.IsSet("tls-key")) {
      return errors.New("tls-key not set")
    } else if (!c.IsSet("tls-cert") && c.IsSet("tls-key")) {
      return errors.New("tls-cert not set")
    } else if(c.IsSet("tls-cert") && c.IsSet("tls-key")) {
        // Verify both provided files exist
        if _, err := os.Stat(c.String("tls-cert")); errors.Is(err, os.ErrNotExist) {
            log.Fatal(err.Error())
        }
        if _, err := os.Stat(c.String("tls-key")); errors.Is(err, os.ErrNotExist) {
            log.Fatal(err.Error())
        }
        config.Config.Tls = true
    }
    // Get commands from the cli context
    if len(c.StringSlice("command")) > 0 {
        // StringSliceFlag doesn't support Destination
        for _, path := range c.StringSlice("command") {
            // Verify file exists before adding it
            if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
                log.Fatal(err.Error())
            }
            name := filepath.Base(path)
            config.Config.Commands[name] = path
        }
    }
    return nil
}
