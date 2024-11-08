package doc

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"strings"
)

type Document struct {
	v3.Document
}

func NewDocument() *Document {
	paths := &v3.Paths{
		PathItems: orderedmap.New[string, *v3.PathItem](),
	}
	d := v3.Document{
		Version: "3.1.0",
		Info: &base.Info{
			Title:   "Bookshop API",
			Version: "1.0.0",
		},
		Paths: paths,
	}
	return &Document{Document: d}
}

func (d *Document) AddPathItem(httpMethod, endpoint string, op *Operation) {
	switch strings.ToLower(httpMethod) {
	case "post":
		d.Paths.PathItems.Set(endpoint, &v3.PathItem{
			Post: &op.Operation,
		})
	case "get":
		d.Paths.PathItems.Set(endpoint, &v3.PathItem{
			Get: &op.Operation,
		})
	case "put":
		d.Paths.PathItems.Set(endpoint, &v3.PathItem{
			Put: &op.Operation,
		})
	case "patch":
		d.Paths.PathItems.Set(endpoint, &v3.PathItem{
			Patch: &op.Operation,
		})
	case "delete":
		d.Paths.PathItems.Set(endpoint, &v3.PathItem{
			Delete: &op.Operation,
		})
	}
}

type Operation struct {
	v3.Operation
}

func NewOperation(id, desc string) *Operation {
	requestContent := orderedmap.New[string, *v3.MediaType]()
	responseCodes := orderedmap.New[string, *v3.Response]()
	op := v3.Operation{
		OperationId: id,
		Description: desc,
		RequestBody: &v3.RequestBody{
			Content: requestContent,
		},
		Responses: &v3.Responses{
			Codes: responseCodes,
		},
	}
	return &Operation{Operation: op}
}

func (o *Operation) AddResponseSchema(schema *Schema, code, contentType, desc string) {
	content := orderedmap.New[string, *v3.MediaType]()
	content.Set(contentType, &v3.MediaType{
		Schema: &schema.SchemaProxy,
	})
	o.Responses.Codes.Set(code, &v3.Response{
		Description: desc,
		Content:     content,
	})
}

func (o *Operation) AddRequestSchema(schema *Schema, contentType string, required bool) {
	o.RequestBody.Content.Set(contentType, &v3.MediaType{
		Schema: &schema.SchemaProxy,
	})
	o.RequestBody.Required = &required
}

type Schema struct {
	base.SchemaProxy
}

func NewSchema() *Schema {
	sp := *base.CreateSchemaProxy(&base.Schema{
		Type:       []string{"object"},
		Properties: orderedmap.New[string, *base.SchemaProxy](),
	})
	return &Schema{SchemaProxy: sp}
}

func (s *Schema) AddProperty(name, propType, desc string, required bool) {
	s.SchemaProxy.Schema().Properties.Set(name, base.CreateSchemaProxy(&base.Schema{
		Type:        []string{propType},
		Description: desc,
	}))
	if required {
		s.SchemaProxy.Schema().Required = append(s.SchemaProxy.Schema().Required, name)
	}
}

func GenerateOpenAPISpec() *v3.Document {
	d := NewDocument()
	d.AddPathItem("post", "/v1/authors", createAuthorOperation())
	return &d.Document
}
