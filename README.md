# HandlerX

**HandlerX** is a lightweight, framework-agnostic HTTP handler abstraction for Go.

It allows you to write your business logic once and plug it into different HTTP frameworks like **Fiber**, **Echo**, or even the standard **net/http** — without rewriting your handlers.

---

## ✨ Features

* 🔌 Framework-agnostic handlers
* 🧩 Simple and composable API
* 🔁 Middleware support (`Next()`)
* 📦 Unified response model
* 🧪 Easy to test (mockable context)
* ⚡ Minimal and idiomatic Go

---

## 📦 Installation

```bash
go get github.com/atendi9/handlerx
```

---

## 🚀 Basic Example

```go
func HelloHandler(c handlerx.Context) handlerx.Response {
	name := c.Query("name", "guest")

	return handlerx.Response{
		Data: map[string]string{
			"message": "Hello " + name,
		},
	}
}
```

---

## 🔌 Integrations

### ⚡ Fiber

```go
package main

import (
	"mime/multipart"
	"time"

	"github.com/atendi9/handlerx"
	"github.com/gofiber/fiber/v2"
)

// ===== Context Implementation =====

type FiberContext struct {
	Ctx *fiber.Ctx
}

func (f FiberContext) Headers() map[string][]string {
	return f.Ctx.GetReqHeaders()
}

func (f FiberContext) BodyParser(v any) error {
	return f.Ctx.BodyParser(v)
}

func (f FiberContext) QueryParser(v any) error {
	return f.Ctx.QueryParser(v)
}

func (f FiberContext) ParamsParser(v any) error {
	return nil
}

func (f FiberContext) ReqHeaderParser(v any) error {
	return nil
}

func (f FiberContext) Header(key string) string {
	return f.Ctx.Get(key)
}

func (f FiberContext) Method() string {
	return f.Ctx.Method()
}

func (f FiberContext) IP() string {
	return f.Ctx.IP()
}

func (f FiberContext) IPs() []string {
	return f.Ctx.IPs()
}

func (f FiberContext) Body() []byte {
	return f.Ctx.Body()
}

func (f FiberContext) Query(name string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return f.Ctx.Query(name, defaultValue[0])
	}
	return f.Ctx.Query(name)
}

func (f FiberContext) Params(name string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return f.Ctx.Params(name, defaultValue[0])
	}
	return f.Ctx.Params(name)
}

func (f FiberContext) FormFile(key string) (*multipart.FileHeader, error) {
	return f.Ctx.FormFile(key)
}

func (f FiberContext) SendStatus(status int) error {
	return f.Ctx.SendStatus(status)
}

func (f FiberContext) Send(data []byte) error {
	return f.Ctx.Send(data)
}

func (f FiberContext) JSON(data any) error {
	return f.Ctx.JSON(data)
}

func (f FiberContext) Next() error {
	return f.Ctx.Next()
}

func (f FiberContext) Now() time.Time {
	return time.Now()
}

func (f FiberContext) Path(defaultValue ...string) string {
	return f.Ctx.Path()
}

// ===== Converter =====

type FiberConverter struct{}

func (f FiberConverter) Convert(h handlerx.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := FiberContext{Ctx: c}
		res := h(handlerx.Atendi9Context{Context: ctx})

		if res.GoNext() {
			return c.Next()
		}

		if len(res.FilePath) > 0 {
			return c.SendFile(res.FilePath)
		}

		if err := res.Err; err != nil {
			return c.Status(res.Status()).JSON(fiber.Map{
				"err": err.Error(),
			})
		}

		if v, ok := res.Data.(string); ok {
			return c.Status(res.Status()).SendString(v)
		}

		return c.Status(res.Status()).JSON(res.Data)
	}
}

// ===== Handler =====

func Hello(c handlerx.Context) handlerx.Response {
	return handlerx.Response{
		Data: map[string]string{
			"message": "Hello from Fiber",
		},
	}
}

// ===== Main =====

func main() {
	app := fiber.New()
	conv := FiberConverter{}

	app.Get("/", conv.Convert(Hello))

	app.Listen(":3000")
}
```

---

### 🌐 Echo

