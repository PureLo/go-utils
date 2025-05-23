package wsserver

import (
	"net/http"
	"testing"
)

func TestHub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		t.Fatalf("ListenAndServe: %v", err)
	}
}
