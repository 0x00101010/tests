package derivation

import (
	"context"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrInvalidVersion = errors.New("invalid version byte")
	ErrInvalidData    = errors.New("invalid data")
)

type Frame struct {
	channelId    []byte
	frameNumber  uint16
	frameDataLen uint32
	frameData    []byte
	isLast       bool
}

type Derivation struct {
	startBlock      uint64
	batchSenderAddr common.Address
	batchInboxAddr  common.Address
	l2GenesisTime   uint64
	l2BlockTime     uint64
	rpcUrl          string
	rpcClient       *ethclient.Client
}

func NewDerivation(
	startBlock uint64,
	batchSenderAddr common.Address,
	batchInboxAddr common.Address,
	l2GenesisTime uint64,
	l2BlockTime uint64,
	rpcUrl string,
) *Derivation {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}

	return &Derivation{
		startBlock:      startBlock,
		batchSenderAddr: batchSenderAddr,
		batchInboxAddr:  batchInboxAddr,
		l2GenesisTime:   l2GenesisTime,
		l2BlockTime:     l2BlockTime,
		rpcUrl:          rpcUrl,
		rpcClient:       client,
	}
}

func (d *Derivation) Derive() {

}

func (d *Derivation) DownloadData(ctx context.Context, blockNum uint64) ([][]byte, error) {
	block, err := d.rpcClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNum))
	if err != nil {
		return nil, err
	}

	data := make([][]byte, 0)
	for _, tx := range block.Transactions() {
		if tx.To() != nil && tx.To().Cmp(d.batchInboxAddr) == 0 {
			// from, err := d.ValidateFromAddress(ctx, tx, block.Hash(), uint(idx))
			// if err != nil {
			// 	return nil, err
			// }

			// if from.Cmp(d.batchSenderAddr) == 0 {
			// 	data = append(data, tx.Data())
			// }
			data = append(data, tx.Data())
		}
	}

	return data, nil
}

func (d *Derivation) ValidateFromAddress(ctx context.Context, tx *types.Transaction, block common.Hash, idx uint) (common.Address, error) {
	from, err := d.rpcClient.TransactionSender(ctx, tx, block, idx)
	if err != nil {
		return common.Address{}, err
	}

	return from, nil
}

func (d *Derivation) Decode(ctx context.Context, data []byte) ([]Frame, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// check version byte
	if data[0] != 0 {
		return nil, ErrInvalidVersion
	}

	offset := uint32(1)
	frames := make([]Frame, 0)

	for offset < uint32(len(data)) {
		channelId := data[offset : offset+16]
		offset += 16
		frameNumber := binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2
		frameDataLen := binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4
		frameData := data[offset : offset+frameDataLen]
		offset += frameDataLen
		isLast := data[offset] != 0
		offset += 1

		frame := Frame{
			channelId:    channelId,
			frameNumber:  frameNumber,
			frameDataLen: frameDataLen,
			frameData:    frameData,
			isLast:       isLast,
		}
		frames = append(frames, frame)
	}

	return frames, nil
}
