package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	bytes2 "github.com/labstack/gommon/bytes"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func InterfaceToString(value interface{}) string {

	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case *float64:
		ft := value.(*float64)
		if ft != nil {
			key = strconv.FormatFloat(*ft, 'f', -1, 64)
		} else {
			key = ""
		}
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case *float32:
		ft := value.(*float32)
		if ft != nil {
			key = strconv.FormatFloat(float64(*ft), 'f', -1, 64)
		} else {
			key = ""
		}
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case *int:
		it := value.(*int)
		if it != nil {
			key = strconv.Itoa(*it)
		}
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case *uint:
		it := value.(*uint)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case *int8:
		it := value.(*int8)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case *uint8:
		it := value.(*uint8)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case *int16:
		it := value.(*int16)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case *uint16:
		it := value.(*uint16)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case *int32:
		it := value.(*int32)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case *uint32:
		it := value.(*uint32)
		if it != nil {
			key = strconv.Itoa(int(*it))
		}
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case *int64:
		it := value.(*int64)
		if it != nil {
			key = strconv.FormatInt(*it, 10)
		}
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case *uint64:
		it := value.(*uint64)
		if it != nil {
			key = strconv.FormatUint(*it, 10)
		}
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case error:
		key = value.(error).Error()
	case *string:
		t := value.(*string)
		if t != nil {
			key = *t
		}
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func DPrintln(val ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	loc := ""
	if ok {
		loc = fmt.Sprintf("[INFO ts=%s %s:%d] ", time.Now().Format(time.RFC3339), path.Base(file), line)
	}
	var elems []string
	for k, _ := range val {
		elems = append(elems, InterfaceToString(val[k]))
	}
	// INFO ts=2025-05-18T03:46:03Z caller=biz/zoom_meeting_sms_notice_job.go:99
	fmt.Println(loc + strings.Join(elems, " "))
}

// SqlRowsToRow 获取一行数据
func SqlRowsToRow(rows *sql.Rows) (columns []string, row map[string]interface{}, err error) {
	columns, err = rows.Columns()
	if err != nil {
		return columns, row, err
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return columns, row, err
		}

		ret := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = nil
			} else {
				switch val := (*scanArgs[i].(*interface{})).(type) {
				case byte:
					ret[columns[i]] = val
					break
				case []byte:
					v := string(val)
					switch v {
					case "\x00": // 处理数据类型为bit的情况
						ret[columns[i]] = 0
					case "\x01": // 处理数据类型为bit的情况
						ret[columns[i]] = 1
					default:
						ret[columns[i]] = v
						break
					}
					break
				case time.Time:
					if val.IsZero() {
						ret[columns[i]] = nil
					} else {
						ret[columns[i]] = val.Format("2006-01-02 15:04:05")
					}
					break
				default:
					ret[columns[i]] = val
				}
			}
		}
		row = ret
		break
	}
	return columns, row, nil
}

func SqlRowsTrans(rows *sql.Rows) (columns []string, list []map[string]interface{}, err error) {

	columns, err = rows.Columns()
	if err != nil {
		return columns, list, err
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return columns, list, err
		}

		ret := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = nil
			} else {
				switch val := (*scanArgs[i].(*interface{})).(type) {
				case byte:
					ret[columns[i]] = val
					break
				case []byte:
					v := string(val)
					switch v {
					case "\x00": // 处理数据类型为bit的情况
						ret[columns[i]] = 0
					case "\x01": // 处理数据类型为bit的情况
						ret[columns[i]] = 1
					default:
						ret[columns[i]] = v
						break
					}
					break
				case time.Time:
					if val.IsZero() {
						ret[columns[i]] = nil
					} else {
						ret[columns[i]] = val.Format("2006-01-02 15:04:05")
					}
					break
				default:
					ret[columns[i]] = val
				}
			}
		}
		list = append(list, ret)
	}
	return columns, list, nil
}

/*
SqlRowsToMap

使用方法

	for sqlRows.Next() {
		_, row, err := lib.SqlRowsToMap(sqlRows)
	}
*/
func SqlRowsToMap(rows *sql.Rows) (columns []string, typeMap TypeMap, err error) {

	columns, err = rows.Columns()
	if err != nil {
		return columns, typeMap, err
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		return columns, typeMap, err
	}

	ret := make(map[string]interface{})
	for i, col := range values {
		if col == nil {
			ret[columns[i]] = nil
		} else {
			switch val := (*scanArgs[i].(*interface{})).(type) {
			case byte:
				ret[columns[i]] = val
				break
			case []byte:
				v := string(val)
				switch v {
				case "\x00": // 处理数据类型为bit的情况
					ret[columns[i]] = 0
				case "\x01": // 处理数据类型为bit的情况
					ret[columns[i]] = 1
				default:
					ret[columns[i]] = v
					break
				}
				break
			case time.Time:
				if val.IsZero() {
					ret[columns[i]] = nil
				} else {
					ret[columns[i]] = val.Format("2006-01-02 15:04:05")
				}
				break
			default:
				ret[columns[i]] = val
			}
		}
	}
	return columns, ret, nil
}

