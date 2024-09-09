package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-basics/internal/core/domain"
	"github.com/osalomon89/go-basics/internal/core/ports"
)

// type responseError struct {
// 	Message    string
// 	StatusCode int
// }

type handler struct {
	itemService ports.ItemService
}

func NewHandler(itemService ports.ItemService) *handler {
	return &handler{
		itemService: itemService,
	}
}

func (h *handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem domain.Item
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.itemService.AddItem(ctx, newItem)
	if err != nil {
		log.Printf("error inserting item: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//return web.EncodeJSON(w, responseError{Message: "error inserting item", StatusCode: http.StatusInternalServerError}, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item := h.itemService.ReadItem(r.Context(), id)
	if item != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "id not found", http.StatusBadRequest)
}

func (h *handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limitParam := r.URL.Query().Get("limit")
	cursorParam := r.URL.Query().Get("cursor")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var searchAfter []interface{}
	if cursorParam != "" {
		searchAfter = parseCursor(cursorParam)
	}

	items, newCursor, err := h.itemService.GetAllItems(ctx, limit, searchAfter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Items  []domain.Item `json:"items"`
		Cursor string        `json:"cursor,omitempty"`
	}{
		Items: items,
	}

	if newCursor != nil {
		response.Cursor = encodeCursor(newCursor)
	}

	w.Header().Set("Content-Type", "application/json")

	// Retornar la lista de usuarios en formato JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseCursor(cursor string) []interface{} {
	parts := strings.Split(cursor, ",")
	parsed := make([]interface{}, len(parts))
	for i, part := range parts {
		if num, err := strconv.ParseFloat(part, 64); err == nil {
			parsed[i] = num
		} else {
			parsed[i] = part
		}
	}
	return parsed
}

func encodeCursor(cursor []interface{}) string {
	parts := make([]string, len(cursor))
	for i, part := range cursor {
		switch v := part.(type) {
		case float64:
			// Convertimos el número a un string sin formato exponencial
			parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			parts[i] = fmt.Sprintf("%v", part)
		}
	}
	return strings.Join(parts, ",")
}

func (h *handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var existItem domain.Item

	if err := json.NewDecoder(r.Body).Decode(&existItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	existItem.ID = vars["id"]

	result := h.itemService.UpdateItem(r.Context(), existItem)

	if result != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "id not found", http.StatusBadRequest)
}

func (h *handler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s, world!", r.URL.Path[1:])
}
