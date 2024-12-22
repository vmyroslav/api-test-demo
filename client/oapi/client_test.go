package oapi

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/vmyroslav/api-test-demo/tests"
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

				if activity.ApplicationjsonV10200 == nil {
					t.Fatal("expected activity, got nil")
				}

				if activity.ApplicationjsonV10200.Id == nil || *activity.ApplicationjsonV10200.Id != 1 {
					t.Errorf("expected ID 1, got %v", activity.ApplicationjsonV10200.Id)
				}

				if activity.ApplicationjsonV10200.Title == nil || *activity.ApplicationjsonV10200.Title == "" {
					t.Error("expected non-empty title")
				}
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
	}

	client := setupClient(t)

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, client)
		})
	}
}
