package obs

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func (obsClient ObsClient) doAuthTemporary(method, bucketName, objectKey string, params map[string]string,
	headers map[string][]string, expires int64) (requestUrl string, err error) {

	isV4 := obsClient.conf.signature == SignatureV4

	requestUrl, canonicalizedUrl := obsClient.conf.formatUrls(bucketName, objectKey, params)
	parsedRequestUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	encodeHeaders(headers)

	hostName := parsedRequestUrl.Host

	skipAuth := obsClient.prepareHeaders(headers, hostName, isV4)

	if !skipAuth {
		if isV4 {
			date, _ := time.Parse(RFC1123_FORMAT, headers[HEADER_DATE_CAMEL][0])
			shortDate := date.Format(SHORT_DATE_FORMAT)
			longDate := date.Format(LONG_DATE_FORMAT)

			signedHeaders, _headers := getSignedHeaders(headers)

			credential, scope := getCredential(obsClient.conf.securityProvider.ak, obsClient.conf.region, shortDate)
			params[PARAM_ALGORITHM_AMZ_CAMEL] = V4_HASH_PREFIX
			params[PARAM_CREDENTIAL_AMZ_CAMEL] = credential
			params[PARAM_DATE_AMZ_CAMEL] = longDate
			params[PARAM_EXPIRES_AMZ_CAMEL] = Int64ToString(expires)
			params[PARAM_SIGNEDHEADERS_AMZ_CAMEL] = strings.Join(signedHeaders, ";")

			requestUrl, canonicalizedUrl = obsClient.conf.formatUrls(bucketName, objectKey, params)
			parsedRequestUrl, _ = url.Parse(requestUrl)
			stringToSign := getV4StringToSign(method, canonicalizedUrl, parsedRequestUrl.RawQuery, scope, longDate, UNSIGNED_PAYLOAD, signedHeaders, _headers)
			signature := getSignature(stringToSign, obsClient.conf.securityProvider.sk, obsClient.conf.region, shortDate)

			requestUrl += fmt.Sprintf("&%s=%s", PARAM_SIGNATURE_AMZ_CAMEL, UrlEncode(signature, false))

		} else {
			originDate := headers[HEADER_DATE_CAMEL][0]
			date, _ := time.Parse(RFC1123_FORMAT, originDate)
			expires += date.Unix()
			headers[HEADER_DATE_CAMEL] = []string{Int64ToString(expires)}

			stringToSign := getV2StringToSign(method, canonicalizedUrl, headers)
			signature := UrlEncode(Base64Encode(HmacSha1([]byte(obsClient.conf.securityProvider.sk), []byte(stringToSign))), false)
			if strings.Index(requestUrl, "?") < 0 {
				requestUrl += "?"
			} else {
				requestUrl += "&"
			}
			headers[HEADER_DATE_CAMEL] = []string{originDate}
			requestUrl += fmt.Sprintf("AWSAccessKeyId=%s&Expires=%d&Signature=%s", UrlEncode(obsClient.conf.securityProvider.ak, false),
				expires, signature)
		}
	}
	return
}

func (obsClient ObsClient) prepareHeaders(headers map[string][]string, hostName string, isV4 bool) bool {
	headers[HEADER_HOST_CAMEL] = []string{hostName}
	if date, ok := headers[HEADER_DATE_AMZ]; ok {
		flag := false
		if len(date) == 1 {
			if isV4 {
				if t, err := time.Parse(LONG_DATE_FORMAT, date[0]); err == nil {
					headers[HEADER_DATE_CAMEL] = []string{FormatUtcToRfc1123(t)}
					flag = true
				}
			} else {
				if strings.HasSuffix(date[0], "GMT") {
					headers[HEADER_DATE_CAMEL] = []string{date[0]}
					flag = true
				}
			}
		}
		if !flag {
			delete(headers, HEADER_DATE_AMZ)
		}
	}

	if _, ok := headers[HEADER_DATE_CAMEL]; !ok {
		headers[HEADER_DATE_CAMEL] = []string{FormatUtcToRfc1123(time.Now().UTC())}
	}

	if obsClient.conf.securityProvider == nil || obsClient.conf.securityProvider.ak == "" || obsClient.conf.securityProvider.sk == "" {
		doLog(LEVEL_WARN, "No ak/sk provided, skip to construct authorization")
		return true
	}

	if obsClient.conf.securityProvider.securityToken != "" {
		headers[HEADER_STS_TOKEN_AMZ] = []string{obsClient.conf.securityProvider.securityToken}
	}
	return false
}

