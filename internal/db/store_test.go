package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	schema "github.com/keygen-sh/keygen-relay/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FIXME(ezekg) dup of internal/testutils/memory_store to prevent import cycle
func newMemoryStore(t *testing.T) (*Store, *sql.DB) {
	conn, err := sql.Open("sqlite3", ":memory:?_pragma=foreign_keys(on)")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	_, err = conn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	migrations, err := iofs.New(schema.Migrations, "migrations")
	if err != nil {
		t.Fatalf("failed to initialize migrations fs: %v", err)
	}

	migrator, err := NewMigrator(conn, migrations)
	if err != nil {
		t.Fatalf("failed to initialize migrations: %v", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	store := NewStore(New(conn), conn)

	return store, conn
}

func closeMemoryStore(conn *sql.DB) {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close in-memory database connection: %v", err)
	}
}

func TestStore_GetLicenses(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	// create test pool
	testPool, err := store.CreatePool(ctx, "test-pool")
	require.NoError(t, err)

	// insert test licenses (some with pool, some without)
	pooledLicense, err := store.InsertLicense(ctx, testPool, "pooled-guid", []byte("pooled-file"), "pooled-key")
	require.NoError(t, err)

	unpooledLicense, err := store.InsertLicense(ctx, nil, "unpooled-guid", []byte("unpooled-file"), "unpooled-key")
	require.NoError(t, err)

	t.Run("without any predicates", func(t *testing.T) {
		licenses, err := store.GetLicenses(ctx)
		require.NoError(t, err)
		assert.Len(t, licenses, 2)
	})

	t.Run("with any pool predicate", func(t *testing.T) {
		licenses, err := store.GetLicenses(ctx, WithAnyPool())
		require.NoError(t, err)
		assert.Len(t, licenses, 2)
	})

	t.Run("with named pool predicate", func(t *testing.T) {
		licenses, err := store.GetLicenses(ctx, WithPool(testPool))
		require.NoError(t, err)
		assert.Len(t, licenses, 1)
		assert.Equal(t, pooledLicense.ID, licenses[0].ID)
		assert.Equal(t, &testPool.ID, licenses[0].PoolID)
	})

	t.Run("with nil pool predicate", func(t *testing.T) {
		licenses, err := store.GetLicenses(ctx, WithPool(nil))
		require.NoError(t, err)
		assert.Len(t, licenses, 1)
		assert.Equal(t, unpooledLicense.ID, licenses[0].ID)
		assert.Nil(t, licenses[0].PoolID)
	})

	t.Run("without pool predicate", func(t *testing.T) {
		licenses, err := store.GetLicenses(ctx, WithoutPool())
		require.NoError(t, err)
		assert.Len(t, licenses, 1)
		assert.Equal(t, unpooledLicense.ID, licenses[0].ID)
		assert.Nil(t, licenses[0].PoolID)
	})
}

func TestStore_GetLicenseByGUID(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	testPool, err := store.CreatePool(ctx, "test-pool")
	require.NoError(t, err)

	// create licenses
	pooledLicense, err := store.InsertLicense(ctx, testPool, "pooled-guid", []byte("pooled-file"), "pooled-key")
	require.NoError(t, err)

	unpooledLicense, err := store.InsertLicense(ctx, nil, "unpooled-guid", []byte("unpooled-file"), "unpooled-key")
	require.NoError(t, err)

	t.Run("without any predicates", func(t *testing.T) {
		license, err := store.GetLicenseByGUID(ctx, pooledLicense.Guid)
		require.NoError(t, err)
		assert.Equal(t, pooledLicense.ID, license.ID)
		assert.Equal(t, pooledLicense.Guid, license.Guid)

		license, err = store.GetLicenseByGUID(ctx, unpooledLicense.Guid)
		require.NoError(t, err)
		assert.Equal(t, unpooledLicense.ID, license.ID)
		assert.Equal(t, unpooledLicense.Guid, license.Guid)
	})

	t.Run("with any pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByGUID(ctx, pooledLicense.Guid, WithAnyPool())
		require.NoError(t, err)
		assert.Equal(t, pooledLicense.ID, license.ID)
	})

	t.Run("with named pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByGUID(ctx, pooledLicense.Guid, WithPool(testPool))
		require.NoError(t, err)
		assert.Equal(t, pooledLicense.ID, license.ID)
		assert.Equal(t, &testPool.ID, license.PoolID)

		// should not find unpooled license with pool predicate
		_, err = store.GetLicenseByGUID(ctx, unpooledLicense.Guid, WithPool(testPool))
		assert.Error(t, err)
	})

	t.Run("with nil pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByGUID(ctx, unpooledLicense.Guid, WithPool(nil))
		require.NoError(t, err)
		assert.Equal(t, unpooledLicense.ID, license.ID)
		assert.Nil(t, license.PoolID)

		// should not find pooled license with WithoutPool predicate
		_, err = store.GetLicenseByGUID(ctx, pooledLicense.Guid, WithoutPool())
		assert.Error(t, err)
	})

	t.Run("without pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByGUID(ctx, unpooledLicense.Guid, WithoutPool())
		require.NoError(t, err)
		assert.Equal(t, unpooledLicense.ID, license.ID)
		assert.Nil(t, license.PoolID)

		// should not find pooled license with WithoutPool predicate
		_, err = store.GetLicenseByGUID(ctx, pooledLicense.Guid, WithoutPool())
		assert.Error(t, err)
	})

	t.Run("license not found", func(t *testing.T) {
		_, err := store.GetLicenseByGUID(ctx, "nonexistent-guid")
		assert.Error(t, err)
	})
}

