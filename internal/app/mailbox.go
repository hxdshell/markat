package app

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (a *App) MailboxListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &ApiResponse{}

	mbnames, err := a.Core.MailBoxList()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Data = nil
		response.Message = "could not fetch mailboxes"
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Data = mbnames
	response.Message = "mailboxes fetched successfully"
	json.NewEncoder(w).Encode(response)
}
