package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	classQuery = `query ClassQuery ($limit:Int!, $offset: Int!) {
		class(offset: $offset, limit: $limit, order_by: {id: asc}) {
		  hash
		}
	}`
)

// JsonSchemaTester -
type JsonSchemaTester struct {
	client     *http.Client
	grpcClient *grpc.Client
	baseUrl    string
}

// NewJsonSchemaTester -
func NewJsonSchemaTester(grpcClient *grpc.Client, baseUrl string) JsonSchemaTester {
	tester := JsonSchemaTester{}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	tester.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: t,
	}
	tester.grpcClient = grpcClient
	tester.baseUrl = baseUrl
	return tester
}

// String -
func (js JsonSchemaTester) String() string {
	return "json schema tester"
}

// Test -
func (js JsonSchemaTester) Test(ctx context.Context) error {
	log.Info().Msg("start testing json schema...")
	var (
		limit  = 100
		offset = 0
		end    = false
	)
	for !end {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		classes, err := js.getClasses(ctx, limit, offset)
		if err != nil {
			return errors.Wrap(err, "receiving classes")
		}

		for i := range classes.Data.Classes {
			hash, err := hex.DecodeString(strings.TrimPrefix(classes.Data.Classes[i].Hash, "\\x"))
			if err != nil {
				return errors.Wrap(err, "decoding hash")
			}

			select {
			case <-ctx.Done():
				return nil
			default:
				schemaBytes, err := js.grpcClient.JsonSchemaForClass(ctx, &pb.Bytes{
					Data: hash,
				})
				if err != nil {
					return errors.Wrap(err, "receiving json schema")
				}

				var schema abi.JsonSchema
				if err := json.Unmarshal(schemaBytes.Data, &schema); err != nil {
					return errors.Wrap(err, "decoding json schema")
				}

				log.Info().Str("hash", classes.Data.Classes[i].Hash).Msg("success")
			}
		}

		offset += len(classes.Data.Classes)
		end = len(classes.Data.Classes) < limit
	}

	return nil
}

// Close -
func (js JsonSchemaTester) Close() error {
	return nil
}

// ClassResponse -
type ClassResponse struct {
	Data struct {
		Classes []struct {
			Hash string `json:"hash"`
		} `json:"class"`
	} `json:"data"`
}

// GraphQLRequest -
type GraphQLRequest struct {
	Name      string         `json:"operationName"`
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func (js JsonSchemaTester) getClasses(ctx context.Context, limit, offset int) (*ClassResponse, error) {
	requestBody := GraphQLRequest{
		Name:  "ClassQuery",
		Query: classQuery,
		Variables: map[string]any{
			"limit":  limit,
			"offset": offset,
		},
	}
	reader := new(bytes.Buffer)
	if err := json.NewEncoder(reader).Encode(requestBody); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, js.baseUrl, reader)
	if err != nil {
		return nil, errors.Errorf("makeGetRequest.NewRequest: %v", err)
	}
	resp, err := js.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var response ClassResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}
