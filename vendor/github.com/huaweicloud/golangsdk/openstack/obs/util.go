package obs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var regex = regexp.MustCompile("^[\u4e00-\u9fa5]$")

func StringToInt(value string, def int) int {
	ret, err := strconv.Atoi(value)
	if err != nil {
		ret = def
	}
	return ret
}

func StringToInt64(value string, def int64) int64 {
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		ret = def
	}
	return ret
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func FormatUtcNow(format string) string {
	return time.Now().UTC().Format(format)
}

func FormatUtcToRfc1123(t time.Time) string {
	ret := t.UTC().Format(time.RFC1123)
	return ret[:strings.LastIndex(ret, "UTC")] + "GMT"
}

func Md5(value []byte) []byte {
	m := md5.New()
	m.Write(value)
	return m.Sum(nil)
}

func HmacSha1(key, value []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}

func HmacSha256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}

func Base64Encode(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func HexMd5(value []byte) string {
	return Hex(Md5(value))
}

func Base64Md5(value []byte) string {
	return Base64Encode(Md5(value))
}

func Sha256Hash(value []byte) []byte {
	hash := sha256.New()
	hash.Write(value)
	return hash.Sum(nil)
}

func ParseXml(value []byte, result interface{}) error {
	if len(value) == 0 {
		return nil
	}
	return xml.Unmarshal(value, result)
}

func TransToXml(value interface{}) ([]byte, error) {
	if value == nil {
		return []byte{}, nil
	}
	return xml.Marshal(value)
}

func Hex(value []byte) string {
	return hex.EncodeToString(value)
}

func HexSha256(value []byte) string {
	return Hex(Sha256Hash(value))
}

func UrlDecode(value string) (string, error) {
	ret, err := url.QueryUnescape(value)
	if err == nil {
		return ret, nil
	}
	return "", err
}

func UrlEncode(value string, chineseOnly bool) string {
	if chineseOnly {
		values := make([]string, 0, len(value))
		for _, val := range value {
			_value := string(val)
			if regex.MatchString(_value) {
				_value = url.QueryEscape(_value)
			}
			values = append(values, _value)
		}
		return strings.Join(values, "")
	}
	return url.QueryEscape(value)
}
