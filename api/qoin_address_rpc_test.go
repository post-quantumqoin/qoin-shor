package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/go-jsonrpc"
	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/api/mocks"
	"github.com/post-quantumqoin/qoin-shor/core/types/ethtypes"
)

func TestQoinAddressRPCNames(t *testing.T) {
	service := &api.FullNodeStruct{}
	service.Internal.EthAddressToQoinAddress = func(_ context.Context, input ethtypes.EthAddress) (address.Address, error) {
		return input.ToQoinAddress()
	}
	service.Internal.QoinAddressToEthAddress = func(_ context.Context, input address.Address) (ethtypes.EthAddress, error) {
		return ethtypes.EthAddressFromQoinAddress(input)
	}
	server := jsonrpc.NewServer()
	server.Register("Qoin", service)
	input := ethtypes.EthAddress{0x12, 0x34, 0x56}
	qoinAddress, err := input.ToQoinAddress()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		method  string
		param   interface{}
		want    interface{}
		missing bool
	}{
		{"Qoin.EthAddressToQoinAddress", input, qoinAddress, false},
		{"Qoin.QoinAddressToEthAddress", qoinAddress, input, false},
		{"Qoin.EthAddressToFilecoinAddress", input, nil, true},
		{"Qoin.FilecoinAddressToEthAddress", qoinAddress, nil, true},
	} {
		t.Run(test.method, func(t *testing.T) {
			body, err := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": test.method, "params": []interface{}{test.param}})
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest("POST", "/rpc/v1", bytes.NewReader(body)))
			var decoded struct {
				Result json.RawMessage `json:"result"`
				Error  *struct {
					Code int `json:"code"`
				} `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
				t.Fatal(err)
			}
			if test.missing {
				if decoded.Error == nil || decoded.Error.Code != -32601 {
					t.Fatalf("old RPC must not exist: %s", response.Body.String())
				}
				return
			}
			if decoded.Error != nil {
				t.Fatalf("new RPC failed: %s", response.Body.String())
			}
			want, err := json.Marshal(test.want)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(decoded.Result, want) {
				t.Fatalf("got %s, want %s", decoded.Result, want)
			}
		})
	}
	internal := reflect.TypeOf(service.Internal)
	for _, name := range []string{"EthAddressToQoinAddress", "QoinAddressToEthAddress"} {
		field, ok := internal.FieldByName(name)
		if !ok || field.Tag.Get("perm") != "read" {
			t.Fatalf("RPC %s must retain read permission", name)
		}
	}
}

func TestQoinAddressRPCMocks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	client := mocks.NewMockFullNode(controller)
	ctx := context.Background()
	input := ethtypes.EthAddress{0x12, 0x34}
	qoinAddress, err := input.ToQoinAddress()
	if err != nil {
		t.Fatal(err)
	}
	client.EXPECT().EthAddressToQoinAddress(ctx, input).Return(qoinAddress, nil)
	if got, err := client.EthAddressToQoinAddress(ctx, input); err != nil || got != qoinAddress {
		t.Fatalf("new mock conversion failed: %v", err)
	}
	client.EXPECT().QoinAddressToEthAddress(ctx, qoinAddress).Return(input, nil)
	if got, err := client.QoinAddressToEthAddress(ctx, qoinAddress); err != nil || got != input {
		t.Fatalf("new mock conversion failed: %v", err)
	}
}
