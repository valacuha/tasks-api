package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct {
	Store storage.Storage
}

func New(s storage.Storage) *Handler {
	return &Handler{Store: s}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := h.Store.List()
		writeJSON(w, http.StatusOK, tasks)

	case http.MethodPost:
		var input models.Task
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "некорректный формат JSON")
			return
		}
		if strings.TrimSpace(input.Title) == "" {
			writeError(w, http.StatusBadRequest, "поле title обязательно")
			return
		}

		created, err := h.Store.Create(input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "ошибка создания задачи")
			return
		}
		writeJSON(w, http.StatusCreated, created)

	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if idStr == "" || strings.Contains(idStr, "/") {
		writeError(w, http.StatusBadRequest, "некорректный идентификатор задачи")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id должен быть положительным числом")
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, ok := h.Store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "задача не найдена")
			return
		}
		writeJSON(w, http.StatusOK, task)

	case http.MethodPut:
		var input models.Task
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "некорректный формат JSON")
			return
		}
		if strings.TrimSpace(input.Title) == "" {
			writeError(w, http.StatusBadRequest, "поле title обязательно")
			return
		}

		updated, err := h.Store.Update(id, input)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeError(w, http.StatusNotFound, "задача не найдена")
				return
			}
			writeError(w, http.StatusInternalServerError, "ошибка обновления задачи")
			return
		}
		writeJSON(w, http.StatusOK, updated)

	case http.MethodDelete:
		err := h.Store.Delete(id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeError(w, http.StatusNotFound, "задача не найдена")
				return
			}
			writeError(w, http.StatusInternalServerError, "ошибка удаления задачи")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
