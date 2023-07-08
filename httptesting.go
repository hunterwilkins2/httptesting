// Package httptesting Go HTTP testing library. Simplifies chaining and asserting HTTP REST calls for unit and integration testing
package httptesting

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	urlpkg "net/url"
	"reflect"

	"github.com/hunterwilkins2/httptesting/internal/util"
)

// State stores the state of the previous request and is used for chaining headers and cookies for each request.
// Use *WithState functions to get access to the state. Helpful when the server generates an ID or uuid that needs to be chained with another request
type State struct {
	// Request previous http request made
	Request *http.Request

	// Response previous http response result
	Response *http.Response

	// ResponseResult stores the value of the decoded json body from the response result
	// ResponseResult will be nil until AssertStruct or AssertStructDeepEquals is called
	ResponseResult interface{}

	// Values key-value store to save values needed later in the test
	Values map[string]any
}

// Httptester struct for chaining REST calls together
// Uses the builder pattern for constructing and chaining requests
type Httptester struct {
	t util.TestingT
	// handler http.Handler to run tests against
	handler http.Handler
	// state internal State used for chaining requests
	state State

	// requestExecuted is set to true then Execute() is called,
	// and set back to false when a new request is initialized
	// If Execute() is not called before an assertion is made then the test will fail
	requestExecuted bool
}

// New returns a new httptester. Create a new httptester for each test for concurrent use
func New(t util.TestingT, h http.Handler) *Httptester {
	return &Httptester{
		t:       t,
		handler: h,
		state: State{
			Values: make(map[string]any),
		},
	}
}

// getRequest helper function for getting the current state of the request being build
func (ht *Httptester) getRequest() *http.Request {
	ht.requestExecuted = false
	ht.state.ResponseResult = nil
	if ht.state.Request == nil {
		ht.state.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	}
	return ht.state.Request
}

// setBodyReader helper function to convert an io.Reader to an io.ReadCloser to set the body of the request
func (ht *Httptester) setBodyReader(reader io.Reader) {
	rc, ok := reader.(io.ReadCloser)
	if !ok && reader != nil {
		rc = io.NopCloser(reader)
	}
	ht.getRequest().Body = rc
}

// NewRequest creates a new httptester Request the same as http.NewRequest
func (ht *Httptester) NewRequest(method string, url string, reader io.Reader) {
	var err error
	req := ht.getRequest()
	req.Method = method
	req.URL, err = urlpkg.Parse(url)
	if err != nil {
		ht.t.Fatalf(err.Error())
	}
	ht.setBodyReader(reader)
}

// NewRequestWithState creates a new httptester Request the same as http.NewRequest.
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) NewRequestWithState(f func(s State) (method string, url string, reader io.Reader)) {
	ht.NewRequest(f(ht.state))
}

// Get creates a new Get request
func (ht *Httptester) Get(url string) {
	ht.NewRequest(http.MethodGet, url, nil)
}

// GetWithState creates a new Get request.
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) GetWithState(f func(s State) (url string)) {
	ht.Get(f(ht.state))
}

// Post creates a new Post request with a url and request body
func (ht *Httptester) Post(url string, reader io.Reader) {
	ht.NewRequest(http.MethodPost, url, reader)
}

// PostWithState creates a new Post request with a url and request body
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) PostWithState(f func(s State) (url string, reader io.Reader)) {
	ht.Post(f(ht.state))
}

// Put creates a new Put request with a url and request body
func (ht *Httptester) Put(url string, reader io.Reader) {
	ht.NewRequest(http.MethodPut, url, reader)
}

// PutWithState creates a new Put request with a url and request body
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) PutWithState(f func(s State) (url string, reader io.Reader)) {
	ht.Put(f(ht.state))
}

// Patch creates a new Patch request with a url and request body
func (ht *Httptester) Patch(url string, reader io.Reader) {
	ht.NewRequest(http.MethodPatch, url, reader)
}

