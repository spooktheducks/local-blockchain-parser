package main

import (
	"fmt"
	"os"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/querycmds"
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
						cmd := querycmds.NewTxInfoCommand(datFileDir, dbFile, txHash)
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
						cmd := querycmds.NewBlockInfoCommand(datFileDir, dbFile, blockHash)
						return cmd.RunCommand()
					},
				},
			},
		},

		{
			Name: "builddb",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "datFileDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "dbFile", Usage: "The database file", Value: "blockchain.db"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, datFileDir, dbFile := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("datFileDir"), c.String("dbFile")
				cmd := querycmds.NewBuildBlockDBCommand(startBlock, endBlock, datFileDir, dbFile)
				return cmd.RunCommand()
			},
		},

		{
			Name: "suspicious-txs",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "inDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, inDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("inDir"), c.String("outDir")
				return cmds.FindSuspiciousTxs(startBlock, endBlock, inDir, outDir)
			},
		},

		{
			Name: "find-plaintext",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "inDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, inDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("inDir"), c.String("outDir")
				return cmds.SearchForPlaintext(startBlock, endBlock, inDir, outDir)
			},
		},

		{
			Name: "find-file-headers",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "inDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, inDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("inDir"), c.String("outDir")
				return cmds.FindFileHeaders(startBlock, endBlock, inDir, outDir)
			},
		},

		{
			Name: "opreturns",
			Flags: []cli.Flag{
				cli.Uint64Flag{Name: "startBlock", Usage: "The block number to start from"},
				cli.Uint64Flag{Name: "endBlock", Usage: "The block number to end on"},
				cli.StringFlag{Name: "inDir", Usage: "The directory containing blockchain blk00XXX.dat files"},
				cli.StringFlag{Name: "outDir", Usage: "The output directory", Value: "output"},
			},
			Action: func(c *cli.Context) error {
				startBlock, endBlock, inDir, outDir := c.Uint64("startBlock"), c.Uint64("endBlock"), c.String("inDir"), c.String("outDir")
				return cmds.PrintOpReturns(startBlock, endBlock, inDir, outDir)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	// flag.Parse()

	// if *flagInDir == "" {
	// 	panic("Missing --inDir param")
	// } else if *flagEndBlock == 0 {
	// 	panic("Must specify --endBlock param")
	// }

	// cmd := flag.Arg(0)

	// startBlock := uint64(*flagStartBlock)
	// endBlock := uint64(*flagEndBlock)

	// switch cmd {
	// case "querydb":
	// 	cmd := cmds.NewQueryBlockDBCommand("output/blockdb", flag.Args())
	// 	err := cmd.RunCommand()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "opreturns":
	// 	err := cmds.PrintOpReturns(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "scripts":
	// 	err := cmds.PrintBlockScripts(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "script-patterns":
	// 	err := cmds.CheckScriptPatterns(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "build-blockdb":
	// 	cmd, err := cmds.NewBuildBlockDBCommand(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	err = cmd.RunCommand()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "suspicious-txs":
	// 	err := cmds.FindSuspiciousTxs(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// case "search-plaintext":
	// 	err := cmds.SearchForPlaintext(startBlock, endBlock, *flagInDir, *flagOutDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// default:
	// 	panic("unknown subcommand '" + cmd + "'")
	// }

}
