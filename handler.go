package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Request/response DTOs
type createReq struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type updateReq struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
}

// List response with pagination metadata
type listResp struct {
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
	Pages int     `json:"pages"`
	Tasks []*Task `json:"tasks"`
}

type Handler struct {
	svc *TaskService
}

func NewHandler(svc *TaskService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", h.health).Methods("GET")
	r.HandleFunc("/tasks", h.create).Methods("POST")
	r.HandleFunc("/tasks", h.list).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", h.get).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", h.update).Methods("PUT")
	r.HandleFunc("/tasks/{id:[0-9]+}", h.delete).Methods("DELETE")
	return r
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		writeErr(w, http.StatusBadRequest, "title required")
		return
	}
	t, err := h.svc.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func parseIntQuery(r *http.Request, key string, defaultVal int) int {
	vals := r.URL.Query()
	s := vals.Get(key)
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}

func parseStatusQuery(r *http.Request) (*TaskStatus, error) {
	s := strings.TrimSpace(r.URL.Query().Get("status"))
	if s == "" {
		return nil, nil
	}
	ts := TaskStatus(s)
	switch ts {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return &ts, nil
	default:
		return nil, errors.New("invalid status")
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	page := parseIntQuery(r, "page", 1)
	size := parseIntQuery(r, "size", 10)
	if size > 100 {
		size = 100
	}
	status, err := parseStatusQuery(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid status filter")
		return
	}
	tasks, total, err := h.svc.ListTasks(r.Context(), page, size, status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	pages := (total + size - 1) / size
	writeJSON(w, http.StatusOK, listResp{
		Total: total,
		Page:  page,
		Size:  size,
		Pages: pages,
		Tasks: tasks,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	t, err := h.svc.GetTask(context.Background(), id)
	if err != nil {
		if err == ErrNotFound {
			writeErr(w, http.StatusNotFound, "task not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	current, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		if err == ErrNotFound {
			writeErr(w, http.StatusNotFound, "task not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	title := current.Title
	if req.Title != nil {
		if strings.TrimSpace(*req.Title) == "" {
			writeErr(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		title = *req.Title
	}
	desc := current.Description
	if req.Description != nil {
		desc = *req.Description
	}
	status := current.Status
	if req.Status != nil {
		status = *req.Status
	}

	updated, err := h.svc.UpdateTask(r.Context(), id, title, desc, status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := h.svc.DeleteTask(r.Context(), id)
	if err != nil {
		if err == ErrNotFound {
			writeErr(w, http.StatusNotFound, "task not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
