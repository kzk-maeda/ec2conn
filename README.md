# ec2conn

`ec2conn` is a CLI tool for starting a Session Manager session to an AWS EC2 instance. This tool is written in Go and uses libraries such as Cobra, Viper, and PromptUI.

## Installation

To install `ec2conn`, follow these steps:

1. Install the Go development environment.
2. Clone this repository.
3. Move to the cloned repository directory and run the following command:

```
go install
```


4. After installation is complete, `ec2conn` will be installed in the `$GOPATH/bin` directory. Please set your PATH accordingly.

## Usage

To run `ec2conn`, enter the following command:

```
ec2conn <env> <region>
```


- `<env>` : Specify your AWS profile name.
- `<region>` : Specify the AWS region where your EC2 instances are located.

`ec2conn` will execute the following steps:

1. Log in to AWS and get a list of EC2 instances.
2. Display a list of instances and prompt the user to select one.
3. Start a Session Manager session to the selected instance.

## Options

`ec2conn` has the following options:

- `-h`, `--help` : Display help.

## Example

Here's an example of running `ec2conn`:

```
ec2conn dev ap-northeast-1
```


In this example, it's assumed that your AWS profile name is `dev` and your EC2 instances are located in the `ap-northeast-1` region. `ec2conn` will get a list of instances from the `dev` profile and prompt the user to select one. When the user selects an instance, `ec2conn` will start a Session Manager session to the selected instance.

## Notes

- To run `ec2conn`, you must have AWS CLI configured. For more information on configuring AWS CLI, please refer to the official documentation: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html.
