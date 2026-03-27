package errors

import (
	stderrors "errors"
	"net/http"
	"testing"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	if Message(CodeSuccess) != "success" {
		t.Fatal("unexpected success message")
	}
	if Message(999999) != "未知错误" {
		t.Fatal("unexpected unknown message")
	}
}

func TestWrapAndExtract(t *testing.T) {
	t.Parallel()

	origin := stderrors.New("origin")
	e := Wrap(CodeInternalError, origin)
	if e == nil || e.Code != CodeInternalError {
		t.Fatal("wrap failed")
	}
	if !IsAppError(e) {
		t.Fatal("should be app error")
	}
	ae := GetAppError(e)
	if ae == nil || ae.Err == nil {
		t.Fatal("extract failed")
	}
	if !stderrors.Is(e, origin) {
		t.Fatal("unwrap chain not working")
	}
}

func TestWrapKeepExistingAppError(t *testing.T) {
	t.Parallel()

	base := ErrUnauthorized()
	got := Wrap(CodeInternalError, base)
	if got != base {
		t.Fatal("expected existing app error returned directly")
	}
}

func TestHTTPStatusMapping(t *testing.T) {
	t.Parallel()

	cases := []struct {
		code int
		want int
	}{
		{CodeSuccess, http.StatusOK},
		{CodeInvalidParams, http.StatusBadRequest},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeTokenInvalid, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeUserDisabled, http.StatusForbidden},
		{CodeUserAlreadyExists, http.StatusConflict},
		{CodePasswordIncorrect, http.StatusUnauthorized},
		{CodeCategoryNotFound, http.StatusNotFound},
		{CodeTooManyRequests, http.StatusTooManyRequests},
		{123456, http.StatusInternalServerError},
	}
	for _, c := range cases {
		if got := HTTPStatus(c.code); got != c.want {
			t.Fatalf("code=%d got=%d want=%d", c.code, got, c.want)
		}
	}
}

func TestConstructors(t *testing.T) {
	t.Parallel()

	e := New(CodeInvalidParams, "bad")
	if e.Code != CodeInvalidParams || e.Message != "bad" {
		t.Fatal("new failed")
	}
	if ErrPasswordIncorrect().Code != CodePasswordIncorrect {
		t.Fatal("shortcut constructor failed")
	}
}
