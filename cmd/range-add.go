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
var name string
var iface string
var localIp string
var netMask string
var vlanId int
var startIp string
var endIp string
var gatewayIp string

func init() {
	rangeAddCmd.Flags().StringVarP(&name, "name", "n", "", "Name to give scan range.")
	rangeAddCmd.Flags().StringVarP(&iface, "iface", "i", "", "Name of interface to scan range.")
	rangeAddCmd.Flags().StringVarP(&localIp, "local", "l", "", "local IP address.")
	rangeAddCmd.Flags().StringVarP(&netMask, "mask", "m", "", "Subnet Mask.")
	rangeAddCmd.Flags().IntVarP(&vlanId, "vlan", "v", -1, "Vlan ID.")
	rangeAddCmd.Flags().StringVarP(&startIp, "start", "s", "", "Start IP range.")
	rangeAddCmd.Flags().StringVarP(&endIp, "end", "e", "", "End IP range.")
	rangeAddCmd.Flags().StringVarP(&gatewayIp, "gateway", "g", "", "Gateway IP.")
	RootCmd.AddCommand(rangeAddCmd)
}

var rangeAddCmd = &cobra.Command{
	Use:   "range-add",
	Short: "Add a scan range",
	Long: `Please specify all range data. For example:
languardctl range-add -n range1 -i eth0 -l 192.168.99.100 -m 24 -v 10 -s 192.168.99.1 -e 192.168.99.253 -g 192.168.99.254
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check inputs
		if name == "" || iface == "" || localIp == "" || netMask =="" || vlanId == -1 || startIp == "" || endIp == "" || gatewayIp == "" {
			cmd.Help()
			return
		}

		db, err := sql.Open("sqlite3", "/tmp/languard.db")
		if err != nil {
			log.Fatal(err)
		}
		sqlStmt := `create table if not exists range(
			name varchar(40) not null primary key,
			iface varchar(40),
			local_ip varchar(40),
			mask varchar(40),
			vlan_id interger default 0,
			start_ip varchar(40),
			end_ip varchar(40),
			gateway_ip varchar(40),
			unique(name) on conflict replace
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		// insert range
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert or replace into range(name, local_ip, mask, vlan_id, start_ip, end_ip, gateway_ip) values(?,?,?,?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(name, localIp, netMask, vlanId, startIp, endIp, gatewayIp)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		fmt.Println("Done.")
	},
}
