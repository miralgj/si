package router

import (
    "os"
    "fmt"
    "time"
    "bytes"
    "errors"
    "os/exec"
    "net/http"
    "math/rand"

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

    r.Get("/", ShowConfigHandler)
    //r.Post("/", RunCommandWithArgsHandler)
    r.Group(func(r chi.Router) {
        r.Use(middleware.Timeout(time.Duration(config.Config.Timeout) * time.Second))

        r.Post("/", RunCommandWithArgsHandler)

        r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
            rand.Seed(time.Now().Unix())

            // Processing will take 1-5 seconds.
            processTime := time.Duration(rand.Intn(4)+1) * time.Second

            select {
            case <-r.Context().Done():
                return

            case <-time.After(processTime):
                // The above channel simulates some hard work.
            }

            w.Write([]byte(fmt.Sprintf("Processed in %v seconds\n", processTime)))
        })
    })
    return r
}

func RunCommand(data *CommandRequest, resp *CommandResponse, done chan<- bool) {
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
