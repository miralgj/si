package config

import (
    "github.com/urfave/cli/v2"
)

type Options struct {
    Commands map[string]string `json:"commands"`
    Listen string `json:"listen-host"`
    Port string `json:"port"`
    Timeout int `json:"timeout"`
    Tls bool `json:"tls"`
    TlsCert string `json:"tls-cert"`
    TlsKey string `json:"tls-key"`
}

var Config *Options

func GetFlags() []cli.Flag {
    f := []cli.Flag{
        &cli.StringSliceFlag{
            Name:   "command",
            Usage:  "command to expose",
            EnvVars: []string{"COMMANDS"},
            Required: true,
        },
        &cli.StringFlag{
            Name:   "listen-host",
            Usage:  "specifies the host to listen on",
            EnvVars: []string{"LISTEN_HOST"},
            Value: "0.0.0.0",
            Destination: &Config.Listen,
        },
        &cli.StringFlag{
            Name:   "port",
            Usage:  "port to listen on",
            EnvVars: []string{"PORT"},
            Value: "3000",
            Destination: &Config.Port,
        },
        &cli.IntFlag{
            Name:   "timeout",
            Usage:  "timeout for commands",
            EnvVars: []string{"TIMEOUT"},
            Value: 90,
            Destination: &Config.Timeout,
        },
        &cli.StringFlag{
            Name:   "tls-cert",
            Usage:  "path to tls certificate chain file",
            EnvVars: []string{"TLS_CERT"},
            Destination: &Config.TlsCert,
        },
        &cli.StringFlag{
            Name:   "tls-key",
            Usage:  "path to tls key file",
            EnvVars: []string{"TLS_KEY"},
            Destination: &Config.TlsKey,
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
