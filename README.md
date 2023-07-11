![Unit tests](https://github.com/hunterwilkins2/httptesting/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/hunterwilkins2/httptesting/actions/workflows/linter.yml/badge.svg)

# Go HTTP Testing Library

Go HTTP Testing Library simplies unit testing HTTP routes in Go. httptesting handles build/executing requests, chaining cookies, and asserting response statuses and bodies.

## Getting Started

Use this library with

```sh
$ go get github.com/hunterwilkins2/httptesting
```

Documentation can be found [here](https://pkg.go.dev/github.com/hunterwilkins2/httptesting).

### Examples

#### Simple requests

```go
func TestRoute(t *testing.T) {
  tester := httptesting.New(t, routes()) // routes() returns a http.Handler
  tester.Post("/todo", strings.New(`{"name": "Get Groceries"}`))
  tester.Execute()
  tester.AssertStatusCode(http.StatusCreated)
}
```

#### Chaining requests

```go
func TestSessionCookie(t *testing.T) {
  tester := httptesting.New(t, routes()) // routes() returns a http.Handler
  tester.Post("/login", strings.New(`{"username": "john.doe@gmail.com", "password": "secret_password"}`)) // Creates a session cookie
  tester.Execute()
  tester.AssertStatusCode(http.StatusOK)

  tester.Get("/user") // Chains the session cookie created in the previous request
  tester.Execute()
  tester.AssertStatusCode(http.StatusOK)
  tester.AssertStructDeepEquals(&User{}, &User{
    Username: "john.doe@gmail.com",
  })
}
```

#### Use state from previous requests

```go
func TestGetRoom(t *testing.T) {
  tester := httptesting.New(t, routes()) // routes return a http.Handler
  tester.Post("/room", strings.NewReader(`{"name": "Test Room"}`))
  tester.Execute()
  tester.AssertStatusCode(http.StatusCreated)

  tester.GetWithState(func (s httptesting.State) (url string) {
    room, ok := s.ResponseResult.(Room)
    if !ok {
      t.Fatal("Could not cast response result to type Room")
    }
    return fmt.Sprintf("/room/%d", room.ID)
  })
  tester.Execute()
  tester.AssertStatusCode(http.StatusOK)
  tester.AssertStruct(&Room{}, func (responseBody interface{}) bool {
    room, ok := s.ResponseResult.(Room)
    if !ok {
      t.Fatal("Could not cast response result to type Room")
    }
    return room.Name == "Test Room"
  })
}
```
