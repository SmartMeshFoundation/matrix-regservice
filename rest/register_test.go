package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/SmartMeshFoundation/matrix-regservice/models"

	"github.com/SmartMeshFoundation/matrix-regservice/params"

	"github.com/SmartMeshFoundation/Photon/utils"
)

type mockWriter struct {
	t *testing.T
}

func (m *mockWriter) Header() http.Header {
	return make(http.Header)
}
func (m *mockWriter) WriteJson(v interface{}) error {
	m.t.Logf("writejson=%s", utils.StringInterface(v, 3))
	return nil
}
func (m *mockWriter) EncodeJson(interface{}) ([]byte, error) {
	return nil, nil
}
func (m *mockWriter) WriteHeader(code int) {
	if code != http.StatusOK {
		m.t.Errorf("code is %d", code)
	}
}
func TestRegisterUser(t *testing.T) {
	models.SetupTestDB()
	key, addr := utils.MakePrivateKeyAddress()
	w := &mockWriter{
		t: t,
	}
	reg := &reg{
		LocalPart: strings.ToLower(addr.String()),
	}
	userID := fmt.Sprintf("@%s:%s", reg.LocalPart, params.MatrixDomain)
	sig, err := utils.SignData(key, []byte(userID))
	if err != nil {
		t.Error(err)
		return
	}
	reg.DisplayName = fmt.Sprintf("%s-%s", utils.APex2(addr), hex.EncodeToString(sig))
	sig, err = utils.SignData(key, []byte(params.MatrixDomain))
	if err != nil {
		t.Error(err)
		return
	}
	reg.Password = hex.EncodeToString(sig)
	verifyReg(reg, w)
}
