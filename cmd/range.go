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
	"os"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func init() {
	RootCmd.AddCommand(rangeCmd)
}

var rangeCmd = &cobra.Command{
	Use:   "range",
	Short: "A brief description of scan range",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 10, 2, 2, ' ', 0)

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
		fmt.Fprintln(w, "NAME\tLOCALIP\tNETMASK\tVLANID\tSTART\tEND\tGATEWAY")
		rows, err := db.Query("select name, local_ip, mask, vlan_id, start_ip, end_ip, gateway_ip from range")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			var local_ip string
			var mask string
			var vlan_id int
			var start string
			var end string
			var gateway string
			err = rows.Scan(&name, &local_ip, &mask, &vlan_id, &start, &end, &gateway)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n", name, local_ip, mask, vlan_id, start, end, gateway)
		}
		w.Flush()
	},
}
