
# local-blockchain-parser

Parses blockchain .dat files and spits out the information contained in them.

Faster alternative to using a local RPC implementation.

## Installation

Two options:

### Install a pre-built executable for your platform

- **Windows:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.0/local-blockchain-parser-windowsamd64-0.1.0.exe>
- **Linux:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.0/local-blockchain-parser-linuxamd64-0.1.0>
- **OSX:** <https://github.com/WikiLeaksFreedomForce/local-blockchain-parser/releases/download/0.1.0/local-blockchain-parser-osxamd64-0.1.0>

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

Acquire some blockchain `.dat` files and put them in a subdirectory of this repo.

Assuming you create a directory called `data` inside this repo and place a single .dat file called `blk00689.dat` in it, you can run one of the following commands.

### Viewing basic block data

```sh
$ local-blockchain-parser --inDir ./data --startBlock 689 --endBlock 689 blockdata
```

No file output currently — just logs block info to the console.

### Viewing transaction scripts as strings

```sh
$ local-blockchain-parser --inDir ./data --startBlock 689 --endBlock 689 scripts
```

Script strings will be dumped as .txt files.

### Viewing `OP_RETURN` data

```sh
$ local-blockchain-parser --inDir ./data --startBlock 689 --endBlock 689 opreturns
```

Each time the script finds an `OP_RETURN`, it will create a .dat file containing the raw bytes from the associated data field.

## Output

Output files will be located in the `output` subdirectory (unless you specified an `--outDir` param).

