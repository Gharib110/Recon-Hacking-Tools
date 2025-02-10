package ntlmssp

import (
	"crypto/hmac"
	"crypto/md5"
	"golang.org/x/crypto/md4"
	"smb/smb/encoder"
	"strings"
)

func NtowFv1(pass string) []byte {
	hash := md4.New()
	hash.Write(encoder.ToUnicode(pass))
	return hash.Sum(nil)
}

func NtowFv2(pass, user, domain string) []byte {
	h := hmac.New(md5.New, NtowFv1(pass))
	h.Write(encoder.ToUnicode(strings.ToUpper(user) + domain))
	return h.Sum(nil)
}

func Lmowfv2(pass, user, domain string) []byte {
	return NtowFv2(pass, user, domain)
}

func ComputeResponseNTLMv2(nthash, lmhash, clientChallenge, serverChallenge, timestamp, serverName []byte) []byte {

	temp := []byte{1, 1}
	temp = append(temp, 0, 0, 0, 0, 0, 0)
	temp = append(temp, timestamp...)
	temp = append(temp, clientChallenge...)
	temp = append(temp, 0, 0, 0, 0)
	temp = append(temp, serverName...)
	temp = append(temp, 0, 0, 0, 0)

	h := hmac.New(md5.New, nthash)
	h.Write(append(serverChallenge, temp...))
	ntproof := h.Sum(nil)
	return append(ntproof, temp...)
}
