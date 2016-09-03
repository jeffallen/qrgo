package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// The finalized QR-encoding of an input string
// containing all preliminary steps the lead to the
// result.
type QR struct {
	data     string
	length   int
	mode     int
	version  int
	encoding string
	canvas   [][]*Cell
}

type Cell struct {
	color int
	data  bool
}

type mask func(row, col int) bool

const (
	regexNumeric = "^.[0-9]*$"                   // [0,9]
	regexAlpha   = "^.[0-9A-Z /.:%+\\$\\-\\*]*$" // [0,9] | [A,Z] | {/.:%+$-*}

	numeric = 1
	alpha   = 2
	bytes   = 4

	indicatorNumeric = "0001"
	indicatorAlpha   = "0010"
	indicatorBytes   = "0100"

	versions = 40

	white = "\033[47m  \033[0m"
	black = "\033[40m  \033[0m"
)

var (
	maxCharsNumeric = []int{
		41, 77, 127, 187, 255, 322, 370, 461, 552, 652, 772, 883, 1022,
		1101, 1250, 1408, 1548, 1725, 1903, 2061, 2232, 2409, 2620, 2812,
		3057, 3283, 3517, 3669, 3909, 4158, 4417, 4686, 4965, 5253, 5529,
		5836, 6153, 6479, 6743, 7089}

	maxCharsAlpha = []int{
		25, 47, 77, 114, 154, 195, 224, 279, 335, 395, 468, 535, 619,
		667, 758, 854, 938, 1046, 1153, 1249, 1352, 1460, 1588, 1704,
		1853, 1990, 2132, 2223, 2369, 2520, 2677, 2840, 3009, 3183,
		3351, 3537, 3729, 3927, 4087, 4296}

	maxCharsBytes = []int{
		17, 32, 53, 78, 106, 134, 154, 192, 230, 271, 321, 367, 425, 458,
		520, 586, 644, 718, 792, 858, 929, 779, 1091, 1171, 1273, 1367,
		1465, 1528, 1628, 1732, 1840, 1952, 2068, 2188, 2303, 2431, 2563,
		2699, 2809, 2953}

	alphaTable = map[rune]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8,
		'9': 9, 'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16,
		'H': 17, 'I': 18, 'J': 19, 'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24,
		'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30, 'V': 31, 'W': 32,
		'X': 33, 'Y': 34, 'Z': 35, ' ': 36, '$': 37, '%': 38, '*': 39, '+': 40,
		'-': 41, '.': 42, '/': 43, ':': 44}

	// 1: # Total Codewords
	// 2: # Error Correction Words
	// 3: # Blocks in Group1
	// 4: # Codewords in Blocks of Group1
	// 5: # Blocks in Group2
	// 6: # Codewords in Blocks of Group2
	// 7: # Remainder Bits
	blockInfo = map[int][7]int{
		1: {19, 7, 1, 19, 0, 0, 0}, 2: {34, 10, 1, 34, 0, 0, 7},
		3: {55, 15, 1, 55, 0, 0, 7}, 4: {80, 20, 1, 80, 0, 0, 7},
		5: {108, 26, 1, 108, 0, 0, 7}, 6: {136, 18, 2, 68, 0, 0, 7},
		7: {156, 20, 2, 78, 0, 0, 0}, 8: {194, 24, 2, 97, 0, 0, 0},
		9: {232, 30, 2, 116, 0, 0, 0}, 10: {274, 18, 2, 68, 2, 69, 0},
		11: {324, 20, 4, 81, 0, 0, 0}}

	alignmentPatterns = map[int][]int{
		2: {18, 18}, 3: {22, 22}, 4: {26, 26}, 5: {30, 30}, 6: {34, 34},
		7:  {6, 22, 22, 6, 22, 22, 22, 38, 38, 22, 38, 38},
		8:  {6, 24, 24, 6, 24, 24, 24, 42, 42, 24, 42, 42},
		9:  {6, 26, 26, 6, 26, 26, 26, 46, 46, 26, 46, 46},
		10: {6, 28, 28, 6, 28, 28, 28, 50, 50, 28, 50, 50},
		11: {6, 30, 30, 6, 30, 30, 30, 54, 54, 30, 54, 54},
		12: {6, 32, 32, 6, 32, 32, 32, 58, 58, 32, 58, 58}}

	terminatorPads   = []string{"11101100", "00010001"}
	penaltySequences = []string{"10111010000", "00001011101"}

	masks = []mask{mask0, mask1, mask2, mask3, mask4, mask5, mask6, mask7}

	formatInformationStrings = []string{
		"111011111000100", "111001011110011", "111110110101010",
		"111100010011101", "110011000101111", "110001100011000",
		"110110001000001", "110100101110110"}

	versionInformationStrings = []string{
		"000111110010010100", "001000010110111100", "001001101010011001",
		"001010010011010011", "001011101111110110", "001100011101100010"}
)