```go
package main

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/atendi9/handlerx"
	"github.com/labstack/echo/v4"
)

// ===== Context =====

type EchoContext struct {
	Ctx echo.Context
}

func (e EchoContext) Headers() map[string][]string {
	return e.Ctx.Request().Header
}

func (e EchoContext) BodyParser(v any) error {
	return e.Ctx.Bind(v)
}

func (e EchoContext) QueryParser(v any) error {
	return e.Ctx.Bind(v)
}

func (e EchoContext) ParamsParser(v any) error {
	return e.Ctx.Bind(v)
}

func (e EchoContext) ReqHeaderParser(v any) error {
	return nil
}

func (e EchoContext) Header(key string) string {
	return e.Ctx.Request().Header.Get(key)
}

func (e EchoContext) Method() string {
	return e.Ctx.Request().Method
}

func (e EchoContext) IP() string {
	return e.Ctx.RealIP()
}

func (e EchoContext) IPs() []string {
	return []string{e.Ctx.RealIP()}
}

func (e EchoContext) Body() []byte {
	return nil
}

func (e EchoContext) Query(name string, defaultValue ...string) string {
	val := e.Ctx.QueryParam(name)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func (e EchoContext) Params(name string, defaultValue ...string) string {
	val := e.Ctx.Param(name)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func (e EchoContext) FormFile(key string) (*multipart.FileHeader, error) {
	return e.Ctx.FormFile(key)
}

func (e EchoContext) SendStatus(status int) error {
	return e.Ctx.NoContent(status)
}

func (e EchoContext) Send(data []byte) error {
	return e.Ctx.Blob(http.StatusOK, "application/octet-stream", data)
}

func (e EchoContext) JSON(data any) error {
	return e.Ctx.JSON(http.StatusOK, data)
}

func (e EchoContext) Next() error {
	return nil
}

func (e EchoContext) Now() time.Time {
	return time.Now()
}

func (e EchoContext) Path(defaultValue ...string) string {
	return e.Ctx.Path()
}

// ===== Converter =====

type EchoConverter struct{}

func (e EchoConverter) Convert(h handlerx.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := EchoContext{Ctx: c}
		res := h(handlerx.Atendi9Context{Context: ctx})

		if res.GoNext() {
			return nil
		}

		if len(res.FilePath) > 0 {
			return c.File(res.FilePath)
		}

		if err := res.Err; err != nil {
			return c.JSON(res.Status(), map[string]string{
				"err": err.Error(),
			})
		}

		if v, ok := res.Data.(string); ok {
			return c.String(res.Status(), v)
		}

		return c.JSON(res.Status(), res.Data)
	}
}

// ===== Handler =====

func Hello(c handlerx.Context) handlerx.Response {
	return handlerx.Response{
		Data: map[string]string{
			"message": "Hello from Echo",
		},
	}
}

// ===== Main =====

func main() {
	e := echo.New()
	conv := EchoConverter{}

	e.GET("/", conv.Convert(Hello))

	e.Start(":3000")
}

```

---

### 🧱 net/http (Standard Library)

```go
package main

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/atendi9/handlerx"
)

// ===== Context =====

type HTTPContext struct {
	Req *http.Request
	Res http.ResponseWriter
}

func (h HTTPContext) Headers() map[string][]string {
	return h.Req.Header
}

func (h HTTPContext) BodyParser(v any) error {
	return json.NewDecoder(h.Req.Body).Decode(v)
}

func (h HTTPContext) QueryParser(v any) error {
	return nil
}

func (h HTTPContext) ParamsParser(v any) error {
	return nil
}

func (h HTTPContext) ReqHeaderParser(v any) error {
	return nil
}

func (h HTTPContext) Header(key string) string {
	return h.Req.Header.Get(key)
}

func (h HTTPContext) Method() string {
	return h.Req.Method
}

func (h HTTPContext) IP() string {
	return h.Req.RemoteAddr
}

func (h HTTPContext) IPs() []string {
	return []string{h.Req.RemoteAddr}
}

func (h HTTPContext) Body() []byte {
	return nil
}

func (h HTTPContext) Query(name string, defaultValue ...string) string {
	val := h.Req.URL.Query().Get(name)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func (h HTTPContext) Params(name string, defaultValue ...string) string {
	return ""
}

func (h HTTPContext) FormFile(key string) (*multipart.FileHeader, error) {
	return nil, nil
}

func (h HTTPContext) SendStatus(status int) error {
	h.Res.WriteHeader(status)
	return nil
}

func (h HTTPContext) Send(data []byte) error {
	h.Res.Write(data)
	return nil
}

func (h HTTPContext) JSON(data any) error {
	h.Res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(h.Res).Encode(data)
}

func (h HTTPContext) Next() error {
	return nil
}

func (h HTTPContext) Now() time.Time {
	return time.Now()
}

func (h HTTPContext) Path(defaultValue ...string) string {
	return h.Req.URL.Path
}

// ===== Converter =====

type HTTPConverter struct{}

func (h HTTPConverter) Convert(fn handlerx.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := HTTPContext{Req: r, Res: w}
		res := fn(handlerx.Atendi9Context{Context: ctx})

		if res.GoNext() {
			return
		}

		if len(res.FilePath) > 0 {
			http.ServeFile(w, r, res.FilePath)
			return
		}

		// ===== ERROR =====
		if err := res.Err; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(res.Status())
			json.NewEncoder(w).Encode(map[string]string{
				"err": err.Error(),
			})
			return
		}

		// ===== STRING =====
		if v, ok := res.Data.(string); ok {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(res.Status())
			w.Write([]byte(v))
			return
		}

		// ===== JSON =====
		w.Header().Set("Content-Type", "application/json") // FIX
		w.WriteHeader(res.Status())
		json.NewEncoder(w).Encode(res.Data)
	}
}

// ===== Handler =====

func Hello(c handlerx.Context) handlerx.Response {
	return handlerx.Response{
		Data: map[string]string{
			"message": "Hello from net/http",
		},
	}
}

// ===== Main =====

func main() {
	mux := http.NewServeMux()
	conv := HTTPConverter{}

	mux.HandleFunc("/", conv.Convert(Hello))

	http.ListenAndServe(":3000", mux)
}
```

---

## 🧠 Philosophy

HandlerX separates:

* **Transport layer** (Fiber, Echo, HTTP)
* **Business logic** (your handlers)

This makes your code:

* Easier to test 🧪
* Easier to migrate 🔄
* Easier to maintain 🧼

---

## 📄 License

MIT
