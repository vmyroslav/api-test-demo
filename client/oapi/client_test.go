package oapi

import (
	"context"
	"github.com/vmyroslav/api-test-demo/tests"
	"testing"
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