// Return number of modules depending on version.
//func size(version int) int {
//return ((version - 1) * 4) + 21
//}

var size int

// Canonical integer max function.
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// Left pad a binary string up until to final length.
// Also applies to empty string but leaves string unchanged
// if its length is already at least the final length.
func padLeft(binary string, final int) string {
	length := len(binary)
	if length >= final {
		return binary
	}
	padded := binary
	for i := 0; i < final-length; i++ {
		padded = "0" + padded
	}
	return padded
}

// Right pad a binary string up until to final length.
// Also applies to empty string but leaves string unchanged
// if its length is already at least the final length.
func padRight(binary string, final int) string {
	length := len(binary)
	if length >= final {
		return binary
	}
	padded := binary
	for i := 0; i < final-length; i++ {
		padded = padded + "0"
	}
	return padded
}

// Converts a string number to its binary representation.
// Returns an empty string if input number is empty as well.
func stringToBinary(number string) string {
	if len(number) == 0 {
		return ""
	}
	bin, _ := strconv.ParseInt(number, 10, 64)
	return strconv.FormatInt(bin, 2)
}

// Convert a binary string to the corresponing byte array
// representation.
//
//		"0000001011111111" -> [2, 255]
//
func encodingToByteArray(encoding string) []byte {
	length := len(encoding)
	array := make([]byte, length/8)
	for i, j := 0, 0; i < length; i, j = i+8, j+1 {
		b, _ := strconv.ParseInt(encoding[i:i+8], 2, 64)
		array[j] = byte(b)
	}
	return array
}

// Inverse of encodingToByteArray. Create a binary string
// out of a byte array.
//
//		[2, 255] -> "0000001011111111"
//
func byteArrayToEncoding(array []byte) string {
	encoding := ""
	for _, b := range array {
		encoding += padLeft(strconv.FormatInt(int64(b), 2), 8)
	}
	return encoding
}

// Returns the appropriate mode for a given
// data string.
//
//		Numeric (1): All string only consisting of numbers.
//      Alpha   (2): All strings consisting of numbers, upper-
//					 case letters and {/.:%+$-*}.
//		Bytes	(4): All strings containing ASCII-characters.
//
func getMode(data string) int {
	rNumeric, _ := regexp.Compile(regexNumeric)
	rAlpha, _ := regexp.Compile(regexAlpha)

	if rNumeric.MatchString(data) {
		return numeric
	} else if rAlpha.MatchString(data) {
		return alpha
	}
	return bytes
}

// Returns the qr-version corresponding to the length and
// mode of the data string. The version lies in [1, 40]
// the length doesn't fit a version zero is returned.
func getVersion(data string, mode int) int {
	length := len(data)
	var array []int
	if mode == numeric {
		array = maxCharsNumeric
	} else if mode == alpha {
		array = maxCharsAlpha
	} else {
		array = maxCharsBytes
	}

	// Binary search through version table.
	low, high := 0, versions
	for low != high {
		mid := (low + high) / 2
		if array[mid] < length {
			low = mid + 1
		} else {
			high = mid
		}
	}

	if low == versions {
		return 0
	} else {
		return low + 1
	}
}

