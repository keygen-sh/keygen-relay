package schema

import (
	"embed"
)

//go:embed schema.sql
var Schema string

//go:embed migrations/*.sql
var Migrations embed.FS
