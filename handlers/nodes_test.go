package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/4serviceSoftware/tech-task/nodes"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestPostPositive(t *testing.T) {
	logger := log.New(os.Stdout, "tech-task-one-test", log.LstdFlags)

	dbUrl := "postgres://kbnq:root@localhost:5432/techtaskone_test"
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		t.Fatal("DB conn: " + err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	// creating nodes repository
	nodesRepo := nodes.NewRepositoryPostgres(db, ctx)

	nh := NewNodes(nodesRepo, logger)

	// Create a request to pass to our handler.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for i, fileName := range []string{"nodes-1.xlsx", "nodes-2.xlsx", "nodes-3.xlsx"} {
		fieldName := fmt.Sprintf("file%d", i)
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			t.Fatal(err)
		}
		file1, err := os.Open("./../" + fileName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, file1)
		if err != nil {
			t.Fatal(err)
		}
	}
	boundary := writer.Boundary()
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/nodes", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nh.Post)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",
			status, http.StatusOK, rr.Body)
	}
}

func TestPostNegative(t *testing.T) {
	logger := log.New(os.Stdout, "tech-task-one-test", log.LstdFlags)

	dbUrl := "postgres://kbnq:root@localhost:5432/techtaskone_test"
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		t.Fatal("DB conn: " + err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	// creating nodes repository
	nodesRepo := nodes.NewRepositoryPostgres(db, ctx)

	nh := NewNodes(nodesRepo, logger)

	// Create a request to pass to our handler.
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for i, fileName := range []string{"nodes-circular.xlsx", "nodes-2.xlsx", "nodes-3.xlsx"} {
		fieldName := fmt.Sprintf("file%d", i)
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			t.Fatal(err)
		}
		file1, err := os.Open("./../" + fileName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, file1)
		if err != nil {
			t.Fatal(err)
		}
	}
	boundary := writer.Boundary()
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/nodes", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nh.Post)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want some error. Body: %s",
			status, rr.Body)
	}
}

func TestGet(t *testing.T) {
	logger := log.New(os.Stdout, "tech-task-one-test", log.LstdFlags)

	dbUrl := "postgres://kbnq:root@localhost:5432/techtaskone_test"
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		t.Fatal("DB conn: " + err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	// creating nodes repository
	nodesRepo := nodes.NewRepositoryPostgres(db, ctx)

	nh := NewNodes(nodesRepo, logger)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nh.Get)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is valid json
	if !json.Valid([]byte(rr.Body.String())) {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}