// The count indicator follows the mode indicator in the
// encoding. The indicator's length differs in connection
// to the given length, mode and version of the data string.
//
//		Version [1, 9]:
//			Numeric:	10 bytes
//			Alpha: 		9 bytes
//			Bytes:		8 bytes
//
//		Version [10, 26]:
//			Numeric:	12 bytes
//			Alpha:		11 bytes
//			Bytes:		16 bytes
//
//		Version [26, 40]:
//			Numeric:	14 bytes
//			Alpha:		13 bytes
//			Bytes:		16 bytes
//
func getCountIndicator(length, mode, version int) string {
	count := strconv.FormatInt(int64(length), 2)
	if version >= 1 && version <= 9 {
		if mode == numeric {
			return padLeft(count, 10)
		} else if mode == alpha {
			return padLeft(count, 9)
		} else {
			return padLeft(count, 8)
		}
	} else if version >= 10 && version <= 26 {
		if mode == numeric {
			return padLeft(count, 12)
		} else if mode == alpha {
			return padLeft(count, 11)
		} else {
			return padLeft(count, 16)
		}
	} else {
		if mode == numeric {
			return padLeft(count, 14)
		} else if mode == alpha {
			return padLeft(count, 13)
		} else {
			return padLeft(count, 16)
		}
	}
}

// The numeric encoding converts every three-digit number
// in the data string into its (unsigned) binary representation.
// Hanging numbers at the end are equally turned into binary,
// whereas zeros in front of a number are ignored.
//
//		8675309:
//			867 -> 1101100011
//			530 -> 1000010010
//			9 	-> 1001
//
func encodingNumeric(data string) string {
	i, encoding := 0, ""
	for ; i <= len(data)-3; i += 3 {
		encoding += stringToBinary(data[i : i+3])
	}
	return encoding + stringToBinary(data[i:])
}

// The alphanumeric encoding takes groups of two chars,
// respectively their decimal value from the alphaTable and
// multiplies the first number with 45, then adds the second
// value. The result is converted to 11-bit binary representation.
// If there is a char left at the end, its value is turned into a
// 6-bit binary.
//
//		HELLO WORLD:
//			H -> 17
//			E -> 14
//			(45 * 17) + 14 = 779 = 01100001011
//
func encodingAlpha(data string) string {
	i, encoding := 0, ""
	for ; i <= len(data)-2; i += 2 {
		pair := data[i : i+2]
		num := alphaTable[rune(pair[0])]*45 + alphaTable[rune(pair[1])]
		encoding += padLeft(strconv.FormatInt(int64(num), 2), 11)
	}
	tail := data[i:] // Possible hanging char.
	if tail == "" {
		return encoding
	} else {
		end := int64(alphaTable[rune(tail[0])])
		return encoding + padLeft(strconv.FormatInt(end, 2), 6)
	}
}

// The byte encoding simply turns every char in the data
// string into its ASCII 8-bit binary representation.
//
//		H -> 0x48 -> 01001000
//		e -> 0x65 -> 01100101
//
func encodingBytes(data string) string {
	encoding := ""
	for _, c := range data {
		ascii := int64(c)
		encoding += padLeft(strconv.FormatInt(ascii, 2), 8)
	}
	return encoding
}

// After encoding the data string it needs to be right-padded
// until the length is a multiple of eight.
func addTerminator(encoding string, version int) string {
	length := len(encoding)
	bits := blockInfo[version][0] * 8

	if bits-length == 0 {
		return encoding
	}

	rest := 8 - (length % 8)
	return padRight(encoding, length+rest)
}

// If after appending the terminator bits the string is still not
// at the full capacity the terminator pads {11101100, 00010001}
// are added alternately until the maximum length is reached.
func addTerminatorPads(encoding string, version int) string {
	numPads := blockInfo[version][0] - len(encoding)/8
	padding := encoding
	for i := 0; i < numPads; i++ {
		padding = padding + terminatorPads[i%2]
	}
	return padding
}

