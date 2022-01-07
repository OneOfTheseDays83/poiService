package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEchoHandlerWith(data string, status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if data != "" {
			w.Write([]byte(data))
		}
	})
}

func TestRequester_OkCaseForGet(t *testing.T) {
	// given
	testData := "abc"
	responseStatus := http.StatusOK
	testKeyServer := httptest.NewServer(getEchoHandlerWith(testData, responseStatus))
	requester := NewHttpRequester(http.DefaultClient)

	// when
	receivedData, err := requester.GetContent(testKeyServer.URL)

	// then
	require.NoError(t, err)
	assert.Equal(t, receivedData, testData)
}

func TestRequester_GetRequestWithUnsupportedProtocolScheme(t *testing.T) {
	// when
	requester := NewHttpRequester(http.DefaultClient)
	receivedData, err := requester.GetContent("")

	// then
	assert.Equal(t, receivedData, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestRequester_GetWithErrorStatusCode(t *testing.T) {
	// given
	testData := "abc"
	responseStatus := http.StatusInternalServerError
	testKeyServer := httptest.NewServer(getEchoHandlerWith(testData, responseStatus))
	requester := NewHttpRequester(http.DefaultClient)

	// when
	receivedData, err := requester.GetContent(testKeyServer.URL)

	// then
	assert.Equal(t, receivedData, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("status: %d", responseStatus))
}
