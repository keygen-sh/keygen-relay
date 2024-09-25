package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"log/slog"
	"net/http"
)

type RequestBodyPayload struct {
	Fingerprint string `json:"fingerprint"`
}

type ClaimLicenseResponse struct {
	LicenseFile []byte `json:"license_file"`
	LicenseKey  string `json:"license_key"`
}

type Handler interface {
	RegisterRoutes(r *mux.Router)
}

type handler struct {
	licenseManager licenses.Manager
	Server         *Server
}

func NewHandler(m licenses.Manager) Handler {
	return &handler{
		licenseManager: m,
	}
}

func (h *handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/nodes/{fingerprint}", h.ClaimLicense).Methods("PUT")
	r.HandleFunc("/v1/nodes/{fingerprint}", h.ReleaseLicense).Methods("DELETE")
}

func (h *handler) ClaimLicense(w http.ResponseWriter, r *http.Request) {
	fingerprint := mux.Vars(r)["fingerprint"]

	slog.Debug("Incoming request", "headers", r.Header, "method", r.Method, "request_uri", r.RequestURI)

	result, err := h.licenseManager.ClaimLicense(r.Context(), fingerprint)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to claim license"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch result.Status {
	case licenses.OperationStatusCreated:
		w.WriteHeader(http.StatusCreated)
	case licenses.OperationStatusExtended:
		w.WriteHeader(http.StatusAccepted)
	case licenses.OperationStatusConflict:
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "License claim conflict, heartbeat disabled"})
		return
	case licenses.OperationStatusNoLicensesAvailable:
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "No licenses available"})
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unknown claim status"})
		return
	}

	if result.License != nil {
		resp := ClaimLicenseResponse{
			LicenseFile: result.License.File,
			LicenseKey:  result.License.Key,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (h *handler) ReleaseLicense(w http.ResponseWriter, r *http.Request) {
	fingerprint := mux.Vars(r)["fingerprint"]

	result, err := h.licenseManager.ReleaseLicense(r.Context(), fingerprint)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to release license"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch result.Status {
	case licenses.OperationStatusSuccess:
		w.WriteHeader(http.StatusNoContent)
	case licenses.OperationStatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Claim not found"})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unknown release status"})
	}
}
