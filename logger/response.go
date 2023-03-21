package logger

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type StatusWriter struct {
	http.ResponseWriter
	Status int
}

func (s *StatusWriter) WriteHeader(status int) {
	s.Status = status
	s.ResponseWriter.WriteHeader(status)
}

func (w *StatusWriter) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = 200
	}
	n, e := w.ResponseWriter.Write(b)
	return n, e
}

func ProcessResponseBody(b bool, res interface{}) []byte {
	response := Response{
		Success: b,
		Result:  res,
	}

	bz, _ := json.Marshal(response)
	return bz
}
