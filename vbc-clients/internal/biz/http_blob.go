package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type HttpBlobUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	JWTUsecase         *JWTUsecase
	BlobCommentUsecase *BlobCommentUsecase
	UserUsecase        *UserUsecase
	LogUsecase         *LogUsecase
	BoxUsecase         *BoxUsecase
	BlobbuzUsecase     *BlobbuzUsecase
	BlobUsecase        *BlobUsecase
	ClientCaseUsecase  *ClientCaseUsecase
	TUsecase           *TUsecase
	HttpTUsecase       *HttpTUsecase
	AzstorageUsecase   *AzstorageUsecase
	BlobSliceUsecase   *BlobSliceUsecase
	BoxbuzUsecase      *BoxbuzUsecase
	FileUsecase        *FileUsecase
}

func NewHttpBlobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	BlobCommentUsecase *BlobCommentUsecase,
	UserUsecase *UserUsecase,
	LogUsecase *LogUsecase,
	BoxUsecase *BoxUsecase,
	BlobbuzUsecase *BlobbuzUsecase,
	BlobUsecase *BlobUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	TUsecase *TUsecase,
	HttpTUsecase *HttpTUsecase,
	AzstorageUsecase *AzstorageUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	FileUsecase *FileUsecase) *HttpBlobUsecase {
	uc := &HttpBlobUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		JWTUsecase:         JWTUsecase,
		BlobCommentUsecase: BlobCommentUsecase,
		UserUsecase:        UserUsecase,
		LogUsecase:         LogUsecase,
		BoxUsecase:         BoxUsecase,
		BlobbuzUsecase:     BlobbuzUsecase,
		BlobUsecase:        BlobUsecase,
		ClientCaseUsecase:  ClientCaseUsecase,
		TUsecase:           TUsecase,
		HttpTUsecase:       HttpTUsecase,
		AzstorageUsecase:   AzstorageUsecase,
		BlobSliceUsecase:   BlobSliceUsecase,
		BoxbuzUsecase:      BoxbuzUsecase,
		FileUsecase:        FileUsecase,
	}

	return uc
}

