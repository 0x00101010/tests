docker exec opt_conductor_1 sh -c "
curl \
    --location http://localhost:8545 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"conductor_leader",
    "params":[
    ],
    "id": 1
}' | jq"

docker exec opt_conductor_1 sh -c '
curl \
    -sf \
    --location http://localhost:8545 \
    --header "Content-Type: application/json" \
    --data '"'"'{
        "jsonrpc":"2.0",
        "method":"conductor_clusterMembership",
        "params":[],
        "id":1
    }'"'"' | jq'

docker exec opt_conductor_1 sh -c '
curl \
    -sf \
    --location http://localhost:8545 \
    --header "Content-Type: application/json" \
    --data '"'"'{
        "jsonrpc":"2.0",
        "method":"conductor_clusterMembership",
        "params":[
            "{argument name="block hash"}"
        ],
        "id":1
    }'"'"' | jq'

curl \
    --location http://localhost:7545 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"optimism_syncStatus",
    "params":[
    ],
    "id": 1
}' |
    jq

curl \
    --location http://localhost:8000 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"eth_getBlockByNumber",
    "params":[
        "latest",
        false
    ],
    "id": 1
}' |
    jq

curl \
    --location http://localhost:8000 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"debug_executionWitness",
    "params":[
        "0x65e705",
        true
    ],
    "id": 1
}' |
    jq >reth-witness-0x65e705-after.json

curl \
    --location http://localhost:6000 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"eth_getBlockByNumber",
    "params":[
        "latest",
        false
    ],
    "id": 1
}' |
    jq .result.number

curl \
    --location http://localhost:6000 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"debug_executionWitness",
    "params":[
        "0xE8AD6F",
        true
    ],
    "id": 1
}' |
    jq >op-reth-witness-0xE8AD6F-after.json

curl \
    --location https://base-sepolia-archive-k8s-dev.cbhq.net:7545 \
    --header 'Content-Type: application/json' \
    --data '{
    "jsonrpc":"2.0",
    "method":"optimism_outputAtBlock",
    "params":[
        "0xE7D2D3"
    ],
    "id": 1
}' |
    jq
