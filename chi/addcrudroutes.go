package chi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/h"
)

// Remember you can overshadow routes from returned here router (authed)
func AddCRUDRoutes[CRUDClass any, CreateCRUDClassDto any, UpdateCRUDClassDto any, User any](
	router chi.Router,
	className string,
	listRecords func(ctx context.Context) ([]CRUDClass, error),
	getRecord func(ctx context.Context, id string) (CRUDClass, error),
	createRecord func(ctx context.Context, params CreateCRUDClassDto) (CRUDClass, error),
	deleteRecord func(ctx context.Context, id string) error,
	updateRecord func(ctx context.Context, arg UpdateCRUDClassDto) (CRUDClass, error),
) chi.Router {
	router.Get(fmt.Sprintf("/%ss", className), func(w http.ResponseWriter, r *http.Request) {
		records, err := listRecords(r.Context())
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, records)
	})

	router.Get(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		record, err := getRecord(r.Context(), id)
		if err != nil {
			h.ResNotFound(w, className)
			return
		}
		h.ResSuccess(w, record)
	})

	router.Post(fmt.Sprintf("/%ss", className), func(w http.ResponseWriter, r *http.Request) {
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

	router.Put(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
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

	router.Delete(fmt.Sprintf("/%ss/{id}", className), func(w http.ResponseWriter, r *http.Request) {
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

	return router
}
