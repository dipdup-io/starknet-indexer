package indexer

import (
	"os"
	"testing"

	starknet "github.com/dipdup-io/starknet-go-api/pkg/api"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Test_decimalFromHex(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want decimal.Decimal
	}{
		{
			name: "test 1",
			s:    "0xa",
			want: decimal.RequireFromString("10"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decimalFromHex(tt.s); !got.Equal(tt.want) {
				t.Errorf("decimalFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInternalModels(t *testing.T) {
	tests := []struct {
		filename string
		want     models.Block
		wantErr  bool
	}{
		{
			filename: "./tests/1.json",
			want: models.Block{
				Height:             2,
				Time:               1637084470,
				Hash:               "0x4e1f77f39545afe866ac151ac908bd1a347a2a8a7d58bef1276db4f06fdf2f6",
				ParentHash:         "0x2a70fb03fe363a2d6be843343a1d81ce6abeda1e9bd5cc6ad8fa9f45e30fdeb",
				NewRoot:            "0x3ceee867d50b5926bb88c0ec7e0b9c20ae6b537e74aac44b8fcf6bb6da138d9",
				SequencerAddress:   "0x37b2cd6baaa515f520383bee7b7094f892f4c770695fc329a8973e841a971ae",
				Status:             models.StatusAcceptedOnL1,
				InvokeV0Count:      1,
				InvokeV1Count:      1,
				DeclareCount:       1,
				DeployCount:        1,
				DeployAccountCount: 1,
				L1HandlerCount:     1,
				TxCount:            6,

				InvokeV0: []models.InvokeV0{
					{
						Height:             2,
						Time:               1637084470,
						Status:             models.StatusAcceptedOnL1,
						Hash:               "0x2e530fe2f39ba92380de33cfca060f68c2f50b8af954dae7370c97bf97e1e55",
						MaxFee:             decimal.RequireFromString("0"),
						Signature:          []string{},
						Nonce:              decimal.RequireFromString("0"),
						ContractAddress:    "0x2d6c9569dea5f18628f1ef7c15978ee3093d2d3eec3b893aac08004e678ead3",
						EntrypointSelector: "0x12ead94ae9d3f9d2bdb6b847cf255f1f398193a1f88884a0ae8e18f24a037b6",
						CallData: []string{
							"0xdaee7b1ac98d5d3fa7cf5dcfa0dd5f47dc8728fc",
						},
					},
				},
				InvokeV1: []models.InvokeV1{
					{
						Height: 2,
						Time:   1637084470,
						Status: models.StatusAcceptedOnL1,
						Hash:   "0x6e9e4cba46dd68732525d0eeca23214f40e98dabfd96b2cf65c19b6a4dabb70",
						MaxFee: decimal.RequireFromString("8999999999999999"),
						Signature: []string{
							"0x58c9ac52935f0ca6024c2da1941c49f6efc0d64b61a1f91a50e225d3a8f777d",
							"0x75a7e4f1fdb5868e9eaf12729f3a9cd715c142a6c660b7834665eb2ede53ae",
						},
						Nonce:         decimal.RequireFromString("29007"),
						SenderAddress: "0x7b393627bd514d2aa4c83e9f0c468939df15ea3c29980cd8e7be3ec847795f0",
						CallData: []string{
							"0x1",
							"0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
							"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
							"0x0",
							"0x3",
							"0x3",
							"0x64bcb38d84e46553ca4bd7c2eee05b9e97103153fffd1cfb5e249b7533df2c4",
							"0x50417c3d79e059",
							"0x0",
						},
					},
				},
				Declare: []models.Declare{
					{
						Height:        2,
						Time:          1637084470,
						Status:        models.StatusAcceptedOnL1,
						Hash:          "0x5e764bdf83525030bed0bceef7b31dd87542a200755f0fa48cb51234246ca71",
						MaxFee:        decimal.RequireFromString("2000000000000000"),
						SenderAddress: "0x7cd2dad3730cc956eb2587ff1bb94742300adcac26f4bfd98d2a4b43c1d1b82",
						ClassHash:     "0x5a55d60ea6d6515b5aab71027828dcf49225e3ceb9f87771ba1c8a12afcc381",
						Signature: []string{
							"0x613378b52f28a06ea8b1fde3e9d7c66685b7fe05a0140b4a294317c5b96c69f",
							"0x63d3c7e08953e17305106214fcd0f8593fc86878fd6b9223790f7b9e1291a78",
						},
						Nonce: decimal.RequireFromString("8"),
					},
				},
				Deploy: []models.Deploy{
					{
						Height:              2,
						Time:                1637084470,
						Status:              models.StatusAcceptedOnL1,
						Hash:                "0x5a8629d7852d3c8f4fda51d83b48cc8b2184763c46383419c1beeadaea1e66e",
						ClassHash:           "0x10455c752b86932ce552f2b0fe81a880746649b9aee7e0d842bf3f52378f9f8",
						ContractAddressSalt: "0x23a93d3a3463ac1539852fcb9dbf58ed9581e4abbb4a828889768fbbbdb9bcd",
						ConstructorCalldata: []string{
							"0x7f93985c1baa5bd9b2200dd2151821bd90abb87186d0be295d7d4b9bc8ca41f",
							"0x127cd00a078199381403a33d315061123ce246c8e5f19aa7f66391a9d3bf7c6",
						},
					},
				},
				DeployAccount: []models.DeployAccount{
					{
						Height: 2,
						Time:   1637084470,
						Status: models.StatusAcceptedOnL1,
						Hash:   "0x4dbfc1059b7c7710417af8ade1cd76384cf185a872de80c031b9df99aa3b2f",
						MaxFee: decimal.RequireFromString("489564191908984"),
						Nonce:  decimal.RequireFromString("0"),
						Signature: []string{
							"0x6376672530f83d9aea93594759cb6f6f47978cf63e8cba5b9a5cc16e02a7f98",
							"0x7957c54e83dccb4568c61bcc9c97c502e265bb706c00e76c5e530b55464164a",
						},
						ContractAddressSalt: "0x4ddd1c9b8c3607b34801a86bfd970974ea34bbb69edc97e73435c007e3db11",
						ConstructorCalldata: []string{
							"0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2",
							"0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463",
							"0x2",
							"0x4ddd1c9b8c3607b34801a86bfd970974ea34bbb69edc97e73435c007e3db11",
							"0x0",
						},
						ClassHash: "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918",
					},
				},
				L1Handler: []models.L1Handler{
					{
						Height:             2,
						Time:               1637084470,
						Status:             models.StatusAcceptedOnL1,
						Hash:               "0x3ee1a1881f293154ae375175462eea8a3eea09320267e7d7c1e8da5410da2b",
						Nonce:              decimal.RequireFromString("165827"),
						ContractAddress:    "0x73314940630fd6dcda0d772d4c972c4e0a9946bef9dabf4ef84eda8ef542b82",
						EntrypointSelector: "0x2d757788a8d8d6f21d1cd40bce38a8222d70654214e96ff95d8086e684fbee5",
						CallData: []string{
							"0xae0ee0a63a2ce6baeeffe56e7714fb4efe48d419",
							"0x59e37c6c739dbd94a19c3de69345c5e73d540ad27dfd4de69d242936f1e8e6c",
							"0x470de4df820000",
							"0x0",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			f, err := os.Open(tt.filename)
			if err != nil {
				t.Error(err)
				return
			}

			var block starknet.BlockWithTxs
			if err := json.NewDecoder(f).Decode(&block); err != nil {
				t.Error(err)
				return
			}

			if err := f.Close(); err != nil {
				t.Error(err)
				return
			}

			got, err := getInternalModels(block)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInternalModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
