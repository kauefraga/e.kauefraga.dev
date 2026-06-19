package shortener

import (
	"net/http"
)

func RetrieveLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if code == "" {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}

	li, ok := LinksDb[code]

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", li.LongUrl)
	w.WriteHeader(http.StatusPermanentRedirect)
}
