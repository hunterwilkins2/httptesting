package httptesting

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hunterwilkins2/httptesting/internal/util"
)

func assertRequest(t *testing.T, ht *Httptester, expectedMethod string, expectedPath string, expectedBody io.Reader) {
	t.Helper()
	if ht.state.Request.Method != expectedMethod {
		t.Errorf("Expected method to be %s; got %s", expectedMethod, ht.state.Request.Method)
	}
	if ht.state.Request.URL.Path != expectedPath {
		t.Errorf("Expected path to be%s; got %s;", expectedPath, ht.state.Request.URL.Path)
	}
	if expectedBody == nil && ht.state.Request.Body != nil {
		t.Errorf("Expected body to be nil")
	} else if expectedBody != nil && ht.state.Request.Body == nil {
		t.Errorf("Expected body to not be nil")
	}
}

func assertBody(t *testing.T, body io.ReadCloser, expectedBody string) {
	b, err := io.ReadAll(body)
	if err := body.Close(); err != nil {
		t.Fatalf("Error closing body: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Unexpected error reading body: %s", err.Error())
		return
	}
	if string(b) != expectedBody {
		t.Errorf("Expected body to contain %s; got %s", expectedBody, string(b))
	}
}

func assertFatal(t *testing.T) {
	if err := recover(); err == nil {
		t.Errorf("Expected Fatalf to be called during test.")
	}
}

func TestNewRequest(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPost
	expectedPath := "/test"
	expectedBody := strings.NewReader("test body")

	tester.NewRequest(expectedMethod, expectedPath, expectedBody)
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestNewRequestURLErrorHandling(t *testing.T) {
	t.Parallel()
	mockT := util.MockTestingT{}
	tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPost
	expectedPath := ":"
	expectedBody := strings.NewReader("test body")

	defer assertFatal(t)
	tester.NewRequest(expectedMethod, expectedPath, expectedBody)
}

func TestNewRequestWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPost
	expectedPath := "/test/123"
	expectedBody := strings.NewReader("test body")

	tester.SetValue("id", 123)
	tester.NewRequestWithState(func(s State) (method string, url string, reader io.Reader) {
		return expectedMethod, fmt.Sprintf("/test/%d", s.Values["id"]), expectedBody
	})
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestGet(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedPath := "/test"
	tester.Get(expectedPath)
	assertRequest(t, tester, http.MethodGet, expectedPath, nil)
}

func TestGetWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedPath := "/test/123"
	tester.SetValue("id", 123)
	tester.GetWithState(func(s State) (url string) {
		return fmt.Sprintf("/test/%d", s.Values["id"])
	})
	assertRequest(t, tester, http.MethodGet, expectedPath, nil)
}

func TestPost(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPost
	expectedPath := "/test"
	expectedBody := strings.NewReader("test body")

	tester.Post(expectedPath, expectedBody)
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestPostWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPost
	expectedPath := "/test/123"
	expectedBody := strings.NewReader("test body")

	tester.SetValue("id", 123)
	tester.PostWithState(func(s State) (url string, reader io.Reader) {
		return fmt.Sprintf("/test/%d", s.Values["id"]), expectedBody
	})
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestPut(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPut
	expectedPath := "/test"
	expectedBody := strings.NewReader("test body")

	tester.Put(expectedPath, expectedBody)
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestPutWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPut
	expectedPath := "/test/123"
	expectedBody := strings.NewReader("test body")

	tester.SetValue("id", 123)
	tester.PutWithState(func(s State) (url string, reader io.Reader) {
		return fmt.Sprintf("/test/%d", s.Values["id"]), expectedBody
	})
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestPatch(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPatch
	expectedPath := "/test"
	expectedBody := strings.NewReader("test body")

	tester.Patch(expectedPath, expectedBody)
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestPatchWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedMethod := http.MethodPatch
	expectedPath := "/test/123"
	expectedBody := strings.NewReader("test body")

	tester.SetValue("id", 123)
	tester.PatchWithState(func(s State) (url string, reader io.Reader) {
		return fmt.Sprintf("/test/%d", s.Values["id"]), expectedBody
	})
	assertRequest(t, tester, expectedMethod, expectedPath, expectedBody)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedPath := "/test"

	tester.Delete(expectedPath)
	assertRequest(t, tester, http.MethodDelete, expectedPath, nil)
}

func TestDeleteWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedPath := "/test/123"

	tester.SetValue("id", 123)
	tester.DeleteWithState(func(s State) (url string) {
		return fmt.Sprintf("/test/%d", s.Values["id"])
	})
	assertRequest(t, tester, http.MethodDelete, expectedPath, nil)
}

func TestSetBody(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedBody := "test body"

	tester.SetBody(strings.NewReader(expectedBody))
	assertBody(t, tester.state.Request.Body, expectedBody)
}

func TestRequestBodyJSON(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedBody := `{"name":"Bob"}`

	tester.SetRequestBodyJSON(struct {
		Name string `json:"name"`
	}{
		Name: "Bob",
	})
	assertBody(t, tester.state.Request.Body, expectedBody)
}

func TestSetRequestBodyJSONErrorHandling(t *testing.T) {
	t.Parallel()
	mockT := util.MockTestingT{}
	tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	defer assertFatal(t)
	c := make(chan int)
	tester.SetRequestBodyJSON(c)
}

func TestSetBodyWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedBody := `test 123`

	tester.SetValue("body val", 123)
	tester.SetBodyWithState(func(s State) (reader io.Reader) {
		return strings.NewReader(fmt.Sprintf("test %d", s.Values["body val"]))
	})
	assertBody(t, tester.state.Request.Body, expectedBody)
}

