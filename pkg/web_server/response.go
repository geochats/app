package web_server

import (
	"encoding/json"
	"net/http"
)

// responseWithSuccessJSON function send body like json with status 200
func (s *WebServer) responseWithSuccessJSON(w http.ResponseWriter, payload interface{}) {
	s.sendResponse(w, http.StatusOK, &JSONResponse{
		Success: true,
		Data:    payload,
	})
}

// responseWithErrorJSON function send body like json with status 200
func (s *WebServer) responseWithErrorJSON(w http.ResponseWriter, err error) {
	s.sendResponse(w, http.StatusInternalServerError, &JSONResponse{
		Success: false,
		Error:   err.Error(),
	})
}

// sendResponse method send body like json
func (s *WebServer) sendResponse(w http.ResponseWriter, httpStatus int, body *JSONResponse) {
	data, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorf("errors encoding body to json, body: %v", body)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorf("failed to write data to http.ResponseWriter, body: %v", data)
	}
}

type JSONResponse struct {
	Success bool               `json:"success"`
	Data    interface{}           `json:"data"`
	Error   string `json:"error"`
}

