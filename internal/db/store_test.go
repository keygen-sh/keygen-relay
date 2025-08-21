package db

import (
	"context"
	"database/sql"
	"testing"
)

// QuerierInterface defines the methods we need from Queries for testing
type QuerierInterface interface {
	GetLicenses(ctx context.Context) ([]License, error)
	GetLicensesWithPool(ctx context.Context, poolID *int64) ([]License, error)
	GetLicensesWithoutPool(ctx context.Context) ([]License, error)
	GetLicenseByGUID(ctx context.Context, guid string) (License, error)
	GetLicenseWithPoolByGUID(ctx context.Context, params GetLicenseWithPoolByGUIDParams) (License, error)
	GetLicenseWithoutPoolByGUID(ctx context.Context, guid string) (License, error)
	ReleaseLicenseWithPoolByNodeID(ctx context.Context, params ReleaseLicenseWithPoolByNodeIDParams) error
	ReleaseLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) error
	ClaimLicenseWithPoolFIFO(ctx context.Context, params ClaimLicenseWithPoolFIFOParams) (License, error)
	ClaimLicenseWithoutPoolFIFO(ctx context.Context, nodeID *int64) (License, error)
	GetLicenseWithPoolByNodeID(ctx context.Context, params GetLicenseWithPoolByNodeIDParams) (License, error)
	GetLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) (License, error)
	WithTx(tx *sql.Tx) *Queries
}

// mockQueries implements the query interface for testing
type mockQueries struct {
	getLicensesFunc                    func(ctx context.Context) ([]License, error)
	getLicensesWithPoolFunc            func(ctx context.Context, poolID *int64) ([]License, error)
	getLicensesWithoutPoolFunc         func(ctx context.Context) ([]License, error)
	getLicenseByGUIDFunc               func(ctx context.Context, guid string) (License, error)
	getLicenseWithPoolByGUIDFunc       func(ctx context.Context, params GetLicenseWithPoolByGUIDParams) (License, error)
	getLicenseWithoutPoolByGUIDFunc    func(ctx context.Context, guid string) (License, error)
	releaseLicenseWithPoolByNodeIDFunc func(ctx context.Context, params ReleaseLicenseWithPoolByNodeIDParams) error
	releaseLicenseWithoutPoolByNodeIDFunc func(ctx context.Context, nodeID *int64) error
	claimLicenseWithPoolFIFOFunc       func(ctx context.Context, params ClaimLicenseWithPoolFIFOParams) (License, error)
	claimLicenseWithoutPoolFIFOFunc    func(ctx context.Context, nodeID *int64) (License, error)
	getLicenseWithPoolByNodeIDFunc     func(ctx context.Context, params GetLicenseWithPoolByNodeIDParams) (License, error)
	getLicenseWithoutPoolByNodeIDFunc  func(ctx context.Context, nodeID *int64) (License, error)
}

func (m *mockQueries) GetLicenses(ctx context.Context) ([]License, error) {
	if m.getLicensesFunc != nil {
		return m.getLicensesFunc(ctx)
	}
	return []License{}, nil
}

func (m *mockQueries) GetLicensesWithPool(ctx context.Context, poolID *int64) ([]License, error) {
	if m.getLicensesWithPoolFunc != nil {
		return m.getLicensesWithPoolFunc(ctx, poolID)
	}
	return []License{}, nil
}

func (m *mockQueries) GetLicensesWithoutPool(ctx context.Context) ([]License, error) {
	if m.getLicensesWithoutPoolFunc != nil {
		return m.getLicensesWithoutPoolFunc(ctx)
	}
	return []License{}, nil
}

func (m *mockQueries) GetLicenseByGUID(ctx context.Context, guid string) (License, error) {
	if m.getLicenseByGUIDFunc != nil {
		return m.getLicenseByGUIDFunc(ctx, guid)
	}
	return License{}, nil
}

func (m *mockQueries) GetLicenseWithPoolByGUID(ctx context.Context, params GetLicenseWithPoolByGUIDParams) (License, error) {
	if m.getLicenseWithPoolByGUIDFunc != nil {
		return m.getLicenseWithPoolByGUIDFunc(ctx, params)
	}
	return License{}, nil
}

