
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Installation

Two options:

### Install a pre-built executable for your platform

- **Windows (amd64):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.3/local-blockchain-parser-windowsamd64.exe>
- **Windows (386):** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.3/local-blockchain-parser-windows386.exe>
- **Linux:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.3/local-blockchain-parser-linuxamd64>
- **OSX:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.3/local-blockchain-parser-osxamd64>

Either rename the executable to `local-blockchain-parser` or use the existing executable name for the commands listed below under "Usage".

### Build/install from source

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
- Run `./init.sh` to set up the project (this will build and install the binary into your `$PATH`)


## Usage

Acquire some blockchain `.dat` files: <https://mega.nz/#!Y0g3TZxZ!Dgx9bew6hx7gT2s1vE1SRFBjWETOh6HjccC9YL4DH5s>

Assuming you have a single .dat file called `blk00689.dat` in it, you can run one of the following commands:


### Viewing satoshi-downloader encoded data

This is based on the encoding/decoding method from the satoshi python scripts used for cablegate.

```sh
$ local-blockchain-parser --inDir /path/to/data/dir --startBlock 689 --endBlock 689 opreturns
```

Each time the program finds non-`OP_` tokens in the TxOut scripts, it will create a .dat file containing the raw bytes from the associated data field (the data is concatenated across all TxOuts for the given transaction).

As the program runs, it will print to the console when it finds data that matches a known file header or footer (jpeg, pdf, etc.).  For example:

```
- file header match (type: pdf) (block hash: 00000000000000ecbbff6bafb7efa2f7df05b227d5c73dca8f2635af32a2e949) (tx hash: 54e48e5f5c656b26c3bca14a8c95aa583d07ebe84dde3b7dd4a78f4e4186e713)
- file footer match (type: pdf) (block hash: 00000000000000ecbbff6bafb7efa2f7df05b227d5c73dca8f2635af32a2e949) (tx hash: 54e48e5f5c656b26c3bca14a8c95aa583d07ebe84dde3b7dd4a78f4e4186e713)
```

You can verify this file by renaming `./output/op-returns/00000000000000ecbbff6bafb7efa2f7df05b227d5c73dca8f2635af32a2e949/txouts-combined-54e48e5f5c656b26c3bca14a8c95aa583d07ebe84dde3b7dd4a78f4e4186e713.dat` to `blah.pdf` and trying to open it in a regular PDF viewer.

Not everything will be detected by the magic header/footer search.  If you run `./scan-opreturn-data.sh` in this repo after running this command, it will try to identify all valid files among the output using the `file` command.

### Searching for plaintext (instructions?)

```sh
$ local-blockchain-parser --inDir /path/to/data/dir --startBlock 689 --endBlock 689 search-plaintext
```

Output will be generated to `./output/search-plaintext/*.csv`.


### Other subcommands

Viewing basic block data (no file output currently — just logs block info to the console):

```sh
$ local-blockchain-parser --inDir /path/to/data/dir --startBlock 689 --endBlock 689 blockdata
```

Viewing transaction scripts as strings (script strings will be dumped as .txt files):

```sh
$ local-blockchain-parser --inDir /path/to/data/dir --startBlock 689 --endBlock 689 scripts
```


## Output

Output files will be located in the `output` subdirectory (unless you specified an `--outDir` param).

