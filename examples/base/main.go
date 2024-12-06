package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		var json_body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&json_body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(json_body)

		return_data := map[string]interface{}{
			"status": "ok",
		}
		json.NewEncoder(w).Encode(return_data)
	})
	log.Println("Server started on port 7900")
	http.ListenAndServe(":7900", nil)
}
