package main

// Workaround for mockery support for third party interfaces

import "github.com/hysem/top-word-service/client"

type TopWordServiceClient interface {
	client.Client
}