func (c *HttpBlobUsecase) Delete(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	user, _ := c.JWTUsecase.JWTUser(ctx)
	data, err := c.BizDelete(user, body.GetString("gid"), body.GetString("blob_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizDelete(tUser TData, Gid string, BlobGid string) (lib.TypeMap, error) {

	//data := make(lib.TypeMap)
	if Gid == "" {
		return nil, errors.New("Gid cannot be empty")
	}
	if BlobGid == "" {
		return nil, errors.New("BlobGid cannot be empty")
	}

	BlobComment, err := c.BlobCommentUsecase.GetByCond(Eq{"deleted_at": 0, "blob_gid": BlobGid, "gid": Gid})
	if err != nil {
		return nil, err
	}
	if BlobComment == nil {
		return nil, errors.New("BlobComment is nil")
	}
	BlobComment.DeletedAt = time.Now().Unix()

	err = c.CommonUsecase.DB().Save(&BlobComment).Error
	if err != nil {
		return nil, err
	}
	c.LogUsecase.SaveLog(BlobComment.ID, "BlobCommentDelete", map[string]interface{}{
		"UserId": tUser.CustomFields.NumberValueByNameBasic("id"),
	})
	return nil, nil
}

func (c *HttpBlobUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	user, _ := c.JWTUsecase.JWTUser(ctx)
	data, err := c.BizSave(user, body.GetString("gid"),
		body.GetString("blob_gid"),
		body.GetString("content"),
		body.GetString("json_data"),
		body.GetInt("page"),
		body.GetInt("type"),
	)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizSave(tUser TData, Gid string, BlobGid string,
	Content string, JsonData string, Page int32, Type int32) (lib.TypeMap, error) {

	//data := make(lib.TypeMap)
	if Content == "" {
		return nil, errors.New("Content cannot be empty")
	}
	if JsonData == "" {
		return nil, errors.New("JsonData cannot be empty")
	}
	if Page <= 0 {
		return nil, errors.New("The page number is wrong")
	}
	if Type != BlobComment_Type_Default && Type != BlobComment_Type_OnlyText {
		return nil, errors.New("The type is wrong")
	}
	var BlobComment *BlobCommentEntity
	var err error
	if Gid != "" {
		BlobComment, err = c.BlobCommentUsecase.GetByCond(Eq{"deleted_at": 0, "blob_gid": BlobGid, "gid": Gid})
		if err != nil {
			return nil, err
		}
		if BlobComment == nil {
			return nil, errors.New("BlobComment is nil")
		}
	} else {
		BlobComment = &BlobCommentEntity{
			Gid:       uuid.UuidWithoutStrike(),
			BlobGid:   BlobGid,
			CreatedAt: time.Now().Unix(),
			Type:      int(Type),
		}
	}
	BlobComment.UpdatedAt = time.Now().Unix()
	BlobComment.Content = Content
	BlobComment.JsonData = JsonData
	BlobComment.Page = Page
	BlobComment.UserGid = tUser.CustomFields.TextValueByNameBasic("gid")
	err = c.CommonUsecase.DB().Save(&BlobComment).Error
	if err != nil {
		return nil, err
	}
	c.LogUsecase.SaveLog(BlobComment.ID, "BlobCommentSave", map[string]interface{}{
		"UserGid": tUser.CustomFields.TextValueByNameBasic("gid"),
	})

	userCaches := lib.CacheInit[*TData]()
	comment := BlobComment.BlobCommentToApi(userCaches, c.UserUsecase, c.log)

	r := make(lib.TypeMap)
	r.Set("comment", comment)
	return r, nil
}

func (c *HttpBlobUsecase) Comments(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizComments(body.GetString("blob_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizComments(BlobGid string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	records, err := c.BlobCommentUsecase.AllByCond(Eq{"deleted_at": 0, "blob_gid": BlobGid})
	if err != nil {
		return nil, err
	}
	var Comments lib.TypeList
	userCaches := lib.CacheInit[*TData]()
	for _, v := range records {
		Comments = append(Comments, v.BlobCommentToApi(userCaches, c.UserUsecase, c.log))
	}
	data.Set("comments", Comments)
	return data, nil
}

func (c *HttpBlobUsecase) CreateTask(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	user, _ := c.JWTUsecase.JWTUser(ctx)
	data, err := c.BizCreateTask(ctx, user, body.GetString("box_file_id"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizCreateTask(ctx context.Context, tUser TData, BoxFileId string) (lib.TypeMap, error) {

	if BoxFileId == "" {
		return nil, errors.New("BoxFileId is empty")
	}

	fileInfo, _, err := c.BoxUsecase.GetFileInfoForTypeMap(BoxFileId)
	lib.DPrintln(fileInfo)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	tCase, err := c.BlobbuzUsecase.GetCaseByBoxFileInfo(fileInfo)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	fileVersionID := fileInfo.GetString("file_version.id")
	fileName := fileInfo.GetString("name")

	_, suffix := lib.FileExt(fileName, true)
	if suffix != BlobType_pdf {
		return nil, errors.New("The file format is not supported")
	}

	uniqblob := GenUniqblob(BoxFileId, fileVersionID)
	tBlob, err := c.TUsecase.Data(Kind_blobs, Eq{"uniqblob": uniqblob, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tBlob == nil {
		tBlob, err = c.BlobbuzUsecase.HandleBoxFile(ctx, fileInfo,
			tCase.Gid(), tUser.Gid())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("The task already exists")
	}

	data := make(lib.TypeMap)
	userCaches := lib.CacheInit[*TData]()
	caseCaches := lib.CacheInit[*TData]()
	data.Set(Fab_Blob, BlobToApi(tBlob, userCaches, c.UserUsecase, caseCaches, c.ClientCaseUsecase, c.AzstorageUsecase, c.log))
	return data, nil
}

func (c *HttpBlobUsecase) OcrDetail(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizOcrDetail(ctx, userFacade, body.GetString("gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)

}

func (c *HttpBlobUsecase) BizOcrDetail(ctx context.Context, userFacade UserFacade, gid string) (lib.TypeMap, error) {

	entity, _ := c.BlobSliceUsecase.GetByCond(Eq{"gid": gid})
	if entity == nil {
		return nil, errors.New("Parameter error")
	}
	data := make(lib.TypeMap)
	data.Set("ocr_result", entity.OcrResult)
	return data, nil
}

func (c *HttpBlobUsecase) RecordReviewFiles(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizRecordReviewFiles(ctx, userFacade, body.GetString("gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type RecordReviewFileList []RecordReviewFileVo

type RecordReviewFileVo struct {
	SourceId    string `json:"source_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	CanReview   bool   `json:"can_review"`
}

func (c *HttpBlobUsecase) BizRecordReviewFiles(ctx context.Context, userFacade UserFacade, caseGid string) (lib.TypeMap, error) {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("Parameter error")
	}
	recordReviewFolderId, err := c.BoxbuzUsecase.DCRecordReviewFolderId(tCase)
	if err != nil {
		return nil, err
	}
	if recordReviewFolderId == "" {
		return nil, errors.New("Cannot find Record Review Folder ID")
	}

	files, _ := c.FileUsecase.AllByCond(Eq{"parent_folder_id": recordReviewFolderId, "biz_deleted_at": 0})

	var folderIds []string
	for _, v := range files {
		if v.Type == "folder" {
			folderIds = append(folderIds, v.SourceId)
		}
	}

	subFiles, err := c.FileUsecase.AllByCond(And(Eq{"biz_deleted_at": 0, "type": "file"}, In("parent_folder_id", folderIds)))
	if err != nil {
		return nil, err
	}

	var recordReviewFileList RecordReviewFileList
	for _, v := range files {
		if v.Type == "file" {
			recordReviewFileList = append(recordReviewFileList, RecordReviewFileVo{
				SourceId:    v.SourceId,
				Name:        v.SourceName,
				DisplayName: v.SourceName,
				CanReview:   v.CanReview(),
			})
		}
	}

	getParentFolder := func(sourceId string) *FileEntity {
		for k, v := range files {
			if v.SourceId == sourceId {
				return files[k]
			}
		}
		return nil
	}

	for _, v := range subFiles {
		temp := getParentFolder(v.ParentFolderId)

		displayName := v.SourceName
		if temp != nil {
			displayName = temp.SourceName + "/" + displayName
		}
		recordReviewFileList = append(recordReviewFileList, RecordReviewFileVo{
			SourceId:    v.SourceId,
			Name:        v.SourceName,
			DisplayName: displayName,
			CanReview:   v.CanReview(),
		})
	}

	data := make(lib.TypeMap)
	data.Set("files", recordReviewFileList)
	return data, nil
}

//
//func (c *HttpBlobUsecase) RecordReviewTasksProgress(ctx *gin.Context) {
//	reply := CreateReply()
//	rawData, _ := ctx.GetRawData()
//	body := lib.ToTypeMapByString(string(rawData))
//
//	blobSliceGids, _ := lib.ConvertTypeListInterface[string](body.GetTypeListInterface("blob_slice_gids"))
//
//	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
//	data, err := c.BlobbuzUsecase.BizRecordReviewTasksProgress(ctx, userFacade, blobSliceGids)
//	if err != nil {
//		reply.CommonError(err)
//	} else {
//		reply.Merge(data)
//		reply.Success()
//	}
//	ctx.JSON(200, reply)
//}

func (c *HttpBlobUsecase) SliceJoinOcr(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BlobbuzUsecase.BizSliceJoinOcr(ctx, userFacade, body.GetString("blob_slice_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) RecordReviewDetail(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	blobSliceGids, _ := lib.ConvertTypeListInterface[string](body.GetTypeListInterface("blob_slice_gids"))
	data, err := c.BizRecordReviewDetail(ctx, userFacade, body.GetString("gid"), blobSliceGids)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizRecordReviewDetail(ctx context.Context, userFacade UserFacade, gid string, blobSliceGids []string) (lib.TypeMap, error) {

	if gid == "" {
		return nil, errors.New("gid is empty")
	}

	tBlob, err := c.TUsecase.Data(Kind_blobs, Eq{"gid": gid})
	if err != nil {
		return nil, err
	}
	if tBlob == nil {
		return nil, errors.New("tBlob is nil")
	}

	var gidCond Cond
	if len(blobSliceGids) > 0 {
		gidCond = In("gid", blobSliceGids)
	}

	slices, err := c.BlobSliceUsecase.AllByCondWithOrderBySelect(BlobSlice_ToApi_Columns,
		And(Eq{"blob_gid": gid, "deleted_at": 0}, gidCond),
		"slice_id", 5000)
	if err != nil {
		return nil, err
	}
	var slicesRes lib.TypeList
	for _, v := range slices {
		slicesRes = append(slicesRes, v.ToApi(c.CommonUsecase, c.log, c.AzstorageUsecase))
	}
	data := make(lib.TypeMap)
	data.Set("slices", slicesRes)
	return data, nil
}

func (c *HttpBlobUsecase) Detail(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	user, _ := c.JWTUsecase.JWTUser(ctx)
	data, err := c.BizDetail(ctx, user, body.GetString("gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizDetail(ctx context.Context, tUser TData, gid string) (lib.TypeMap, error) {

	if gid == "" {
		return nil, errors.New("gid is empty")
	}

	tBlob, err := c.TUsecase.Data(Kind_blobs, Eq{"gid": gid})
	if err != nil {
		return nil, err
	}
	if tBlob == nil {
		return nil, errors.New("tBlob is nil")
	}
	data := make(lib.TypeMap)
	userCaches := lib.CacheInit[*TData]()
	caseCaches := lib.CacheInit[*TData]()
	data.Set(Fab_Blob, BlobToApi(tBlob, userCaches, c.UserUsecase, caseCaches, c.ClientCaseUsecase, c.AzstorageUsecase, c.log))
	return data, nil
}

func (c *HttpBlobUsecase) TaskList(ctx *gin.Context) {

	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//user, _ := c.JWTUsecase.JWTUser(ctx)
	data, err := c.BizTaskList()
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpBlobUsecase) BizTaskList() (lib.TypeMap, error) {
	tList, total, page, pageSize, err := c.HttpTUsecase.DoBizList(Kind_blobs, nil)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)

	var destList lib.TypeList
	userCaches := lib.CacheInit[*TData]()
	caseCaches := lib.CacheInit[*TData]()
	for k, _ := range tList {
		//  v.CustomFields.ToApiMap()
		destList = append(destList, BlobToApi(&tList[k], userCaches, c.UserUsecase, caseCaches, c.ClientCaseUsecase, c.AzstorageUsecase, c.log))
	}

	data.Set(Fab_TData+"."+Fab_TList, destList)
	data.Set(Fab_TData+"."+Fab_TTotal, int32(total))
	data.Set(Fab_TData+"."+Fab_TPage, page)
	data.Set(Fab_TData+"."+Fab_TPageSize, pageSize)
	return data, nil
}
