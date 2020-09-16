package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		desc      string
		plainText string
		key       string
	}{
		{
			desc:      "Should success",
			plainText: "data",
			key:       "password",
		},
		{
			desc:      "Should success",
			plainText: "abcdef",
			key:       "password01",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			passphrase := CreateHash(tC.key)
			encryptedData, err := Encrypt(tC.plainText, passphrase)
			require.NoError(t, err)

			decryptedData, err := Decrypt(encryptedData, passphrase)

			require.NoError(t, err)
			assert.Equal(t, tC.plainText, decryptedData)
		})
	}
}
