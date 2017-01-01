package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/dbcmds"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils/aeskeyfind"
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						dbFile := c.String("dbFile")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxInfoCommand(cfg.DatFileDir, dbFile, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "block-info",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
						cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
					},
					Action: func(c *cli.Context) error {
						dbFile, outDir := c.String("dbFile"), c.String("outDir")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxChainCommand(cfg.DatFileDir, dbFile, outDir, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "duplicates",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, "blocks")
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, cfg.DatFileDir, dbFile, "transactions")
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
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
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("dbFile")
						cmd := dbcmds.NewBuildSpentTxOutIndexCommand(startBlock, endBlock, cfg.DatFileDir, dbFile)
						return cmd.RunCommand()
					},
				},
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
}

var configFilename = filepath.Join(os.Getenv("HOME"), ".wlff-blockchain")

func getConfig() (Config, error) {
	bs, err := ioutil.ReadFile(configFilename)
	if err, is := err.(*os.PathError); is {
		return createConfig()
	} else if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not parse config: %v", err)
	}

	return cfg, nil
}

func createConfig() (Config, error) {
	cfg := Config{}

	for cfg.DatFileDir == "" {
		fmt.Printf("Enter the path to the directory containing your blockchain .dat files: ")
		datFileDir, err := scanStr()
		if err != nil {
			return Config{}, err
		}
		cfg.DatFileDir = datFileDir
	}

	f, err := os.Create(configFilename)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	j := json.NewEncoder(f)
	j.SetIndent("", "    ")

	err = j.Encode(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
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
