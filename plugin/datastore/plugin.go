// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package datastore

import (
	"fmt"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/mwitkow/go-dbprotos/plugin"
	"github.com/mwitkow/go-dbprotos"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/proto"
)

func init() {
	generator.RegisterPlugin(NewDatastorePlugin())
}

type datastorePlugin struct {
	*generator.Generator
	generator.PluginImports
	regexPkg        generator.Single
	fmtPkg          generator.Single
	protoPkg        generator.Single
	dbprotosPkg     generator.Single
	datastorePkg    generator.Single
	golangPtypesPkg generator.Single

	timePkg generator.Single
}

func NewDatastorePlugin() generator.Plugin {
	return &datastorePlugin{}
}

func (p *datastorePlugin) Name() string {
	return "dbprotos_datastore"
}

func (p *datastorePlugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *datastorePlugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	//p.regexPkg = p.NewImport("regexp")
	p.fmtPkg = p.NewImport("fmt")
	p.dbprotosPkg = p.NewImport("github.com/mwitkow/go-dbprotos")
	p.datastorePkg = p.NewImport("google.golang.org/cloud/datastore")
	p.timePkg = p.NewImport("time")
	p.golangPtypesPkg = p.NewImport("github.com/golang/protobuf/ptypes")

	for _, msg := range file.Messages() {
		if plugin.GetEntityOptIfAny(msg.DescriptorProto) == nil {
			continue
		}
		if !gogoproto.IsProto3(file.FileDescriptorProto) {
			p.Error(fmt.Errorf("The dbprotos_datastore plugin only works on proto3 files. %s", file.GetName()))
			return
		}
		p.generateDatastoreKind(msg.DescriptorProto)
		p.generateDatastoreLoader(file, msg)
		p.generateDatastoreSaver(file, msg)

	}
}

// getDatastoreFieldOptIfAny returns the Field-level option describing the datastore access info.
func getDatastoreFieldOptIfAny(field *descriptor.FieldDescriptorProto) *dbprotos.DatastoreFieldOpt {
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

