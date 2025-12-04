package chii

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/h"
)

type Auther[User any] interface {
	GetUserFromToken(ctx context.Context, token string) (*User, error)
}

// Remember you can overshadow routes from returned here router (authed)
func AddAuthedCRUDRoutes[CRUDClass any, CreateCRUDClassDto any, UpdateCRUDClassDto any, User any](
	router chi.Router,
	auther Auther[User],
	className string,
	listRecords func(ctx context.Context) ([]CRUDClass, error),
	getRecord func(ctx context.Context, id string) (CRUDClass, error),
	createRecord func(ctx context.Context, params CreateCRUDClassDto) (CRUDClass, error),
	deleteRecord func(ctx context.Context, id string) error,
	updateRecord func(ctx context.Context, arg UpdateCRUDClassDto) (CRUDClass, error),
) chi.Router {
	auth := router.With(WithAuth(auther))

	auth.Get(fmt.Sprintf("/%ss", className), func(w http.ResponseWriter, r *http.Request) {
		records, err := listRecords(r.Context())
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, records)
	})

	auth.Get(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		record, err := getRecord(r.Context(), id)
		if err != nil {
			h.ResNotFound(w, className)
			return
		}
		h.ResSuccess(w, record)
	})

	auth.Post(fmt.Sprintf("/%ss", className), func(w http.ResponseWriter, r *http.Request) {
		dto, err := h.GetDto[CreateCRUDClassDto](r)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}
		record, err := createRecord(r.Context(), *dto)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, record)
	})

	auth.Put(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		dto, err := h.GetDto[UpdateCRUDClassDto](r)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}

		_, err = getRecord(r.Context(), id)
		if err != nil {
			h.ResNotFound(w, className)
			return
		}

		updated, err := updateRecord(r.Context(), *dto)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, updated)
	})

	auth.Delete(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, err := getRecord(r.Context(), id)
		if err != nil {
			h.ResNotFound(w, className)
			return
		}

		err = deleteRecord(r.Context(), id)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, nil)
	})

	return auth
}
