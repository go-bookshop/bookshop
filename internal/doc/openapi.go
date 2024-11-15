package doc

import (
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type Document struct {
	v3.Document
}

func NewDocument() *Document {
	paths := &v3.Paths{
		PathItems: orderedmap.New[string, *v3.PathItem](),
	}
	components := &v3.Components{
		Schemas: orderedmap.New[string, *base.SchemaProxy](),
	}
	d := v3.Document{
		Version: "3.1.0",
		Info: &base.Info{
			Title:   "Bookshop API",
			Version: "1.0.0",
		},
		Paths:      paths,
		Components: components,
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

func (d *Document) AddComponent(compName string, compSchema *Schema) {
	d.Components.Schemas.Set(compName, &compSchema.SchemaProxy)
}

type Operation struct {
	v3.Operation
}

func NewOperation(id, desc string, tags []string) *Operation {
	requestContent := orderedmap.New[string, *v3.MediaType]()
	responseCodes := orderedmap.New[string, *v3.Response]()
	op := v3.Operation{
		OperationId: id,
		Description: desc,
		Tags:        tags,
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

func (o *Operation) AddQueryParameter(name, paramType, desc string, required bool) {
	schema := &base.Schema{
		Type: []string{paramType},
	}

	parameter := &v3.Parameter{
		Name:        name,
		In:          "query",
		Description: desc,
		Required:    &required,
		Schema:      base.CreateSchemaProxy(schema),
	}

	o.Parameters = append(o.Parameters, parameter)
}

type Schema struct {
	base.SchemaProxy
}

var NoBodySchema = &Schema{SchemaProxy: *base.CreateSchemaProxy(
	&base.Schema{
		Type: []string{"object"},
	})}

func NewSchema() *Schema {
	sp := *base.CreateSchemaProxy(&base.Schema{
		Type:       []string{"object"},
		Properties: orderedmap.New[string, *base.SchemaProxy](),
	})
	return &Schema{SchemaProxy: sp}
}

func NewRefSchema(schemaName string) *Schema {
	schemaProxyRef := base.CreateSchemaProxyRef(fmt.Sprintf("#/components/schemas/%s", schemaName))

	return &Schema{*schemaProxyRef}
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

func (s *Schema) AddSimpleArrayProperty(name, itemType, desc string, required bool) {
	itemSchema := &base.Schema{
		Type: []string{itemType},
	}

	s.SchemaProxy.Schema().Properties.Set(name, base.CreateSchemaProxy(&base.Schema{
		Type:        []string{"array"},
		Description: desc,

		Items: &base.DynamicValue[*base.SchemaProxy, bool]{
			A: base.CreateSchemaProxy(itemSchema),
		},
	}))
	if required {
		s.SchemaProxy.Schema().Required = append(s.SchemaProxy.Schema().Required, name)
	}
}

func (s *Schema) AddSchemaArrayProperty(name, desc string, itemSchema *Schema, required bool) {
	s.SchemaProxy.Schema().Properties.Set(name, base.CreateSchemaProxy(&base.Schema{
		Type:        []string{"array"},
		Description: desc,

		Items: &base.DynamicValue[*base.SchemaProxy, bool]{
			A: &itemSchema.SchemaProxy,
		},
	}))
	if required {
		s.SchemaProxy.Schema().Required = append(s.SchemaProxy.Schema().Required, name)
	}
}

func GenerateOpenAPISpec() *v3.Document {
	d := NewDocument()

	d.AddComponent("BookItem", getBookItemComponentSchema())
	d.AddComponent("BookProperty", getBookPropertySchema())

	d.AddPathItem("post", "/v1/authors", createAuthorOperation())
	d.AddPathItem("post", "/v1/books/categories", createCategoryOperation())
	d.AddPathItem("post", "/v1/users", registerUserOperation())
	d.AddPathItem("put", "/v1/users/activate", activateUserOperation())
	d.AddPathItem("post", "/v1/users/activation/resend-token", resendTokenOperation())
	d.AddPathItem("get", "/v1/books", getBooksOperation())

	return &d.Document
}
