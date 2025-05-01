package encryptor

import (
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type SkipEncryptionChecker func(message.Message) bool

// NewProtocolSkipChecker creates a checker that skips encryption for specific protocols
func NewProtocolSkipChecker(protocolsToSkip ...message.Protocol) SkipEncryptionChecker {
	skipMap := make(map[message.Protocol]bool)
	for _, protocol := range protocolsToSkip {
		skipMap[protocol] = true
	}

	return func(msg message.Message) bool {
		return skipMap[msg.GetProtocol()]
	}
}

type SkipperEncryptor struct {
	wrapped interfaces.Encryptor
	skip    SkipEncryptionChecker
}

func NewSkipperEncryptor(wrapped interfaces.Encryptor, skip SkipEncryptionChecker) *SkipperEncryptor {
	return &SkipperEncryptor{
		wrapped: wrapped,
		skip:    skip,
	}
}

func (e *SkipperEncryptor) SetSessionID(id types.EncryptionSessionID) {
	e.wrapped.SetSessionID(id)
}

func (e *SkipperEncryptor) SetKeys(encryptorKey, decryptorKey types.Key) error {
	return e.wrapped.SetKeys(encryptorKey, decryptorKey)
}

func (e *SkipperEncryptor) Encrypt(msg message.Message) (message.Message, error) {
	if e.skip(msg) {
		return msg, nil
	}

	return e.wrapped.Encrypt(msg)
}

func (e *SkipperEncryptor) Decrypt(msg message.Message) (message.Message, error) {
	if !e.skip(msg) {
		return e.wrapped.Decrypt(msg)
	}
	return msg, nil
}

func (e *SkipperEncryptor) Ready() bool {
	return e.wrapped.Ready()
}

func (e *SkipperEncryptor) Close() error {
	return e.wrapped.Close()
}