func (obsClient ObsClient) doAuth(method, bucketName, objectKey string, params map[string]string,
	headers map[string][]string, hostName string) (requestUrl string, err error) {

	requestUrl, canonicalizedUrl := obsClient.conf.formatUrls(bucketName, objectKey, params)
	parsedRequestUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	encodeHeaders(headers)

	if hostName == "" {
		hostName = parsedRequestUrl.Host
	}

	isV4 := obsClient.conf.signature == SignatureV4

	skipAuth := obsClient.prepareHeaders(headers, hostName, isV4)

	if !skipAuth {
		if isV4 {
			headers[HEADER_CONTENT_SHA256_AMZ] = []string{EMPTY_CONTENT_SHA256}
			err = obsClient.v4Auth(method, canonicalizedUrl, parsedRequestUrl.RawQuery, headers)
		} else {
			err = obsClient.v2Auth(method, canonicalizedUrl, headers)
		}
	}
	return
}

func encodeHeaders(headers map[string][]string) {
	for key, values := range headers {
		for index, value := range values {
			values[index] = UrlEncode(value, true)
		}
		headers[key] = values
	}
}

func attachHeaders(headers map[string][]string) string {
	length := len(headers)
	_headers := make(map[string][]string, length)
	keys := make([]string, 0, length)

	for key, value := range headers {
		_key := strings.ToLower(strings.TrimSpace(key))
		if _key != "" {
			if _key == "content-md5" || _key == "content-type" || _key == "date" || strings.HasPrefix(_key, HEADER_PREFIX) {
				keys = append(keys, _key)
				_headers[_key] = value
			}
		} else {
			delete(headers, key)
		}
	}

	for _, interestedHeader := range interested_headers {
		if _, ok := _headers[interestedHeader]; !ok {
			_headers[interestedHeader] = []string{""}
			keys = append(keys, interestedHeader)
		}
	}

	sort.Strings(keys)

	stringToSign := make([]string, 0, len(keys))
	for _, key := range keys {
		var value string
		if strings.HasPrefix(key, HEADER_PREFIX) {
			if strings.HasPrefix(key, HEADER_PREFIX_META) {
				for index, v := range _headers[key] {
					value += strings.TrimSpace(v)
					if index != len(_headers[key])-1 {
						value += ","
					}
				}
			} else {
				value = strings.Join(_headers[key], ",")
			}
			value = fmt.Sprintf("%s:%s", key, value)
		} else {
			value = strings.Join(_headers[key], ",")
		}
		stringToSign = append(stringToSign, value)
	}
	return strings.Join(stringToSign, "\n")
}

func getV2StringToSign(method, canonicalizedUrl string, headers map[string][]string) string {
	stringToSign := strings.Join([]string{method, "\n", attachHeaders(headers), "\n", canonicalizedUrl}, "")
	doLog(LEVEL_DEBUG, "The v2 auth stringToSign:\n%s", stringToSign)
	return stringToSign
}

func (obsClient ObsClient) v2Auth(method, canonicalizedUrl string, headers map[string][]string) error {
	stringToSign := getV2StringToSign(method, canonicalizedUrl, headers)
	signature := Base64Encode(HmacSha1([]byte(obsClient.conf.securityProvider.sk), []byte(stringToSign)))

	headers[HEADER_AUTH_CAMEL] = []string{fmt.Sprintf("%s %s:%s", V2_HASH_PREFIX, obsClient.conf.securityProvider.ak, signature)}
	return nil
}