func TestAddHeader(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedHeader := "application/json"
	tester.AddHeader("Content-Type", expectedHeader)
	gotHeader := tester.state.Request.Header.Get("Content-Type")
	if gotHeader != expectedHeader {
		t.Errorf("Expected header %s; got %s", expectedHeader, gotHeader)
	}
}

func TestAddHeaderWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedHeader := "https://foo.example"
	tester.SetValue("OriginURL", expectedHeader)
	tester.AddHeaderWithState(func(s State) (key string, value string) {
		return "Origin", fmt.Sprint(s.Values["OriginURL"])
	})
	gotHeader := tester.state.Request.Header.Get("Origin")
	if gotHeader != expectedHeader {
		t.Errorf("Expected header %s; got %s", expectedHeader, gotHeader)
	}
}

func TestAddCookie(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedCookie := http.Cookie{
		Name:  "TestCookie",
		Value: "1234",
		Path:  "/",
	}
	tester.AddCookie(&expectedCookie)
	gotCookie := getCookie(tester.state.Request.Cookies(), "TestCookie")
	if gotCookie.String() == expectedCookie.String() {
		t.Errorf("Expected cookie %v; got %v", expectedCookie, gotCookie)
	}
}

func TestAddCookieWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	expectedCookie := http.Cookie{
		Name:  "TestCookie",
		Value: "1234",
		Path:  "/",
	}
	tester.SetValue("path", "/")
	tester.AddCookieWithState(func(s State) *http.Cookie {
		return &http.Cookie{
			Name:  "TestCookie",
			Value: "1234",
			Path:  fmt.Sprint(s.Values["path"]),
		}
	})
	gotCookie := getCookie(tester.state.Request.Cookies(), "TestCookie")
	if gotCookie.String() == expectedCookie.String() {
		t.Errorf("Expected cookie %v; got %v", expectedCookie, gotCookie)
	}
}

func TestSetValue(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	tester.SetValue("key", "id")
	value, ok := tester.state.Values["key"].(string)
	if !ok {
		t.Errorf("Did not find key")
		return
	}
	if value != "id" {
		t.Errorf("Expected id; got %s", value)
	}
}

func TestSetValueWithState(t *testing.T) {
	t.Parallel()
	tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	tester.Get("/route")
	tester.SetValueWithState(func(s State) (key string, value any) {
		return "path", s.Request.URL.Path
	})
	value, ok := tester.state.Values["path"].(string)
	if !ok {
		t.Errorf("Did not find key path")
		return
	}
	if value != "/route" {
		t.Errorf("Expected id; got %s", value)
	}
}

func TestExecute(t *testing.T) {
	t.Parallel()

	t.Run("Execute executes request", func(t *testing.T) {
		t.Parallel()
		tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertStatusCode(http.StatusOK)
		tester.AssertBody([]byte("Ok"))
	})

	t.Run("Execute chains cookies from previous request", func(t *testing.T) {
		t.Parallel()
		tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()

			mux.Handle("/set-cookie", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cookie := http.Cookie{
					Name:  "TestCookie",
					Value: "123",
				}
				http.SetCookie(w, &cookie)
				_, err := w.Write([]byte("cookie set"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			mux.Handle("/assert-cookie", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cookie := getCookie(r.Cookies(), "TestCookie")
				if cookie == nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				_, err := w.Write([]byte("cookie retrieved"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))

			mux.ServeHTTP(w, r)
		}))

		tester.Get("/set-cookie")
		tester.Execute()
		tester.AssertStatusCode(http.StatusOK)
		tester.AssertBody([]byte("cookie set"))
		tester.AssertCookieExists("TestCookie")
		tester.Get("/assert-cookie")
		tester.Execute()
		tester.AssertStatusCode(http.StatusOK)
	})

	t.Run("Execute resets state", func(t *testing.T) {
		t.Parallel()
		tester := New(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()

		if !tester.requestExecuted {
			t.Errorf("Expected 'requestExecuted' to be true")
		}
		if tester.state.Response == nil {
			t.Errorf("Expected state.Response to not be nil")
		}
		if tester.state.Request != nil {
			t.Errorf("Expected request to be set to nil")
		}
	})
}

func TestAssertStatus(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertStatus("200 OK")
	})

	t.Run("test assertion fails", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStatus("201 Created")
	})

	t.Run("test assertion successeds", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertStatus("200 OK")
	})
}

