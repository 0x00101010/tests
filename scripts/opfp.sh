./target/debug/opfp \
    from-op-program \
    --op-program ~/src/optimism/op-program/bin/op-program \
    --l2-block 15192787 \
    --l1-block 6678352 \
    --l1-rpc-url http://44.192.45.93:8000 \
    --l2-rpc-url https://base-sepolia-dev.cbhq.net:8545 \
    --beacon-url http://44.192.45.93:8010 \
    --rollup-url https://base-sepolia-archive-k8s-dev.cbhq.net:7545 \
    --chain-name base-sepolia \
    --output base-sepolia.out
