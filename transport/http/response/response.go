package response

import (
	"encoding/json"
	"net/http"

	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/helpers/logger"
)

// Base is the base object of all responses
type Base struct {
	Data       *interface{} `json:"data,omitempty"`
	Metadata   *interface{} `json:"metadata,omitempty"`
	Error      *string      `json:"error,omitempty"`
	Message    *string      `json:"message,omitempty"`
	Code       *int         `json:"code,omitempty"`
	Page       *int         `json:"page,omitempty"`
	Limit      *int         `json:"limit,omitempty"`
	Total      *int         `json:"total,omitempty"`
	TotalPage  *int         `json:"total_page,omitempty"`
	CodeStatus *string      `json:"code_status,omitempty"`
}

// NoContent sends a response without any content
func NoContent(w http.ResponseWriter) {
	respond(w, http.StatusNoContent, nil)
}

// WithMessage sends a response with a simple text message
func WithMessage(w http.ResponseWriter, code int, message string) {
	respond(w, code, Base{Message: &message})
}

// WithJSON sends a response containing a JSON object
func WithJSON(w http.ResponseWriter, code int, jsonPayload interface{}, message string) {
	respond(w, code, Base{Data: &jsonPayload, Message: &message, Code: &code})
}

// WITHJSON data sends a response with pagination
func WithJSONPagination(w http.ResponseWriter, code int, jsonPayload interface{}, page int, limit int, total int, totalPage int, message string) {
	respond(w, code, Base{Data: &jsonPayload, Message: &message, Code: &code, Page: &page, Limit: &limit, Total: &total, TotalPage: &totalPage})
}

// WithMetadata sends a response containing a JSON object with metadata
func WithMetadata(w http.ResponseWriter, code int, jsonPayload interface{}, metadata interface{}) {
	respond(w, code, Base{Data: &jsonPayload, Metadata: &metadata})
}

// WithError sends a response with an error message
func WithError(w http.ResponseWriter, err error) {
	code := failure.GetCode(err)
	errMsg := err.Error()
	respond(w, code, Base{Error: &errMsg})
}

// WithPreparingShutdown sends a default response for when the server is preparing to shut down
func WithPreparingShutdown(w http.ResponseWriter) {
	WithMessage(w, http.StatusServiceUnavailable, "SERVER PREPARING TO SHUT DOWN")
}

// WithUnhealthy sends a default response for when the server is unhealthy
func WithUnhealthy(w http.ResponseWriter) {
	WithMessage(w, http.StatusServiceUnavailable, "SERVER UNHEALTHY")
}

func WithJSONCodeStatus(w http.ResponseWriter, code int, jsonPayload interface{}, message string, codeStatus string) {
	respond(w, code, Base{Data: &jsonPayload, Message: &message, Code: &code, CodeStatus: &codeStatus})
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		logger.ErrorWithStack(err)
	}
}
