package topword_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hysem/top-word-api-gateway/mocks"
	"github.com/hysem/top-word-api-gateway/topword"
	st "github.com/hysem/top-word-service/topword"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type handlerMocks struct {
	topWordServiceClient mocks.TopWordServiceClientMock
}

func (m *handlerMocks) assertExpectations(t *testing.T) {
	m.topWordServiceClient.AssertExpectations(t)
}
func newHandler(t *testing.T) (*topword.Handler, *handlerMocks) {
	m := handlerMocks{}
	h := topword.NewHandler(&m.topWordServiceClient)
	return h, &m
}

func TestFindTopWords(t *testing.T) {
	const testParagraph = "paragraph to test"
	var testResponse = []*st.WordInfo{{
		Word:  "paragraph",
		Count: 1,
	}, {
		Word:  "to",
		Count: 1,
	}, {
		Word:  "test",
		Count: 1,
	}}
	testCases := map[string]struct {
		requestMethod  string
		requestText    string
		expectedStatus int
		expectedBody   interface{}
		setMocks       func(m *handlerMocks)
	}{
		`error case: invalid request method`: {
			requestMethod:  http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		`error case: failed to find top words`: {
			requestMethod: http.MethodPost,
			requestText:   testParagraph,
			setMocks: func(m *handlerMocks) {
				m.topWordServiceClient.On("FindTopWords", mock.Anything, &st.FindTopWordsRequest{
					Text: testParagraph,
				}).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		`success case: found top words`: {
			requestMethod: http.MethodPost,
			requestText:   testParagraph,
			setMocks: func(m *handlerMocks) {
				m.topWordServiceClient.On("FindTopWords", mock.Anything, &st.FindTopWordsRequest{
					Text: testParagraph,
				}).Return(testResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testResponse,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			h, m := newHandler(t)
			defer m.assertExpectations(t)
			if tc.setMocks != nil {
				tc.setMocks(m)
			}

			form := url.Values{}
			form.Add("text", tc.requestText)

			req, err := http.NewRequest(tc.requestMethod, "/endpoint", nil)
			require.NoError(t, err)
			req.PostForm = form
			rec := httptest.NewRecorder()
			h.FindTopWords(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Result().StatusCode)
			if tc.expectedBody != nil {
				var actualResponse []*st.WordInfo
				err := json.NewDecoder(rec.Body).Decode(&actualResponse)
				require.NoError(t, err)
				require.Equal(t, actualResponse, testResponse)
			} else {
				require.Empty(t, rec.Body.String())
			}
		})
	}
}
