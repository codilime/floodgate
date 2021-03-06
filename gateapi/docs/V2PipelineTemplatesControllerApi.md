# \V2PipelineTemplatesControllerApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateUsingPOST1**](V2PipelineTemplatesControllerApi.md#CreateUsingPOST1) | **Post** /v2/pipelineTemplates/create | (ALPHA) Create a pipeline template.
[**DeleteUsingDELETE1**](V2PipelineTemplatesControllerApi.md#DeleteUsingDELETE1) | **Delete** /v2/pipelineTemplates/{id} | Delete a pipeline template.
[**GetUsingGET2**](V2PipelineTemplatesControllerApi.md#GetUsingGET2) | **Get** /v2/pipelineTemplates/{id} | (ALPHA) Get a pipeline template.
[**ListPipelineTemplateDependentsUsingGET1**](V2PipelineTemplatesControllerApi.md#ListPipelineTemplateDependentsUsingGET1) | **Get** /v2/pipelineTemplates/{id}/dependents | (ALPHA) List all pipelines that implement a pipeline template
[**ListUsingGET1**](V2PipelineTemplatesControllerApi.md#ListUsingGET1) | **Get** /v2/pipelineTemplates | (ALPHA) List pipeline templates.
[**ListVersionsUsingGET**](V2PipelineTemplatesControllerApi.md#ListVersionsUsingGET) | **Get** /v2/pipelineTemplates/versions | List pipeline templates with versions
[**PlanUsingPOST**](V2PipelineTemplatesControllerApi.md#PlanUsingPOST) | **Post** /v2/pipelineTemplates/plan | (ALPHA) Plan a pipeline template configuration.
[**UpdateUsingPOST1**](V2PipelineTemplatesControllerApi.md#UpdateUsingPOST1) | **Post** /v2/pipelineTemplates/update/{id} | (ALPHA) Update a pipeline template.


# **CreateUsingPOST1**
> map[string]interface{} CreateUsingPOST1(ctx, pipelineTemplate, optional)
(ALPHA) Create a pipeline template.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **pipelineTemplate** | [**interface{}**](interface{}.md)| pipelineTemplate | 
 **optional** | ***CreateUsingPOST1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CreateUsingPOST1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **tag** | **optional.String**| tag | 

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteUsingDELETE1**
> map[string]interface{} DeleteUsingDELETE1(ctx, id, optional)
Delete a pipeline template.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| id | 
 **optional** | ***DeleteUsingDELETE1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeleteUsingDELETE1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **application** | **optional.String**| application | 
 **digest** | **optional.String**| digest | 
 **tag** | **optional.String**| tag | 

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUsingGET2**
> map[string]interface{} GetUsingGET2(ctx, id, optional)
(ALPHA) Get a pipeline template.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| id | 
 **optional** | ***GetUsingGET2Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetUsingGET2Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **digest** | **optional.String**| digest | 
 **tag** | **optional.String**| tag | 

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListPipelineTemplateDependentsUsingGET1**
> []interface{} ListPipelineTemplateDependentsUsingGET1(ctx, id)
(ALPHA) List all pipelines that implement a pipeline template

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| id | 

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUsingGET1**
> []interface{} ListUsingGET1(ctx, optional)
(ALPHA) List pipeline templates.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ListUsingGET1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scopes** | [**optional.Interface of []string**](string.md)| scopes | 

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListVersionsUsingGET**
> interface{} ListVersionsUsingGET(ctx, optional)
List pipeline templates with versions

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListVersionsUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ListVersionsUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scopes** | [**optional.Interface of []string**](string.md)| scopes | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PlanUsingPOST**
> map[string]interface{} PlanUsingPOST(ctx, pipeline)
(ALPHA) Plan a pipeline template configuration.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **pipeline** | [**interface{}**](interface{}.md)| pipeline | 

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateUsingPOST1**
> map[string]interface{} UpdateUsingPOST1(ctx, id, pipelineTemplate, optional)
(ALPHA) Update a pipeline template.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| id | 
  **pipelineTemplate** | [**interface{}**](interface{}.md)| pipelineTemplate | 
 **optional** | ***UpdateUsingPOST1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UpdateUsingPOST1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **skipPlanDependents** | **optional.Bool**| skipPlanDependents | [default to false]
 **tag** | **optional.String**| tag | 

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

