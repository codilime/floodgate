# \V2CanaryControllerApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCanaryResultUsingGET**](V2CanaryControllerApi.md#GetCanaryResultUsingGET) | **Get** /v2/canaries/canary/{canaryConfigId}/{canaryExecutionId} | (DEPRECATED) Retrieve a canary result
[**GetCanaryResultUsingGET1**](V2CanaryControllerApi.md#GetCanaryResultUsingGET1) | **Get** /v2/canaries/canary/{canaryExecutionId} | Retrieve a canary result
[**GetCanaryResultsByApplicationUsingGET**](V2CanaryControllerApi.md#GetCanaryResultsByApplicationUsingGET) | **Get** /v2/canaries/{application}/executions | Retrieve a list of an application&#39;s canary results
[**GetMetricSetPairListUsingGET**](V2CanaryControllerApi.md#GetMetricSetPairListUsingGET) | **Get** /v2/canaries/metricSetPairList/{metricSetPairListId} | Retrieve a metric set pair list
[**InitiateCanaryUsingPOST**](V2CanaryControllerApi.md#InitiateCanaryUsingPOST) | **Post** /v2/canaries/canary/{canaryConfigId} | Start a canary execution
[**InitiateCanaryWithConfigUsingPOST**](V2CanaryControllerApi.md#InitiateCanaryWithConfigUsingPOST) | **Post** /v2/canaries/canary | Start a canary execution with the supplied canary config
[**ListCredentialsUsingGET**](V2CanaryControllerApi.md#ListCredentialsUsingGET) | **Get** /v2/canaries/credentials | Retrieve a list of configured Kayenta accounts
[**ListJudgesUsingGET**](V2CanaryControllerApi.md#ListJudgesUsingGET) | **Get** /v2/canaries/judges | Retrieve a list of all configured canary judges
[**ListMetricsServiceMetadataUsingGET**](V2CanaryControllerApi.md#ListMetricsServiceMetadataUsingGET) | **Get** /v2/canaries/metadata/metricsService | Retrieve a list of descriptors for use in populating the canary config ui


# **GetCanaryResultUsingGET**
> interface{} GetCanaryResultUsingGET(ctx, canaryConfigId, canaryExecutionId, optional)
(DEPRECATED) Retrieve a canary result

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **canaryConfigId** | **string**| canaryConfigId | 
  **canaryExecutionId** | **string**| canaryExecutionId | 
 **optional** | ***GetCanaryResultUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetCanaryResultUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCanaryResultUsingGET1**
> interface{} GetCanaryResultUsingGET1(ctx, canaryExecutionId, optional)
Retrieve a canary result

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **canaryExecutionId** | **string**| canaryExecutionId | 
 **optional** | ***GetCanaryResultUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetCanaryResultUsingGET1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCanaryResultsByApplicationUsingGET**
> []interface{} GetCanaryResultsByApplicationUsingGET(ctx, application, limit, optional)
Retrieve a list of an application's canary results

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **application** | **string**| application | 
  **limit** | **int32**| limit | 
 **optional** | ***GetCanaryResultsByApplicationUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetCanaryResultsByApplicationUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **statuses** | **optional.String**| Comma-separated list of statuses, e.g.: RUNNING, SUCCEEDED, TERMINAL | 
 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMetricSetPairListUsingGET**
> []interface{} GetMetricSetPairListUsingGET(ctx, metricSetPairListId, optional)
Retrieve a metric set pair list

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **metricSetPairListId** | **string**| metricSetPairListId | 
 **optional** | ***GetMetricSetPairListUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetMetricSetPairListUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **InitiateCanaryUsingPOST**
> interface{} InitiateCanaryUsingPOST(ctx, canaryConfigId, executionRequest, optional)
Start a canary execution

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **canaryConfigId** | **string**| canaryConfigId | 
  **executionRequest** | [**interface{}**](interface{}.md)| executionRequest | 
 **optional** | ***InitiateCanaryUsingPOSTOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a InitiateCanaryUsingPOSTOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **application** | **optional.String**| application | 
 **configurationAccountName** | **optional.String**| configurationAccountName | 
 **metricsAccountName** | **optional.String**| metricsAccountName | 
 **parentPipelineExecutionId** | **optional.String**| parentPipelineExecutionId | 
 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **InitiateCanaryWithConfigUsingPOST**
> interface{} InitiateCanaryWithConfigUsingPOST(ctx, adhocExecutionRequest, optional)
Start a canary execution with the supplied canary config

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **adhocExecutionRequest** | [**interface{}**](interface{}.md)| adhocExecutionRequest | 
 **optional** | ***InitiateCanaryWithConfigUsingPOSTOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a InitiateCanaryWithConfigUsingPOSTOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **application** | **optional.String**| application | 
 **metricsAccountName** | **optional.String**| metricsAccountName | 
 **parentPipelineExecutionId** | **optional.String**| parentPipelineExecutionId | 
 **storageAccountName** | **optional.String**| storageAccountName | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListCredentialsUsingGET**
> []interface{} ListCredentialsUsingGET(ctx, )
Retrieve a list of configured Kayenta accounts

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListJudgesUsingGET**
> []interface{} ListJudgesUsingGET(ctx, )
Retrieve a list of all configured canary judges

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListMetricsServiceMetadataUsingGET**
> []interface{} ListMetricsServiceMetadataUsingGET(ctx, optional)
Retrieve a list of descriptors for use in populating the canary config ui

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ListMetricsServiceMetadataUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ListMetricsServiceMetadataUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| filter | 
 **metricsAccountName** | **optional.String**| metricsAccountName | 

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

