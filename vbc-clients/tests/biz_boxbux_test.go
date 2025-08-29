package tests

import (
	"fmt"
	"strings"
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_box"
	"vbc/lib"
)

/*
err: 263252423092 263407432521 ok1
err: 261849118832 263409220477 ok1
err: 263424822143 263409383635 ok1
err: 262831893453 263409085591 ok1 (文件夹冲突，重点关注)
err: 260944375133 263407544468 ok1

262836135891 263407893513 ok1
260937542087 263408315542 ok1
260939600452 263409090018 ok1
260938403992 263407789701 ok1
260943786958 263409003856 ok1
260939834336 263407513317 ok1
260939008536 263409085859 ok1
262564917168 263409189404 ok1

err: 261514026402 263408941892 ok1
err: 261717551273 263408622377 ok1
err: 260937927949 263408562555 ok1
err: 260939356241 263407717756 ok1
err: 260940515467 263409584864 ok1
err: 260940914337 263408383080 ok1
err: 262601317810 263409011236 ok1
err: 260942313293 263407378937 ok1
err: 261815812608 263407759064 ok1
err: 260942586693 263407622264 ok
err: 260939281799 263409369154 ok1
err: 263250198764 263408536697 ok1
err: 261957561306 263407539990 ok1
err: 261683043974 263407549727 ok1
err: 261821085474 263407777243 ok1
err: 261632097331 263408291548 ok1
err: 261719322823 263409109595 ok1 Andrew, Smith
err: 261962365453 263408466241 ok1 (有多个cases)  Smith, Christopher Manuel #5026
err: 260940454986 263408486215 ok1
err: 263257386442 263408173982 ok1 (重名问题) Andrew, Smith
err: 260937572816 263408879881 ok1
err: 260939913450 263408877126 ok1
*/
func Test_BoxbuzUsecase_MergeFolder(t *testing.T) {
	//err := UT.BoxbuzUsecase.MergeFolder("260944375133", "263407544468")

	//err := UT.BoxbuzUsecase.MergeFolder("262836135891", "263407893513")
	//err := UT.BoxbuzUsecase.MergeFolder("260937542087", "263408315542")
	//err := UT.BoxbuzUsecase.MergeFolder("260939600452", "263409090018")
	//err := UT.BoxbuzUsecase.MergeFolder("260938403992", "263407789701")
	//err := UT.BoxbuzUsecase.MergeFolder("260943786958", "263409003856")
	//err := UT.BoxbuzUsecase.MergeFolder("260939834336", "263407513317")
	//err := UT.BoxbuzUsecase.MergeFolder("260939008536", "263409085859")
	//err := UT.BoxbuzUsecase.MergeFolder("262564917168", "263409189404")

	//err := UT.BoxbuzUsecase.MergeFolder("261514026402", "263408941892") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261717551273", "263408622377") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260937927949", "263408562555") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260939356241", "263407717756") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260940515467", "263409584864") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260940914337", "263408383080") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("262601317810", "263409011236") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260942313293", "263407378937") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261815812608", "263407759064") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260942586693", "263407622264") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260939281799", "263409369154") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("263250198764", "263408536697") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261957561306", "263407539990") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261683043974", "263407549727") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261821085474", "263407777243") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261632097331", "263408291548") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261719322823", "263409109595") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("261962365453", "263408466241") // ok (有多个cases)  Smith, Christopher Manuel #5026
	//err := UT.BoxbuzUsecase.MergeFolder("260940454986", "263408486215") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("263257386442", "263408173982") // ok (重名问题) Andrew, Smith
	//err := UT.BoxbuzUsecase.MergeFolder("260937572816", "263408879881") // ok
	//err := UT.BoxbuzUsecase.MergeFolder("260939913450", "263408877126") // ok

	//lib.DPrintln("error:", err)
}

func Test_BoxbuzUsecase_RenameDataCollection(t *testing.T) {
	//err := UT.BoxbuzUsecase.RenameDataCollection()
	//lib.DPrintln(err)
}

func Test_BoxbuzUsecase_RenameClientCasesFolderName(t *testing.T) {
	//err := UT.BoxbuzUsecase.RenameClientCasesFolderName()
	//if err != nil {
	//	panic(err)
	//}
	//lib.DPrintln(err)
}

