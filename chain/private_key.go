package chain

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/go-bip39"
)

const (
	DefaultHDPath = "m/44’/118’/0/0/0"
)

// KeyRing 包含区块链账户的密钥信息
type KeyRing struct {
	Address  string             `json:"address"`  //cosmos地址
	Mnemonic string             `json:"mnemonic"` //24词助记词
	Private  *secp256k1.PrivKey `json:"private"`  //私钥
	Public   *secp256k1.PubKey  `json:"public"`   //公钥
	HDPath   string             `json:"hd_path"`  //HD路径
}

// NewDefaultKey 创建一个使用默认前缀和HD路径的新密钥
// 返回生成的密钥信息和可能的错误
func NewDefaultKey() (KeyRing, error) {
	return NewKey(DefaultAddressPrefix, DefaultHDPath)
}

// NewKey 根据指定的地址前缀和HD路径创建新的密钥对
// prefix: 地址前缀，如"cosmos"
// hdPath: HD钱包路径，如"m/44'/118'/0/0/0"
// 返回生成的密钥信息和可能的错误
func NewKey(prefix, hdPath string) (KeyRing, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return KeyRing{}, err
	}
	// 生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return KeyRing{}, err
	}
	// 根据助记词和HD路径生成私钥
	pBytes, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return KeyRing{}, err
	}

	// 根据助记词生成私钥
	private := &secp256k1.PrivKey{Key: pBytes}

	// 获取公钥
	pubKey := private.PubKey()

	//  公钥 -> 地址 (Bech32 格式)
	addr := pubKey.Address()
	bech32Addr, _ := bech32.ConvertAndEncode(prefix, addr)

	return KeyRing{
		Address:  bech32Addr,
		Mnemonic: mnemonic,
		Private:  private,
		Public:   pubKey.(*secp256k1.PubKey),
		HDPath:   hdPath,
	}, nil
}
