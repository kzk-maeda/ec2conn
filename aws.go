package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/manifoldco/promptui"
)

type InstanceInfo struct {
	Name string
	ID   string
}

func createSession(profileName, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile:           profileName,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func getInstanceIDs(awsSession *session.Session) ([]InstanceInfo, error) {
	svc := ec2.New(awsSession)

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances := make([]InstanceInfo, 0)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceName := ""
			for _, tag := range instance.Tags {
				if aws.StringValue(tag.Key) == "Name" {
					instanceName = aws.StringValue(tag.Value)
					break
				}
			}
			instances = append(instances, InstanceInfo{Name: instanceName, ID: aws.StringValue(instance.InstanceId)})
		}
	}

	return instances, nil
}

func decideConnectInstance(instances []InstanceInfo) (InstanceInfo, error) {
	prompt := promptui.Select{
		Label: "Instances",
		Items: instances,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Active:   "* {{ .Name | cyan }} ({{ .ID | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .ID | red }})",
			Selected: "Selected: {{ .Name | cyan }}",
		},
		Size: len(instances),
	}

	selectedIndex, _, err := prompt.Run()

	if err != nil {
		return InstanceInfo{}, fmt.Errorf("Prompt failed: %v", err)
	}

	selectedInstance := instances[selectedIndex]

	return selectedInstance, nil
}

func startSessionWithCmd(instanceID, profileName, region string) error {
	cmd := exec.Command("aws", "ssm", "start-session", "--target", instanceID, "--profile", profileName, "--region", region)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create a new signal catcher for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Caught sig: %+v", sig)
			if cmd.Process != nil {
				cmd.Process.Signal(sig)
			}
		}
	}()

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	return nil
}
