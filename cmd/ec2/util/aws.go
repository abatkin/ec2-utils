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
		var err error

		// AWS Profile
		if awsOptions.Profile != "" {
			options.SharedConfigProfile = awsOptions.Profile
		}

		// AWS Region
		if awsOptions.Region != "" {
			options.Region = awsOptions.Region
		}

		// HTTP Proxy
		httpClient := http.NewBuildableClient()
		if awsOptions.Proxy != "" {
			if httpClient, err = setupProxy(awsOptions.Proxy, httpClient); err != nil {
				return err
			}
		}
		options.HTTPClient = httpClient

		return nil
	})

	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := ec2.NewFromConfig(cfg)
	return client
}

func setupProxy(proxy string, httpClient *http.BuildableClient) (*http.BuildableClient, error) {
	var proxyFunc func(*http2.Request) (*url.URL, error)
	if proxy == "none" {
		proxyFunc = forceNoProxy()
	} else {
		if proxyUrl, err := url.Parse(proxy); err != nil {
			return nil, fmt.Errorf("bad proxy url %s: %v", proxy, err)
		} else {
			proxyFunc = http2.ProxyURL(proxyUrl)
		}
	}

	httpClient = httpClient.WithTransportOptions(func(transport *http2.Transport) {
		transport.Proxy = proxyFunc
	})
	return httpClient, nil
}

func forceNoProxy() func(request *http2.Request) (*url.URL, error) {
	return func(request *http2.Request) (*url.URL, error) {
		return nil, nil
	}
}
