# Si

Si takes provided commands and exposes a simple RESTful API in order to trigger them.

Is this a bad idea? Si!

## Command-Line Options

* `--command`\
  Command to expose (can be used multiple times). The command base name is used as the command name

* `--basic-auth-user`\
  Username for basic http authentication

* `--basic-auth-pass`\
  Password for basic http authentication

* `--files-dir`\
  Path to directory to serve files under `/files` route

* `--listen-host`\
  Specifies the host to listen on\
  Default: `0.0.0.0`

* `--port`\
  Port to listen on\
  Default: `3000`

* `--timeout`\
  Timeout in seconds before command is cancelled\
  Default: `90`

* `--token-auth`\
  Use token authentication

* `--token-key`\
  Secret key for token authentication

* `--tls-cert`\
  Path to tls certificate chain file

* `--tls-key`\
  Path to tls key file

## Usage

### Exposing Commands

In this example, we'll expose the `/usr/bin/ps` command and set a timeout of `10` seconds.

```
si --command /usr/bin/ps --timeout 10
```

### Get Configuration

You can get the server configuration with a GET request on `/`.

```
$ curl -s http://10.0.0.1:3000/ | jq .
{
  "commands": {
    "ps": "/usr/bin/ps"
  },
  "listen-host": "0.0.0.0",
  "port": "3000",
  "timeout": 10
}
```

### Run Commands

You can POST to `/` to run a defined command.

```
curl -s -X POST -d '{"name":"ps", "args": ["-l"]}' http://10.0.0.1:3000/
{
  "cmd": "ps",
  "msg": "",
  "rc": 0,
  "stderr": "",
  "stdout": "F S   UID     PID    PPID  C PRI  NI ADDR SZ WCHAN  TTY          TIME CMD\n0 S  1001    3109    3046  0  80   0 -  3422 -      pts/1    00:00:01 zsh\n0 S  1001    9457    3109  0  80   0 - 326535 -     pts/1    00:00:00 go\n0 S  1001    9521    9457  0  80   0 - 288997 -     pts/1    00:00:00 main\n4 R  1001    9618    9521  0  80   0 -  2554 -      pts/1    00:00:00 ps\n"
}
```

### Serving Output Files

If an exposed command generates files, you can expose a directory with the `--files-dir` option to make them available for download at http://${IP}:${PORT}/files/.

```
si --command /usr/bin/customscript --files-dir /tmp/outputdir
```