func getScope(region, shortDate string) string {
	return fmt.Sprintf("%s/%s/%s/%s", shortDate, region, V4_SERVICE_NAME, V4_SERVICE_SUFFIX)
}

func getCredential(ak, region, shortDate string) (string, string) {
	scope := getScope(region, shortDate)
	return fmt.Sprintf("%s/%s", ak, scope), scope
}

func getV4StringToSign(method, canonicalizedUrl, queryUrl, scope, longDate, payload string, signedHeaders []string, headers map[string][]string) string {

	canonicalRequest := make([]string, 0, 10+len(signedHeaders)*4)
	canonicalRequest = append(canonicalRequest, method)
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, canonicalizedUrl)
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, queryUrl)
	canonicalRequest = append(canonicalRequest, "\n")

	for _, signedHeader := range signedHeaders {
		values, _ := headers[signedHeader]
		for _, value := range values {
			canonicalRequest = append(canonicalRequest, signedHeader)
			canonicalRequest = append(canonicalRequest, ":")
			canonicalRequest = append(canonicalRequest, value)
			canonicalRequest = append(canonicalRequest, "\n")
		}
	}
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, strings.Join(signedHeaders, ";"))
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, payload)

	_canonicalRequest := strings.Join(canonicalRequest, "")
	doLog(LEVEL_DEBUG, "The v4 auth canonicalRequest:\n%s", _canonicalRequest)

	stringToSign := make([]string, 0, 7)
	stringToSign = append(stringToSign, V4_HASH_PREFIX)
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, longDate)
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, scope)
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, HexSha256([]byte(_canonicalRequest)))

	_stringToSign := strings.Join(stringToSign, "")

	doLog(LEVEL_DEBUG, "The v4 auth stringToSign:\n%s", _stringToSign)
	return _stringToSign
}

func getSignedHeaders(headers map[string][]string) ([]string, map[string][]string) {
	length := len(headers)
	_headers := make(map[string][]string, length)
	signedHeaders := make([]string, 0, length)
	for key, value := range headers {
		_key := strings.ToLower(strings.TrimSpace(key))
		if _key != "" {
			signedHeaders = append(signedHeaders, _key)
			_headers[_key] = value
		} else {
			delete(headers, key)
		}
	}
	sort.Strings(signedHeaders)
	return signedHeaders, _headers
}

func getSignature(stringToSign, sk, region, shortDate string) string {
	key := HmacSha256([]byte(V4_HASH_PRE+sk), []byte(shortDate))
	key = HmacSha256(key, []byte(region))
	key = HmacSha256(key, []byte(V4_SERVICE_NAME))
	key = HmacSha256(key, []byte(V4_SERVICE_SUFFIX))
	return Hex(HmacSha256(key, []byte(stringToSign)))
}

func (obsClient ObsClient) v4Auth(method, canonicalizedUrl, queryUrl string, headers map[string][]string) error {
	t, err := time.Parse(RFC1123_FORMAT, headers[HEADER_DATE_CAMEL][0])
	if err != nil {
		t = time.Now().UTC()
	}
	shortDate := t.Format(SHORT_DATE_FORMAT)
	longDate := t.Format(LONG_DATE_FORMAT)

	signedHeaders, _headers := getSignedHeaders(headers)

	credential, scope := getCredential(obsClient.conf.securityProvider.ak, obsClient.conf.region, shortDate)

	stringToSign := getV4StringToSign(method, canonicalizedUrl, queryUrl, scope, longDate, EMPTY_CONTENT_SHA256, signedHeaders, _headers)

	signature := getSignature(stringToSign, obsClient.conf.securityProvider.sk, obsClient.conf.region, shortDate)
	headers[HEADER_AUTH_CAMEL] = []string{fmt.Sprintf("%s Credential=%s,SignedHeaders=%s,Signature=%s", V4_HASH_PREFIX, credential, strings.Join(signedHeaders, ";"), signature)}
	return nil
}
