// Package aes256gcm contains the libsodium bindings for AES256-GCM.
package aes256gcm

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import "github.com/romanserikov/libsodium-go/support"

// Sodium should always be initialised
func init() {
	C.sodium_init()
}

// Sizes of nonces, key and mac.
const (
	KeyBytes   int = C.crypto_aead_aes256gcm_KEYBYTES  // Size of a secret key in bytes
	NSecBytes  int = C.crypto_aead_aes256gcm_NSECBYTES // Size of a secret nonce in bytes
	NonceBytes int = C.crypto_aead_aes256gcm_NPUBBYTES // Size of a nonce in bytes
	ABytes     int = C.crypto_aead_aes256gcm_ABYTES    // Size of an authentication tag in bytes
)

// IsAvailable returns true if AES256 is available on the current CPU
func IsAvailable() bool {
	return C.crypto_aead_aes256gcm_is_available() != 0
}

// GenerateKey generates a secret key
func GenerateKey() *[KeyBytes]byte {
	k := new([KeyBytes]byte)
	C.crypto_aead_aes256gcm_keygen((*C.uchar)(&k[0]))
	return k
}

// Encrypt a message `m` with additional data `ad` using a nonce `npub` and a secret key `k`.
// A ciphertext (including authentication tag) and encryption status are returned.
func Encrypt(m, ad []byte, nonce *[NonceBytes]byte, k *[KeyBytes]byte) (c []byte) {
	support.NilPanic(k == nil, "secret key")
	support.NilPanic(nonce == nil, "nonce")

	c = make([]byte, len(m)+ABytes)

	C.crypto_aead_aes256gcm_encrypt(
		(*C.uchar)(support.BytePointer(c)),
		(*C.ulonglong)(nil),
		(*C.uchar)(support.BytePointer(m)),
		(C.ulonglong)(len(m)),
		(*C.uchar)(support.BytePointer(ad)),
		(C.ulonglong)(len(ad)),
		(*C.uchar)(nil),
		(*C.uchar)(&nonce[0]),
		(*C.uchar)(&k[0]))

	return
}

// Decrypt and verify a ciphertext `c` using additional data `ad`, nonce `npub` and secret key `k`.
// Returns the decrypted message and verification status.
func Decrypt(c, ad []byte, nonce *[NonceBytes]byte, k *[KeyBytes]byte) (m []byte, err error) {
	support.NilPanic(k == nil, "secret key")
	support.NilPanic(nonce == nil, "nonce")
	support.CheckSizeMin(c, ABytes, "ciphertext")

	m = make([]byte, len(c)-ABytes)

	exit := C.crypto_aead_aes256gcm_decrypt(
		(*C.uchar)(support.BytePointer(m)),
		(*C.ulonglong)(nil),
		(*C.uchar)(nil),
		(*C.uchar)(&c[0]),
		(C.ulonglong)(len(c)),
		(*C.uchar)(support.BytePointer(ad)),
		(C.ulonglong)(len(ad)),
		(*C.uchar)(&nonce[0]),
		(*C.uchar)(&k[0]))

	if exit != 0 {
		err = &support.VerificationError{}
	}

	return
}

// EncryptDetached encrypts a message `m` with additional data `ad` using
// a nonce `npub` and a secret key `k`.
// A ciphertext, authentication tag and encryption status are returned.
func EncryptDetached(m, ad []byte, nonce *[NonceBytes]byte, k *[KeyBytes]byte) (c, mac []byte) {
	support.NilPanic(k == nil, "secret key")
	support.NilPanic(nonce == nil, "nonce")

	c = make([]byte, len(m))
	mac = make([]byte, ABytes)

	C.crypto_aead_aes256gcm_encrypt_detached(
		(*C.uchar)(support.BytePointer(c)),
		(*C.uchar)(&mac[0]),
		(*C.ulonglong)(nil),
		(*C.uchar)(support.BytePointer(m)),
		(C.ulonglong)(len(m)),
		(*C.uchar)(support.BytePointer(ad)),
		(C.ulonglong)(len(ad)),
		(*C.uchar)(nil),
		(*C.uchar)(&nonce[0]),
		(*C.uchar)(&k[0]))

	return
}

// DecryptDetached decrypts and verifies a ciphertext `c` with authentication tag `mac`
// using additional data `ad`, nonce `npub` and secret key `k`.
// Returns the decrypted message and verification status.
func DecryptDetached(c, mac, ad []byte, nonce *[NonceBytes]byte, k *[KeyBytes]byte) (m []byte, err error) {
	support.NilPanic(k == nil, "secret key")
	support.NilPanic(nonce == nil, "nonce")
	support.CheckSize(mac, ABytes, "mac")

	m = make([]byte, len(c))

	exit := C.crypto_aead_aes256gcm_decrypt_detached(
		(*C.uchar)(support.BytePointer(m)),
		(*C.uchar)(nil),
		(*C.uchar)(support.BytePointer(c)),
		(C.ulonglong)(len(c)),
		(*C.uchar)(&mac[0]),
		(*C.uchar)(support.BytePointer(ad)),
		(C.ulonglong)(len(ad)),
		(*C.uchar)(&nonce[0]),
		(*C.uchar)(&k[0]))

	if exit != 0 {
		err = &support.VerificationError{}
	}

	return
}
