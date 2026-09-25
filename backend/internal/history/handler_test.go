package history

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeStore is an in-memory Store used to test Handler without a live
// Postgres connection.
type fakeStore struct {
	createFn func(ctx context.Context, expression, result string) (History, error)
	listFn   func(ctx context.Context) ([]History, error)

	createCalled bool
	createExpr   string
	createResult string
}

func (f *fakeStore) Create(ctx context.Context, expression, result string) (History, error) {
	f.createCalled = true
	f.createExpr = expression
	f.createResult = result
	return f.createFn(ctx, expression, result)
}

func (f *fakeStore) List(ctx context.Context) ([]History, error) {
	return f.listFn(ctx)
}

func TestHandlerCalculate(t *testing.T) {
	fixedTime := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name           string
		body           string
		store          *fakeStore
		wantStatus     int
		wantBodyPart   string
		wantStoreCalls bool
	}{
		{
			name: "valid expression is evaluated and persisted",
			body: `{"expression":"3 + 4 * 2"}`,
			store: &fakeStore{
				createFn: func(ctx context.Context, expression, result string) (History, error) {
					return History{ID: 1, Expression: expression, Result: result, CreatedAt: fixedTime}, nil
				},
			},
			wantStatus:     http.StatusCreated,
			wantBodyPart:   `"result":"11"`,
			wantStoreCalls: true,
		},
		{
			name:           "invalid expression is rejected without persisting",
			body:           `{"expression":"2 +"}`,
			store:          &fakeStore{createFn: failIfCalled(t)},
			wantStatus:     http.StatusBadRequest,
			wantStoreCalls: false,
		},
		{
			name:           "division by zero is rejected without persisting",
			body:           `{"expression":"5 / 0"}`,
			store:          &fakeStore{createFn: failIfCalled(t)},
			wantStatus:     http.StatusBadRequest,
			wantBodyPart:   `"error":"division by zero"`,
			wantStoreCalls: false,
		},
		{
			name:           "malformed JSON body",
			body:           `not json`,
			store:          &fakeStore{createFn: failIfCalled(t)},
			wantStatus:     http.StatusBadRequest,
			wantStoreCalls: false,
		},
		{
			name: "store error surfaces as 500",
			body: `{"expression":"1 + 1"}`,
			store: &fakeStore{
				createFn: func(ctx context.Context, expression, result string) (History, error) {
					return History{}, errors.New("db unavailable")
				},
			},
			wantStatus:     http.StatusInternalServerError,
			wantStoreCalls: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler(tc.store)
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			h.Calculate(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantBodyPart != "" && !strings.Contains(rec.Body.String(), tc.wantBodyPart) {
				t.Fatalf("body = %s, want it to contain %s", rec.Body.String(), tc.wantBodyPart)
			}
			if tc.store.createCalled != tc.wantStoreCalls {
				t.Fatalf("store.Create called = %v, want %v", tc.store.createCalled, tc.wantStoreCalls)
			}
		})
	}
}

func TestHandlerList(t *testing.T) {
	fixedTime := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

	t.Run("returns items from the store", func(t *testing.T) {
		store := &fakeStore{
			listFn: func(ctx context.Context) ([]History, error) {
				return []History{
					{ID: 2, Expression: "2 + 2", Result: "4", CreatedAt: fixedTime},
					{ID: 1, Expression: "1 + 1", Result: "2", CreatedAt: fixedTime},
				}, nil
			},
		}
		h := NewHandler(store)
		req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var got []History
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if len(got) != 2 || got[0].ID != 2 || got[1].ID != 1 {
			t.Fatalf("unexpected history items: %+v", got)
		}
	})

	t.Run("store error surfaces as 500", func(t *testing.T) {
		store := &fakeStore{
			listFn: func(ctx context.Context) ([]History, error) {
				return nil, errors.New("db unavailable")
			},
		}
		h := NewHandler(store)
		req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}

func failIfCalled(t *testing.T) func(ctx context.Context, expression, result string) (History, error) {
	t.Helper()
	return func(ctx context.Context, expression, result string) (History, error) {
		t.Fatalf("store.Create should not be called for an invalid expression (got expression=%q, result=%q)", expression, result)
		return History{}, nil
	}
}
