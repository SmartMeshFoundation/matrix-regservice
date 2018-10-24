package utils

import (
	"testing"
)

func TestHashPasswordWrapper(t *testing.T) {
	hb := HashPasswordWrapper("123456")
	t.Logf("hb=%s", string(hb))
}

func TestVerifyPassword(t *testing.T) {
	err := VerifyPassword("123456", "$2a$12$Fr4ahv6xju/38iEzZ88UW.nkQCCHQrSuuOGYSePgP2wrnPfSBKNGu")
	if err != nil {
		t.Error(err)
		return
	}
	err = VerifyPassword("123456", "$2b$12$NPgM9FqO2HZDxlL/cbinAOqqwX4D/RvvFG.CaFdegzfwq6aJ9cJEe")
	if err != nil {
		t.Error(err)
		return
	}
}
