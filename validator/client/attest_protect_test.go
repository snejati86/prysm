package client

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/config/features"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func Test_slashableAttestationCheck(t *testing.T) {
	config := &features.Flags{
		RemoteSlasherProtection: true,
	}
	reset := features.InitWithReset(config)
	defer reset()
	validator, m, validatorKey, finish := setup(t)
	defer finish()
	pubKey := [48]byte{}
	copy(pubKey[:], validatorKey.PublicKey().Marshal())
	att := &ethpb.IndexedAttestation{
		AttestingIndices: []uint64{1, 2},
		Data: &ethpb.AttestationData{
			Slot:            5,
			CommitteeIndex:  2,
			BeaconBlockRoot: bytesutil.PadTo([]byte("great block"), 32),
			Source: &ethpb.Checkpoint{
				Epoch: 4,
				Root:  bytesutil.PadTo([]byte("good source"), 32),
			},
			Target: &ethpb.Checkpoint{
				Epoch: 10,
				Root:  bytesutil.PadTo([]byte("good target"), 32),
			},
		},
	}

	m.slasherClient.EXPECT().IsSlashableAttestation(
		gomock.Any(), // ctx
		att,
	).Return(&ethpb.AttesterSlashingResponse{AttesterSlashings: []*ethpb.AttesterSlashing{{
		Attestation_1: &ethpb.IndexedAttestation{},
		Attestation_2: &ethpb.IndexedAttestation{},
	}}}, nil /*err*/)

	err := validator.slashableAttestationCheck(context.Background(), att, pubKey, [32]byte{1})
	require.ErrorContains(t, failedPostAttSignExternalErr, err)

	m.slasherClient.EXPECT().IsSlashableAttestation(
		gomock.Any(), // ctx
		att,
	).Return(&ethpb.AttesterSlashingResponse{}, nil /*err*/)

	err = validator.slashableAttestationCheck(context.Background(), att, pubKey, [32]byte{1})
	require.NoError(t, err, "Expected allowed attestation not to throw error")
}

func Test_slashableAttestationCheck_UpdatesLowestSignedEpochs(t *testing.T) {
	config := &features.Flags{
		RemoteSlasherProtection: true,
	}
	reset := features.InitWithReset(config)
	defer reset()
	validator, m, validatorKey, finish := setup(t)
	defer finish()
	ctx := context.Background()
	pubKey := [48]byte{}
	copy(pubKey[:], validatorKey.PublicKey().Marshal())
	att := &ethpb.IndexedAttestation{
		AttestingIndices: []uint64{1, 2},
		Data: &ethpb.AttestationData{
			Slot:            5,
			CommitteeIndex:  2,
			BeaconBlockRoot: bytesutil.PadTo([]byte("great block"), 32),
			Source: &ethpb.Checkpoint{
				Epoch: 4,
				Root:  bytesutil.PadTo([]byte("good source"), 32),
			},
			Target: &ethpb.Checkpoint{
				Epoch: 10,
				Root:  bytesutil.PadTo([]byte("good target"), 32),
			},
		},
	}

	m.slasherClient.EXPECT().IsSlashableAttestation(
		gomock.Any(), // ctx
		att,
	).Return(&ethpb.AttesterSlashingResponse{}, nil /*err*/)

	m.validatorClient.EXPECT().DomainData(
		gomock.Any(), // ctx
		&ethpb.DomainRequest{Epoch: 10, Domain: []byte{1, 0, 0, 0}},
	).Return(&ethpb.DomainResponse{SignatureDomain: make([]byte, 32)}, nil /*err*/)
	_, sr, err := validator.getDomainAndSigningRoot(ctx, att.Data)
	require.NoError(t, err)

	err = validator.slashableAttestationCheck(context.Background(), att, pubKey, sr)
	require.NoError(t, err)
	differentSigningRoot := [32]byte{2}

	err = validator.slashableAttestationCheck(context.Background(), att, pubKey, differentSigningRoot)
	require.ErrorContains(t, "could not sign attestation", err)

	e, exists, err := validator.db.LowestSignedSourceEpoch(context.Background(), pubKey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	require.Equal(t, types.Epoch(4), e)
	e, exists, err = validator.db.LowestSignedTargetEpoch(context.Background(), pubKey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	require.Equal(t, types.Epoch(10), e)
}

func Test_slashableAttestationCheck_OK(t *testing.T) {
	config := &features.Flags{
		RemoteSlasherProtection: true,
	}
	reset := features.InitWithReset(config)
	defer reset()
	ctx := context.Background()
	validator, mocks, _, finish := setup(t)
	defer finish()
	att := &ethpb.IndexedAttestation{
		AttestingIndices: []uint64{1, 2},
		Data: &ethpb.AttestationData{
			Slot:            5,
			CommitteeIndex:  2,
			BeaconBlockRoot: []byte("great block"),
			Source: &ethpb.Checkpoint{
				Epoch: 4,
				Root:  []byte("good source"),
			},
			Target: &ethpb.Checkpoint{
				Epoch: 10,
				Root:  []byte("good target"),
			},
		},
	}
	sr := [32]byte{1}
	fakePubkey := bytesutil.ToBytes48([]byte("test"))

	mocks.slasherClient.EXPECT().IsSlashableAttestation(
		gomock.Any(), // ctx
		att,
	).Return(&ethpb.AttesterSlashingResponse{}, nil /*err*/)

	err := validator.slashableAttestationCheck(ctx, att, fakePubkey, sr)
	require.NoError(t, err, "Expected allowed attestation not to throw error")
}

func Test_slashableAttestationCheck_GenesisEpoch(t *testing.T) {
	config := &features.Flags{
		RemoteSlasherProtection: true,
	}
	reset := features.InitWithReset(config)
	defer reset()
	ctx := context.Background()
	validator, mocks, _, finish := setup(t)
	defer finish()
	att := &ethpb.IndexedAttestation{
		AttestingIndices: []uint64{1, 2},
		Data: &ethpb.AttestationData{
			Slot:            5,
			CommitteeIndex:  2,
			BeaconBlockRoot: bytesutil.PadTo([]byte("great block root"), 32),
			Source: &ethpb.Checkpoint{
				Epoch: 0,
				Root:  bytesutil.PadTo([]byte("great root"), 32),
			},
			Target: &ethpb.Checkpoint{
				Epoch: 0,
				Root:  bytesutil.PadTo([]byte("great root"), 32),
			},
		},
	}

	mocks.slasherClient.EXPECT().IsSlashableAttestation(
		gomock.Any(), // ctx
		att,
	).Return(&ethpb.AttesterSlashingResponse{}, nil /*err*/)

	fakePubkey := bytesutil.ToBytes48([]byte("test"))
	err := validator.slashableAttestationCheck(ctx, att, fakePubkey, [32]byte{})
	require.NoError(t, err, "Expected allowed attestation not to throw error")
	e, exists, err := validator.db.LowestSignedSourceEpoch(context.Background(), fakePubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	require.Equal(t, types.Epoch(0), e)
	e, exists, err = validator.db.LowestSignedTargetEpoch(context.Background(), fakePubkey)
	require.NoError(t, err)
	require.Equal(t, true, exists)
	require.Equal(t, types.Epoch(0), e)
}
