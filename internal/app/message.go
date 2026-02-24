package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MarkSeenUnseenRequest struct {
	Mb   string   `json:"mb"`
	Uids []uint32 `json:"uids"`
	Seen bool     `json:"seen"`
}
type MoveRequest struct {
	Mb   string   `json:"mb"`
	Uids []uint32 `json:"uids"`
	Dest string   `json:"dest"`
}

func (a *App) FetchMeta(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	response := &ApiResponse{}

	vars := mux.Vars(r)
	mb := vars["mb"]
	uid, err := strconv.Atoi(vars["uid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	meta, err := a.Core.FetchMeta(ctx, mb, uint32(uid))
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			response.Message = "message not found"
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "internal server error while fetching message metadata"
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Message = "metadata fetched successfully"
	response.Data = meta
	json.NewEncoder(w).Encode(response)
}

func (a *App) FetchMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mb := vars["mb"]
	msg, ct, err := a.Core.FetchMessageText(ctx, mb, uint32(uid))
	if err != nil {
		log.Println(err)
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Write(msg)
}
func (a *App) FetchAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mb := vars["mb"]
	specifier := vars["specifier"]
	h, b, err := a.Core.FetchAttachment(ctx, mb, uint32(uid), specifier)
	if err != nil {
		log.Println(err)
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", h.Type)
	w.Header().Set("Content-Disposition", h.Disposition)

	w.Header().Set("Content-Length", strconv.FormatUint(uint64(len(b)), 10))
	w.Write(b)
}

func (a *App) MarkSeenUnseen(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var markReq MarkSeenUnseenRequest
	err := json.NewDecoder(r.Body).Decode(&markReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = a.Core.MarkSeenUnseen(ctx, markReq.Mb, markReq.Uids, markReq.Seen)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var moveReq MoveRequest
	err := json.NewDecoder(r.Body).Decode(&moveReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Core.Move(ctx, moveReq.Mb, moveReq.Uids, moveReq.Dest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
