package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"log"
	netHttp "net/http"
	"net/url"
	"time"
)

type Options struct {
	Profile string
	Region  string
	Timeout time.Duration
	Proxy   string
}

func BuildAwsOptions(rootCmd *cobra.Command) *Options {
	awsOptions := &Options{}
	rootCmd.PersistentFlags().StringVar(&awsOptions.Profile, "profile", "", "AWS Profile")
	rootCmd.PersistentFlags().StringVar(&awsOptions.Region, "region", "", "AWS Region")
	rootCmd.PersistentFlags().DurationVar(&awsOptions.Timeout, "timeout", 0, "Timeout")
	rootCmd.PersistentFlags().StringVar(&awsOptions.Proxy, "proxy", "", "Proxy (set to 'none' to force-disable)")
	return awsOptions
}

func (awsOptions *Options) BuildAwsClient() *ec2.Client {
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
	var proxyFunc func(*netHttp.Request) (*url.URL, error)
	if proxy == "none" {
		proxyFunc = forceNoProxy()
	} else {
		if proxyUrl, err := url.Parse(proxy); err != nil {
			return nil, fmt.Errorf("bad proxy url %s: %v", proxy, err)
		} else {
			proxyFunc = netHttp.ProxyURL(proxyUrl)
		}
	}

	httpClient = httpClient.WithTransportOptions(func(transport *netHttp.Transport) {
		transport.Proxy = proxyFunc
	})
	return httpClient, nil
}

func forceNoProxy() func(request *netHttp.Request) (*url.URL, error) {
	return func(request *netHttp.Request) (*url.URL, error) {
		return nil, nil
	}
}

func (awsOptions *Options) BuildRequestContext() context.Context {
	var requestContext = context.Background()
	if awsOptions.Timeout > 0 {
		requestContext, _ = context.WithTimeout(context.Background(), awsOptions.Timeout)
	}
	return requestContext
}
