package apiserver

import (
	"encoding/json"
	"mjcomparer/internal/app/diffimage"
	"mjcomparer/internal/app/store"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/compare", s.authMiddleware(s.handleCompare())).Methods("POST")
}

func (s *server) authMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Api-key")

		keyExists := s.store.Auth().KeyExists(apiKey)

		if keyExists == false {
			s.error(w, r, http.StatusForbidden, store.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"Error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) handleCompare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, 500, err)
			return
		}
		src := r.Form.Get("src")
		dst := r.Form.Get("dst")

		result, err := diffimage.CompareImagesByUrls(src, dst)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, result)
	}
}
