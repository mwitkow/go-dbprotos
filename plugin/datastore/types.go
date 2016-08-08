// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package datastore

import "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

func getGolangProtoType(field *descriptor.FieldDescriptorProto) string {
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	default:
		// TODO(michal): Figure a nicer flow here, so we check errors.
		return ""
	}
}

func getDatastoreSimpleGolangType(field *descriptor.FieldDescriptorProto) string {
	switch *(field.Type) {
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	}
	return ""
}
