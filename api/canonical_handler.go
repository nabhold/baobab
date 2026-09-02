package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/service"
)

type canonicalHandler struct {
	service service.CanonicalEntityService
}

func (h canonicalHandler) create(w http.ResponseWriter, r *http.Request) {
	var entity domain.CanonicalEntity
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entity); err != nil {
		problem(w, r, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid canonical entity JSON", false)
		return
	}
	if entity.ID == "" {
		entity.ID = newUUID()
	}
	created, err := h.service.Create(r.Context(), entity)
	if err != nil {
		problem(w, r, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), false)
		return
	}
	w.Header().Set("Location", "/v1/canonical-entities/"+created.ID)
	w.Header().Set("ETag", `"1"`)
	writeJSON(w, http.StatusCreated, created)
}

func (h canonicalHandler) get(w http.ResponseWriter, r *http.Request) {
	entity, err := h.service.Get(r.Context(), chi.URLParam(r, "entityID"))
	if err != nil {
		problem(w, r, http.StatusNotFound, "CANONICAL_ENTITY_NOT_FOUND", err.Error(), false)
		return
	}
	w.Header().Set("ETag", `"`+formatVersion(entity.Version)+`"`)
	writeJSON(w, http.StatusOK, entity)
}

func (h canonicalHandler) lifecycle(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := service.ParseExpectedVersion(r.Header.Get("If-Match"))
		if err != nil {
			problem(w, r, http.StatusPreconditionRequired, "IF_MATCH_REQUIRED", err.Error(), false)
			return
		}
		id := chi.URLParam(r, "entityID")
		var entity domain.CanonicalEntity
		switch action {
		case "validate":
			entity, err = h.service.Validate(r.Context(), id, version)
		case "activate":
			entity, err = h.service.Activate(r.Context(), id, version)
		case "suspend":
			entity, err = h.service.Suspend(r.Context(), id, version)
		case "retire":
			entity, err = h.service.Retire(r.Context(), id, version)
		}
		if err != nil {
			problem(w, r, http.StatusConflict, "CANONICAL_ENTITY_LIFECYCLE_CONFLICT", err.Error(), false)
			return
		}
		w.Header().Set("ETag", `"`+formatVersion(entity.Version)+`"`)
		writeJSON(w, http.StatusOK, entity)
	}
}

func formatVersion(version int64) string {
	return strconv.FormatInt(version, 10)
}
