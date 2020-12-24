package web

import (
	"github.com/procyon-projects/goo"
	"reflect"
	"strings"
)

type valueConverterFunction func(val string) interface{}

var requestObjectMetadataMap = make(map[reflect.Type]*RequestObjectMetadata, 0)

type RequestObjectMetadata struct {
	typ            reflect.Type
	hasOnlyBody    bool
	bodyMetadata   *requestBodyMetadata
	paramMetadata  *requestParamMetadata
	pathMetadata   *requestPathMetadata
	headerMetadata *requestHeaderMetadata
}

func newRequestObjectMetadata() *RequestObjectMetadata {
	return &RequestObjectMetadata{
		bodyMetadata:   newRequestBodyMetadata(),
		paramMetadata:  newRequestParamMetadata(),
		pathMetadata:   newRequestPathMetadata(),
		headerMetadata: newRequestHeaderMetadata(),
	}
}

type fieldMetadata struct {
	index     int
	name      string
	typ       goo.Type
	extra     int
	converter valueConverterFunction
}

type requestBodyMetadata struct {
	fieldIndex int
}

func newRequestBodyMetadata() *requestBodyMetadata {
	return &requestBodyMetadata{
		fieldIndex: -1,
	}
}

type requestParamMetadata struct {
	fieldIndex int
	paramMap   map[string]*fieldMetadata
}

func newRequestParamMetadata() *requestParamMetadata {
	return &requestParamMetadata{
		fieldIndex: -1,
		paramMap:   make(map[string]*fieldMetadata, 0),
	}
}

type requestPathMetadata struct {
	fieldIndex      int
	pathVariableMap map[string]*fieldMetadata
}

func newRequestPathMetadata() *requestPathMetadata {
	return &requestPathMetadata{
		fieldIndex:      -1,
		pathVariableMap: make(map[string]*fieldMetadata, 0),
	}
}

type requestHeaderMetadata struct {
	fieldIndex int
	headerMap  map[string]*fieldMetadata
}

func newRequestHeaderMetadata() *requestHeaderMetadata {
	return &requestHeaderMetadata{
		fieldIndex: -1,
		headerMap:  make(map[string]*fieldMetadata, 0),
	}
}

func ScanRequestObjectMetadata(requestObject interface{}) *RequestObjectMetadata {
	requestObjType := goo.GetType(requestObject)
	if !requestObjType.IsStruct() {
		panic("Request object must be a type struct")
	}

	structType := requestObjType.ToStructType()
	structFields := structType.GetFields()

	hasUntypedStructField := false
	hasFields := false

	requestObjectMetadata := newRequestObjectMetadata()

	for index, structField := range structFields {
		fieldType := structField.GetType()

		if fieldType.IsStruct() && strings.HasPrefix(fieldType.GetName(), "struct") {
			hasUntypedStructField = true
		} else {
			hasFields = true
		}

		if hasUntypedStructField && hasFields {
			panic("RequestObject must only consist of untyped struct or fields completely")
		}

		if hasFields {
			continue
		}

		requestTag, err := structField.GetTagByName("request")

		if err != nil {
			panic("Untyped struct in Request Object must has a 'request' tag ")
		}

		structFieldType := structField.GetType().ToStructType()

		switch requestTag.Value {
		case "param":
			requestObjectMetadata.paramMetadata.fieldIndex = index
			traverseFields(structFieldType, requestObjectMetadata.paramMetadata.paramMap)
		case "body":
			requestObjectMetadata.bodyMetadata.fieldIndex = index
		case "path":
			requestObjectMetadata.pathMetadata.fieldIndex = index
			traverseFields(structFieldType, requestObjectMetadata.pathMetadata.pathVariableMap)
		case "header":
			requestObjectMetadata.headerMetadata.fieldIndex = index
			traverseFields(structFieldType, requestObjectMetadata.headerMetadata.headerMap)
		default:
			panic("Invalid request tag value")
		}

	}

	if hasFields {
		requestObjectMetadata.hasOnlyBody = true
	}

	requestObjectMetadata.typ = structType.GetGoType()
	requestObjectMetadataMap[structType.GetGoType()] = requestObjectMetadata

	return requestObjectMetadata
}

func traverseFields(requestStruct goo.Struct, fieldMap map[string]*fieldMetadata) int {
	if requestStruct == nil || fieldMap == nil {
		return 0
	}

	fieldCount := requestStruct.GetExportedFieldCount()
	for index, field := range requestStruct.GetExportedFields() {
		fieldType := field.GetType()

		if !fieldType.IsString() && !fieldType.IsBoolean() && !fieldType.IsNumber() {
			panic("Fields could be string, boolean and number types")
		}

		fieldMetadata := &fieldMetadata{
			index:     index,
			name:      field.GetName(),
			typ:       fieldType,
			converter: getConverterFunction(fieldType),
			extra:     -1,
		}

		jsonTag, err := field.GetTagByName("json")
		if err == nil {
			if jsonTag.Value == "" {
				panic("wtf : tag value cannot be empty, " + field.GetName())
			}
			fieldMap[jsonTag.Value] = fieldMetadata
		} else {
			yamlTag, err := field.GetTagByName("yaml")
			if err != nil {
				panic("wtf : you have to add a json or yaml tag, " + field.GetName())
			}

			if yamlTag.Value == "" {
				panic("wtf : tag value cannot be empty, " + field.GetName())
			}
			fieldMap[yamlTag.Value] = fieldMetadata
		}

	}
	return fieldCount
}

func getConverterFunction(typ goo.Type) valueConverterFunction {
	if typ.IsNumber() {
		numberType := typ.ToNumberType()
		if numberType.GetType() == goo.IntegerType {
			integerType := numberType.(goo.Integer)
			return getIntegerConverterFunction(integerType)
		} else if numberType.GetType() == goo.FloatType {
			floatType := numberType.(goo.Float)
			return getFloatConverterFunction(floatType)
		}
		panic("type must be integer or float")
	} else if typ.IsBoolean() {
		return strToBool
	} else if typ.IsString() {
		return nil
	}
	panic("type must be string, number or boolean")
}

func getIntegerConverterFunction(integerType goo.Integer) valueConverterFunction {
	if integerType.IsSigned() {
		if integerType.GetName() == "int" {
			return strToInt
		}
		switch integerType.GetBitSize() {
		case goo.BitSize8:
			return strToInt8
		case goo.BitSize16:
			return strToInt16
		case goo.BitSize32:
			return strToInt32
		case goo.BitSize64:
			return strToInt64
		default:
			panic("Wtf!")
		}
	} else {
		if integerType.GetName() == "int" {
			return strToUInt
		}
		switch integerType.GetBitSize() {
		case goo.BitSize8:
			return strToUInt8
		case goo.BitSize16:
			return strToUInt16
		case goo.BitSize32:
			return strToUInt32
		case goo.BitSize64:
			return strToUInt64
		default:
			panic("Wtf!")
		}
	}
}

func getFloatConverterFunction(floatType goo.Float) valueConverterFunction {
	if floatType.GetBitSize() == goo.BitSize8 {
		return strToFloat32
	}
	return strToFloat64
}
