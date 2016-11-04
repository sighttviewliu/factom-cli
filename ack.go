// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FactomProject/factom"
)

var status = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli status TxID|FullTx"
	cmd.description = "Returns information about a factoid transaction, or an" +
		" entry / entry credit transaction"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		tx := args[0]

		txID := ""
		fullTx := ""

		_, err := hex.DecodeString(strings.Replace(tx, "\"", "", -1))
		if len(tx) == 64 && err == nil {
			txID = tx
		} else {
			if len(tx) < 64 || err != nil {
				h, err := factom.TransactionHash(tx)
				if err != nil {
					errorln(err)
					return
				}
				txID = h
			} else {
				fullTx = strings.Replace(tx, "\"", "", -1)
			}
		}

		fcack, err1 := factom.FactoidACK(txID, fullTx)
		ecack, err2 := factom.EntryACK(txID, fullTx)
		if err1 != nil && err2 != nil {
			errorln(err1)
			return
		}

		if fcack != nil {
			if fcack.Status != "Unknown" {
				fmt.Println(fcack)
				return
			}
		}

		if ecack != nil {
			if ecack.CommitTxID != "" || ecack.EntryHash != "" {
				fmt.Println(ecack)
				return
			}
		}

		fmt.Printf("Entry / transaction not found.\n")
	}
	help.Add("status", cmd)
	return cmd
}()
