package tests

var FilebuzTriggerFIleUploaded = `{"type":"webhook_event","id":"06d58ee4-9bf9-4a7b-8f71-184a26a4c06f","created_at":"2025-05-11T20:22:42-07:00","trigger":"FILE.UPLOADED","webhook":{"id":"4220307672","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"source":{"id":"1859497976800","type":"file","file_version":{"type":"file_version","id":"2050398896800","sha1":"0629d7678c0025f646309833745c18bc6366507b"},"sequence_id":"0","etag":"0","sha1":"0629d7678c0025f646309833745c18bc6366507b","name":"4444.jpg","description":"","size":32245,"path_collection":{"total_count":7,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183605424","sequence_id":"4","etag":"4","name":"VBC All Teams"},{"type":"folder","id":"263406803830","sequence_id":"1","etag":"1","name":"Data Collection"},{"type":"folder","id":"311268435813","sequence_id":"4","etag":"4","name":"TestL2, Test1 #5511"},{"type":"folder","id":"311268440613","sequence_id":"0","etag":"0","name":"Record Review"},{"type":"folder","id":"315732993509","sequence_id":"0","etag":"0","name":"VA Medical Records"},{"type":"folder","id":"320740557022","sequence_id":"0","etag":"0","name":"sub folder"}]},"created_at":"2025-05-11T20:22:42-07:00","modified_at":"2025-05-11T20:22:42-07:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-08-14T17:26:33-07:00","content_modified_at":"2024-08-14T17:26:33-07:00","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"320740557022","sequence_id":"0","etag":"0","name":"sub folder"},"item_status":"active"},"additional_info":[]}`
var FilebuzTriggerFileRename = `{
	"type": "webhook_event",
	"id": "04b49601-064e-475a-87f5-c90d79eed37d",
	"created_at": "2025-05-11T18:00:55-07:00",
	"trigger": "FILE.RENAMED",
	"webhook": {
		"id": "4220307672",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "1859359115815",
		"type": "file",
		"file_version": {
			"type": "file_version",
			"id": "2050237715815",
			"sha1": "633586e64063af33ac643f42cbbec2367ce3b075"
		},
		"sequence_id": "1",
		"etag": "1",
		"sha1": "633586e64063af33ac643f42cbbec2367ce3b075",
		"name": "111122_new_name.jpg",
		"description": "",
		"size": 1207935,
		"path_collection": {
			"total_count": 6,
			"entries": [{
				"type": "folder",
				"id": "0",
				"sequence_id": null,
				"etag": null,
				"name": "All Files"
			}, {
				"type": "folder",
				"id": "241183605424",
				"sequence_id": "4",
				"etag": "4",
				"name": "VBC All Teams"
			}, {
				"type": "folder",
				"id": "263406803830",
				"sequence_id": "1",
				"etag": "1",
				"name": "Data Collection"
			}, {
				"type": "folder",
				"id": "311268435813",
				"sequence_id": "4",
				"etag": "4",
				"name": "TestL2, Test1 #5511"
			}, {
				"type": "folder",
				"id": "311268440613",
				"sequence_id": "0",
				"etag": "0",
				"name": "Record Review"
			}, {
				"type": "folder",
				"id": "315732993509",
				"sequence_id": "0",
				"etag": "0",
				"name": "VA Medical Records"
			}]
		},
		"created_at": "2025-05-11T17:58:26-07:00",
		"modified_at": "2025-05-11T18:00:55-07:00",
		"trashed_at": null,
		"purged_at": null,
		"content_created_at": "2025-05-09T03:42:23-07:00",
		"content_modified_at": "2025-05-09T03:42:23-07:00",
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
			"id": "315732993509",
			"sequence_id": "0",
			"etag": "0",
			"name": "VA Medical Records"
		},
		"item_status": "active"
	},
	"additional_info": {
		"old_name": "111122.jpg"
	}
}`

var FilebuzTriggerFolderMove = `{
	"type": "webhook_event",
	"id": "6d6da1b4-5244-49da-bc7c-8c9030dae507",
	"created_at": "2025-05-11T20:30:00-07:00",
	"trigger": "FOLDER.MOVED",
	"webhook": {
		"id": "4220307672",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "320740557022",
		"type": "folder",
		"sequence_id": "1",
		"etag": "1",
		"name": "sub folder",
		"created_at": "2025-05-11T18:02:09-07:00",
		"modified_at": "2025-05-11T20:25:24-07:00",
		"description": "",
		"size": 32245,
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
				"id": "241183605424",
				"sequence_id": "4",
				"etag": "4",
				"name": "VBC All Teams"
			}, {
				"type": "folder",
				"id": "263406803830",
				"sequence_id": "1",
				"etag": "1",
				"name": "Data Collection"
			}, {
				"type": "folder",
				"id": "311268435813",
				"sequence_id": "4",
				"etag": "4",
				"name": "TestL2, Test1 #5511"
			}, {
				"type": "folder",
				"id": "311268440613",
				"sequence_id": "0",
				"etag": "0",
				"name": "Record Review"
			}, {
				"type": "folder",
				"id": "315732993509",
				"sequence_id": "0",
				"etag": "0",
				"name": "VA Medical Records"
			}, {
				"type": "folder",
				"id": "320766422256",
				"sequence_id": "0",
				"etag": "0",
				"name": "sub folder2"
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
		"content_created_at": "2025-05-11T18:02:09-07:00",
		"content_modified_at": "2025-05-11T20:25:24-07:00",
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
			"id": "320766422256",
			"sequence_id": "0",
			"etag": "0",
			"name": "sub folder2"
		},
		"item_status": "active"
	},
	"additional_info": {
		"after": {
			"id": "320766422256",
			"type": "folder"
		},
		"before": {
			"id": "315732993509",
			"type": "folder"
		}
	}
}`

var FilebuzTriggerTrashed = `{
	"type": "webhook_event",
	"id": "1f00344e-e4ea-4391-a612-f79c5a8d4f57",
	"created_at": "2025-05-11T22:35:32-07:00",
	"trigger": "FILE.TRASHED",
	"webhook": {
		"id": "4220307672",
		"type": "webhook"
	},
	"created_by": {
		"type": "user",
		"id": "30888625898",
		"name": "VBC Team",
		"login": "info@vetbenefitscenter.com"
	},
	"source": {
		"id": "1859597022426",
		"type": "file"
	},
	"additional_info": []
}`
