package keys

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s Store) List() (keys Keys, err error) {
	for _, key := range s.Keys.Keys {
		keys.Keys = append(keys.Keys, Key{
			Name:     key.Name,
			Operator: key.Operator,
			Address:  key.Address,
		})
	}

	return keys, err
}

func (s *Store) AddRecover(req AddRecoverRequest) (err error) {
	if _, err = s.RecoverAddressFromMnemonic(req.Mnemonic, req.Name); err != nil {
		return err
	}

	return nil
}

func newKeyringFromBackend(ctx client.Context, backend string) (keyring.Keyring, error) {
	if ctx.GenerateOnly || ctx.Simulate {
		return keyring.New(sdk.KeyringServiceName(), keyring.BackendFile, ctx.KeyringDir, ctx.Input)
	}
	return keyring.New(sdk.KeyringServiceName(), backend, ctx.KeyringDir, ctx.Input)
}

func (s *Store) addKeyring(mnemonic, name string) (keyring.Info, error) {
	var info keyring.Info
	var bip39Passphrase string
	var coinType = uint32(cosmosTypes.CoinType)
	var account = uint32(0)
	var index = uint32(0)
	var hdPath = hd.CreateHDPath(coinType, account, index).String()
	var algoStr = string(hd.Secp256k1Type)

	keyringAlgos, _ := s.context.Keyring.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	if err != nil {
		return nil, err
	}

	info, err = s.context.Keyring.NewAccount(name, mnemonic, bip39Passphrase, hdPath, algo)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s Store) RecoverAddressFromMnemonic(mnemonic, name string) (string, error) {
	var stringAddr string

	info, err := s.addKeyring(mnemonic, name)
	if err != nil {
		return stringAddr, err
	}

	stringAddr = info.GetAddress().String()

	return stringAddr, nil
}
