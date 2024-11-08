package doc

import (
	"bookshop/internal/assert"
	"strings"
	"testing"
)

func TestCreateAuthorOperation(t *testing.T) {
	got := createAuthorOperation()
	assert.NotNil(t, got)
	assert.Contains(t, strings.ToLower(got.OperationId), "create")
	assert.NotNil(t, got.RequestBody)
	assert.NotNil(t, got.Responses)
	assert.NotNil(t, got.RequestBody.Content)
	assert.NotNil(t, got.Responses.Codes)
}
