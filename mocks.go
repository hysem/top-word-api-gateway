package main

// Workaround for mockery support for third party interfaces

import "github.com/hysem/top-word-service/client"

//go:generate mockery --name=TopWordServiceClient --structname=TopWordServiceClientMock --filename=top_word_service_client_mock.go

type TopWordServiceClient interface {
	client.Client
}
