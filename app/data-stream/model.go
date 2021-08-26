package data_stream

var (
	_ Buffer = DataBuffer{}
	_ Buffer = NewDataBuffer{}
	_ Buffer = EOFBuffer{}
	_ Buffer = CloseBuffer{}
)

const (
	Data Command = iota
	NewData
	EOF
	CloseData
)

type Command int

type Buffer interface {
	Bytes() []byte
	Command() Command
}

func NewBuffer(cmd Command, buf []byte) Buffer {
	switch cmd {
	case NewData:
		return NewDataBuffer(buf)
	case EOF:
		return EOFBuffer(buf)
	case CloseData:
		return CloseBuffer(buf)
	default:
		return DataBuffer(buf)
	}
}

type DataBuffer []byte

func (d DataBuffer) Bytes() []byte {
	return d
}

func (d DataBuffer) Command() Command {
	return Data
}

type NewDataBuffer []byte

func (d NewDataBuffer) Bytes() []byte {
	return d
}

func (d NewDataBuffer) Command() Command {
	return NewData
}

type EOFBuffer []byte

func (d EOFBuffer) Bytes() []byte {
	return d
}

func (d EOFBuffer) Command() Command {
	return EOF
}

type CloseBuffer []byte

func (d CloseBuffer) Bytes() []byte {
	return d
}

func (d CloseBuffer) Command() Command {
	return CloseData
}
