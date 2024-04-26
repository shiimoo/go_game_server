package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

/*
DES 对称加密组件
*/

// 加密
func EncryptDES(key string, src []byte) ([]byte, error) {
	keybs := amendKey(key)
	block, err := des.NewTripleDESCipher(keybs) // 创建加密块
	if err != nil {
		// todo 创建失败错误码
		return nil, err
	}
	length := block.BlockSize()
	//填充最后一组数据
	src = paddingData(src, length)
	//创建cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, keybs[:length])
	//创建切片，用于存储加密后的数据
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return dst, nil
}

// 解密
func DecryptDES(key string, src []byte) ([]byte, error) {
	keybs := amendKey(key)
	//创建解密块
	block, err := des.NewTripleDESCipher(keybs)
	if nil != err {
		// todo 创建失败错误码
		return nil, err
	}
	//创建cbc解密模式
	blockMode := cipher.NewCBCDecrypter(block, keybs[:block.BlockSize()])
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return unPaddingData(dst), nil
}

// ------------------字节加工方法------------------ //

// 填充补全 数据字节
func paddingData(bs []byte, blockSize int) []byte {
	//求出最后一个分组需要填充的字节数，至少填充1，整除默认填充一个单位
	padding := blockSize - len(bs)%blockSize
	// 原始字符字节流+填充字节，填充数量作字节数据
	return append(bs, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// 取出冗余 数据字节
func unPaddingData(bs []byte) []byte {
	len := len(bs)
	return bs[:len-int(bs[len-1])]
}

const keyStandardLen = 24 // 密钥标准长度

// 修正 密钥字节
func amendKey(key string) []byte {
	bs := []byte(key)
	keyLength := len(bs)
	diffLen := keyStandardLen - keyLength
	if diffLen == 0 {
		return bs
	} else if diffLen > 0 { // 补全
		return append(bs, bytes.Repeat([]byte{byte(diffLen)}, diffLen)...)
	} else { // 截断
		return bs[:keyStandardLen]
	}
}