// The data byte-array has to be interleaved according to the QR-Code
// specification. In simplified terms the first bytes of every blocks
// are aligned before moving to the second bytes. Data strings only
// consisting of a single block are not interleaved.
//
//		Group1, Block1: [1, 2, 3, 4]
//		Group1, Block2: [5, 6, 7, 8]
//		-> [1, 5, 2, 6, 3, 7, 4, 8]
//
//		Group1, Block1: [1, 2, 3, 4]
//		Group1, Block2: [5, 6, 7, 8]
//		Group2, Block1: [9, 10, 11]
//		Group2, Block2: [12, 13, 14]
//		-> [1, 5, 9, 12, 2, 6, 10, 13, 4, 8, 11, 14]
//
func interleaveData(array []byte, blocks1, blocks2, words1, words2 int) []byte {
	inter := make([]byte, len(array))

	// Iteration through all the words in the same block column
	// before advancing to the next column.
	x, b1, b2 := 0, 0, 0
	for i := 0; i < max(words1, words2); i++ {
		y := i
		if b1 < words1 {
			for j := 0; j < blocks1; j++ {
				inter[x] = array[y]
				y += words1
				x++
			}
			b1++
		} else {
			// Blocks in Group1 are smaller than Blocks in Group2,
			// thus skipping empty cells to advance to Group2, when
			// all columns of Group1 are traversed.
			y += blocks1 * words1
		}

		if b2 < words2 {
			for j := 0; j < blocks2; j++ {
				inter[x] = array[y]
				y += words2
				x++
			}
			b2++
		}
	}
	return inter
}

// For each block in a data string a certain number error correction
// words are generated applying the Reed-Solomon standard.
func errorCorrection(array []byte, numErr, numB1, numB2, numW1, numW2 int) []byte {
	correction := make([]byte, numErr*(numB1+numB2))
	field := NewField(0x11d, 2) // Reed-Solomon parameters for QR-Codes.
	enc := NewRSEncoder(field, numErr)

	x, y, b1, b2 := 0, 0, 0, 0
	for i := 0; i < numB1+numB2; i++ {
		check := make([]byte, numErr)
		if b1 < numB1 {
			enc.ECC(array[x:x+numW1], check)
			copy(correction[y:y+numErr], check)
			x += numW1
			y += numErr
			b1++
		} else if b2 < numB2 {
			enc.ECC(array[x:x+numW2], check)
			copy(correction[y:y+numErr], check)
			x += numW2
			y += numErr
			b2++
		}
	}
	return correction
}

// The error correction words are interleaved in the same way as the
// data bytes. Since all blocks have the same number of error words, it
// boils down to a transposition.
//
//		[1, 2, 3, 4][5, 6, 7, 8]
//		-> [1, 5, 2, 6 ,3, 7, 4, 8]
//
func interleaveError(array []byte, numErr, numB1, numB2 int) []byte {
	inter := make([]byte, len(array))

	y := 0
	for i := 0; i < numErr; i++ {
		x := i
		for j := 0; j < numB1+numB2; j++ {
			inter[y] = array[x]
			x += numErr
			y++
		}
	}
	return inter
}

// Every QR-Code contains three finder patterns located in the
// two upper and bottom-left corners.
func placeFinderPatterns(canvas [][]*Cell, version int) {
	drawPattern(canvas, 0, 0, 7)
	drawPattern(canvas, size-7, 0, 7)
	drawPattern(canvas, 0, size-7, 7)
}

