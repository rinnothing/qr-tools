package qr_tools

import (
	"errors"
	"math"
	"strconv"
	"unicode"
)

const (
	// just ones for making bit masks
	allOnes = byte(math.MaxUint8)
)

var (
	bitBeyondError      = errors.New("bit number goes beyond number of bytes")
	wrongFormatError    = errors.New("content format is not suitable for marshaler")
	wrongQRVersionError = errors.New("wrong qr version")

	// just some magic numbers, hope later I'll calculate them by myself (taken from https://www.thonky.com/qr-code-tutorial/character-capacities)
	numericCapacities = [][]uint{
		{41, 77, 127, 187, 255, 322, 370, 461, 552, 652, 772, 883, 1022, 1101, 1250, 1408, 1548, 1725, 1903, 2061, 2232, 2409, 2620, 2812, 3057, 3283, 3517, 3669, 3909, 4158, 4417, 4686, 4965, 5253, 5529, 5836, 6153, 6479, 6743, 7089}, // L
		{34, 63, 101, 149, 202, 255, 293, 365, 432, 513, 604, 691, 796, 871, 991, 1082, 1212, 1346, 1500, 1600, 1708, 1872, 2059, 2188, 2395, 2544, 2701, 2857, 3035, 3289, 3486, 3693, 3909, 4134, 4343, 4588, 4775, 5039, 5313, 5596},    // M
		{27, 48, 77, 111, 144, 178, 207, 259, 312, 364, 427, 489, 580, 621, 703, 775, 876, 948, 1063, 1159, 1224, 1358, 1468, 1588, 1718, 1804, 1933, 2085, 2181, 2358, 2473, 2670, 2805, 2949, 3081, 3244, 3417, 3599, 3791, 3993},        // Q
		{17, 34, 58, 82, 106, 139, 154, 202, 235, 288, 331, 374, 427, 468, 530, 602, 674, 746, 813, 919, 969, 1056, 1108, 1228, 1286, 1425, 1501, 1581, 1677, 1782, 1897, 2022, 2157, 2301, 2361, 2524, 2625, 2735, 2927, 3057},            // H
	}
	alphanumericCapacities = [][]uint{
		{25, 47, 77, 114, 154, 195, 224, 279, 335, 395, 468, 535, 619, 667, 758, 854, 938, 1046, 1153, 1249, 1352, 1460, 1588, 1704, 1853, 1990, 2132, 2223, 2369, 2520, 2677, 2840, 3009, 3183, 3351, 3537, 3729, 3927, 4087, 4296}, // L
		{20, 38, 61, 90, 122, 154, 178, 221, 262, 311, 366, 419, 483, 528, 600, 656, 734, 816, 909, 970, 1035, 1134, 1248, 1326, 1451, 1542, 1637, 1732, 1839, 1994, 2113, 2238, 2369, 2506, 2632, 2780, 2894, 3054, 3220, 3391},     // M
		{16, 29, 47, 67, 87, 108, 125, 157, 189, 221, 259, 296, 352, 376, 426, 470, 531, 574, 644, 702, 742, 823, 890, 963, 1041, 1094, 1172, 1263, 1322, 1429, 1499, 1618, 1700, 1787, 1867, 1966, 2071, 2181, 2298, 2420},          // Q
		{10, 20, 35, 50, 64, 84, 93, 122, 143, 174, 200, 227, 259, 283, 321, 365, 408, 452, 493, 557, 587, 640, 672, 744, 779, 864, 910, 958, 1016, 1080, 1150, 1226, 1307, 1394, 1431, 1530, 1591, 1658, 1774, 1852},                // H
	}
	byteCapacities = [][]uint{
		{17, 32, 53, 78, 106, 134, 154, 192, 230, 271, 321, 367, 425, 458, 520, 586, 644, 718, 792, 858, 929, 1003, 1091, 1171, 1273, 1367, 1465, 1528, 1628, 1732, 1840, 1952, 2068, 2188, 2303, 2431, 2563, 2699, 2809, 2953}, // L
		{14, 26, 42, 62, 84, 106, 122, 152, 180, 213, 251, 287, 331, 362, 412, 450, 504, 560, 624, 666, 711, 779, 857, 911, 997, 1059, 1125, 1190, 1264, 1370, 1452, 1538, 1628, 1722, 1809, 1911, 1989, 2099, 2213, 2331},      // M
		{11, 20, 32, 46, 60, 74, 86, 108, 130, 151, 177, 203, 241, 258, 292, 322, 364, 394, 442, 482, 509, 565, 611, 661, 715, 751, 805, 868, 908, 982, 1030, 1112, 1168, 1228, 1283, 1351, 1423, 1499, 1579, 1663},             // Q
		{7, 14, 24, 34, 44, 58, 64, 84, 98, 119, 137, 155, 177, 194, 220, 250, 280, 310, 338, 382, 403, 439, 461, 511, 535, 593, 625, 658, 698, 742, 790, 842, 898, 958, 983, 1051, 1093, 1139, 1219, 1273},                     // H
	}
	kanjiCapacities = [][]uint{
		{10, 20, 32, 48, 65, 82, 95, 118, 141, 167, 198, 226, 262, 282, 320, 361, 397, 442, 488, 528, 572, 618, 672, 721, 784, 842, 902, 940, 1002, 1066, 1132, 1201, 1273, 1347, 1417, 1496, 1577, 1661, 1729, 1817}, // L
		{8, 16, 26, 38, 52, 65, 75, 93, 111, 131, 155, 177, 204, 223, 254, 277, 310, 345, 384, 410, 438, 480, 528, 561, 614, 652, 692, 732, 778, 843, 894, 947, 1002, 1060, 1113, 1176, 1224, 1292, 1362, 1435},       // M
		{7, 12, 20, 28, 37, 45, 53, 66, 80, 93, 109, 125, 149, 159, 180, 198, 224, 243, 272, 297, 314, 348, 376, 407, 440, 462, 496, 534, 559, 604, 634, 684, 719, 756, 790, 832, 876, 923, 972, 1024},                // Q
		{4, 8, 15, 21, 27, 36, 39, 52, 60, 74, 85, 96, 109, 120, 136, 154, 173, 191, 208, 235, 248, 270, 284, 315, 330, 365, 385, 405, 430, 457, 486, 518, 553, 590, 605, 647, 673, 701, 750, 784},                    // H
	}

	//additional alphanumeric chars which codes are too strange to make ifs for them
	excessAlphanumerics = map[int32]int{' ': 36, '$': 37, '%': 38, '*': 39, '+': 40, '-': 41, '.': 42, '/': 43, ':': 44}
)

