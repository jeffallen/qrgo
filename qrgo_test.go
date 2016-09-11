package qrgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	assert.Equal(t, 0, max(0, 0))
	assert.Equal(t, 1, max(0, 1))
	assert.Equal(t, 1, max(1, 0))
	assert.Equal(t, 2, max(1, 2))
	assert.Equal(t, 2, max(2, 1))
}

func TestPadLeft(t *testing.T) {
	assert.Equal(t, "000", padLeft("000", 3))
	assert.Equal(t, "000", padLeft("", 3))
	assert.Equal(t, "000", padLeft("0", 3))
	assert.Equal(t, "000", padLeft("00", 3))
	assert.Equal(t, "001", padLeft("1", 3))
	assert.Equal(t, "011", padLeft("11", 3))
}

func TestPadRight(t *testing.T) {
	assert.Equal(t, "000", padRight("000", 3))
	assert.Equal(t, "000", padRight("", 3))
	assert.Equal(t, "000", padRight("0", 3))
	assert.Equal(t, "000", padRight("00", 3))
	assert.Equal(t, "100", padRight("1", 3))
	assert.Equal(t, "110", padRight("11", 3))
}

func TestStringToBinary(t *testing.T) {
	assert.Equal(t, "", stringToBinary(""))
	assert.Equal(t, "0", stringToBinary("0"))
	assert.Equal(t, "1", stringToBinary("1"))
	assert.Equal(t, "10", stringToBinary("2"))
	assert.Equal(t, "11", stringToBinary("3"))
	assert.Equal(t, "100", stringToBinary("4"))
}

func TestEncodingToByteArray(t *testing.T) {
	assert.Equal(t, []byte{0}, encodingToByteArray("00000000"))
	assert.Equal(t, []byte{1}, encodingToByteArray("00000001"))
	assert.Equal(t, []byte{2}, encodingToByteArray("00000010"))
	assert.Equal(t, []byte{0, 1}, encodingToByteArray("0000000000000001"))
	assert.Equal(t, []byte{1, 1}, encodingToByteArray("0000000100000001"))
	assert.Equal(t, []byte{1, 2}, encodingToByteArray("0000000100000010"))
}

func TestByteArrayToEncoding(t *testing.T) {
	assert.Equal(t, "00000000", byteArrayToEncoding([]byte{0}))
	assert.Equal(t, "00000001", byteArrayToEncoding([]byte{1}))
	assert.Equal(t, "00000010", byteArrayToEncoding([]byte{2}))
	assert.Equal(t, "0000000000000001", byteArrayToEncoding([]byte{0, 1}))
	assert.Equal(t, "0000000100000001", byteArrayToEncoding([]byte{1, 1}))
	assert.Equal(t, "0000000100000010", byteArrayToEncoding([]byte{1, 2}))
}

func TestBinarySearch(t *testing.T) {
	assert.Equal(t, 0, binarySearch(maxCharsNumeric, 20))
	assert.Equal(t, 0, binarySearch(maxCharsNumeric, 41))
	assert.Equal(t, 6, binarySearch(maxCharsNumeric, 360))
	assert.Equal(t, 39, binarySearch(maxCharsNumeric, 7089))
	assert.Equal(t, 0, binarySearch(maxCharsAlpha, 20))
	assert.Equal(t, 0, binarySearch(maxCharsAlpha, 25))
	assert.Equal(t, 6, binarySearch(maxCharsAlpha, 200))
	assert.Equal(t, 39, binarySearch(maxCharsAlpha, 4296))
	assert.Equal(t, 0, binarySearch(maxCharsBytes, 15))
	assert.Equal(t, 0, binarySearch(maxCharsBytes, 17))
	assert.Equal(t, 6, binarySearch(maxCharsBytes, 150))
	assert.Equal(t, 39, binarySearch(maxCharsBytes, 2953))
}

func TestCountIndicator(t *testing.T) {
	assert.Equal(t, "0000000001", indCount(1, numeric, 1))
	assert.Equal(t, "000000000010", indCount(2, numeric, 15))
	assert.Equal(t, "00000000000011", indCount(3, numeric, 40))
	assert.Equal(t, "000000001", indCount(1, alpha, 1))
	assert.Equal(t, "00000000010", indCount(2, alpha, 15))
	assert.Equal(t, "0000000000011", indCount(3, alpha, 40))
	assert.Equal(t, "00000001", indCount(1, byteMode, 1))
	assert.Equal(t, "0000000000000010", indCount(2, byteMode, 15))
	assert.Equal(t, "0000000000000011", indCount(3, byteMode, 40))
}

func TestEncodingNumeric(t *testing.T) {
	assert.Equal(t, "0", encNumeric("0"))
	assert.Equal(t, "1", encNumeric("1"))
	assert.Equal(t, "10", encNumeric("2"))
	assert.Equal(t, "10", encNumeric("02"))
	assert.Equal(t, "1010", encNumeric("10"))
	assert.Equal(t, "1100100", encNumeric("100"))
	// Doc example
	assert.Equal(t, "110110001110000100101001", encNumeric("8675309"))
}

func TestEncodingAlpha(t *testing.T) {
	assert.Equal(t, "000000", encAlpha("0"))
	assert.Equal(t, "001010", encAlpha("A"))
	assert.Equal(t, "001011", encAlpha("B"))
	assert.Equal(t, "00111001101", encAlpha("AB"))
	// Doc example
	assert.Equal(t, "01100001011", encAlpha("HE"))
}

func TestEncodingBytes(t *testing.T) {
	assert.Equal(t, "01100001", encBytes("a"))
	assert.Equal(t, "01100010", encBytes("b"))
	assert.Equal(t, "01100011", encBytes("c"))
	assert.Equal(t, "0110000101100010", encBytes("ab"))
	//Doc example
	assert.Equal(t, "0100100001100101", encBytes("He"))
}

func TestMain(t *testing.T) {
	qr, _ := NewQR("EPFLLAUSANNE2016SWITZERLAND")
	qr.OutputTerminal()
}
