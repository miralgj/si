package config

import (
    "github.com/urfave/cli/v2"
    //"github.com/spf13/viper"
)

type config struct {
    Commands map[string]string
    CommandNames []string 
    Listen string
    Port string
}

var Config config

func GetFlags() []cli.Flag {
    f := []cli.Flag{
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

func GetConfig() *config {
    return &Config
}

func initConfig() {
    // Initialize config with defaults
    Config = config{
        Commands: make(map[string]string),
        Listen: "0.0.0.0",
        Port: "3000",
    }
}

func init() {
    initConfig()
}
