package chatgpt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadBody(t *testing.T) {
	body := []byte("ab\n\ndata: {}  \nfghil\ndata: {a, b}\nel")
	// body := []byte("ab\n\n cde  \nfghil\nel")
	// lines, err := readBody(bytes.NewReader(body))
	res := readBody(bytes.NewReader(body))
	assert.Equal(t, []byte("{a, b}"), res)

	// assert.Equal(t, 4, len(lines))
	// assert.Equal(t, []byte("ab"), lines[0])
	// assert.Equal(t, []byte(" cde  "), lines[1])
	// assert.Equal(t, []byte("fghil"), lines[2])
	// assert.Equal(t, []byte("el"), lines[3])
}
