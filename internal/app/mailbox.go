package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Cause   any    `json:"cause"`
}

type SelectMbRequest struct {
	Mailbox string `json:"mailbox"`
}

func (a *App) MailboxListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	response := &ApiResponse{}

	mbnames, err := a.Core.MailBoxList(ctx)

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

func (a *App) SelectMailBoxHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	response := &ApiResponse{}
	response.Data = nil

	var mbReq SelectMbRequest
	err := json.NewDecoder(r.Body).Decode(&mbReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "could not select mailbox"
		response.Cause = "invalid request data"
		json.NewEncoder(w).Encode(response)
		return
	}

	mbInfo, err := a.Core.SelectMailBox(ctx, mbReq.Mailbox)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "could not select mailbox"
		response.Cause = "invalid mailbox name"
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Message = "mailbox is now set to " + mbReq.Mailbox
	response.Data = mbInfo
	json.NewEncoder(w).Encode(response)
}

func (a *App) FetchEnvelopes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	response := &ApiResponse{}

	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	envelopeData, err := a.Core.FetchEnvelopes(ctx, page, 15)
	if err != nil {
		if err.Error() == "404" {
			w.WriteHeader(http.StatusNotFound)
			response.Message = "page limit exceeded"
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = "could not fetch envelopes"
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Data = *envelopeData
	response.Message = "envelopes fetched"
	json.NewEncoder(w).Encode(response)
}
