package bindingsfactory

import (
	"encoding/json"

	filetreetypes "github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/types/filetree"
	storagetypes "github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/types/storage"
)

type InstantiateMsg struct {
	BindingsCodeId int `json:"bindings_code_id"`
}

// not sure if 'create_bindings_v2' is correct
type ExecuteMsg struct {
	CreateBindings      *ExecuteMsg_CreateBindings      `json:"create_bindings,omitempty"`
	FundBindings        *ExecuteMsg_FundBindings        `json:"fund_bindings,omitempty"`
	CallBindings        *ExecuteMsg_CallBindings        `json:"call_bindings,omitempty"`
	CallStorageBindings *ExecuteMsg_CallStorageBindings `json:"call_storage_bindings,omitempty"`
	AddToWhiteList      *ExecuteMsg_AddToWhiteList      `json:"add_to_white_list,omitempty"`
}

type ExecuteMsg_AddToWhiteList struct {
	JKLAddress *string `json:"jkl_address,omitempty"`
}

type ExecuteMsg_CreateBindings struct {
	UserEvmAddress *string `json:"user_evm_address,omitempty"`
}

type ExecuteMsg_CallBindings struct {
	EvmAddress *string                   `json:"evm_address,omitempty"`
	Msg        *filetreetypes.ExecuteMsg `json:"msg,omitempty"`
}

type ExecuteMsg_FundBindings struct {
	EvmAddress *string `json:"evm_address,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
}

type ExecuteMsg_CallStorageBindings struct {
	EvmAddress *string                  `json:"evm_address,omitempty"`
	Msg        *storagetypes.ExecuteMsg `json:"msg,omitempty"`
}

// ToString returns a string representation of the message
func (m *ExecuteMsg) ToString() string {
	return toString(m)
}

func toString(v any) string {
	jsonBz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(jsonBz)
}
