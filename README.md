
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Running it

- Install Go
- Acquire some blockchain .dat files and put them in a subdirectory
- Run `./init.sh` to set up the project (this will call `go get` for the script's dependencies)

Now, assuming you create a directory called `data` inside this repo and place a single .dat file called `blk00689.dat` in it, you can run one of the following commands.

### Viewing basic block data

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 blockdata
```

No file output currently — just logs block info to the console.

### Viewing transaction scripts as strings

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 scripts
```

Script strings will be dumped as .txt files.

### Viewing `OP_RETURN` data

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 opreturns
```

Each time the script finds an `OP_RETURN`, it will create a .dat file containing the raw bytes from the associated data field.

## Output

Output files will be located in the `output` subdirectory (unless you specified an `--outDir` param).

