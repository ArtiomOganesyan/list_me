package server

import (
	"fmt"
	"list_me/db"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr  string
	Store *db.Storage
}

func NewServer(addr string, store *db.Storage) *Server {
	return &Server{Addr: addr, Store: store}
}

func (s *Server) Init() *http.Server {
	fmt.Println("Starting server on", s.Addr)
	r := mux.NewRouter()
	r.Use(panicRecover)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve HTML
	r.HandleFunc("/{element}", s.HandleHTML).Methods("GET")

	// Serve API
	r.HandleFunc("/api/list", s.CreateList).Methods("POST")
	r.HandleFunc("/api/list", s.GetList).Methods("GET")
	r.HandleFunc("/api/list/{id}/row", s.CreateListRow).Methods("POST")
	r.HandleFunc("/api/row/{id}", s.ChangeRowState).Methods("PATCH")
	r.HandleFunc("/api/row/{id}", s.DeleteRow).Methods("DELETE")

	return &http.Server{Addr: ":" + s.Addr, Handler: r}
}
