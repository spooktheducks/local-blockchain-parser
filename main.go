package main

import (
	"fmt"
	"os"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/dbcmds"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "local blockchain parser"
	app.Commands = []cli.Command{
		{
			Name: "querydb",
			Subcommands: []cli.Command{
				{
					Name: "tx-info",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						datFileDir, dbFile := c.String("datFileDir"), c.String("dbFile")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxInfoCommand(datFileDir, dbFile, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "block-info",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						datFileDir, dbFile := c.String("datFileDir"), c.String("dbFile")
						blockHash := c.Args().Get(0)
						if blockHash == "" {
							return fmt.Errorf("must specify block hash")
						}
						cmd := dbcmds.NewBlockInfoCommand(datFileDir, dbFile, blockHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "tx-chain",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
						cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
					},
					Action: func(c *cli.Context) error {
						datFileDir, dbFile, outDir := c.String("datFileDir"), c.String("dbFile"), c.String("outDir")
						txHash := c.Args().Get(0)
						if txHash == "" {
							return fmt.Errorf("must specify tx hash")
						}
						cmd := dbcmds.NewTxChainCommand(datFileDir, dbFile, outDir, txHash)
						return cmd.RunCommand()
					},
				},
				{
					Name: "duplicates",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						datFileDir, dbFile := c.String("datFileDir"), c.String("dbFile")
						cmd := dbcmds.NewScanDupesIndexCommand(datFileDir, dbFile)
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
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, datFileDir, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("dbFile")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, datFileDir, dbFile, "blocks")
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
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, datFileDir, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("dbFile")
						cmd, err := dbcmds.NewBuildBlockDBCommand(startBlock, endBlock, datFileDir, dbFile, "transactions")
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
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, datFileDir, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("dbFile")
						cmd := dbcmds.NewBuildDupesIndexCommand(startBlock, endBlock, datFileDir, dbFile)
						return cmd.RunCommand()
					},
				},
				{
					Name: "spent-txouts",
					Flags: []cli.Flag{
						cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
						cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
						cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
						cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
					},
					Action: func(c *cli.Context) error {
						startBlock, endBlock, datFileDir, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("dbFile")
						cmd := dbcmds.NewBuildSpentTxOutIndexCommand(startBlock, endBlock, datFileDir, dbFile)
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
				cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, datFileDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("outDir")
				return cmds.FindSuspiciousTxs(startBlock, endBlock, datFileDir, outDir)
			},
		},

		{
			Name: "find-plaintext",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, datFileDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("outDir")
				cmd := cmds.NewFindPlaintextCommand(startBlock, endBlock, datFileDir, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "find-file-headers",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, datFileDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("outDir")
				cmd := cmds.NewFindFileHeadersCommand(startBlock, endBlock, datFileDir, outDir)
				return cmd.RunCommand()
			},
		},

		{
			Name: "opreturns",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, datFileDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("outDir")
				return cmds.PrintOpReturns(startBlock, endBlock, datFileDir, outDir)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}
