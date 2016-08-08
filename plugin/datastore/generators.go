// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package datastore

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/mwitkow/go-dbprotos/plugin"
)

func (p *datastorePlugin) generateDatastoreKind(message *descriptor.DescriptorProto) {
	ccTypeName := generator.CamelCase(message.GetName())
	emo := plugin.GetEntityOptIfAny(message)
	if emo == nil {
		return
	}
	kind := emo.GetDatastore().GetKind()
	if kind == "" {
		return
	}
	p.P(`var `, ccTypeName, `_DatastoreKind = "`, kind, `"`)
}

func (p *datastorePlugin) generateDatastoreLoader(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())

	entityOpt := plugin.GetEntityOptIfAny(message.DescriptorProto)
	isStrict := entityOpt.Datastore.GetStrictReading()

	p.P(`// Save implements the Google Datastore Entity Property interpreter for this type.`)
	p.P(`func (this *`, ccTypeName, `) Load(props []`, p.datastorePkg.Use(), `.Property) error {`)
	p.In()
	p.P(`for _, prop := range props {`)
	p.In()

	for _, field := range message.Field {
		//p.P("// processing field: ", field.GetName(), ` type `, field.Type.String(), "type name", field.GetTypeName())

		datastoreOpt := getDatastoreFieldOptIfAny(field)
		if datastoreOpt == nil || datastoreOpt.GetName() == "" {
			p.P(`// field "`, field.GetName(), `" ignored due to no datastore option or no datastore name`)
			continue
		}

		p.P(`if prop.Name == "`, datastoreOpt.GetName(), `" {`)
		p.In()

		datastoreType := getDatastoreSimpleGolangType(field)
		protoType := getGolangProtoType(field)
		if field.IsRepeated() {
			p.P(`if v, ok := (prop.Value).([]interface{}); ok {`)
			p.In()
			p.P(`for _, item := range v {`)
			p.In()
			// TODO(michal): allow reading a single value field into a repated for optional -> repeated upgrades.
			if datastoreType != "" {
				p.P(`if castItem, ok := (item).(`, datastoreType, `); ok {`)
				p.In()
				p.P(`this.`, field.GetName(), ` = append(this.`, field.GetName(), `, (`, protoType, `)(castItem))`)
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

func (p *datastorePlugin) generateDatastoreSaver(file *generator.FileDescriptor, message *generator.Descriptor) {
	ccTypeName := generator.CamelCaseSlice(message.TypeName())

	p.P(`// Load implements the Google Datastore Entity Property interpreter for this type.`)
	p.P(`func (this *`, ccTypeName, `) Save() ([]`, p.datastorePkg.Use(), `.Property, error) {`)
	p.In()
	p.P(`props := []`, p.datastorePkg.Use(), `.Property{}`)

	for _, field := range message.Field {

		indexOpt := plugin.GetIndexFieldOptIfAny(field)
		datastoreOpt := getDatastoreFieldOptIfAny(field)
		if datastoreOpt == nil || datastoreOpt.GetName() == "" {
			p.P(`// field "`, field.GetName(), `" ignored due to no datastore option or no datastore name`)
			continue
		}
		if datastoreOpt.GetNotWriteable() {
			p.P(`// IGNORED field "`, field.GetName(), `" due to not writeable option set.`)
			continue
		}

		p.P(`{`)
		p.In()
		p.P(`prop := `, p.datastorePkg.Use(), `.Property{Name: "`, datastoreOpt.GetName(), `", NoIndex: `, !indexOpt.GetSingle(), `}`)
		datastoreType := getDatastoreSimpleGolangType(field)

		if field.IsRepeated() {
			p.P(`arrValue := []interface{}{}`)
			p.P(`for _, item := range this.`, field.GetName(), ` {`)
			p.In()
			if datastoreType != "" {
				p.P(`value := (`, datastoreType, `)(item)`)
			} else if field.GetTypeName() == ".google.protobuf.Timestamp" {
				p.P(`value, _ :=`, p.golangPtypesPkg.Use(), `.Timestamp(item)`)
			}
			p.P(`arrValue = append(arrValue, (interface{})(value))`)
			p.Out()
			p.P(`}`)
			p.P(`prop.Value = arrValue`)
		} else {
			if datastoreType != "" {
				p.P(`prop.Value = (`, datastoreType, `)(this.`, field.GetName(), `)`)
			} else if field.GetTypeName() == ".google.protobuf.Timestamp" {
				p.P(`prop.Value, _ =`, p.golangPtypesPkg.Use(), `.Timestamp(this.`, field.GetName(), `)`)
			}
		}
		p.P(`props = append(props, prop)`)
		p.Out()
		p.P(`}`)

	}
	p.P(`return props, nil`)
	p.Out()
	p.P(`}`)
}
