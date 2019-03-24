package model

type FinsOperateInterface interface {
	ReadBites(memoryArea byte, address uint16, bitOffset byte, readCount uint16) ([]bool, error)
	ReadWords(memoryArea byte, address uint16, readCount uint16) ([]uint16, error)
	ReadString(memoryArea byte, address uint16, readCount uint16) (*string, error)
	ReadBytes(memoryArea byte, address uint16, readCount uint16) ([]byte, error)
	WriteBites(memoryArea byte, address uint16, bitOffset byte, data []bool) error
	WriteWords(memoryArea byte, address uint16, data []uint16) error
	WriteString(memoryArea byte, address uint16, itemCount uint16, s string) error
}
