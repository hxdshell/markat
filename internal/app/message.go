package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) FetchMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mb := vars["mb"]

	parts, err := a.Core.FetchMessage(ctx, mb, uint32(uid))
	if err != nil {
		if err.Error() == "404" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
		}
		return
	}
	if len(parts) > 0 {
		w.WriteHeader(http.StatusOK)
		// Only get text/plain (the first item) for now
		_, err := w.Write(parts[0].Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", parts[0].Header.Get("Content-Type"))
	}
}
