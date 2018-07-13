package onetimesecret_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jasosa/onetimesecret"
)

var serverAnswer = `{"custid":"user","metadata_key":"934523j02349rjkf","secret_key":"k5o8n8cyvbt2xwasgbx6r3b5whs35s",
					"ttl":7200,"metadata_ttl":7200,"secret_ttl":3600,"state":"new","updated":1531399574,"created":1531399574,"recipient":[],
					"value":"zpuOBYWPa6*3","passphrase_required":false}`
var unmarshalingResponseErrorAnswer = `{error}`

type testHandler struct {
	calls []string
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.calls = append(h.calls, req.URL.String())
	if req.URL.String() == "/generate?ttl=0" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(unmarshalingResponseErrorAnswer))
	} else if req.URL.String() == "/generate?ttl=1" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{0})
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverAnswer))
	}
}

func TestGenerateSuccesfully(t *testing.T) {

	th := &testHandler{}
	ts := httptest.NewServer(th)

	expectedSecretKey := "k5o8n8cyvbt2xwasgbx6r3b5whs35s"
	expectedSecretValue := "zpuOBYWPa6*3"
	expectedServerRequest := "/generate?ttl=3600"
	client := onetimesecret.NewClient("user", "token", ts.URL)
	secretKey, secretvalue, err := client.Generate(3600)

	if th.calls[0] != expectedServerRequest {
		t.Fatalf("Expected server request key was %q but got %q", expectedServerRequest, th.calls[0])
	}

	if err != nil {
		t.Fatalf("Error was not expected but got %v", err)
	}

	if secretKey != expectedSecretKey {
		t.Fatalf("Expected secret key was %q but got %q", expectedSecretKey, secretKey)
	}
	if secretvalue != expectedSecretValue {
		t.Fatalf("Expected secret value was %q but got %q", expectedSecretValue, secretvalue)
	}
}

func TestGenerateErrorExecutingRequest(t *testing.T) {
	expectedErrorMessage := "error executing request"

	client := onetimesecret.NewClient("user", "token", "")
	_, _, err := client.Generate(3600)

	if err != nil && !strings.HasPrefix(err.Error(), expectedErrorMessage) {
		t.Fatalf("Error %q was expected but got %q", expectedErrorMessage, err.Error())
	}
}
func TestGenerateErrorCreatingRequest(t *testing.T) {
	expectedErrorMessage := "error creating request"

	client := onetimesecret.NewClient("user", "token", "://127.0.0.1:3245")
	_, _, err := client.Generate(3600)

	if err != nil && !strings.HasPrefix(err.Error(), expectedErrorMessage) {
		t.Fatalf("Error %q was expected but got %q", expectedErrorMessage, err.Error())
	}
}

func TestGenerateErrorUnmarshalingResponse(t *testing.T) {
	th := &testHandler{}
	ts := httptest.NewServer(th)

	expectedServerRequest := "/generate?ttl=0"
	expectedErrorMessage := "error unmarshaling response"

	client := onetimesecret.NewClient("user", "errorToken", ts.URL)
	_, _, err := client.Generate(0)

	if th.calls[0] != expectedServerRequest {
		t.Fatalf("Expected server request key was %q but got %q", expectedServerRequest, th.calls[0])
	}

	if err != nil && !strings.HasPrefix(err.Error(), expectedErrorMessage) {
		t.Fatalf("Error %q was expected but got %q", expectedErrorMessage, err.Error())
	}
}
