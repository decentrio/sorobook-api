package contract

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/stellar/go/xdr"
)

const (
	cst_bool         = "bool"
	cst_Uint32       = "Uint32"
	cst_Int32        = "Int32"
	cst_Uint64       = "Uint64 "
	cst_Int64        = "Int64 "
	cst_TimePoint    = "TimePoint"
	cst_Duration     = "Duration"
	cst_UInt128Parts = "U128"
	cst_Int128Parts  = "I128"
	cst_UInt256Parts = "U256"
	cst_Int256Parts  = "I256"
	cst_ScBytes      = "Bytes"
	cst_ScString     = "String "
	cst_ScSymbol     = "Symbol"
	cst_ScNonceKey   = "NonceKey "
	cst_ScVec        = "Vec"
	cst_ScMap        = "Map"
)

func convertToData(key_type string, key_name string) (xdr.ScValType, interface{}, error) {
	values := strings.Split(key_name, "#")

	switch key_type {
	case cst_bool:
		if len(values) == 1 {
			return convertToDataBool(values[0])
		}
	case cst_Uint32:
		if len(values) == 1 {
			return convertToDataUint32(values[0])
		}
	case cst_Int32:
		if len(values) == 1 {
			return convertToDataInt32(values[0])
		}
	case cst_Uint64:
		if len(values) == 1 {
			return convertToDataUint64(values[0])
		}
	case cst_Int64:
		if len(values) == 1 {
			return convertToDataInt64(values[0])
		}
	case cst_TimePoint:
		if len(values) == 1 {
			return convertToDataTimePoint(values[0])
		}
	case cst_Duration:
		if len(values) == 1 {
			return convertToDataDuration(values[0])
		}
	case cst_UInt128Parts:
		if len(values) == 2 {
			return convertToDataUInt128Parts(values[0], values[1])
		}
	case cst_Int128Parts:
		if len(values) == 2 {
			return convertToDataInt128Parts(values[0], values[1])
		}
	case cst_UInt256Parts:
		if len(values) == 4 {
			return convertToDataUInt256Parts(values[0], values[1], values[2], values[3])
		}
	case cst_Int256Parts:
		if len(values) == 4 {
			return convertToDataInt256Parts(values[0], values[1], values[2], values[3])
		}
	case cst_ScBytes:
		if len(values) == 1 {
			return convertToDataScBytes(values[0])
		}
	case cst_ScString:
		if len(values) == 1 {
			return convertToDataScString(values[0])
		}
	case cst_ScSymbol:
		if len(values) == 1 {
			return convertToDataScSymbol(values[0])
		}
	case cst_ScNonceKey:
		if len(values) == 1 {
			return convertToDataScNonceKey(values[0])
		}
	case cst_ScVec:
		if len(values) == 1 {
			return convertToDataScVec(values[0])
		}
	case cst_ScMap:
		if len(values) == 1 {
			return convertToDataScMap(values[0])
		}
	default:
		return 0, nil, errors.New("convert false")
	}

	return 0, nil, errors.New("convert false")
}

func convertToDataBool(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseBool(value)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvBool, data, err
}

func convertToDataUint32(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvU32, xdr.Uint32(data), err
}

func convertToDataInt32(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvI32, xdr.Int32(data), err
}

func convertToDataUint64(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvU64, xdr.Uint64(data), err
}

func convertToDataInt64(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvI64, xdr.Int64(data), err
}

func convertToDataTimePoint(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvTimepoint, xdr.TimePoint(xdr.Uint64(data)), err
}

func convertToDataDuration(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvDuration, xdr.Duration(xdr.Uint64(data)), err
}

func convertToDataUInt128Parts(value1 string, value2 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseUint(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.UInt128Parts = xdr.UInt128Parts{
		Hi: xdr.Uint64(data1),
		Lo: xdr.Uint64(data2),
	}

	return xdr.ScValTypeScvU128, data, err
}

func convertToDataInt128Parts(value1 string, value2 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.Int128Parts = xdr.Int128Parts{
		Hi: xdr.Int64(data1),
		Lo: xdr.Uint64(data2),
	}

	return xdr.ScValTypeScvI128, data, err
}

func convertToDataUInt256Parts(value1 string, value2 string, value3 string, value4 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseUint(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data3, err := strconv.ParseUint(value3, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data4, err := strconv.ParseUint(value4, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.UInt256Parts = xdr.UInt256Parts{
		HiHi: xdr.Uint64(data1),
		HiLo: xdr.Uint64(data2),
		LoHi: xdr.Uint64(data3),
		LoLo: xdr.Uint64(data4),
	}

	return xdr.ScValTypeScvU256, data, err
}

func convertToDataInt256Parts(value1 string, value2 string, value3 string, value4 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data3, err := strconv.ParseUint(value3, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data4, err := strconv.ParseUint(value4, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.Int256Parts = xdr.Int256Parts{
		HiHi: xdr.Int64(data1),
		HiLo: xdr.Uint64(data2),
		LoHi: xdr.Uint64(data3),
		LoLo: xdr.Uint64(data4),
	}

	return xdr.ScValTypeScvI256, data, err
}

func convertToDataScBytes(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvBytes, xdr.ScBytes([]byte(value)), nil
}

func convertToDataScString(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvString, xdr.ScString(value), nil
}

func convertToDataScSymbol(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvSymbol, xdr.ScSymbol(value), nil
}

func convertToDataScNonceKey(value string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.ScNonceKey = xdr.ScNonceKey{
		Nonce: xdr.Int64(data1),
	}

	return xdr.ScValTypeScvLedgerKeyNonce, data, err
}

func convertToDataScVec(value string) (xdr.ScValType, interface{}, error) {
	items := strings.Split(value, "#")
	var dataVec xdr.ScVec

	for _, item := range items {
		s := strings.Split(item, ",")
		if len(s) != 2 {
			return 0, nil, errors.New("split s false " + fmt.Sprint(len(s)))
		}

		xdrType, data, err := convertToData(s[0], s[1])
		if err != nil {
			return 0, nil, err
		}

		xdrKey, err := xdr.NewScVal(xdrType, data)
		if err != nil {
			return 0, nil, err
		}

		dataVec = append(dataVec, xdrKey)
	}

	return xdr.ScValTypeScvVec, dataVec, nil
}

func convertToDataScMap(value string) (xdr.ScValType, interface{}, error) {
	items := strings.Split(value, "#")
	var dataMap xdr.ScMap

	for _, item := range items {
		s := strings.Split(item, ",")
		if len(s) != 4 {
			return 0, nil, errors.New("split s false " + fmt.Sprint(len(s)))
		}

		xdrType, data1, err := convertToData(s[0], s[1])
		if err != nil {
			return 0, nil, err
		}

		xdrKey1, err := xdr.NewScVal(xdrType, data1)
		if err != nil {
			return 0, nil, err
		}

		xdrType, data2, err := convertToData(s[2], s[3])
		if err != nil {
			return 0, nil, err
		}

		xdrKey2, err := xdr.NewScVal(xdrType, data2)
		if err != nil {
			return 0, nil, err
		}

		var itemMap xdr.ScMapEntry = xdr.ScMapEntry{
			Key: xdrKey1,
			Val: xdrKey2,
		}

		dataMap = append(dataMap, itemMap)
	}

	return xdr.ScValTypeScvMap, dataMap, nil
}
