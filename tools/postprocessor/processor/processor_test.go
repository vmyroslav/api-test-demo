package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultProcessor_ProcessSimulation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		simulation   *Simulation
		expectedBody string
		expectedResp string
		expectChange bool
	}{
		{
			name: "processes request with fake-api-test pattern",
			simulation: &Simulation{
				Data: struct {
					Pairs []Pair `json:"pairs"`
				}{
					Pairs: []Pair{
						{
							Request: Request{
								Body: []Matcher{
									{
										Matcher: "exact",
										Value:   `{"title":"fake-api-test-activity-title-12345"}`,
									},
								},
							},
							Response: Response{
								Status: 200,
								Body:   `{"title":"fake-api-test-activity-title-12345"}`,
							},
						},
					},
				},
			},
			expectedBody: `{"title":"fake-api-test-activity-title-*"}`,
			expectedResp: `{"title":"{{ Request.Body 'jsonpath' '$.title' }}"}`,
			expectChange: true,
		},
		{
			name: "processes request with timestamp",
			simulation: &Simulation{
				Data: struct {
					Pairs []Pair `json:"pairs"`
				}{
					Pairs: []Pair{
						{
							Request: Request{
								Body: []Matcher{
									{
										Matcher: "exact",
										Value:   `{"dueDate":"2024-12-29T21:15:18Z"}`,
									},
								},
							},
							Response: Response{
								Status: 200,
								Body:   `{"dueDate":"2024-12-29T21:15:18Z"}`,
							},
						},
					},
				},
			},
			expectedBody: `{"dueDate":"*"}`,
			expectedResp: `{"dueDate":"{{ Request.Body 'jsonpath' '$.dueDate' }}"}`,
			expectChange: true,
		},
		{
			name: "does not modify request without patterns",
			simulation: &Simulation{
				Data: struct {
					Pairs []Pair `json:"pairs"`
				}{
					Pairs: []Pair{
						{
							Request: Request{
								Body: []Matcher{
									{
										Matcher: "exact",
										Value:   `{"title":"normal-title"}`,
									},
								},
							},
							Response: Response{
								Status: 200,
								Body:   `{"title":"normal-title"}`,
							},
						},
					},
				},
			},
			expectedBody: `{"title":"normal-title"}`,
			expectedResp: `{"title":"normal-title"}`,
			expectChange: false,
		},
	}

	processor := NewDefaultProcessor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := processor.Process(tt.simulation)
			require.NoError(t, err)

			if tt.expectChange {
				assert.Equal(t, "glob", tt.simulation.Data.Pairs[0].Request.Body[0].Matcher)
				assert.Equal(t, tt.expectedBody, tt.simulation.Data.Pairs[0].Request.Body[0].Value)
				assert.Equal(t, tt.expectedResp, tt.simulation.Data.Pairs[0].Response.Body)
				assert.True(t, tt.simulation.Data.Pairs[0].Response.Templated)
			} else {
				assert.Equal(t, "exact", tt.simulation.Data.Pairs[0].Request.Body[0].Matcher)
				assert.Equal(t, tt.expectedBody, tt.simulation.Data.Pairs[0].Request.Body[0].Value)
				assert.Equal(t, tt.expectedResp, tt.simulation.Data.Pairs[0].Response.Body)
				assert.False(t, tt.simulation.Data.Pairs[0].Response.Templated)
			}
		})
	}
}
