package router

import (
    "os"
    //"fmt"
    "log"
    "time"
    "bytes"
    "errors"
    "strings"
    "os/exec"
    "net/http"
    "math/rand"
    //"math/rand"

    "github.com/miralgj/si/pkg/config"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "github.com/go-chi/jwtauth/v5"
    "github.com/go-chi/chi/v5/middleware"
)

type CommandRequest struct {
    Name string `json:"name"`
    Args []string `json:"args,omitempty"`
}

func (c *CommandRequest) Bind(r *http.Request) error {
   if c.Name == "" {
      return errors.New("name not defined")
   }
   return nil
}

type CommandResponse struct {
    Cmd string `json:"cmd"`
    //Delta
    //End
    Msg string `json:"msg"`
    Rc int `json:"rc"`
    //Start
    Stderr string `json:"stderr"`
    StderrLines []string `json:"stderr_lines"`
    Stdout string `json:"stdout"`
    StdoutLines []string `json:"stdout_lines"`
}

func (c *CommandResponse) Render(w http.ResponseWriter, r *http.Request) error {
   return nil
}

func NewCommandResponse() *CommandResponse {
    resp := &CommandResponse{}
    return resp
}

type ShowConfigResponse struct {
    *config.Options
}

func (c *ShowConfigResponse) Render(w http.ResponseWriter, r *http.Request) error {
   // Don't have to do anything special to render
   return nil
}

func NewShowConfigResponse() *ShowConfigResponse {
    resp := &ShowConfigResponse{config.Config}
    return resp
}

func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // Set up basic middlewares
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Set up basic auth
    if (config.Config.BasicAuth) {
        log.Println("Basic authentication is enabled")
        creds := make(map[string]string)
        creds[config.Config.BasicAuthUser] = config.Config.BasicAuthPass
        r.Use(middleware.BasicAuth("si", creds))
    }

    // Set up token auth
    if (config.Config.TokenAuth) {
        var key []byte
        if (config.Config.TokenKey != "") {
            key = []byte(config.Config.TokenKey)
        } else {
            key = RandomString(32)
        }
        tokenAuth := jwtauth.New("HS256", key, nil)
        _, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"authenticated": true})
        log.Println("Token authentication is enabled")
        log.Println("Bearer token: "+tokenString)
        r.Use(jwtauth.Verifier(tokenAuth))
        r.Use(jwtauth.Authenticator)
    }

    // Set up files route
    if (config.Config.Files) {
        r.Get("/files", http.RedirectHandler("/files/", 301).ServeHTTP)
        r.Get("/files/*", FilesHandler)
    }

    r.Group(func(r chi.Router) {
        r.Use(render.SetContentType(render.ContentTypeJSON))
        r.Get("/", ShowConfigHandler)
    })

    r.Group(func(r chi.Router) {
        r.Use(render.SetContentType(render.ContentTypeJSON))
        r.Use(middleware.Timeout(time.Duration(config.Config.Timeout) * time.Second))
        r.Post("/", RunCommandWithArgsHandler)
    })
    return r
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
    fs := http.StripPrefix("/files", http.FileServer(http.Dir(config.Config.FilesDir)))
    fs.ServeHTTP(w, r)
    return
}

func RunCommand(data *CommandRequest, resp *CommandResponse, done chan<- bool) {
    cmd := exec.Command(config.Config.Commands[data.Name], data.Args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    // Run requested command
    err := cmd.Run()

    resp.Stderr = string(stderr.Bytes())
    resp.Stdout = string(stdout.Bytes())
    resp.StdoutLines = strings.Split(resp.Stdout, "\n")
    resp.StderrLines = strings.Split(resp.Stderr, "\n")

    var (
        ee *exec.ExitError
        pe *os.PathError
    )

    if errors.As(err, &ee) {
        // Non-zero exit code
        resp.Msg = err.Error()
        resp.Rc = ee.ExitCode()
    } else if errors.As(err, &pe) {
        resp.Msg = err.Error()
        resp.Rc = 1
    } else if err != nil {
        resp.Msg = err.Error()
        resp.Rc = 1
    }
    done <- true
    return
}

func RunCommandWithArgsHandler(w http.ResponseWriter, r *http.Request) {
    data := &CommandRequest{}
    resp := &CommandResponse{}

    // Decode the POST json data
    if err := render.Bind(r, data); err != nil {
        resp.Msg = err.Error()
        resp.Rc = 1
        render.Status(r, http.StatusBadRequest)
        render.Render(w, r, resp)
        return
    }

    // Make sure requested command is defined
    if _, ok := config.Config.Commands[data.Name]; !ok {
        resp.Msg = "command not found - "+data.Name
        resp.Rc = 127
        render.Status(r, http.StatusBadRequest)
        render.Render(w, r, resp)
        return
    }

    resp.Cmd = data.Name
    resp.Rc = 0

    done := make(chan bool, 1)
    go RunCommand(data, resp, done)

    select {
        case <-r.Context().Done():
            resp.Msg = "command timed out"
            resp.Rc = 1
            render.Status(r, http.StatusGatewayTimeout)
            render.Render(w, r, resp)
            return
        case <-done:
            render.Status(r, http.StatusOK)
            render.Render(w, r, resp)
            return
    }

    render.Status(r, http.StatusOK)
    render.Render(w, r, resp)
    return
}

func ShowConfigHandler(w http.ResponseWriter, r *http.Request) {
    render.Status(r, http.StatusOK)
    render.Render(w, r, NewShowConfigResponse())
    return
}

func RandomString(length int) []byte {
    const charset = "abcdefghijklmnopqrstuvwxyz" +
          "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    var seededRand *rand.Rand = rand.New(
        rand.NewSource(time.Now().UnixNano()))

    s := make([]byte, length)
    for i := range s {
        s[i] = charset[seededRand.Intn(len(charset))]
    }
    return s
}
