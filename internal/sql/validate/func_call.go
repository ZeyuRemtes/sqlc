package validate

import (
	"errors"
	"fmt"

	"github.com/ZeyuRemtes/sqlc/internal/config"
	"github.com/ZeyuRemtes/sqlc/internal/sql/ast"
	"github.com/ZeyuRemtes/sqlc/internal/sql/astutils"
	"github.com/ZeyuRemtes/sqlc/internal/sql/catalog"
	"github.com/ZeyuRemtes/sqlc/internal/sql/sqlerr"
)

type funcCallVisitor struct {
	catalog  *catalog.Catalog
	settings config.CombinedSettings
	err      error
}

func (v *funcCallVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	fn := call.Func
	if fn == nil {
		return v
	}

	// Custom validation for sqlc.arg, sqlc.narg and sqlc.slice
	// TODO: Replace this once type-checking is implemented
	if fn.Schema == "sqlc" {
		if !(fn.Name == "arg" || fn.Name == "narg" || fn.Name == "slice" || fn.Name == "embed") {
			v.err = sqlerr.FunctionNotFound("sqlc." + fn.Name)
			return nil
		}

		if len(call.Args.Items) != 1 {
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected 1 parameter to sqlc.%s; got %d", fn.Name, len(call.Args.Items)),
				Location: call.Pos(),
			}
			return nil
		}
		switch n := call.Args.Items[0].(type) {
		case *ast.A_Const:
		case *ast.ColumnRef:
		default:
			v.err = &sqlerr.Error{
				Message:  fmt.Sprintf("expected parameter to sqlc.%s to be string or reference; got %T", fn.Name, n),
				Location: call.Pos(),
			}
			return nil
		}

		// If we have sqlc.arg or sqlc.narg, there is no need to resolve the function call.
		// It won't resolve anyway, sinc it is not a real function.
		return nil
	}

	fun, err := v.catalog.ResolveFuncCall(call)
	if fun != nil {
		return v
	}
	if errors.Is(err, sqlerr.NotFound) && !v.settings.Package.StrictFunctionChecks {
		return v
	}
	v.err = err
	return nil
}

func FuncCall(c *catalog.Catalog, cs config.CombinedSettings, n ast.Node) error {
	visitor := funcCallVisitor{catalog: c, settings: cs}
	astutils.Walk(&visitor, n)
	return visitor.err
}
