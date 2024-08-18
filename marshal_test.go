package qr_tools

import (
	"math"
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
	for i := uint(1); i <= 100; i++ {
		ba.appendByte(byte(i)<<(8-min(i, 8)), i)

		prLn := ln
		ln += min(i, 8)
		if ln != ba.n {
			t.Errorf("Bitset appender length %d does not match actual number of bits appended %d", ba.n, ln)
		}

		var extr byte
		if ba.n/8 == prLn/8 {
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
