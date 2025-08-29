package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_RecordReviewUsecase_Process(t *testing.T) {
	str := `{
	"type": "webhook_event",
	"id": "8614248c-5ae1-473a-bfe3-cd6ed043f301",
	"created_at": "2024-06-06T01:59:18-07:00",
	"trigger": "FILE.UPLOADED",
	"webhook": {
		"id": "2802339849",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "1552371904380",
		"type": "file",
		"file_version": {
			"type": "file_version",
			"id": "1705580056531",
			"sha1": "c4accfa363e9192347afd8477e04baafc8adb868"
		},
		"sequence_id": "3",
		"etag": "3",
		"sha1": "c4accfa363e9192347afd8477e04baafc8adb868",
		"name": "33.boxnote",
		"description": "",
		"size": 591,
		"path_collection": {
			"total_count": 7,
			"entries": [{
				"type": "folder",
				"id": "0",
				"sequence_id": null,
				"etag": null,
				"name": "All Files"
			}, {
				"type": "folder",
				"id": "241183180615",
				"sequence_id": "4",
				"etag": "4",
				"name": "VBC Engineering Team"
			}, {
				"type": "folder",
				"id": "247457173873",
				"sequence_id": "1",
				"etag": "1",
				"name": "Testing"
			}, {
				"type": "folder",
				"id": "255166311971",
				"sequence_id": "1",
				"etag": "1",
				"name": "Test Clients"
			}, {
				"type": "folder",
				"id": "264686374897",
				"sequence_id": "2",
				"etag": "2",
				"name": "[PROD]VBC - TestLiao, TestGary #5076"
			}, {
				"type": "folder",
				"id": "264686394097",
				"sequence_id": "0",
				"etag": "0",
				"name": "VA Medical Records"
			}, {
				"type": "folder",
				"id": "268258622342",
				"sequence_id": "0",
				"etag": "0",
				"name": "leve1_folder"
			}]
		},
		"created_at": "2024-06-06T01:19:50-07:00",
		"modified_at": "2024-06-06T01:59:18-07:00",
		"trashed_at": null,
		"purged_at": null,
		"content_created_at": "2024-06-06T01:19:50-07:00",
		"content_modified_at": "2024-06-06T01:58:58-07:00",
		"created_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"modified_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"owned_by": {
			"type": "user",
			"id": "30690179025",
			"name": "Yannan Wang",
			"login": "ywang@vetbenefitscenter.com"
		},
		"shared_link": null,
		"parent": {
			"type": "folder",
			"id": "268258622342",
			"sequence_id": "0",
			"etag": "0",
			"name": "leve1_folder"
		},
		"item_status": "active"
	},
	"additional_info": []
}`

	str = `{
	"type": "webhook_event",
	"id": "398b0e60-818d-4c55-9a9c-a98cc7c7ae88",
	"created_at": "2024-06-09T05:03:02-07:00",
	"trigger": "FOLDER.CREATED",
	"webhook": {
		"id": "2802339849",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "269014014377",
		"type": "folder",
		"sequence_id": "0",
		"etag": "0",
		"name": "666",
		"created_at": "2024-06-09T05:03:02-07:00",
		"modified_at": "2024-06-09T05:03:02-07:00",
		"description": "",
		"size": 0,
		"path_collection": {
			"total_count": 7,
			"entries": [{
				"type": "folder",
				"id": "0",
				"sequence_id": null,
				"etag": null,
				"name": "All Files"
			}, {
				"type": "folder",
				"id": "241183180615",
				"sequence_id": "4",
				"etag": "4",
				"name": "VBC Engineering Team"
			}, {
				"type": "folder",
				"id": "247457173873",
				"sequence_id": "1",
				"etag": "1",
				"name": "Testing"
			}, {
				"type": "folder",
				"id": "255166311971",
				"sequence_id": "1",
				"etag": "1",
				"name": "Test Clients"
			}, {
				"type": "folder",
				"id": "264686374897",
				"sequence_id": "2",
				"etag": "2",
				"name": "[PROD]VBC - TestLiao, TestGary #5076"
			}, {
				"type": "folder",
				"id": "264686394097",
				"sequence_id": "0",
				"etag": "0",
				"name": "VA Medical Records"
			}, {
				"type": "folder",
				"id": "268258622342",
				"sequence_id": "0",
				"etag": "0",
				"name": "leve1_folder"
			}]
		},
		"created_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"modified_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"trashed_at": null,
		"purged_at": null,
		"content_created_at": "2024-06-09T05:03:02-07:00",
		"content_modified_at": "2024-06-09T05:03:02-07:00",
		"owned_by": {
			"type": "user",
			"id": "30690179025",
			"name": "Yannan Wang",
			"login": "ywang@vetbenefitscenter.com"
		},
		"shared_link": null,
		"folder_upload_email": null,
		"parent": {
			"type": "folder",
			"id": "268258622342",
			"sequence_id": "0",
			"etag": "0",
			"name": "leve1_folder"
		},
		"item_status": "active"
	},
	"additional_info": []
}`

	// 5004
	str = `{
	"type": "webhook_event",
	"id": "8cefa745-cb04-4840-a6da-eaa44bdfcd0c",
	"created_at": "2024-06-09T05:13:48-07:00",
	"trigger": "FILE.UPLOADED",
	"webhook": {
		"id": "2802339849",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "1555186130116",
		"type": "file",
		"file_version": {
			"type": "file_version",
			"id": "1708665244516",
			"sha1": "53ca6d123553cdf83fbbf9c793cad5db33392e08"
		},
		"sequence_id": "0",
		"etag": "0",
		"sha1": "53ca6d123553cdf83fbbf9c793cad5db33392e08",
		"name": "a.pdf",
		"description": "",
		"size": 162709,
		"path_collection": {
			"total_count": 8,
			"entries": [{
				"type": "folder",
				"id": "0",
				"sequence_id": null,
				"etag": null,
				"name": "All Files"
			}, {
				"type": "folder",
				"id": "241183180615",
				"sequence_id": "4",
				"etag": "4",
				"name": "VBC Engineering Team"
			}, {
				"type": "folder",
				"id": "247457173873",
				"sequence_id": "1",
				"etag": "1",
				"name": "Testing"
			}, {
				"type": "folder",
				"id": "255166311971",
				"sequence_id": "1",
				"etag": "1",
				"name": "Test Clients"
			}, {
				"type": "folder",
				"id": "264665862512",
				"sequence_id": "1",
				"etag": "1",
				"name": "VBC - Local #5004"
			}, {
				"type": "folder",
				"id": "264665881712",
				"sequence_id": "0",
				"etag": "0",
				"name": "VA Medical Records"
			}, {
				"type": "folder",
				"id": "269015152061",
				"sequence_id": "0",
				"etag": "0",
				"name": "aaaaa"
			}, {
				"type": "folder",
				"id": "269014396652",
				"sequence_id": "0",
				"etag": "0",
				"name": "ccccc"
			}]
		},
		"created_at": "2024-06-09T05:13:47-07:00",
		"modified_at": "2024-06-09T05:13:47-07:00",
		"trashed_at": null,
		"purged_at": null,
		"content_created_at": "2024-05-18T19:08:28-07:00",
		"content_modified_at": "2024-05-18T19:08:28-07:00",
		"created_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"modified_by": {
			"type": "user",
			"id": "30888625898",
			"name": "VBC Team",
			"login": "info@vetbenefitscenter.com"
		},
		"owned_by": {
			"type": "user",
			"id": "30690179025",
			"name": "Yannan Wang",
			"login": "ywang@vetbenefitscenter.com"
		},
		"shared_link": null,
		"parent": {
			"type": "folder",
			"id": "269014396652",
			"sequence_id": "0",
			"etag": "0",
			"name": "ccccc"
		},
		"item_status": "active"
	},
	"additional_info": []
}`

	// 文件
	str = `{"type":"webhook_event","id":"7a91c72d-4be5-4791-b5ea-147996759421","created_at":"2024-06-09T19:21:44-07:00","trigger":"FILE.UPLOADED","webhook":{"id":"2823532065","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"source":{"id":"1555570084468","type":"file","file_version":{"type":"file_version","id":"1709091942868","sha1":"c5995ef02e12ac636432db8831674f7756e8492c"},"sequence_id":"0","etag":"0","sha1":"c5995ef02e12ac636432db8831674f7756e8492c","name":"a.pdf","description":"","size":85724,"path_collection":{"total_count":5,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183605424","sequence_id":"4","etag":"4","name":"VBC All Teams"},{"type":"folder","id":"241109085470","sequence_id":"1","etag":"1","name":"Clients"},{"type":"folder","id":"264686374897","sequence_id":"4","etag":"4","name":"VBC - TestLN, TestFN #5076"},{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"}]},"created_at":"2024-06-09T19:21:44-07:00","modified_at":"2024-06-09T19:21:44-07:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-06-02T22:17:19-07:00","content_modified_at":"2024-06-02T22:17:19-07:00","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"},"item_status":"active"},"additional_info":[]}`

	// 文件夹
	str = `{"type":"webhook_event","id":"f07a661f-103e-4d83-a28f-c6fa72ddc39e","created_at":"2024-06-09T19:12:12-07:00","trigger":"FILE.UPLOADED","webhook":{"id":"2823532065","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"source":{"id":"1554566204507","type":"file","file_version":{"type":"file_version","id":"1709084399039","sha1":"c5995ef02e12ac636432db8831674f7756e8492c"},"sequence_id":"2","etag":"2","sha1":"c5995ef02e12ac636432db8831674f7756e8492c","name":"a.pdf","description":"","size":85724,"path_collection":{"total_count":7,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183605424","sequence_id":"4","etag":"4","name":"VBC All Teams"},{"type":"folder","id":"241109085470","sequence_id":"1","etag":"1","name":"Clients"},{"type":"folder","id":"264686374897","sequence_id":"4","etag":"4","name":"VBC - TestLN, TestFN #5076"},{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"},{"type":"folder","id":"268906213262","sequence_id":"2","etag":"2","name":"FolderTest"},{"type":"folder","id":"268905519956","sequence_id":"0","etag":"0","name":"FolderTestSub"}]},"created_at":"2024-06-08T03:41:30-07:00","modified_at":"2024-06-09T19:12:12-07:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-05-18T19:08:28-07:00","content_modified_at":"2024-06-02T22:17:19-07:00","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"268905519956","sequence_id":"0","etag":"0","name":"FolderTestSub"},"item_status":"active"},"additional_info":[]}`

	typeMap := lib.ToTypeMapByString(str)
	a, err := UT.RecordReviewUsecase.Process(context.TODO(), typeMap, 0)
	lib.DPrintln("Process result:", a, err)
}
