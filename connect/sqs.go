package connect

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SqsConf TODO
type SqsConf struct {
}

// NewSqsConn new an aws sqs client
func NewSqsConn(_ *SqsConf) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return client, nil
}
