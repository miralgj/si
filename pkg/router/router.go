package router

import (
    //"log"
    "fmt"
    "errors"
    "net/http"
    //"encoding/json"

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
      return errors.New("missing required fields")
   }
   return nil
}

type ListResponse struct {
    String string `json:"string"`
    Map map[string]string `json:"map"`
}

func NewListResponse() *ListResponse {
    resp := &ListResponse{
        String:"test",
        Map: map[string]string{
            "key": "value",
        },
    }
    return resp
}

func (lr *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
   // Pre-processing before a response is marshalled and sent across the wire
   //rd.Elapsed = 10
   return nil
}

func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // Set up middlewares
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(render.SetContentType(render.ContentTypeJSON))

    r.Get("/", ListCommands)
    r.Post("/", RunCommand)
    return r
}

func RunCommand(w http.ResponseWriter, r *http.Request) {
    conf := config.GetConfig()
    data := &CommandRequest{}
    if err := render.Bind(r, data); err != nil {
        w.Write([]byte("bad request"))
        return
    }

    fmt.Println("name="+data.Name+",command="+conf.Commands[data.Name])

    render.Status(r, http.StatusCreated)
    w.Write([]byte("good"))
}

func ListCommands(w http.ResponseWriter, r *http.Request) {
    render.Render(w, r, NewListResponse())
    return
}
