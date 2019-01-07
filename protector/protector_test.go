package protector

import (
	"../protector"
	"testing"
)

func TestProtector_GenerateNextSessionKey_KeyIsNotEmpty(t *testing.T) {
	t.Parallel()

	var prot = protector.New("32165")
	var sKey = prot.GenerateNextSessionKey("6733578957")
	if sKey == "" {
		t.Fatal("Generated key is empty")
	}
}
