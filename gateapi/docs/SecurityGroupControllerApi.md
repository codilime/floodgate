# \SecurityGroupControllerApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AllByAccountUsingGET1**](SecurityGroupControllerApi.md#AllByAccountUsingGET1) | **Get** /securityGroups/{account} | Retrieve a list of security groups for a given account, grouped by region
[**AllUsingGET5**](SecurityGroupControllerApi.md#AllUsingGET5) | **Get** /securityGroups | Retrieve a list of security groups, grouped by account, cloud provider, and region
[**GetSecurityGroupUsingGET1**](SecurityGroupControllerApi.md#GetSecurityGroupUsingGET1) | **Get** /securityGroups/{account}/{region}/{name} | Retrieve a security group&#39;s details


# **AllByAccountUsingGET1**
> interface{} AllByAccountUsingGET1(ctx, account, optional)
Retrieve a list of security groups for a given account, grouped by region

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | **string**| account | 
 **optional** | ***AllByAccountUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AllByAccountUsingGET1Opts struct

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

# **AllUsingGET5**
> interface{} AllUsingGET5(ctx, optional)
Retrieve a list of security groups, grouped by account, cloud provider, and region

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AllUsingGET5Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AllUsingGET5Opts struct

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

# **GetSecurityGroupUsingGET1**
> interface{} GetSecurityGroupUsingGET1(ctx, account, name, region, optional)
Retrieve a security group's details

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | **string**| account | 
  **name** | **string**| name | 
  **region** | **string**| region | 
 **optional** | ***GetSecurityGroupUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetSecurityGroupUsingGET1Opts struct

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

