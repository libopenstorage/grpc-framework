# Protocol Buffer Style Guide

Please first read Google's [Protocol Buffer Style Guide](https://developers.google.com/protocol-buffers/docs/style):

!!! quote

    This document provides a style guide for `.proto` files. By following these
    conventions, you'll make your protocol buffer message definitions and their
    corresponding classes consistent and easy to read.
    Note that protocol buffer style has evolved over time, so it is likely that you
    will see `.proto` files written in different conventions or styles. Please respect
    the existing style when you modify these files. Consistency is key. However, it
    is best to adopt the current best style when you are creating a new `.proto`file.

The following documentation is provided as a set of guidelines to help you in your gRPC APIs.

## Types

* `string` types should be used only for ids, messages, or opaque values. They are not meant to marshal information as a `yaml`. Instead create a concrete _message_.
* Only use `map<string, string>` for opaque values like labels, key-value pairs, etc. Do not use them for operations. Use enums instead.
* Value options should not be passed as `string`. Instead of passing "Done", or "paused", use enums for these value, making it clear to the reader.
* Try not to use `uint64`. Instead try to use signed `int64`. (See [CSI #172](https://github.com/container-storage-interface/spec/issues/172))

## Services

* See [CSI](https://github.com/container-storage-interface/spec/blob/master/csi.proto) as an example.
* Use Camelcase
* Services should be in the format `AppName<Service Name>`.
* Note that the service is a collection of APIs and are grouped as such in the documentation.
    * Here is an example for [OpenStorageVolume](https://libopenstorage.github.io/w/release-9.5.generated-api.html#serviceopenstorageapiopenstoragevolume)

## RPCs

* All APIs should have a single message for the request and a single message for the response with the following style: `[App]<Service Type><Api Name>Request|Response`
  * See [CSI](https://github.com/container-storage-interface/spec/blob/master/csi.proto) as an example.
* RPCs will be created as _methods_ to the service _object_, therefore there is
  no need to add the service name as part of the RPC. For example,
  use `Foo`, or `Bar` instead or `ServiceFoo` or `ServiceBar` as RPC names.

## Enums

* Follow the [Google protobuf style for enums](https://developers.google.com/protocol-buffers/docs/style#enums)
* According to the Google guide, the enum of zero value should be labeled as `UNSPECIFIED` to check if it was not set since `0` is the default value set when the client does not provide it.
* Wrap enums in messages so that their string values are clearer. Wrapping an enum in a message also has the benefit of not needing to prefix the enums with namespaced information. For example, instead of using the enum `XATTR_UNSPECIFIED`, the example above uses just `UNSPECIFIED` since it is inide the `Xattr` message. The generated code will be namepaced:

Proto:

```proto
// Xattr defines implementation specific volume attribute
message Xattr {
  enum Value {
    // Value is uninitialized or unknown
    UNSPECIFIED = 0;
    // Enable on-demand copy-on-write on the volume
    COW_ON_DEMAND = 1;
  }
}
```

Using the enum in a Proto

```proto
message VolumeSpec {
  // Holds the extended attributes for the volume
  Xattr.Value xattr = 1;
}
```

Notice the namepaced and string values in the generated output code:

```go
type Xattr_Value int32

const (
	// Value is uninitialized or unknown
	Xattr_UNSPECIFIED Xattr_Value = 0
	// Enable on-demand copy-on-write on the volume
	Xattr_COW_ON_DEMAND Xattr_Value = 1
)

var Xattr_Value_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "COW_ON_DEMAND",
}
var Xattr_Value_value = map[string]int32{
	"UNSPECIFIED":   0,
	"COW_ON_DEMAND": 1,
}

typedef VolueSpec struct {
  // Holds the extended attributes for the volume
	Xattr Xattr_Value `protobuf:"varint,36,opt,name=xattr,enum=openstorage.api.Xattr_Value" json:"xattr,omitempty"`
}
```

## Messages

* If at all possible, APIs _must_ be supported forever once released. 
* They will almost never be deprecated since at some point you may have many versions of the clients. So please be clear and careful on the API you create.
* If we need to change or update, you can always **add** values.

### Field Numbers

* If it is a new message, start with the field number of `1`.
* If it is an addition to a message, continue the field number sequence by one.
* If you are using `oneof` you may want to start with a large value for the
  field number so that they do not interfere with other values in the message:

```proto
  string s3_storage_class = 7;

  // Start at field number 200 to allow for expansion
  oneof credential_type {
    // Credentials for AWS/S3
    SdkAwsCredentialRequest aws_credential = 200;
    // Credentials for Azure
    SdkAzureCredentialRequest azure_credential = 201;
    // Credentials for Google
    SdkGoogleCredentialRequest google_credential = 202;
  }
```

### Deprecation

Here is the process if you would like to deprecate:

1. According to [proto3 Language Guide](https://developers.google.com/protocol-buffers/docs/proto3) set the value in the message to deprecated and add a `(deprecated)` string to the comment as follows:

```proto
// (deprecated) Field documentation here
int32 field = 6 [deprecated = true];
```

2. Comment in the a changelog that the value is deprecated.
3. Provide at least two releases before removing support for that value in the message. Make sure to document in the release notes of the product the deprecation.
4. Once at least two releases have passed. Reserve the field number as shown in the [proto3 Language Guide](https://developers.google.com/protocol-buffers/docs/proto3#reserved):

```proto
message Foo {
  reserved 6;
}
```

It is essential that no values override the field number when updating or replacing. From Google's guide:

!!! warning 

    If you update a message type by entirely removing a field, or commenting it
    out, future users can reuse the field number when making their own updates
    to the type. This can cause severe issues if they later load old versions of
    the same .proto, including data corruption, privacy bugs, and so on.

## REST

REST endpoints are autogenerated from the protofile by the
[grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/) protoc compiler.
All APIs should add the appropriate information to generate a
REST endpoint for their service. Here is an example:

```proto
  rpc Inspect(RoleInspectRequest)
    returns (RoleInspectResponse){
      option(google.api.http) = {
        get: "/v1/roles/{name}"
      };
    }

  // Delete an existing role
  rpc Delete(RoleDeleteRequest)
    returns (RoleDeleteResponse){
      option(google.api.http) = {
        delete: "/v1/roles/{name}"
      };
    }

  // Update an existing role
  rpc Update(RoleUpdateRequest)
    returns (RoleUpdateResponse){
      option(google.api.http) = {
        put: "/v1/roles"
        body: "*"
      };
    }
```

Here are the guidelines for REST commands:

* Endpoint must be prefixed as follows: `/v1/<service name>/<rpc name if needed>/{any variables if needed}`.
* Use the appropriate HTTP method. Here are some guidelines:
    * For _Create_ RPCs use the `post` http method
    * For _Inspect_ and _List_ RPCs use the `get` http method
    * For _Update_ RPCs use the `put` http method
    * For _Delete_ RPCs use the `delete` http method
* Use `get` for immutable calls.
* Use `put` with `body: "*"` most calls that need to send a message to the SDK server.

Please see [grpc-gateway documentation](https://grpc-ecosystem.github.io/grpc-gateway/) for more information.

## Documentation

* All APIs, messages, and types should be documented if possible. The `grpc-framework` utilizes
[protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) to automatically generate documentation from
your protocol buffers file.

* **Documenting Messages**
    * Document each value of the message.
    * Do not use Golang style. Do not repeat the name of the variable in _Golang Camel Format_ in the comment to document it since the variable could be in other styles in other languages. For example:

```proto
// Provides volume's exclusive bytes and its total usage. This cannot be
// retrieved individually and is obtained as part node's usage for a given
// node.
message VolumeUsage {
   // id for the volume/snapshot
  string volume_id = 1;
  // name of the volume/snapshot
  string volume_name = 2;
  // uuid of the pool that this volume belongs to
  string pool_uuid = 3;
  // size in bytes exclusively used by the volume/snapshot
  uint64 exclusive_bytes = 4;
  //  size in bytes by the volume/snapshot
  uint64 total_bytes = 5;
  // set to true if this volume is snapshot created by cloudbackups
  bool local_cloud_snapshot = 6;
}
```

### Here is an example: 

* Protocol buffers file: [hello.proto](https://github.com/libopenstorage/grpc-framework/blob/master/test/app/api/hello.proto)
* Documentation in markdown format: [hello.pb.md](hello.pb.md)

