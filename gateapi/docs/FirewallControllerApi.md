# \FirewallControllerApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AllByAccountAndRegionUsingGET**](FirewallControllerApi.md#AllByAccountAndRegionUsingGET) | **Get** /firewalls/{account}/{region} | Retrieve a list of firewalls for a given account and region
[**AllByAccountUsingGET**](FirewallControllerApi.md#AllByAccountUsingGET) | **Get** /firewalls/{account} | Retrieve a list of firewalls for a given account, grouped by region
[**AllUsingGET1**](FirewallControllerApi.md#AllUsingGET1) | **Get** /firewalls | Retrieve a list of firewalls, grouped by account, cloud provider, and region
[**GetSecurityGroupUsingGET**](FirewallControllerApi.md#GetSecurityGroupUsingGET) | **Get** /firewalls/{account}/{region}/{name} | Retrieve a firewall&#39;s details


# **AllByAccountAndRegionUsingGET**
> []interface{} AllByAccountAndRegionUsingGET(ctx, account, region, optional)
Retrieve a list of firewalls for a given account and region

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | **string**| account | 
  **region** | **string**| region | 
 **optional** | ***AllByAccountAndRegionUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AllByAccountAndRegionUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xRateLimitApp** | **optional.String**| X-RateLimit-App | 
 **provider** | **optional.String**| provider | [default to aws]

### Return type

[**[]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AllByAccountUsingGET**
> interface{} AllByAccountUsingGET(ctx, account, optional)
Retrieve a list of firewalls for a given account, grouped by region

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | **string**| account | 
 **optional** | ***AllByAccountUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AllByAccountUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRateLimitApp** | **optional.String**| X-RateLimit-App | 
 **provider** | **optional.String**| provider | [default to aws]

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AllUsingGET1**
> interface{} AllUsingGET1(ctx, optional)
Retrieve a list of firewalls, grouped by account, cloud provider, and region

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AllUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AllUsingGET1Opts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRateLimitApp** | **optional.String**| X-RateLimit-App | 
 **id** | **optional.String**| id | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSecurityGroupUsingGET**
> interface{} GetSecurityGroupUsingGET(ctx, account, name, region, optional)
Retrieve a firewall's details

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | **string**| account | 
  **name** | **string**| name | 
  **region** | **string**| region | 
 **optional** | ***GetSecurityGroupUsingGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetSecurityGroupUsingGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRateLimitApp** | **optional.String**| X-RateLimit-App | 
 **provider** | **optional.String**| provider | [default to aws]
 **vpcId** | **optional.String**| vpcId | 

### Return type

[**interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

