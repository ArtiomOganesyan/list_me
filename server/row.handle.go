package server

import (
	"encoding/json"
	"fmt"
	"list_me/db"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) CreateListRow(w http.ResponseWriter, r *http.Request) {
	row := db.ListRow{}

	err := json.NewDecoder(r.Body).Decode(&row)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	row.ListID = vars["id"]

	fmt.Println(row)

	id, err := s.Store.CreateListRow(row)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusCreated, id)
}

func (s *Server) ChangeRowState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rowID := vars["id"]
	done := r.URL.Query().Get("state")

	err := s.Store.ChangeRowState(rowID, done == "true")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	statusResponse(w, http.StatusOK)
}

func (s *Server) DeleteRow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rowID := vars["id"]

	err := s.Store.DeleteRow(rowID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	statusResponse(w, http.StatusOK)
}