// SqlRowsToEntities lib.SqlRowsToEntities[biz.TaskEntity](UT.CommonUsecase.DB(), rows)
func SqlRowsToEntities[T any](db *gorm.DB, rows *sql.Rows) ([]*T, error) {

	var entities []*T
	for rows.Next() {
		var entity T
		err := db.ScanRows(rows, &entity)
		//err := rows.Scan(&entity)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}

// LastRune 获取最后几个字符
func LastRune(str string, num int) string {

	runs := []rune(str)
	if len(runs) < num {
		num = len(runs)
	}
	if num > 0 {
		return string(runs[len(runs)-num:])
	}
	return ""
}

// HTTPJsonWithHeaders POST GET
func HTTPJsonWithHeaders(method string, url string, params []byte, headers map[string]string) (res *string, httpCode int, err error) {

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	for k := range headers {
		header[k] = headers[k]
	}

	return Request(method, url, params, header)
}

// HTTPJson method POST GET
func HTTPJson(method string, url string, params []byte) (res *string, httpCode int, err error) {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return Request(method, url, params, headers)
}

func HTTPJsonWithBasicAuth(method string, url string, params []byte, username string, password string) (res *string, httpCode int, err error) {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Basic " + BasicAuth(username, password)
	//req.SetBasicAuth(username, password)
	return Request(method, url, params, headers)
}

func HTTPGetWithBasicAuth(url string, username string, password string) (res *string, httpCode int, err error) {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Basic " + BasicAuth(username, password)
	//req.SetBasicAuth(username, password)
	return Request("GET", url, nil, headers)
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// RequestGet Get
func RequestGet(url string, query url.Values, headers map[string]string) (*string, int, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}
	req.URL.RawQuery = query.Encode()
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)

		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		if res == "" {
			return nil, resp.StatusCode, errors.New(msg)
		} else {
			return &res, resp.StatusCode, errors.New(msg)
		}
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	result := string(bs)
	return &result, resp.StatusCode, nil
}

// Request method POST GET
func Request(method string, url string, params []byte, headers map[string]string) (res *string, httpCode int, err error) {

	body := bytes.NewBuffer(params)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	client := &http.Client{
		Timeout: time.Second * 60 * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	// StatusNoContent:204 - 服务器已经成功处理了客户端的请求,但没有返回任何内容,这种状态码通常用于不需要返回数据的情况,例如DELETE请求
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)

		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		if res == "" {
			return nil, resp.StatusCode, errors.New(msg)
		} else {
			return &res, resp.StatusCode, errors.New(msg)
		}
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	result := string(bs)
	return &result, resp.StatusCode, nil
}

// RequestStream method POST GET PUT
func RequestStream(method string, url string, body io.Reader, headers map[string]string, contentLength int64) (res *string, httpCode int, err error) {
	//body := bytes.NewBuffer(params)
	req, err := http.NewRequest(method, url, body)
	if contentLength > 0 {
		req.ContentLength = contentLength
	}
	if err != nil {
		return nil, 0, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	client := &http.Client{
		Timeout: time.Second * 60 * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	// StatusNoContent:204 - 服务器已经成功处理了客户端的请求,但没有返回任何内容,这种状态码通常用于不需要返回数据的情况,例如DELETE请求
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)

		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		if res == "" {
			return nil, resp.StatusCode, errors.New(msg)
		} else {
			return &res, resp.StatusCode, errors.New(msg)
		}
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	result := string(bs)
	return &result, resp.StatusCode, nil
}

// RequestWithQuery method POST GET
func RequestWithQuery(method string, url string, query url.Values, bodyParams []byte, headers map[string]string) (res *string, httpCode int, err error) {

	body := bytes.NewBuffer(bodyParams)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	client := &http.Client{
		Timeout: time.Second * 60 * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	// StatusNoContent:204 - 服务器已经成功处理了客户端的请求,但没有返回任何内容,这种状态码通常用于不需要返回数据的情况,例如DELETE请求
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)

		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		if res == "" {
			return nil, resp.StatusCode, errors.New(msg)
		} else {
			return &res, resp.StatusCode, errors.New(msg)
		}
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	result := string(bs)
	return &result, resp.StatusCode, nil
}

// RequestDo method POST GET
func RequestDo(method string, url string, params []byte, headers map[string]string) (*http.Response, error) {

	body := bytes.NewBuffer(params)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	client := &http.Client{
		Timeout: time.Second * 60,
	}
	return client.Do(req)
}

// RequestDoTimeout method POST GET
func RequestDoTimeout(method string, url string, params []byte, headers map[string]string, timeout time.Duration) (*http.Response, error) {

	body := bytes.NewBuffer(params)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	client := &http.Client{
		Timeout: timeout,
	}
	return client.Do(req)
}

// RuntimePath 返回：/Users/.../cmd/vbc/bin  最后没有斜杠
func RuntimePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0:i]), nil
}

