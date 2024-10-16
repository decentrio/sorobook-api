package contract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/stellar/go/xdr"
)

const (
	XDR_BOOL       = "bool"
	XDR_U32        = "u32"
	XDR_I32        = "i32"
	XDR_U64        = "u64 "
	XDR_I64        = "i64 "
	XDR_TIME_POINT = "time_point"
	XDR_DURATION   = "duration"
	XDR_U128       = "u128"
	XDR_I128       = "i128"
	XDR_U256       = "u256"
	XDR_I256       = "i256"
	XDR_BYTES      = "bytes"
	XDR_STRING     = "string "
	XDR_SYM        = "symbol"
	XDR_NONCE      = "nonce"
	XDR_VEC        = "vec"
	XDR_MAP        = "map"
	XDR_ADDRESS    = "address"
)

func convertToData(keyType string, keyValue string) (xdr.ScVal, error) {
	switch keyType {
	case XDR_BOOL:
		return convertToDataBool(keyValue)
	case XDR_U32:
		return convertToDataUint32(keyValue)
	case XDR_I32:
		return convertToDataInt32(keyValue)
	case XDR_U64:
		return convertToDataUint64(keyValue)
	case XDR_I64:
		return convertToDataInt64(keyValue)
	case XDR_TIME_POINT:
		return convertToDataTimePoint(keyValue)
	case XDR_DURATION:
		return convertToDataDuration(keyValue)
	case XDR_U128:
		return convertToDataUInt128Parts(keyValue)
	case XDR_I128:
		return convertToDataInt128Parts(keyValue)
	case XDR_U256:
		return convertToDataUInt256Parts(keyValue)
	case XDR_I256:
		return convertToDataInt256Parts(keyValue)
	case XDR_BYTES:
		return convertToDataScBytes(keyValue)
	case XDR_STRING:
		return convertToDataScString(keyValue)
	case XDR_SYM:
		return convertToDataScSymbol(keyValue)
	case XDR_NONCE:
		return convertToDataScNonceKey(keyValue)
	case XDR_ADDRESS:
		return convertToDataScAddress(keyValue)
	default:
		return xdr.ScVal{}, errors.New("convert false")
	}
}

func convertToDataBool(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseBool(value)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvBool, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUint32(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt32(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvI32, xdr.Int32(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUint64(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt64(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvI64, xdr.Int64(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataTimePoint(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvTimepoint, xdr.TimePoint(xdr.Uint64(data)))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataDuration(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvDuration, xdr.Duration(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUInt128Parts(value string) (xdr.ScVal, error) {
	data, err := parseU128String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvU128, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt128Parts(value string) (xdr.ScVal, error) {
	data, err := parseI128String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvI128, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUInt256Parts(value string) (xdr.ScVal, error) {
	data, err := parseU256String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvU256, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt256Parts(value string) (xdr.ScVal, error) {
	data, err := parseI256String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvI256, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScBytes(value string) (xdr.ScVal, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScString(value string) (xdr.ScVal, error) {
	res, err := xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString(value))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScSymbol(value string) (xdr.ScVal, error) {
	res, err := xdr.NewScVal(xdr.ScValTypeScvSymbol, xdr.ScSymbol(value))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScNonceKey(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvLedgerKeyNonce, xdr.ScNonceKey{
		Nonce: xdr.Int64(data),
	})

	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScAddress(value string) (xdr.ScVal, error) {
	aid := xdr.MustAddress(value)
	res, err := xdr.NewScVal(xdr.ScValTypeScvAddress, xdr.ScAddress{
		AccountId: &aid,
	})

	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}


func parseU128String(s string) (xdr.UInt128Parts, error) {
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10)
	if !ok {
		return xdr.UInt128Parts{}, fmt.Errorf("invalid number format")
	}

	// Mask for lower 64 bits
	lowMask := new(big.Int).SetUint64(^uint64(0))

	// Extract the lower 64 bits
	low := new(big.Int).And(bigInt, lowMask).Uint64()

	// Extract the higher 64 bits by shifting right 64 bits
	high := new(big.Int).Rsh(bigInt, 64).Uint64()

	return xdr.UInt128Parts{
		Hi: xdr.Uint64(high),
		Lo: xdr.Uint64(low),
	}, nil
}

func parseI128String(s string) (xdr.Int128Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.Int128Parts{}, fmt.Errorf("invalid number format")
	}

	// Handle negative numbers for signed 128-bit integers
	negative := bigInt.Sign() < 0
	if negative {
		bigInt = bigInt.Abs(bigInt) // Convert to positive for bitwise operations
	}

	// Mask for lower 64 bits
	lowMask := new(big.Int).SetUint64(^uint64(0))

	// Extract the lower 64 bits
	low := new(big.Int).And(bigInt, lowMask).Uint64()

	// Extract the higher 64 bits and cast to int64 for signed interpretation
	high := new(big.Int).Rsh(bigInt, 64).Int64()
	if negative {
		high = -high
	}

	return xdr.Int128Parts{
		Hi: xdr.Int64(high),
		Lo: xdr.Uint64(low),
	}, nil
}

func parseU256String(s string) (xdr.UInt256Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.UInt256Parts{}, fmt.Errorf("invalid number format")
	}

	// Mask for 64 bits
	mask64 := new(big.Int).SetUint64(^uint64(0))

	// Extract the four 64-bit parts
	lowLow := new(big.Int).And(bigInt, mask64).Uint64()
	lowHigh := new(big.Int).Rsh(bigInt, 64).And(bigInt, mask64).Uint64()
	highLow := new(big.Int).Rsh(bigInt, 128).And(bigInt, mask64).Uint64()
	highHigh := new(big.Int).Rsh(bigInt, 192).Uint64()

	return xdr.UInt256Parts{
		HiHi: xdr.Uint64(highHigh),
		HiLo: xdr.Uint64(highLow),
		LoHi: xdr.Uint64(lowHigh),
		LoLo: xdr.Uint64(lowLow),
	}, nil
}

func parseI256String(s string) (xdr.Int256Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.Int256Parts{}, fmt.Errorf("invalid number format")
	}

	// Handle negative numbers
	negative := bigInt.Sign() < 0
	if negative {
		bigInt = bigInt.Abs(bigInt)
	}

	// Mask for 64 bits
	mask64 := new(big.Int).SetUint64(^uint64(0))

	// Extract the four 64-bit parts
	lowLow := new(big.Int).And(bigInt, mask64).Uint64()
	lowHigh := new(big.Int).Rsh(bigInt, 64).And(bigInt, mask64).Uint64()
	highLow := new(big.Int).Rsh(bigInt, 128).And(bigInt, mask64).Uint64()
	highHigh := new(big.Int).Rsh(bigInt, 192).Int64()

	if negative {
		highHigh = -highHigh
	}

	return xdr.Int256Parts{
		HiHi: xdr.Int64(highHigh),
		HiLo: xdr.Uint64(highLow),
		LoHi: xdr.Uint64(lowHigh),
		LoLo: xdr.Uint64(lowLow),
	}, nil
}
