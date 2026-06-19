package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LinkBody struct {
	Url string `json:"url"`
}

type LinkStatus string

const (
	StatusActive   LinkStatus = "active"
	StatusInactive LinkStatus = "inactive"
)

type Link struct {
	Id        string     `json:"id"`
	Code      string     `json:"code"`
	LongUrl   string     `json:"long_url"`
	Status    LinkStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
}

func generateShortCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

var LinksDb = make(map[string]Link)

func CreateLink(w http.ResponseWriter, r *http.Request) {
	var b LinkBody

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if b.Url == "" {
		http.Error(w, "empty string", http.StatusBadRequest)
		return
	}

	c, err := generateShortCode()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	linkId, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	li := Link{
		Id:        linkId.String(),
		LongUrl:   b.Url,
		Status:    StatusActive,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 5),
		Code:      c,
	}

	LinksDb[li.Code] = li

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(li)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
