package golang

import (
	"log"
	"strings"

	"github.com/ZeyuRemtes/sqlc/internal/codegen/sdk"
	"github.com/ZeyuRemtes/sqlc/internal/plugin"
)

func sqliteType(req *plugin.CodeGenRequest, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull || col.IsArray

	switch dt {

	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint", "unsignedbigint", "int2", "int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "blob":
		return "[]byte"

	case "real", "double", "doubleprecision", "float":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "boolean", "bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "date", "datetime", "timestamp":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "any":
		return "interface{}"

	}

	switch {

	case strings.HasPrefix(dt, "character"),
		strings.HasPrefix(dt, "varchar"),
		strings.HasPrefix(dt, "varyingcharacter"),
		strings.HasPrefix(dt, "nchar"),
		strings.HasPrefix(dt, "nativecharacter"),
		strings.HasPrefix(dt, "nvarchar"),
		dt == "text",
		dt == "clob":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case strings.HasPrefix(dt, "decimal"), dt == "numeric":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	default:
		log.Printf("unknown SQLite type: %s\n", dt)
		return "interface{}"

	}
}