// ErrorCorrectionLevel is enum that
// shows some unmarshalers and marshalers how strictly to encode data
type ErrorCorrectionLevel uint

const (
	L = iota
	M
	Q
	H
)

// QRVersion is enum that
// shows what version of QR code we're using
type QRVersion uint

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

// NewNumericMarshaler returns NumericMarshaler
// with chosen ErrorCorrectionLevel
func NewNumericMarshaler(lvl ErrorCorrectionLevel, ver QRVersion) *NumericMarshaler {
	return &NumericMarshaler{lvl: lvl, ver: ver}
}

func isNumeric(str string) bool {
	for _, ch := range str {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func clearLeadingZeroes(str string) string {
	var i int
	for i = 0; i < len(str)-1; i++ {
		if str[i] != '0' {
			break
		}
	}

	return str[i:]
}

// MarshalString marshals the given numeric string
func (nm *NumericMarshaler) MarshalString(str string) ([]byte, error) {
	if !isNumeric(str) {
		return nil, wrongFormatError
	}

	ba := newBitsetAppender()

	//adding mode indicator - numeric
	ba.appendByte(0b00010000, 4)
	//adding size indicator
	if err := addCharacterCount([3]uint{10, 12, 14}, ba, nm.ver, len(str)); err != nil {
		return nil, err
	}

	// splitting digits into triplets and cutting away leading zeroes
	// the size of encoded depends on number of digits
	for i := 0; i < len(str); i += 3 {
		piece := str[i:min(len(str), i+3)]
		piece = clearLeadingZeroes(piece)

		bits := uint(1 + 3*len(piece))

		pint, _ := strconv.Atoi(piece)
		ba.appendUint16(uint16(pint)<<(16-bits), bits)
	}

	// padding information
	bitsNum := numericCapacities[nm.lvl][nm.ver-1] * 8
	addPadding(ba, bitsNum)

	return ba.getData(), nil
}

// An AlphanumericMarshaler can marshal alphanumeric data (0-9, A-Z,' ', S$, %, *, +, -, ., /, :)
// with respect to ErrorCorrectionLevel
type AlphanumericMarshaler struct {
	lvl ErrorCorrectionLevel
	ver QRVersion
}

// NewAlphanumericMarshaler returns AlphanumericMarshaller
// with chosen ErrorCorrectionLevel
func NewAlphanumericMarshaler(lvl ErrorCorrectionLevel, ver QRVersion) *AlphanumericMarshaler {
	return &AlphanumericMarshaler{lvl: lvl, ver: ver}
}

// returns the char equivalent in alphanumeric encoding (or -1 in case char isn't alphanumeric)
func getAlphanumericNumber(ch int32) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	} else if ch >= 'A' && ch <= 'Z' {
		return int(ch - 'A' + 10)
	} else if n, is := excessAlphanumerics[ch]; is {
		return n
	}

	return -1
}