func StringToTE[T any](val string, def T) (T, error) {
	r := StringToT[T](val)
	if r.IsOk() {
		return r.Unwrap(), nil
	} else {
		return def, r.Err()
	}
}

func StringToTDef[T any](val string, def T) T {
	r := StringToT[T](val)
	return r.UnwrapOr(def)
}

func StringToT[T any](val string) Result[T] {
	return TryCatch[T](func() Result[T] {

		var r T
		err := json.Unmarshal([]byte(val), &r)
		if err != nil {
			return Err[T](err)
		}
		return Ok(r)
	})
}

func BytesToTDef[T any](val []byte, def T) T {
	r := BytesToT[T](val)
	return r.UnwrapOr(def)
}

func BytesToT[T any](val []byte) Result[T] {
	return TryCatch[T](func() Result[T] {

		var r T
		err := json.Unmarshal(val, &r)
		if err != nil {
			return Err[T](err)
		}
		return Ok(r)
	})
}

func InterfaceToTE[T any](val interface{}, t T) (T, error) {
	a := InterfaceToT[T](val)
	if a.IsErr() {
		return t, a.Err()
	}
	return a.Unwrap(), nil
}

func InterfaceToTDef[T any](val interface{}, t T) T {
	a := InterfaceToT[T](val)
	return a.UnwrapOr(t)
}

func InterfaceToT[T any](val interface{}) Result[T] {
	return TryCatch[T](func() Result[T] {
		var b []byte
		var err error
		b, err = json.Marshal(val)
		if err != nil {
			return Err[T](err)
		}
		var r T
		err = json.Unmarshal(b, &r)
		if err != nil {
			return Err[T](err)
		}
		return Ok(r)
	})
}

func TryCatch[T any](try func() Result[T]) (result Result[T]) {
	defer func() {
		if err := recover(); err != nil {
			result = Err[T](errors.New(fmt.Sprintf("%v", err)))
		}
	}()
	result = try()
	return
}

type Result[T any] struct {
	value T
	err   error
}

// Ok 创建一个成功的 Result
func Ok[T any](v T) Result[T] {
	return Result[T]{value: v}
}

// Err 创建一个失败的 Result
func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}

