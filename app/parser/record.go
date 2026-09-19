package parser

type Record struct {
	HeaderSize int
	SerialNumbers []SerialNumber
	Values []Value
}

const (
	TableType = "table"
)

type SerialNumber struct {
	Type int
	Size int
}

type Value struct {
	Type int
	String string
	Integer int
}

const (
	nullSerialType = 0
	intSerialType = 1
	blobSerialType = 12
	stringSerialType = 13
)

func countSerialNumber(n uint64) (int, int) {
	switch true {
	case n >= 12 && n % 2 == 0:
		return blobSerialType, int((n - 12) / 2)
	
	case n >= 13 && n % 2 != 0:
		return stringSerialType, int((n - 13) / 2)

	case n == intSerialType:
		return intSerialType, 1
	}

	return 0, 0
}

func getValue(t int, buff []byte) Value {
	switch t {
	case stringSerialType:
		return Value{
			Type: stringSerialType,
			String: string(buff),
		}
	
	case blobSerialType:
		return Value{
			Type: blobSerialType,
			String: string(buff),
		}
	case intSerialType:
		return Value{
			Type: intSerialType,
			Integer: int(buff[0]),
		}
	}

	return Value{}
}

func (p *Parser) ReadRecord(offset int64) Record {
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

	var values []Value
	for _, serialNumber := range serialNumbers {
		buff, _ := p.GetBytes(serialNumber.Size, offset)
		values = append(values, getValue(serialNumber.Type, buff))
		offset += int64(serialNumber.Size)
	}

	return Record{
		HeaderSize: int(headerSize),
		SerialNumbers: serialNumbers,
		Values: values,
	}

}