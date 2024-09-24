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
		http.Error(w, "Failed to claim license", http.StatusInternalServerError)
		return
	}

	switch result.Status {
	case licenses.OperationStatusCreated:
		w.WriteHeader(http.StatusCreated)
	case licenses.OperationStatusExtended:
		w.WriteHeader(http.StatusAccepted)
	case licenses.OperationStatusConflict:
		http.Error(w, "License claim conflict, heartbeat disabled", http.StatusConflict)
		return
	case licenses.OperationStatusNoLicensesAvailable:
		http.Error(w, "No licenses available", http.StatusGone)
		return
	default:
		http.Error(w, "Unknown claim status", http.StatusInternalServerError)
		return
	}

	if result.License != nil {
		resp := ClaimLicenseResponse{
			LicenseFile: result.License.File,
			LicenseKey:  result.License.Key,
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error("Failed to write response: " + err.Error())
		}
	}
}

func (h *handler) ReleaseLicense(w http.ResponseWriter, r *http.Request) {
	fingerprint := mux.Vars(r)["fingerprint"]

	result, err := h.licenseManager.ReleaseLicense(r.Context(), fingerprint)
	if err != nil {
		http.Error(w, "Failed to release license", http.StatusInternalServerError)
		return
	}

	switch result.Status {
	case licenses.OperationStatusSuccess:
		w.WriteHeader(http.StatusNoContent)
	case licenses.OperationStatusNotFound:
		http.Error(w, "Claim not found", http.StatusNotFound)
	default:
		http.Error(w, "Unknown release status", http.StatusInternalServerError)
	}
}
