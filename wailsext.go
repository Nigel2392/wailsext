//go:build js && wasm
// +build js,wasm

package wailsext

import (
	"syscall/js"
)

// Wails functions

var (
	windowRuntime js.Value = js.Global().Get("window").Get("runtime")
	WindowGo      js.Value = js.Global().Get("window").Get("go")
)

type Mapper map[string]interface{}

func NewMapper(arg js.Value) Mapper {
	var m = make(Mapper)
	m.parseArg(arg)
	return m
}

func (m Mapper) parseArg(arg js.Value) {
	var object = js.Global().Get("Object")
	var keys = object.Call("keys", arg)
	for i := 0; i < keys.Length(); i++ {
		var key = keys.Index(i).String()
		// Cast to interface{} to avoid type assertion panic
		m[key] = arg.Get(key)
	}
}

func (m Mapper) Get(key string) interface{} {
	return m[key]
}

func (m Mapper) GetInt(key string) int {
	switch val := m[key].(type) {
	case int:
		return val
	case int64:
		return int(val)
	case int32:
		return int(val)
	case int16:
		return int(val)
	case int8:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	case js.Value:
		return val.Int()
	default:
		panic("unknown type")
	}
}

func (m Mapper) GetFloat(key string) float64 {
	switch val := m[key].(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case js.Value:
		return val.Float()
	default:
		panic("unknown type")
	}
}

func (m Mapper) GetString(key string) string {
	switch val := m[key].(type) {
	case string:
		return val
	case js.Value:
		return val.String()
	default:
		panic("unknown type")
	}
}

func (m Mapper) GetBool(key string) bool {
	switch val := m[key].(type) {
	case bool:
		return val
	case js.Value:
		return val.Bool()
	default:
		panic("unknown type")
	}
}

func (m Mapper) GetArray(key string) []interface{} {
	switch val := m[key].(type) {
	case []interface{}:
		return val
	case js.Value:
		var arr = make([]interface{}, val.Length())
		for i := 0; i < val.Length(); i++ {
			arr[i] = val.Index(i)
		}
		return arr
	default:
		panic("unknown type")
	}
}

func (m Mapper) GetMap(key string) Mapper {
	switch val := m[key].(type) {
	case Mapper:
		return val
	case js.Value:
		return NewMapper(val)
	default:
		panic("unknown type")
	}
}

func GetStructure(pkgName, structName string) js.Value {
	var pkg = WindowGo.Get(pkgName)
	if !pkg.Truthy() {
		panic("package not found: " + pkgName)
	}
	var structure = pkg.Get(structName)
	if !structure.Truthy() {
		panic("structure not found: " + structName)
	}
	return structure
}

func WailsCall(pkgName, structName, funcName string, cb func(this js.Value, args []js.Value) any, args ...any) js.Value {
	var structure = GetStructure(pkgName, structName)
	// Function is a promise, so we need to call the callback when it resolves
	if cb == nil {
		return structure.Call(funcName, args...)
	}
	var function = structure.Get(funcName)
	if !function.Truthy() {
		panic("function not found: " + funcName)
	}
	return function.Invoke(args...).Call("then", js.FuncOf(cb))
}

func MainCall(strctName, funcName string, cb func(this js.Value, args []js.Value) any, args ...any) js.Value {
	return WailsCall("main", strctName, funcName, cb, args...)
}

func EventsOn(eventName string, callback func(this js.Value, args []js.Value) any) {
	windowRuntime.Call("EventsOn", eventName, js.FuncOf(callback))
}

func EventsOff(eventName string) {
	windowRuntime.Call("EventsOff", eventName)
}

func EventsOnce(eventName string, callback func(this js.Value, args []js.Value) any) {
	windowRuntime.Call("EventsOnce", eventName, js.FuncOf(callback))
}

func EventsOnMultiple(eventNames string, maxCnt int, callback func(this js.Value, args []js.Value) any) {
	windowRuntime.Call("EventsOnMultiple", eventNames, js.FuncOf(callback), maxCnt)
}

