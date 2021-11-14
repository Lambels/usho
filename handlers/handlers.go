package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Lambels/usho/encoding/json"
	"github.com/Lambels/usho/repo"
	"github.com/gorilla/mux"
)

type Service struct {
	r   repo.Repo
	ctx context.Context
}

func NewService(r repo.Repo, ctx context.Context) http.Handler {
	m := mux.NewRouter()

	s := Service{r: r, ctx: ctx}

	postPath := m.PathPrefix("/url").Methods(http.MethodPost).Subrouter()

	postPath.Use(MiddlewareValidateIn)
	postPath.HandleFunc("/new", s.New)

	m.HandleFunc("/{to}", s.Redirect).Methods(http.MethodGet)

	return m
}

func (s Service) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	out, err := s.r.Get(s.ctx, vars["to"])
	switch err {
	case nil:
		http.Redirect(w, r, out.Initial, http.StatusMovedPermanently)

	case repo.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.Encode(map[string]string{"error": err.Error()}, w)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.Encode(map[string]string{"error": err.Error()}, w)
	}
}

func (s Service) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("ContentType", "application/json")

	in := r.Context().Value(repo.URLKey{}).(repo.URLRequest)

	in.Intial = strings.TrimSpace(in.Intial)

	out, err := s.r.New(s.ctx, in)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.Encode(map[string]string{"error": err.Error()}, w)
		return
	}

	err = json.Encode(out, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.Encode(map[string]string{"error": err.Error()}, w)
		return
	}
}