func TestAssertStatusCode(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertStatusCode(http.StatusOK)
	})

	t.Run("test assertion fails", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStatusCode(http.StatusCreated)
	})

	t.Run("test assertion successeds", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertStatusCode(http.StatusOK)
	})
}

func TestAssertHeader(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertHeader("Content-Type", "application/json")
	})

	t.Run("test assertion fails", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertHeader("Content-Type", "application/json")
	})

	t.Run("test assertion successeds", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertHeader("Content-Type", "application/json")
	})
}

func TestAssertCookieExists(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertCookieExists("TestCookie")
	})

	t.Run("test cookie is not found fails test", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieExists("TestCookie")
	})

	t.Run("test cookie is found passes assertion", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:  "TestCookie",
				Value: "123",
			})
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieExists("TestCookie")
	})
}

func TestAssertCookieValue(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertCookieValue("TestCookie", "123")
	})

	t.Run("test cookie is not found fails test", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieValue("TestCookie", "123")
	})

	t.Run("test value assertion fails then fail test", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:  "TestCookie",
				Value: "456",
			})
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieValue("TestCookie", "123")
	})

	t.Run("test cookie is found passes assertion", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:  "TestCookie",
				Value: "123",
			})
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieValue("TestCookie", "123")
	})
}

func TestAssertCookieDeepEquals(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		cookie := http.Cookie{
			Name:  "TestCookie",
			Value: "123",
		}
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertCookieDeepEquals(&cookie)
	})

	t.Run("test expected cookie is nil then fail the test", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieDeepEquals(nil)
	})

	t.Run("test if expected cookie name is empty then fail the test", func(t *testing.T) {
		t.Parallel()
		cookie := http.Cookie{
			Value: "123",
		}
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieDeepEquals(&cookie)
	})

	t.Run("test cookie not found fails test", func(t *testing.T) {
		t.Parallel()
		cookie := http.Cookie{
			Name:  "TestCookie",
			Value: "123",
		}
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieDeepEquals(&cookie)
	})

	t.Run("test cookie fails deep assertion", func(t *testing.T) {
		t.Parallel()
		cookie := http.Cookie{
			Name:  "TestCookie",
			Value: "123",
		}
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:  "TestCookie",
				Value: "456",
			})
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieDeepEquals(&cookie)
	})

	t.Run("test cookie is found passes assertion", func(t *testing.T) {
		t.Parallel()
		cookie := http.Cookie{
			Name:  "TestCookie",
			Value: "123",
		}
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &cookie)
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertCookieDeepEquals(&cookie)
	})
}

type mockReadCloser struct{}

func (m *mockReadCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("read fail")
}

func (m *mockReadCloser) Close() error {
	return nil
}

func TestAssertBody(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertBody([]byte("Ok"))
	})

	t.Run("test empty body fails test", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.state.Response.Body = &mockReadCloser{}
		tester.AssertBody([]byte("Ok"))
	})

	t.Run("test response body does not equal expectd", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertBody([]byte("Not equal"))
	})

	t.Run("test response body are equal", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		tester.Get("/get")
		tester.Execute()
		tester.AssertBody([]byte("Ok"))
	})
}

type testStruct struct {
	Value string `json:"value"`
}

func TestAssertStruct(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertStruct(&testStruct{}, func(responseBody interface{}) bool {
			return true
		})
	})

	t.Run("test decode json error handling", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := make(chan string)
		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStruct(expected, func(responseBody interface{}) bool {
			return true
		})
	})

	t.Run("test predicate fails", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"value": "123"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStruct(&testStruct{}, func(responseBody interface{}) bool {
			return false
		})
	})

	t.Run("test predicate passes", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"value": "123"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		var expected testStruct
		tester.Get("/get")
		tester.Execute()
		tester.AssertStruct(&expected, func(responseBody interface{}) bool {
			return true
		})

		if expected.Value != "123" {
			t.Fatalf("Expected %s; got %s", "123", expected.Value)
		}
	})
}

func TestAssertStructDeepEquals(t *testing.T) {
	t.Parallel()
	t.Run("test execute must be called before assert", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Ok"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := &testStruct{
			Value: "123",
		}
		defer assertFatal(t)
		tester.Get("/get")
		tester.AssertStructDeepEquals(&testStruct{}, expected)
	})

	t.Run("test decode json error handling", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"value": "123"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := testStruct{
			Value: "123",
		}

		receiver := make(chan int)
		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStructDeepEquals(&receiver, expected)
	})

	t.Run("test deep equals fails", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"value": "123"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := testStruct{
			Value: "456",
		}
		defer assertFatal(t)
		tester.Get("/get")
		tester.Execute()
		tester.AssertStructDeepEquals(&testStruct{}, &expected)
	})

	t.Run("test deep equals passes", func(t *testing.T) {
		t.Parallel()
		mockT := util.MockTestingT{}
		tester := New(&mockT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"value": "123"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))

		expected := testStruct{
			Value: "123",
		}
		tester.Get("/get")
		tester.Execute()
		tester.AssertStructDeepEquals(&testStruct{}, &expected)
	})
}
