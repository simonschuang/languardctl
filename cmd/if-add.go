package cmd

import (
	"fmt"
	"net"
	"log"
	"github.com/spf13/cobra"
	_ "github.com/mattn/go-sqlite3"
	"github.com/docker/libcontainer/netlink"
)

var ifname string
var cidr string
var ifVlanId int
var ifGatewayIp string

func init() {
	ifAddCmd.Flags().StringVarP(&ifname, "iface", "i", "", "Name of interface")
	ifAddCmd.Flags().StringVarP(&cidr, "cidr", "c", "", "local CIDR address.")
	ifAddCmd.Flags().IntVarP(&ifVlanId, "vlan", "v", -1, "Vlan ID.")
	ifAddCmd.Flags().StringVarP(&ifGatewayIp, "gateway", "g", "", "Gateway IP.")
	RootCmd.AddCommand(ifAddCmd)
}

var ifAddCmd = &cobra.Command{
	Use:   "if-add",
	Short: "Add an interface",
	Long: `Please specify all range data. For example:
languardctl if-add -i eth0 -c 192.168.20.100/24 -v 20 -g 192.168.20.254
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check inputs
		if ifname == "" || cidr == "" || ifVlanId == -1 || ifVlanId > 4096 || ifGatewayIp == "" {
			cmd.Help()
			return
		}
		var name string
		var Id uint16
		Id = uint16(ifVlanId)
		name = fmt.Sprintf("%s.%d", ifname, ifVlanId)
		err := netlink.NetworkLinkAddVlan(ifname, name, Id)
		if err != nil {
			log.Fatal(err)
		}
		iface, err := net.InterfaceByName(name)
		if err != nil {
			log.Fatal(err)
		}
		err = netlink.NetworkLinkUp(iface)
		if err != nil {
			log.Fatal(err)
		}
		IP, NET, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatal(err)
		}
		err = netlink.NetworkLinkAddIp(iface, IP, NET)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done.")
	},
}
