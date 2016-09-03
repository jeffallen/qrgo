package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSize(t *testing.T) {
	//t.Skip()
	s := size(1)
	assert.Equal(t, 21, s)

	s = size(2)
	assert.Equal(t, 25, s)

	s = size(32)
	assert.Equal(t, 145, s)
}

func TestMax(t *testing.T) {
	//t.Skip()
	m := max(1, 2)
	assert.Equal(t, 2, m)

	m = max(1, 1)
	assert.Equal(t, 1, m)

	m = max(17, 16)
	assert.Equal(t, 17, m)
}

func TestPadLeft(t *testing.T) {
	//t.Skip()
	pad := padLeft("", 4)
	assert.Equal(t, "0000", pad)

	pad = padLeft("1111", 2)
	assert.Equal(t, "1111", pad)

	pad = padLeft("1111", 4)
	assert.Equal(t, "1111", pad)

	pad = padLeft("111", 10)
	assert.Equal(t, "0000000111", pad)
}

func TestPadRight(t *testing.T) {
	//t.Skip()
	pad := padRight("", 4)
	assert.Equal(t, "0000", pad)

	pad = padRight("1111", 2)
	assert.Equal(t, "1111", pad)

	pad = padRight("1111", 4)
	assert.Equal(t, "1111", pad)

	pad = padRight("111", 10)
	assert.Equal(t, "1110000000", pad)
}

func TestStringToBinary(t *testing.T) {
	//t.Skip()
	bin := stringToBinary("")
	assert.Equal(t, "", bin)

	bin = stringToBinary("0")
	assert.Equal(t, "0", bin)

	bin = stringToBinary("1")
	assert.Equal(t, "1", bin)

	bin = stringToBinary("10")
	assert.Equal(t, "1010", bin)

	bin = stringToBinary("100")
	assert.Equal(t, "1100100", bin)

	bin = stringToBinary("255")
	assert.Equal(t, "11111111", bin)
}

func TestEncodingToByteArray(t *testing.T) {
	//t.Skip()
	array := encodingToByteArray("00000000")
	assert.Equal(t, []byte{0}, array)

	array = encodingToByteArray("0000001011111111")
	assert.Equal(t, []byte{2, 255}, array)

	array = encodingToByteArray("000000010000001000000100")
	assert.Equal(t, []byte{1, 2, 4}, array)
}

func TestByteArrayToEncoding(t *testing.T) {
	//t.Skip()
	en := byteArrayToEncoding([]byte{0})
	assert.Equal(t, "00000000", en)

	en = byteArrayToEncoding([]byte{2, 255})
	assert.Equal(t, "0000001011111111", en)

	en = byteArrayToEncoding([]byte{1, 2, 4})
	assert.Equal(t, "000000010000001000000100", en)
}

func TestCountSubstringOccurrences(t *testing.T) {
	//t.Skip()
	c := countSubstringOccurrences("ABCDEF", "CDE")
	assert.Equal(t, 1, c)

	c = countSubstringOccurrences("101011101", "101")
	assert.Equal(t, 3, c)

	c = countSubstringOccurrences("ABCEF", "GG")
	assert.Equal(t, 0, c)
}

func TestGetMode(t *testing.T) {
	//t.Skip()
	assert.Equal(t, 1, getMode("123"))
	assert.Equal(t, 2, getMode("ABC"))
	assert.Equal(t, 4, getMode("abc"))

	assert.Equal(t, 1, getMode("00000"))
	assert.Equal(t, 2, getMode("AB1C"))
	assert.Equal(t, 4, getMode("ABc"))

	assert.Equal(t, 1, getMode("001234"))
	assert.Equal(t, 2, getMode("1+B*C$/"))
	assert.Equal(t, 4, getMode("LOCAL@HOST.COM"))

	assert.NotEqual(t, 1, getMode("1111 1"))
	assert.NotEqual(t, 2, getMode("12AbC"))
	assert.NotEqual(t, 2, getMode("äöü"))
}

func TestGetVersion(t *testing.T) {
	//t.Skip()
	//v := getVersion(11, alpha)
	//assert.Equal(t, 2, v)

	//v = getVersion(1219, bytes)
	//assert.Equal(t, 39, v)

	//v = getVersion(2, numeric)
	//assert.Equal(t, 1, v)

	//v = getVersion(1853, alpha)
	//assert.Equal(t, 0, v)
}

func TestGetCountIndicator(t *testing.T) {
	//t.Skip()
	c := getCountIndicator(11, numeric, 2)
	assert.Equal(t, "0000001011", c)

	c = getCountIndicator(11, alpha, 2)
	assert.Equal(t, "000001011", c)

	c = getCountIndicator(4, bytes, 1)
	assert.Equal(t, "00000100", c)

	c = getCountIndicator(40, bytes, 4)
	assert.Equal(t, "00101000", c)

	c = getCountIndicator(2370, numeric, 38)
	assert.Equal(t, "00100101000010", c)

	c = getCountIndicator(300, alpha, 15)
	assert.Equal(t, "00100101100", c)
}

