package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCanvas(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/conversations.canvases.create":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok": true, "canvas_id": "C061EG9SL3D"}`))
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

func TestCreateUserCanvas_withChannelID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/canvases.create":
			body, _ := ioutil.ReadAll(r.Body)
			var data map[string]interface{}
			json.Unmarshal(body, &data)

			if data["channel_id"] != "C12345" {
				t.Errorf("expected channel_id to be C12345, got %s", data["channel_id"])
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok": true, "canvas_id": "F061EG9SL3D"}`))
		default:
			t.Errorf("Unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.httpClient = server.Client()

	_, err := client.CreateUserCanvas("Hello, world!", "C12345", false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestReadCanvas(t *testing.T) {
	// No-op test
}