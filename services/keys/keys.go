package keys

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	"os"
)

type Key struct {
	Name     string `json:"Name"`
	Operator string `json:"Operator"`
	Address  string `json:"Address"`
}

type Keys struct {
	Keys []Key `json:"Keys"`
}

type AddRecoverRequest struct {
	Name     string `json:"Name"`
	Mnemonic string `json:"Mnemonic"`
}

type DeleteRequest struct {
	KeyNames []string `json:"Names"`
}

type Store struct {
	Keys    Keys
	context *client.Context
	homeDir string
}

var Wallet Store

func Load(sentinelPath string) (err error) {
	Wallet, err = New(sentinelPath)
	if err != nil {
		return err
	}

	return nil
}

func New(sentinelPath string) (Store, error) {
	var store Store

	store.homeDir = sentinelPath
	if err := store.initBackend(); err != nil {
		return store, err
	}

	return store, nil
}

func (s *Store) initBackend() error {
	hubtypes.GetConfig().Seal()

	var config = hub.MakeEncodingConfig()
	var clientCtx = client.Context{}.
		WithJSONMarshaler(config.Marshaler).
		WithInterfaceRegistry(config.InterfaceRegistry).
		WithTxConfig(config.TxConfig).
		WithLegacyAmino(config.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(s.homeDir).
		WithKeyringDir(s.homeDir)

	kr, err := newKeyringFromBackend(clientCtx, keyring.BackendTest)
	if err != nil {
		return err
	}

	clientCtx = clientCtx.WithKeyring(kr)
	s.context = &clientCtx

	return nil
}
