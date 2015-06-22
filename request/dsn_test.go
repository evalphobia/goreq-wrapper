package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var record map[string]string

func TestCall(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	// GET request
	qs := url.Values{}
	qs.Add("hoge", "foo")
	r, err := DSN{
		Method: "GET",
		Uri:    ts.URL,
		Params: qs,
	}.Call()
	assert.Nil(err)
	assert.Equal("GET", record["method"])
	assert.Equal("", record["content_type"])
	assert.Equal("", record["body"])
	assert.Equal("hoge=foo", record["query"])

	resp, err := r.ToMap()
	assert.Nil(err)
	assert.False(resp["error"].(bool))

	// POST request
	r, err = DSN{
		Method: "POST",
		Uri:    ts.URL,
		Params: "hoge=foo",
	}.Call()

	assert.Nil(err)
	assert.Equal("POST", record["method"])
	assert.Equal("application/x-www-form-urlencoded", record["content_type"])
	assert.Equal("hoge=foo", record["body"])
	assert.Equal("", record["query"])

	resp, err = r.ToMap()
	assert.Nil(err)
	assert.False(resp["error"].(bool))
}

func TestGET(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	// GET request
	qs := url.Values{}
	qs.Add("hoge", "foo")
	r, err := DSN{
		Uri:    ts.URL,
		Params: qs,
	}.GET()
	assert.Nil(err)
	assert.Equal("GET", record["method"])
	assert.Equal("", record["content_type"])
	assert.Equal("", record["body"])
	assert.Equal("hoge=foo", record["query"])

	resp, err := r.ToMap()
	assert.Nil(err)
	assert.False(resp["error"].(bool))
}


func TestPOST(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	// POST request
	r, err := DSN{
		Uri:    ts.URL,
		Params: "hoge=foo",
	}.POST()

	assert.Nil(err)
	assert.Equal("POST", record["method"])
	assert.Equal("application/x-www-form-urlencoded", record["content_type"])
	assert.Equal("hoge=foo", record["body"])
	assert.Equal("", record["query"])

	resp, err := r.ToMap()
	assert.Nil(err)
	assert.False(resp["error"].(bool))
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	record = make(map[string]string)
	body, err := ioutil.ReadAll(r.Body)

	record["body"] = string(body)
	record["query"] = r.URL.Query().Encode()
	record["method"] = r.Method
	record["content_type"] = r.Header.Get("Content-Type")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, `{"error": true}`)
	} else {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"error": false}`)
	}
	return
}
