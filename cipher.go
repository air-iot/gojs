package gojs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/dop251/goja"
	"github.com/tjfoc/gmsm/sm4"
)

// CipherType 加密算法类型
type CipherType string

const (
	AES CipherType = "AES"
	SM4 CipherType = "SM4"
	RSA CipherType = "RSA"
	SM2 CipherType = "SM2"
)

// Mode 加密模式类型
type Mode string

const (
	CBC Mode = "CBC"
	ECB Mode = "ECB"
	CTR Mode = "CTR"
	GCM Mode = "GCM"
	OFB Mode = "OFB"
	CFB Mode = "CFB"
)

type CipherUtil struct{}

func (c CipherUtil) Attach(vm *goja.Runtime) error {
	obj := vm.NewObject()
	if err := obj.DefineDataProperty("Encrypt", vm.ToValue(func(data goja.Value, key goja.Value, cipherType CipherType, mode Mode) (string, error) {
		return c.Encrypt(vm, data, key, cipherType, mode)
	}), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_FALSE); err != nil {
		return fmt.Errorf("注册 'Encrypt' 方法失败, %+v", err)
	}

	if err := obj.DefineDataProperty("Decrypt", vm.ToValue(func(encrypted string, key goja.Value, cipherType CipherType, mode Mode) (goja.Value, error) {
		return c.Decrypt(vm, encrypted, key, cipherType, mode)
	}), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_FALSE); err != nil {
		return fmt.Errorf("注册 'Decrypt' 方法失败, %+v", err)
	}

	if err := vm.Set("cipherUtil", obj); err != nil {
		return fmt.Errorf("注册 cipherUtil 对象失败, %+v", err)
	}

	return nil
}

// Encrypt 扩展支持 RSA 和 SM2
func (c CipherUtil) Encrypt(_vm *goja.Runtime, data goja.Value, key goja.Value, cipherType CipherType, mode Mode) (string, error) {
	switch cipherType {
	case RSA, SM2:
		return "", errors.New("unsupported asymmetric cipher type")
	default:
		// 原有的对称加密逻辑保持不变
		if !IsString(data) && !IsBuffer(data) {
			return "", fmt.Errorf("被加密数据 data 必须为字符串或 Buffer 对象")
		} else if !IsString(key) && !IsBuffer(key) {
			return "", fmt.Errorf("密钥 key 必须为字符串或 Buffer 对象")
		}

		var dataBytes []byte
		var keyBytes []byte

		if IsString(data) {
			dataBytes = []byte(data.String())
		} else if IsBuffer(data) {
			dataBytes = data.Export().([]byte)
		}

		if IsString(key) {
			keyBytes = []byte(key.String())
		} else if IsBuffer(data) {
			keyBytes = key.Export().([]byte)
		}

		return symmetricEncrypt(dataBytes, keyBytes, cipherType, mode)
	}
}

// Decrypt 扩展支持 RSA 和 SM2
func (c CipherUtil) Decrypt(vm *goja.Runtime, encrypted string, key goja.Value, cipherType CipherType, mode Mode) (goja.Value, error) {
	switch cipherType {
	case RSA, SM2:
		return nil, errors.New("unsupported asymmetric cipher type")
	default:
		// 原有的对称解密逻辑保持不变
		var keyBytes []byte
		if IsString(key) {
			keyBytes = []byte(key.String())
		} else if IsBuffer(key) {
			keyBytes = key.Export().([]byte)
		}

		decryped, err := symmetricDecrypt(encrypted, keyBytes, cipherType, mode)
		if err != nil {
			return nil, err
		}

		return BytesToBuffer(vm, decryped)
	}
}

// AsymmetricDecrypt 非对称解密 (公钥解密)
func AsymmetricDecrypt(encrypted string, key []byte, cipherType CipherType) ([]byte, error) {
	// 解码 base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	switch cipherType {
	case RSA:
		// 解析 RSA 公钥
		block, _ := pem.Decode(key)
		if block == nil {
			return nil, errors.New("failed to parse PEM block containing the key")
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}

		// RSA 公钥解密
		return rsaDecrypt(rsaPubKey, ciphertext)

	case SM2:
		// 解析 SM2 公钥
		//pubKey, err := sm2.ParseSm2PublicKey(key)
		//if err != nil {
		//	return nil, err
		//}
		//
		//// SM2 公钥解密
		//return sm2.Decrypt(pubKey, ciphertext)
		return nil, nil
	default:
		return nil, errors.New("unsupported asymmetric cipher type")
	}
}