// PatchWithState creates a new Patch request with a url and request body
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) PatchWithState(f func(s State) (url string, reader io.Reader)) {
	ht.Patch(f(ht.state))
}

// Delete creates a new Delete request
func (ht *Httptester) Delete(url string) {
	ht.NewRequest(http.MethodDelete, url, nil)
}

// DeleteWithState creates a new Delete request
// Takes a func of the current state and returns the parameters for a NewRequest
func (ht *Httptester) DeleteWithState(f func(s State) (url string)) {
	ht.Delete(f(ht.state))
}

// SetBody sets the body of the current request
func (ht *Httptester) SetBody(reader io.Reader) {
	ht.setBodyReader(reader)
}

// SetRequestBodyJSON encodes the struct passed in as JSON and sets the resulting []byte as the request body
func (ht *Httptester) SetRequestBodyJSON(body interface{}) {
	jsonBody, err := util.EncodeJSON(&body)
	if err != nil {
		ht.t.Fatalf("Error encoding request body: %s", err.Error())
	}
	ht.setBodyReader(bytes.NewReader(jsonBody))
}

// SetBodyWithState encodes the struct passed in as JSON and sets the resulting []byte as the request body.
// Able to use the values from previous request to update the body
func (ht *Httptester) SetBodyWithState(f func(s State) (reader io.Reader)) {
	ht.SetBody(f(ht.state))
}

// AddHeader adds a header to the current request
func (ht *Httptester) AddHeader(key, value string) {
	ht.getRequest().Header.Set(key, value)
}

// AddHeaderWithState adds a header to the current request.
// Able to use the values from previous requests to create a new header
func (ht *Httptester) AddHeaderWithState(f func(s State) (key, value string)) {
	ht.AddHeader(f(ht.state))
}

// AddCookie adds a cookie to the current request. This cookie will be chained through all subsuquent requests made.
func (ht *Httptester) AddCookie(cookie *http.Cookie) {
	ht.getRequest().AddCookie(cookie)
}

// AddCookieWithState adds a cookie to the current request. This cookie will be chained through all subsuquent requests made.
// Able to use the values from previous requests to create the cookie
func (ht *Httptester) AddCookieWithState(f func(s State) *http.Cookie) {
	ht.AddCookie(f(ht.state))
}

// SetValue sets a value in State to be referenced later
func (ht *Httptester) SetValue(key string, value any) {
	ht.state.Values[key] = value
}

// SetValueWithState get access to the current state of the store to set a value to be referenced later
func (ht *Httptester) SetValueWithState(f func(s State) (key string, value any)) {
	key, value := f(ht.state)
	ht.state.Values[key] = value
}

// Execute executes the current request that was build and resets the state of Response and ResponseResult.
// This method must be called before any assertions are made.
func (ht *Httptester) Execute() {
	if ht.state.Response != nil {
		for _, cookie := range ht.state.Response.Cookies() {
			ht.state.Request.AddCookie(cookie)
		}
	}
	response := httptest.NewRecorder()
	ht.handler.ServeHTTP(response, ht.getRequest())

	ht.requestExecuted = true
	ht.state.Response = response.Result()
	ht.state.Request = nil
}

// assertRequestExecuted helper fuction to assert the current request was executed
func (ht *Httptester) assertRequestExecuted() {
	if !ht.requestExecuted {
		ht.t.Fatalf("Request %q was not executed", ht.getRequest().URL.String())
	}
}

// AssertStatus asserts the status of the response to the previous request
func (ht *Httptester) AssertStatus(expectedStatus string) {
	ht.assertRequestExecuted()
	if ht.state.Response.Status != expectedStatus {
		ht.t.Fatalf("Expected status %q; got %q", ht.state.Response.Status, expectedStatus)
	}
}

// AssertStatusCode asserts the status code of the response to the previous request
func (ht *Httptester) AssertStatusCode(statusCode int) {
	ht.assertRequestExecuted()
	if ht.state.Response.StatusCode != statusCode {
		ht.t.Fatalf("Expected %d; got %d", ht.state.Response.StatusCode, statusCode)
	}
}

