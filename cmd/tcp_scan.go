/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bcarnazzi/bfg9000/tcp"
	"github.com/spf13/cobra"
)

// tcpScanCmd represents the tcp command
var tcpScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A *brutal* TCP port scanner",
	Long:  `A *brutal* TCP port scanner`,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		tcp.HostsSpec = args[0]
		tcp.PortsSpec, _ = cmd.Flags().GetString("ports")
		tcp.Workers, _ = cmd.Flags().GetInt("worker")
		tcp.Timeout, _ = cmd.Flags().GetInt("timeout")
		tcp.Scanners, _ = cmd.Flags().GetInt("scanner")
		tcp.TcpConnectScan()
	},
}

func init() {
	tcpCmd.AddCommand(tcpScanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tcpScanCmd.Flags().StringP("ports", "p", "1-1024", "Ports to scan")
	tcpScanCmd.Flags().IntP("worker", "w", 200, "Number of workers par host")
	tcpScanCmd.Flags().IntP("scanner", "s", 253, "Number of conccurent host scanner")
	tcpScanCmd.Flags().IntP("timeout", "t", 70, "Timeout in ms")
}
