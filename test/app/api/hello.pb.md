# gRPC API Reference

## Contents


- Services
    - [HelloGreeter](#servicehellohellogreeter)
    - [HelloIdentity](#servicehellohelloidentity)
  


- Messages
    - [HelloGreeterSayHelloRequest](#hellogreetersayhellorequest)
    - [HelloGreeterSayHelloResponse](#hellogreetersayhelloresponse)
    - [HelloIdentityVersionRequest](#helloidentityversionrequest)
    - [HelloIdentityVersionResponse](#helloidentityversionresponse)
    - [HelloVersion](#helloversion)
  



- [Scalar Value Types](#scalar-value-types)




## HelloGreeter {#servicehellohellogreeter}
The greeting service definition.

### SayHello {#methodhellohellogreetersayhello}

> **rpc** SayHello([HelloGreeterSayHelloRequest](#hellogreetersayhellorequest))
    [HelloGreeterSayHelloResponse](#hellogreetersayhelloresponse)

Sends a greeting
 <!-- end methods -->

## HelloIdentity {#servicehellohelloidentity}


### Version {#methodhellohelloidentityversion}

> **rpc** Version([HelloIdentityVersionRequest](#helloidentityversionrequest))
    [HelloIdentityVersionResponse](#helloidentityversionresponse)


 <!-- end methods -->
 <!-- end services -->

## Messages


### HelloGreeterSayHelloRequest {#hellogreetersayhellorequest}
The request message containing the user's name.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | none |
 <!-- end Fields -->
 <!-- end HasFields -->


### HelloGreeterSayHelloResponse {#hellogreetersayhelloresponse}
The response message containing the greetings


| Field | Type | Description |
| ----- | ---- | ----------- |
| message | [ string](#string) | none |
 <!-- end Fields -->
 <!-- end HasFields -->


### HelloIdentityVersionRequest {#helloidentityversionrequest}
Empty request

 <!-- end HasFields -->


### HelloIdentityVersionResponse {#helloidentityversionresponse}
Defines the response to version


| Field | Type | Description |
| ----- | ---- | ----------- |
| hello_version | [ HelloVersion](#helloversion) | Hello application version |
 <!-- end Fields -->
 <!-- end HasFields -->


### HelloVersion {#helloversion}
Hello version in Major.Minor.Patch format. The goal of this
message is to provide clients a method to determine the server
and client versions.


| Field | Type | Description |
| ----- | ---- | ----------- |
| major | [ int32](#int32) | Version major number |
| minor | [ int32](#int32) | Version minor number |
| patch | [ int32](#int32) | Version patch number |
| version | [ string](#string) | String representation of the version. Must be in `major.minor.patch` format. |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

## Enums


### HelloVersion.Version {#helloversionversion}
These values are constants that can be used by the
client and server applications

| Name | Number | Description |
| ---- | ------ | ----------- |
| MUST_HAVE_ZERO_VALUE | 0 | Must be set in the proto file; ignore. |
| MAJOR | 0 | Version major value of this specification |
| MINOR | 0 | Version minor value of this specification |
| PATCH | 1 | Version patch value of this specification |


 <!-- end Enums -->
 <!-- end Files -->

## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <div><h4 id="double" /></div><a name="double" /> double |  | double | double | float |
| <div><h4 id="float" /></div><a name="float" /> float |  | float | float | float |
| <div><h4 id="int32" /></div><a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <div><h4 id="int64" /></div><a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <div><h4 id="uint32" /></div><a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <div><h4 id="uint64" /></div><a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <div><h4 id="sint32" /></div><a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <div><h4 id="sint64" /></div><a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <div><h4 id="fixed32" /></div><a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <div><h4 id="fixed64" /></div><a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <div><h4 id="sfixed32" /></div><a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <div><h4 id="sfixed64" /></div><a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <div><h4 id="bool" /></div><a name="bool" /> bool |  | bool | boolean | boolean |
| <div><h4 id="string" /></div><a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <div><h4 id="bytes" /></div><a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

