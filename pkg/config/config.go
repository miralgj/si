package config

import (
    "github.com/urfave/cli/v2"
)

var Flags []cli.Flag

func init() {
    Flags = []cli.Flag{
        &cli.StringFlag{
            Name:   "config",
            Aliases: []string{"c"},
            Usage:  "load configuration from `FILE`",
            EnvVars: []string{"CONFIG"},
        },
        &cli.StringSliceFlag{
            Name:   "cmd",
            Usage:  "command to expose `FILE`",
            EnvVars: []string{"CMDS"},
        },
    }
}
