package numutils

func SwapInt(a, b *int) {
	*a += *b
	*b = *a - *b
	*a -= *b
}

func ExchanggeSymbol(a int) (ret int) {
	ret = ^a + 1
	return
}

func CheckOdd(a int) bool {
	if a%2 > 0 {
		return true
	}
	return false
}

func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}

/*
	calculate the number of 1 in a
*/
func CalcOneNum(a int) int {
	a = ((a & 0xAAAA) >> 1) + (a & 0x5555)
	a = ((a & 0xCCCC) >> 2) + (a & 0x3333)
	a = ((a & 0xF0F0) >> 4) + (a & 0x0F0F)
	a = ((a & 0xFF00) >> 8) + (a & 0x00FF)
	return a
}

/*
	reverse the order according to the int's binary byte
*/
func ByteReverse(a int) int {
	a = ((a & 0xAAAA) >> 1) | ((a & 0x5555) << 1)
	a = ((a & 0xCCCC) >> 2) | ((a & 0x3333) << 2)
	a = ((a & 0xF0F0) >> 4) | ((a & 0x0F0F) << 4)
	a = ((a & 0xFF00) >> 8) | ((a & 0x00FF) << 8)
	return a
}

// reverse digits of an interger
// such as 123 reverse to 321, -123 reverse to -321
// if reversed integer overflows, then return 0
func ReverseInt(x int) int {
	var result int
	for x != 0 {
		tail := x % 10
		newResult := result*10 + tail
		// check overflows
		if (newResult-tail)/10 != result {
			return 0
		}
		result = newResult
		x /= 10
	}
	return result
}

// check whether x is palindrome
func IsPalindrome(x int) bool {
	y, t := 0, x
	for t > 0 {
		y = y*10 + t%10
		t = t / 10
	}
	return x == y
}

var (
	romanDigits = [4][10]string{
		{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"},
		{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"},
		{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"},
		{"", "M", "MM", "MMM"},
	}
	romanBit = map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}
)

// convert integer to roman string
// number ranges from 1 to 3999
func IntToRoman(num int) string {
	var roman string
	roman += romanDigits[3][num/1000%10]
	roman += romanDigits[2][num/100%10]
	roman += romanDigits[1][num/10%10]
	roman += romanDigits[0][num%10]
	return roman
}

// convert roman string to integer
// number ranges from 1 to 3999
// I（1）、V（5）、X（10）、L（50）、C（100）、D（500）、 M（1000）
func RomanToInt(s string) int {
	var num int

	for i := 0; i < len(s)-1; i++ {
		if romanBit[s[i]] < romanBit[s[i+1]] {
			num -= romanBit[s[i]]
		} else {
			num += romanBit[s[i]]
		}
	}
	num += romanBit[s[len(s)-1]]
	return num
}