func Test_BoxbuzUsecase_DCRecordReviewFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	dCRecordReviewFolderId, err := UT.BoxbuzUsecase.DCRecordReviewFolderId(tCase)
	lib.DPrintln(dCRecordReviewFolderId)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_HandleClientFolder(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	err := UT.BoxbuzUsecase.HandleClientFolder(tCase)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_HandleDCFolder(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	err := UT.BoxbuzUsecase.HandleDCFolder(tCase)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_GetDCSubFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	aaa, err := UT.BoxbuzUsecase.GetDCSubFolderId(biz.MapKeyBuildAutoBoxDCQuestionnairesFolderId(
		tCase.CustomFields.NumberValueByNameBasic("id")), tCase)
	lib.DPrintln(aaa, err)
}

func Test_BoxbuzUsecase_GetDCSubFolderId_DCPrivateExamsFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	aaa, err := UT.BoxbuzUsecase.GetDCSubFolderId(biz.MapKeyBuildAutoBoxDCPrivateExamsFolderId(
		tCase.CustomFields.NumberValueByNameBasic("id")), tCase)
	lib.DPrintln(aaa, err)
}

func Test_BoxbuzUsecase_GetOrMakeClaimsAnalysisFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5446)
	aaa, err := UT.BoxbuzUsecase.GetOrMakeClaimsAnalysisFolderId(tCase)
	lib.DPrintln(aaa, err)
}

func Test_BoxbuzUsecase_DCPersonalStatementsFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5293)
	a1, err := UT.BoxbuzUsecase.DCPersonalStatementsFolderId(tCase)
	lib.DPrintln(a1, err)
}

func Test_BoxbuzUsecase_PersonalStatementDocFileBoxFileId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5293)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	a, b, err := UT.BoxbuzUsecase.PersonalStatementDocFileBoxFileId(tClient, tCase)
	lib.DPrintln(a, b, err)
}

func Test_BoxbuzUsecase_CPersonalStatementsFolderIdByAnyCase(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	a1, err := UT.BoxbuzUsecase.CPersonalStatementsFolderIdByAnyCase(*tCase)
	lib.DPrintln(a1, err)
}

func Test_BoxbuzUsecase_CPersonalStatementsFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5649)
	a1, err := UT.BoxbuzUsecase.CPersonalStatementsFolderId(tCase)
	lib.DPrintln(a1, err)
}