// A finder pattern is 7x7 square with a black border followed by
// a white ring embeding a 3x3 black square in the center.
// An alignment pattern is a 5x5 square with a black border followed
// by white surrounding a black module in the middle.
//
//		Finder pattern:
//		1111111
//		1000001
//		1011101
//		1011101
//		1011101
//		1000001
//		1111111
//
//		Alignment pattern:
//		11111
//		10001
//		10101
//		10001
//		11111
//
func drawPattern(canvas [][]*Cell, row, col, size int) {
	inner := size - 2
	for r := row; r < row+size; r++ {
		for c := col; c < col+size; c++ {
			if c >= col+1 && c <= col+inner && r == row+1 ||
				c >= col+1 && c <= col+inner && r == row+inner ||
				r >= row+1 && r <= row+inner && c == col+1 ||
				r >= row+1 && r <= row+inner && c == col+inner {
				canvas[r][c].color = 0
			} else {
				canvas[r][c].color = 1
			}
			canvas[r][c].data = false
		}
	}
}

// The seperators are white lines surrounding the three finder
// patterns.
func placeSeperators(canvas [][]*Cell) {
	drawSeperator(canvas, 7, 7, -1, -1)
	drawSeperator(canvas, 7, len(canvas)-8, -1, 1)
	drawSeperator(canvas, len(canvas)-8, 7, 1, -1)
}

// Example of the seperator for the top-left finder pattern.
//
//		11111110
//		10000010
//		10111010
//		10111010
//		10111010
//		10000010
//		11111110
//		00000000
//
func drawSeperator(canvas [][]*Cell, row, col, dr, dc int) {
	c := col
	for i := 0; i <= 7; i++ {
		canvas[row][c].color = 0
		canvas[row][c].data = false
		c += dc
	}

	r := row
	for i := 0; i <= 7; i++ {
		canvas[r][col].color = 0
		canvas[r][col].data = false
		r += dr
	}
}

// Every QR-Code that is not of version 1, has one or more
// alignmentPatterns. The positions of those patterns can
// be found in the alignmentPattern table.
func placeAlignmentPatterns(canvas [][]*Cell, version int) {
	patterns := alignmentPatterns[version]

	for i := 0; i < len(patterns)-1; i += 2 {
		drawPattern(canvas, patterns[i]-2, patterns[i+1]-2, 5)
	}
}

// The timing pattern is an alternating line on the seventh
// row and seventh column between the connecting the finder
// patterns.
//
//		Left timing pattern:
//		1111111
//		1000001
//		1011101
//		1011101
//		1011101
//		1000001
//		1111111
//			  0
//			  1
//			  0
//		1111111
//		1000001
//		1011101
//		1011101
//		1011101
//		1000001
//		1111111
//
func drawTimingPattern(canvas [][]*Cell, version int) {
	length := size - 14

	for i := 6; i < 6+length; i++ {
		if i%2 == 0 {
			canvas[6][i].color = 1
			canvas[i][6].color = 1
		} else {
			canvas[6][i].color = 0
			canvas[i][6].color = 0
		}

		canvas[6][i].data = false
		canvas[i][6].data = false
	}
}

// The dark module is always placed at coordinates ((4 * V) +9, 8)
func drawDarkModule(canvas [][]*Cell, version int) {
	r := (4 * version) + 9
	canvas[r][8].color = 1
	canvas[r][8].data = false
}

func reserveFormatInformationArea(canvas [][]*Cell, version int) {
	canvas[8][0].data = false
	canvas[8][1].data = false
	canvas[8][2].data = false
	canvas[8][3].data = false
	canvas[8][4].data = false
	canvas[8][5].data = false
	canvas[8][7].data = false
	canvas[8][8].data = false

	canvas[0][8].data = false
	canvas[1][8].data = false
	canvas[2][8].data = false
	canvas[3][8].data = false
	canvas[4][8].data = false
	canvas[5][8].data = false
	canvas[7][8].data = false

	canvas[8][size-8].data = false
	canvas[8][size-7].data = false
	canvas[8][size-6].data = false
	canvas[8][size-5].data = false
	canvas[8][size-4].data = false
	canvas[8][size-3].data = false
	canvas[8][size-2].data = false
	canvas[8][size-1].data = false

	canvas[size-7][8].data = false
	canvas[size-6][8].data = false
	canvas[size-5][8].data = false
	canvas[size-4][8].data = false
	canvas[size-3][8].data = false
	canvas[size-2][8].data = false
	canvas[size-1][8].data = false
}

