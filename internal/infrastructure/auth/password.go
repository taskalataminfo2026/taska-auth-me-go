package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2")

var (
	// ErrInvalidHash es retornado cuando el hash proporcionado no es válido
	ErrInvalidHash = errors.New("el hash proporcionado no es válido")
	// ErrIncompatibleVersion es retornado cuando la versión de Argon2 no es compatible
	ErrIncompatibleVersion = errors.New("versión de Argon2 incompatible")
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// PasswordHasher implementa funciones para hashear y verificar contraseñas
// usando el algoritmo Argon2id, que es resistente a ataques de fuerza bruta
// y es considerado uno de los algoritmos más seguros actualmente disponibles.
type PasswordHasher struct {
	params *params
}

// NewPasswordHasher crea una nueva instancia de PasswordHasher con parámetros seguros por defecto
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		params: &params{
			memory:      64 * 1024, // 64 MB
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		},
	}
}

// GenerateHash genera un hash seguro de la contraseña usando Argon2id
func (h *PasswordHasher) GenerateHash(password string) (string, error) {
	// Generar una sal aleatoria
	salt := make([]byte, h.params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("error al generar la sal: %v", err)
	}

	// Generar el hash usando Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.iterations,
		h.params.memory,
		h.params.parallelism,
		h.params.keyLength,
	)

	// Codificar el hash y la sal en un string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Formato: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.params.memory,
		h.params.iterations,
		h.params.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Verify compara una contraseña con un hash existente
func (h *PasswordHasher) Verify(password, encodedHash string) (bool, error) {
	// Extraer los parámetros del hash codificado
	p, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generar el hash de la contraseña proporcionada con los mismos parámetros
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// Comparar los hashes de manera segura contra ataques de tiempo
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// decodeHash extrae los parámetros, la sal y el hash de un hash codificado
func (h *PasswordHasher) decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	// Verificar el formato del hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Verificar el identificador del algoritmo
	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("solo se soporta Argon2id")
	}

	// Verificar la versión
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Extraer parámetros
	p := &params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decodificar la sal
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Decodificar el hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
