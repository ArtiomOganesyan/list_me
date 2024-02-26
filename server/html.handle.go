package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) HandleHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	element := vars["element"]
	fmt.Println("Handling HTML request for", element)
	http.ServeFile(w, r, "html/"+element+".html")
}
