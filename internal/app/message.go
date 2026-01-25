package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) FetchMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mb := vars["mb"]
	msg, err := a.Core.FetchMessageText(ctx, mb, uint32(uid))
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(msg)
}
