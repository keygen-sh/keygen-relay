package server_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimLicense_NewNode_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		ClaimLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				License: &licenses.License{
					File: []byte("test_license_file"),
					Key:  "test_license_key",
				},
				Status: licenses.OperationStatusCreated,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodPut, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp server.ClaimLicenseResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test_license_file"), resp.LicenseFile)
	assert.Equal(t, "test_license_key", resp.LicenseKey)
}

func TestClaimLicense_ExistingNode_Extended(t *testing.T) {
	manager := &testutils.FakeManager{
		ClaimLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				License: &licenses.License{
					File: []byte("test_license_file"),
					Key:  "test_license_key",
				},
				Status: licenses.OperationStatusExtended,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodPut, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)

	var resp server.ClaimLicenseResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test_license_file"), resp.LicenseFile)
	assert.Equal(t, "test_license_key", resp.LicenseKey)
}

func TestClaimLicense_HeartbeatDisabled_Conflict(t *testing.T) {
	manager := &testutils.FakeManager{
		ClaimLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				Status: licenses.OperationStatusConflict,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodPut, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), "License claim conflict, heartbeat disabled")
}

func TestClaimLicense_NoLicensesAvailable(t *testing.T) {
	manager := &testutils.FakeManager{
		ClaimLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				Status: licenses.OperationStatusNoLicensesAvailable,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodPut, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusGone, rr.Code)
	assert.Contains(t, rr.Body.String(), "No licenses available")
}

func TestClaimLicense_InternalServerError(t *testing.T) {
	manager := &testutils.FakeManager{
		ClaimLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return nil, errors.New("database error")
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodPut, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to claim license")
}

func TestReleaseLicense_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		ReleaseLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				Status: licenses.OperationStatusSuccess,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodDelete, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.Bytes())
}

func TestReleaseLicense_NotFound(t *testing.T) {
	manager := &testutils.FakeManager{
		ReleaseLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return &licenses.LicenseOperationResult{
				Status: licenses.OperationStatusNotFound,
			}, nil
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodDelete, "/v1/nodes/non_existent_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Claim not found")
}

func TestReleaseLicense_InternalServerError(t *testing.T) {
	manager := &testutils.FakeManager{
		ReleaseLicenseFn: func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
			return nil, errors.New("database error")
		},
	}

	handler := server.NewHandler(manager)

	req := httptest.NewRequest(http.MethodDelete, "/v1/nodes/test_fingerprint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to release license")
}
