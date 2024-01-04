/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// feesCmd represents the fees command
var feesCmd = &cobra.Command{
	Use:     "fees",
	Short:   "Print the average fees for transactions at a requested block height",
	Long:    `Print the average fees for transactions at a requested block height`,
	Example: "fees [tendermint rpc endpoint] [block height]",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the tendermint endpoint and block height
		endpoint := args[0]
		block_height, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("could not parse block height")
		}

		fmt.Println("fees called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(feesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// feesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// feesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// feesCmd.Flags().String("endpoint", "http://0.0.0.0:26657", "Tendermint endpoint for the chain")
}
