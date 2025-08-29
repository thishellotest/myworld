package biz

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

const ()

/*
	{
	  "id": "F971964745A5CD0C001BBE4E58196BFD",
	  "type": "upload_session",
	  "num_parts_processed": 455,
	  "part_size": 1024,
	  "session_endpoints": {
	    "abort": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD",
	    "commit": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD/commit",
	    "list_parts": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD/parts",
	    "log_event": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD/log",
	    "status": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD",
	    "upload_part": "https://{box-upload-server}/api/2.0/files/upload_sessions/F971964745A5CD0C001BBE4E58196BFD"
	  },
	  "session_expires_at": "2012-12-12T10:53:43-08:00",
	  "total_parts": 1000
	}
*/
type BoxUploadSession struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	NumPartsProcessed int64  `json:"num_parts_processed"`
	PartSize          int64  `json:"part_size"`
	SessionEndpoints  struct {
		Abort      string `json:"abort"`
		Commit     string `json:"commit"`
		ListParts  string `json:"list_parts"`
		LogEvent   string `json:"log_event"`
		Status     string `json:"status"`
		UploadPart string `json:"upload_part"`
	} `json:"session_endpoints"`
	TotalParts int64 `json:"total_parts"`
	TotalSize  int64
}

type BoxUploadPart struct {
	PartId string `json:"part_id"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
	SHA1   string `json:"sha1"`
}

type ZoomUploadBoxUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	BoxUsecase    *BoxUsecase
	ZoomUsecase   *ZoomUsecase
}

func NewZoomUploadBoxUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxUsecase *BoxUsecase,
	ZoomUsecase *ZoomUsecase,
) *ZoomUploadBoxUsecase {
	uc := &ZoomUploadBoxUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		BoxUsecase:    BoxUsecase,
		ZoomUsecase:   ZoomUsecase,
	}

	return uc
}

func (c *ZoomUploadBoxUsecase) TestUpload() error {
	url := "https://us06web.zoom.us/rec/download/Au386IkDQnYnDbYvq5dXsw1cvyNHH2CdgWBtB2eHatpEfvc8mkYwOri4BxP2KCMqNzVZoDmeUP7ZxIM.mumk6MD1Lup9BT4F"
	boxFileId, err := c.UploadToBoxAndCheck(url,
		"267683828516",
		"GMT20240813-203606_Recording_1920x1080.mp4",
		376674194,
	)
	c.log.Info(err, " : ", boxFileId)
	return err
}

func (c *ZoomUploadBoxUsecase) UploadToBoxAndCheck(downloadUrl string, boxFolderId string, boxFileName string, fileTotalSize int64) (boxFileId string, err error) {

	err = c.UploadToBox(downloadUrl, boxFolderId, boxFileName, fileTotalSize)
	if err != nil {
		return "", err
	}
	items, err := c.BoxUsecase.ListItemsInFolderFormat(boxFolderId)
	if err != nil {
		return "", err
	}
	for _, v := range items {
		if v.GetString("type") == "file" && v.GetString("name") == boxFileName {
			return v.GetString("id"), nil
		}
	}
	return "", errors.New("The corresponding file was not found")
}

func (c *ZoomUploadBoxUsecase) UploadToBox(downloadUrl string, boxFolderId string, boxFileName string, fileTotalSize int64) error {
	// Step 1: 获取远程文件大小
	//headResp, err := http.Head(downloadUrl)
	//if err != nil {
	//	return err
	//}
	//totalSize := headResp.ContentLength
	headers, err := c.ZoomUsecase.Headers()
	if err != nil {
		return err
	}

	// Step 2: 创建上传会话
	session, err := c.BoxUsecase.CreateUploadSession(boxFolderId, boxFileName, fileTotalSize)
	if err != nil {
		return err
	}
	session.TotalSize = fileTotalSize
	c.log.Info(fmt.Sprintf("🚀 创建上传会话成功，session：%s", InterfaceToString(session)))
	c.log.Info(fmt.Sprintf("🚀 创建上传会话成功，分片大小：%d", session.PartSize))

	// Step 3: 流式下载 + 分片上传
	//resp, err := http.Get(downloadUrl)

	resp, err := lib.RequestDoTimeout("GET", downloadUrl, nil, headers, time.Hour)
	//httpResponse, err := http.Get(task.DownloadUrl)
	if err != nil {
		c.log.Error(err)
		return err
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.log.Info(fmt.Sprintf("下载失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes)))
		return fmt.Errorf("下载失败: %s", resp.Status)
	}

	hash := sha1.New()
	teeReader := io.TeeReader(resp.Body, hash)

	var parts []BoxUploadPart
	//var actualSize int64
	buf := make([]byte, session.PartSize)
	var offset int64 = 0
	for {
		n, err := io.ReadFull(teeReader, buf)
		//n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		if n == 0 {
			break
		}
		part, err := c.BoxUsecase.UploadPart(session, buf[:n], offset)
		if err != nil {
			return err
		}
		parts = append(parts, part)
		offset += int64(n)

		// 累计上传的实际字节数
		//actualSize = offset

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	if offset != fileTotalSize {
		errMsg := fmt.Sprintf("上传数据大小与预期不符，实际上传 %d 字节，预期 %d 字节", offset, fileTotalSize)
		c.log.Info(errMsg)
		return fmt.Errorf(errMsg)
	}
	digest := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	c.log.Infof("📦 SHA1 Digest: %s", digest)
	// Step 4: 提交上传
	return c.BoxUsecase.CommitUploadForDownload(session.SessionEndpoints.Commit, parts, digest)
}

/*
UploadFileToBoxUseChunks

filePath      = "/tmp/VSCode-darwin-arm64.zip"
fileName      = "VSCode-darwin-arm64.zip"
boxFolderId      = "267683828516"
*/
func (c *BoxUsecase) UploadFileToBoxUseChunks(filePath string, fileName string, boxFolderId string) error {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	session, err := c.CreateUploadSession(boxFolderId, fileName, fileSize)
	if err != nil {
		return err
	}
	session.TotalSize = fileSize
	var parts []BoxUploadPart
	buf := make([]byte, session.PartSize)
	var offset int64 = 0
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		part, err := c.UploadPart(session, buf[:n], offset)
		if err != nil {
			return err
		}
		parts = append(parts, part)
		offset += int64(n)
		if err == io.EOF {
			break
		}
	}
	err = c.CommitUpload(session.SessionEndpoints.Commit, parts)
	c.log.Info("✅ 文件上传完成！", err)
	return nil
}

func (c *BoxUsecase) CreateUploadSession(folderID string, fileName string, fileSize int64) (uploadSession BoxUploadSession, err error) {
	//payload := map[string]interface{}{
	//	"folder_id": folderID,
	//	"file_name": fileName,
	//	"file_size": fileSize,
	//}

	token, err := c.Token()
	if err != nil {
		return uploadSession, err
	}

	api := fmt.Sprintf("%s/api/2.0/files/upload_sessions", c.conf.Box.UploadUrl)
	c.log.Info(api)
	params := make(lib.TypeMap)
	params.Set("folder_id", folderID)
	params.Set("file_name", fileName)
	params.Set("file_size", fileSize)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CreateUploadSession"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	//lib.DPrintln(res)
	c.log.Info(InterfaceToString(res))
	if err != nil {
		return uploadSession, err
	}
	if res != nil {
		//resMap := lib.ToTypeMapByString(*res)
		//lib.DPrintln(resMap)
		err = json.Unmarshal([]byte(*res), &uploadSession)
		return uploadSession, err
	}
	return uploadSession, err

	//body, _ := json.Marshal(payload)
	//req, _ := http.NewRequest("POST", uploadInitURL, bytes.NewReader(body))
	//req.Header.Set("Authorization", "Bearer "+accessToken)
	//req.Header.Set("Content-Type", "application/json")
	//
	//resp, _ := http.DefaultClient.Do(req)
	//defer resp.Body.Close()
	//
	//var session UploadSession
	//json.NewDecoder(resp.Body).Decode(&session)
	//return session
}

func (c *BoxUsecase) UploadPart(session BoxUploadSession, data []byte, offset int64) (boxUploadPart BoxUploadPart, err error) {

	token, err := c.Token()
	if err != nil {
		return boxUploadPart, err
	}

	h := sha1.New()
	h.Write(data)
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req, _ := http.NewRequest("PUT", session.SessionEndpoints.UploadPart, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Digest", "SHA="+digest)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", offset, offset+int64(len(data))-1, session.TotalSize))

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	//lib.DPrintln(resp.Body)
	//lib.DPrintln("err:", err)
	if err != nil {
		return BoxUploadPart{}, err
	}

	var partInfo struct {
		Part BoxUploadPart `json:"part"`
	}
	err = json.NewDecoder(resp.Body).Decode(&partInfo)
	if err != nil {
		return boxUploadPart, err
	}
	c.log.Info(fmt.Sprintf("✅ Chunk uploaded: offset=%d, size=%d\n", offset, len(data)))
	return partInfo.Part, nil
}

func (c *BoxUsecase) CommitUploadForDownload(commitUrl string, parts []BoxUploadPart, digest string) error {

	token, err := c.Token()
	if err != nil {
		return err
	}
	//https://upload.app.box.com/api/2.0/files/upload_sessions/5B4A8BD7E088EEEBB3231E2A3292C2B2/commit
	body := map[string]interface{}{
		"parts": parts,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.log.Info(InterfaceToString(body))
	req, err := http.NewRequest("POST", commitUrl, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Digest", "SHA="+digest)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	c.log.Info(resp)
	if err != nil {
		c.log.Info(err)
		return err
	}
	c.log.Info("🚀 文件上传提交完成")
	return nil

}

func (c *BoxUsecase) CommitUpload(commitUrl string, parts []BoxUploadPart) error {

	filePath := "/tmp/VSCode-darwin-arm64.zip"
	//fileName := "VSCode-darwin-arm64.zip"
	token, err := c.Token()
	if err != nil {
		return err
	}
	//https://upload.app.box.com/api/2.0/files/upload_sessions/5B4A8BD7E088EEEBB3231E2A3292C2B2/commit
	body := map[string]interface{}{
		"parts": parts,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.log.Info(InterfaceToString(body))
	req, err := http.NewRequest("POST", commitUrl, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Digest", "SHA="+computeFileSHA1(filePath))

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	c.log.Info(resp)
	if err != nil {
		c.log.Info(err)
		return err
	}
	c.log.Info("🚀 文件上传提交完成")
	return nil

}

// 可选：整个文件 hash（用于 commit 校验）
func computeFileSHA1(path string) string {
	f, _ := os.Open(path)
	defer f.Close()

	h := sha1.New()
	io.Copy(h, f)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
