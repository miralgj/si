package config

import (
    "github.com/urfave/cli/v2"
    //"github.com/spf13/viper"
)

type Options struct {
    Commands map[string]string `json:"commands"`
    Config string `json:"config"`
    Listen string `json:"listen-host"`
    Port string `json:"port"`
}

var Config *Options

func GetFlags() []cli.Flag {
    f := []cli.Flag{
        &cli.StringFlag{
            Name:   "config",
            Aliases: []string{"c"},
            Usage:  "si configuration file",
            EnvVars: []string{"CONFIG"},
            Destination: &Config.Config,
        },
        &cli.StringFlag{
            Name:   "listen-host",
            Aliases: []string{"l"},
            Usage:  "specifies the host to listen on",
            EnvVars: []string{"LISTEN_HOST"},
            Value: "0.0.0.0",
            Destination: &Config.Listen,
        },
        &cli.StringFlag{
            Name:   "port",
            Aliases: []string{"p"},
            Usage:  "port to listen on",
            EnvVars: []string{"PORT"},
            Value: "3000",
            Destination: &Config.Port,
        },
        &cli.StringSliceFlag{
            Name:   "command",
            Aliases: []string{"cmd"},
            Usage:  "command to expose",
            EnvVars: []string{"COMMANDS", "CMDS"},
            Required: true,
        },
    }
    return f
}

func initConfig() {
    // Initialoize config with defaults
    Config = &Options{
        Commands: make(map[string]string),
        Listen: "0.0.0.0",
        Port: "3000",
    }
}

func init() {
    initConfig()
}
