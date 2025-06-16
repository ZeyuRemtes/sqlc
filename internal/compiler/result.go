package compiler

import (
	"github.com/ZeyuRemtes/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog *catalog.Catalog
	Queries []*Query
}
