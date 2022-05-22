package config

import (
    "github.com/urfave/cli/v2"
    //"github.com/spf13/viper"
)

var Flags []cli.Flag

type config struct {
    Listen string
    Port string
}

var defaults struct {
}

func New() *config {
    c := config{
        Listen: "0.0.0.0",
        Port: "3000",
    }
    return &c
}

func init() {
    defaults := New()
    Flags = []cli.Flag{
        &cli.StringFlag{
            Name:   "config",
            Aliases: []string{"c"},
            Usage:  "si configuration file",
            EnvVars: []string{"CONFIG"},
        },
        &cli.StringFlag{
            Name:   "listen-host",
            Aliases: []string{"l"},
            Usage:  "specifies the host to listen on",
            DefaultText: defaults.Listen,
            EnvVars: []string{"LISTEN_HOST"},
        },
        &cli.IntFlag{
            Name:   "port",
            Aliases: []string{"p"},
            Usage:  "port to listen on",
            DefaultText: defaults.Port,
            EnvVars: []string{"PORT"},
        },
        &cli.StringSliceFlag{
            Name:   "cmd",
            Usage:  "command to expose",
            EnvVars: []string{"CMDS"},
        },
    }
}