func TestStore_ReleaseLicenseByNodeID(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	testPool, err := store.CreatePool(ctx, "test-pool")
	require.NoError(t, err)

	// create separate nodes since each node can only have one license
	pooledNode, err := store.ActivateNode(ctx, "pooled-test-fingerprint")
	require.NoError(t, err)

	unpooledNode, err := store.ActivateNode(ctx, "unpooled-test-fingerprint")
	require.NoError(t, err)

	// create licenses and claim them
	pooledLicense, err := store.InsertLicense(ctx, testPool, "pooled-guid", []byte("pooled-file"), "pooled-key")
	require.NoError(t, err)

	unpooledLicense, err := store.InsertLicense(ctx, nil, "unpooled-guid", []byte("unpooled-file"), "unpooled-key")
	require.NoError(t, err)

	// claim licenses for specific nodes
	_, err = conn.ExecContext(ctx, "UPDATE licenses SET node_id = ?, last_claimed_at = strftime('%s', 'now') WHERE id = ?", pooledNode.ID, pooledLicense.ID)
	require.NoError(t, err)

	_, err = conn.ExecContext(ctx, "UPDATE licenses SET node_id = ?, last_claimed_at = strftime('%s', 'now') WHERE id = ?", unpooledNode.ID, unpooledLicense.ID)
	require.NoError(t, err)

	t.Run("with any pool predicate", func(t *testing.T) {
		err := store.ReleaseLicenseByNodeID(ctx, &pooledNode.ID, WithAnyPool())
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("without any predicates", func(t *testing.T) {
		err := store.ReleaseLicenseByNodeID(ctx, &pooledNode.ID)
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("with named pool predicate", func(t *testing.T) {
		err := store.ReleaseLicenseByNodeID(ctx, &pooledNode.ID, WithPool(testPool))
		require.NoError(t, err)

		// verify the pooled license was released
		license, err := store.GetLicenseByGUID(ctx, pooledLicense.Guid, WithPool(testPool))
		require.NoError(t, err)
		assert.Nil(t, license.NodeID)
		assert.NotNil(t, license.LastReleasedAt)
	})

	t.Run("without pool predicate", func(t *testing.T) {
		err := store.ReleaseLicenseByNodeID(ctx, &unpooledNode.ID, WithoutPool())
		require.NoError(t, err)

		// verify the unpooled license was released
		license, err := store.GetLicenseByGUID(ctx, unpooledLicense.Guid, WithoutPool())
		require.NoError(t, err)
		assert.Nil(t, license.NodeID)
		assert.NotNil(t, license.LastReleasedAt)
	})
}

func TestStore_ClaimLicenseByStrategy(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	testPool, err := store.CreatePool(ctx, "test-pool")
	require.NoError(t, err)

	node, err := store.ActivateNode(ctx, "test-fingerprint")
	require.NoError(t, err)

	// create available licenses
	_, err = store.InsertLicense(ctx, testPool, "pooled-guid", []byte("pooled-file"), "pooled-key")
	require.NoError(t, err)

	_, err = store.InsertLicense(ctx, nil, "unpooled-guid", []byte("unpooled-file"), "unpooled-key")
	require.NoError(t, err)

	t.Run("with any pool predicate", func(t *testing.T) {
		_, err := store.ClaimLicenseByStrategy(ctx, "fifo", &node.ID, WithAnyPool())
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("without any predicates", func(t *testing.T) {
		_, err := store.ClaimLicenseByStrategy(ctx, "fifo", &node.ID)
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("with named pool predicate", func(t *testing.T) {
		license, err := store.ClaimLicenseByStrategy(ctx, "fifo", &node.ID, WithPool(testPool))
		require.NoError(t, err)
		assert.NotNil(t, license)
		assert.Equal(t, "pooled-guid", license.Guid)
		assert.Equal(t, &node.ID, license.NodeID)
		assert.NotNil(t, license.LastClaimedAt)
	})

	t.Run("without pool predicate", func(t *testing.T) {
		// create another node for this test
		node2, err := store.ActivateNode(ctx, "test-fingerprint-2")
		require.NoError(t, err)

		license, err := store.ClaimLicenseByStrategy(ctx, "fifo", &node2.ID, WithoutPool())
		require.NoError(t, err)
		assert.NotNil(t, license)
		assert.Equal(t, "unpooled-guid", license.Guid)
		assert.Equal(t, &node2.ID, license.NodeID)
		assert.NotNil(t, license.LastClaimedAt)
	})

	t.Run("strategies", func(t *testing.T) {
		// create more licenses and test different strategies
		node3, err := store.ActivateNode(ctx, "test-fingerprint-3")
		require.NoError(t, err)

		for i := range 5 {
			_, err = store.InsertLicense(ctx, testPool, fmt.Sprintf("strategy-test-%d", i), []byte(fmt.Sprintf("file-%d", i)), fmt.Sprintf("key-%d", i))
			require.NoError(t, err)
		}

		// test FIFO
		license, err := store.ClaimLicenseByStrategy(ctx, "fifo", &node3.ID, WithPool(testPool))
		require.NoError(t, err)
		assert.NotNil(t, license)

		// test LIFO
		node4, err := store.ActivateNode(ctx, "test-fingerprint-4")
		require.NoError(t, err)

		license, err = store.ClaimLicenseByStrategy(ctx, "lifo", &node4.ID, WithPool(testPool))
		require.NoError(t, err)
		assert.NotNil(t, license)

		// test random
		node5, err := store.ActivateNode(ctx, "test-fingerprint-5")
		require.NoError(t, err)

		license, err = store.ClaimLicenseByStrategy(ctx, "rand", &node5.ID, WithPool(testPool))
		require.NoError(t, err)
		assert.NotNil(t, license)

		// invalid
		node6, err := store.ActivateNode(ctx, "test-fingerprint-6")
		require.NoError(t, err)

		license, err = store.ClaimLicenseByStrategy(ctx, "invalid-strategy", &node6.ID, WithPool(testPool))
		require.Error(t, err)
		assert.Nil(t, license)
	})
}

func TestStore_GetLicenseByNodeID(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	testPool, err := store.CreatePool(ctx, "test-pool")
	require.NoError(t, err)

	// create separate nodes since each node can only have one license
	pooledNode, err := store.ActivateNode(ctx, "pooled-node-fingerprint")
	require.NoError(t, err)

	unpooledNode, err := store.ActivateNode(ctx, "unpooled-node-fingerprint")
	require.NoError(t, err)

	// create and claim licenses
	pooledLicense, err := store.InsertLicense(ctx, testPool, "pooled-guid", []byte("pooled-file"), "pooled-key")
	require.NoError(t, err)

	unpooledLicense, err := store.InsertLicense(ctx, nil, "unpooled-guid", []byte("unpooled-file"), "unpooled-key")
	require.NoError(t, err)

	// claim licenses for specific nodes
	_, err = conn.ExecContext(ctx, "UPDATE licenses SET node_id = ?, last_claimed_at = strftime('%s', 'now') WHERE id = ?", pooledNode.ID, pooledLicense.ID)
	require.NoError(t, err)

	_, err = conn.ExecContext(ctx, "UPDATE licenses SET node_id = ?, last_claimed_at = strftime('%s', 'now') WHERE id = ?", unpooledNode.ID, unpooledLicense.ID)
	require.NoError(t, err)

	t.Run("with any pool predicate", func(t *testing.T) {
		_, err := store.GetLicenseByNodeID(ctx, &pooledNode.ID, WithAnyPool())
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("without any predicates", func(t *testing.T) {
		_, err := store.GetLicenseByNodeID(ctx, &pooledNode.ID)
		assert.ErrorIs(t, err, ErrAnyPoolNotSupported)
	})

	t.Run("with named pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByNodeID(ctx, &pooledNode.ID, WithPool(testPool))
		require.NoError(t, err)
		assert.NotNil(t, license)
		assert.Equal(t, pooledLicense.ID, license.ID)
		assert.Equal(t, "pooled-guid", license.Guid)
		assert.Equal(t, &pooledNode.ID, license.NodeID)
	})

	t.Run("without pool predicate", func(t *testing.T) {
		license, err := store.GetLicenseByNodeID(ctx, &unpooledNode.ID, WithoutPool())
		require.NoError(t, err)
		assert.NotNil(t, license)
		assert.Equal(t, unpooledLicense.ID, license.ID)
		assert.Equal(t, "unpooled-guid", license.Guid)
		assert.Equal(t, &unpooledNode.ID, license.NodeID)
	})

	t.Run("node has no license in specific pool", func(t *testing.T) {
		emptyPool, err := store.CreatePool(ctx, "empty-pool")
		require.NoError(t, err)

		_, err = store.GetLicenseByNodeID(ctx, &pooledNode.ID, WithPool(emptyPool))
		assert.Error(t, err)
	})
}

func TestStore_AdditionalMethods(t *testing.T) {
	store, conn := newMemoryStore(t)
	defer closeMemoryStore(conn)
	ctx := context.Background()

	t.Run("InsertLicense", func(t *testing.T) {
		testPool, err := store.CreatePool(ctx, "insert-test-pool")
		require.NoError(t, err)

		// test with pool
		license, err := store.InsertLicense(ctx, testPool, "insert-test-guid", []byte("test-file"), "test-key")
		require.NoError(t, err)
		assert.Equal(t, "insert-test-guid", license.Guid)
		assert.Equal(t, "test-key", license.Key)
		assert.Equal(t, &testPool.ID, license.PoolID)

		// test without pool
		license2, err := store.InsertLicense(ctx, nil, "insert-test-guid-2", []byte("test-file-2"), "test-key-2")
		require.NoError(t, err)
		assert.Equal(t, "insert-test-guid-2", license2.Guid)
		assert.Nil(t, license2.PoolID)
	})

	t.Run("DeleteLicenseByGUID", func(t *testing.T) {
		// insert a license to delete
		license, err := store.InsertLicense(ctx, nil, "delete-test-guid", []byte("delete-file"), "delete-key")
		require.NoError(t, err)

		// delete it
		deleted, err := store.DeleteLicenseByGUID(ctx, license.Guid)
		require.NoError(t, err)
		assert.Equal(t, license.ID, deleted.ID)

		// verify it's gone
		_, err = store.GetLicenseByGUID(ctx, license.Guid)
		assert.Error(t, err)
	})

	t.Run("Nodes", func(t *testing.T) {
		// activate node
		node, err := store.ActivateNode(ctx, "node-test-fingerprint")
		require.NoError(t, err)
		assert.Equal(t, "node-test-fingerprint", node.Fingerprint)

		// get node
		retrievedNode, err := store.GetNodeByFingerprint(ctx, "node-test-fingerprint")
		require.NoError(t, err)
		assert.Equal(t, node.ID, retrievedNode.ID)

		// ping heartbeat
		err = store.PingNodeHeartbeatByFingerprint(ctx, "node-test-fingerprint")
		require.NoError(t, err)

		// deactivate node
		err = store.DeactivateNodeByFingerprint(ctx, "node-test-fingerprint")
		require.NoError(t, err)

		// should not be able to get it now
		_, err = store.GetNodeByFingerprint(ctx, "node-test-fingerprint")
		assert.Error(t, err)
	})

	t.Run("Pools", func(t *testing.T) {
		// create pool
		pool, err := store.CreatePool(ctx, "pool-operations-test")
		require.NoError(t, err)
		assert.Equal(t, "pool-operations-test", pool.Name)

		// get by ID
		retrievedPool, err := store.GetPoolByID(ctx, pool.ID)
		require.NoError(t, err)
		assert.Equal(t, pool.Name, retrievedPool.Name)

		// get by name
		retrievedPool2, err := store.GetPoolByName(ctx, "pool-operations-test")
		require.NoError(t, err)
		assert.Equal(t, pool.ID, retrievedPool2.ID)

		// get all pools
		pools, err := store.GetPools(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(pools), 0)

		// delete pool
		deleted, err := store.DeletePoolByID(ctx, pool.ID)
		require.NoError(t, err)
		assert.Equal(t, pool.ID, deleted.ID)

		// should not be able to get it now
		_, err = store.GetPoolByID(ctx, pool.ID)
		assert.Error(t, err)
	})
}
