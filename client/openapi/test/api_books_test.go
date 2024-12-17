/*
FakeRESTApi.Web V1

Testing BooksAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package openapi

import (
	"context"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_openapi_BooksAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	configuration.Host = "fakerestapi.azurewebsites.net"
	configuration.Scheme = "https"
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test BooksAPIService ApiV1BooksGet", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.BooksAPI.ApiV1BooksGet(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

		assert.Equal(t, 200, len(resp))
		assert.Equal(t, int32(100), resp[1].GetPageCount()-resp[0].GetPageCount())
	})

	t.Run("Test BooksAPIService ApiV1BooksIdDelete", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		httpRes, err := apiClient.BooksAPI.ApiV1BooksIdDelete(context.Background(), id).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test BooksAPIService ApiV1BooksIdGet", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		resp, httpRes, err := apiClient.BooksAPI.ApiV1BooksIdGet(context.Background(), id).Execute()

		if id == 0 {
			assert.Error(t, err)
			assert.Equal(t, 404, httpRes.StatusCode)
		} else {
			require.Nil(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, httpRes.StatusCode)

			assert.Equal(t, int32(1), resp.GetId())
			assert.Equal(t, int32(1)*100, resp.GetPageCount())
		}

	})

	t.Run("Test BooksAPIService ApiV1BooksIdPut", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var id int32

		var book = openapiclient.Book{
			Id:          openapiclient.PtrInt32(id),
			Title:       *openapiclient.NewNullableString(nil),
			Description: *openapiclient.NewNullableString(nil),
			PageCount:   openapiclient.PtrInt32(0),
			Excerpt:     *openapiclient.NewNullableString(nil),
			PublishDate: openapiclient.PtrTime(time.Now()),
		}

		httpRes, err := apiClient.BooksAPI.ApiV1BooksIdPut(context.Background(), id).
			Book(book).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test BooksAPIService ApiV1BooksPost", func(t *testing.T) {

		//t.Skip("skip test") // remove to run test

		var book = openapiclient.Book{
			Id:          openapiclient.PtrInt32(1),
			Title:       *openapiclient.NewNullableString(nil),
			Description: *openapiclient.NewNullableString(nil),
			PageCount:   openapiclient.PtrInt32(0),
			Excerpt:     *openapiclient.NewNullableString(nil),
			PublishDate: openapiclient.PtrTime(time.Now()),
		}

		httpRes, err := apiClient.BooksAPI.ApiV1BooksPost(context.Background()).
			Book(book).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)
	})

}