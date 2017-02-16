
# local-blockchain-parser

Parses blockchain .dat files and spits out various types of information contained in them.

## Installation

Two options:

### Install a pre-built executable for your platform

- **Linux:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.2.1/local-blockchain-parser-linuxamd64>
- **Windows (amd64):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.2.1/local-blockchain-parser-windowsamd64.exe>
- **Windows (386):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.2.1/local-blockchain-parser-windows386.exe>
- **OSX:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.2.1/local-blockchain-parser-osxamd64>

Either rename the executable to `local-blockchain-parser` or use the existing executable name for the commands listed below under "Usage".

### Build/install from source

If you do not have Go installed, here are some instructions:

- Install Go (see <https://golang.org/doc/install> for more information)
    - Windows
        - <https://storage.googleapis.com/golang/go1.7.4.windows-amd64.msi>
    - Linux
        - <https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz>
        - Add to your `~/.profile`: `export GOPATH=$HOME/go`
        - Add to your `~/.profile`: `export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin`
        - `source ~/.profile`
        - `mkdir $HOME/go`
    - OS X / macOS
        - <https://storage.googleapis.com/golang/go1.7.4.darwin-amd64.pkg>
- Run `./init.sh` to install the binary into your `$PATH` so you can use it from any directory.

## Usage

The following example will demonstrate how to decode WikiLeaks' "cablegate" release, which is stored in the blockchain.

### 1. First run

The first time you run `local-blockchain-parser`, it will ask you a few questions:

- the location of your blockchain .dat files
- the location where you want to store the database file created by this program

If you ever need to change these values, they are saved to `~/.wlff-blockchain`, a simple JSON file.

### 2. Acquire some blockchain `.dat` files

- Here are some .dat files to use with the instructions below: <https://mega.nz/#!BkJB3KhI!wuL3Zr_3XNHAgVTiZnWwOLSDz9JnbEkOeULBnlId_JQ>
- You can also use the .dat files downloaded by:
    - [Bitcoin Core](https://bitcoin.org/en/download)
    - [Bitcoin Unlimited](https://www.bitcoinunlimited.info)

### 3. Build the block index

```sh
$ local-blockchain-parser builddb blocks --startBlock 52 --endBlock 52
```

This will index the blocks contained in blk00052.dat (note that you can give the same number for `startBlock` and `endBlock` if you only want to index a single .dat file).

*Note: if you index a .dat file, you can run this command again and it will simply skip indexing.  You can force it to re-index with `--force`.*

Once the index is built, you can ask it about any block contained in the .dat files you indexed:

```sh
$ local-blockchain-parser querydb block-info 000000000000015c28163515610010a24f6469e7741f83a9186393ff25bb8637
```

(where `000000000000015c28163515610010a24f6469e7741f83a9186393ff25bb8637` is the block hash you want to query).

This gives some basic, rudimentary information about the block.

### 4. Build the transaction index

```sh
$ local-blockchain-parser builddb transactions --startBlock 52 --endBlock 52
```

*Note: if you index a .dat file, you can run this command again and it will simply skip indexing.  You can force it to re-index with `--force`.*

Once the index is built, you can ask it about any transaction contained in the .dat files you indexed:

```sh
$ local-blockchain-parser querydb tx-info 5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5
```

(where `5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5` is the transaction hash you want to query).

This command does several things:

- It prints some rudimentary output about the given transaction to the console.
- It creates several output files in `output/tx-chain/<tx hash>` containing the input and output script data for the given transaction.
- It runs a full suite of "transaction checks" that look inside transactions for hidden data.  It currently searches for plaintext, known file headers, PGP keys, Satoshi-encoded data, and a few other things.  If you're querying the transaction given in the example above, the tool should report that it found a `7z header` and "Satoshi data".

### 5. Build the "spent transaction" index

```sh
$ local-blockchain-parser builddb spent-txouts --startBlock 52 --endBlock 52
```

*Note: if you index a .dat file, you can run this command again and it will simply skip indexing.  You can force it to re-index with `--force`.*

This allows us to crawl forward through chains of transactions.

### 6. Decode the Cablegate files from the blockchain

You can do this with only blk00052.dat.  You have to build the block + transaction indices as explained in the examples above.

Then run this command:

```sh
$ local-blockchain-parser querydb tx-chain 5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5 --direction forward --limit 130
```

(`--direction` can be shortened to `-d`, and `--limit` can be shortened to `-l`).

You will notice that the following folder has been created:

```
output/tx-chain/5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5
```

A file called `all-outputs-satoshi-concatenated.dat` is in this folder.  Rename that file to `cablegate.7z` and unzip it with a 7zip extractor.  Tada!  You have the entire Cablegate release.


## Other commands

Note: you can run `local-blockchain-parser --help` (or append `--help` after any command) to see a list of flags and subcommands.

### Dump transaction input/output scripts into files

```sh
$ local-blockchain-parser dump-tx-data --startBlock 645 --endBlock 655
```

Options:

- `--coalesce`: Merge all input scripts for a transaction into a single file, and all output scripts into another file.  Without this flag, every TxIn and TxOut script will be dumped to a separate file.
- `--groupBy <grouptype>`: `grouptype` must be "alpha", "dat", or "blockDate".  "alpha" will group transactions by the first few letters of the transaction hash.  "dat" will group them by .dat file.  "blockDate" will group them by the date the block was accepted into the blockchain.

### Dump transaction fees into a CSV file

```sh
$ local-blockchain-parser dump-tx-fees --startBlock 645 --endBlock 655
```

Note: you should probably run the `builddb transactions` command for the entire blockchain before trying this command.  It has to look up transactions outside of the `--startBlock`/`--endBlock` range to calculate the fees.

### Grep transaction script data for a given hex pattern

This is helpful if you're searching for known file headers or strings inside of transaction scripts.

```sh
$ local-blockchain-parser binary-grep 73706f6f6b7468656475636b73 --block 771 --block 772 --block 773
```

The command above searches for the string "spooktheducks".  You should find it in blk00772.dat.

The output will give you the exact block + transaction + script where the hex pattern was found.

### Searching for known file headers encoded into the blockchain

```sh
$ local-blockchain-parser find-file-headers --startBlock 52 --endBlock 52
```

Searches for a predefined set of file headers (including gzip, 7zip, plaintext PGP packets, JPG, zip, PDF, torrent, etc.) in the specified .dat files.

----

There are other commands.  Use the `--help` flag to find them, or email me at spooktheducks {at} protonmail.com and I'll try to assist.