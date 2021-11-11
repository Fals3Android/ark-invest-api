package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateLocalClient(platform string, port string) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String(fmt.Sprintf("http://docker.for.%s.localhost:%s", platform, port))},
	)
	return sess
}