func TestEncodingNumeric(t *testing.T) {
	//t.Skip()
	en := encodingNumeric("123")
	assert.Equal(t, "1111011", en)

	en = encodingNumeric("1232")
	assert.Equal(t, "111101110", en)

	en = encodingNumeric("12310")
	assert.Equal(t, "11110111010", en)

	en = encodingNumeric("001")
	assert.Equal(t, "1", en)

	en = encodingNumeric("015")
	assert.Equal(t, "1111", en)

	en = encodingNumeric("123010")
	assert.Equal(t, "11110111010", en)

	en = encodingNumeric("12301")
	assert.Equal(t, "11110111", en)

	en = encodingNumeric("8675309")
	assert.Equal(t, "110110001110000100101001", en)
}

func TestEncodingAlpha(t *testing.T) {
	//t.Skip()
	en := encodingAlpha("AA")
	assert.Equal(t, "00111001100", en)

	en = encodingAlpha("AAA")
	assert.Equal(t, "00111001100001010", en)

	en = encodingAlpha("ABC")
	assert.Equal(t, "00111001101001100", en)

	en = encodingAlpha("HE")
	assert.Equal(t, "01100001011", en)

	en = encodingAlpha("HELLO WORLD")
	assert.Equal(t, "01100001011011110001101000101"+
		"11001011011100010011010100001101", en)
}

func TestEncodingBytes(t *testing.T) {
	//t.Skip()
	en := encodingBytes("H")
	assert.Equal(t, "01001000", en)

	en = encodingBytes("e")
	assert.Equal(t, "01100101", en)

	en = encodingBytes("Hello, world!")
	assert.Equal(t, "010010000110010101101100011011000"+
		"1101111001011000010000001110111011011110111001001"+
		"1011000110010000100001", en)
}

func TestAddTerminator(t *testing.T) {
	//t.Skip()
	ter := addTerminator("111", 1)
	assert.Equal(t, "11100000", ter)

	ter = addTerminator("111111111", 1)
	assert.Equal(t, "1111111110000000", ter)

	ter = addTerminator("111", 40)
	assert.Equal(t, "11100000", ter)
}

func TestAddTerminatorPads(t *testing.T) {
	//t.Skip()
	ter := addTerminatorPads("11100000", 1)
	assert.Equal(t, "11100000111011000001000111101100000100011110"+
		"1100000100011110110000010001", ter)

	ter = addTerminatorPads("00100000010110110000101101111000110100"+
		"010111001011011100010011010100001101000000", 2)
	assert.Equal(t, "0010000001011011000010110111100011010001011100"+
		"101101110001001101010000110100000011101100000100011110110000010"+
		"0011110110000010001", ter)
}

func seq(end int) []byte {
	seq := make([]byte, end)
	for i := 0; i < end; i++ {
		seq[i] = byte(i + 1)
	}
	return seq
}

func TestInterleaveData(t *testing.T) {
	//t.Skip()
	data := seq(9)
	inter := interleaveData(data, 1, 0, 9, 0)
	assert.Equal(t, data, inter)

	data = seq(16)
	inter = interleaveData(data, 1, 0, 16, 0)
	assert.Equal(t, data, inter)

	data = seq(26)
	result := []byte{
		1, 14, 2, 15, 3, 16, 4, 17, 5, 18, 6, 19, 7, 20, 8, 21, 9, 22, 10,
		23, 11, 24, 12, 25, 13, 26}
	inter = interleaveData(data, 2, 0, 13, 0)
	assert.Equal(t, result, inter)

	data = seq(36)
	result = []byte{
		1, 10, 19, 28, 2, 11, 20, 29, 3, 12, 21, 30, 4, 13, 22, 31, 5, 14, 23,
		32, 6, 15, 24, 33, 7, 16, 25, 34, 8, 17, 26, 35, 9, 18, 27, 36}
	inter = interleaveData(data, 4, 0, 9, 0)
	assert.Equal(t, result, inter)

	data = seq(46)
	result = []byte{
		1, 12, 23, 35, 2, 13, 24, 36, 3, 14, 25, 37, 4, 15, 26, 38, 5, 16,
		27, 39, 6, 17, 28, 40, 7, 18, 29, 41, 8, 19, 30, 42, 9, 20, 31, 43,
		10, 21, 32, 44, 11, 22, 33, 45, 34, 46}
	inter = interleaveData(data, 2, 2, 11, 12)
	assert.Equal(t, result, inter)
}

