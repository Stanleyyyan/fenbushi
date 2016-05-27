package libkademlia

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	mathrand "math/rand"
	"time"
	"sss"
	"fmt"
	//"strconv"
	"bytes"
)

type VanashingDataObject struct {
	AccessKey  int64
	Ciphertext []byte
	NumberKeys byte
	Threshold  byte
}

func GenerateRandomCryptoKey() (ret []byte) {
	for i := 0; i < 32; i++ {
		ret = append(ret, uint8(mathrand.Intn(256)))
	}
	return
}

func GenerateRandomAccessKey() (accessKey int64) {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	accessKey = r.Int63()
	return
}

func CalculateSharedKeyLocations(accessKey int64, count int64) (ids []ID) {
	r := mathrand.New(mathrand.NewSource(accessKey))
	ids = make([]ID, count)
	for i := int64(0); i < count; i++ {
		for j := 0; j < IDBytes; j++ {
			ids[i][j] = uint8(r.Intn(256))
		}
	}
	return
}

func encrypt(key []byte, text []byte) (ciphertext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext = make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return
}

func decrypt(key []byte, ciphertext []byte) (text []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext is not long enough")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func (k *Kademlia) VanishData(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	K := GenerateRandomCryptoKey()//[]byte
	Ciphertext := encrypt(K, data)
	Shares, err := sss.Split(numberKeys, threshold, K)
	if err != nil {
		fmt.Println("AA")
	}
	L := GenerateRandomAccessKey()//int64
	ids := CalculateSharedKeyLocations(L, int64(numberKeys))
	// for i := 0; i < len(Ciphertext); i++ {
		
	// }
	i := 0
	for key, value := range Shares {
		all := append([]byte{key}, value...)
		k.DoIterativeStore(ids[i], all)
		i++
	}

	fmt.Println("Length shares is :", len(Shares))
	fmt.Println("ids is :", len(ids))


	vdo.AccessKey = L
	vdo.Ciphertext = Ciphertext
	vdo.NumberKeys = numberKeys
	vdo.Threshold = threshold
	return vdo
}

func (k *Kademlia) UnvanishData(vdo VanashingDataObject) (data []byte) {
	ids := CalculateSharedKeyLocations(vdo.AccessKey, int64(vdo.NumberKeys))
	fmt.Println("AccessKey IS ", vdo.AccessKey)
	fmt.Println("Ciphertext IS ", vdo.Ciphertext)
	fmt.Println("NumberKeys IS ", vdo.NumberKeys)
	fmt.Println("Threshold IS ", vdo.Threshold)
	var Shares map[byte][]byte
	Shares = make(map[byte][]byte)
	count := 0;
	threshold := int(vdo.Threshold)
	fmt.Println("ID LENGTH IS :",len(ids))
	for i := 0; i<len(ids); i++ {
		if count == threshold {
			break
		}
		all, _ := k.DoIterativeFindValue(ids[i])
		fmt.Println("find shars is :", all)
		if !bytes.Equal(all, []byte("")) {
			Shares[all[0]] = all[1:]
			count++
		}

	}
	K := sss.Combine(Shares)
	ciphertext := decrypt(K, vdo.Ciphertext)
	return ciphertext
}


