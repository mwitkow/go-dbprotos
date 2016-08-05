// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package plugin

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/mwitkow/go-dbprotos"
)

func (p *plugin) generateDatastoreKind(message *descriptor.DescriptorProto) {
	ccTypeName := generator.CamelCase(message.GetName())

	emo := getEntityOptIfAny(message)
	if emo == nil {
		return
	}
	kind := emo.GetDatastore().GetKind()
	if kind == "" {
		return
	}
	p.P(`var `, ccTypeName, `_DatastoreKind = "`, kind, `"`)
}

func (p *plugin) generateDatastoreLoader(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())

	entityOpt := getEntityOptIfAny(message.DescriptorProto)
	isStrict := entityOpt.Datastore.GetStrictReading()

	p.P(`// Load implements the Google Datastore Entity Property interpreter for this type.`)
	p.P(`func (this *`, ccTypeName, `) Load(props []`, p.datastorePkg.Use(), `.Property) error {`)
	p.In()
	p.P(`for _, prop := range props {`)
	p.In()

	for _, field := range message.Field {
		//p.P("// processing field: ", field.GetName(), ` type `, field.Type.String(), "type name", field.GetTypeName())

		datastoreOpt := getDatastoreFieldOpt(field)
		if datastoreOpt == nil || datastoreOpt.GetName() == "" {
			p.P(`// field "`, field.GetName(), `" ignored due to no datastore option or no datastore name`)
			continue
		}

		p.P(`if prop.Name == "`, datastoreOpt.GetName(), `" {`)
		p.In()

		datastoreType := getDatastoreSimpleGolangType(field)
		dstType := p.getDestinationType(field)
		if field.IsRepeated() {
			p.P(`if v, ok := (prop.Value).([]interface{}); ok {`)
			p.In()
			p.P(`for _, item := range v {`)
			p.In()
			// TODO(michal): allow reading a single value field into a repated for optional -> repeated upgrades.
			if datastoreType != "" {
				p.P(`if castItem, ok := (item).(`, datastoreType, `); ok {`)
				p.In()
				p.P(`this.`, field.GetName(), ` = append(this.`, field.GetName(), `, (`, dstType, `)(castItem))`)
				p.Out()
			} else if field.GetTypeName() == ".google.protobuf.Timestamp" {
				p.P(`if castItem, ok := (item).(time.Time); ok {`)
				p.In()
				p.P(`o, _ :=`, p.golangPtypesPkg.Use(), `.TimestampProto(castItem)`)
				p.P(`this.`, field.GetName(), ` = append(this.`, field.GetName(), `, o)`)
				p.Out()
			}
			p.P(`} else {`)
			p.In()
			p.P(`return `, p.fmtPkg.Use(), `.Errorf("bad type '%t' of `, field.GetName(), ` when parsing `, ccTypeName, `", item)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
		} else {
			// Field generation part.
			if datastoreType != "" {
				p.P(`if v, ok := (prop.Value).(`, datastoreType, `); ok {`)
				p.In()
				p.P(`this.`, field.GetName(), ` = v`)
				p.Out()
			} else if field.GetTypeName() == ".google.protobuf.Timestamp" {
				p.P(`if v, ok := (prop.Value).(`, p.timePkg.Use(), `.Time); ok {`)
				p.In()
				p.P(`this.`, field.GetName(), `, _ = `, p.golangPtypesPkg.Use(), `.TimestampProto(v)`)
				p.Out()
			}
			p.P(`} else {`)
			p.In()
			p.P(`return `, p.fmtPkg.Use(), `.Errorf("bad type '%t' of `, field.GetName(), ` when parsing `, ccTypeName, `", prop.Value)`)
			p.Out()
			p.P(`}`)
		}


		// End of field generation part.
		p.P(`continue`)
		p.Out()
		p.P(`}`)

	}
	if isStrict {
		p.P(`return `, p.fmtPkg.Use(), `.Errorf("Property %v unknown for `, ccTypeName, `", prop.Name)`)
	} else {
		p.P(p.dbprotosPkg.Use(), `.UnknownFieldCallback("`, ccTypeName, `", prop.Name)`)
	}
	p.Out()
	p.P(`}`)
	p.P(`return nil`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) getDestinationType(field *descriptor.FieldDescriptorProto) string {
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
		// TODO(michal): YOLO
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

func getDatastoreFieldOpt(field *descriptor.FieldDescriptorProto) *dbprotos.DatastoreFieldOpt {
	if field.GetOptions() == nil {
		return nil
	}
	e, err := proto.GetExtension(field.GetOptions(), dbprotos.E_Datastore)
	if err != nil {
		return nil
	}
	if emo, ok := e.(*dbprotos.DatastoreFieldOpt); ok {
		return emo
	}
	return nil
}
