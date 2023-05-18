# \ProjectApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1ProjectsGet**](ProjectApi.md#V1ProjectsGet) | **Get** /v1/projects | List existing projects
[**V1ProjectsPost**](ProjectApi.md#V1ProjectsPost) | **Post** /v1/projects | Create new project
[**V1ProjectsProjectIdGet**](ProjectApi.md#V1ProjectsProjectIdGet) | **Get** /v1/projects/{project_id} | Get project
[**V1ProjectsProjectIdPut**](ProjectApi.md#V1ProjectsProjectIdPut) | **Put** /v1/projects/{project_id} | Update project


# **V1ProjectsGet**
> []Project V1ProjectsGet(ctx, optional)
List existing projects

Projects can be filtered by optional `name` parameter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ProjectApiV1ProjectsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProjectApiV1ProjectsGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **optional.String**|  | 

### Return type

[**[]Project**](Project.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsPost**
> Project V1ProjectsPost(ctx, body)
Create new project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Project**](Project.md)| Project object that has to be added | 

### Return type

[**Project**](Project.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdGet**
> Project V1ProjectsProjectIdGet(ctx, projectId)
Get project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**| project id of the project to be retrieved | 

### Return type

[**Project**](Project.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **V1ProjectsProjectIdPut**
> Project V1ProjectsProjectIdPut(ctx, projectId, body)
Update project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectId** | **int32**| project id of the project to be updated | 
  **body** | [**Project**](Project.md)| Project object that has to be updated | 

### Return type

[**Project**](Project.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

