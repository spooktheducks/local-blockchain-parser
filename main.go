package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"

	"github.com/spooktheducks/local-blockchain-parser/cmds"
	"github.com/spooktheducks/local-blockchain-parser/cmds/dbcmds"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	app := cli.NewApp()

	app.Name = "local blockchain parser"
	app.Commands = []cli.Command{
		{
			Name: "querydb",
			Subcommands: []cli.Command{
				{
					Name: "tx-info",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
					},
					Action: func(c *cli.Context) error {
						dbFile, outDir := c.String("dbFile"), c.String("outDir")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxInfoCommand(cfg.DatFileDir, dbFile, outDir, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "block-info",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
					},
					Action: func(c *cli.Context) error {
						dbFile := c.String("dbFile")
						blockHash := c.Args().Get(0)
						if blockHash == "" {
							return fmt.Errorf("must specify block hash")
						}
						cmd := dbcmds.NewBlockInfoCommand(cfg.DatFileDir, dbFile, blockHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "tx-chain",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
						cli.StringFlag{Name: "direction, d", Usage: "'forward', 'backward', or 'both'", Value: "both"},
						cli.UintFlag{Name: "limit, l", Usage: "Limits the number of transactions crawled", Value: 0},
					},
					Action: func(c *cli.Context) error {
						dbFile, outDir, direction, limit := c.String("dbFile"), c.String("outDir"), c.String("direction"), c.Uint("limit")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxChainCommand(cfg.DatFileDir, dbFile, outDir, direction, limit, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "scan-address",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
					},
					Action: func(c *cli.Context) error {
						dbFile, outDir := c.String("dbFile"), c.String("outDir")
						address := c.Args().Get(0)
						if address == "" {
							return fmt.Errorf("must specify address")
						}
						cmd := dbcmds.NewScanAddressCommand(cfg.DatFileDir, dbFile, outDir, address)
						return cmd.RunCommand()
					},
				},
				{
					Name: "duplicates",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
					},
					Action: func(c *cli.Context) error {
						dbFile := c.String("dbFile")
						cmd := dbcmds.NewScanDupesIndexCommand(cfg.DatFileDir, dbFile)
						return cmd.RunCommand()
					},
				},
			},
		},

		{
			Name: "builddb",
			Subcommands: []cli.Command{
				{
					Name: "blocks",
					Flags: []cli.Flag{
						cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
						cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.BoolFlag{Name: "force, f", Usage: "Force the indexer to re-index blocks that have already been indexed"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile, force := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile"), c.Bool("force")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, "blocks", force)
						if err != nil {
							return err
						}
						return cmd.RunCommand()
					},
				},
				{
					Name: "transactions",
					Flags: []cli.Flag{
						cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
						cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.BoolFlag{Name: "force, f", Usage: "Force the indexer to re-index blocks that have already been indexed"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile, force := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile"), c.Bool("force")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, "transactions", force)
						if err != nil {
							return err
						}
						return cmd.RunCommand()
					},
				},
				{
					Name: "duplicates",
					Flags: []cli.Flag{
						cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
						cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile")
						cmd := dbcmds.NewBuildDupesIndexCommand(startBlock, endBlock, cfg.DatFileDir, dbFile)
						return cmd.RunCommand()
					},
				},
				{
					Name: "spent-txouts",
					Flags: []cli.Flag{
						cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
						cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
						cli.BoolFlag{Name: "force, f", Usage: "Force the indexer to re-index blocks that have already been indexed"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile, force := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile"), c.Bool("force")
						cmd := dbcmds.NewBuildSpentTxOutIndexCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, force)
						return cmd.RunCommand()
					},
				},
			},
		},

		{
			Name: "binary-grep",
			Flags: []cli.Flag{
				// cli.Uint64Flag{Name: "startBlock, s", Usage: "The block number to start from"},
				// cli.Uint64Flag{Name: "endBlock, e", Usage: "The block number to end on"},
				cli.IntSliceFlag{Name: "block, b"},
				cli.StringFlag{Name: "outDir, out", Usage: "The directory where carved files will be saved", Value: "output"},
				cli.Uint64Flag{Name: "carveLen, len", Usage: "The amount of data to carve after each match"},
				cli.StringFlag{Name: "carveExt, ext", Usage: "The extension of the files that are carved", Value: "dat"},
			},
			Action: func(c *cli.Context) error {
				/*startBlock, endBlock,*/ blocks, outDir, carveLen, carveExt := c.IntSlice("block"), c.String("outDir") /*c.Uint64("startBlock"), c.Uint64("endBlock"),*/, c.Uint64("carveLen"), c.String("carveExt")
				hexPattern := c.Args().Get(0)
				if hexPattern == "" {
					return fmt.Errorf("must specify hex pattern to search for")
				}
				cmd := cmds.NewBinaryGrepCommand(blocks, carveLen, carveExt, outDir, cfg.DatFileDir, hexPattern)
				return cmd.RunCommand()
			},
		},

		{
			Name: "graph",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				dbFile, outDir := c.String("dbFile"), c.String("outDir")
				cmd := dbcmds.NewGraphCommand(cfg.DatFileDir, dbFile, outDir, "")
				return cmd.RunCommand()
			},
		},

		{
			Name: "dump-tx-fees",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: cfg.DBFile},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, dbFile, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile"), c.String("outDir")
				cmd := cmds.NewDumpTxFeesCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "dump-tx-data",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
				cli.BoolFlag{Name: "coalesce, c", Usage: "Only output one txin.dat and one txout.dat per transaction"},
				cli.StringFlag{Name: "groupBy, g", Usage: "How to group the output files (if specified, must be either 'alpha' or 'dat')"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir, coalesce, groupBy := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir"), c.Bool("coalesce"), c.String("groupBy")
				cmd, err := cmds.NewDumpTxDataCommand(startBlock, endBlock, cfg.DatFileDir, outDir, coalesce, groupBy)
				if err != nil {
					return err
				}
				return cmd.RunCommand()
			},
		},

		{
			Name: "suspicious-txs",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir")
				return cmds.FindSuspiciousTxs(startBlock, endBlock, cfg.DatFileDir, outDir)
			},
		},

		{
			Name: "find-plaintext",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir")
				cmd := cmds.NewFindPlaintextCommand(startBlock, endBlock, cfg.DatFileDir, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "find-aes-keys",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir")
				cmd := cmds.NewFindAESKeysCommand(startBlock, endBlock, cfg.DatFileDir, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "find-file-headers",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir")
				cmd := cmds.NewFindFileHeadersCommand(startBlock, endBlock, cfg.DatFileDir, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "opreturns",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("outDir")
				return cmds.PrintOpReturns(startBlock, endBlock, cfg.DatFileDir, outDir)
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Config struct {
	DatFileDir string `json:"datFileDir"`
	DBFile     string `json:"dbFile"`
}

var configFilename = filepath.Join(os.Getenv("HOME"), ".wlff-blockchain")

func getConfig() (Config, error) {
	cfg := Config{}

	bs, err := ioutil.ReadFile(configFilename)
	if err, is := err.(*os.PathError); is {
		// no-op
	} else if err != nil {
		return cfg, err
	} else {
		err := json.Unmarshal(bs, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("Could not parse config (%v).  Try deleting the file \"~/.wlff-blockchain\" and running this program again to regenerate it.", err)
		}
	}

	if cfg.DatFileDir == "" {
		cfg.DatFileDir, err = promptInput("Enter the path to the directory containing your blockchain .dat files: ")
		if err != nil {
			return cfg, err
		}
	}

	if cfg.DBFile == "" {
		cfg.DBFile, err = promptInput("Enter the path to your blockchain.db file (or where you want it to go): ")
		if err != nil {
			return cfg, err
		}
	}

	err = saveConfig(cfg)
	if err != nil {
		return cfg, err
	}

	// replace ~ with home dir
	cfg.DatFileDir = strings.Replace(cfg.DatFileDir, "~", os.Getenv("HOME"), 1)
	cfg.DBFile = strings.Replace(cfg.DBFile, "~", os.Getenv("HOME"), 1)

	return cfg, nil
}

func saveConfig(cfg Config) error {
	f, err := os.Create(configFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	j := json.NewEncoder(f)
	j.SetIndent("", "    ")

	err = j.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func promptInput(prompt string) (string, error) {
	var input string
	var err error
	for input == "" {
		fmt.Printf(prompt)
		input, err = scanStr()
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

func scanStr() (string, error) {
	str := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] == '\n' {
			break
		} else {
			str = append(str, b[0])
		}
	}
	return string(str), nil
}
