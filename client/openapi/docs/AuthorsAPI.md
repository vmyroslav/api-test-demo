# \AuthorsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1AuthorsAuthorsBooksIdBookGet**](AuthorsAPI.md#ApiV1AuthorsAuthorsBooksIdBookGet) | **Get** /api/v1/Authors/authors/books/{idBook} | 
[**ApiV1AuthorsGet**](AuthorsAPI.md#ApiV1AuthorsGet) | **Get** /api/v1/Authors | 
[**ApiV1AuthorsIdDelete**](AuthorsAPI.md#ApiV1AuthorsIdDelete) | **Delete** /api/v1/Authors/{id} | 
[**ApiV1AuthorsIdGet**](AuthorsAPI.md#ApiV1AuthorsIdGet) | **Get** /api/v1/Authors/{id} | 
[**ApiV1AuthorsIdPut**](AuthorsAPI.md#ApiV1AuthorsIdPut) | **Put** /api/v1/Authors/{id} | 
[**ApiV1AuthorsPost**](AuthorsAPI.md#ApiV1AuthorsPost) | **Post** /api/v1/Authors | 



## ApiV1AuthorsAuthorsBooksIdBookGet

> []Author ApiV1AuthorsAuthorsBooksIdBookGet(ctx, idBook).Execute()



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
	resp, r, err := apiClient.AuthorsAPI.ApiV1AuthorsAuthorsBooksIdBookGet(context.Background(), idBook).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsAuthorsBooksIdBookGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1AuthorsAuthorsBooksIdBookGet`: []Author
	fmt.Fprintf(os.Stdout, "Response from `AuthorsAPI.ApiV1AuthorsAuthorsBooksIdBookGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**idBook** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1AuthorsAuthorsBooksIdBookGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Author**](Author.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1AuthorsGet

> []Author ApiV1AuthorsGet(ctx).Execute()



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
	resp, r, err := apiClient.AuthorsAPI.ApiV1AuthorsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1AuthorsGet`: []Author
	fmt.Fprintf(os.Stdout, "Response from `AuthorsAPI.ApiV1AuthorsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1AuthorsGetRequest struct via the builder pattern


### Return type

[**[]Author**](Author.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1AuthorsIdDelete

> ApiV1AuthorsIdDelete(ctx, id).Execute()



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
	r, err := apiClient.AuthorsAPI.ApiV1AuthorsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsIdDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiApiV1AuthorsIdDeleteRequest struct via the builder pattern


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


## ApiV1AuthorsIdGet

> Author ApiV1AuthorsIdGet(ctx, id).Execute()



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
	resp, r, err := apiClient.AuthorsAPI.ApiV1AuthorsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1AuthorsIdGet`: Author
	fmt.Fprintf(os.Stdout, "Response from `AuthorsAPI.ApiV1AuthorsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1AuthorsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Author**](Author.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1AuthorsIdPut

> Author ApiV1AuthorsIdPut(ctx, id).Author(author).Execute()



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
	author := *openapiclient.NewAuthor() // Author |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthorsAPI.ApiV1AuthorsIdPut(context.Background(), id).Author(author).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1AuthorsIdPut`: Author
	fmt.Fprintf(os.Stdout, "Response from `AuthorsAPI.ApiV1AuthorsIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1AuthorsIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **author** | [**Author**](Author.md) |  | 

### Return type

[**Author**](Author.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json; v=1.0, text/json; v=1.0, application/*+json; v=1.0
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1AuthorsPost

> Author ApiV1AuthorsPost(ctx).Author(author).Execute()



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
	author := *openapiclient.NewAuthor() // Author |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthorsAPI.ApiV1AuthorsPost(context.Background()).Author(author).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthorsAPI.ApiV1AuthorsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1AuthorsPost`: Author
	fmt.Fprintf(os.Stdout, "Response from `AuthorsAPI.ApiV1AuthorsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1AuthorsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **author** | [**Author**](Author.md) |  | 

### Return type

[**Author**](Author.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json; v=1.0, text/json; v=1.0, application/*+json; v=1.0
- **Accept**: text/plain; v=1.0, application/json; v=1.0, text/json; v=1.0

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

