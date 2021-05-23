package webapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage_ServeHTTP(t *testing.T) {
	expected := "hello there"
	ip := indexPage(expected)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()
	ip.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
