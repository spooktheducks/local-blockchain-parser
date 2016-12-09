
# local-blockchain-searcher

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Running it

- Install Go
- Acquire some blockchain .dat files and put them in a subdirectory

Assuming you create a directory called `data` inside this repo and place a single .dat file called `blk00689.dat` in it, you can run the command:

```sh
$ go run main.go --infile ./data --startBlock 689 --endBlock 689
```

