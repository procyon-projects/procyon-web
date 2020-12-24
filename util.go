package web

import (
	"math/bits"
	"strconv"
)

func strToBool(value string) interface{} {
	result, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}
	return result
}

func strToFloat32(value string) interface{} {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return float32(result)
}

func strToFloat64(value string) interface{} {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return result
}

func strToInt(value string) interface{} {
	result, err := strconv.ParseInt(value, 10, bits.UintSize)
	if err != nil {
		panic(err)
	}
	return int(result)
}

func strToInt8(value string) interface{} {
	result, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		panic(err)
	}
	return int8(result)
}

func strToInt16(value string) interface{} {
	result, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		panic(err)
	}
	return int16(result)
}

func strToInt32(value string) interface{} {
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(result)
}

func strToInt64(value string) interface{} {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func strToUInt(value string) interface{} {
	result, err := strconv.ParseUint(value, 10, bits.UintSize)
	if err != nil {
		panic(err)
	}
	return uint(result)
}

func strToUInt8(value string) interface{} {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		panic(err)
	}
	return uint8(result)
}

func strToUInt16(value string) interface{} {
	result, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		panic(err)
	}
	return uint16(result)
}

func strToUInt32(value string) interface{} {
	result, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(result)
}

func strToUInt64(value string) interface{} {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}
