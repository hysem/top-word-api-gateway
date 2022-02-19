package topword_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hysem/top-word-api-gateway/mocks"
	"github.com/hysem/top-word-api-gateway/topword"
	st "github.com/hysem/top-word-service/topword"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("FindTopWords", func() {
	var (
		h                        *topword.Handler
		topwordServiceClientMock *mocks.TopWordServiceClient
		req                      *http.Request
		rec                      *httptest.ResponseRecorder
		err                      error
	)
	BeforeEach(func() {
		topwordServiceClientMock = &mocks.TopWordServiceClient{}
		h = topword.NewHandler(topwordServiceClientMock)
		req, err = http.NewRequest(http.MethodPost, "", nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(req).ShouldNot(BeNil())
		rec = httptest.NewRecorder()
	})
	AfterEach(func() {
		topwordServiceClientMock.AssertExpectations(GinkgoT())
	})

	It("cannot process requests other than post", func() {
		req.Method = http.MethodGet
		h.FindTopWords(rec, req)
		Expect(rec.Result().StatusCode).To(BeEquivalentTo(http.StatusMethodNotAllowed))
	})

	It("failed to get top words", func() {
		form := url.Values{}
		form.Add("text", "test")
		req.PostForm = form
		topwordServiceClientMock.On("FindTopWords", mock.Anything, &st.FindTopWordsRequest{
			Text: "test",
		}).Return(nil, errors.New("failed to get top words"))
		h.FindTopWords(rec, req)
		Expect(rec.Result().StatusCode).To(BeEquivalentTo(http.StatusInternalServerError))
	})
	It("retrieved top words", func() {
		form := url.Values{}
		form.Add("text", "test")
		req.PostForm = form

		topwordServiceClientMock.On("FindTopWords", mock.Anything, &st.FindTopWordsRequest{
			Text: "test",
		}).Return([]*st.WordInfo{
			{Word: "test", Count: 1},
		}, nil)

		h.FindTopWords(rec, req)
		Expect(rec.Result().StatusCode).To(BeEquivalentTo(http.StatusOK))

		var topWords []*st.WordInfo
		err := json.NewDecoder(rec.Body).Decode(&topWords)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(topWords).Should(BeEquivalentTo([]*st.WordInfo{
			{Word: "test", Count: 1},
		}))
	})
})

func TestTopword(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Topword Suite")
}
