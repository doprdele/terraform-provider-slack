package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCanvas(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/conversations.canvases.create":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok": true, "id": "C061EG9SL3D"}`))
		default:
			t.Errorf("Unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.httpClient = server.Client()

	id, err := client.CreateCanvas("Hello, world!", "C12345")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedID := "C061EG9SL3D"
	if id != expectedID {
		t.Errorf("expected canvas ID '%s', got '%s'", expectedID, id)
	}
}

func TestReadCanvas(t *testing.T) {
	// No-op test
}
