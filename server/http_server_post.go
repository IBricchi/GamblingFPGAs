package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *HttpServer) handlePostDynamicTest() http.HandlerFunc {
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

		fmt.Println("Received data: ", data.Info, data.Data)
	}
}
