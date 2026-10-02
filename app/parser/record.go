package parser

import "strconv"

type Record struct {
	HeaderSize int
	RowId int
	SerialNumbers []SerialNumber
	Values []*Value
}

const (
	TableType = "table"
	IndexType = "index"
)

type SerialNumber struct {
	Type int
	Size int
}

type Value struct {
	Type int
	String string
	Integer int64
}

const (
	nullSerialType = 0
	intSerialType = 1
	floatSerialType = 7
	zeroSerialType = 8
	oneSerialType = 9
	blobSerialType = 12
	stringSerialType = 13
)

// размеры целых serial type 1..6 в байтах
var intSizes = []int{1, 2, 3, 4, 6, 8}

func countSerialNumber(n uint64) (int, int) {
	switch true {
	case n >= 12 && n % 2 == 0:
		return blobSerialType, int((n - 12) / 2)
	
	case n >= 13 && n % 2 != 0:
		return stringSerialType, int((n - 13) / 2)

	case n >= 1 && n <= 6:
		return intSerialType, intSizes[n - 1]

	case n == floatSerialType:
		return floatSerialType, 8

	case n == zeroSerialType:
		return zeroSerialType, 0

	case n == oneSerialType:
		return oneSerialType, 0
	}

	return nullSerialType, 0
}

func newIntValue(v int64) *Value {
	return &Value{
		Type: intSerialType,
		String: strconv.FormatInt(v, 10),
		Integer: v,
	}
}

func getValue(t int, buff []byte) *Value {
	switch t {
	case stringSerialType:
		return &Value{
			Type: stringSerialType,
			String: string(buff),
		}
	
	case blobSerialType:
		return &Value{
			Type: blobSerialType,
			String: string(buff),
		}
	case intSerialType:
		// big-endian, знак берётся из старшего байта
		v := int64(int8(buff[0]))
		for _, b := range buff[1:] {
			v = v<<8 | int64(b)
		}
		return newIntValue(v)
	case zeroSerialType:
		return newIntValue(0)
	case oneSerialType:
		return newIntValue(1)
	}

	return &Value{}
}

func (p *Parser) ReadRecord(offset int64) *Record {
	headerSize, n1 := p.GetVarint(offset)
	offset += int64(n1)

	var serialNumbers []SerialNumber
	var start uint64 = uint64(n1)
	for start < headerSize {
		serialNumber, n := p.GetVarint(offset)
		serialType, size := countSerialNumber(serialNumber)
		serialNumbers = append(serialNumbers, SerialNumber{Type: serialType, Size: size})
		offset = offset + int64(n)
		start += uint64(n)
	}

	var values []*Value
	for _, serialNumber := range serialNumbers {
		buff, _ := p.GetBytes(serialNumber.Size, offset)
		values = append(values, getValue(serialNumber.Type, buff))
		offset += int64(serialNumber.Size)
	}

	return &Record{
		HeaderSize: int(headerSize),
		SerialNumbers: serialNumbers,
		Values: values,
	}
}