func TestErrorCorrection(t *testing.T) {
	//t.Skip()
	data := seq(9)
	result := []byte{
		48, 91, 177, 195, 164, 76, 88, 163, 118, 109, 85, 153, 22, 5, 66, 122,
		247}
	correction := errorCorrection(data, 17, 1, 0, 9, 0)
	assert.Equal(t, result, correction)

	data = seq(16)
	result = []byte{
		164, 102, 147, 155, 85, 236, 194, 153, 26, 186, 202, 157, 234, 245,
		221, 19, 232, 248, 229, 173, 171, 47, 250, 135, 174, 17, 203, 203}
	correction = errorCorrection(data, 28, 1, 0, 16, 0)
	assert.Equal(t, result, correction)

	data = seq(26)
	result = []byte{
		108, 198, 169, 133, 107, 122, 23, 61, 170, 4, 240, 137, 36, 218, 167,
		73, 75, 243, 140, 7, 102, 46, 61, 246, 34, 63, 18, 3, 37, 139, 140,
		40, 21, 168, 44, 195, 134, 38, 181, 167, 97, 104, 218, 228}
	correction = errorCorrection(data, 22, 2, 0, 13, 0)
	assert.Equal(t, result, correction)

	data = seq(36)
	result = []byte{
		177, 172, 114, 149, 225, 181, 117, 6, 10, 141, 19, 137, 16, 48, 72,
		169, 208, 89, 208, 15, 151, 75, 173, 17, 157, 177, 179, 25, 230, 49,
		129, 244, 144, 68, 189, 136, 109, 242, 34, 116, 204, 12, 60, 76, 6,
		33, 141, 33, 199, 89, 153, 158, 184, 183, 163, 92, 79, 147, 93, 154,
		126, 131, 188, 23}
	correction = errorCorrection(data, 16, 4, 0, 9, 0)
	assert.Equal(t, result, correction)
}

func TestInterleaveError(t *testing.T) {
	//t.Skip()
	err := seq(17)
	inter := interleaveError(err, 17, 1, 0)
	assert.Equal(t, err, inter)

	err = seq(28)
	inter = interleaveError(err, 28, 1, 0)
	assert.Equal(t, err, inter)

	err = seq(44)
	result := []byte{
		1, 23, 2, 24, 3, 25, 4, 26, 5, 27, 6, 28, 7, 29, 8, 30, 9, 31, 10, 32,
		11, 33, 12, 34, 13, 35, 14, 36, 15, 37, 16, 38, 17, 39, 18, 40, 19,
		41, 20, 42, 21, 43, 22, 44}
	inter = interleaveError(err, 22, 2, 0)
	assert.Equal(t, result, inter)
}

func dataEncoding() string {
	data := "Lausanne"
	length := len(data)
	mode := getMode(data)
	version := getVersion(data, mode)

	enc := ""
	if mode == numeric {
		enc = indicatorNumeric + getCountIndicator(length, mode, version) + encodingNumeric(data)
	} else if mode == alpha {
		enc = indicatorAlpha + getCountIndicator(length, mode, version) + encodingAlpha(data)
	} else {
		enc = indicatorBytes + getCountIndicator(length, mode, version) + encodingBytes(data)
	}
	enc = addTerminatorPads(addTerminator(enc, version), version)

	numErr := blockInfo[version][1]
	numB1 := blockInfo[version][2]
	numW1 := blockInfo[version][3]
	numB2 := blockInfo[version][4]
	numW2 := blockInfo[version][5]

	dataBytes := encodingToByteArray(enc)
	errorBytes := errorCorrection(dataBytes, numErr, numB1, numB2, numW1, numW2)
	interData := interleaveData(dataBytes, numB1, numB2, numW1, numW2)
	interError := interleaveError(errorBytes, numErr, numB1, numB2)

	enc = byteArrayToEncoding(interData) + byteArrayToEncoding(interError)
	enc = padRight(enc, len(enc)+blockInfo[version][6])

	return enc
}

//func TestDrawFinderPattern(t *testing.T) {
//enc := dataEncoding()
//canvas := make([][]*Cell, 25)
//for i, _ := range canvas {
//canvas[i] = make([]*Cell, 25)
//}

//for i := 0; i < len(canvas); i++ {
//for j := 0; j < len(canvas); j++ {
//canvas[i][j] = &Cell{0, true}
//}
//}

//version := 2

//placeFinderPatterns(canvas, version)
//placeSeperators(canvas)
//placeAlignmentPatterns(canvas, version)
//drawTimingPattern(canvas, version)
//drawDarkModule(canvas, version)
//reserveFormatInformationArea(canvas, version)

//if version >= 7 {
//reserveVersionInformationData(canvas, version)
//}

//drawDataBits(canvas, enc, version)

//winner, mask := dataMasking(canvas)
//println(mask)
//drawFormatInformationString(winner, mask)

//if version >= 7 {
//drawVersionInformationString(winner, version)
//}
//outputTerminal(winner)
//}
