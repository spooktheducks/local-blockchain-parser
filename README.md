
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Installation

Two options:

### Install a pre-built executable for your platform

- **Windows (amd64):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.5/local-blockchain-parser-windowsamd64.exe>
- **Windows (386):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.5/local-blockchain-parser-windows386.exe>
- **Linux:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.5/local-blockchain-parser-linuxamd64>
- **OSX:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.5/local-blockchain-parser-osxamd64>

Either rename the executable to `local-blockchain-parser` or use the existing executable name for the commands listed below under "Usage".

### Build/install from source

If you already have Go and `git` installed, just run the following command in the terminal:

```sh
curl https://raw.githubusercontent.com/WikiLeaksFreedomForce/local-blockchain-parser/master/setup.sh  | bash
```

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
- Run `./init.sh` to set up the `local-blockchain-parser` project (this will build and install the binary into your `$PATH`).


## Usage

Acquire some blockchain `.dat` files:

- Here are some .dat files to use with the instructions below: <https://mega.nz/#!BkJB3KhI!wuL3Zr_3XNHAgVTiZnWwOLSDz9JnbEkOeULBnlId_JQ>
- You can also use the .dat files downloaded by [Bitcoin Core](https://bitcoin.org/en/download).

For the following examples, we assume that you put the .dat files in a folder called `data` in this repo.  Note, however, that you can provide a `--datFileDir` argument to any command, specifying a different location.

Assuming you have a single .dat file called `blk00052.dat` in your `data` folder, you can run one of the following commands:

### 1. Build the block index

```sh
local-blockchain-parser builddb blocks --datFileDir ./data --startBlock 52 --endBlock 53
```

This will index the blocks contained in blk00052.dat and blk00053.dat (note that you can give the same number for `startBlock` and `endBlock` if you only want to index a single .dat file).

Once the index is built, you can ask it about any block contained in the .dat files you indexed:

```sh
local-blockchain-parser querydb block-info --datFileDir ./data 000000000000015c28163515610010a24f6469e7741f83a9186393ff25bb8637
```

(where `000000000000015c28163515610010a24f6469e7741f83a9186393ff25bb8637` is the block hash you want to query).

This gives some basic, rudimentary information about the block.

### 2. Build the transaction index

```sh
local-blockchain-parser builddb transactions --datFileDir ./data --startBlock 52 --endBlock 53
```

Once the index is built, you can ask it about any transaction contained in the .dat files you indexed:

```sh
local-blockchain-parser querydb tx-info --datFileDir ./data 5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5
```

(where `5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5` is the transaction hash you want to query).

This will run the full suite of "transaction checks" that look inside transactions for hidden data.  It currently searches for plaintext, known file headers, PGP keys, Satoshi-encoded data.  If you're querying the transaction given in the example above, the tool should report that it found a `7z header`.

### 3. Decode the Cablegate files from the blockchain

You can do this with only blk00052.dat.  You have to build the block + transaction indices as explained in the examples above.

```sh
local-blockchain-parser querydb tx-chain --datFileDir ./data 2c9e766020d9e93bea3a1d149313ab224d3c375ad9341594331fa9c48bce13b8
```

You will notice that an `output` folder has been created.  A file called `txchain-output` is in this folder.  Rename that file to `cablegate.7z` and unzip it with a 7zip extractor.  Tada, you have the entire cablegate release.

The transaction hash you have to specify for this command can be any transaction in the cablegate release.  The tool will crawl forwards and backwards through the relevant transactions and combine all of the data it finds.


## Other commands

### Searching for plaintext encoded into the blockchain

```sh
$ local-blockchain-parser find-plaintext --datFileDir ./data --startBlock 52 --endBlock 52
```

This command generates a lot of false positives (and therefore, a LOT of output).  We're working on improving it.

Output will be generated to `./output/find-plaintext/*.csv`.


### Searching for known file headers encoded into the blockchain

```sh
$ local-blockchain-parser find-file-headers --datFileDir ./data --startBlock 52 --endBlock 52
```

Output will be generated to `./output/find-file-headers/*.txt`.

