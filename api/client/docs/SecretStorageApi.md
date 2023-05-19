# \SecretStorageApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1ProjectsProjectIdSecretStoragesGet**](SecretStorageApi.md#V1ProjectsProjectIdSecretStoragesGet) | **Get** /v1/projects/{project_id}/secret_storages | List secret storage
[**V1ProjectsProjectIdSecretStoragesPost**](SecretStorageApi.md#V1ProjectsProjectIdSecretStoragesPost) | **Post** /v1/projects/{project_id}/secret_storages | Create secret storage
[**V1ProjectsProjectIdSecretStoragesSecretStorageIdDelete**](SecretStorageApi.md#V1ProjectsProjectIdSecretStoragesSecretStorageIdDelete) | **Delete** /v1/projects/{project_id}/secret_storages/{secret_storage_id} | Delete secret storage
[**V1ProjectsProjectIdSecretStoragesSecretStorageIdGet**](SecretStorageApi.md#V1ProjectsProjectIdSecretStoragesSecretStorageIdGet) | **Get** /v1/projects/{project_id}/secret_storages/{secret_storage_id} | Get secret storage
[**V1ProjectsProjectIdSecretStoragesSecretStorageIdPatch**](SecretStorageApi.md#V1ProjectsProjectIdSecretStoragesSecretStorageIdPatch) | **Patch** /v1/projects/{project_id}/secret_storages/{secret_storage_id} | Update secret storage


# **V1ProjectsProjectIdSecretStoragesGet**
> []SecretStorage V1ProjectsProjectIdSecretStoragesGet(ctx, projectID)
List secret storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectID** | **int32**|  | 

### Return type

[**[]SecretStorage**](SecretStorage.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretStoragesPost**
> SecretStorage V1ProjectsProjectIdSecretStoragesPost(ctx, projectID, body)
Create secret storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectID** | **int32**|  | 
  **body** | [**SecretStorage**](SecretStorage.md)|  | 

### Return type

[**SecretStorage**](SecretStorage.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretStoragesSecretStorageIdDelete**
> V1ProjectsProjectIdSecretStoragesSecretStorageIdDelete(ctx, projectID, secretStorageID)
Delete secret storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectID** | **int32**|  | 
  **secretStorageID** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretStoragesSecretStorageIdGet**
> SecretStorage V1ProjectsProjectIdSecretStoragesSecretStorageIdGet(ctx, projectID, secretStorageID)
Get secret storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectID** | **int32**|  | 
  **secretStorageID** | **int32**|  | 

### Return type

[**SecretStorage**](SecretStorage.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretStoragesSecretStorageIdPatch**
> SecretStorage V1ProjectsProjectIdSecretStoragesSecretStorageIdPatch(ctx, projectID, secretStorageID, optional)
Update secret storage

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectID** | **int32**|  | 
  **secretStorageID** | **int32**|  | 
 **optional** | ***SecretStorageApiV1ProjectsProjectIdSecretStoragesSecretStorageIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SecretStorageApiV1ProjectsProjectIdSecretStoragesSecretStorageIdPatchOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**optional.Interface of SecretStorage**](SecretStorage.md)|  | 

### Return type

[**SecretStorage**](SecretStorage.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

