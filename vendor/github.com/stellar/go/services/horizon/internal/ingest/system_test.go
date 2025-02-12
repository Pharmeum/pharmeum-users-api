package ingest

import (
	"fmt"
	"testing"

	"github.com/stellar/go/network"
	"github.com/stellar/go/services/horizon/internal/test"
)

func TestBackfill(t *testing.T) {
	tt := test.Start(t).ScenarioWithoutHorizon("kahuna")
	defer tt.Finish()
	is := sys(tt, Config{EnableAssetStats: false, CursorName: "HORIZON"})

	err := is.ReingestSingle(10)
	tt.Require.NoError(err)
	tt.UpdateLedgerState()

	// ensure 1 ledger
	var found int
	err = tt.HorizonSession().GetRaw(&found, "SELECT COUNT(*) FROM history_ledgers")
	tt.Require.NoError(err)
	tt.Require.Equal(1, found)

	err = is.Backfill(3)
	if tt.Assert.NoError(err) {
		err = tt.HorizonSession().GetRaw(&found, "SELECT COUNT(*) FROM history_ledgers")
		tt.Require.NoError(err)

		tt.Assert.Equal(4, found, "expected 4 ledgers to be in history database, but got %d", found)
	}
}

func TestClearAll(t *testing.T) {
	tt := test.Start(t).Scenario("kahuna")
	defer tt.Finish()
	is := sys(tt, Config{EnableAssetStats: false, CursorName: "HORIZON"})

	err := is.ClearAll()

	tt.Require.NoError(err)

	// ensure all tables are cleared
	tables := []TableName{
		AssetStatsTableName,
		AccountsTableName,
		AssetsTableName,
		EffectsTableName,
		LedgersTableName,
		OperationParticipantsTableName,
		OperationsTableName,
		TradesTableName,
		TransactionParticipantsTableName,
		TransactionsTableName,
	}

	for _, tableName := range tables {
		ensureEmpty(tt, tableName)
	}
}

func ensureEmpty(tt *test.T, tableName TableName) {
	var found int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", string(tableName))
	err := tt.HorizonSession().GetRaw(&found, query)
	tt.Require.NoError(err)
	tt.Assert.Equal(0, found)
}

func TestValidation(t *testing.T) {
	tt := test.Start(t).Scenario("kahuna")
	defer tt.Finish()

	sys := New(network.TestNetworkPassphrase, "", tt.CoreSession(), tt.HorizonSession(), Config{})

	// intact chain
	for i := int32(2); i <= 57; i++ {
		tt.Assert.NoError(sys.validateLedgerChain(i))
	}

	// missing cur
	_, err := tt.CoreSession().ExecRaw(
		`DELETE FROM ledgerheaders WHERE ledgerseq = ?`, 5,
	)
	tt.Require.NoError(err)

	err = sys.validateLedgerChain(5)
	if tt.Assert.Error(err) {
		tt.Assert.Contains(err.Error(), "failed to load cur ledger")
	}

	// missing prev
	_, err = tt.HorizonSession().ExecRaw(
		`DELETE FROM history_ledgers WHERE sequence = ?`, 5,
	)
	tt.Require.NoError(err)

	err = sys.validateLedgerChain(6)
	if tt.Assert.Error(err) {
		tt.Assert.Contains(err.Error(), "failed to load prev ledger")
	}

	// mismatched header in core
	_, err = tt.CoreSession().ExecRaw(`
		UPDATE ledgerheaders
		SET prevhash = ?
		WHERE ledgerseq = ?`, "00000", 8)
	tt.Require.NoError(err)

	err = sys.validateLedgerChain(8)
	if tt.Assert.Error(err) {
		tt.Assert.Contains(err.Error(), "cur and prev ledger hashes don't match")
	}

	// mismatched header in horizon
	_, err = tt.HorizonSession().ExecRaw(`
		UPDATE history_ledgers
		SET ledger_hash = ?
		WHERE sequence = ?`, "00000", 10)
	tt.Require.NoError(err)

	err = sys.validateLedgerChain(11)
	if tt.Assert.Error(err) {
		tt.Assert.Contains(err.Error(), "cur and prev ledger hashes don't match")
	}
}
