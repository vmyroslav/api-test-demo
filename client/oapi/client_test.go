package oapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmyroslav/api-test-demo/tests"
)

const (
	baseURL    = "https://fakerestapi.azurewebsites.net"
	fakePrefix = "fake-api-test"
)

type testCase struct {
	name     string
	testFunc func(*testing.T, *ClientWithResponses)
}

func setupClient(t *testing.T) *ClientWithResponses {
	httpClient := tests.NewHttpClient(t)
	client, err := NewClientWithResponses(baseURL, WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("failed to inistialize client: %v", err)
	}

	return client
}

// Generate a random title to test the possibility of random data in Hoverfly
func randomTitleForResource(resource string) string {
	return fmt.Sprintf("%s-%s-title-%s", fakePrefix, resource, tests.RandomString(5))
}

func TestActivities(t *testing.T) {
	client := setupClient(t)

	scenarios := []testCase{
		{
			name: "get existing activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.GetApiV1ActivitiesIdWithResponse(ctx, 1)
				if err != nil {
					t.Fatalf("getting resp: %v", err)
				}

				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200.Id, "expected ID, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, int32(1), *resp.ApplicationjsonV10200.Id)
				assert.NotEmpty(t, *resp.ApplicationjsonV10200.Title, "expected non-empty title")
			},
		},
		{
			name: "get non-existing activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.GetApiV1ActivitiesIdWithResponse(ctx, 0)

				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode())
			},
		},
		{
			name: "create activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx     = context.Background()
					id      = int32(1)
					title   = randomTitleForResource("activity")
					dueDate = time.Now()
				)

				resp, err := client.PostApiV1ActivitiesWithApplicationJSONV10BodyWithResponse(
					ctx,
					PostApiV1ActivitiesApplicationJSONV10RequestBody{
						Id:        tests.ToPtr(id),
						Title:     tests.ToPtr(title),
						DueDate:   tests.ToPtr(dueDate),
						Completed: tests.ToPtr(false),
					},
				)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp body, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
				assert.Equal(t, title, *resp.ApplicationjsonV10200.Title)
				assert.Equal(
					t,
					dueDate.UTC().Format(time.RFC3339),
					resp.ApplicationjsonV10200.DueDate.UTC().Format(time.RFC3339),
				)
				assert.Equal(t, false, *resp.ApplicationjsonV10200.Completed)
			},
		},
		{
			name: "update activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx     = context.Background()
					id      = int32(10)
					title   = randomTitleForResource("activity")
					dueDate = time.Now()
				)

				resp, err := client.PutApiV1ActivitiesIdWithApplicationJSONV10BodyWithResponse(
					ctx,
					id,
					PutApiV1ActivitiesIdApplicationJSONV10RequestBody{
						Id:        tests.ToPtr(id),
						Title:     tests.ToPtr(title),
						DueDate:   tests.ToPtr(dueDate),
						Completed: tests.ToPtr(true),
					},
				)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp body, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
				assert.Equal(t, title, *resp.ApplicationjsonV10200.Title)
				assert.Equal(
					t,
					dueDate.UTC().Format(time.RFC3339),
					resp.ApplicationjsonV10200.DueDate.UTC().Format(time.RFC3339),
				)
				assert.Equal(t, true, *resp.ApplicationjsonV10200.Completed)
			},
		},
		{
			name: "delete activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.DeleteApiV1ActivitiesIdWithResponse(ctx, 1)

				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode())
			},
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestAuthors(t *testing.T) {
	scenarios := []testCase{
		{
			name: "get author",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx = context.Background()
					id  = int32(1)
				)

				resp, err := client.GetApiV1AuthorsIdWithResponse(ctx, id)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp body, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
				assert.NotEmpty(t, *resp.ApplicationjsonV10200.FirstName, "expected non-empty first name")
				assert.NotEmpty(t, *resp.ApplicationjsonV10200.LastName, "expected non-empty last name")
			},
		},
		{
			name: "get non-existing author",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.GetApiV1AuthorsIdWithResponse(ctx, 0)

				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode())
			},
		},
		{
			name: "create new author",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx       = context.Background()
					id        = int32(1)
					firstName = randomTitleForResource("author")
					lastName  = randomTitleForResource("author")
				)

				resp, err := client.PostApiV1AuthorsWithApplicationJSONV10BodyWithResponse(
					ctx,
					PostApiV1AuthorsApplicationJSONV10RequestBody{
						Id:        tests.ToPtr(id),
						FirstName: tests.ToPtr(firstName),
						LastName:  tests.ToPtr(lastName),
					},
				)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp body, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
				assert.Equal(t, firstName, *resp.ApplicationjsonV10200.FirstName)
				assert.Equal(t, lastName, *resp.ApplicationjsonV10200.LastName)
			},
		},
	}

	client := setupClient(t)

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestBooks(t *testing.T) {
	scenarios := []testCase{
		{
			name: "get book",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx = context.Background()
					id  = int32(1)
				)

				resp, err := client.GetApiV1BooksIdWithResponse(ctx, id)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")
				require.NotNil(t, resp.ApplicationjsonV10200, "expected resp body, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
				assert.NotEmpty(t, *resp.ApplicationjsonV10200.Title, "expected non-empty title")
			},
		},
		{
			name: "get non-existing book",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.GetApiV1BooksIdWithResponse(ctx, 0)

				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode())
			},
		},
		{
			name: "create new book",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx = context.Background()
					id  = int32(1)
					// Generate a random title to test the possibility of random data in Hoverfly
					title = randomTitleForResource("book")
				)

				resp, err := client.PostApiV1BooksWithApplicationJSONV10BodyWithResponse(
					ctx,
					PostApiV1BooksApplicationJSONV10RequestBody{
						Id:    tests.ToPtr(id),
						Title: tests.ToPtr(title),
					},
				)
				require.NoError(t, err)
				require.NotNil(t, resp, "expected resp, got nil")

				assert.Equal(t, http.StatusOK, resp.StatusCode())
			},
		},
	}

	client := setupClient(t)

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestUsers(t *testing.T) {
	scenarios := []testCase{
		{
			name: "PUT user",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				activity, err := client.PutApiV1UsersIdWithApplicationJSONV10BodyWithResponse(
					ctx,
					1,
					PutApiV1UsersIdApplicationJSONV10RequestBody{
						Id:       tests.ToPtr(int32(1)),
						Password: tests.ToPtr("test-password"),
						UserName: tests.ToPtr("test-username"),
					},
				)
				require.NoErrorf(t, err, "putting user: %v", err)
				require.NotNil(t, activity)
			},
		},
		{
			name: "get non-existing user",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				resp, err := client.GetApiV1UsersIdWithResponse(ctx, 0)

				require.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode())
			},
		},
	}

	client := setupClient(t)

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestUUIDPatterns(t *testing.T) {
	client := setupClient(t)

	// Test multiple UUID cases
	scenarios := []testCase{
		{
			name: "create activity with UUID in title",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()

				// Create multiple activities with different UUIDs
				for i := 0; i < 3; i++ {
					id := int32(12356 + i)
					title := uuid.New().String()

					resp, err := client.PostApiV1ActivitiesWithApplicationJSONV10BodyWithResponse(
						ctx,
						PostApiV1ActivitiesApplicationJSONV10RequestBody{
							Id:        tests.ToPtr(id),
							Title:     tests.ToPtr(title),
							DueDate:   tests.ToPtr(time.Now()),
							Completed: tests.ToPtr(false),
						},
					)
					require.NoError(t, err)
					require.NotNil(t, resp)
					require.NotNil(t, resp.ApplicationjsonV10200)
					assert.Equal(t, http.StatusOK, resp.StatusCode())
					assert.NotNil(t, resp.ApplicationjsonV10200)
					assert.Equal(t, title, *resp.ApplicationjsonV10200.Title)
				}
			},
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestDemoEndpointReplacement(t *testing.T) {
	client := setupClient(t)

	// Test multiple UUID cases
	scenarios := []testCase{
		{
			name: "create activity with UUID in title",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()

				resp, err := client.GetApiV1ActivitiesIdWithResponse(ctx, 30)
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.ApplicationjsonV10200)

				assert.Equal(t, http.StatusOK, resp.StatusCode())
				// Uncomment these lines during the simulation to check the static response replacement
				// assert.Equal(t, int32(77), *resp.ApplicationjsonV10200.Id)
				// assert.Equal(t, "John Doe", *resp.ApplicationjsonV10200.Title)
				// assert.Equal(t, false, *resp.ApplicationjsonV10200.Completed)
			},
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}

func TestActivitiesWithRetries(t *testing.T) {
	serverUrl, ok := os.LookupEnv("API_SERVER_URL")
	if !ok || serverUrl == "" {
		ts := tests.NewTestServer() // Start our test server
		defer ts.Close()

		serverUrl = ts.URL()
	}

	httpClient := tests.NewHttpClient(t)
	client, err := NewClientWithResponses(serverUrl, WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("failed to inistialize client: %v", err)
	}
	defer func() {
		// Reset the server state
		_, _ = httpClient.Post(serverUrl+"/reset", "application/json", nil)
	}()

	scenarios := []testCase{
		{
			name: "create activity with retries",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				var (
					ctx   = context.Background()
					id    = int32(7474)
					title = "Test Activity"
				)

				var lastErr error
				for attempt := 1; attempt <= 4; attempt++ {
					resp, err := client.GetApiV1ActivitiesIdWithResponse(
						ctx,
						id,
					)
					if err != nil {
						lastErr = err
						continue
					}

					// Success!
					if resp.StatusCode() == http.StatusOK {
						require.NotNil(t, resp.ApplicationjsonV10200)
						assert.Equal(t, id, *resp.ApplicationjsonV10200.Id)
						assert.Equal(t, title, *resp.ApplicationjsonV10200.Title)
						return
					}

					// Continue if it's a retryable status
					if resp.StatusCode() == http.StatusServiceUnavailable ||
						resp.StatusCode() == http.StatusGatewayTimeout {
						continue
					}

					// Non-retryable error
					t.Fatalf("Got non-retryable status: %d", resp.StatusCode())
				}

				// If we got here, all retries failed
				t.Fatalf("All retries failed. Last error: %v", lastErr)
			},
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}
