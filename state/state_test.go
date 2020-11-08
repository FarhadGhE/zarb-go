package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func mockState(t *testing.T) (State, crypto.Address) {
	_, pb, _ := crypto.RandomKeyPair()
	addr := pb.Address()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	val := validator.NewValidator(pb, 1)
	gen := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	loggerConfig := logger.DefaultConfig()
	loggerConfig.Levels["default"] = "error"
	logger.InitLogger(loggerConfig)
	stateConfig := DefaultConfig()
	stateConfig.Store.Path = util.TempDirName()
	txPoolConfig := txpool.DefaultConfig()
	txPool, err := txpool.NewTxPool(txPoolConfig, make(chan *message.Message, 10))
	require.NoError(t, err)
	st, err := LoadOrNewState(stateConfig, gen, val.Address(), txPool)
	require.NoError(t, err)
	return st, addr
}

func TestBlockValidate(t *testing.T) {
	st, _ := mockState(t)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}