/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bcarnazzi/bfg9000/scan"
	"github.com/spf13/cobra"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "TCP port scanner",
	Long:  `TCP port scanner`,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		scan.HostsSpec = args[0]
		scan.PortsSpec, _ = cmd.Flags().GetString("ports")
		scan.Workers, _ = cmd.Flags().GetInt("worker")
		scan.Timeout, _ = cmd.Flags().GetInt("timeout")
		scan.TcpConnectScan()
	},
}

func init() {
	scanCmd.AddCommand(tcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tcpCmd.Flags().StringP("ports", "p", "1-1024", "Ports to scan")
	tcpCmd.Flags().IntP("worker", "w", 100, "Number of workers")
	tcpCmd.Flags().IntP("timeout", "t", 1000, "Timeout in ms")
}
