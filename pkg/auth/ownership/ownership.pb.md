# gRPC API Reference

## Contents



- Messages
    - [Ownership](#ownership)
    - [Ownership.AccessControl](#ownershipaccesscontrol)
    - [Ownership.AccessControl.CollaboratorsEntry](#ownershipaccesscontrolcollaboratorsentry)
    - [Ownership.AccessControl.GroupsEntry](#ownershipaccesscontrolgroupsentry)
    - [Ownership.PublicAccessControl](#ownershippublicaccesscontrol)
  



- [Scalar Value Types](#scalar-value-types)



 <!-- end services -->

## Messages


### Ownership {#ownership}
Ownership information for resource.
Administrators are users who belong to the group `*`, meaning, every group.


| Field | Type | Description |
| ----- | ---- | ----------- |
| owner | [ string](#string) | Username of owner.

The storage system uses the username taken from the security authorization token and is saved on this field. Only users with system administration can edit this value. |
| acls | [ Ownership.AccessControl](#ownershipaccesscontrol) | Permissions to share resource which can be set by the owner.

NOTE: To create an "admin" user which has access to any resource set the group value in the token of the user to `*`. |
 <!-- end Fields -->
 <!-- end HasFields -->


### Ownership.AccessControl {#ownershipaccesscontrol}



| Field | Type | Description |
| ----- | ---- | ----------- |
| groups | [map Ownership.AccessControl.GroupsEntry](#ownershipaccesscontrolgroupsentry) | Group access to resource which must match the group set in the authorization token. Can be set by the owner or the system administrator only. Possible values are: 1. no groups: Means no groups are given access. 2. `["*"]`: All groups are allowed. 3. `["group1", "group2"]`: Only certain groups are allowed. In this example only _group1_ and _group2_ are allowed. |
| collaborators | [map Ownership.AccessControl.CollaboratorsEntry](#ownershipaccesscontrolcollaboratorsentry) | Collaborator access to resource gives access to other user. Must be the username (unique id) set in the authorization token. The owner or the administrator can set this value. Possible values are: 1. no collaborators: Means no users are given access. 2. `["*"]`: All users are allowed. 3. `["username1", "username2"]`: Only certain usernames are allowed. In this example only _username1_ and _username2_ are allowed. |
| public | [ Ownership.PublicAccessControl](#ownershippublicaccesscontrol) | Public access to resource may be assigned for access by the public userd |
 <!-- end Fields -->
 <!-- end HasFields -->


### Ownership.AccessControl.CollaboratorsEntry {#ownershipaccesscontrolcollaboratorsentry}



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | none |
| value | [ Ownership.AccessType](#ownershipaccesstype) | none |
 <!-- end Fields -->
 <!-- end HasFields -->


### Ownership.AccessControl.GroupsEntry {#ownershipaccesscontrolgroupsentry}



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | none |
| value | [ Ownership.AccessType](#ownershipaccesstype) | none |
 <!-- end Fields -->
 <!-- end HasFields -->


### Ownership.PublicAccessControl {#ownershippublicaccesscontrol}
PublicAccessControl allows assigning public ownership


| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ Ownership.AccessType](#ownershipaccesstype) | AccessType declares which level of public access is allowed |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

## Enums


### Ownership.AccessType {#ownershipaccesstype}
Access types can be set by owner to have different levels of access to
a resource.

It is up to the resource to interpret what the types mean and are
used for.

| Name | Number | Description |
| ---- | ------ | ----------- |
| READ | 0 | Read access only and cannot affect the resource. |
| WRITE | 1 | Write access and can affect the resource. This type automatically provides Read access also. |
| ADMIN | 2 | Administrator access. This type automatically provides Read and Write access also. |


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

