package handlerx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConverter(t *testing.T) {
	converter := testHTTPConverter{}
	t.Run("SendString", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{Data: "converted"}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api", nil))

		if w.Body.String() != "converted" {
			t.Fatalf("expected 'converted', got %v", w.Body.String())
		}
	})

	t.Run("SendJSON", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{Data: map[string]string{"msg": "ok"}}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("expected application/json, got %s", ct)
		}

		var body map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}

		if body["msg"] != "ok" {
			t.Fatalf("expected ok, got %v", body["msg"])
		}
	})

	t.Run("SendError", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{
				Err:        errors.New("fail"),
				StatusCode: 400,
			}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if w.Code != 400 {
			t.Fatalf("expected 400, got %d", w.Code)
		}

		var body map[string]string
		json.Unmarshal(w.Body.Bytes(), &body)

		if body["err"] != "fail" {
			t.Fatalf("expected fail, got %v", body["err"])
		}
	})

	t.Run("CustomStatus", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{
				Data:       "created",
				StatusCode: 201,
			}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", nil))

		if w.Code != 201 {
			t.Fatalf("expected 201, got %d", w.Code)
		}
	})

	t.Run("GoNext", func(t *testing.T) {
		called := false

		handler := func(c Context) Response {
			called = true
			return Response{}.Next()
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if !called {
			t.Fatal("handler was not called")
		}

		if w.Body.Len() != 0 {
			t.Fatal("expected empty response when GoNext is true")
		}
	})

	t.Run("FilePath", func(t *testing.T) {
		tmp := t.TempDir() + "/file.txt"
		content := []byte("file-content")

		if err := os.WriteFile(tmp, content, 0644); err != nil {
			t.Fatal(err)
		}

		handler := func(c Context) Response {
			return Response{FilePath: tmp}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if w.Body.String() != "file-content" {
			t.Fatalf("expected file-content, got %s", w.Body.String())
		}
	})

	t.Run("DefaultStatusIs200", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{Data: "ok"}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if w.Code != 200 {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("ContentTypeString", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{Data: "plain"}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if ct := w.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
			t.Fatalf("unexpected content-type: %s", ct)
		}
	})

	t.Run("ContentTypeJSONFallback", func(t *testing.T) {
		handler := func(c Context) Response {
			return Response{Data: 123}
		}

		fn := converter.Convert(handler)
		w := httptest.NewRecorder()
		fn.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("unexpected content-type: %s", ct)
		}
	})

}

type testHTTPConverter struct{}

func (h testHTTPConverter) Convert(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := testHTTPContext{Req: r, Res: w}
		res := fn(Atendi9Context{Context: ctx})

		if res.GoNext() {
			return
		}

		if len(res.FilePath) > 0 {
			http.ServeFile(w, r, res.FilePath)
			return
		}

		if err := res.Err; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(res.Status())
			json.NewEncoder(w).Encode(map[string]string{
				"err": err.Error(),
			})
			return
		}

		if v, ok := res.Data.(string); ok {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(res.Status())
			w.Write([]byte(v))
			return
		}

		w.Header().Set("Content-Type", "application/json") // FIX
		w.WriteHeader(res.Status())
		json.NewEncoder(w).Encode(res.Data)
	}
}
