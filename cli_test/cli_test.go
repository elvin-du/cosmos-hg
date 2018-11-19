package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"
)

var (
	hashgarddHome = ""
	hashgardcliHome   = ""
)

func init() {
	hashgarddHome, hashgardcliHome = getTestingHomeDirs()
}

func TestInitStartSequence(t *testing.T) {
	os.RemoveAll(hashgarddHome)
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	executeInit(t)
	executeStart(t, servAddr, port)
}

func executeInit(t *testing.T) {
	var (
		chainID string
		initRes map[string]json.RawMessage
	)
	_, stderr := tests.ExecuteT(t, fmt.Sprintf("hashgardd --home=%s --home-client=%s init --name=test", hashgarddHome, hashgardcliHome), app.DefaultKeyPass)
	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(t, err)
	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)
}

func executeStart(t *testing.T, servAddr, port string) {
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("hashgardd start --home=%s --rpc.laddr=%v", hashgarddHome, servAddr))
	defer proc.Stop(false)
	tests.WaitForTMStart(port)
}

func getTestingHomeDirs() (string, string) {
	tmpDir := os.TempDir()
	hashgarddHome := fmt.Sprintf("%s%s.test_hashgardd", tmpDir, string(os.PathSeparator))
	hashgardcliHome := fmt.Sprintf("%s%s.test_hashgardcli", tmpDir, string(os.PathSeparator))
	return hashgarddHome, hashgardcliHome
}
