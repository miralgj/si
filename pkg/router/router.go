package router

import (
    "log"
    "fmt"
    "errors"
    "net/http"
    "path/filepath"
    "encoding/json"

    "github.com/miralgj/si/pkg/config"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "github.com/go-chi/chi/v5/middleware"
)

type Command struct {
    Name string `json:"name"`
    //Args []string `json:"args,omitempty"`
}

type CommandRequest struct {
   *Command
}

func (c *CommandRequest) Bind(r *http.Request) error {
   if c.Command == nil {
      return errors.New("missing required fields")
   }
   return nil
}

var cmdParams map[string]string

func initCommandRoutes() {
}

func NewRouter() *chi.Mux {
    conf := config.GetConfig()
    cmdParams = make(map[string]string)
    for _, cmd := range conf.Commands {
        key := filepath.Base(cmd)
        cmdParams[key] = cmd
    }

    r := chi.NewRouter()

    // Set up middlewares
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(render.SetContentType(render.ContentTypeJSON))

    r.Route("/exec", func(r chi.Router) {
        //r.Get("/{cmd}", exec)
        r.Post("/", execArgs)
    })

    r.Get("/list", list)
    r.Get("/ping", ping)
    return r
}

func execArgs(w http.ResponseWriter, r *http.Request) {
    data := &CommandRequest{}
    if err := render.Bind(r, data); err != nil {
        //render.Render(w, r, networkError.ErrInvalidRequest(err))
        w.Write([]byte("bad request"))
        return
    }

    cmd := data.Command   // Handles logic here
    fmt.Println("name="+cmd.Name)

    // generates a response with 201 status code and json-encoded 
    // booking data
    render.Status(r, http.StatusCreated)
    w.Write([]byte("good"))
}

func exec(w http.ResponseWriter, r *http.Request) {
    cmd := chi.URLParam(r, "cmd")
    if _, ok := cmdParams[cmd]; ok {
        w.Write([]byte("yes"))
    } else {
        w.Write([]byte("no"))
    }
}

func list(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    conf := config.GetConfig()
    resp := make(map[string]interface{})
    resp["commands"] = conf.Commands
    jsonResp, err := json.Marshal(resp)
    if err != nil {
        log.Fatalf("Error happened in JSON marshal. Err: %s", err)
    }
    w.Write(jsonResp)
    return
}

func ping(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    resp := make(map[string]string)
    resp["msg"] = "pong"
    jsonResp, err := json.Marshal(resp)
    if err != nil {
        log.Fatalf("Error happened in JSON marshal. Err: %s", err)
    }
    w.Write(jsonResp)
    return
}
