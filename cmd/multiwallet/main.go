package main

import (
	"github.com/OpenBazaar/multiwallet"
	"github.com/OpenBazaar/multiwallet/config"
	"fmt"
	"sync"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/OpenBazaar/multiwallet/api"
	wi "github.com/OpenBazaar/wallet-interface"
	"os"
	"os/signal"
	"github.com/jessevdk/go-flags"
	"github.com/OpenBazaar/multiwallet/cli"
)

const WALLET_VERSION = "0.1.0"

var parser = flags.NewParser(nil, flags.Default)

type Start struct {
	Testnet            bool   `short:"t" long:"testnet" description:"use the test network"`
}
type Version struct{}

var start Start
var version Version
var mw multiwallet.MultiWallet

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("SPVWallet shutting down...")
			os.Exit(1)
		}
	}()
	parser.AddCommand("start",
		"start the wallet",
		"The start command starts the wallet daemon",
		&start)
	parser.AddCommand("version",
		"print the version number",
		"Print the version number and exit",
		&version)
	cli.SetupCli(parser)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

func (x *Version) Execute(args []string) error {
	fmt.Println(WALLET_VERSION)
	return nil
}

func (x *Start) Execute(args []string) error {
	m := make(map[wi.CoinType]bool)
	m[wi.Bitcoin] = true
	params := &chaincfg.MainNetParams
	if x.Testnet {
		params = &chaincfg.TestNet3Params
	}
	cfg := config.NewDefaultConfig(m, params)
	cfg.Mnemonic = "design author ability expose illegal saddle antique setup pledge wife innocent treat"
	mw, err := multiwallet.NewMultiWallet(cfg)
	if err != nil {
		return err
	}
	go api.ServeAPI(mw)
	var wg sync.WaitGroup
	wg.Add(1)
	mw.Start()
	wg.Wait()
	return nil
}