func Test_BoxbuzUsecase_HandlePersonalStatementsFile(t *testing.T) {
	err := UT.BoxbuzUsecase.DoPersonalStatementsFile(5369)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_DoClaimsAnalysisFile(t *testing.T) {
	err := UT.BoxbuzUsecase.DoClaimsAnalysisFile(5446)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_DoDocEmailFile(t *testing.T) {
	err := UT.BoxbuzUsecase.DoDocEmailFile(5511)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_SameNameFolderOrFile(t *testing.T) {
	aaa, err := UT.BoxbuzUsecase.SameNameFolderOrFile("folder", "Test Jotform", "264658751993")
	lib.DPrintln(aaa, err)
}
func Test_BoxbuzUsecase_SameNameFolderOrFile_file(t *testing.T) {
	aaa, err := UT.BoxbuzUsecase.SameNameFolderOrFile("file", "b.pdf", "264924117433")
	lib.DPrintln(aaa, err)
}

/*
__ 260391738256 345
__ 257399924492 355
__ 260435037274 351
__ 260395319179 66
__ 257604633337 391
*/
func Test_aaaa_ListItemsInFolder(t *testing.T) {
	destRes, err := UT.BoxUsecase.ListItemsInFolder("241109085470")
	if err != nil {
		panic(err)
	}
	destResMap := lib.ToTypeMapByString(*destRes)
	destResEntries := destResMap.GetTypeList("entries")
	for _, v := range destResEntries {
		sourceName := v.GetString("name")
		if strings.Index(sourceName, "#") >= 0 {
			res := strings.Split(sourceName, "#")
			if len(res) != 2 {
				panic("sourceName error: " + sourceName)
			}
			a := lib.InterfaceToInt32(res[1])
			if a < 5000 {
				fmt.Println("__", v.GetString("id"), a)
			}
		}
	}
}

func Test_dcFolderSubs(t *testing.T) {

	//dcFolderSubs := `[{"etag":"2","id":"263407033477","name":"Abutin, Niko Ralphluis","sequence_id":"2","type":"folder"},{"etag":"1","id":"263407629244","name":"Acuario, Edralin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409690475","name":"Albrecht, Keith Richard","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407432521","name":"Alcantar, Francisco Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409699906","name":"Alexander, Keith David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409220477","name":"Alexander, Troy Don","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409383635","name":"Allen, Robert Joseph","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408142981","name":"Ancho, Romulo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409340623","name":"Anderson, Jilleah","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408260360","name":"Anderson, Trever Shaw","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409604748","name":"Andrews, Jamaal","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085591","name":"Angulo, Don Clark","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408638849","name":"Arellano, Hector Gibram","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409608463","name":"Ayala, Joe","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407544468","name":"Ayuyao, Bernard Elijah","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407775301","name":"Bailey, Bernard","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408325137","name":"Baker, Jeffrey Stephen Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407895271","name":"Balavram, Jason Sy","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408186009","name":"Baldemeca, Noel Genetia","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407002935","name":"Ballesteros, Judeasar Galapon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408572372","name":"Banuex, Noemi","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409027699","name":"Barajas, Emery Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407850270","name":"Barbosa, Brendan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407864388","name":"Barnes, Santana Venique","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408646454","name":"Barnett, Jovan Eric","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408332367","name":"Battle, Rodney Terez","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408192682","name":"Bautista, Alexander Clement","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407518222","name":"Becerra, Jonathan Contreras","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408711489","name":"Beeler, Rebecca #5019","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408505251","name":"Blaine, Christopher Charles","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408101803","name":"Boden, Boyce Robert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408584289","name":"Bolino, Louis Alto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408994347","name":"Briggs, Scott Christopher","sequence_id":"1","type":"folder"},{"etag":"0","id":"263966723631","name":"Brooks, Roy #5069","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408480938","name":"Brown, Terence Lewis","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408941892","name":"Burgess, Dominique Farrell","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407893513","name":"Camacho, Pedro","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408432830","name":"Campbell, Dailyn","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408809561","name":"Canseco, Manuel Valiente","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408946196","name":"Carr, Christopher","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407168436","name":"Carreon, Lovelito Flores","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408315542","name":"Carrillo, Adrian","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407400559","name":"Carter, Gabrielle Alexandria","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409575210","name":"Castelluccio, Dillon Thomas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409090018","name":"Castillo, Jacinto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408282334","name":"Castro, Roman Lucas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408953729","name":"Cavaliere, Robynn","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408315294","name":"Chacon, Guy #5052","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409193805","name":"Claire, Seth Alexander","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409058831","name":"Cobian Jr., Jose","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409335804","name":"Coley, Alvin Lorenzo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407511190","name":"Collins, Samantha Rose","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408706463","name":"Cook, Robert #5020","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408775705","name":"D'Alessandro, James Scott","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408622377","name":"De Leon, Cecilia","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409573530","name":"DelPrete, Desiree","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408560027","name":"Demps, Kelvin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409486296","name":"Deocampo, Teddy Baysan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408562555","name":"Devine, Justin Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409697401","name":"DiBenedetto, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408598259","name":"Dickerson, Ramon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408170942","name":"Dickey, James #5005","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408262134","name":"Dishmon, Varian Dione","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407789701","name":"Dodd, Brent","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408883616","name":"Dukes, Spencer Patrick","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409395451","name":"Dunkin, John Steven","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408255849","name":"Edwards, Addison Chase","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407679255","name":"FaisonLanier, Terrica","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408082641","name":"Fannin, Quentin Dekote","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407698819","name":"Farmer, Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407825270","name":"Flores, Jacinth Aaron","sequence_id":"1","type":"folder"},{"etag":"2","id":"263409236508","name":"Fowler, Casey #345","sequence_id":"2","type":"folder"},{"etag":"1","id":"263409256573","name":"Fowler, Jerry Joseph","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407897776","name":"Franco, Alec Robert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409003856","name":"Galac, Cesar","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408776073","name":"Garcia, Rey #5041","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408716059","name":"Garlejo, Jason Doctolero","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409049595","name":"Gilmore, Earl Glenn Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408759623","name":"Gonzales, Beau Matthew","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407379970","name":"Goodson, Augustus Ivan IV","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409028000","name":"Green, Antionette","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407516742","name":"Green, Donnell","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409561453","name":"Green, Sandra","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407261628","name":"Griffin, Major Pete III","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407014935","name":"Haley, George Walter Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407237574","name":"Harris, Debra Lee","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408234122","name":"Herron, Leslie Rhea","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409282194","name":"Hiers, Tommy Lamar Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408977216","name":"Ho, Duyet #5043","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407717756","name":"Houston, Keith Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407411164","name":"Howard, Nicole Marie","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408903510","name":"Huerta Jr., Edward David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408925455","name":"Huerta Jr., Edward David #96","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407497275","name":"Huerta Sr., Edward David #260","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407386390","name":"Hutchinson, Derek","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408325009","name":"Ibarrondo, William","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408646453","name":"Ibarrondo, William Basean","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407665390","name":"Inzer, Russell Dustin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408617842","name":"Jacob, Melvin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408202281","name":"James, Dillon Randall","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409584864","name":"Johnson, Christopher #5009","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409280743","name":"Johnson, Jermaine","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409239761","name":"Johnson, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407513317","name":"Johnson, Robin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408383080","name":"Jones, Cyril Evans Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407472517","name":"Jones, Naurice","sequence_id":"1","type":"folder"},{"etag":"0","id":"264597201436","name":"Kallmeyer, Michael #5072","sequence_id":"0","type":"folder"},{"etag":"1","id":"263409011236","name":"Keith, Kristopher #5025","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407636399","name":"Keller, Tony Jordan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407378937","name":"Kennedy, Payton Gary","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407672521","name":"Kubas III, William Philip","sequence_id":"1","type":"folder"},{"etag":"0","id":"264171116583","name":"Kunkowski, James #5070","sequence_id":"0","type":"folder"},{"etag":"1","id":"263407746455","name":"Lane, Michael John","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408048756","name":"Lang, Shanise Eileen","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407759064","name":"Lastrella, Amando Llorin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407079735","name":"Laxa, Eduardo Dulu","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407578819","name":"Le, Khanh Si","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409409666","name":"Lepe, Christopher","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407774875","name":"Liggans, Nyjerus Lavondai Onijar","sequence_id":"1","type":"folder"},{"etag":"0","id":"264108189026","name":"Long, Jonathan #5068","sequence_id":"0","type":"folder"},{"etag":"0","id":"263965342942","name":"Luna, Larzen #355","sequence_id":"0","type":"folder"},{"etag":"0","id":"264296809567","name":"Maldonado, Jimmy #5060","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408704030","name":"Mangra, Robbi #5023","sequence_id":"1","type":"folder"},{"etag":"2","id":"263409358994","name":"Marcial, Janry #351","sequence_id":"2","type":"folder"},{"etag":"1","id":"263408360692","name":"Mendez, Marco Antonio","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407919167","name":"Montoya, Mario","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408174156","name":"Montoya, Mary","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407556962","name":"Moreno, Michael Joey","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408493117","name":"Morris, Lee Roy","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407593367","name":"Murrell, Brashaad","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408197832","name":"Myers, Luis Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408754259","name":"Nellis, Lawrence #5024","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085859","name":"Netemeyer, Aaron","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409335873","name":"Newman, James Wesley","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407708055","name":"Olivetti, Anthony Ryan Borja","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409189404","name":"Orias, Ricardo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407676838","name":"Padua, Michael Daniel","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409107750","name":"Peca, Jason","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407622264","name":"Perez, Joaquin Xavier","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409369154","name":"Perry, Fred Douglas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408536697","name":"Petit-Frere, Alexandre Freud","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408613185","name":"Pharnes, Eric Dwayne","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408872130","name":"Pierre, Gilbert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408335098","name":"Prado, Martius Oris","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409357404","name":"Pratko, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407645792","name":"Provasek, Jared #5046","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407539990","name":"Ralat, Carlos Alberto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407549727","name":"Reynoso, Algis #5036","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407439423","name":"Rios, Jonathan David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408483369","name":"Rivera, Louie","sequence_id":"1","type":"folder"},{"etag":"2","id":"263408128360","name":"Rosales, Juan #66","sequence_id":"2","type":"folder"},{"etag":"1","id":"263409671231","name":"Rutledge, Ronnie #5042","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408531552","name":"Salb, Austin Reid","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407996012","name":"Santillan-Mondaca, Kristy #5050","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407479863","name":"Sayles, Matthew Evans","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408282344","name":"Serrano, Ronald","sequence_id":"1","type":"folder"},{"etag":"0","id":"263964486190","name":"Sese, Maria #5002","sequence_id":"0","type":"folder"},{"etag":"1","id":"263407777243","name":"Shelrud, Cierra Kay","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408291548","name":"Sida, Andrew Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409143090","name":"Slater, Jamie","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408173982","name":"Smith, Andrew","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409109595","name":"Smith, Andrew #5011","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409138194","name":"Smith, Austin Cole","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408466241","name":"Smith, Christopher Manuel #5026","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408730987","name":"Smith, Max","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409028559","name":"Smith, Rhyheime","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408486215","name":"Smith, Zane Eugene","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407441338","name":"Smolinski, Donald Jay","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085790","name":"Stacks, Taylor","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408843374","name":"Stewart, Robert #5038","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408879881","name":"Stuart, James Francis","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409189108","name":"Summer, Aaron Michael Blake","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408183176","name":"Sutton, Derrick","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408877126","name":"Tanquilut, Remigio Dimalanta","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408462333","name":"Taylor, Bobbee Nykole","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408290941","name":"Terrell, Conley II","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408989699","name":"TestLi, TestShi #5057","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408579201","name":"TestLi, TestShi #5058","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407556960","name":"TestLiao, TestGary #5061","sequence_id":"1","type":"folder"},{"etag":"0","id":"263470914421","name":"TestLiao, TestGary #5064","sequence_id":"0","type":"folder"},{"etag":"1","id":"263409010923","name":"Thompson, Jason Scott","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408644030","name":"Thrower, Tony Lee","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408560526","name":"Tran, Danny Minh","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408893699","name":"Tran, Joanne Perea","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409387704","name":"Tran, Kenny #5045","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408661172","name":"Valdez, Jacob","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408233836","name":"Valdez, Joshua Raul","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407744765","name":"Valli, Matthew Lawrence","sequence_id":"1","type":"folder"},{"etag":"2","id":"263408943728","name":"Valli, Ronald #391","sequence_id":"2","type":"folder"},{"etag":"1","id":"263408636453","name":"Vargas, Gabrian","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408008709","name":"Velasquez, Jose David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407730033","name":"Walker, Ronald Stanley","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408847871","name":"Warren, Olga","sequence_id":"1","type":"folder"},{"etag":"0","id":"264219349997","name":"Warren, Olga #5071","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408262984","name":"Watts, Patrick Levon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407622121","name":"Webster, Craig","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408430635","name":"West, Melissa Theresa","sequence_id":"1","type":"folder"},{"etag":"0","id":"263723175862","name":"Westerveld, Michael #5017","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408356225","name":"White, Michael #5015","sequence_id":"1","type":"folder"},{"etag":"0","id":"263570678121","name":"Wirth, David #5029","sequence_id":"0","type":"folder"}]`
	//res := lib.ToTypeListByString(dcFolderSubs)
	//i := 0
	//for _, v := range res {
	//
	//	folderId := v.GetString("id")
	//	name := v.GetString("name")
	//	if strings.Index(name, "#") < 0 {
	//		fmt.Println(folderId, name)
	//		var entity biz.MapEntity
	//		err := UT.MapUsecase.CommonUsecase.DB().Where("mval=?", folderId).Take(&entity).Error
	//		if err != nil {
	//			panic(err)
	//			fmt.Println("folderId:", folderId, err, "====")
	//		} else {
	//			aa := strings.Split(entity.Mkey, ":")
	//			if len(aa) == 3 && aa[1] == "DataCollectionFolderId" {
	//
	//				if aa[2] == "103" || aa[2] == "40" || aa[2] == "290" || aa[2] == "50" || aa[2] == "342" || aa[2] == "49" {
	//					continue
	//				}
	//
	//				newName := name + " #" + aa[2]
	//				fmt.Println(entity.Mkey, "folderId:", folderId, "oldName:", name, "newName:", newName)
	//				_, err = UT.BoxUsecase.UpdateFolderName(folderId, newName)
	//				if err != nil {
	//					panic(err)
	//				}
	//				i++
	//				if i >= 2 {
	//					//break
	//				}
	//			} else {
	//				panic("err")
	//			}
	//
	//		}
	//
	//	}
	//}
}

func Test_BoxbuzUsecase_GetBoxResId_PatientPaymentForm(t *testing.T) {
	boxResId, err := UT.BoxbuzUsecase.GetBoxResId("264658754393",
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PE_PatientPaymentForm_File),
		5004,
	)
	lib.DPrintln(boxResId, err)
}

func Test_BoxbuzUsecase_GetBoxResId(t *testing.T) {
	boxResId, err := UT.BoxbuzUsecase.GetBoxResId("264658749593",
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_RV_PrivateMedicalRecords_Folder),
		5004,
	)
	lib.DPrintln(boxResId, err)
}

func Test_BoxbuzUsecase_GetBoxResId_PatientPaymentForm_File(t *testing.T) {
	boxResId, err := UT.BoxbuzUsecase.GetBoxResId("264658754393",
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PE_PatientPaymentForm_File),
		5004,
	)
	lib.DPrintln(boxResId, err)
}

