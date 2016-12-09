package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds"
	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/utils"
)

var (
	flagInDir      = flag.String("inDir", "", "The .dat file containing blockchain input data")
	flagStartBlock = flag.Int64("startBlock", 0, "The block number to start from")
	flagEndBlock   = flag.Int64("endBlock", 0, "The block number to end on")
	flagOutDir     = flag.String("outDir", "output", "Output directory")
)

type (
	Cmd int
)

const (
	CmdBlockData Cmd = iota
	CmdScripts
	CmdOpReturns
)

func getCmd(arg string) Cmd {
	switch arg {
	case "blockdata":
		return CmdBlockData
	case "scripts":
		return CmdScripts
	case "opreturns":
		return CmdOpReturns

	case "":
		fallthrough
	default:
		panic("Must specify a command (blockdata, scripts, or opreturns)")
	}
}

func main() {
	flag.Parse()

	if *flagInDir == "" {
		panic("Missing --inDir param")
	} else if *flagEndBlock == 0 {
		panic("Must specify --endBlock param")
	}

	cmd := getCmd(flag.Arg(0))

	startBlock := uint64(*flagStartBlock)
	endBlock := uint64(*flagEndBlock)

	for i := int(startBlock); i < int(endBlock)+1; i++ {
		filename := fmt.Sprintf("blk%05d.dat", i)

		blocks, err := utils.LoadBlockFile(filepath.Join(*flagInDir, filename))
		if err != nil {
			panic(err)
		}

		for _, bl := range blocks {
			var err error

			switch cmd {
			case CmdBlockData:
				err = cmds.PrintBlockData(bl)
			case CmdScripts:
				err = cmds.PrintBlockScripts(bl, *flagOutDir)
			case CmdOpReturns:
				err = cmds.PrintBlockScriptsOpReturns(bl, *flagOutDir)
			}

			if err != nil {
				panic(err)
			}
		}
	}
}
