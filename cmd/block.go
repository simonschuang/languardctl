// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)
var blockIp string
var blockMac string
var blockvlanId int

func init() {
	blockCmd.Flags().StringVarP(&blockIp, "ip", "i", "", "Block IP address.")
	blockCmd.Flags().StringVarP(&blockMac, "mac", "m", "", "Block MAC address.")
	blockCmd.Flags().IntVarP(&blockvlanId, "vlan", "v", 0, "Vlan ID.")
	RootCmd.AddCommand(blockCmd)
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "set block node",
	Long: `Please specify all range data. For example:
languardctl block -v 10 -i 192.168.99.1 -m 02:42:1b:3c:69:ce
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check inputs
		if blockIp == "" || blockMac == "" {
			cmd.Help()
			return
		}

		db, err := sql.Open("sqlite3", "/tmp/languard.db")
		if err != nil {
			log.Fatal(err)
		}

		// set block node
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert or replace into node(vlan_id, ipv4, mac, block) values(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(blockvlanId, blockIp, blockMac, 1)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		fmt.Println("Done.")
	},
}
