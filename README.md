# Frugal-Hero

This service will help you detect resources created in your AWS account that are not being used. 
Currently, the only output is a list of resources in the console. Future versions will support different types of output, 
making possible to integrate the output with automation tools. 

![](https://media4.giphy.com/media/Ti22D4CDvb5fvUtNPC/giphy.gif?cid=790b7611186d72ab4214d8197f00edd1ad4ea6d1ec0b9b3b&rid=giphy.gif&ct=g)

## Requirements

1. You must have Go 1.17+ installed in your machine. Follow [these](https://go.dev/doc/install) instructions to install it.

2. You must have AWS-CLI installed and configured too. Follow [these](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) instructions to install it.

## Installing

Go to the source folder, and type the following command
```
go build -o fh
```

## Running

This program will get the credentials that are your .aws folder to communicate with AWS services.

The following services are available

### S3

This service checks if there is empty buckets in your account (it will only look for buckets that are in the same region configured in your AWS-CLI)

```
./fh s3
```

### Lambda

This service returns a list of all functions that are not invoked in the past 7 days.

```
./fh lambda
```