// IsOk 检查 Result 是否成功
func (r *Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr 检查 Result 是否失败
func (r *Result[T]) IsErr() bool {
	return r.err != nil
}

// UnwrapOr 获取 Result 的值，如果失败则返回默认值
func (r *Result[T]) UnwrapOr(defaultVal T) T {
	if r.err != nil {
		return defaultVal
	}
	return r.value
}

// Unwrap 获取 Result 的值，如果失败则 panic
func (r *Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

func (r *Result[T]) Err() error {
	return r.err
}

func (r *Result[T]) Value() T {
	return r.value
}

func GetFromMapToInt64(maps map[string]interface{}, key string) int64 {
	if maps == nil {
		return 0
	}
	if _, ok := maps[key]; ok {
		r, _ := strconv.ParseInt(InterfaceToString(maps[key]), 10, 64)
		return r
	}
	return 0
}

func VerifyEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func EmptyString(str *string) bool {
	if str == nil || *str == "" {
		return true
	}
	return false
}

func HTTPPostFormData(url string, values map[string]interface{}) (*string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		w.WriteField(key, InterfaceToString(r))
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)

		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		if res == "" {
			return nil, errors.New(msg)
		} else {
			return &res, errors.New(msg)
		}
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := string(bs)
	return &result, nil

}

// NumberEnglishPrinter input: 1222 output: 1,222
func NumberEnglishPrinter(val int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", val)
}

func InterfaceToFloat32(val interface{}) float32 {
	str := InterfaceToString(val)
	a, _ := strconv.ParseFloat(str, 32)
	return float32(a)
}

func InterfaceToInt32(val interface{}) int32 {
	str := InterfaceToString(val)
	a, _ := strconv.ParseInt(str, 10, 32)
	return int32(a)
}

func InterfaceToInt64(val interface{}) int64 {
	str := InterfaceToString(val)
	a, _ := strconv.ParseInt(str, 10, 64)
	return a
}

// FormatNameWithVersion
// 输入：name v2024-06-16_11-33 输出：name_v2024-06-16_11-33
// 输入：ccc我是中文.pdf.pdf.doc  v2024-06-16_11-33 输出：ccc我是中文.pdf.pdf_v2024-06-16_11-33.doc
func FormatNameWithVersion(name string, version string) string {
	nameArr := strings.Split(name, ".")
	if len(nameArr) == 1 {
		return fmt.Sprintf("%s_%s", name, version)
	}
	name1 := strings.Join(nameArr[0:len(nameArr)-1], ".")
	return fmt.Sprintf("%s_%s.%s", name1, version, nameArr[len(nameArr)-1])
}

func StringToBytesFormat(s string) string {
	return bytes2.Format(int64(len([]byte(s))))
}

func CalDimensions(sourceW, sourceH, MaxDestW, MaxDestH float64) (float64, float64) {
	a1 := float64(sourceW) / float64(sourceH)
	a2 := float64(MaxDestW) / float64(MaxDestH)
	if a1 > a2 {
		aaa := (float64(MaxDestW) / float64(sourceW)) * float64(sourceH)
		return MaxDestW, math.Floor(aaa)
	} else {
		aaa := (float64(MaxDestH) / float64(sourceH)) * float64(sourceW)
		return math.Floor(aaa), MaxDestH
	}
}

func IsValidEmail(email string) bool {
	// 定义一个标准的 email 正则表达式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// IsValidURL 验证是否为有效的 URL
func IsValidURL(val string) bool {
	// 使用 url.Parse 解析
	u, err := url.Parse(val)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func MD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//func MD5Hash1(s string) string {
//	hash := md5.Sum([]byte(s))
//	return hex.EncodeToString(hash[:])
//}

func StringToInt32(str string) int32 {
	r, _ := strconv.ParseInt(str, 10, 32)
	return int32(r)
}

// GetExecPath 执行的目录
func GetExecPath() string {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(execPath)
}

// GetSourceRootPath 此项目源码的根目录
func GetSourceRootPath() string {
	return filepath.Dir(GetCurrentFilePath())
}

// GetCurrentFilePath 当前文件夹的路径
func GetCurrentFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}

func EqualIgnoreOrder[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	counts := make(map[T]int)
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
	}
	for _, v := range counts {
		if v != 0 {
			return false
		}
	}
	return true
}

func GetCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1) // 获取调用者的程序计数器（PC）
	fn := runtime.FuncForPC(pc)      // 根据PC获取函数信息
	return fn.Name()                 // 返回函数名
}

func GenerateSafePassword(length int) string {
	const letters = "abcdefghjkmnpqrstuvwxyz" // i l 0 O
	const upper = "ABCDEFGHJKLMNPQRSTUVWXYZ"  //  L 1
	const digits = "23456789"

	all := letters + upper + digits

	// 创建一个独立的随机源
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	password := make([]byte, length)

	// 确保包含不同类型的字符
	password[0] = letters[r.Intn(len(letters))]
	password[1] = upper[r.Intn(len(upper))]
	password[2] = digits[r.Intn(len(digits))]

	for i := 3; i < length; i++ {
		password[i] = all[r.Intn(len(all))]
	}

	// 打乱顺序
	r.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func DayWithSuffix(t time.Time) string {
	day := t.Day()
	suffix := "th"

	if day%10 == 1 && day != 11 {
		suffix = "st"
	} else if day%10 == 2 && day != 12 {
		suffix = "nd"
	} else if day%10 == 3 && day != 13 {
		suffix = "rd"
	}
	return fmt.Sprintf("%d%s", day, suffix)
}

func TruncateStringWithRune(str string, maxLen int) string {
	runes := []rune(str)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return str
}

func AESEncrypt(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func AESDecrypt(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func fixedIV() []byte {
	iv := make([]byte, aes.BlockSize)
	binary.LittleEndian.PutUint64(iv, 1) // 设置一个非零初始值
	return iv
}

// AES-CTR 加密/解密
func aesCTRTransform(input, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := fixedIV()
	stream := cipher.NewCTR(block, iv)

	output := make([]byte, len(input))
	stream.XORKeyStream(output, input)
	return output, nil
}

// EncryptToBase64 加密：明文 → Base64字符串
func EncryptToBase64(plaintext, key []byte) (string, error) {
	ciphertext, err := aesCTRTransform(plaintext, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromBase64 解密：Base64字符串 → 明文
func DecryptFromBase64(ciphertextBase64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}
	return aesCTRTransform(ciphertext, key)
}
