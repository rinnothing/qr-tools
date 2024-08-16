package qr_tools

import "errors"

var (
	bitBeyondError = errors.New("bit number goes beyond number of bytes")
)

// ErrorCorrectionLevel is enum that
// shows some unmarshalers and marshalers how strictly to encode data
type ErrorCorrectionLevel int

const (
	L = iota
	M
	Q
	H
)

// QRVersion is enum that
// shows what version of QR code we're using
type QRVersion int

// Marshaler is the interface implemented by types that
// can marshal a string into a sequence of bytes.
type Marshaler interface {
	MarshalString(str string) (data []byte, err error)
}

// A NumericMarshaler can marshal numeric data
// with respect to ErrorCorrectionLevel
type NumericMarshaler struct {
	lvl ErrorCorrectionLevel
	ver QRVersion
}

// NewNumericMarshaler returns NumericMarshaller
// with chosen ErrorCorrectionLevel
func NewNumericMarshaler(lvl ErrorCorrectionLevel, ver QRVersion) *NumericMarshaler {
	return &NumericMarshaler{lvl: lvl, ver: ver}
}

// MarshalString marshals the given numeric string
func (nm *NumericMarshaler) MarshalString(str string) ([]byte, error) {
	ba := newBitsetAppender()
	ba.appendByte(0b00010000, 4)

	var cntSize uint
	switch {
	case nm.ver >= 1 && nm.ver <= 9:
		cntSize = 10
	case nm.ver >= 10 && nm.ver <= 26:
		cntSize = 12
	case nm.ver >= 27 && nm.ver <= 40:
		cntSize = 14
	}
	ba.appendUint16(uint16(len(str))<<(16-cntSize), cntSize)

	//to be continued...

	return nil, nil
}

// Unmarshaler is the interface implemented by types that
// can unmarshal a string from a sequence of bytes
type Unmarshaler interface {
	UnmarshalToString(data []byte) (str string, err error)
}

type outOfBoundsError struct {
	bytes, bits int
}

// bitsetAppender is a structure that
// helps to connect long sequences of bits
//
// note that n should always be not greater than len(data) * 8
type bitsetAppender struct {
	data []byte
	n    uint
}

// newBitsetAppender creates empty bitsetAppender
func newBitsetAppender() *bitsetAppender {
	return &bitsetAppender{data: make([]byte, 0), n: 0}
}

// appendByte appends n bits to the sequence of bits inside bitAppender
// if n is greater than 8 appendByte append just 8 bits
func (ba *bitsetAppender) appendByte(data byte, n uint) {
	if ba.n%8 != 0 {
		ba.data[ba.n/8] |= data >> ba.n % 8
	}
	ba.data = append(ba.data, data<<(8-ba.n%8)%8)

	if n > 8 {
		ba.n += 8
	}
	ba.n += n
}

// append appends n bits to the sequence that is already in bitsetAppender
// note that n shall not be greater than len(date) * 8, otherwise you'll get a bitBeyondError
// but previous data would be written
func (ba *bitsetAppender) append(data []byte, n uint) error {
	for i := 0; n > 0; n, i = n-8, i+1 {
		if i >= len(data) {
			return bitBeyondError
		}

		ba.appendByte(data[i], n)
	}

	return nil
}

// appendUint16 simply splits uint into 4 bytes and appends them
// if n is greater than 16 appendUint16 appends just 16 bits
func (ba *bitsetAppender) appendUint16(data uint16, n uint) {
	n = max(n, 16)

	m := 8
	for ; n > 0; n, m = n-8, m-8 {
		ba.appendByte(byte(data>>m), n)
	}
}