func EventsEmit(eventName string, data any) {
	windowRuntime.Call("EventEmit", eventName, data)
}

func WindowSetTitle(title string) {
	windowRuntime.Call("WindowSetTitle", title)
}

func WindowFullscreen() {
	windowRuntime.Call("WindowFullscreen")
}

func WindowUnFullScreen() {
	windowRuntime.Call("WindowUnfullscreen")
}

func WindowIsFullScreen(cb func(v bool)) {
	var promise = windowRuntime.Call("WindowIsFullscreen")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Bool())
		return nil
	}))
}

func WindowCenter() {
	windowRuntime.Call("WindowCenter")
}

func WindowReload() {
	windowRuntime.Call("WindowReload")
}

func WindowReloadApp() {
	windowRuntime.Call("WindowReloadApp")
}

func WindowSetSystemDefaultTheme() {
	windowRuntime.Call("WindowSetSystemDefaultTheme")
}

func WindowSetLightTheme() {
	windowRuntime.Call("WindowSetLightTheme")
}

func WindowSetDarkTheme() {
	windowRuntime.Call("WindowSetDarkTheme")
}

func WindowShow() {
	windowRuntime.Call("WindowShow")
}

func WindowHide() {
	windowRuntime.Call("WindowHide")
}

func WindowisNormal(cb func(v bool)) {
	var promise = windowRuntime.Call("WindowIsNormal")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Bool())
		return nil
	}))
}

func WindowGetSize(cb func(width, height int)) {
	var promise = windowRuntime.Call("WindowGetSize")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Get("w").Int(), args[0].Get("h").Int())
		return nil
	}))
}

func WindowSetMinSize(width, height int) {
	windowRuntime.Call("WindowSetMinSize", width, height)
}

func WindowSetMaxSize(width, height int) {
	windowRuntime.Call("WindowSetMaxSize", width, height)
}

func WindowSetAlwaysOnTop(flag bool) {
	windowRuntime.Call("WindowSetAlwaysOnTop", flag)
}

func WindowSetPosition(x, y int) {
	windowRuntime.Call("WindowSetPosition", x, y)
}

func WindowGetPosition(cb func(x, y int)) {
	var promise = windowRuntime.Call("WindowGetPosition")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Get("x").Int(), args[0].Get("y").Int())
		return nil
	}))

}

func WindowMaximise() {
	windowRuntime.Call("WindowMaximise")
}

func WindowUnmaximise() {
	windowRuntime.Call("WindowUnmaximise")
}

func WindowIsMaximised(cb func(v bool)) {
	var promise = windowRuntime.Call("WindowIsMaximised")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Bool())
		return nil
	}))
}

func WindowToggleMaximise() {
	windowRuntime.Call("WindowToggleMaximise")
}

func WindowUnminimise() {
	windowRuntime.Call("WindowUnminimise")
}

func WindowIsMinimised(cb func(v bool)) {
	var promise = windowRuntime.Call("WindowIsMinimised")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(args[0].Bool())
		return nil
	}))
}

func WindowMinimise() {
	windowRuntime.Call("WindowMinimise")
}

func WindowSetbackgroundColor(r, g, b, a int8) {
	windowRuntime.Call("WindowSetbackgroundColor", r, g, b, a)
}

func BrowserOpen(url string) {
	windowRuntime.Call("BrowserOpenURL", url)
}

func ExitApp() {
	windowRuntime.Call("Quit")
}

func Hide() {
	windowRuntime.Call("Hide")
}

func Show() {
	windowRuntime.Call("Show")
}

type EnvInfo struct {
	BuildType string
	Platform  string
	Arch      string
}

func Environment(cb func(v EnvInfo)) {
	var promise = windowRuntime.Call("Environment")
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		cb(EnvInfo{
			BuildType: args[0].Get("buildType").String(),
			Platform:  args[0].Get("platform").String(),
			Arch:      args[0].Get("arch").String(),
		})
		return nil
	}))
}
