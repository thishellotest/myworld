package tests

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"os"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)
import drive "google.golang.org/api/drive/v3"

func Test_GoogleDrive_Files(t *testing.T) {

	a, err := UT.Oauth2TokenUsecase.GetByAppId(biz.Oauth2_AppId_google)
	lib.DPrintln(err)
	fmt.Println("accessToken:", a.AccessToken)
	tokenSource, err := a.TokenSourceMod()
	aa, err := drive.NewService(context.TODO(), option.WithTokenSource(tokenSource))
	//lib.DPrintln(err)
	ccc := aa.Files.List()
	//r, err := ccc.Q("mimeType='application/vnd.google-apps.folder'").Do()
	r, err := ccc.Q("'1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR' in parents").Do()
	lib.DPrintln(r, err)

	//r2, err := aa.Files.Get("1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR").Do()
	//lib.DPrintln(r2, err)
	//
	//r3, err := aa.Files.ListLabels("1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR").Do()
	//lib.DPrintln(r3, err)

	//drives := aa.Drives.List()
	//r1, err := drives.UseDomainAdminAccess(false).Do()
	//for _, v := range r1.Drives {
	//	lib.DPrintln(v.Id)
	//}
	//lib.DPrintln(r1, len(r1.Drives), err)

}

// https://stackoverflow.com/questions/73981337/how-to-upload-files-to-googledrive-and-share-with-anyone-using-serviceaccount-a
func Test_GoogleDrive_upload(t *testing.T) {

	/*


		filename := "./lemon.txt"
		    file, err := os.Open(filename)
		    if err != nil {
		        log.Fatalln(err)
		    }


			res, err := srv.Files.Create(
			    &drive.File{
			        Parents: []string{"17n-EpJcGg0DmmWqSoJ75iIUdXDP7neoH"},
			        Name:    "banana.txt",
			        Permissions: []*drive.Permission{
			            {
			                Role: "reader",
			                Type: "anyone",
			            },
			        },
			    },
			).Media(file, googleapi.ChunkSize(int(stat.Size()))).Do()
			if err != nil {
			    log.Fatalln(err)
			}

			fmt.Printf("%s\n", res.Id)


				res, err := srv.Files.Create(
				    &drive.File{
				        Parents: []string{"17n-EpJcGg0DmmWqSoJ75iIUdXDP7neoH"},
				        Name:    "banana.txt",
				    },
				).Media(file, googleapi.ChunkSize(int(stat.Size()))).Do()
				if err != nil {
				    log.Fatalln(err)
				}
				fmt.Printf("%s\n", res.Id)

				res2, err := srv.Permissions.Create(res.Id, &drive.Permission{
				    Role: "reader",
				    Type: "anyone",
				}).Do()


	*/

	filename := "res/b.pdf"
	fileData, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	a, err := UT.Oauth2TokenUsecase.GetByAppId(biz.Oauth2_AppId_google)
	lib.DPrintln(err)
	fmt.Println("accessToken:", a.AccessToken)
	tokenSource, err := a.TokenSourceMod()
	aa, err := drive.NewService(context.TODO(), option.WithTokenSource(tokenSource))
	lib.DPrintln(err)
	file := drive.File{
		MimeType: "application/pdf",
		Name:     "test.pdf",
		Parents:  []string{"1KZo9l1sOENV4s8DJbL2vIx8FBiBIznww"},
	}

	f, err := aa.Files.Create(&file).Media(fileData).Do()
	lib.DPrintln(f, err)
}

func Test_GoogleDrive_CreateFolder(t *testing.T) {

	a, err := UT.Oauth2TokenUsecase.GetByAppId(biz.Oauth2_AppId_google)
	lib.DPrintln(err)
	fmt.Println("accessToken:", a.AccessToken)
	tokenSource, err := a.TokenSourceMod()
	aa, err := drive.NewService(context.TODO(), option.WithTokenSource(tokenSource))
	lib.DPrintln(err)
	file := drive.File{
		MimeType: "application/vnd.google-apps.folder",
		Name:     "newFolder",
		Parents:  []string{"1KZo9l1sOENV4s8DJbL2vIx8FBiBIznww"},
	}

	f, err := aa.Files.Create(&file).Do()
	lib.DPrintln(f, err)
}
