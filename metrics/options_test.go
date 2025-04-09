package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTags(t *testing.T) {
	options := &options{
		tags: []string{"a", "b"},
	}

	err := WithGlobalTags("c")(options)
	assert.NoError(t, err)

	assert.Equal(t, []string{"a", "b", "c"}, options.tags)

	err = WithGlobalTags("d")(options)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d"}, options.tags)
}
