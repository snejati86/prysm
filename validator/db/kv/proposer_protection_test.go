package kv

import (
	"context"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func TestProposalHistoryForSlot_InitializesNewPubKeys(t *testing.T) {
	pubkeys := [][48]byte{{30}, {25}, {20}}
	db := setupDB(t, pubkeys)

	for _, pub := range pubkeys {
		signingRoot, _, err := db.ProposalHistoryForSlot(context.Background(), pub, 0)
		require.NoError(t, err)
		expected := bytesutil.PadTo([]byte{}, 32)
		require.DeepEqual(t, expected, signingRoot[:], "Expected proposal history slot signing root to be empty")
	}
}

func TestNewProposalHistoryForSlot_ReturnsNilIfNoHistory(t *testing.T) {
	valPubkey := [48]byte{1, 2, 3}
	db := setupDB(t, [][48]byte{})

	_, proposalExists, err := db.ProposalHistoryForSlot(context.Background(), valPubkey, 0)
	require.NoError(t, err)
	assert.Equal(t, false, proposalExists)
}

func TestSaveProposalHistoryForSlot_OK(t *testing.T) {
	pubkey := [48]byte{3}
	db := setupDB(t, [][48]byte{pubkey})

	slot := types.Slot(2)

	err := db.SaveProposalHistoryForSlot(context.Background(), pubkey, slot, []byte{1})
	require.NoError(t, err, "Saving proposal history failed: %v")
	signingRoot, _, err := db.ProposalHistoryForSlot(context.Background(), pubkey, slot)
	require.NoError(t, err, "Failed to get proposal history")

	require.NotNil(t, signingRoot)
	require.DeepEqual(t, bytesutil.PadTo([]byte{1}, 32), signingRoot[:], "Expected DB to keep object the same")
}

func TestNewProposalHistoryForPubKey_ReturnsEmptyIfNoHistory(t *testing.T) {
	valPubkey := [48]byte{1, 2, 3}
	db := setupDB(t, [][48]byte{})

	proposalHistory, err := db.ProposalHistoryForPubKey(context.Background(), valPubkey)
	require.NoError(t, err)
	assert.DeepEqual(t, make([]*Proposal, 0), proposalHistory)
}

func TestSaveProposalHistoryForPubKey_OK(t *testing.T) {
	pubkey := [48]byte{3}
	db := setupDB(t, [][48]byte{pubkey})

	slot := types.Slot(2)

	root := [32]byte{1}
	err := db.SaveProposalHistoryForSlot(context.Background(), pubkey, slot, root[:])
	require.NoError(t, err, "Saving proposal history failed: %v")
	proposalHistory, err := db.ProposalHistoryForPubKey(context.Background(), pubkey)
	require.NoError(t, err, "Failed to get proposal history")

	require.NotNil(t, proposalHistory)
	want := []*Proposal{
		{
			Slot:        slot,
			SigningRoot: root[:],
		},
	}
	require.DeepEqual(t, want[0], proposalHistory[0])
}

func TestSaveProposalHistoryForSlot_Overwrites(t *testing.T) {
	pubkey := [48]byte{0}
	tests := []struct {
		signingRoot []byte
	}{
		{
			signingRoot: bytesutil.PadTo([]byte{1}, 32),
		},
		{
			signingRoot: bytesutil.PadTo([]byte{2}, 32),
		},
		{
			signingRoot: bytesutil.PadTo([]byte{3}, 32),
		},
	}

	for _, tt := range tests {
		db := setupDB(t, [][48]byte{pubkey})
		err := db.SaveProposalHistoryForSlot(context.Background(), pubkey, 0, tt.signingRoot)
		require.NoError(t, err, "Saving proposal history failed")
		proposalHistory, err := db.ProposalHistoryForPubKey(context.Background(), pubkey)
		require.NoError(t, err, "Failed to get proposal history")

		require.NotNil(t, proposalHistory)
		require.DeepEqual(t, tt.signingRoot, proposalHistory[0].SigningRoot, "Expected DB to keep object the same")
		require.NoError(t, db.Close(), "Failed to close database")
	}
}

func TestPruneProposalHistoryBySlot_OK(t *testing.T) {
	slotsPerEpoch := params.BeaconConfig().SlotsPerEpoch
	wsPeriod := params.BeaconConfig().WeakSubjectivityPeriod
	pubKey := [48]byte{0}
	tests := []struct {
		slots        []types.Slot
		storedSlots  []types.Slot
		removedSlots []types.Slot
	}{
		{
			// Go 2 epochs past pruning point.
			slots:        []types.Slot{slotsPerEpoch / 2, slotsPerEpoch*5 + 6, slotsPerEpoch.Mul(uint64(wsPeriod+3)) + 8},
			storedSlots:  []types.Slot{slotsPerEpoch*5 + 6, slotsPerEpoch.Mul(uint64(wsPeriod+3)) + 8},
			removedSlots: []types.Slot{slotsPerEpoch / 2},
		},
		{
			// Go 10 epochs past pruning point.
			slots: []types.Slot{
				slotsPerEpoch + 4,
				slotsPerEpoch * 2,
				slotsPerEpoch * 3,
				slotsPerEpoch * 4,
				slotsPerEpoch * 5,
				slotsPerEpoch.Mul(uint64(wsPeriod+10)) + 8,
			},
			storedSlots: []types.Slot{slotsPerEpoch.Mul(uint64(wsPeriod+10)) + 8},
			removedSlots: []types.Slot{
				slotsPerEpoch + 4,
				slotsPerEpoch * 2,
				slotsPerEpoch * 3,
				slotsPerEpoch * 4,
				slotsPerEpoch * 5,
			},
		},
		{
			// Prune none.
			slots:       []types.Slot{slotsPerEpoch + 4, slotsPerEpoch*2 + 3, slotsPerEpoch*3 + 4, slotsPerEpoch*4 + 3, slotsPerEpoch*5 + 3},
			storedSlots: []types.Slot{slotsPerEpoch + 4, slotsPerEpoch*2 + 3, slotsPerEpoch*3 + 4, slotsPerEpoch*4 + 3, slotsPerEpoch*5 + 3},
		},
	}
	signedRoot := bytesutil.PadTo([]byte{1}, 32)

	for _, tt := range tests {
		db := setupDB(t, [][48]byte{pubKey})
		for _, slot := range tt.slots {
			err := db.SaveProposalHistoryForSlot(context.Background(), pubKey, slot, signedRoot)
			require.NoError(t, err, "Saving proposal history failed")
		}

		signingRootsBySlot := make(map[types.Slot][]byte)
		proposalHistory, err := db.ProposalHistoryForPubKey(context.Background(), pubKey)
		require.NoError(t, err)

		for _, hist := range proposalHistory {
			signingRootsBySlot[hist.Slot] = hist.SigningRoot
		}

		for _, slot := range tt.removedSlots {
			_, ok := signingRootsBySlot[slot]
			require.Equal(t, false, ok)
		}
		for _, slot := range tt.storedSlots {
			root, ok := signingRootsBySlot[slot]
			require.Equal(t, true, ok)
			require.DeepEqual(t, signedRoot, root, "Unexpected difference in bytes for epoch %d", slot)
		}
		require.NoError(t, db.Close(), "Failed to close database")
	}
}

func TestStore_ProposedPublicKeys(t *testing.T) {
	ctx := context.Background()
	validatorDB, err := NewKVStore(ctx, t.TempDir(), &Config{})
	require.NoError(t, err, "Failed to instantiate DB")
	t.Cleanup(func() {
		require.NoError(t, validatorDB.Close(), "Failed to close database")
		require.NoError(t, validatorDB.ClearDB(), "Failed to clear database")
	})

	keys, err := validatorDB.ProposedPublicKeys(ctx)
	require.NoError(t, err)
	assert.DeepEqual(t, make([][48]byte, 0), keys)

	pubKey := [48]byte{1}
	dummyRoot := [32]byte{}
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubKey, 1, dummyRoot[:])
	require.NoError(t, err)

	keys, err = validatorDB.ProposedPublicKeys(ctx)
	require.NoError(t, err)
	assert.DeepEqual(t, [][48]byte{pubKey}, keys)
}

