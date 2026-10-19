package httpapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
)

func TestWriteUseCaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		err        error
		statusCode int
		message    string
	}{
		{name: "invalid input", err: domain.ErrInvalidDocument, statusCode: http.StatusBadRequest, message: domain.ErrInvalidDocument.Error()},
		{name: "account not found", err: domain.ErrAccountNotFound, statusCode: http.StatusNotFound, message: domain.ErrAccountNotFound.Error()},
		{name: "internal error", err: errors.New("database credentials leaked"), statusCode: http.StatusInternalServerError, message: "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			writeUseCaseError(ctx, tt.err)

			if recorder.Code != tt.statusCode {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.statusCode)
			}
			if got := recorder.Body.String(); got != `{"error":"`+tt.message+`"}` {
				t.Errorf("body = %s, want error %q", got, tt.message)
			}
		})
	}
}
