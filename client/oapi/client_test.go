package oapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmyroslav/api-test-demo/tests"
	"testing"
	"time"
)

const baseURL = "https://fakerestapi.azurewebsites.net"

type testCase struct {
	name     string
	testFunc func(*testing.T, *ClientWithResponses)
}

func setupClient(t *testing.T) *ClientWithResponses {
	httpClient := tests.NewHttpClient(t)
	client, err := NewClientWithResponses(baseURL, WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	return client
}

func TestActivities(t *testing.T) {
	scenarios := []testCase{
		{
			name: "get existing activity",
			testFunc: func(t *testing.T, client *ClientWithResponses) {
				ctx := context.Background()
				activity, err := client.GetApiV1ActivitiesIdWithResponse(ctx, 1)
				if err != nil {
					t.Fatalf("getting activity: %v", err)
				}

				require.NotNil(t, activity.ApplicationjsonV10200, "expected activity, got nil")
				require.NotNil(t, activity.ApplicationjsonV10200.Id, "expected ID, got nil")

				assert.Equal(t, int32(1), *activity.ApplicationjsonV10200.Id)
				assert.NotEmpty(t, *activity.ApplicationjsonV10200.Title, "expected non-empty title")
			},
		},
	}

	client := setupClient(t)

	//measure time for each test
	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			tt.testFunc(t, client)
			t.Logf("test %s took %v", tt.name, time.Since(now))
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
	}

	client := setupClient(t)

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}
