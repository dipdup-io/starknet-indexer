package api

import (
	"context"
	"fmt"
	"github.com/dipdup-net/go-lib/config"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
	"strconv"
)

type Subsquid struct {
	httpClient fastshot.ClientHttpMethods
}

func NewSubsquid(cfg config.DataSource) *Subsquid {
	var httpClient = fastshot.NewClient(cfg.URL).
		Build()

	return &Subsquid{
		httpClient: httpClient,
	}
}

func (s *Subsquid) GetWorkerUrl(_ context.Context, startLevel uint64) (string, error) {
	path := fmt.Sprintf("/%d/worker", startLevel)
	response, err := s.httpClient.
		GET(path).
		Send()

	if err != nil {
		return "", err
	}

	return response.Body().AsString()
}

func (s *Subsquid) GetBlocks(ctx context.Context, from, to uint64, workerUrl string) ([]*SqdBlockResponse, error) {
	var workerClient = fastshot.NewClient(workerUrl).
		Build()

	response, err := workerClient.POST("").
		Context().Set(ctx).
		Header().AddContentType(mime.JSON).
		Body().AsJSON(NewRequest(from, to)).
		Send()

	if err != nil {
		return nil, err
	}

	var result []*SqdBlockResponse
	err = response.Body().AsJSON(&result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *Subsquid) GetBlankBlocks(ctx context.Context, startLevel uint64, workerUrl string) ([]*SqdBlockResponse, error) {
	var workerClient = fastshot.NewClient(workerUrl).
		Build()

	response, err := workerClient.POST("").
		Context().Set(ctx).
		Header().AddContentType(mime.JSON).
		Body().AsJSON(NewSimpleRequest(startLevel)).
		Send()

	if err != nil {
		return nil, err
	}

	var result []*SqdBlockResponse
	err = response.Body().AsJSON(&result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *Subsquid) GetHead(_ context.Context) (uint64, error) {
	response, err := s.httpClient.
		GET("/height").
		Send()

	if err != nil {
		return 0, err
	}

	stringResponse, err := response.Body().AsString()
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(stringResponse, 10, 64)
}
