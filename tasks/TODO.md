# done

* add ExecutionWitness type in geth => done
* sort keys before applying deletions on accounts / nodes => done
* incorporate block_executor => done

# Work needed

* reth
  * add codes to the return => in progress
  * add codes in debug_execution_witness
* geth
  * implement debug_execution_witness
  
* op-program
  * add hints
    * l2-account-proof u64(block_number) ++ address
    * l2-execution-witness ++ block_hash
  * add a client with configurable timeout
  * allow switch between behaviors (read from debug_dbGet or execution_witness and proof)
  * Merge v1.14.10 into op-geth => in progress

1. run base sepolia geth archival node