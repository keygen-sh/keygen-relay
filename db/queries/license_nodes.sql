-- name: GetLicenseNodeCount :one
SELECT COUNT(*) FROM license_nodes WHERE license_id = ?;