func Test_BoxbuzUsecase_DoCopyDocEmailFile(t *testing.T) {
	err := UT.BoxbuzUsecase.DoCopyDocEmailFile(5369)
	lib.DPrintln("err: ", err)
}

func Test_BoxbuzUsecase_DoCopyReadPriorToYourDoctorVisitFile(t *testing.T) {
	err := UT.BoxbuzUsecase.DoCopyReadPriorToYourDoctorVisitFile(5369)
	lib.DPrintln("err: ", err)
}

func Test_BoxbuzUsecase_PSDocEmailFileId(t *testing.T) {

	caseId := int32(5511)
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, caseId)
	fileId, err := UT.BoxbuzUsecase.DCPSDocEmailFileId(tCase)
	lib.DPrintln(fileId, err)
}

func Test_BoxbuzUsecase_GetBoxResId_PatientPaymentForm_File1(t *testing.T) {

	caseId := int32(5369)
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, caseId)

	psFolderId, err := UT.BoxbuzUsecase.DCPersonalStatementsFolderId(tCase)
	lib.DPrintln("psFolderId: ", psFolderId, err)
	boxResId, err := UT.BoxbuzUsecase.GetBoxResId(psFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PS_DocEmail_File),
		caseId,
	)

	lib.DPrintln("CaseRelaBox_DC_PS_DocEmail_File:", boxResId, err)
}