func isAlphaNumeric(str string) bool {
	for _, ch := range str {
		if getAlphanumericNumber(ch) == -1 {
			return false
		}
	}

	return true
}

// MarshalString marshals the given alphanumeric string
func (am *AlphanumericMarshaler) MarshalString(str string) ([]byte, error) {
	if !isAlphaNumeric(str) {
		return nil, wrongFormatError
	}

	ba := newBitsetAppender()

	//adding mode indicator - alphanumeric
	ba.appendByte(0b00100000, 4)
	//adding size indicator
	if err := addCharacterCount([3]uint{9, 11, 13}, ba, am.ver, len(str)); err != nil {
		return nil, err
	}

	//splitting string in duos and encoding
	i := 0
	for ; i+2 <= len(str); i += 2 {
		nm := uint16(getAlphanumericNumber(int32(str[i]))*45 + getAlphanumericNumber(int32(str[i+1])))

		ba.appendUint16(nm<<5, 11)
	}
	if i == len(str)-1 {
		nm := uint16(getAlphanumericNumber(int32(str[i])))

		ba.appendUint16(nm<<10, 6)
	}

	//applying padding
	bitsNum := alphanumericCapacities[am.lvl][am.ver-1] * 8
	addPadding(ba, bitsNum)

	return ba.getData(), nil
}

// A ByteMarshaler can marshal byte data
// with respect to ErrorCorrectionLevel
type ByteMarshaler struct {
	lvl ErrorCorrectionLevel
	ver QRVersion
}

// NewByteMarshaler returns ByteMarshaler
// with chosen ErrorCorrectionLevel
func NewByteMarshaler(lvl ErrorCorrectionLevel, ver QRVersion) *ByteMarshaler {
	return &ByteMarshaler{lvl: lvl, ver: ver}
}

// MarshalString marshals the given byte string
func (bm *ByteMarshaler) MarshalString(str string) ([]byte, error) {
	ba := newBitsetAppender()

	//adding mode indicator - byte
	ba.appendByte(0b01000000, 4)
	//adding size indicator
	if err := addCharacterCount([3]uint{8, 16, 16}, ba, bm.ver, len(str)); err != nil {
		return nil, err
	}

	//just past data (no way there would be an error)
	_ = ba.append([]byte(str), uint(len(str)*8))

	//padding information
	bitsNum := byteCapacities[bm.lvl][bm.ver-1] * 8
	addPadding(ba, bitsNum)

	return ba.getData(), nil
}

