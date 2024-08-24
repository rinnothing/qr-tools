package qr_tools

import (
	"bytes"
	cryptoRand "crypto/rand"
	"errors"
	"math"
	"math/rand"
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
	testCapacity(t, codewordsCapacities, "codewords")
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

	var arr [100]byte
	for l := 1; l < 100; l++ {
		_, err := cryptoRand.Read(arr[:l-1])
		if err != nil {
			t.Fatalf("Some strange error %s", err.Error())
		}
		arr[l] = 'W'

		s := string(arr[:l])
		if isNumeric(s) {
			t.Errorf("Not numeric string %s is recognized as one", s)
		}
	}
}

func TestGetAlphaNumericNumber(t *testing.T) {
	anSymbols := make(map[int32]int)
	for ch := '0'; ch <= '9'; ch++ {
		anSymbols[ch] = int(ch - '0')
	}
	for ch := 'A'; ch <= 'Z'; ch++ {
		anSymbols[ch] = int(ch-'A') + 10
	}
	for ch, val := range excessAlphanumerics {
		anSymbols[ch] = val
	}

	for i := int32(0); i < 256; i++ {
		correct, exists := anSymbols[i]
		answer := getAlphanumericNumber(i)
		if !((!exists && answer == -1) || (correct == answer)) {
			t.Errorf("Get alphanumeric number %c answered %d instead of %c", i, answer, anSymbols[i])
		}
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	symbols := make([]int32, 0, 45)
	for ch := '0'; ch <= '9'; ch++ {
		symbols = append(symbols, ch)
	}
	for ch := 'A'; ch <= 'Z'; ch++ {
		symbols = append(symbols, ch)
	}
	for ch := range excessAlphanumerics {
		symbols = append(symbols, ch)
	}

	{
		var arr [100]int32
		for l := 1; l < 100; l++ {
			for i := 0; i < l; i++ {
				arr[i] = symbols[rand.Intn(len(symbols))]
			}

			s := string(arr[:l])
			if !isAlphaNumeric(s) {
				t.Errorf("Alphanumeric string %s is not recognized as one", s)
			}
		}
	}

	{
		var arr [100]byte
		for l := 1; l < 100; l++ {
			if _, err := cryptoRand.Read(arr[:l-1]); err != nil {
				t.Fatalf("Some strange error %s", err.Error())
			}

			s := string(arr[:l])
			if isAlphaNumeric(s) {
				t.Errorf("Non alphanumeric string %s is recognized as one", s)
			}
		}
	}
}

func TestClearLeadingZeroes(t *testing.T) {
	for i := 0; i < 100; i++ {
		sb := strings.Builder{}
		for j := 0; j < i; j++ {
			sb.WriteByte('0')
		}

		s := strconv.Itoa(rand.Int())
		sb.WriteString(s)

		sWithZeroes := sb.String()
		sWithoutZeroes := clearLeadingZeroes(sWithZeroes)
		if sWithoutZeroes != s {
			t.Errorf("Wrongly cleared leading zeroes in string %s, %s should be equal %s", sWithZeroes, sWithoutZeroes, s)
		}
	}
}

func TestAddCharacterCount(t *testing.T) {
	for lvl := QRVersion(0); lvl < 100; lvl++ {
		chCnt := rand.Int() % 100

		chCntBa := newBitsetAppender()
		err := addCharacterCount([3]uint{10, 12, 14}, chCntBa, lvl, chCnt)

		ba := newBitsetAppender()

		var bitsNum uint
		switch {
		case lvl >= 1 && lvl <= 9:
			bitsNum = 10
		case lvl >= 10 && lvl <= 26:
			bitsNum = 12
		case lvl >= 27 && lvl <= 40:
			bitsNum = 14
		default:
			if err == nil {
				t.Errorf("Lvl %d should give an error", lvl)
			}
			continue
		}
		ba.appendUint16(uint16(chCnt)<<(16-bitsNum), bitsNum)

		dataA := ba.getData()
		dataM := chCntBa.getData()

		for i := 0; i < len(dataA)-1; i++ {
			if dataA[i] != dataM[i] {
				t.Errorf("Byte %d doesn't match: %d != %d", i, dataA[i], dataM[i+1])
			}
		}

		lastByteBits := ba.n % 8
		if lastByteBits == 0 {
			lastByteBits = 8
		}
		lastByteM := dataM[len(dataA)-1] & (allOnes << (8 - lastByteBits))
		lastByteA := dataA[len(dataA)-1] & (allOnes << (8 - lastByteBits))
		if lastByteA != lastByteM {
			t.Errorf("Last byte doesn't match: %d != %d", lastByteA, lastByteM)
		}
	}

}

func TestAddPadding(t *testing.T) {
	data := make([]byte, 100)
	if _, err := cryptoRand.Read(data[:]); err != nil {
		t.Errorf("Failed to generate random data %v", err)
	}

	const l = 104
	for num := uint(1); num < l; num++ {
		ba := newBitsetAppender()
		_ = ba.append(data, num)

		addPadding(ba, l)

		if ba.n != l {
			t.Errorf("Not all space is filled")
		}

		baCopy := newBitsetAppender()
		_ = baCopy.append(data, num)

		if baCopy.n < l-4 {
			baCopy.appendByte(0, 4)
		} else {
			baCopy.appendByte(0, l-baCopy.n)
		}

		for baCopy.n%8 != 0 {
			baCopy.appendByte(0, 1)
		}

		flag := true
		for baCopy.n < l {
			if flag {
				baCopy.appendByte(236, 8)
			} else {
				baCopy.appendByte(17, 8)
			}
			flag = !flag
		}

		if !bytes.Equal(ba.data, baCopy.data) {
			t.Errorf("Incorrect padding: %x != %x", ba.data, baCopy.data)
		}
	}
}

func TestNewNumericMarshaler(t *testing.T) {
	var lvl ErrorCorrectionLevel = H
	var ver QRVersion = 10

	nm := NewNumericMarshaler(lvl, ver)
	if nm.lvl != lvl || ver != nm.ver {
		t.Errorf("NewNumericMarshaller's arguments are wrong")
	}
}

func TestNumericMarshaler_MarshalString(t *testing.T) {
	{
		s := "abcde"
		nm := NewNumericMarshaler(L, 1)
		if _, err := nm.MarshalString(s); !errors.Is(err, wrongFormatError) {
			t.Errorf("%s is not numeric, but error doesn't appear", s)
		}
	}

	lvl := ErrorCorrectionLevel(L)
	for _, ver := range []QRVersion{1, 45} {
		s := "8675309"
		nm := NewNumericMarshaler(lvl, ver)
		data, err := nm.MarshalString(s)
		if err != nil {
			if ver <= 40 && ver >= 1 {
				t.Errorf("Level is %d, but error arises", ver)
			}
			continue
		}

		var sMarshaled []byte
		{
			ba := newBitsetAppender()
			ba.appendByte(0b0001<<4, 4)
			_ = addCharacterCount(numericBitCounts, ba, ver, len(s))
			ba.appendUint16(0b1101100011<<6, 10)
			ba.appendUint16(0b1000010010<<6, 10)
			ba.appendUint16(0b1001<<12, 4)
			addPadding(ba, codewordsCapacities[lvl][ver-1]*8)

			sMarshaled = ba.getData()
		}

		if !bytes.Equal(data, sMarshaled) {
			t.Errorf("Data isn't marshaled properly")
			t.Log(sMarshaled)
			t.Log(data)
		}
	}
}

func TestNewAlphanumericMarshaler(t *testing.T) {
	var lvl ErrorCorrectionLevel = H
	var ver QRVersion = 10

	nm := NewAlphanumericMarshaler(lvl, ver)
	if nm.lvl != lvl || ver != nm.ver {
		t.Errorf("NewAlphanumericMarshaller's arguments are wrong")
	}
}

func TestAlphanumericMarshaler_MarshalString(t *testing.T) {
	{
		s := "abcde!!!"
		nm := NewAlphanumericMarshaler(L, 1)
		if _, err := nm.MarshalString(s); !errors.Is(err, wrongFormatError) {
			t.Errorf("%s is not alphanumeric, but error doesn't appear", s)
		}
	}

	lvl := ErrorCorrectionLevel(L)
	for _, ver := range []QRVersion{1, 45} {
		s := "HELLO WORLD"
		nm := NewAlphanumericMarshaler(lvl, ver)
		data, err := nm.MarshalString(s)
		if err != nil {
			if ver <= 40 && ver >= 1 {
				t.Errorf("Level is %d, but error arises", ver)
				t.Log(err)
			}
			continue
		}

		var sMarshaled []byte
		{
			sSplitted := []string{"HE", "LL", "O ", "WO", "RL", "D"}
			ba := newBitsetAppender()
			ba.appendByte(0b0010<<4, 4)
			_ = addCharacterCount(alphanumericBitCounts, ba, ver, len(s))
			for _, pair := range sSplitted[:len(sSplitted)-1] {
				num := getAlphanumericNumber(int32(pair[0]))*45 + getAlphanumericNumber(int32(pair[1]))
				ba.appendUint16(uint16(num)<<5, 11)
			}
			ba.appendByte(uint8(getAlphanumericNumber(int32(sSplitted[len(sSplitted)-1][0])))<<2, 6)
			addPadding(ba, codewordsCapacities[lvl][ver-1]*8)

			sMarshaled = ba.getData()
		}

		if !bytes.Equal(data, sMarshaled) {
			t.Errorf("Data isn't marshaled properly")
			t.Log(sMarshaled)
			t.Log(data)
		}
	}
}

func TestNewByteMarshaler(t *testing.T) {
	var lvl ErrorCorrectionLevel = H
	var ver QRVersion = 10

	nm := NewByteMarshaler(lvl, ver)
	if nm.lvl != lvl || ver != nm.ver {
		t.Errorf("NewByteMarshaller's arguments are wrong")
	}
}

func TestByteMarshaler_MarshalString(t *testing.T) {
	lvl := ErrorCorrectionLevel(L)
	for _, ver := range []QRVersion{1, 45} {
		s := "Hello, world!"
		nm := NewByteMarshaler(lvl, ver)
		data, err := nm.MarshalString(s)
		if err != nil {
			if ver <= 40 && ver >= 1 {
				t.Errorf("Level is %d, but error arises", ver)
				t.Log(err)
			}
			continue
		}

		var sMarshaled []byte
		{
			ba := newBitsetAppender()
			ba.appendByte(0b0100<<4, 4)
			_ = addCharacterCount(byteBitCounts, ba, ver, len(s))
			for _, ch := range s {
				ba.appendByte(byte(ch), 8)
			}
			addPadding(ba, codewordsCapacities[lvl][ver-1]*8)

			sMarshaled = ba.getData()
		}

		if !bytes.Equal(data, sMarshaled) {
			t.Errorf("Data isn't marshaled properly")
			t.Log(sMarshaled)
			t.Log(data)
		}
	}
}

func TestNewQRMarshaler(t *testing.T) {
	var lvl ErrorCorrectionLevel = H
	var ver QRVersion = 10

	nm := NewQRMarshaler(lvl, ver)
	if nm.lvl != lvl || ver != nm.ver {
		t.Errorf("NewQRMarshaller's arguments are wrong")
	}
}

func TestQRMarshaler_MarshalString(t *testing.T) {
	lvl := ErrorCorrectionLevel(L)
	for _, ver := range []QRVersion{1, 45} {
		qr := NewQRMarshaler(lvl, ver)

		s := "12345"
		nm := NewNumericMarshaler(lvl, ver)
		dataNm, errNm := nm.MarshalString(s)
		dataQr, errQr := qr.MarshalString(s)
		if errNm == nil && !bytes.Equal(dataNm, dataQr) {
			t.Errorf("QRMarshaller didn't use NumericMarshaller to encode %s", s)
		} else if errNm != nil && errNm != errQr {
			t.Errorf("Different errors arose: %s != %s", errNm, errQr)
		}

		s = "BABA 1234.++"
		anm := NewAlphanumericMarshaler(lvl, ver)
		dataAnm, errAnm := anm.MarshalString(s)
		dataQr, errQr = qr.MarshalString(s)
		if errAnm == nil && !bytes.Equal(dataAnm, dataQr) {
			t.Errorf("QRMarshaller didn't use AlphanumericMarshaller to encode %s", s)
		} else if errAnm != nil && errAnm != errQr {
			t.Errorf("Different errors arose: %s != %s", errAnm, errQr)
		}

		s = "MoNeY!!!$$$"
		bm := NewByteMarshaler(lvl, ver)
		dataBm, errBm := bm.MarshalString(s)
		dataQr, errQr = qr.MarshalString(s)
		if errBm == nil && !bytes.Equal(dataBm, dataQr) {
			t.Errorf("QRMarshaller didn't use ByteMarshaller to encode %s", s)
		} else if errBm != nil && errBm != errQr {
			t.Errorf("Different errors arose: %s != %s", errBm, errQr)
		}

	}
}
