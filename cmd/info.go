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
	"net"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func init() {
	RootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of languard",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var ipv4 string
		var machine_uuid string
		var serial_number string
		w := tabwriter.NewWriter(os.Stdout, 10, 2, 2, ' ', 0)

		db, err := sql.Open("sqlite3", "/tmp/languard.db")
		if err != nil {
			log.Fatal(err)
		}
		sqlStmt := `create table if not exists general(
			id integer default 0,
			name varchar(40),
			ipv4 varchar(40),
			machine_uuid varchar(200),
			serial_number varchar(200),
			unique(id) on conflict replace

		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}
		sqlStmt = `create table if not exists interface(
			if_name varchar(40),
			ipv4 varchar(40),
			mask varchar(40),
			ipv6 varchar(40),
			mac varchar(40),
			unique(if_name) on conflict replace
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}
		rows, err := db.Query("select name, ipv4, machine_uuid, serial_number from general")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&name, &ipv4, &machine_uuid, &serial_number)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprintf(w, "Name\t%s\n", name)
		fmt.Fprintf(w, "ipv4\t%s\n", ipv4)
		fmt.Fprintf(w, "machine_uuid\t%s\n", machine_uuid)
		fmt.Fprintf(w, "serial_number\t%s\n", serial_number)
		w.Flush()

		if err := scanInterface(db); err != nil {
			log.Printf("%v", err)
		}
	},
}

func scanInterface(db *sql.DB) error {
	// Get a list of all interfaces.
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 2, ' ', 0)
	fmt.Fprintln(w, "INTERFACE\tIPV4\tMASK\tMAC")
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			var mask net.IPMask
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
                        // ip = ip.To16() This can both retrive ipv6 & ipv4
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", iface.Name, ip.String(), mask.String(), iface.HardwareAddr.String())
		}
	}
	w.Flush()
	return nil
}
