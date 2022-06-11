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

var (
    conf = config.Config

    die = log.Fatal

    // Compile with --ldflags="-w -X main.Version=$(VERSION)"
    // to set version number
    Version string
)

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
        die(err)
    }
}

func cliActionHandler(c *cli.Context) error {
    r := router.NewRouter()
    var err error
    log.Println("Starting si server on "+conf.Listen+":"+conf.Port)
    if (conf.Tls) {
        err = http.ListenAndServeTLS(conf.Listen+":"+conf.Port, conf.TlsCert, conf.TlsKey, r)
    } else {
        err = http.ListenAndServe(conf.Listen+":"+conf.Port, r)
    }
    if err != nil {
        return err
    }
    return nil
}

func cliBeforeHandler(c *cli.Context) error {
    // Verify basic and jwt auth weren't used together
    if ((c.IsSet("basic-auth-user") || c.IsSet("basic-auth-pass")) && c.IsSet("jwt-auth")) {
        die("Basic auth and JWT auth are mutually exclusive")
    }
    // Verify both basic-auth-user and basic-auth-pass were used together
    if (c.IsSet("basic-auth-user") || c.IsSet("basic-auth-pass")) {
        dieIfFlagsMissing(c, []string{"basic-auth-user", "basic-auth-pass"})
        conf.BasicAuth = true
    }
    // Verify both tls-cert and tls-key were used together
    if (c.IsSet("tls-cert") || c.IsSet("tls-key")) {
        dieIfFlagsMissing(c, []string{"tls-cert", "tls-key"})
        dieIfFilesMissing([]string{c.String("tls-cert"), c.String("tls-key")})
        conf.Tls = true
    }
    // Get commands from the cli context
    if len(c.StringSlice("command")) > 0 {
        cmds := c.StringSlice("command")
        dieIfFilesMissing(cmds)
        // StringSliceFlag doesn't support Destination so I need
        // to add these to the Config
        for _, path := range cmds {
            name := filepath.Base(path)
            conf.Commands[name] = path
        }
    }
    return nil
}

// Verifies all provided files exist
func dieIfFilesMissing(files []string) {
    for _, f := range files {
        if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
            die(err.Error())
        }
    }
}

// Verifies that all the provided flags are set
func dieIfFlagsMissing(c *cli.Context, flags []string) {
    for _, f := range flags {
        if (!c.IsSet(f)) {
            die("Missing required flag: --"+f)
        }
    }
}
