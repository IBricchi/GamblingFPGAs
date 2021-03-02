package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *HttpServer) handlePostDynamicTest(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		data := staticTestData{}
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if correct data format was send
		if data.Info == "" && data.Data == nil {
			http.Error(w, "Error: Invalid data was send", http.StatusBadRequest)
			return
		}

		if err := h.db.insertTestData(ctx, data); err != nil {
			http.Error(w, "Error: Failed to insert data in DB", http.StatusInternalServerError)
			return
		}

		fmt.Println("Received data: ", data.Info, data.Data)
	}
}
