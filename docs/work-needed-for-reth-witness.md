# Work needed

* reth
  * incorporate block_executor => done
  * add codes to the return => in progress
  * add codes in debug_execution_witness
* geth
  * implement debug_execution_witness
  * add ExecutionWitness type in geth => done
  * sort keys before applying deletions on accounts / nodes => done
* op-program
  * add hints
    * l2-account-proof u64(block_number) ++ address
    * l2-execution-witness ++ block_hash
  * add a client with configurable timeout
  * allow switch between behaviors (read from debug_dbGet or execution_witness and proof)
  * cannon change for input parameters

So it sounds like your ideal path is:
* Chat w/ peter, get commitment around execution standardization on L1 between reth/geth.
* If they'll not support it, dead in water.
* If they will support it, we start by supporting it as a primary but leave the debug_dbGet hints in as a retry.
* Gain confidence over a few weeks on the vm runner + contrived cases in actions, swap over soon after.