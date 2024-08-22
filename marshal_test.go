package qr_tools

import (
	"bytes"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

func testCapacity(t *testing.T, cap [][]uint, name string) {
	if len(cap) != 4 {
		t.Errorf("Capacity %s supports only %d error correction modes from 4", name, len(cap))
	}

	for i := 0; i < len(cap); i++ {
		if len(cap[i]) != 40 {
			t.Errorf("Capacity %s's error correction level %d supports only %d qr versions from 40", name, i, len(cap[i]))
		}
	}
}

func TestCapacities(t *testing.T) {
	testCapacity(t, numericCapacities, "numeric")
	testCapacity(t, alphanumericCapacities, "alphanumeric")
	testCapacity(t, byteCapacities, "byte")
	testCapacity(t, kanjiCapacities, "kanji")
}

func TestNewBitsetAppender(t *testing.T) {
	ba := newBitsetAppender()

	if ba.n != 0 || len(ba.data) != 0 {
		t.Errorf("Newly created bitset appender is not empty")
	}
}

func TestBitsetAppendByte(t *testing.T) {
	allOnes := byte(math.MaxUint8)
	ba := newBitsetAppender()

	ln := uint(0)
	for i := uint(0); i <= 100; i++ {
		ba.appendByte(byte(i)<<(8-min(i, 8)), i)

		prLn := ln
		ln += min(i, 8)
		if ln != ba.n {
			t.Errorf("Bitset appender length %d does not match actual number of bits appended %d", ba.n, ln)
		}

		var extr byte
		if i == 0 {
			extr = 0
		} else if ba.n/8 == prLn/8 {
			extr = ba.data[len(ba.data)-1] >> ((8 - ba.n%8) % 8)
			extr &= allOnes >> (8 - (ba.n - prLn))
		} else {
			extr = ba.data[len(ba.data)-1] >> ((8 - ba.n%8) % 8)
			if ba.n%8 != 0 {
				extr |= ba.data[len(ba.data)-2] << (ba.n % 8)
			}
			extr &= allOnes >> (8 - (ba.n - prLn))
		}

		actual := byte(i) & (allOnes >> (8 - min(i, 8)))
		if extr != actual {
			t.Errorf("Extracted bites %b aren't equal with actual bites %b, i = %d", extr, actual, i)
		}
	}
}

func TestBitsetAppendUint16(t *testing.T) {
	ba1 := newBitsetAppender()
	ba2 := newBitsetAppender()

	for i, num := 0, uint16(math.MaxUint8+1); i < 100; i, num = i+1, num+1 {
		bits := i % 25

		ba1.appendUint16(num, uint(bits))

		ba2.appendByte(byte(num>>8), uint(min(bits, 8)))
		ba2.appendByte(byte(num), uint(max(bits-8, 0)))

		if ba1.n != ba2.n || !bytes.Equal(ba1.data, ba2.data) {
			t.Errorf("Arrays aren't equal after adding %d bits from %d", bits, i)
		}
	}
}

func TestBitsetAppend(t *testing.T) {
	ba1 := newBitsetAppender()
	ba2 := newBitsetAppender()

	data := make([]byte, 0, 12)
	for i := 0; i < 12; i += 4 {
		r := rand.Uint32()
		data = append(data, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
	}

	for i := 0; i <= 12*8+100; i++ {
		err := ba1.append(data, uint(i))
		if err != nil && i <= 12*8 {
			t.Errorf("Got sudden error after adding %d bits", i)
		} else if err == nil && i > 12*8 {
			t.Errorf("Didn't get error after adding %d bits", i)
		}

		for j, y := min(i, 12*8), 0; j > 0; j, y = j-8, y+1 {
			ba2.appendByte(data[y], uint(min(j, 8)))
		}

		if ba1.n != ba2.n || ba1.n > 0 && !bytes.Equal(ba1.data, ba2.data) {
			t.Errorf("Arrays aren't equal after adding %d bits", i)
		}
	}
}

func TestBitsetGetData(t *testing.T) {
	ba := newBitsetAppender()
	data := make([]byte, 0, 12)
	for i := 0; i < 12; i += 4 {
		r := rand.Uint32()
		data = append(data, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
	}

	bits := uint(12*5 + 6)

	_ = ba.append(data, bits)
	if ba.n != bits || !bytes.Equal(ba.getData(), data[:(bits-1)/8+1]) {
		t.Errorf("getData's value and data aren't equal")
	}
}

func TestIsNumeric(t *testing.T) {
	for l := 1; l < 100; l++ {
		sb := strings.Builder{}
		var w string
		for w = strconv.Itoa(int(rand.Uint32())); sb.Len()+len(w) <= l; {
			sb.WriteString(w)
		}
		sb.WriteString(w[:l-sb.Len()])

		s := sb.String()
		if !isNumeric(s) {
			t.Errorf("Numeric string %sb is not recognized as one", s)
		}
	}

	chacha8seed := [32]byte([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456"))
	r := rand.NewChaCha8(chacha8seed)

	var arr [100]byte
	for l := 1; l < 100; l++ {
		_, _ = r.Read(arr[:l-1])
		arr[l] = 'W'

		s := string(arr[:l])
		if isNumeric(s) {
			t.Errorf("Not numeric string %s is recognized as one", s)
		}
	}
}
