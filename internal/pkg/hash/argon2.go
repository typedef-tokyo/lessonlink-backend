package hash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	// ソルトの生成
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Argon2id パラメータの設定
	time := uint32(1)
	memory := uint32(64 * 1024) // 64 MB
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// エンコードして保存しやすい形式にする
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// フォーマット: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, threads, b64Salt, b64Hash)

	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, log.WrapErrorWithStackTrace(errors.New("invalid hash format"))
	}

	// パラメータの解析
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(hash))
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// 定数時間比較
	if len(computedHash) != len(hash) {
		return false, nil
	}
	var result uint8
	for i := 0; i < len(computedHash); i++ {
		result |= computedHash[i] ^ hash[i]
	}
	return result == 0, nil
}
