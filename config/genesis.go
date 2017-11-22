package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/bytom/protocol/bc/legacy"
	"github.com/bytom/protocol/bc"
	"github.com/bytom/consensus"
	"github.com/bytom/protocol/state"
)

/*
initial_block_hex := (
{
	{1 // version
	 1 // block height
	 {0 0 0 0} // prev block hash
	 1508832880206 // timestamp
	{
			{7711099753227061847 10258754681536048093 9974865291001675724 174641040502263804} //  tx merkle root
			{5269591548479411549 11645895074199058355 6636635844902830556 14344152548198350030} // asset merkle root
	}
	417239 // nonce
	2161727821138738707 // bits
	} // block header
}
)*/

func GennerateGenesisBlock() *legacy.Block {
//want_initial_block_hex := "0301010000000000000000000000000000000000000000000000000000000000000000cecccaebf42b406b03545ed2b38a578e5e6b0796d4ebdd8a6dd72210873fcc026c7319de578ffc492159980684155da19e87de0d1b37b35c1a1123770ec1dcc710aabe77607cced7bb1993fcb680808080801e0107010700cecccaebf42b000001012cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8080ccdee2a69fb314010151000000"

	txData := legacy.TxData{
		Version: 1,
		SerializedSize: 60,
		Inputs: []*legacy.TxInput{},
		Outputs:[]*legacy.TxOutput{
			&legacy.TxOutput{
				AssetVersion: 1,
				OutputCommitment: legacy.OutputCommitment{
					AssetAmount: bc.AssetAmount{
						AssetId: consensus.BTMAssetID,
						Amount:  1470000000000000000,
					},
					VMVersion:      1,
					ControlProgram: []byte{81},
				},
			},
		},
		MinTime: 0,
		MaxTime: 1508832880206,
	}

	genesisCoinbaseTx := legacy.NewTx(txData)

    merkleRoot, err := bc.MerkleRoot([]*bc.Tx{genesisCoinbaseTx.Tx})
	if err != nil {
		log.Errorf("Fatal create merkelRoot")
	}
	snap := state.Empty()
	if err := snap.ApplyTx(genesisCoinbaseTx.Tx); err != nil {
		log.Errorf("Fatal ApplyTx")
	}

	genesisBlock := &legacy.Block{
		BlockHeader:  legacy.BlockHeader{
			Version: 1,
			Height: 1,
			TimestampMS: 1508832880206,
			BlockCommitment: legacy.BlockCommitment{
				TransactionsMerkleRoot: merkleRoot,
				AssetsMerkleRoot:       snap.Tree.RootHash(),
			},
			Bits: 2161727821138738707,
		},
		Transactions: []*legacy.Tx{genesisCoinbaseTx},
	}
	for i := uint64(0); i <= 10000000000000; i++ {
		genesisBlock.Nonce = i
		hash := genesisBlock.Hash()

		if consensus.CheckProofOfWork(&hash, genesisBlock.Bits) {
			break
		}
	}

	log.Infof("genesisBlock:%v", genesisBlock)
	return genesisBlock
}
