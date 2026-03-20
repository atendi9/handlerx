package handlerx

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAtendi9Context(t *testing.T) {
	t.Run("Context", func(t *testing.T) {
		_, ok := any(Atendi9Context{}).(Context)
		if !ok {
			t.Fatal("Atendi9Context does not implement Context (API changed)")
		}
	})
	t.Run("Test", func(t *testing.T) {
		httpContext := testHTTPContext{Req: httptest.NewRequest(http.MethodGet, "https://google.com", nil), Res: httptest.NewRecorder()}
		ctx := NewContext(httpContext)
		testContext := testHTTPContext{Req: httptest.NewRequest(http.MethodGet, "https://www.atendi9.com.br", nil), Res: httptest.NewRecorder()}
		ctx = ctx.Test(testContext)
		if ctx.Context != testContext || ctx.Context == httpContext {
			t.Fail()
		}
	})
}

type testHTTPContext struct {
	Req *http.Request
	Res http.ResponseWriter
}

func (h testHTTPContext) Headers() map[string][]string {
	return h.Req.Header
}

func (h testHTTPContext) BodyParser(v any) error {
	return json.NewDecoder(h.Req.Body).Decode(v)
}

func (h testHTTPContext) QueryParser(v any) error {
	return nil
}

func (h testHTTPContext) ParamsParser(v any) error {
	return nil
}

func (h testHTTPContext) ReqHeaderParser(v any) error {
	return nil
}

func (h testHTTPContext) Header(key string) string {
	return h.Req.Header.Get(key)
}

func (h testHTTPContext) Method() string {
	return h.Req.Method
}

func (h testHTTPContext) IP() string {
	return h.Req.RemoteAddr
}

func (h testHTTPContext) IPs() []string {
	return []string{h.Req.RemoteAddr}
}

func (h testHTTPContext) Body() []byte {
	return nil
}

func (h testHTTPContext) Query(name string, defaultValue ...string) string {
	val := h.Req.URL.Query().Get(name)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func (h testHTTPContext) Params(name string, defaultValue ...string) string {
	return ""
}

func (h testHTTPContext) FormFile(key string) (*multipart.FileHeader, error) {
	return nil, nil
}

func (h testHTTPContext) SendStatus(status int) error {
	h.Res.WriteHeader(status)
	return nil
}

func (h testHTTPContext) Send(data []byte) error {
	h.Res.Write(data)
	return nil
}

func (h testHTTPContext) JSON(data any) error {
	h.Res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(h.Res).Encode(data)
}

func (h testHTTPContext) Next() error {
	return nil
}

func (h testHTTPContext) Now() time.Time {
	return time.Now()
}

func (h testHTTPContext) Path(defaultValue ...string) string {
	return h.Req.URL.Path
}
