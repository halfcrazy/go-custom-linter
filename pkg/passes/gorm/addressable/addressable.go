package addressable

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

type Kind int

const (
	KindNone Kind = iota
	KindDest
	KindValue
)

var Analyzer = &analysis.Analyzer{
	Name:       "gormaddressable",
	Doc:        Doc,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	Run:        run,
	ResultType: reflect.TypeOf((*Result)(nil)),
	FactTypes:  []analysis.Fact{},
}

const Doc = `检查gorm的传入参数是否为可寻址`

// TODO: recognize wrapper func, enhance accuracy
var needAddressableDest = map[string]bool{
	"(*gorm.io/gorm.DB).First":         true,
	"(*gorm.io/gorm.DB).Take":          true,
	"(*gorm.io/gorm.DB).Last":          true,
	"(*gorm.io/gorm.DB).Find":          true,
	"(*gorm.io/gorm.DB).FindInBatches": true,
	"(*gorm.io/gorm.DB).FirstOrInit":   true,
	"(*gorm.io/gorm.DB).FirstOrCreate": true,
	"(*gorm.io/gorm.DB).Scan":          true,
	"(*gorm.io/gorm.DB).Pluck":         true,
	"(*gorm.io/gorm.DB).ScanRows":      true,
}

var needAddressableValue = map[string]bool{
	"(*gorm.io/gorm.DB).Create": true,
}

// Result is the printf analyzer's result type. Clients may query the result
// to learn whether a function behaves like fmt.Print or fmt.Printf.
type Result struct {
	funcs map[*types.Func]string
}

// Kind reports whether fn behaves like fmt.Print or fmt.Printf.
func (r *Result) Kind(fn *types.Func) string {
	return fn.FullName()
}

func run(pass *analysis.Pass) (interface{}, error) {
	res := &Result{
		funcs: make(map[*types.Func]string),
	}
	checkCall(pass)
	return res, nil
}

// checkCall triggers the print-specific checks if the call invokes a print function.
func checkCall(pass *analysis.Pass) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		fn, kind := getNeedAddressableFnKind(pass, call)
		if kind != KindNone {
			checkParamAddressable(pass, call, fn, kind)
		}
	})
}

func getNeedAddressableFnKind(pass *analysis.Pass, call *ast.CallExpr) (fn *types.Func, kind Kind) {
	fn, _ = typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if fn == nil {
		return nil, KindNone
	}
	if _, ok := needAddressableDest[fn.FullName()]; ok {
		return fn, KindDest
	}
	if _, ok := needAddressableValue[fn.FullName()]; ok {
		return fn, KindValue
	}
	return fn, KindNone
}

func checkParamAddressable(pass *analysis.Pass, call *ast.CallExpr, fn *types.Func, kind Kind) {
	// just report runtime potential error, discard compile time error
	// recognize dest interface{} or value interface{} arg then check if it's addressable
	argsNum := len(call.Args)
	sig := fn.Type().(*types.Signature)
	if sig.Params() != nil {
		for idx := 0; idx < argsNum; idx++ {
			param := sig.Params().At(idx)
			if param.Type().String() == "interface{}" &&
				((kind == KindValue && param.Name() == "value") || (kind == KindDest && param.Name() == "dest")) {
				typ := pass.TypesInfo.Types[call.Args[idx]].Type
				for typ.String() != typ.Underlying().String() {
					typ = typ.Underlying()
				}
				// ignore interface
				if typ.String() != "" && !strings.HasPrefix(typ.String(), "interface{") &&
					!strings.HasPrefix(typ.String(), "*") {
					pass.Reportf(call.Lparen, "not addressable param passed to gorm %s", fn.Name())
				}
			}
		}
	}
}
