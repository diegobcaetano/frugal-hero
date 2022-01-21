#Frugal-Hero

This service will help you detect any waste of resources in your AWS account. The policy is: if it is not useful, delete it!

![](https://media4.giphy.com/media/Ti22D4CDvb5fvUtNPC/giphy.gif?cid=790b7611186d72ab4214d8197f00edd1ad4ea6d1ec0b9b3b&rid=giphy.gif&ct=g)

##Requirements

1. You must have Go 1.17+ installed in your machine. Follow [these](https://go.dev/doc/install) instructions to install it.

2. You must have AWS-CLI installed and configured too. Follow [these](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) instructions to install it.

##Installing

Go to the source folder, and type the following command
```
go build -o bin/fh main.go
```

##Running

This program will get the credentials that are your .aws folder to communicate with AWS services.

The following services are available

###S3

This service checks if there is empty buckets in your account

```
bin/fh s3
```
