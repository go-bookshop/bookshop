package doc

import (
	"bookshop/internal/assert"
	"strings"
	"testing"
)

func TestGetBooksOperation(t *testing.T) {
	got := getBooksOperation()
	assert.NotNil(t, got)
	assert.StringContains(t, strings.ToLower(got.OperationId), "get")
	assert.SliceContains(t, got.Tags, BooksTag)
	assert.NotNil(t, got.RequestBody)
	assert.NotNil(t, got.Responses)
	assert.NotNil(t, got.RequestBody.Content)
	assert.NotNil(t, got.Responses.Codes)
}
