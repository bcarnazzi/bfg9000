/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tcpCmd represents the scan command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Various TCP tools",
	Long:  `Various TCP tools`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("scan called")
		//fmt.Printf("%-10s %-10s %-10s\n", "PORT", "STATE", "SERVICE")
		//cmd.HasAvailableSubCommands()
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