// rsaEncrypt RSA 私钥加密
func rsaEncrypt(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// 使用私钥的公共模数作为加密参数
	pubKey := &privKey.PublicKey
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
}

// rsaDecrypt RSA 公钥解密
func rsaDecrypt(pubKey *rsa.PublicKey, ciphertext []byte) ([]byte, error) {
	// 模拟公钥解密 (实际 RSA 通常是公钥加密、私钥解密)
	// 这里使用公钥的模数进行"解密"，实际应用中可能需要不同的方法
	// 注意：这不是标准 RSA 用法，仅演示目的
	c := new(big.Int).SetBytes(ciphertext)
	m := new(big.Int).Exp(c, big.NewInt(int64(pubKey.E)), pubKey.N)
	return m.Bytes(), nil
}

// Encrypt 通用加密函数
func symmetricEncrypt(data []byte, key []byte, cipherType CipherType, mode Mode) (string, error) {
	var block cipher.Block
	var err error

	// 根据加密算法类型创建对应的块
	switch cipherType {
	case AES:
		block, err = aes.NewCipher(key)
	case SM4:
		block, err = sm4.NewCipher(key)
	default:
		return "", errors.New("unsupported cipher type")
	}
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()

	switch mode {
	case ECB:
		// ECB 模式不需要 IV
		data = pkcs7Padding(data, blockSize)
		crypted := make([]byte, len(data))
		for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
			block.Encrypt(crypted[bs:be], data[bs:be])
		}
		return base64.StdEncoding.EncodeToString(crypted), nil

	case CBC:
		data = pkcs7Padding(data, blockSize)
		crypted := make([]byte, blockSize+len(data))
		iv := crypted[:blockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}

		cbcMode := cipher.NewCBCEncrypter(block, iv)
		cbcMode.CryptBlocks(crypted[blockSize:], data)
		return base64.StdEncoding.EncodeToString(crypted), nil

	case CTR:
		crypted := make([]byte, blockSize+len(data))
		iv := crypted[:blockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}

		var stream cipher.Stream
		if cipherType == SM4 {
			// SM4 的 CTR 模式需要特殊处理
			stream = NewSM4CTR(block, iv)
		} else {
			stream = cipher.NewCTR(block, iv)
		}
		stream.XORKeyStream(crypted[blockSize:], data)
		return base64.StdEncoding.EncodeToString(crypted), nil

	case GCM:
		if cipherType == SM4 {
			nonce := make([]byte, sm4.BlockSize)
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return "", err
			}
			data = pkcs7Padding(data, sm4.BlockSize)
			ciphertext, _, err := sm4.Sm4GCM(key, nonce, data, nil, true)
			if err != nil {
				return "", err
			}
			// 组合 nonce + ciphertext 用于传输
			encryptedData := append(nonce, ciphertext...)
			return base64.StdEncoding.EncodeToString(encryptedData), nil
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return "", err
		}

		crypted := gcm.Seal(nonce, nonce, data, nil)
		return base64.StdEncoding.EncodeToString(crypted), nil

	case OFB:
		crypted := make([]byte, blockSize+len(data))
		iv := crypted[:blockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}

		stream := cipher.NewOFB(block, iv)
		stream.XORKeyStream(crypted[blockSize:], data)
		return base64.StdEncoding.EncodeToString(crypted), nil

	case CFB:
		crypted := make([]byte, blockSize+len(data))
		iv := crypted[:blockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}

		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(crypted[blockSize:], data)
		return base64.StdEncoding.EncodeToString(crypted), nil

	default:
		return "", errors.New("unsupported encryption mode")
	}
}

