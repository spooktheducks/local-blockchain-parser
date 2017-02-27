# TODO

- [ ] Add CSV dump to `dump-tx-data` without `--coalesce`
- [ ] Remove spaces from `dump-tx-data` output folders
- [ ] Flag to avoid API calls
- [ ] Better command line help
- [ ] Make `tx-chain` subcommand able to use different, pluggable algorithms for detecting a valid "next transaction"
- [ ] Create a `RunFullSuite(tx *btcutil.Tx)` function that implements all known checks on a given transaction (see `cmds/utils/extract-data.go` for the current set of checks).  This function should output a `struct` representing the "scores" for a given Tx (on a scale of not-suspicious to very-suspicious)
- [ ] Improve plaintext detection to filter more irrelevant data (see `cmds/cmd-find-plaintext.go` and `cmds/utils/extract-data.go`)
- [ ] Improve PGP/GPG data checker (see `cmd-txinfo.go`, should be abstracted out into `cmds/utils/extract-data.go`)
- [ ] Add a transaction input/output checker for hex data matching known WikiLeaks file hashes and other known hex strings


## done

- [x] Implement forward crawling in `cmds/cmd-tx-chain.go`
- [x] Re-architect the code so that we can run any set of checks on any "transaction source" and "data source".  Use interfaces.
    - Transaction sources:
        - Single transaction
        - Transaction chain (see `cmds/cmd-txchain.go`)
        - Full scan of a given .dat file (or set of .dat files)
    - Input sources:
        - TxIn scripts raw byte data
        - TxIn scripts raw byte data, "Satoshi encoded" (i.e., with length+checksum prefix)
        - TxIn scripts interpreted as hex (or scraped for valid hex values)
        - TxOut scripts raw byte data
        - TxOut scripts raw byte data, "Satoshi encoded" (i.e., with length+checksum prefix)
        - TxOut scripts interpreted as hex (or scraped for valid hex values)
- [x] Rewrite readme
    - Mention config file
    - Mention that you should do full tx index (and mention that after indexing, it won't reindex, so it's safe for scripts that you run over and over)
    - Update cablegate example - mention tx-chain `-l` and `-d` flags