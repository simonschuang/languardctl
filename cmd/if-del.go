package cmd

import (
	"fmt"
	"log"
	"strings"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
	"github.com/docker/libcontainer/netlink"
)

var ifname_del string

func init() {
	ifDelCmd.Flags().StringVarP(&ifname_del, "iface", "i", "", "Name of interface")
	RootCmd.AddCommand(ifDelCmd)
}

var ifDelCmd = &cobra.Command{
	Use:   "if-del",
	Short: "Delete an interface",
	Long: `Please specify all range data. For example:
languardctl if-del -i eth0.20
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check inputs
		if ifname_del == "" || !strings.Contains(ifname_del, ".") {
			cmd.Help()
			return
		}
		err := netlink.NetworkLinkDel(ifname_del)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done.")
	},
}
