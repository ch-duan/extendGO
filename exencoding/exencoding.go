package exencoding

import (
	"encoding/hex"
	"errors"
	"extend/exruntime"
	"fmt"
	"strconv"
)

func DecHex(n int64) string {
	if n < 0 {
		return ""
	}
	if n == 0 {
		return "0"
	}
	s := fmt.Sprintf("%x", n)
	return s
}

func EnOneByte(str string) byte {
	s, _ := strconv.ParseInt(str, 10, 64)
	val := DecHex(s)
	if len(val) == 1 {
		val = "0" + val
	} else if len(val) == 0 {
		val = "00"
	}
	result, _ := hex.DecodeString(val)
	return result[0]
}

func Float2Float(num float64) float64 {
	float_num, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", num), 64)
	return float_num
}

//FromHexChar converts a hex character into its value and a success flag.
func FromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

//Complete8Byte Complete 8 byte
func Complete8Byte(rawData []byte) (result []byte) {
	if len(rawData) >= 8 {
		result = rawData
		return
	}
	l := len(rawData)
	res := make([]byte, 8)
	for i := 0; i < 8; i++ {
		if i < l {
			res[i] = rawData[i]
		} else {
			if rawData[l-1]&0x80 == 0x80 {
				res[i] = 255
			} else {
				res[i] = 0
			}
		}
	}
	result = res
	return
}

//Complete8Byte Complete 8 byte
func Complete8ByteBigEndian(rawData []byte) (result []byte) {
	if len(rawData) >= 8 {
		result = rawData
		return
	}
	l := len(rawData)
	res := make([]byte, 8)
	for i, j := 7, l-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		res[i] = rawData[j]
	}
	result = res
	return
}

//Uint64ToByte params{b *[]byte,[]byte长度,isBigEndian,v uint64}
func Uint64ToByte(b *[]byte, length int, isBigEndian bool, v uint64) {
	if !isBigEndian {
		for i := 0; i < length; i++ {
			(*b)[i] = byte(v >> (i * 8))
		}
	} else {
		for i, j := 0, length-1; i < length && j >= 0; i, j = i+1, j-1 {
			(*b)[i] = byte(v >> (j * 8))
		}
	}
}

//Uint64ToDecString Uint64转换十进制字符串
func Uint64ToDecString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

//Uint64ToHexString Uint64转换十六进制字符串
func Uint64ToHexString(v uint64) string {
	return strconv.FormatUint(v, 16)
}

//ByteToUint64 byte根据大小端转换为uint64,params{b []byte,[]byte长度,isBigEndian},返回uint64
func ByteToUint64(b []byte, length int, isBigEndian bool) (result uint64) {
	result = 0
	if !isBigEndian {
		for i := 0; i < length; i++ {
			result |= uint64(b[i]) << (i * 8)
		}
	} else {
		for i, j := 0, length-1; i < length && j >= 0; i, j = i+1, j-1 {
			result |= uint64(b[i]) << (j * 8)
		}
	}
	return
}

//ByteToFloat64 byte根据大小端转换为uint64,params{b []byte,[]byte长度,isBigEndian},返回uint64
func ByteToFloat64(b []byte, length int, isBigEndian bool) float64 {
	return float64(ByteToUint64(b, length, isBigEndian))
}

//ByteToFloatString byte根据指定大小端转float字符串
func ByteToFloatString(bts []byte, isBigEndian bool, div, n float64) (string, error) {
	return DivToFloatString(float64(ByteToUint64(bts, len(bts), isBigEndian)), div, n)
}

func ByteToDecString(bts []byte, isBigEndian bool) string {
	return strconv.FormatInt(int64(ByteToUint64(bts, len(bts), isBigEndian)), 10)
}

//Float64ToString float64转成string并保留n位小数
func Float64ToString(f, n float64) string {
	return strconv.FormatFloat(f, 'f', int(n), 64)
}

func StringToByte(b *[]byte, s string, base int, isBigEndian bool, length int) error {
	v, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return err
	}
	Uint64ToByte(b, length, isBigEndian, v)
	return nil
}

func DecStringToHex8Byte(s string) ([]byte, error) {
	//先将数字转float64
	a, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New(exruntime.RunFuncName() + err.Error())
	}
	//再转成2字节无符号数然后转成字符串
	val := strconv.FormatUint(uint64(a), 16)
	if len(val)%2 != 0 {
		val = "0" + val
	}
	//字符串转byte
	res, err := hex.DecodeString(val)
	if err != nil {
		return nil, errors.New(exruntime.RunFuncName() + err.Error())
	}
	ReserveByteSlice(&res)
	return res, nil
}

//Divide 除法f/div
func Divide(f, div float64) (float64, error) {
	if div == 0 {
		return 0, errors.New(exruntime.RunFuncName() + " div is 0")
	}
	return f / div, nil
}

//Mutiply 乘法f/multiplier
func Mutiply(f, multiplier float64) float64 {
	return f * multiplier
}

//DivToByte f/div后转换为指定长度和大小端的[]byte
func DivToByte(res *[]byte, isBigEndian bool, length int, f, div float64) error {
	v, err := Divide(f, div)
	if err != nil {
		return err
	}
	Uint64ToByte(res, length, isBigEndian, uint64(int(v)))
	return nil
}

//DivToFloatByte f/div后转换为指定小数点的string
func DivToFloatString(f, div, n float64) (string, error) {
	v, err := Divide(f, div)
	if err != nil {
		return "", err
	}
	return Float64ToString(v, n), nil
}

//MultiplyToFloatByte f*div后转换为指定小数点的string
func MultiplyToFloatString(f, multiplier, n float64) string {
	return Float64ToString(Mutiply(f, multiplier), n)
}

//Multiply f*multiplier后转换为指定长度和大小端的[]byte
func MultiplyToByte(res *[]byte, isBigEndian bool, length int, f, multiplier float64) {
	Uint64ToByte(res, length, isBigEndian, uint64(int(Mutiply(f, multiplier))))
}

func ReserveByteSlice(b *[]byte) {
	if len(*b) <= 1 {
		return
	}
	length := len(*b)
	for i := 0; i < length/2; i++ {
		temp := (*b)[length-1-i]
		(*b)[length-1-i] = (*b)[i]
		(*b)[i] = temp
	}
}
