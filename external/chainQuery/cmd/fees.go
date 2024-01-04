/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/SDavidson1177/chainQuery/types"
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
		block_height := args[1]

		// Ensure there is a trailing slash for the endpoint
		if endpoint[len(endpoint)-1] != '/' {
			endpoint = endpoint + "/"
		}

		block_height = "tx.height=" + block_height
		query_endpoint := fmt.Sprintf("%s%s\"%v\"", endpoint, "tx_search?query=", url.QueryEscape(block_height))

		resp, err := http.Get(query_endpoint)
		if err != nil {
			return err
		}

		const SIZE = 100

		resp_body_bytes := make([]byte, SIZE)
		n, err := resp.Body.Read(resp_body_bytes)
		for n == SIZE {
			if err != nil {
				return err
			}

			tmp := make([]byte, SIZE)
			n, err = resp.Body.Read(tmp)
			resp_body_bytes = append(resp_body_bytes, tmp...)
		}

		// truncate trailing empty bytes
		resp_body_bytes = resp_body_bytes[:len(resp_body_bytes)-(100-n)]

		// Received the response as bytes (Array of transactions)
		// Now, extract the gas used
		var result types.FeesResponse
		if err := json.Unmarshal(resp_body_bytes, &result); err != nil {
			return err
		}

		// Take the average "gas used" for transactions in this block
		// TODO: determine if this is reasonable for determining gas prices

		total_gas := 0
		avg_gas := 0
		len := 0

		for _, tx := range result.Result.Txs {
			val, _ := strconv.Atoi(tx.TxResult.GasUsed)
			total_gas += val
			len++
		}

		if len > 0 {
			avg_gas = total_gas / len
		}

		fmt.Printf("Total Gas Used: %v\nAverage Gas Used: %v\nNumber of Txs: %v\n", total_gas, avg_gas, len)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(feesCmd)
}
