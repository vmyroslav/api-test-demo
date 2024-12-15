# \CoverPhotosAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1CoverPhotosBooksCoversIdBookGet**](CoverPhotosAPI.md#ApiV1CoverPhotosBooksCoversIdBookGet) | **Get** /api/v1/CoverPhotos/books/covers/{idBook} | 
[**ApiV1CoverPhotosGet**](CoverPhotosAPI.md#ApiV1CoverPhotosGet) | **Get** /api/v1/CoverPhotos | 
[**ApiV1CoverPhotosIdDelete**](CoverPhotosAPI.md#ApiV1CoverPhotosIdDelete) | **Delete** /api/v1/CoverPhotos/{id} | 
[**ApiV1CoverPhotosIdGet**](CoverPhotosAPI.md#ApiV1CoverPhotosIdGet) | **Get** /api/v1/CoverPhotos/{id} | 
[**ApiV1CoverPhotosIdPut**](CoverPhotosAPI.md#ApiV1CoverPhotosIdPut) | **Put** /api/v1/CoverPhotos/{id} | 
[**ApiV1CoverPhotosPost**](CoverPhotosAPI.md#ApiV1CoverPhotosPost) | **Post** /api/v1/CoverPhotos | 



## ApiV1CoverPhotosBooksCoversIdBookGet

> []CoverPhoto ApiV1CoverPhotosBooksCoversIdBookGet(ctx, idBook).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	idBook := int32(56) // int32 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosBooksCoversIdBookGet(context.Background(), idBook).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosBooksCoversIdBookGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1CoverPhotosBooksCoversIdBookGet`: []CoverPhoto
	fmt.Fprintf(os.Stdout, "Response from `CoverPhotosAPI.ApiV1CoverPhotosBooksCoversIdBookGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**idBook** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosBooksCoversIdBookGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]CoverPhoto**](CoverPhoto.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1CoverPhotosGet

> []CoverPhoto ApiV1CoverPhotosGet(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1CoverPhotosGet`: []CoverPhoto
	fmt.Fprintf(os.Stdout, "Response from `CoverPhotosAPI.ApiV1CoverPhotosGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosGetRequest struct via the builder pattern


### Return type

[**[]CoverPhoto**](CoverPhoto.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1CoverPhotosIdDelete

> ApiV1CoverPhotosIdDelete(ctx, id).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := int32(56) // int32 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1CoverPhotosIdGet

> CoverPhoto ApiV1CoverPhotosIdGet(ctx, id).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := int32(56) // int32 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1CoverPhotosIdGet`: CoverPhoto
	fmt.Fprintf(os.Stdout, "Response from `CoverPhotosAPI.ApiV1CoverPhotosIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CoverPhoto**](CoverPhoto.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1CoverPhotosIdPut

> CoverPhoto ApiV1CoverPhotosIdPut(ctx, id).CoverPhoto(coverPhoto).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := int32(56) // int32 | 
	coverPhoto := *openapiclient.NewCoverPhoto() // CoverPhoto |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosIdPut(context.Background(), id).CoverPhoto(coverPhoto).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1CoverPhotosIdPut`: CoverPhoto
	fmt.Fprintf(os.Stdout, "Response from `CoverPhotosAPI.ApiV1CoverPhotosIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **coverPhoto** | [**CoverPhoto**](CoverPhoto.md) |  | 

### Return type

[**CoverPhoto**](CoverPhoto.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json; v=1.0, text/json; v=1.0, application/*+json; v=1.0
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1CoverPhotosPost

> CoverPhoto ApiV1CoverPhotosPost(ctx).CoverPhoto(coverPhoto).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	coverPhoto := *openapiclient.NewCoverPhoto() // CoverPhoto |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CoverPhotosAPI.ApiV1CoverPhotosPost(context.Background()).CoverPhoto(coverPhoto).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CoverPhotosAPI.ApiV1CoverPhotosPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1CoverPhotosPost`: CoverPhoto
	fmt.Fprintf(os.Stdout, "Response from `CoverPhotosAPI.ApiV1CoverPhotosPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1CoverPhotosPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **coverPhoto** | [**CoverPhoto**](CoverPhoto.md) |  | 

### Return type

[**CoverPhoto**](CoverPhoto.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json; v=1.0, text/json; v=1.0, application/*+json; v=1.0
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

