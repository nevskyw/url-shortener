package delete_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "NorFound",
			alias:     "test_alias",
			respError: "not found",
			mockError: storage.ErrURLNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleteMock := mocks.NewURLDeleter(t)

			// Expect DeleteURL to be called with the alias from the test case
			urlDeleteMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()

			r := chi.NewRouter()
			r.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleteMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := api.DeleteURL(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Check thr response status
			if tc.respError == "" {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				errorResponse := api.ParseErrorResponse(resp)
				assert.Equal(t, tc.respError, errorResponse.Error)
			}
		})
	}
}
