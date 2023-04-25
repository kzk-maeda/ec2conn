package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockEC2Client struct {
	ec2iface.EC2API
	Reservations []*ec2.Reservation
}

func (m *mockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{Reservations: m.Reservations}, nil
}

func TestCreateSession(t *testing.T) {
	_, err := createSession("default", "us-west-2")
	assert.Nil(t, err)
}

func TestGetInstanceIDs(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	assert.Nil(t, err)

	mockClient := &mockEC2Client{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("i-1234567890abcdef0"),
						Tags: []*ec2.Tag{
							{Key: aws.String("Name"), Value: aws.String("test-instance")},
						},
					},
				},
			},
		},
	}

	instances, err := getInstanceIDsWithClient(sess, mockClient)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(instances))
	assert.Equal(t, "test-instance", instances[0].Name)
	assert.Equal(t, "i-1234567890abcdef0", instances[0].ID)
}

func getInstanceIDsWithClient(awsSession *session.Session, client ec2iface.EC2API) ([]InstanceInfo, error) {
	// Use provided client instead of creating a new one.
	svc := client

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

// Please note that it is not easy to test functions that rely on external libraries like promptui,
// so for decideConnectInstance and startSessionWithCmd, it's suggested to implement them in such a way that
// they can be tested more easily by separating the business logic from the input/output or using a mock