func Test_BoxbuzUsecase_CopyBoxResItemsToFolder(t *testing.T) {

	items, err := UT.BoxUsecase.ListItemsInFolderFormat("267683828516")
	lib.DPrintln("err:", err)

	err = UT.BoxbuzUsecase.CopyBoxResItemsToFolder("271859979627", items, biz.CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
	lib.DPrintln(err)
}

func Test_BoxbuzUsecase_CreateFolderByEntries(t *testing.T) {
	res, code, err := UT.BoxUsecase.GetFileInfoForTypeMap("1559300046485")
	lib.DPrintln(code, err)
	lib.DPrintln(res)

	i, path, err := UT.BoxbuzUsecase.CreateFolderByEntries(5, res.GetTypeList("path_collection.entries"), "267683828516")
	lib.DPrintln(i, path, err)
	//269734300185
}

func Test_BoxbuzUsecase_CopyFileToFolderNoCover(t *testing.T) {
	f, newFile, err := UT.BoxbuzUsecase.CopyFileToFolderNoCover("267683828516", "1559300046485", "a.pdf")
	lib.DPrintln(f, newFile, err)
}

func Test_BoxbuzUsecase_DC_PE_PsychFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)

	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))

	psychFolderId, err := UT.BoxbuzUsecase.FolderIdDC_PE_Psych(tCase, tClient)
	lib.DPrintln(err)
	lib.DPrintln(psychFolderId)
}