func reserveVersionInformationData(canvas [][]*Cell, version int) {
	for i := 0; i < 6; i++ {
		canvas[size-11][i].data = false
		canvas[size-10][i].data = false
		canvas[size-9][i].data = false

		canvas[i][size-9].data = false
		canvas[i][size-10].data = false
		canvas[i][size-11].data = false
	}
}

func drawDataBits(canvas [][]*Cell, data string, version int) {
	i, up := 0, true
	for c := size - 1; c > 0; c -= 2 {
		if c == 6 {
			c++
		} else {
			if up {
				for r := size - 1; r >= 0; r-- {
					if canvas[r][c].data {
						wb, _ := strconv.Atoi(string(data[i]))
						canvas[r][c].color = wb
						i++
					}
					if canvas[r][c-1].data {
						wb, _ := strconv.Atoi(string(data[i]))
						canvas[r][c-1].color = wb
						i++
					}
				}
			} else {
				for r := 0; r < size; r++ {
					if canvas[r][c].data {
						wb, _ := strconv.Atoi(string(data[i]))
						canvas[r][c].color = wb
						i++
					}
					if canvas[r][c-1].data {
						wb, _ := strconv.Atoi(string(data[i]))
						canvas[r][c-1].color = wb
						i++
					}
				}
			}
			up = !up
		}
	}
}

func initCanvas(length int) [][]*Cell {
	canvas := make([][]*Cell, length)
	for i, _ := range canvas {
		canvas[i] = make([]*Cell, length)
	}

	for i := 0; i < len(canvas); i++ {
		for j := 0; j < len(canvas); j++ {
			canvas[i][j] = &Cell{0, true}
		}
	}
	return canvas
}

func deepCopyCanvas(canvas, copy [][]*Cell) {
	length := len(canvas)
	for r := 0; r < length; r++ {
		for c := 0; c < length; c++ {
			color := canvas[r][c].color
			data := canvas[r][c].data
			copy[r][c] = &Cell{color: color, data: data}
		}
	}
}

func mask0(row, col int) bool {
	return (row+col)%2 == 0
}

func mask1(row, col int) bool {
	return row%2 == 0
}

func mask2(row, col int) bool {
	return col%3 == 0
}

func mask3(row, col int) bool {
	return (row+col)%3 == 0
}

func mask4(row, col int) bool {
	return ((row/2)+(col/3))%2 == 0
}

func mask5(row, col int) bool {
	return ((row*col)%2)+((row*col)%3) == 0
}

func mask6(row, col int) bool {
	return (((row*col)%2)+((row*col)%3))%2 == 0
}

func mask7(row, col int) bool {
	return (((row+col)%2)+((row*col)%3))%2 == 0
}

func maskCanvas(canvas [][]*Cell, fn mask) [][]*Cell {
	length := len(canvas)
	masked := initCanvas(length)
	deepCopyCanvas(canvas, masked)

	for r := 0; r < length; r++ {
		for c := 0; c < length; c++ {
			if masked[r][c].data && fn(r, c) {
				masked[r][c].color ^= 1
			}
		}
	}
	return masked
}

func penalty1(masked [][]*Cell) int {
	length := len(masked)

	total := 0
	for r := 0; r < length; r++ {
		count, prev, diff := 0, -1, 0
		for c := 0; c < length; c++ {
			color := masked[r][c].color
			if color == prev {
				count++
				prev = color
			} else {
				diff = count - 5
				prev = color
				if diff >= 0 {
					total += (3 + diff)
				}
				count = 1
			}
		}
		diff = count - 5
		if diff >= 0 {
			total += (3 + diff)
		}
	}
	return total
}

