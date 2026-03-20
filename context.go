package handlerx

import (
	"mime/multipart"
	"time"
)

// Context defines an abstraction over an HTTP request/response lifecycle.
//
// It provides a unified API to access request data (headers, body, params),
// parse input into structs, and build responses.
//
// This interface is designed to be framework-agnostic, allowing different
// HTTP engines (Fiber, Echo, net/http, etc.) to implement it.
type Context interface {
	// Headers returns all request headers.
	//
	// The returned map follows the standard Go format:
	// map[string][]string.
	Headers() map[string][]string

	// BodyParser parses the request body into the given struct.
	//
	// It supports formats such as JSON, XML, or form data,
	// depending on the implementation.
	//
	// Example:
	//	var req MyStruct
	//	if err := c.BodyParser(&req); err != nil {
	//	    return err
	//	}
	BodyParser(v any) error

	// QueryParser parses query string parameters into the given struct.
	//
	// Example:
	//	var query QueryParams
	//	_ = c.QueryParser(&query)
	QueryParser(v any) error

	// ParamsParser parses route (path) parameters into the given struct.
	//
	// Example:
	//	// Route: /users/:id
	//	var params struct { ID string `param:"id"` }
	//	_ = c.ParamsParser(&params)
	ParamsParser(v any) error

	// ReqHeaderParser parses request headers into the given struct.
	//
	// This is useful for binding headers to typed structures.
	ReqHeaderParser(v any) error

	// Header returns the value of a specific request header.
	//
	// If the header does not exist, an empty string is returned.
	Header(key string) string

	// Method returns the HTTP method used in the request
	// (e.g., GET, POST, PUT, DELETE).
	Method() string

	// IP returns the client IP address.
	//
	// Depending on the implementation, this may consider proxy headers
	// such as X-Forwarded-For.
	IP() string

	// IPs returns all IP addresses associated with the request,
	// including proxy chain addresses.
	IPs() []string

	// Body returns the raw request body.
	//
	// Useful when manual parsing is needed.
	Body() []byte

	// Query returns the value of a query parameter.
	//
	// If the parameter is not present, it returns the optional default value.
	//
	// Example:
	//	name := c.Query("name", "guest")
	Query(name string, defaultValue ...string) string

	// Params returns the value of a route (path) parameter.
	//
	// If the parameter is not present, it returns the optional default value.
	//
	// Example:
	//	id := c.Params("id", "0")
	Params(name string, defaultValue ...string) string

	// FormFile retrieves a file uploaded via multipart form.
	//
	// Returns a pointer to multipart.FileHeader, which can be used
	// to open and read the file.
	//
	// Example:
	//	file, err := c.FormFile("avatar")
	FormFile(key string) (*multipart.FileHeader, error)

	// SendStatus sets the HTTP status code and sends the response.
	//
	// Example:
	//	return c.SendStatus(404)
	SendStatus(status int) error

	// Send writes raw bytes as the response body.
	//
	// Example:
	//	return c.Send([]byte("ok"))
	Send(data []byte) error

	// JSON serializes the given data as JSON and writes it to the response.
	//
	// Example:
	//	return c.JSON(map[string]string{"status": "ok"})
	JSON(data any) error

	// Next passes control to the next middleware/handler in the chain.
	//
	// Commonly used in middleware pipelines.
	Next() error

	// Now returns the current time.
	//
	// This abstraction allows easier testing by mocking time.
	Now() time.Time

	// Path returns the request path.
	//
	// If no path is available, it may return the optional default value.
	Path(defaultValue ...string) string
}

// Atendi9Context is a wrapper around Context that allows
// extending or customizing behavior without modifying the original implementation.
//
// It can be used to add helper methods specific to your application.
type Atendi9Context struct {
	Context
}

// NewContext creates a new Atendi9Context wrapping the given Context.
//
// Example:
//	ctx := NewContext(originalCtx)
func NewContext(ctx Context) *Atendi9Context {
	return &Atendi9Context{
		Context: ctx,
	}
}

// Test replaces the underlying Context instance.
//
// This method is primarily useful for testing or dynamically swapping
// context implementations.
//
// Example:
//	ctx.Test(mockCtx)
func (c *Atendi9Context) Test(ctx Context) *Atendi9Context {
	c.Context = ctx
	return c
}