func (m *mockQueries) GetLicenseWithoutPoolByGUID(ctx context.Context, guid string) (License, error) {
	if m.getLicenseWithoutPoolByGUIDFunc != nil {
		return m.getLicenseWithoutPoolByGUIDFunc(ctx, guid)
	}
	return License{}, nil
}

func (m *mockQueries) ReleaseLicenseWithPoolByNodeID(ctx context.Context, params ReleaseLicenseWithPoolByNodeIDParams) error {
	if m.releaseLicenseWithPoolByNodeIDFunc != nil {
		return m.releaseLicenseWithPoolByNodeIDFunc(ctx, params)
	}
	return nil
}

func (m *mockQueries) ReleaseLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) error {
	if m.releaseLicenseWithoutPoolByNodeIDFunc != nil {
		return m.releaseLicenseWithoutPoolByNodeIDFunc(ctx, nodeID)
	}
	return nil
}

func (m *mockQueries) ClaimLicenseWithPoolFIFO(ctx context.Context, params ClaimLicenseWithPoolFIFOParams) (License, error) {
	if m.claimLicenseWithPoolFIFOFunc != nil {
		return m.claimLicenseWithPoolFIFOFunc(ctx, params)
	}
	return License{}, nil
}

func (m *mockQueries) ClaimLicenseWithoutPoolFIFO(ctx context.Context, nodeID *int64) (License, error) {
	if m.claimLicenseWithoutPoolFIFOFunc != nil {
		return m.claimLicenseWithoutPoolFIFOFunc(ctx, nodeID)
	}
	return License{}, nil
}

func (m *mockQueries) GetLicenseWithPoolByNodeID(ctx context.Context, params GetLicenseWithPoolByNodeIDParams) (License, error) {
	if m.getLicenseWithPoolByNodeIDFunc != nil {
		return m.getLicenseWithPoolByNodeIDFunc(ctx, params)
	}
	return License{}, nil
}

func (m *mockQueries) GetLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) (License, error) {
	if m.getLicenseWithoutPoolByNodeIDFunc != nil {
		return m.getLicenseWithoutPoolByNodeIDFunc(ctx, nodeID)
	}
	return License{}, nil
}

func (m *mockQueries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{}
}

// Add other required methods to satisfy the interface
func (m *mockQueries) InsertLicense(ctx context.Context, params InsertLicenseParams) (License, error) {
	return License{}, nil
}

func (m *mockQueries) DeleteLicenseByGUID(ctx context.Context, guid string) (License, error) {
	return License{}, nil
}

func (m *mockQueries) ActivateNode(ctx context.Context, fingerprint string) (Node, error) {
	return Node{}, nil
}

func (m *mockQueries) DeactivateNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return nil
}

func (m *mockQueries) GetNodeByFingerprint(ctx context.Context, fingerprint string) (Node, error) {
	return Node{}, nil
}

func (m *mockQueries) PingNodeHeartbeatByFingerprint(ctx context.Context, fingerprint string) error {
	return nil
}

func (m *mockQueries) CreatePool(ctx context.Context, name string) (Pool, error) {
	return Pool{}, nil
}

func (m *mockQueries) GetPoolByID(ctx context.Context, id int64) (Pool, error) {
	return Pool{}, nil
}

func (m *mockQueries) GetPoolByName(ctx context.Context, name string) (Pool, error) {
	return Pool{}, nil
}

func (m *mockQueries) DeletePoolByID(ctx context.Context, id int64) (Pool, error) {
	return Pool{}, nil
}

func (m *mockQueries) InsertAuditLog(ctx context.Context, params InsertAuditLogParams) error {
	return nil
}

func (m *mockQueries) ClaimLicenseWithPoolLIFO(ctx context.Context, params ClaimLicenseWithPoolLIFOParams) (License, error) {
	return License{}, nil
}

func (m *mockQueries) ClaimLicenseWithPoolRandom(ctx context.Context, params ClaimLicenseWithPoolRandomParams) (License, error) {
	return License{}, nil
}

func (m *mockQueries) ClaimLicenseWithoutPoolLIFO(ctx context.Context, nodeID *int64) (License, error) {
	return License{}, nil
}

func (m *mockQueries) ClaimLicenseWithoutPoolRandom(ctx context.Context, nodeID *int64) (License, error) {
	return License{}, nil
}

