package server

import (
	"encoding/json"
	"list_me/db"
	"net/http"
	"sort"
)

func (s *Server) CreateList(w http.ResponseWriter, r *http.Request) {
	newList := db.List{}
	err := json.NewDecoder(r.Body).Decode(&newList)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = s.Store.CreateList(&newList)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	statusResponse(w, http.StatusCreated)
}

func (s *Server) GetList(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	secret := r.URL.Query().Get("secret")

	list, err := s.Store.GetList(id, secret)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	//get list rows
	rows, err := s.Store.GetListRows(list.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	list.Rows = rows

	sort.Slice(list.Rows, func(i, j int) bool {
		return list.Rows[i].CreatedAt < list.Rows[j].CreatedAt
	})

	jsonResponse(w, http.StatusOK, list)
}
