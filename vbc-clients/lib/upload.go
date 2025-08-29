package lib

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
	"vbc/lib/to"
)

//
//func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
//	// Prepare a form that you will submit to that URL.
//	var b bytes.Buffer
//	w := multipart.NewWriter(&b)
//	for key, r := range values {
//		var fw io.Writer
//		if x, ok := r.(io.Closer); ok {
//			defer x.Close()
//		}
//		// Add an image file
//		if x, ok := r.(*os.File); ok {
//			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
//				return
//			}
//		} else {
//			// Add other fields
//			if fw, err = w.CreateFormField(key); err != nil {
//				return
//			}
//		}
//		if _, err = io.Copy(fw, r); err != nil {
//			return err
//		}
//
//	}
//	// Don't forget to close the multipart writer.
//	// If you don't close it, your request will be missing the terminating boundary.
//	w.Close()
//
//	// Now that you have a form, you can submit it to your handler.
//	req, err := http.NewRequest("POST", url, &b)
//	if err != nil {
//		return
//	}
//	// Don't forget to set the content type, this will contain the boundary.
//	req.Header.Set("Content-Type", w.FormDataContentType())
//
//	// Submit the request
//	res, err := client.Do(req)
//	if err != nil {
//		return
//	}
//
//	// Check the response
//	if res.StatusCode != http.StatusOK {
//		err = fmt.Errorf("bad status: %s", res.Status)
//	}
//	return
//}

func PostUpload(url string, values []*UploadReader, headers map[string]string) (body *string, err error) {

	client := &http.Client{
		Timeout: time.Second * 300,
	}
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		if r == nil {
			continue
		}
		var fw io.Writer
		if x, ok := r.Reader.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		//if x, ok := r.(*os.File); ok {
		//	if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
		//		return
		//	}
		//} else

		if len(r.FileName) == 0 {
			// Add other fields
			if fw, err = w.CreateFormField(values[key].FieldName); err != nil {
				return
			}
		} else {
			if fw, err = w.CreateFormFile(values[key].FieldName, r.FileName); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r.Reader); err != nil {
			return nil, err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	for k := range headers {
		req.Header.Set(k, headers[k])
	}

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	if res != nil {
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil
		}
		body = to.Ptr(string(bs))
	}
	return
}

type UploadReader struct {
	FieldName string // 上传的字段名称
	Reader    io.Reader
	FileName  string //  有值说明的文件
}

func NewUploadReader(FieldName string, Reader io.Reader, FileName string) *UploadReader {
	return &UploadReader{
		FieldName: FieldName,
		Reader:    Reader,
		FileName:  FileName,
	}
}