func Test_BoxbuzUsecase_DC_PE_GeneralFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)

	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))

	psychFolderId, err := UT.BoxbuzUsecase.FolderIdDC_PE_General(tCase, tClient)
	lib.DPrintln(err)
	lib.DPrintln(psychFolderId)
}

func Test_BoxbuzUsecase_RealtimeCPersonalStatementsDocxFileId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	cPersonalStatementsFolderId, name, boxFileId, err := UT.BoxbuzUsecase.RealtimeCPersonalStatementsDocxFileId(*tClient, *tCase)
	lib.DPrintln(cPersonalStatementsFolderId, name, boxFileId, err)
}

func Test_BoxbuzUsecase_DCQuestionnairesFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	a, err := UT.BoxbuzUsecase.DCQuestionnairesFolderId(*tCase)
	lib.DPrintln(a, err)
}

func Test_BoxbuzUsecase_DCGetOrMakeUpdateQuestionnairesFolderId(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	a, err := UT.BoxbuzUsecase.DCGetOrMakeUpdateQuestionnairesFolderId(*tCase)
	lib.DPrintln(a, err)
}

func Test_BoxbuzUsecase_CopyFolderSubsToFolder(t *testing.T) {
	err := UT.BoxbuzUsecase.CopyFolderSubsToFolder("241112327202", "329274083103")
	lib.DPrintln(err)
}
