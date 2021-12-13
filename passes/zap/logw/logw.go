package logw

import (
	"go/ast"
	"go/constant"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
	"reflect"
)

var Analyzer = &analysis.Analyzer{
	Name:       "logw",
	Doc:        Doc,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	Run:        run,
	ResultType: reflect.TypeOf((*Result)(nil)),
	FactTypes:  []analysis.Fact{},
}

const Doc = `检查zap的w系列函数的参数个数是否为msg, kv对的形式`

// TODO: recognize wrapper func, enhance accuracy
var isWStyle = map[string]bool{
	"Debugw":  true,
	"Infow":   true,
	"Warnw":   true,
	"Errorw":  true,
	"DPanicw": true,
	"Panicw":  true,
	"Fatalw":  true,
}

// Result is the printf analyzer's result type. Clients may query the result
// to learn whether a function behaves like fmt.Print or fmt.Printf.
type Result struct {
	funcs map[*types.Func]string
}

// Kind reports whether fn behaves like fmt.Print or fmt.Printf.
func (r *Result) Kind(fn *types.Func) string {
	return fn.Name()
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
		fn, kind := getWstyleFnKind(pass, call)
		if isWStyle[kind] {
			checkPrintw(pass, call, fn)
		}
	})
}

func getWstyleFnKind(pass *analysis.Pass, call *ast.CallExpr) (fn *types.Func, kind string) {
	fn, _ = typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if fn == nil {
		return nil, ""
	}

	_, ok := isWStyle[fn.Name()]
	if ok {
		return fn, fn.Name()
	}
	return fn, ""
}

func checkPrintw(pass *analysis.Pass, call *ast.CallExpr, fn *types.Func) {
	// check 4 cases
	// 1) msg, k, v, k1
	// 2) msg, none-string key, v
	// 3) no args
	// 4) msg, k, v, k, v1
	argsNum := len(call.Args)
	if argsNum < 1 {
		pass.Reportf(call.Lparen, "no args provided call to %s", fn.FullName())
		return
	}
	if (argsNum-1)%2 == 1 {
		pass.Reportf(call.Lparen, "invalid pair provided call to %s", fn.FullName())
	}
	var keySet = make(map[string]bool)
	for idx := 1; idx < argsNum; idx += 2 {
		keyArg := call.Args[idx]
		val := pass.TypesInfo.Types[keyArg].Value
		if val != nil {
			if val.Kind() != constant.String {
				pass.Reportf(call.Lparen, "none string key in pair provided call to %s", fn.FullName())
				return
			} else {
				if keySet[constant.StringVal(val)] {
					pass.Reportf(call.Lparen, "duplicate key %s in pair provided call to %s",
						constant.StringVal(val), fn.FullName())
					return
				} else {
					keySet[constant.StringVal(val)] = true
				}
			}
		}
	}
}
