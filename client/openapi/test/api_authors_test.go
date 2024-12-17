/*
FakeRESTApi.Web V1

Testing AuthorsAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package openapi

import (
	"context"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_openapi_AuthorsAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test AuthorsAPIService ApiV1AuthorsAuthorsBooksIdBookGet", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var idBook int32

		resp, httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsAuthorsBooksIdBookGet(context.Background(), idBook).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test AuthorsAPIService ApiV1AuthorsGet", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsGet(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test AuthorsAPIService ApiV1AuthorsIdDelete", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsIdDelete(context.Background(), id).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test AuthorsAPIService ApiV1AuthorsIdGet", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		resp, httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsIdGet(context.Background(), id).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test AuthorsAPIService ApiV1AuthorsIdPut", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		resp, httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsIdPut(context.Background(), id).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test AuthorsAPIService ApiV1AuthorsPost", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.AuthorsAPI.ApiV1AuthorsPost(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
