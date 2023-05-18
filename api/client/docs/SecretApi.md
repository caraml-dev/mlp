# \SecretApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1ProjectsProjectIdSecretsGet**](SecretApi.md#V1ProjectsProjectIdSecretsGet) | **Get** /v1/projects/{project_id}/secrets | List secret
[**V1ProjectsProjectIdSecretsPost**](SecretApi.md#V1ProjectsProjectIdSecretsPost) | **Post** /v1/projects/{project_id}/secrets | Create secret
[**V1ProjectsProjectIdSecretsSecretIdDelete**](SecretApi.md#V1ProjectsProjectIdSecretsSecretIdDelete) | **Delete** /v1/projects/{project_id}/secrets/{secret_id} | Delete secret
[**V1ProjectsProjectIdSecretsSecretIdGet**](SecretApi.md#V1ProjectsProjectIdSecretsSecretIdGet) | **Get** /v1/projects/{project_id}/secrets/{secret_id} | Get secret
[**V1ProjectsProjectIdSecretsSecretIdPatch**](SecretApi.md#V1ProjectsProjectIdSecretsSecretIdPatch) | **Patch** /v1/projects/{project_id}/secrets/{secret_id} | Update secret


# **V1ProjectsProjectIdSecretsGet**
> []Secret V1ProjectsProjectIdSecretsGet(ctx, projectId)
List secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 

### Return type

[**[]Secret**](Secret.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretsPost**
> Secret V1ProjectsProjectIdSecretsPost(ctx, projectId, body)
Create secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 
  **body** | [**Secret**](Secret.md)|  | 

### Return type

[**Secret**](Secret.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretsSecretIdDelete**
> V1ProjectsProjectIdSecretsSecretIdDelete(ctx, projectId, secretId)
Delete secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 
  **secretId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretsSecretIdGet**
> Secret V1ProjectsProjectIdSecretsSecretIdGet(ctx, projectId, secretId)
Get secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 
  **secretId** | **int32**|  | 

### Return type

[**Secret**](Secret.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdSecretsSecretIdPatch**
> Secret V1ProjectsProjectIdSecretsSecretIdPatch(ctx, projectId, secretId, optional)
Update secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 
  **secretId** | **int32**|  | 
 **optional** | ***SecretApiV1ProjectsProjectIdSecretsSecretIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SecretApiV1ProjectsProjectIdSecretsSecretIdPatchOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**optional.Interface of Secret**](Secret.md)|  | 

### Return type

[**Secret**](Secret.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

