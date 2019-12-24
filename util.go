package main

import (
	"encoding/json"
	"net/http"
)

func jsonTo(w http.ResponseWriter, status bool, message string, data ...interface{}) {
	result := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	_bytes, _ := json.Marshal(result)

	w.Write(_bytes)
}