func penalty2(masked [][]*Cell) int {
	length := len(masked)

	total := 0
	for r := 0; r < 1; r++ {
		for c := 0; c < length-1; c++ {
			sum := masked[r][c].color + masked[r][c+1].color +
				masked[r+1][c].color + masked[r+1][c+1].color
			if sum == 0 || sum == 4 {
				total += 3
			}
		}
	}
	return total
}

// Counts the number of overlapping substring within an
// other string.
func countSubstringOccurrences(s, sub string) int {
	sLength, subLength := len(s), len(sub)
	total := 0
	for i := 0; i < sLength-subLength+1; i++ {
		if strings.Contains(s[i:i+subLength], sub) {
			total++
		}
	}
	return total
}

func penalty3(masked [][]*Cell) int {
	length := len(masked)

	total := 0
	for i := 0; i < length; i++ {
		row, col := "", ""
		for j := 0; j < length; j++ {
			row += strconv.Itoa(masked[i][j].color)
			col += strconv.Itoa(masked[j][i].color)
		}
		total += countSubstringOccurrences(row, penaltySequences[0]) * 40
		total += countSubstringOccurrences(col, penaltySequences[1]) * 40
	}
	return total
}

func penalty4(masked [][]*Cell) int {
	length := len(masked)

	numModules := length * length
	numBlack := 0
	for r := 0; r < length; r++ {
		for c := 0; c < length; c++ {
			numBlack += masked[r][c].color
		}
	}

	ratio := float64(numBlack) / float64(numModules) * 100
	low := int(ratio/5) * 5
	up := (int(ratio/5) + 1) * 5
	down := math.Abs(float64(low)-50) / 5.0
	top := math.Abs(float64(up)-50) / 5.0

	if down <= top {
		return int(down * 10)
	} else {
		return int(top * 10)
	}
}

func dataMasking(canvas [][]*Cell) ([][]*Cell, int) {
	winner, mask, min := [][]*Cell{}, -1, math.MaxInt64
	for i := 0; i < 8; i++ {
		masked := maskCanvas(canvas, masks[i])
		penalty := penalty1(masked)
		penalty += penalty2(masked)
		penalty += penalty3(masked)
		penalty += penalty4(masked)
		if penalty < min {
			winner = masked
			mask = i
			min = penalty
		}
	}
	return winner, mask
}

func drawFormatInformationString(masked [][]*Cell, mask int) {
	length := len(masked)
	fis := formatInformationStrings[mask]

	wb, _ := strconv.Atoi(string(fis[0]))
	masked[8][0].color = wb
	masked[length-1][8].color = wb

	wb, _ = strconv.Atoi(string(fis[1]))
	masked[8][1].color = wb
	masked[length-2][8].color = wb

	wb, _ = strconv.Atoi(string(fis[2]))
	masked[8][2].color = wb
	masked[length-3][8].color = wb

	wb, _ = strconv.Atoi(string(fis[3]))
	masked[8][3].color = wb
	masked[length-4][8].color = wb

	wb, _ = strconv.Atoi(string(fis[4]))
	masked[8][4].color = wb
	masked[length-5][8].color = wb

	wb, _ = strconv.Atoi(string(fis[5]))
	masked[8][5].color = wb
	masked[length-6][8].color = wb

	wb, _ = strconv.Atoi(string(fis[6]))
	masked[8][7].color = wb
	masked[length-7][8].color = wb

	wb, _ = strconv.Atoi(string(fis[7]))
	masked[8][8].color = wb
	masked[8][length-8].color = wb

	wb, _ = strconv.Atoi(string(fis[8]))
	masked[7][8].color = wb
	masked[8][length-7].color = wb

	wb, _ = strconv.Atoi(string(fis[9]))
	masked[5][8].color = wb
	masked[8][length-6].color = wb

	wb, _ = strconv.Atoi(string(fis[10]))
	masked[4][8].color = wb
	masked[8][length-5].color = wb

	wb, _ = strconv.Atoi(string(fis[11]))
	masked[3][8].color = wb
	masked[8][length-4].color = wb

	wb, _ = strconv.Atoi(string(fis[12]))
	masked[2][8].color = wb
	masked[8][length-3].color = wb

	wb, _ = strconv.Atoi(string(fis[13]))
	masked[1][8].color = wb
	masked[8][length-2].color = wb

	wb, _ = strconv.Atoi(string(fis[14]))
	masked[0][8].color = wb
	masked[8][length-1].color = wb
}

