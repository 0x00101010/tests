package derivation

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDownloadTransaction(t *testing.T) {
	ctx := context.Background()

	derivation := NewDerivation(
		4800000,
		common.HexToAddress("0x6CDEbe940BC0F26850285cacA097C11c33103E47"),
		common.HexToAddress("0xfF00000000000000000000000000000000084532"),
		1695768288,
		2,
		"https://gateway.tenderly.co/public/sepolia",
	)

	fmt.Println(derivation.DownloadData(ctx, 4800003))
}

func TestDerive(t *testing.T) {
	ctx := context.Background()

	derivation := NewDerivation(
		4800000,
		common.HexToAddress("0x6CDEbe940BC0F26850285cacA097C11c33103E47"),
		common.HexToAddress("0xfF00000000000000000000000000000000084532"),
		1695768288,
		2,
		"https://gateway.tenderly.co/public/sepolia",
	)

	datas, err := derivation.DownloadData(ctx, 4800003)
	require.NoError(t, err)

	for _, data := range datas {
		frames, err := derivation.Decode(ctx, data)
		require.NoError(t, err)
		for _, frame := range frames {
			fmt.Println(frame.channelId, frame.frameNumber, frame.frameDataLen, frame.isLast)
		}
	}
}

func TestDeriveMore(t *testing.T) {
	ctx := context.Background()

	derivation := NewDerivation(
		4800000,
		common.HexToAddress("0x6CDEbe940BC0F26850285cacA097C11c33103E47"),
		common.HexToAddress("0xfF00000000000000000000000000000000084532"),
		1695768288,
		2,
		// "https://c3-chainproxy-eth-sepolia-archive.cbhq.net",
		"https://c3-chainproxy-eth-mainnet-archive-dev.cbhq.net",
	)

	// 4800000 -> 4805000
	// 4100000 -> 4800000 => running
	for i := uint64(19005200); i < 20000000; i++ {
		datas, err := derivation.DownloadData(ctx, i)
		require.NoError(t, err)
		if len(datas) > 0 {
			fmt.Println(i, len(datas))
		}
		if i%100 == 0 {
			fmt.Println(i)
		}

		for _, data := range datas {
			frames, err := derivation.Decode(ctx, data)
			require.NoError(t, err)
			for _, frame := range frames {
				if !frame.isLast || frame.frameNumber != 0 {
					fmt.Println(i, frame.frameNumber, frame.isLast)
				}
			}
		}
	}
}
