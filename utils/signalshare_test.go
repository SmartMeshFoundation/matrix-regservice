package utils

import (
	"encoding/hex"
	"testing"

	"github.com/SmartMeshFoundation/matrix-regservice/params"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"

	"fmt"
)

func TestCreatePasswordAndVerify(t *testing.T) {
	key, addr := utils.MakePrivateKeyAddress()
	userID := addr.String()
	data := fmt.Sprintf("@%s:%s", userID, params.MatrixDomain)
	sig, err := utils.SignData(key, []byte(data))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("addr=%s,sig=%s\n", addr.String(), hex.EncodeToString(sig))
	err = VerifyDisplayName(addr.String(), sig)
	if err != nil {
		t.Error(err)
	}
	sig, err = utils.SignData(key, []byte(params.MatrixDomain))
	if err != nil {
		t.Error(err)
		return
	}
	err = VerifyPasswordSignature(addr, sig)
	if err != nil {
		t.Error(err)
		return
	}
}