func (m *mockQueries) ReleaseLicensesFromDeadNodes(ctx context.Context, ttl string) ([]License, error) {
	return []License{}, nil
}

func (m *mockQueries) DeactivateDeadNodes(ctx context.Context, ttl string) ([]Node, error) {
	return []Node{}, nil
}

// mockStore is a testable version of Store that uses the interface
type mockStore struct {
	queries    QuerierInterface
	connection *sql.DB
}

func (s *mockStore) GetLicenses(ctx context.Context, predicates ...LicensePredicateFunc) ([]License, error) {
	preds := applyLicensePredicates(predicates...)

	if preds.pool != nil {
		if preds.pool == AnyPool {
			return s.queries.GetLicenses(ctx)
		}

		return s.queries.GetLicensesWithPool(ctx, &preds.pool.ID)
	}

	return s.queries.GetLicensesWithoutPool(ctx)
}

func (s *mockStore) GetLicenseByGUID(ctx context.Context, id string, predicates ...LicensePredicateFunc) (*License, error) {
	preds := applyLicensePredicates(predicates...)

	var license License
	var err error

	if preds.pool != nil {
		if preds.pool == AnyPool {
			license, err = s.queries.GetLicenseByGUID(ctx, id)
		} else {
			license, err = s.queries.GetLicenseWithPoolByGUID(ctx, GetLicenseWithPoolByGUIDParams{id, &preds.pool.ID})
		}
	} else {
		license, err = s.queries.GetLicenseWithoutPoolByGUID(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *mockStore) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) error {
	preds := applyLicensePredicates(predicates...)

	if preds.pool != nil {
		if preds.pool == AnyPool {
			return ErrAnyPoolNotSupported
		}

		return s.queries.ReleaseLicenseWithPoolByNodeID(ctx, ReleaseLicenseWithPoolByNodeIDParams{nodeID, &preds.pool.ID})
	}

	return s.queries.ReleaseLicenseWithoutPoolByNodeID(ctx, nodeID)
}

func (s *mockStore) ClaimLicenseByStrategy(ctx context.Context, strategy string, nodeID *int64, predicates ...LicensePredicateFunc) (*License, error) {
	preds := applyLicensePredicates(predicates...)

	var license License
	var err error

	if preds.pool != nil {
		if preds.pool == AnyPool {
			return nil, ErrAnyPoolNotSupported
		}

		switch strategy {
		case "fifo":
			license, err = s.queries.ClaimLicenseWithPoolFIFO(ctx, ClaimLicenseWithPoolFIFOParams{nodeID, &preds.pool.ID})
		default:
			license, err = s.queries.ClaimLicenseWithPoolFIFO(ctx, ClaimLicenseWithPoolFIFOParams{nodeID, &preds.pool.ID})
		}
	} else {
		switch strategy {
		case "fifo":
			license, err = s.queries.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
		default:
			license, err = s.queries.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
		}
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *mockStore) GetLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) (*License, error) {
	preds := applyLicensePredicates(predicates...)

	var license License
	var err error

	if preds.pool != nil {
		if preds.pool == AnyPool {
			return nil, ErrAnyPoolNotSupported
		}

		license, err = s.queries.GetLicenseWithPoolByNodeID(ctx, GetLicenseWithPoolByNodeIDParams{nodeID, &preds.pool.ID})
	} else {
		license, err = s.queries.GetLicenseWithoutPoolByNodeID(ctx, nodeID)
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func newMockStore() *mockStore {
	return &mockStore{
		queries:    &mockQueries{},
		connection: nil,
	}
}

func TestStore_GetLicenses(t *testing.T) {
	ctx := context.Background()
	testPool := &Pool{ID: 1, Name: "test"}
	
	tests := []struct {
		name        string
		predicates  []LicensePredicateFunc
		setupMock   func(*mockQueries)
		expectError bool
	}{
		{
			name:       "no predicates - defaults to AnyPool",
			predicates: []LicensePredicateFunc{},
			setupMock: func(m *mockQueries) {
				m.getLicensesFunc = func(ctx context.Context) ([]License, error) {
					return []License{{ID: 1}, {ID: 2}}, nil
				}
			},
		},
		{
			name:       "WithAnyPool predicate",
			predicates: []LicensePredicateFunc{WithAnyPool()},
			setupMock: func(m *mockQueries) {
				m.getLicensesFunc = func(ctx context.Context) ([]License, error) {
					return []License{{ID: 1}, {ID: 2}}, nil
				}
			},
		},
		{
			name:       "WithPool predicate",
			predicates: []LicensePredicateFunc{WithPool(testPool)},
			setupMock: func(m *mockQueries) {
				m.getLicensesWithPoolFunc = func(ctx context.Context, poolID *int64) ([]License, error) {
					if *poolID != testPool.ID {
						t.Errorf("expected poolID %d, got %d", testPool.ID, *poolID)
					}
					return []License{{ID: 1, PoolID: poolID}}, nil
				}
			},
		},
		{
			name:       "WithoutPool predicate",
			predicates: []LicensePredicateFunc{WithoutPool()},
			setupMock: func(m *mockQueries) {
				m.getLicensesWithoutPoolFunc = func(ctx context.Context) ([]License, error) {
					return []License{{ID: 1, PoolID: nil}}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockStore()
			mock := store.queries.(*mockQueries)
			tt.setupMock(mock)

			licenses, err := store.GetLicenses(ctx, tt.predicates...)
			
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && len(licenses) == 0 {
				t.Error("expected licenses but got none")
			}
		})
	}
}

func TestStore_GetLicenseByGUID(t *testing.T) {
	ctx := context.Background()
	testGUID := "test-guid-123"
	testPool := &Pool{ID: 1, Name: "test"}
	
	tests := []struct {
		name        string
		predicates  []LicensePredicateFunc
		setupMock   func(*mockQueries)
		expectError bool
	}{
		{
			name:       "no predicates - defaults to AnyPool",
			predicates: []LicensePredicateFunc{},
			setupMock: func(m *mockQueries) {
				m.getLicenseByGUIDFunc = func(ctx context.Context, guid string) (License, error) {
					if guid != testGUID {
						t.Errorf("expected guid %s, got %s", testGUID, guid)
					}
					return License{ID: 1, Guid: testGUID}, nil
				}
			},
		},
		{
			name:       "WithAnyPool predicate",
			predicates: []LicensePredicateFunc{WithAnyPool()},
			setupMock: func(m *mockQueries) {
				m.getLicenseByGUIDFunc = func(ctx context.Context, guid string) (License, error) {
					return License{ID: 1, Guid: testGUID}, nil
				}
			},
		},
		{
			name:       "WithPool predicate",
			predicates: []LicensePredicateFunc{WithPool(testPool)},
			setupMock: func(m *mockQueries) {
				m.getLicenseWithPoolByGUIDFunc = func(ctx context.Context, params GetLicenseWithPoolByGUIDParams) (License, error) {
					if params.Guid != testGUID {
						t.Errorf("expected guid %s, got %s", testGUID, params.Guid)
					}
					if *params.PoolID != testPool.ID {
						t.Errorf("expected poolID %d, got %d", testPool.ID, *params.PoolID)
					}
					return License{ID: 1, Guid: testGUID, PoolID: &testPool.ID}, nil
				}
			},
		},
		{
			name:       "WithoutPool predicate",
			predicates: []LicensePredicateFunc{WithoutPool()},
			setupMock: func(m *mockQueries) {
				m.getLicenseWithoutPoolByGUIDFunc = func(ctx context.Context, guid string) (License, error) {
					if guid != testGUID {
						t.Errorf("expected guid %s, got %s", testGUID, guid)
					}
					return License{ID: 1, Guid: testGUID, PoolID: nil}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockStore()
			mock := store.queries.(*mockQueries)
			tt.setupMock(mock)

			license, err := store.GetLicenseByGUID(ctx, testGUID, tt.predicates...)
			
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && license == nil {
				t.Error("expected license but got nil")
			}
			if !tt.expectError && license.Guid != testGUID {
				t.Errorf("expected license GUID %s, got %s", testGUID, license.Guid)
			}
		})
	}
}

func TestStore_ReleaseLicenseByNodeID_RejectsAnyPool(t *testing.T) {
	ctx := context.Background()
	testNodeID := int64(123)
	
	tests := []struct {
		name       string
		predicates []LicensePredicateFunc
		expectErr  error
	}{
		{
			name:       "WithAnyPool should return error",
			predicates: []LicensePredicateFunc{WithAnyPool()},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "no predicates defaults to AnyPool - should return error",
			predicates: []LicensePredicateFunc{},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "WithPool should work",
			predicates: []LicensePredicateFunc{WithPool(&Pool{ID: 1, Name: "test"})},
			expectErr:  nil,
		},
		{
			name:       "WithoutPool should work",
			predicates: []LicensePredicateFunc{WithoutPool()},
			expectErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockStore()
			mock := store.queries.(*mockQueries)
			
			mock.releaseLicenseWithPoolByNodeIDFunc = func(ctx context.Context, params ReleaseLicenseWithPoolByNodeIDParams) error {
				return nil
			}
			mock.releaseLicenseWithoutPoolByNodeIDFunc = func(ctx context.Context, nodeID *int64) error {
				return nil
			}

			err := store.ReleaseLicenseByNodeID(ctx, &testNodeID, tt.predicates...)
			
			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("expected error %v but got none", tt.expectErr)
				} else if err != tt.expectErr {
					t.Errorf("expected error %v, got %v", tt.expectErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestStore_ClaimLicenseByStrategy_RejectsAnyPool(t *testing.T) {
	ctx := context.Background()
	testNodeID := int64(123)
	
	tests := []struct {
		name       string
		predicates []LicensePredicateFunc
		expectErr  error
	}{
		{
			name:       "WithAnyPool should return error",
			predicates: []LicensePredicateFunc{WithAnyPool()},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "no predicates defaults to AnyPool - should return error",
			predicates: []LicensePredicateFunc{},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "WithPool should work",
			predicates: []LicensePredicateFunc{WithPool(&Pool{ID: 1, Name: "test"})},
			expectErr:  nil,
		},
		{
			name:       "WithoutPool should work",
			predicates: []LicensePredicateFunc{WithoutPool()},
			expectErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockStore()
			mock := store.queries.(*mockQueries)
			
			mock.claimLicenseWithPoolFIFOFunc = func(ctx context.Context, params ClaimLicenseWithPoolFIFOParams) (License, error) {
				return License{ID: 1}, nil
			}
			mock.claimLicenseWithoutPoolFIFOFunc = func(ctx context.Context, nodeID *int64) (License, error) {
				return License{ID: 1}, nil
			}

			_, err := store.ClaimLicenseByStrategy(ctx, "fifo", &testNodeID, tt.predicates...)
			
			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("expected error %v but got none", tt.expectErr)
				} else if err != tt.expectErr {
					t.Errorf("expected error %v, got %v", tt.expectErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestStore_GetLicenseByNodeID_RejectsAnyPool(t *testing.T) {
	ctx := context.Background()
	testNodeID := int64(123)
	
	tests := []struct {
		name       string
		predicates []LicensePredicateFunc
		expectErr  error
	}{
		{
			name:       "WithAnyPool should return error",
			predicates: []LicensePredicateFunc{WithAnyPool()},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "no predicates defaults to AnyPool - should return error",
			predicates: []LicensePredicateFunc{},
			expectErr:  ErrAnyPoolNotSupported,
		},
		{
			name:       "WithPool should work",
			predicates: []LicensePredicateFunc{WithPool(&Pool{ID: 1, Name: "test"})},
			expectErr:  nil,
		},
		{
			name:       "WithoutPool should work",
			predicates: []LicensePredicateFunc{WithoutPool()},
			expectErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockStore()
			mock := store.queries.(*mockQueries)
			
			mock.getLicenseWithPoolByNodeIDFunc = func(ctx context.Context, params GetLicenseWithPoolByNodeIDParams) (License, error) {
				return License{ID: 1}, nil
			}
			mock.getLicenseWithoutPoolByNodeIDFunc = func(ctx context.Context, nodeID *int64) (License, error) {
				return License{ID: 1}, nil
			}

			_, err := store.GetLicenseByNodeID(ctx, &testNodeID, tt.predicates...)
			
			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("expected error %v but got none", tt.expectErr)
				} else if err != tt.expectErr {
					t.Errorf("expected error %v, got %v", tt.expectErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}