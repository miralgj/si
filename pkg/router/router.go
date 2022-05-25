package router

import (
    "os"
    "bytes"
    "errors"
    "os/exec"
    "net/http"

    "github.com/miralgj/si/pkg/config"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
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
    Stdout string `json:"stdout"`
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

    // Set up middlewares
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(render.SetContentType(render.ContentTypeJSON))

    r.Get("/", ShowConfig)
    r.Post("/", RunCommand)
    return r
}

func RunCommand(w http.ResponseWriter, r *http.Request) {
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

    cmd := exec.Command(config.Config.Commands[data.Name], data.Args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    // Run requested command
    err := cmd.Run()

    resp.Stdout = string(stdout.Bytes())
    resp.Stderr = string(stderr.Bytes())

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

    render.Status(r, http.StatusOK)
    render.Render(w, r, resp)
    return
}

func ShowConfig(w http.ResponseWriter, r *http.Request) {
    render.Render(w, r, NewShowConfigResponse())
    return
}
