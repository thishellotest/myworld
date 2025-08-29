package biz

import "fmt"

// ClientContractFolderNameForBox It's box client contract folder id
func ClientContractFolderNameForBox(firstName string, lastName string, clientId int32) string {
	return fmt.Sprintf("%s, %s #%d", lastName, firstName, clientId)
}

// ClientFolderNameForBox It's box client folder id
func ClientFolderNameForBox(firstName string, lastName string) string {
	return fmt.Sprintf("VBC - %s, %s", lastName, firstName)
}

// ClientCaseDataCollectionFolderNameForBox It's box client folder id
func ClientCaseDataCollectionFolderNameForBox(firstName string, lastName string, clientCaseId int32) string {
	return fmt.Sprintf("%s, %s #%d", lastName, firstName, clientCaseId)
}
