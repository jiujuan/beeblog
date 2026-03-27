package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "antblog/pkg/errors"

	"github.com/gin-gonic/gin"
)

func TestNewPageDataNilList(t *testing.T) {
	t.Parallel()

	data := NewPageData[int](nil, 10, 1, 10)
	if data.List == nil || len(data.List) != 0 {
		t.Fatal("expected non-nil empty list")
	}
}

func TestOKResponse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OK(c, gin.H{"a": 1})

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if resp.Code != apperrors.CodeSuccess || resp.Msg != "success" {
		t.Fatal("unexpected response body")
	}
}

func TestFailWithAppError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	FailWithError(c, apperrors.ErrInvalidParams("bad req"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if resp.Code != apperrors.CodeInvalidParams || resp.Msg != "bad req" {
		t.Fatal("unexpected response body")
	}
}

func TestFailWithNonAppError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	FailWithError(c, errors.New("boom"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if resp.Code != apperrors.CodeInternalError {
		t.Fatal("unexpected error code")
	}
}

func TestNoContent(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	NoContent(c)
	if c.Writer.Status() != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", c.Writer.Status())
	}
}
