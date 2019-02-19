package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestMacSecret(t *testing.T) {
	ast := assert.New(t)
	params.MatrixShareSecret = "shared"
	r := &regHelper{
		UserName: "sid1",
		Password: "pass",
		Admin:    false,
	}
	r.Nonce = "b472145df6b07d1b2cfd08e748d432ff1c085f97e6f8f36bbc1ea6205e4e78677d97cd519e812510d95a2610998b1cb1ca3edbf957ab5e3f8f331fc455d37ccf"
	ast.EqualValues("d76fcd0674283fa57f17fb56cbc42b2208144dcb",
		hmacRegHelper(r))
	r.Nonce = "13a07e012864a129352a4aee9e9d0bd34bfbd25df843efe7c23f792121c9569e166e94ba785e3698f76dff3617710b8af41ed796db2a4f345ac0e1430f5a1b69"
	ast.EqualValues("e9b320b3235aa0cb9350e7a8016cd26e5c37a836",
		hmacRegHelper(r))
}
