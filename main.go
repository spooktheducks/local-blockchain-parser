package main

import (
	"flag"
	"fmt"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds"
)

var (
	flagInDir      = flag.String("inDir", "", "The .dat file containing blockchain input data")
	flagStartBlock = flag.Int64("startBlock", 0, "The block number to start from")
	flagEndBlock   = flag.Int64("endBlock", 0, "The block number to end on")
	flagOutDir     = flag.String("outDir", "output", "Output directory")
)

func main() {
	flag.Parse()

	if *flagInDir == "" {
		panic("Missing --inDir param")
	} else if *flagEndBlock == 0 {
		panic("Must specify --endBlock param")
	}

	cmd := flag.Arg(0)

	startBlock := uint64(*flagStartBlock)
	endBlock := uint64(*flagEndBlock)

	switch cmd {
	case "opreturns":
		err := cmds.PrintBlockScriptsOpReturns(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	case "scripts":
		err := cmds.PrintBlockScripts(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	case "blockdata":
		err := cmds.PrintBlockData(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	case "script-patterns":
		err := cmds.CheckScriptPatterns(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	case "build-blockdb":
		err := cmds.BuildBlockDB(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	case "search-plaintext":
		err := cmds.SearchForPlaintext(startBlock, endBlock, *flagInDir, *flagOutDir)
		if err != nil {
			panic(err)
		}

	default:
		panic("unknown subcommand '" + cmd + "'")
	}

	fmt.Println("DONE")
	var s string
	_, _ = fmt.Scanln(&s)
}