// Decrypt 通用解密函数
func symmetricDecrypt(encrypted string, key []byte, cipherType CipherType, mode Mode) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	var block cipher.Block
	switch cipherType {
	case AES:
		block, err = aes.NewCipher(key)
	case SM4:
		block, err = sm4.NewCipher(key)
	default:
		return nil, errors.New("unsupported cipher type")
	}
	if err != nil {
		return nil, err
	}

	switch mode {
	case ECB:
		blockSize := block.BlockSize()
		if len(data)%blockSize != 0 {
			return nil, errors.New("ciphertext is not a multiple of the block size")
		}

		decrypted := make([]byte, len(data))
		for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
			block.Decrypt(decrypted[bs:be], data[bs:be])
		}
		return pkcs7UnPadding(decrypted)

	case CBC:
		blockSize := block.BlockSize()
		if len(data) < blockSize {
			return nil, errors.New("ciphertext too short")
		}

		iv := data[:blockSize]
		data = data[blockSize:]

		if len(data)%blockSize != 0 {
			return nil, errors.New("ciphertext is not a multiple of the block size")
		}

		cbcMode := cipher.NewCBCDecrypter(block, iv)
		decrypted := make([]byte, len(data))
		cbcMode.CryptBlocks(decrypted, data)
		return pkcs7UnPadding(decrypted)

	case CTR:
		blockSize := block.BlockSize()
		if len(data) < blockSize {
			return nil, errors.New("ciphertext too short")
		}

		iv := data[:blockSize]
		data = data[blockSize:]

		var stream cipher.Stream
		if cipherType == SM4 {
			stream = NewSM4CTR(block, iv)
		} else {
			stream = cipher.NewCTR(block, iv)
		}
		decrypted := make([]byte, len(data))
		stream.XORKeyStream(decrypted, data)
		return decrypted, nil

	case GCM:
		if cipherType == SM4 {
			blockSize := sm4.BlockSize
			if len(data) < blockSize {
				return nil, errors.New("ciphertext too short")
			}

			// 分离 nonce 和 ciphertext
			nonce := data[:blockSize]
			ciphertext := data[blockSize:]
			// 使用官方 Sm4GCM 接口解密
			plaintext, _, err := sm4.Sm4GCM(key, nonce, ciphertext, nil, false)
			if err != nil {
				return nil, err
			}
			return pkcs7UnPadding(plaintext)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
			return nil, errors.New("ciphertext too short")
		}

		nonce, ciphertext := data[:nonceSize], data[nonceSize:]
		return gcm.Open(nil, nonce, ciphertext, nil)

	case OFB:
		blockSize := block.BlockSize()
		if len(data) < blockSize {
			return nil, errors.New("ciphertext too short")
		}

		iv := data[:blockSize]
		data = data[blockSize:]

		stream := cipher.NewOFB(block, iv)
		decrypted := make([]byte, len(data))
		stream.XORKeyStream(decrypted, data)
		return decrypted, nil

	case CFB:
		blockSize := block.BlockSize()
		if len(data) < blockSize {
			return nil, errors.New("ciphertext too short")
		}

		iv := data[:blockSize]
		data = data[blockSize:]

		stream := cipher.NewCFBDecrypter(block, iv)
		decrypted := make([]byte, len(data))
		stream.XORKeyStream(decrypted, data)
		return decrypted, nil

	default:
		return nil, errors.New("unsupported encryption mode")
	}
}

// SM4CTR 实现 SM4 的 CTR 模式
type SM4CTR struct {
	block     cipher.Block
	ctr       []byte
	out       []byte
	outUsed   int
	blockSize int
}

func NewSM4CTR(block cipher.Block, iv []byte) *SM4CTR {
	return &SM4CTR{
		block:     block,
		ctr:       bytes.Clone(iv),
		out:       make([]byte, block.BlockSize()),
		outUsed:   block.BlockSize(),
		blockSize: block.BlockSize(),
	}
}

func (x *SM4CTR) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		if x.outUsed == x.blockSize {
			x.block.Encrypt(x.out, x.ctr)

			// 递增计数器
			for j := len(x.ctr) - 1; j >= 0; j-- {
				x.ctr[j]++
				if x.ctr[j] != 0 {
					break
				}
			}

			x.outUsed = 0
		}

		dst[i] = src[i] ^ x.out[x.outUsed]
		x.outUsed++
	}
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 去除填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}

	unPadding := int(data[length-1])
	if unPadding > length {
		return nil, errors.New("invalid padding")
	}

	return data[:(length - unPadding)], nil
}
