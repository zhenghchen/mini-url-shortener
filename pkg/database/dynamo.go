package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "urls"

type Store struct {
	client *dynamodb.Client
}

func NewStore(ctx context.Context, endpoint string) (*Store, error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
    )
    if err != nil {
        return nil, err
    }

	var opts []func(*dynamodb.Options)
	if endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	client := dynamodb.NewFromConfig(cfg, opts...)

	return &Store{client: client}, nil

}

func (s *Store) SaveURL(ctx context.Context, code, url string) error {
	item, err := attributevalue.MarshalMap(map[string]string{
		"code": code,
		"url":  url,
	})

	if err != nil {
		return err
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput {
			TableName: aws.String(tableName),
			Item: item,
		})

	return err
}

func (s *Store) GetURL(ctx context.Context, code string) (string, error) {
	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),

		// The key is matchign primary key to map with { "code" : Value}
		// this is needed because dynamo is stored with explicit type information
		// "code": {"S": "abc123"}
		Key: map[string]types.AttributeValue{
			"code": &types.AttributeValueMemberS{Value: code},
		}, 
	})

	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", nil
	}

	var record map[string]string
	if err := attributevalue.UnmarshalMap(result.Item, &record); err != nil {
		return "", err
	}

	return record["url"], nil
	


}
