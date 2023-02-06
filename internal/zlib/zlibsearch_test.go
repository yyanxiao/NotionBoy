package zlib

import "testing"

func TestParseZlibReq(t *testing.T) {
	tests := []struct {
		query             string
		expectedName      string
		expectedExtension string
	}{
		{"阅读 #pdf #2020", "阅读", "pdf"},
		{"阅读 #epub #2021", "阅读", "epub"},
		{"阅读 #cbz", "阅读", ""},
		{"阅读", "阅读", ""},
	}

	for _, test := range tests {
		req := parseZlibReq(test.query)

		if req.Name != test.expectedName {
			t.Errorf("Expected Name %s, but got %s", test.expectedName, req.Name)
		}
		if req.Externsion != test.expectedExtension {
			t.Errorf("Expected Extension %s, but got %s", test.expectedExtension, req.Externsion)
		}
	}
}
