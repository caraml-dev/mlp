# \SecretApi

All URIs are relative to *http://localhost:8080/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectsProjectIdSecretsGet**](SecretApi.md#ProjectsProjectIdSecretsGet) | **Get** /projects/{project_id}/secrets | List secret
[**ProjectsProjectIdSecretsPost**](SecretApi.md#ProjectsProjectIdSecretsPost) | **Post** /projects/{project_id}/secrets | Create secret
[**ProjectsProjectIdSecretsSecretIdDelete**](SecretApi.md#ProjectsProjectIdSecretsSecretIdDelete) | **Delete** /projects/{project_id}/secrets/{secret_id} | Delete secret
[**ProjectsProjectIdSecretsSecretIdPatch**](SecretApi.md#ProjectsProjectIdSecretsSecretIdPatch) | **Patch** /projects/{project_id}/secrets/{secret_id} | Update secret


# **ProjectsProjectIdSecretsGet**
> []Secret ProjectsProjectIdSecretsGet(ctx, projectId)
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

# **ProjectsProjectIdSecretsPost**
> Secret ProjectsProjectIdSecretsPost(ctx, projectId, body)
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

# **ProjectsProjectIdSecretsSecretIdDelete**
> ProjectsProjectIdSecretsSecretIdDelete(ctx, projectId, secretId)
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

# **ProjectsProjectIdSecretsSecretIdPatch**
> Secret ProjectsProjectIdSecretsSecretIdPatch(ctx, projectId, secretId, optional)
Update secret

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**|  | 
  **secretId** | **int32**|  | 
 **optional** | ***SecretApiProjectsProjectIdSecretsSecretIdPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SecretApiProjectsProjectIdSecretsSecretIdPatchOpts struct

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