//will add kanji later

// addCharacterCount appends character count indicator to string
// throws wrongQRVersionError
func addCharacterCount(bitCounts [3]uint, ba *bitsetAppender, ver QRVersion, chCnt int) error {
	var cntSize uint
	switch {
	case ver >= 1 && ver <= 9:
		cntSize = bitCounts[0]
	case ver >= 10 && ver <= 26:
		cntSize = bitCounts[1]
	case ver >= 27 && ver <= 40:
		cntSize = bitCounts[2]
	default:
		return wrongQRVersionError
	}

	ba.appendUint16(uint16(chCnt<<(16-cntSize)), cntSize)
	return nil
}

// addPadding adds padding after placing information
// used in some marshalers
func addPadding(ba *bitsetAppender, bitsNum uint) {
	// adding 0-terminator
	ba.appendUint16(0, min(bitsNum-ba.n, 4))

	// adding 0s to make a multiple of 8
	ba.appendUint16(0, (8-ba.n%8)%8)

	// adding pad bytes to fill capacity to the end
	for ba.n < bitsNum {
		ba.appendUint16(0b11101100_00010001, min(16, bitsNum-ba.n))
	}
}

// A QRMarshaler can marshal numeric, alphanumeric, byte and kanji (later) effectively
// with respect to ErrorCor
type QRMarshaler struct {
	lvl ErrorCorrectionLevel
	ver QRVersion
}

// NewQRMarshaler returns QRMarshaler
// with chosen ErrorCorrectionLevel
func NewQRMarshaler(lvl ErrorCorrectionLevel, ver QRVersion) *QRMarshaler {
	return &QRMarshaler{lvl: lvl, ver: ver}
}

// MarshalString marshals the given string effectively
func (qm *QRMarshaler) MarshalString(str string) ([]byte, error) {
	var mller Marshaler = NewNumericMarshaler(qm.lvl, qm.ver)
	if data, err := mller.MarshalString(str); err == nil {
		return data, nil
	}

	mller = NewAlphanumericMarshaler(qm.lvl, qm.ver)
	if data, err := mller.MarshalString(str); err == nil {
		return data, nil
	}

	mller = NewByteMarshaler(qm.lvl, qm.ver)
	if data, err := mller.MarshalString(str); err == nil {
		return data, nil
	}

	return nil, wrongFormatError
}

// Unmarshaler is the interface implemented by types that
// can unmarshal a string from a sequence of bytes
type Unmarshaler interface {
	UnmarshalToString(data []byte) (str string, err error)
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
	if n == 0 {
		return
	}
	oldN := n
	n = min(8, n)

	if ba.n%8 != 0 {
		ba.data[ba.n/8] |= data >> (ba.n % 8)
		ba.data[ba.n/8] &= allOnes << (8 - min(ba.n%8+n, 8))
	}
	if (8-ba.n%8)%8 < n {
		ba.data = append(ba.data, data<<((8-ba.n%8)%8))
	}

	if oldN > 8 {
		ba.n += 8
	} else {
		ba.n += n
	}
}

// append appends n bits to the sequence that is already in bitsetAppender
// note that n shall not be greater than len(date) * 8, otherwise you'll get a bitBeyondError
// but previous data would be written
func (ba *bitsetAppender) append(data []byte, n uint) error {
	if n == 0 {
		return nil
	}

	times := (n-1)/8 + 1

	for i := 0; times > 0; times, n, i = times-1, n-8, i+1 {
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
	if n == 0 {
		return
	}

	n = min(n, 16)
	times := (n-1)/8 + 1

	m := 8
	for ; times > 0; times, n, m = times-1, n-8, m-8 {
		ba.appendByte(byte(data>>m), n)
	}
}

// getData gives out written data, cutting all not important stuff away
func (ba *bitsetAppender) getData() []byte {
	inc := uint(0)
	if ba.n%8 != 0 {
		inc = 1
	}
	return ba.data[:ba.n/8+inc]
}
