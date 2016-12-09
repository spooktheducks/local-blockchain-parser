
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Running it

- Install Go
- Acquire some blockchain .dat files and put them in a subdirectory

Assuming you create a directory called `data` inside this repo and place a single .dat file called `blk00689.dat` in it, you can run the command:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689
```

#### --scripts flag

If you want to print the contents of each transaction's script, run the command with the `--scripts` flag:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 --scripts
```

and the `--outDir` flag:

```sh
$ go run main.go --inDir ./data --startBlock 689 --endBlock 689 --scripts --outDir ./output
```

With this flag, .txt file output will be generated in ./output/scripts (or somewhere else, if you used the `--outDir` flag).

