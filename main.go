package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ec2conn [profile_name] [region]",
	Short: "Starts an AWS Systems Manager Session Manager session",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		region := args[1]

		awsSession, err := createSession(profileName, region)
		if err != nil {
			fmt.Println("Error creating AWS session:", err)
			os.Exit(1)
		}

		instances, err := getInstanceIDs(awsSession)
		if err != nil {
			fmt.Printf("Error getting instance ID: %v", err)
		}

		if len(instances) == 0 {
			fmt.Println("No running instances found.")
			os.Exit(0)
		}

		selectedInstance, err := decideConnectInstance(instances)
		if err != nil {
			fmt.Printf("Error deciding connect instance: %v", err)
		}

		fmt.Printf("Selected instance: ID: %s, Name: %s\n", selectedInstance.ID, selectedInstance.Name)

		err = startSessionWithCmd(selectedInstance.ID, profileName, region)
		if err != nil {
			fmt.Println("Error starting session:", err)
			os.Exit(1)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
