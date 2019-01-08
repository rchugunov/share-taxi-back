package logging

import (
	"testing"
)

func TestBase64(t *testing.T) {
	encoded := base64("925_d317185e7ae25f2f:46572b152d4f3bbe79ceba443e994868b1a841c47cf13e1c27761a23d128f158")
	if encoded != "OTI1X2QzMTcxODVlN2FlMjVmMmY6NDY1NzJiMTUyZDRmM2JiZTc5Y2ViYTQ0M2U5OTQ4NjhiMWE4NDFjNDdjZjEzZTFjMjc3NjFhMjNkMTI4ZjE1OA==" {
		t.Errorf("wrong algorithm used. Result: %s\n", encoded)
	}
}
