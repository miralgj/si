package router

import (
    "log"
    "fmt"
    "errors"
    "net/http"
    "encoding/json"

    "github.com/miralgj/si/pkg/config"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "github.com/go-chi/chi/v5/middleware"
)

type Command struct {
}

type CommandRequest struct {
    Name string `json:"name"`
    Args []string `json:"args,omitempty"`
}

func (c *CommandRequest) Bind(r *http.Request) error {
   if c.Command == nil {
      return errors.New("missing required fields")
   }
   return nil
}

type ListResponse struct {
    Commands []Commands `json:"commands"`
}
// Create a response for a booking
func NewListResponse(l *model.Booking) *ListResponse {
   resp := &ListResponse{}
   return resp
}

// Render booking response
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

    cmd := data.Command
    fmt.Println("name="+cmd.Name+",command="+conf.Commands[cmd.Name])

    render.Status(r, http.StatusCreated)
    w.Write([]byte("good"))
}

func ListCommands(w http.ResponseWriter, r *http.Request) {
    conf := config.GetConfig()
    resp := make(map[string]interface{})
    render.Render(w, r, reqmodel.NewListResponse(booking))
    return
}
