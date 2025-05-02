package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sync"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type AES256Encryptor struct {
	encryptor cipher.AEAD
	decryptor cipher.AEAD
	sessionID types.EncryptionSessionID
	mux       sync.RWMutex
}

func (a *AES256Encryptor) SetKeys(encryptorKey, decryptorKey types.Key) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	// Setup encryption AEAD
	encBlock, err := aes.NewCipher(encryptorKey[:])
	if err != nil {
		return fmt.Errorf("%w: %v", encryptionerr.ErrInvalidKey, err)
	}

	encGCM, err := cipher.NewGCM(encBlock)
	if err != nil {
		return fmt.Errorf("%w: %v", encryptionerr.ErrInvalidKey, err)
	}

	// Setup decryption AEAD
	decBlock, err := aes.NewCipher(decryptorKey[:])
	if err != nil {
		return fmt.Errorf("%w: %v", encryptionerr.ErrInvalidKey, err)
	}

	decGCM, err := cipher.NewGCM(decBlock)
	if err != nil {
		return fmt.Errorf("%w: %v", encryptionerr.ErrInvalidKey, err)
	}

	a.encryptor = encGCM
	a.decryptor = decGCM

	return nil
}

func (a *AES256Encryptor) Encrypt(msg message.Message) (message.Message, error) {
	a.mux.Lock()
	defer a.mux.Unlock()

	nonce := types.Nonce{}

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	protocol := msg.GetProtocol()

	data, err := msg.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	encryptedData := a.encryptor.Seal(nil, nonce[:], data, a.sessionID[:])

	return NewEncryptedMessage(encryptedData, protocol, nonce, a.sessionID)
}

func (a *AES256Encryptor) Decrypt(msg message.Message) (message.Message, error) {
	m, ok := msg.(*EncryptedMessage)
	if !ok {
		return nil, encryptionerr.ErrInvalidInterceptor // JUST TO BE SURE
	}

	data, err := a.decryptor.Open(nil, m.Nonce[:], m.NextPayload, a.sessionID[:])
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	m.NextPayload = data

	return m, nil
}

func (a *AES256Encryptor) SetSessionID(id types.EncryptionSessionID) {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.sessionID = id
}

func (a *AES256Encryptor) Ready() bool {
	a.mux.RLock()
	defer a.mux.RUnlock()

	return a.encryptor != nil && a.decryptor != nil
}

func (a *AES256Encryptor) Close() error {
	// TODO implement me
	panic("implement me")
}
