package doc

import (
	"bookshop/internal/assert"
	"strings"
	"testing"
)

func TestCreateCategoryOperation(t *testing.T) {
	got := createCategoryOperation()
	assert.NotNil(t, got)
	assert.StringContains(t, strings.ToLower(got.OperationId), "create")
	assert.SliceContains(t, got.Tags, CategoriesTag)
	assert.NotNil(t, got.RequestBody)
	assert.NotNil(t, got.Responses)
	assert.NotNil(t, got.RequestBody.Content)
	assert.NotNil(t, got.Responses.Codes)
}
