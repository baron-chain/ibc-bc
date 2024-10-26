package errors

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// Error codes
const (
	CodeInvalidSequence = iota + 1
	CodeUnauthorized
	CodeInsufficientFunds
	CodeUnknownRequest
	CodeInvalidAddress
	CodeInvalidCoins
	CodeOutOfGas
	CodeInvalidRequest
	CodeInvalidHeight
	CodeInvalidVersion
	CodeInvalidChainID
	CodeInvalidType
	CodePackAny
	CodeUnpackAny
	CodeLogic
	CodeNotFound
)

// Error categories for grouping related errors
const (
	CategoryValidation = "validation"
	CategoryState     = "state"
	CategoryProtobuf  = "protobuf"
	CategorySecurity  = "security"
	CategoryInternal  = "internal"
)

// Register all errors with their respective codes and messages
var (
	// Authentication and Authorization
	ErrUnauthorized = errorsmod.Register(exported.ModuleName, CodeUnauthorized,
		"unauthorized access attempt")
	
	// Validation Errors
	ErrInvalidSequence = errorsmod.Register(exported.ModuleName, CodeInvalidSequence,
		"invalid sequence number for signature")
	ErrInvalidAddress = errorsmod.Register(exported.ModuleName, CodeInvalidAddress,
		"invalid address format or checksum")
	ErrInvalidCoins = errorsmod.Register(exported.ModuleName, CodeInvalidCoins,
		"invalid coin denomination or amount")
	ErrInvalidHeight = errorsmod.Register(exported.ModuleName, CodeInvalidHeight,
		"invalid block height specified")
	ErrInvalidVersion = errorsmod.Register(exported.ModuleName, CodeInvalidVersion,
		"invalid protocol version")
	ErrInvalidChainID = errorsmod.Register(exported.ModuleName, CodeInvalidChainID,
		"invalid chain identifier")
	ErrInvalidType = errorsmod.Register(exported.ModuleName, CodeInvalidType,
		"invalid type or format")
	ErrInvalidRequest = errorsmod.Register(exported.ModuleName, CodeInvalidRequest,
		"request contains invalid or malformed data")
	
	// State Errors
	ErrInsufficientFunds = errorsmod.Register(exported.ModuleName, CodeInsufficientFunds,
		"insufficient account balance for transaction")
	ErrOutOfGas = errorsmod.Register(exported.ModuleName, CodeOutOfGas,
		"operation exceeded gas limit")
	ErrNotFound = errorsmod.Register(exported.ModuleName, CodeNotFound,
		"requested entity not found in state")
	
	// Protocol Errors
	ErrUnknownRequest = errorsmod.Register(exported.ModuleName, CodeUnknownRequest,
		"unknown or unsupported request type")
	ErrPackAny = errorsmod.Register(exported.ModuleName, CodePackAny,
		"failed to pack protobuf message to Any type")
	ErrUnpackAny = errorsmod.Register(exported.ModuleName, CodeUnpackAny,
		"failed to unpack protobuf message from Any type")
	
	// Internal Errors
	ErrLogic = errorsmod.Register(exported.ModuleName, CodeLogic,
		"internal logic error - possible invariant violation")
)

// IsValidationError checks if the error is related to validation
func IsValidationError(err error) bool {
	code := errorsmod.ABCICode(err)
	return code == CodeInvalidSequence ||
		code == CodeInvalidAddress ||
		code == CodeInvalidCoins ||
		code == CodeInvalidHeight ||
		code == CodeInvalidVersion ||
		code == CodeInvalidChainID ||
		code == CodeInvalidType ||
		code == CodeInvalidRequest
}

// IsStateError checks if the error is related to state operations
func IsStateError(err error) bool {
	code := errorsmod.ABCICode(err)
	return code == CodeInsufficientFunds ||
		code == CodeOutOfGas ||
		code == CodeNotFound
}

// IsProtobufError checks if the error is related to protobuf operations
func IsProtobufError(err error) bool {
	code := errorsmod.ABCICode(err)
	return code == CodePackAny ||
		code == CodeUnpackAny
}

// IsSecurityError checks if the error is related to security
func IsSecurityError(err error) bool {
	code := errorsmod.ABCICode(err)
	return code == CodeUnauthorized
}

// WrapError wraps an existing error with additional context
func WrapError(err error, msg string) error {
	return errorsmod.Wrap(err, msg)
}