// AssertHeader asserts the headers of the response to the previous request contains the expected key and value
func (ht *Httptester) AssertHeader(key, expectedValue string) {
	ht.assertRequestExecuted()
	if ht.state.Response.Header.Get(key) != expectedValue {
		ht.t.Fatalf("Expected %q; got %q", ht.state.Response.Header.Get(key), expectedValue)
	}
}

// getCookie helper function to find a cookie by its name
func getCookie(cookies []*http.Cookie, wantCookie string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == wantCookie {
			return cookie
		}
	}
	return nil
}

// AssertCookieExists asserts that a cookie exists in the response to the previous request with the name of cookieName
func (ht *Httptester) AssertCookieExists(cookieName string) {
	ht.assertRequestExecuted()
	if getCookie(ht.state.Response.Cookies(), cookieName) == nil {
		ht.t.Fatalf("Expected to find cookie %q", cookieName)
	}
}

// AssertCookieValue asserts that a cookie exists and its value is expectedValue in the response to the previous request
func (ht *Httptester) AssertCookieValue(cookieName, expectedValue string) {
	ht.assertRequestExecuted()
	cookie := getCookie(ht.state.Response.Cookies(), cookieName)
	if cookie == nil {
		ht.t.Fatalf("Expected to find cookie %q", cookieName)
	}
	if cookie != nil && cookie.Value != expectedValue {
		ht.t.Fatalf("Expected cookie to have value of %q; got %q", expectedValue, cookie.Value)
	}
}

// AssertCookieDeepEquals asserts that a cookie exists and it deep equals expectedCookie in the response to the previous request
func (ht *Httptester) AssertCookieDeepEquals(expectedCookie *http.Cookie) {
	ht.assertRequestExecuted()
	if expectedCookie == nil {
		ht.t.Fatalf("Expected cookie is nil")
	}
	var cookieName string
	if expectedCookie != nil {
		cookieName = expectedCookie.Name
	}
	if cookieName == "" {
		ht.t.Fatalf("Expected cookie cannot have an empty Name")
	}
	cookie := getCookie(ht.state.Response.Cookies(), cookieName)
	if cookie == nil {
		ht.t.Fatalf("Expected to find cookie %q", cookieName)
	}
	if cookie.String() != expectedCookie.String() {
		ht.t.Fatalf("Expected %v; got %v", expectedCookie, cookie)
	}
}

// AssertBody asserts the body of the response to the previous request matches the []byte provided
func (ht *Httptester) AssertBody(body []byte) {
	ht.assertRequestExecuted()
	resBody, err := io.ReadAll(ht.state.Response.Body)
	if err != nil {
		ht.t.Fatalf(err.Error())
	}
	if string(resBody) != string(body) {
		ht.t.Fatalf("Expected %s; got %s", resBody, body)
	}
}

// AssertStruct decodes the JSON response body into r and asserts the predicate passed in
func (ht *Httptester) AssertStruct(r interface{}, predicate func(responseBody interface{}) bool) {
	ht.assertRequestExecuted()
	err := util.DecodeJSON(ht.state.Response, &r)
	if err != nil {
		ht.t.Fatalf("Error parsing response json: %s", err.Error())
	}
	ht.state.ResponseResult = r
	if !predicate(r) {
		ht.t.Fatalf("Response body was not equal to predicate")
	}
}

// AssertStructDeepEquals decodes the JSON response body into r and asserts r is deeply equatable to expected
func (ht *Httptester) AssertStructDeepEquals(r interface{}, expected interface{}) {
	ht.assertRequestExecuted()
	err := util.DecodeJSON(ht.state.Response, &r)
	if err != nil {
		ht.t.Fatalf("Error parsing response json: %s", err.Error())
	}
	ht.state.ResponseResult = r
	if !reflect.DeepEqual(r, expected) {
		ht.t.Fatalf("Expected %v; got %v", expected, r)
	}
}
