package connect

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsConf struct {
	Region string
	Prefix string
}

type SqsClient struct {
	client  *sqs.Client
	sqsConf *SqsConf
}

// NewSqsConn new an aws sqs client
func NewSqsConn(c *SqsConf) (*SqsClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(c.Region))
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return &SqsClient{client: client, sqsConf: c}, nil
}

func (s *SqsClient) SqsClient() *sqs.Client {
	return s.client
}

func (s *SqsClient) SqsConf() *SqsConf {
	return s.sqsConf
}