func TestStore_LowestSignedProposal(t *testing.T) {
	ctx := context.Background()
	pubkey := [48]byte{3}
	dummySigningRoot := [32]byte{}
	validatorDB := setupDB(t, [][48]byte{pubkey})

	_, exists, err := validatorDB.LowestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, false, exists)

	// We save our first proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 2 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the lowest signed slot is what we just saved.
	slot, exists, err := validatorDB.LowestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(2), slot)

	// We save a higher proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 3 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the lowest signed slot did not change.
	slot, exists, err = validatorDB.LowestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(2), slot)

	// We save a lower proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 1 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the lowest signed slot indeed changed.
	slot, exists, err = validatorDB.LowestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(1), slot)
}

func TestStore_HighestSignedProposal(t *testing.T) {
	ctx := context.Background()
	pubkey := [48]byte{3}
	dummySigningRoot := [32]byte{}
	validatorDB := setupDB(t, [][48]byte{pubkey})

	_, exists, err := validatorDB.HighestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, false, exists)

	// We save our first proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 2 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the highest signed slot is what we just saved.
	slot, exists, err := validatorDB.HighestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(2), slot)

	// We save a lower proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 1 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the lowest signed slot did not change.
	slot, exists, err = validatorDB.HighestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(2), slot)

	// We save a higher proposal history.
	err = validatorDB.SaveProposalHistoryForSlot(ctx, pubkey, 3 /* slot */, dummySigningRoot[:])
	require.NoError(t, err)

	// We expect the highest signed slot indeed changed.
	slot, exists, err = validatorDB.HighestSignedProposal(ctx, pubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	assert.Equal(t, types.Slot(3), slot)
}
