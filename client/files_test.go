package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestFiles_ListFiles(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/files/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dir := r.URL.Query().Get("directory")
		if dir != "/" {
			t.Errorf("directory = %s, want /", dir)
		}

		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "file_object",
				"attributes": {
					"name": "server.jar",
					"mode": "-rw-r--r--",
					"size": 1024,
					"is_file": true,
					"is_symlink": false,
					"mimetype": "application/java-archive"
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	files, err := clientAPI.ListFiles(context.Background(), "d3aac109", "/")
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
	if files[0].Name != "server.jar" {
		t.Errorf("Name = %s, want server.jar", files[0].Name)
	}
}

func TestFiles_CreateDirectory(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/files/create-folder", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.CreateDirectory(context.Background(), "d3aac109", "/", "plugins")
	if err != nil {
		t.Fatalf("CreateDirectory() error = %v", err)
	}
}

func TestFiles_DeleteFiles(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/files/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.DeleteFiles(context.Background(), "d3aac109", "/", []string{"old-file.txt"})
	if err != nil {
		t.Fatalf("DeleteFiles() error = %v", err)
	}
}

func TestFiles_CompressFiles(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/files/compress", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{
			"object": "file_object",
			"attributes": {
				"name": "archive.tar.gz",
				"size": 2048,
				"is_file": true
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	archive, err := clientAPI.CompressFiles(context.Background(), "d3aac109", "/", []string{"world", "plugins"})
	if err != nil {
		t.Fatalf("CompressFiles() error = %v", err)
	}

	if archive.Name != "archive.tar.gz" {
		t.Errorf("Name = %s, want archive.tar.gz", archive.Name)
	}
}
