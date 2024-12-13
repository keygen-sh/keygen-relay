// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type AuditLog struct {
	ID           int64
	EventTypeID  int64
	EntityTypeID int64
	EntityID     string
	CreatedAt    int64
}

type EntityType struct {
	ID   int64
	Name string
}

type EventType struct {
	ID   int64
	Name string
}

type License struct {
	ID             string
	File           []byte
	Key            string
	Claims         int64
	LastClaimedAt  *int64
	LastReleasedAt *int64
	NodeID         *int64
	CreatedAt      int64
}

type Node struct {
	ID              int64
	Fingerprint     string
	ClaimedAt       *int64
	LastHeartbeatAt *int64
	CreatedAt       int64
}
