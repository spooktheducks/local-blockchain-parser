
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Running it

- Install Go
- Acquire some blockchain .dat files and put them in a subdirectory
- Run `./init.sh` to set up the project (this will call `go get` for the script's dependencies)

Now, assuming you create a directory called `data` inside this repo and place a single .dat file called `blk00689.dat` in it, you can run one of the following commands.

To output basic block data, use the `blockdata` subcommand:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 blockdata
```

To output all transaction scripts as strings, use `scripts`:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 scripts
```

To output the data associated with all `OP_RETURN` ops in the transaction scripts, use `opreturns`:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 opreturns
```