func drawVersionInformationString(masked [][]*Cell, version int) {
	length := len(masked)
	vis := versionInformationStrings[version-7]

	x := 0
	for i := 5; i >= 0; i-- {
		for j := 0; j < 3; j++ {
			wb, _ := strconv.Atoi(string(vis[x]))
			masked[length-(9+j)][i].color = wb
			masked[i][length-(9+j)].color = wb
			x++
		}
	}
}

func assembleEncoding(data string, length, mode, version int) string {
	enc := ""
	if mode == numeric {
		enc = indicatorNumeric +
			getCountIndicator(length, mode, version) +
			encodingNumeric(data)
	} else if mode == alpha {
		enc = indicatorAlpha +
			getCountIndicator(length, mode, version) +
			encodingAlpha(data)
	} else {
		enc = indicatorBytes +
			getCountIndicator(length, mode, version) + encodingBytes(data)
	}
	return addTerminatorPads(addTerminator(enc, version), version)
}

func interleaveEncoding(encoding string, version int) string {
	numErr := blockInfo[version][1]
	numB1 := blockInfo[version][2]
	numW1 := blockInfo[version][3]
	numB2 := blockInfo[version][4]
	numW2 := blockInfo[version][5]

	dataBytes := encodingToByteArray(encoding)
	errorBytes := errorCorrection(dataBytes, numErr, numB1, numB2, numW1, numW2)
	interData := interleaveData(dataBytes, numB1, numB2, numW1, numW2)
	interError := interleaveError(errorBytes, numErr, numB1, numB2)

	inter := byteArrayToEncoding(interData) + byteArrayToEncoding(interError)
	return padRight(inter, len(inter)+blockInfo[version][6])
}

func drawQR(canvas [][]*Cell, inter string, version int) [][]*Cell {
	placeFinderPatterns(canvas, version)
	placeSeperators(canvas)
	placeAlignmentPatterns(canvas, version)
	drawTimingPattern(canvas, version)
	drawDarkModule(canvas, version)
	reserveFormatInformationArea(canvas, version)

	if version >= 7 {
		reserveVersionInformationData(canvas, version)
	}

	drawDataBits(canvas, inter, version)
	winner, mask := dataMasking(canvas)
	drawFormatInformationString(winner, mask)

	if version >= 7 {
		drawVersionInformationString(winner, version)
	}
	return winner
}

func upperLowerBorder(length int) string {
	border := ""
	for i := 0; i < length+2; i++ {
		border += white
	}
	return border + "\n"
}

// Print QR-Code to terminal
func (qr *QR) OutputTerminal() {
	length := len(qr.canvas)
	output := upperLowerBorder(length)

	for i := 0; i < length; i++ {
		output += white
		for j := 0; j < length; j++ {
			if qr.canvas[i][j].color == 0 {
				output += white
			} else {
				output += black
			}
		}
		output += white + "\n"
	}
	fmt.Println(output + upperLowerBorder(length))
}

func NewQR(data string) (*QR, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("Empty data input.")
	}

	mode := getMode(data)
	version := getVersion(data, mode)

	size = ((version - 1) * 4) + 21

	encoding := assembleEncoding(data, length, mode, version)
	inter := interleaveEncoding(encoding, version)
	winner := drawQR(initCanvas(size), inter, version)

	return &QR{data: data, length: length, mode: mode, version: version,
		encoding: encoding, canvas: winner}, nil
}

func main() {
	in := os.Args[1]
	qr, _ := NewQR(in)

	qr.OutputTerminal()
}
