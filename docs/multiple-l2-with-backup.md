### Mini design doc

1. Just geth hash archival db
   1. OP_PROGRAM_L2_RPC
2. geth hash archival as backup + new EL
   1. OP_PROGRAM_L2_RPC (geth hash archival db)
   2. OP_PROGRAM_L2_RPC_EXPERIEMENTAL (new EL)
   3. OP_PROGRAM_L2_ENABLE_EXPERIMENTAL_WITH_FALLBACK = true
3. Just the new EL
   1. OP_PROGRAM_L2_RPC


curl -sf --location 'http://localhost:8000' \
--header 'Content-Type: application/json' \
--data '{
	"jsonrpc":"2.0",
	"method":"debug_executionWitness",
	"params":[
        "0x67cfd7",
        false
	],
	"id":1
}'