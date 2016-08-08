// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package plugin

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/mwitkow/go-dbprotos"
)

// GetEntityOptIfAny returns the Message-level option describing the message to entity mapping.
func GetEntityOptIfAny(message *descriptor.DescriptorProto) *dbprotos.EntityMessageOpt {
	if message.GetOptions() == nil {
		return nil
	}
	e, err := proto.GetExtension(message.GetOptions(), dbprotos.E_Entity)
	if err != nil {
		return nil
	}
	if emo, ok := e.(*dbprotos.EntityMessageOpt); ok {
		return emo
	}
	return nil
}

// GetIndexFieldOptIfAny returns the Field-level option describing the indexing properties of the
func GetIndexFieldOptIfAny(field *descriptor.FieldDescriptorProto) *dbprotos.IndexFieldOpt {
	if field.GetOptions() == nil {
		return nil
	}
	e, err := proto.GetExtension(field.GetOptions(), dbprotos.E_Index)
	if err != nil {
		return nil
	}
	if emo, ok := e.(*dbprotos.IndexFieldOpt); ok {
		return emo
	}
	return nil
}
