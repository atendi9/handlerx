package handlerx

import (
	"errors"
	"testing"
)

func TestResponseStatus_Default(t *testing.T) {
	r := Response{}

	if r.Status() != 200 {
		t.Fatalf("expected default status 200, got %d", r.Status())
	}
}

func TestResponseStatus_Custom(t *testing.T) {
	r := Response{StatusCode: 404}

	if r.Status() != 404 {
		t.Fatalf("expected status 404, got %d", r.Status())
	}
}

func TestResponseStatus_Informational(t *testing.T) {
	// 1xx codes are valid and must be preserved, not clamped to 200.
	for _, code := range []int{100, 101, 103} {
		r := Response{StatusCode: code}

		if r.Status() != code {
			t.Fatalf("expected informational status %d to be preserved, got %d", code, r.Status())
		}
	}
}

func TestResponseStatus_ExplicitOK(t *testing.T) {
	// An explicitly set 200 must be returned as-is.
	r := Response{StatusCode: 200}

	if r.Status() != 200 {
		t.Fatalf("expected status 200, got %d", r.Status())
	}
}

func TestResponseStatus_OutOfRange(t *testing.T) {
	// Codes outside [100, 599] are invalid and fall back to 200.
	for _, code := range []int{-1, 99, 600, 1000} {
		r := Response{StatusCode: code}

		if r.Status() != 200 {
			t.Fatalf("expected fallback status 200 for invalid code %d, got %d", code, r.Status())
		}
	}
}

func TestResponseNext(t *testing.T) {
	r := Response{}

	r = r.Next()

	if !r.GoNext() {
		t.Fatal("expected GoNext() to be true after calling Next()")
	}
}

func TestResponseGoNext_Default(t *testing.T) {
	r := Response{}

	if r.GoNext() {
		t.Fatal("expected GoNext() to be false by default")
	}
}

func TestSendStatus(t *testing.T) {
	r := SendStatus(404)

	if r.StatusCode != 404 {
		t.Fatalf("expected status code 404, got %d", r.StatusCode)
	}
}

func TestResponseJSON(t *testing.T) {
	r := Response{StatusCode: 201}

	createdMsg := "created"
	payload := map[string]string{
		"message": createdMsg,
	}

	res := r.JSON(payload)

	if res.StatusCode != 201 {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}

	data, ok := res.Data.(map[string]string)
	if !ok {
		t.Fatalf("expected Data to be map[string]string, got %T", res.Data)
	}

	if data["message"] != createdMsg {
		t.Fatalf("expected message '%s', got '%s'", createdMsg, data["message"])
	}
}

func TestResponseJSON_DefaultStatusFallback(t *testing.T) {
	r := Response{}

	res := r.JSON(nil)

	if res.StatusCode != 200 {
		t.Fatalf("expected default status 200, got %d", res.StatusCode)
	}
}

func TestResponseJSON_PreservesOtherFields(t *testing.T) {
	// JSON must not silently drop next, Err or FilePath already set
	// on the receiver.
	wantErr := errors.New("boom")
	r := Response{
		StatusCode: 201,
		Err:        wantErr,
		FilePath:   "/tmp/file.txt",
	}.Next()

	res := r.JSON(map[string]string{"message": "ok"})

	if !res.GoNext() {
		t.Fatal("expected GoNext() to be preserved through JSON()")
	}
	if res.Err != wantErr {
		t.Fatalf("expected Err to be preserved, got %v", res.Err)
	}
	if res.FilePath != "/tmp/file.txt" {
		t.Fatalf("expected FilePath to be preserved, got %q", res.FilePath)
	}
	if res.StatusCode != 201 {
		t.Fatalf("expected status 201 to be preserved, got %d", res.StatusCode)
	}
	if _, ok := res.Data.(map[string]string); !ok {
		t.Fatalf("expected Data to be set, got %T", res.Data)
	}
}
