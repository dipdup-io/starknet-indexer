package subsquid

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

func (s *Subsquid) GetWorkerUrl(ctx context.Context, startLevel uint64) (string, error) {
	path := fmt.Sprintf("/%d/worker", startLevel)
	response, err := s.httpClient.
		GET(path).
		Send()

	if err != nil {
		return "", err
	}

	return response.Body().AsString()
}

func (s *Subsquid) GetData(ctx context.Context, startLevel uint64) {
	workerUrl, err := s.GetWorkerUrl(ctx, startLevel)
	if err != nil {
		return
	}

	var workerClient = fastshot.NewClient(workerUrl).
		Build()

	response, err := workerClient.POST("").
		Header().AddContentType(mime.JSON).
		Body().AsJSON(NewRequest(startLevel)).
		Send()

	if err != nil {
		return
	}

	var result []SqdBlockResponse
	err = response.Body().AsJSON(&result)
	if err != nil {
		return
	}
	fmt.Println("done")
}

func (s *Subsquid) GetHead(ctx context.Context) (uint64, error) {
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
