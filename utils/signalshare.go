package utils

import (
	"bytes"
	"fmt"

	"github.com/SmartMeshFoundation/matrix-regservice/params"

	"github.com/ethereum/go-ethereum/common"

	"errors"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

//VerifyDisplayName verify user,displayname is right or not
func VerifyDisplayName(userID string, sig []byte) (err error) {
	data := []byte(fmt.Sprintf("@%s:%s", userID, params.MatrixDomain))
	if err != nil {
		return err
	}
	addr := common.HexToAddress(userID)
	hash := crypto.Keccak256Hash(data)
	sender, err := utils.Ecrecover(hash, sig)
	if err != nil {
		return err
	}
	if bytes.Compare(addr[:], sender[:]) != 0 {
		return errors.New("not match")
	}
	return nil
}

//VerifyPasswordSignature verify user,password is right or not
func VerifyPasswordSignature(addr common.Address, sig []byte) (err error) {
	data := []byte(params.MatrixDomain)
	hash := crypto.Keccak256Hash(data)
	sender, err := utils.Ecrecover(hash, sig)
	if err != nil {
		return err
	}
	if bytes.Compare(addr[:], sender[:]) != 0 {
		return errors.New("not match")
	}
	return nil
}
