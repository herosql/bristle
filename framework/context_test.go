package framework

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func setup(method, url string, body []byte) (*Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	w := httptest.NewRecorder()
	return NewContext(w, req), w
}

func TestQueryString(t *testing.T) {
	t.Run("ValueExists", func(t *testing.T) {
		ctx, _ := setup("GET", "/?user=gemini", nil)
		if got := ctx.QueryString("user", "def"); got != "gemini" {
			t.Errorf("want gemini, got %s", got)
		}
	})

	t.Run("ReturnsDefault", func(t *testing.T) {
		ctx, _ := setup("GET", "/", nil)
		if got := ctx.QueryString("user", "admin"); got != "admin" {
			t.Errorf("want default admin, got %s", got)
		}
	})
}

func TestQueryInt(t *testing.T) {
	t.Run("ValidInt", func(t *testing.T) {
		ctx, _ := setup("GET", "/?page=2", nil)
		if got := ctx.QueryInt("page", 1); got != 2 {
			t.Errorf("want 2, got %d", got)
		}
	})

	t.Run("InvalidFormatReturnsDefault", func(t *testing.T) {
		ctx, _ := setup("GET", "/?page=abc", nil)
		if got := ctx.QueryInt("page", 1); got != 1 {
			t.Errorf("want default 1, got %d", got)
		}
	})
}

func TestQueryArray(t *testing.T) {
	t.Run("MultipleValues", func(t *testing.T) {
		ctx, _ := setup("GET", "/?id=1&id=2", nil)
		want := []string{"1", "2"}
		if got := ctx.QueryArray("id", nil); !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func TestBindJson(t *testing.T) {
	t.Run("ValidPayload", func(t *testing.T) {
		ctx, _ := setup("POST", "/", []byte(`{"id":1}`))
		var data struct{ ID int }
		if err := ctx.BindJson(&data); err != nil {
			t.Fatalf("bind failed: %v", err)
		}
	})

	t.Run("MalformedInput", func(t *testing.T) {
		ctx, _ := setup("POST", "/", []byte(`{invalid}`))
		var data struct{ ID int }
		if err := ctx.BindJson(&data); err == nil {
			t.Error("expected error for malformed json, got nil")
		}
	})
}

func TestJson(t *testing.T) {
	t.Run("RenderSuccess", func(t *testing.T) {
		ctx, w := setup("GET", "/", nil)
		_ = ctx.Json(http.StatusOK, map[string]string{"result": "ok"})

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Error("missing content-type header")
		}
	})
}

func TestTimeout(t *testing.T) {
	t.Run("DefaultFalse", func(t *testing.T) {
		ctx, _ := setup("GET", "/", nil)
		if ctx.HasTimeout() {
			t.Error("initial timeout should be false")
		}
	})

	t.Run("SetAndGet", func(t *testing.T) {
		ctx, _ := setup("GET", "/", nil)
		ctx.SetHasTimeout(true)
		if !ctx.HasTimeout() {
			t.Error("timeout should be true after setting")
		}
	})
}
