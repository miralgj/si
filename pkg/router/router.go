package router

import (
    "fmt"
    "errors"
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
    if err := render.Bind(r, data); err != nil {
        w.Write([]byte("bad request"))
        return
    }

    fmt.Println("name="+data.Name+",command="+config.Config.Commands[data.Name])

    render.Status(r, http.StatusCreated)
    w.Write([]byte("good"))
}

func ShowConfig(w http.ResponseWriter, r *http.Request) {
    render.Render(w, r, NewShowConfigResponse())
    return
}
