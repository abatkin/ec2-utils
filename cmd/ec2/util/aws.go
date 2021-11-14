package util

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
	http2 "net/http"
	"net/url"
	"time"
)

type AwsOptions struct {
	Profile string
	Region  string
	Timeout time.Duration
	Proxy   string
}

func BuildAwsClient(awsOptions *AwsOptions) *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(), func(options *config.LoadOptions) error {
		if awsOptions.Profile != "" {
			options.SharedConfigProfile = awsOptions.Profile
		}
		if awsOptions.Region != "" {
			options.Region = awsOptions.Region
		}
		httpClient := http.NewBuildableClient()
		if awsOptions.Timeout > 0 {
			httpClient = httpClient.WithTimeout(awsOptions.Timeout)
		}
		if awsOptions.Proxy != "" {
			if awsOptions.Proxy == "none" {
				log.Printf("Disabling proxy")
				httpClient = httpClient.WithTransportOptions(func(transport *http2.Transport) {
					transport.Proxy = func(request *http2.Request) (*url.URL, error) {
						return nil, nil
					}
				})
			} else {
				log.Printf("Setting proxy to %s", awsOptions.Proxy)
				proxyUrl, err := url.Parse(awsOptions.Proxy)
				if err != nil {
					return fmt.Errorf("bad proxy url %s: %v", awsOptions.Proxy, err)
				}
				httpClient = httpClient.WithTransportOptions(func(transport *http2.Transport) {
					transport.Proxy = http2.ProxyURL(proxyUrl)
				})
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	return client
}
