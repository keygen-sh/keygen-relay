package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
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
	r.HandleFunc("/v1/health", h.HealthCheck).Methods("GET")
	r.HandleFunc("/v1/nodes/{fingerprint}", h.ClaimLicense).Methods("PUT")
	r.HandleFunc("/v1/nodes/{fingerprint}", h.ReleaseLicense).Methods("DELETE")
}

func (h *handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *handler) ClaimLicense(w http.ResponseWriter, r *http.Request) {
	fingerprint := mux.Vars(r)["fingerprint"]

	result, err := h.licenseManager.ClaimLicense(r.Context(), fingerprint)
	if err != nil {
		slog.Error("failed to claim license", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to claim license"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch result.Status {
	case licenses.OperationStatusCreated:
		w.WriteHeader(http.StatusCreated)
		if result.License != nil {
			resp := ClaimLicenseResponse{
				LicenseFile: result.License.File,
				LicenseKey:  result.License.Key,
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	case licenses.OperationStatusExtended:
		w.WriteHeader(http.StatusAccepted)
	case licenses.OperationStatusConflict:
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to claim license due to conflict"})
		return
	case licenses.OperationStatusNoLicensesAvailable:
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "no licenses available"})
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unknown claim status"})
		return
	}
}

func (h *handler) ReleaseLicense(w http.ResponseWriter, r *http.Request) {
	fingerprint := mux.Vars(r)["fingerprint"]

	result, err := h.licenseManager.ReleaseLicense(r.Context(), fingerprint)
	if err != nil {
		slog.Error("failed to release license", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to release license"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch result.Status {
	case licenses.OperationStatusSuccess:
		w.WriteHeader(http.StatusNoContent)
	case licenses.OperationStatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "claim not found"})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unknown release status"})
	}
}
