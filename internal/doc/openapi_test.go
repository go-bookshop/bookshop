package doc

import (
	"bookshop/internal/assert"
	"fmt"
	"testing"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()
	assert.Equal(t, d.Version, "3.1.0")
	assert.Equal(t, d.Info.Title, "Bookshop API")
	assert.Equal(t, d.Info.Version, "1.0.0")
	assert.NotNil(t, d.Paths)
	assert.NotNil(t, d.Paths.PathItems)
	assert.NotNil(t, d.Components)
	assert.NotNil(t, d.Components.Schemas)
}

func TestDocumentAddPathItem(t *testing.T) {
	d := NewDocument()
	o := NewOperation("select", "", []string{})

	d.AddPathItem("post", "/post", o)
	d.AddPathItem("get", "/get", o)
	d.AddPathItem("put", "/put", o)
	d.AddPathItem("patch", "/patch", o)
	d.AddPathItem("delete", "/delete", o)

	val, ok := d.Paths.PathItems.Get("/post")
	assert.True(t, ok)
	assert.Equal(t, val.Post, &o.Operation)

	val, ok = d.Paths.PathItems.Get("/get")
	assert.True(t, ok)
	assert.Equal(t, val.Get, &o.Operation)

	val, ok = d.Paths.PathItems.Get("/put")
	assert.True(t, ok)
	assert.Equal(t, val.Put, &o.Operation)

	val, ok = d.Paths.PathItems.Get("/patch")
	assert.True(t, ok)
	assert.Equal(t, val.Patch, &o.Operation)

	val, ok = d.Paths.PathItems.Get("/delete")
	assert.True(t, ok)
	assert.Equal(t, val.Delete, &o.Operation)
}

func TestNewOperation(t *testing.T) {
	o := NewOperation("select", "desc", []string{"testTag"})
	assert.Equal(t, "select", o.OperationId)
	assert.Equal(t, "desc", o.Description)
	assert.StringContains(t, "testTag", o.Tags[0])
	assert.NotNil(t, o.RequestBody)
	assert.NotNil(t, o.RequestBody.Content)
	assert.NotNil(t, o.Responses)
	assert.NotNil(t, o.Responses.Codes)
}

func TestOperationAddResponseSchema(t *testing.T) {
	o := NewOperation("select", "desc", []string{})
	s := NewSchema()
	o.AddResponseSchema(s, "200", "text/text", "test2")

	code, ok := o.Responses.Codes.Get("200")
	assert.True(t, ok)
	assert.Equal(t, code.Description, "test2")
	assert.NotNil(t, code.Content)

	content, ok := code.Content.Get("text/text")
	assert.True(t, ok)
	assert.Equal(t, content.Schema, &s.SchemaProxy)
}

func TestNewSchema(t *testing.T) {
	s := NewSchema()
	assert.NotNil(t, s.SchemaProxy)
	assert.NotNil(t, s.SchemaProxy.Schema().Properties)
	assert.Equal(t, s.SchemaProxy.Schema().Type[0], "object")
}

func TestNewRefSchema(t *testing.T) {
	schemaName := "TestSchema"

	s := NewRefSchema(schemaName)

	assert.NotNil(t, s.SchemaProxy)
	assert.True(t, s.SchemaProxy.IsReference())
	assert.Equal(t, s.SchemaProxy.GetReference(), fmt.Sprintf("#/components/schemas/%s", schemaName))
}

func TestOperationAddRequestSchema(t *testing.T) {
	o := NewOperation("select", "desc", []string{})
	s := NewSchema()
	o.AddRequestSchema(s, "text/text", true)

	content, ok := o.RequestBody.Content.Get("text/text")
	assert.True(t, ok)
	assert.Equal(t, content.Schema, &s.SchemaProxy)
	assert.True(t, *o.RequestBody.Required)
}

func TestOperationAddQueryParameter(t *testing.T) {
	o := NewOperation("select", "desc", []string{})

	paramName := "test"
	paramType := "integer"
	paramDesc := "test desc"
	required := false
	o.AddQueryParameter(paramName, paramType, paramDesc, required)

	var param *v3.Parameter
	for _, p := range o.Parameters {
		if p.Name == paramName && p.In == "query" {
			param = p
			break
		}
	}

	assert.NotNil(t, param)
	assert.Equal(t, param.Name, paramName)
	assert.Equal(t, param.Description, paramDesc)
	assert.Equal(t, *param.Required, required)
	assert.SliceContains(t, param.Schema.Schema().Type, paramType)
}

func TestSchemaAddProperty(t *testing.T) {
	s := NewSchema()
	s.AddProperty("post", "string", "desc", true)

	val, ok := s.SchemaProxy.Schema().Properties.Get("post")
	assert.True(t, ok)
	assert.Equal(t, val.Schema().Type[0], "string")
	assert.Equal(t, val.Schema().Description, "desc")
	assert.Equal(t, s.SchemaProxy.Schema().Required[0], "post")
}

func TestSchemaAddSimpleArrayProperty(t *testing.T) {
	s := NewSchema()
	propertyName := "tags"
	itemType := "string"
	description := "An array of tags"
	required := true

	s.AddSimpleArrayProperty(propertyName, itemType, description, required)

	val, exists := s.SchemaProxy.Schema().Properties.Get(propertyName)
	assert.True(t, exists)
	assert.SliceContains(t, val.Schema().Type, "array")
	assert.Equal(t, description, val.Schema().Description)
	assert.SliceContains(t, val.Schema().Items.A.Schema().Type, itemType)
	assert.SliceContains(t, s.SchemaProxy.Schema().Required, propertyName)
}

func TestSchemaAddSchemaArrayProperty(t *testing.T) {
	mainSchema := NewSchema()
	itemSchema := NewSchema()

	itemPropertyName := "name"
	itemSchema.AddProperty(itemPropertyName, "string", "Name of the item", true)

	propertyName := "items"
	description := "An array of item schemas"
	required := true

	mainSchema.AddSchemaArrayProperty(propertyName, description, itemSchema, required)

	val, exists := mainSchema.SchemaProxy.Schema().Properties.Get(propertyName)
	assert.True(t, exists)
	assert.SliceContains(t, val.Schema().Type, "array")
	assert.Equal(t, val.Schema().Description, description)
	assert.Equal(t, *val.Schema().Items.A, itemSchema.SchemaProxy)
	assert.SliceContains(t, mainSchema.SchemaProxy.Schema().Required, propertyName)
}

func TestAddEnumProperty(t *testing.T) {
	s := NewSchema()
	s.AddEnumProperty("name", "integer", "desc", true, "1", "2")

	val, ok := s.SchemaProxy.Schema().Properties.Get("name")
	assert.True(t, ok)
	assert.Equal(t, val.Schema().Description, "desc")
	assert.NotNil(t, val.Schema().Enum)
	assert.SliceContains(t, val.Schema().Type, "integer")
	assert.SliceContains(t, s.Schema().Required, "name")
}

func TestGenerateOpenAPISpec(t *testing.T) {
	d := GenerateOpenAPISpec()
	assert.NotNil(t, d.Paths.PathItems)
	assert.NotNil(t, d.Components.Schemas)
}
