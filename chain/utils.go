package chain

import (
	"encoding/base64"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PubToAddress 公钥转换为地址
func PubToAddress(publicKey string) string {
	// base64 -> pubkey bytes  eg:AqbnviTZBMHO0lp7X0S/9uSHgXT3Hjlx5pw5Fhu80K5y
	pubBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Fatalf("base64 decode failed: %v", err)
	}

	// 推导公钥
	pub := secp256k1.PubKey{Key: pubBytes}

	// 公钥转为地址
	return sdk.AccAddress(pub.Address()).String()
}
