package handlerx

import "testing"

func TestHandlerExecution(t *testing.T) {
	h := func(_ Context) Response {
		return Response{
			Data: "ok",
		}
	}
	res := h(testHTTPContext{})

	if res.Data != "ok" {
		t.Fatalf("expected 'ok', got %v", res.Data)
	}
}
