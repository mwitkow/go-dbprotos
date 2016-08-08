# Protobuf Schema for NoSQL databases

Define your NoSQL store (e.g. Datastore) schema in a Protobuf. 

*Status*: Work in Progress.

## Why?

Typically structs for data representations of key value storages (Google Datastore, MongoDB) rely on the data model 
being expressed in the code. In Go, these are usually `struct` field tags. 
 
Protocol buffers provide a way of expressing datastructures that is centred around the data itself, not the 
implementation. The ideas is: if you're changing a database proto file, it may have significant consequences. That's why
this project proposes that all database enties, including their indexes are defined in `.proto` files.

## Example

```proto
syntax = "proto3";

message SomeProto {
  option (dbp.entity).datastore.kind = "MySomeProto";

  string UnusedField = 1;
  string SomeSingleString = 2 [
    (dbp.datastore).name = "single_string",
    (dbp.index) = { single: True, composite: [ {id: "first"}, {id: "second" } ]
  ];
  repeated string SomeMultiString = 4 [
    (dbp.datastore) = {name: "multi_string"},
    (dbp.index) = { single: True }
  ];
  google.protobuf.Timestamp SomeSingleTime = 5 [
    (dbp.datastore) = {name: "single_time"},
    (dbp.index) = { single: True, composite: [ {id: "first"} ] }
  ];
  repeated int32 SomeMultiInt = 6 [
    (dbp.datastore) = {name: "multi_int"},
  ];
  repeated google.protobuf.Timestamp SomeMultiTimes = 7 [
    (dbp.datastore) = {name: "multi_time"},
  ];
  string SomeInteger = 8 [ 
    (dbp.datastore).name = "my_integer",
    (dbp.index) = { single: True, composite: [ {id: "second", descending: True } ]
  ];
  int DeprecatedSizeParameter = 9 [ 
    (dbp.datastore).name = "size",
    (dbp.datastore).not_writeable = true,
  ];
}
```

In this example we see that the given proto Message is tied to a specific entity type. Each field of the message has an
annotation that provides a canonical name for the field. Moreover, the field annotations control indexing: whether the
given field forms an indiviually indexed column (`single`) or is part of a named composite object.

The typical Protobuf concepts apply: if the entity contains fields not present in the proto, they will be ignored 
(unless a message-level option `strict_reading` is used), and all non-default value (by `proto3` semantics) fields will
 be written to the output.

## Roadmap

 [*] - working code generator for simple Datastore types and repeated fields
 [*] - tests for compatibility between proto representation and Datastore canonical librayr
 [ ] - support for Enums and nested messages
 [ ] - grooming tools that use `.proto` files to assess compatibility of database tools
 [ ] - support for other Datastore types (GeoPoint, Key)
 [ ] - support for MongoDB and BSON serialization
 




