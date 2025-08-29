package biz

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

type GopdfUsecase struct {
	log             *log.Helper
	CommonUsecase   *CommonUsecase
	conf            *conf.Data
	ResourceUsecase *ResourceUsecase
}

func NewGopdfUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ResourceUsecase *ResourceUsecase) *GopdfUsecase {
	uc := &GopdfUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		ResourceUsecase: ResourceUsecase,
	}

	return uc
}

type CreateMedicalTeamFormVo struct {
	ClientName          string
	Address             string
	Location            string
	Dob                 string
	Ssn                 string
	Phone               string
	Email               string
	Itf                 string
	PrivateExamsNeededS []string
}

type CreateContractVo struct {
	ClientName  string
	ClientEmail string
	VsName      string
	VsEmail     string
}

//type CreateContractAmVo struct {
//	ClientName           string
//	RepresentativeName   string
//	RepresentativeDate   string
//	RepresentativeNumber string
//}

func (c *GopdfUsecase) CreatePdfFromImg(imgs [][]byte) (pdfFile []byte, err error) {
	if len(imgs) == 0 {
		return nil, errors.New("imgs is empty")
	}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})
	for k, _ := range imgs {
		pdf.AddPage()
		imgH1, err := gopdf.ImageHolderByBytes(imgs[k])
		if err != nil {
			return nil, err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 612, H: 792}); err != nil {
			return nil, err
		}
	}
	return pdf.GetBytesPdf(), nil
}

func (c *GopdfUsecase) CreatePdfFromImgFiles(filePaths []string) (pdfFile []byte, err error) {
	if len(filePaths) == 0 {
		return nil, errors.New("filePaths is empty")
	}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})
	for _, v := range filePaths {
		pdf.AddPage()
		imgH1, err := gopdf.ImageHolderByPath(v)
		if err != nil {
			return nil, err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 612, H: 792}); err != nil {
			return nil, err
		}
	}
	return pdf.GetBytesPdf(), nil
}

func (c *GopdfUsecase) CreateContractAm(contractTime time.Time, contractVetVo ContractVetVo, contractAttorneyVo ContractAttorneyVo) (pdfFile []byte, err error) {

	contractTime = contractTime.In(configs.GetVBCDefaultLocation())
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})

	for i := 1; i <= 9; i++ {
		if i == 6 {
			continue
		}
		pdf.AddPage()
		ttfName := "Arial"
		err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/Arial.ttf")
		if err != nil {
			return nil, err
		}
		err = pdf.SetFont(ttfName, "", 10)
		if err != nil {
			return nil, err
		}
		pageJpgPath := c.ResourceUsecase.ResPath() + "/ContractsAM_V2/page-000" + InterfaceToString(i) + ".jpg"

		//image bytes
		b, err := os.ReadFile(pageJpgPath)
		if err != nil {
			return nil, err
		}
		imgH1, err := gopdf.ImageHolderByBytes(b)
		if err != nil {
			return nil, err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 612, H: 792}); err != nil {
			return nil, err
		}
		//// Color the page
		//pdf.SetLineWidth(0.1)
		//pdf.SetFillColor(124, 252, 0) //setup fill color
		//pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
		//pdf.SetFillColor(0, 0, 0)

		//contractAttorneyVo := ContractAttorneyVo{
		//	Street:              "KLDP, LLP, 250 W 34th St, Floor 3",
		//	City:                "New York",
		//	Province:            "NY",
		//	ZipCode:             "10119",
		//	FirstName:           "Phuong",
		//	LastName:            "Le",
		//	AccreditationNumber: "59046",
		//	Email:               "contact@augustusmiles.com",
		//}
		//
		//contractVetVo := ContractVetVo{
		//	Ssn:      "411-89-1683",
		//	Dob:      "1969-03-22",
		//	Street:   "321, Kensington Dr 321, Kensington Dr 321, Kensington Dr",
		//	City:     "New marketNew marketNew marketNew marketNew market",
		//	Province: "NW",
		//	ZipCode:  "223424242",
		//	Phone:    "8056604465",
		//	Email:    "111aa8056604465@gmail.com",
		//}
		ssnForContract, err := FormatSsnForContract(contractVetVo.Ssn)
		if err != nil {
			return nil, err
		}

		if i == 1 {
			err = c.CreateContractAMPage1(&pdf, contractTime, contractVetVo, contractAttorneyVo)
			if err != nil {
				return nil, err
			}
		} else if i == 5 {
			err = c.CreateContractAmPage6(&pdf)
			if err != nil {
				return nil, err
			}
		} else if i == 7 {
			err = c.CreateContractAMPage7(&pdf, contractAttorneyVo, ssnForContract, contractVetVo)
			if err != nil {
				return nil, err
			}
		} else if i == 8 {
			err = c.CreateContractAMPage8(&pdf, contractAttorneyVo, ssnForContract)
			if err != nil {
				return nil, err
			}
			for k, v := range ssnForContract.First {
				pdf.SetXY(160+(float64(k)*17.25), 27.5)
				pdf.Cell(nil, v)
			}
			for k, v := range ssnForContract.Middle {
				pdf.SetXY(227+(float64(k)*17.25), 27.5)
				pdf.Cell(nil, v)
			}
			for k, v := range ssnForContract.Last {
				pdf.SetXY(273+(float64(k)*17.25), 27.5)
				pdf.Cell(nil, v)
			}
			pdf.SetXY(187, 98)
			pdf.Cell(nil, "U")
			pdf.SetXY(187+16.5, 98)
			pdf.Cell(nil, "S")
		} else if i == 9 {
			for k, v := range ssnForContract.First {
				pdf.SetXY(160+(float64(k)*17.25), 38)
				pdf.Cell(nil, v)
			}
			for k, v := range ssnForContract.Middle {
				pdf.SetXY(227+(float64(k)*17.25), 38)
				pdf.Cell(nil, v)
			}
			for k, v := range ssnForContract.Last {
				pdf.SetXY(273+(float64(k)*17.25), 38)
				pdf.Cell(nil, v)
			}
			err = c.CreateContractAmPage9(&pdf)
			if err != nil {
				return nil, err
			}
		}
	}

	return pdf.GetBytesPdf(), nil
}

func (c *GopdfUsecase) CreateContractAmPage9(pdf *gopdf.GoPdf) error {
	ttfName := "Arial"
	pdf.SetTextColor(255, 255, 255)
	//pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 36)
	if err != nil {
		return err
	}
	pdf.SetXY(70, 140)
	pdf.Cell(nil, "[[s|1]]")

	//err = pdf.SetFont(ttfName, "", 12)
	//if err != nil {
	//	return err
	//}
	//pdf.SetXY(408, 145)
	//pdf.Cell(nil, "[[d|1              ]]")

	err = pdf.SetFont(ttfName, "", 36)
	if err != nil {
		return err
	}
	pdf.SetXY(70, 240)
	pdf.Cell(nil, "[[s|2]]")

	//err = pdf.SetFont(ttfName, "", 12)
	//if err != nil {
	//	return err
	//}
	//pdf.SetXY(408, 245)
	//pdf.Cell(nil, "[[d|2              ]]")
	//pdf.SetTextColor(0, 0, 0)

	pdf.SetTextColor(0, 0, 0)
	err = pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}

	aa := time.Now().In(configs.GetVBCDefaultLocation())
	signTime := aa.Format("01022006")
	signTimes := strings.Split(signTime, "")
	for k, v := range signTimes {
		if k <= 1 {
			pdf.SetXY(411+(float64(k)*17.1), 162)
		} else if k <= 3 {
			pdf.SetXY(464+(float64(k-2)*17.1), 162)
		} else {
			pdf.SetXY(517+(float64(k-4)*17.1), 162)
		}
		pdf.Cell(nil, v)
	}

	for k, v := range signTimes {
		if k <= 1 {
			pdf.SetXY(411+(float64(k)*17.1), 262)
		} else if k <= 3 {
			pdf.SetXY(462+(float64(k-2)*17.1), 262)
		} else {
			pdf.SetXY(516+(float64(k-4)*17.1), 262)
		}
		pdf.Cell(nil, v)
	}

	return nil
}

func (c *GopdfUsecase) CreateContractAMPage8(pdf *gopdf.GoPdf, contractAttorneyVo ContractAttorneyVo, ssnForContract SsnForContract) error {
	c.CreateContractAMPageAttorney(pdf, 8, ssnForContract, contractAttorneyVo)
	return nil
}

func GenContractVetVo(tClient TData, tCase TData) (vo ContractVetVo) {

	vo.FirstName = tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	vo.MiddleName = tClient.CustomFields.TextValueByNameBasic(FieldName_middle_name)
	vo.LastName = tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)

	vo.Street = tCase.CustomFields.TextValueByNameBasic(FieldName_address)
	vo.AptNumber = tCase.CustomFields.TextValueByNameBasic(FieldName_apt_number)
	vo.City = tCase.CustomFields.TextValueByNameBasic(FieldName_city)

	// todo:lgl 需要转化为两个字母
	vo.Province = tCase.CustomFields.TextValueByNameBasic(FieldName_state)
	vo.ZipCode = tCase.CustomFields.TextValueByNameBasic(FieldName_zip_code)
	vo.Ssn = tCase.CustomFields.TextValueByNameBasic(FieldName_ssn)
	vo.Dob = tCase.CustomFields.TextValueByNameBasic(FieldName_dob)
	vo.Email = tCase.CustomFields.TextValueByNameBasic(FieldName_email)

	phone := tCase.CustomFields.TextValueByNameBasic(FieldName_phone)
	phone, _ = USAPhoneHandle(phone)
	vo.Phone = phone
	vo.Branch = tCase.CustomFields.TextValueByNameBasic(FieldName_branch)

	return vo
}

func (c *ContractVetVo) GetAptNumberForContract() string {
	return FormatAptNumber(c.AptNumber)
}

func FormatAptNumber(aptNumber string) string {
	re := regexp.MustCompile(`\D`) // \D 表示非数字
	result := re.ReplaceAllString(aptNumber, "")
	return result
}

type ContractVetVo struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Street     string `json:"street"`
	AptNumber  string `json:"apt_number"`
	City       string `json:"city"`
	Province   string `json:"province"`
	ZipCode    string `json:"zip_code"`
	Ssn        string `json:"ssn"`
	Dob        string `json:"dob"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Branch     string `json:"branch"`
}

func (c *ContractVetVo) GetProvinceForTwo() string {

	if c.Province != "" {
		return GetStateShort(c.Province)
	}
	return ""
}

func (c *ContractVetVo) GetFullName() string {
	return GenFullNameWithMiddleName(c.FirstName, c.MiddleName, c.LastName)
}

func (c *ContractVetVo) GetMiddleNameInitial() string {
	if c.MiddleName == "" {
		return ""
	}
	runes := []rune(c.MiddleName)
	return strings.ToUpper(string(runes[0]))
}

type SsnForContract struct {
	First  []string
	Middle []string
	Last   []string
}

// FormatSsnForContract 411-89-1683 (411)89-1683
func FormatSsnForContract(ssn string) (ssnForContract SsnForContract, err error) {

	ssn = strings.ReplaceAll(ssn, "-", "")
	ssn = strings.ReplaceAll(ssn, "(", "")
	ssn = strings.ReplaceAll(ssn, ")", "")
	ssn = strings.ReplaceAll(ssn, " ", "")
	if len(ssn) != 9 {
		return ssnForContract, errors.New("The format of SSN is incorrect")
	}
	ssns := strings.Split(ssn, "")
	ssnForContract.First = []string{ssns[0], ssns[1], ssns[2]}
	ssnForContract.Middle = []string{ssns[3], ssns[4]}
	ssnForContract.Last = []string{ssns[5], ssns[6], ssns[7], ssns[8]}

	return ssnForContract, nil

}
func (c *GopdfUsecase) CreateContractAMPage7(pdf *gopdf.GoPdf, contractAttorneyVo ContractAttorneyVo,
	ssnForContract SsnForContract, contractVetVo ContractVetVo) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	//firstNames := []string{
	//	"G", "r", "e", "g", "o", "r", "y", "s", "d", "k", "y", "l", "d",
	//}
	//middleLetterName := "K"
	//lastNames := []string{
	//	"G", "a", "e", "g", "o", "r", "y", " ", "d", "k", "y", "l", "3",
	//	//"4", "5", "c", "s", "a", "9",
	//}

	middleLetterName := contractVetVo.GetMiddleNameInitial()
	firstNames := strings.Split(contractVetVo.FirstName, "")
	lastNames := strings.Split(contractVetVo.LastName, "")

	c.CreateContractAMPage7FirstNames(pdf, firstNames)
	if middleLetterName != "" {
		pdf.SetXY(253, 236)
		pdf.Cell(nil, middleLetterName)
	}
	c.CreateContractAMPage7LastNames(pdf, lastNames)
	c.CreateContractAMPage7Ssn(pdf, ssnForContract)
	c.CreateContractAMPage7Dob(pdf, contractVetVo.Dob)
	c.CreateContractAMPage7Branch(pdf, contractVetVo.Branch)
	c.CreateContractAMPage7Street(pdf, contractVetVo.Street)
	c.CreateContractAMPage7Province(pdf, contractVetVo.GetProvinceForTwo())
	c.CreateContractAMPage7City(pdf, contractVetVo.City)
	c.CreateContractAMPage7ZipCode(pdf, contractVetVo.ZipCode)
	c.CreateContractAMPage7Phone(pdf, contractVetVo.Phone)
	c.CreateContractAMPage7Email(pdf, contractVetVo.Email)
	c.CreateContractAMPageAttorney(pdf, 7, ssnForContract, contractAttorneyVo)

	aptNumber := contractVetVo.GetAptNumberForContract()
	if aptNumber != "" {
		aptNumbers := strings.Split(aptNumber, "")
		for k, v := range aptNumbers {
			if k == 5 {
				break
			}
			pdf.SetXY(102+float64(k)*17.1, 365)
			pdf.Cell(nil, v)
		}
	}

	pdf.SetXY(183, 385)
	pdf.Cell(nil, "U")
	pdf.SetXY(183+16, 385)
	pdf.Cell(nil, "S")

	return nil
}

type ContractAttorneyVo struct {
	Street              string `json:"street"`
	City                string `json:"city"`
	Province            string `json:"province"`
	ZipCode             string `json:"zip_code"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	AccreditationDate   string `json:"accreditation_date"`
	AccreditationNumber string `json:"accreditation_number"`
	Email               string `json:"email"`
}

func (c *ContractAttorneyVo) GetAccreditationDate() string {
	if c.AccreditationDate != "" {
		a, err := time.Parse(time.DateOnly, c.AccreditationDate)
		if err != nil {
			return ""
		}
		return a.Format(configs.TimeFormatDate2)
	}
	return ""
}

func (c *ContractAttorneyVo) FullName() string {
	fullName := c.FirstName
	if c.LastName != "" {
		fullName += " " + c.LastName
	}
	return fullName
}

func (c *GopdfUsecase) CreateContractAMPageAttorney(pdf *gopdf.GoPdf, page int, ssnForContract SsnForContract, contractAttorneyVo ContractAttorneyVo) {

	if page == 7 {
		firstNames := strings.Split(contractAttorneyVo.FirstName, "")
		for k, v := range firstNames {
			if k >= 12 {
				break
			}
			pdf.SetXY(39+(float64(k)*16.8), 680)
			pdf.Cell(nil, v)
		}
		lastNames := strings.Split(contractAttorneyVo.LastName, "")
		for k, v := range lastNames {
			if k >= 18 {
				break
			}
			pdf.SetXY(280+(float64(k)*16.8), 680)
			pdf.Cell(nil, v)
		}

		ttfName := "Arial"
		letter := "X"
		pdf.SetFont(ttfName, "", 9)
		pdf.SetXY(38, 713.5)
		pdf.Cell(nil, letter)
		pdf.SetFont(ttfName, "", 10)
	} else if page == 8 {
		street := strings.TrimSpace(contractAttorneyVo.Street)
		arr := strings.Split(street, ",")
		var runes []rune
		for _, v := range arr {
			v = strings.TrimSpace(v)
			if len(v) > 0 {
				t := []rune(v)
				if len(runes) > 0 {
					t = append([]rune{','}, t...)
				}
				runes = append(runes, t...)
			}
		}

		for k, v := range runes {
			if k >= 30 {
				break
			}
			temp := string(v)
			pdf.SetXY(69+(float64(k)*17.1), 60)
			pdf.Cell(nil, temp)
		}

		//accreditationNumbers := strings.Split(contractAttorneyVo.AccreditationNumber, "")
		//for k, v := range accreditationNumbers {
		//	pdf.SetXY(102+(float64(k)*17.1), 79)
		//	pdf.Cell(nil, v)
		//}
		cities := strings.Split(contractAttorneyVo.City, "")
		for k, v := range cities {
			pdf.SetXY(240+(float64(k)*17.1), 79)
			pdf.Cell(nil, v)
		}
		provinces := strings.Split(contractAttorneyVo.Province, "")
		for k, v := range provinces {
			pdf.SetXY(102+(float64(k)*17.1), 98)
			pdf.Cell(nil, v)
		}

		zipCodes := strings.Split(contractAttorneyVo.ZipCode, "")
		for k, v := range zipCodes {
			pdf.SetXY(319.5+(float64(k)*17.1), 98)
			pdf.Cell(nil, v)
		}

		pdf.SetXY(279, 148)
		pdf.Cell(nil, contractAttorneyVo.Email)

		ttfName := "Arial"
		letter := "X"
		pdf.SetFont(ttfName, "", 9)
		pdf.SetXY(46.5, 369)
		pdf.Cell(nil, letter)
		pdf.SetXY(46.5, 442.5)
		pdf.Cell(nil, letter)
		pdf.SetXY(36.5, 542)
		pdf.Cell(nil, letter)
		pdf.SetFont(ttfName, "", 10)

		pdf.SetXY(63, 413)
		pdf.Cell(nil, "Augustus Miles LLC")

		pdf.SetXY(63, 484)
		pdf.Cell(nil, "Edward Bunting Jr., Donald Pratko, Victoria Enriquez")
	}
}

func (c *GopdfUsecase) CreateContractAMPage7Email(pdf *gopdf.GoPdf, email string) {
	pdf.SetXY(287, 434)
	pdf.Cell(nil, email)
}

func (c *GopdfUsecase) CreateContractAMPage7Phone(pdf *gopdf.GoPdf, phone string) {

	phone = strings.TrimSpace(phone)
	runes := []rune(phone)
	for k, v := range runes {
		if k >= 10 {
			break
		}
		temp := string(v)
		if k >= 6 {
			pdf.SetXY(177.5+(float64(k-6)*17.25), 418)
		} else if k >= 3 {
			pdf.SetXY(110.5+(float64(k-3)*17.25), 418)
		} else {
			pdf.SetXY(40+(float64(k)*17.25), 418)
		}
		pdf.Cell(nil, temp)
	}
}

func (c *GopdfUsecase) CreateContractAMPage7ZipCode(pdf *gopdf.GoPdf, zipCode string) {

	zipCode = strings.TrimSpace(zipCode)
	runes := []rune(zipCode)
	for k, v := range runes {
		if k >= 9 {
			break
		}
		temp := string(v)
		if k >= 5 {
			pdf.SetXY(418.5+(float64(k-5)*17.25), 385)
		} else {
			pdf.SetXY(317.5+(float64(k)*17.25), 385)
		}
		pdf.Cell(nil, temp)
	}
}

func (c *GopdfUsecase) CreateContractAMPage7Province(pdf *gopdf.GoPdf, province string) {

	province = strings.TrimSpace(province)
	runes := []rune(province)
	for k, v := range runes {
		if k >= 2 {
			break
		}
		temp := string(v)
		pdf.SetXY(101+(float64(k)*17.1), 384)
		pdf.Cell(nil, temp)
	}
}

func (c *GopdfUsecase) CreateContractAMPage7City(pdf *gopdf.GoPdf, city string) {

	city = strings.TrimSpace(city)
	runes := []rune(city)
	for k, v := range runes {
		if k >= 18 {
			break
		}
		temp := string(v)
		pdf.SetXY(235+(float64(k)*17.25), 366)
		pdf.Cell(nil, temp)
	}
}

var StreetMapping = map[string]string{
	"Floor":     "FL",
	"Apartment": "Apt",
	"Unit":      "#",
	"Suite":     "STE",
	"South":     "S",
	"North":     "N",
	"East":      "E",
	"West":      "W",
	"Road":      "Rd",
	"Street":    "St",
	"Place":     "Pl",
	"Court":     "Ct",
	"Boulevard": "Blvd",
	"Avenue":    "Ave",
	"Circle":    "Cir",
	"Drive":     "Dr",
	"Highway":   "Hwy",
	"Lane":      "Ln",
	"Square":    "Sq",
	"Terrace":   "Ter",
	"Alley":     "Aly",
	"Center":    "Ctr",
	"Creek":     "Crk",
	"Crossing":  "Xing",
	"Extension": "Ext",
	"Freeway":   "Fwy",
	"Grove":     "Grv",
	"Parkway":   "Pkwy",
	"Point":     "Pt",
	"Trail":     "Trl",
}

func StreetShortDo(word string) string {
	if word == "" {
		return ""
	}
	for k, v := range StreetMapping {
		if strings.ToLower(k) == strings.ToLower(word) {
			return v
		}
	}
	for _, v := range StreetMapping {
		if strings.ToLower(v+".") == strings.ToLower(word) {
			return v
		}
	}

	return word
}

func StreetShort(str string) string {
	arr := strings.Split(str, ",")
	var r []string
	lib.DPrintln(arr)
	for _, v := range arr {
		v = strings.TrimSpace(v)
		//if v != "" {
		newV := StreetShortDo(v)
		r = append(r, newV)
		//}
	}
	return strings.Join(r, ",")
}

func StreetClean(street string) []rune {

	street = strings.Trim(street, ",")
	ar := strings.Split(street, " ")
	street = ""
	for _, v := range ar {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		v = StreetShort(v)
		if street == "" {
			street = v
		} else {
			street += " " + v
		}
	}
	re := regexp.MustCompile(`#\s+`)
	street = re.ReplaceAllString(street, "#")

	arr := strings.Split(street, ",")
	var runes []rune
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			t := []rune(v)
			if len(runes) > 0 {
				t = append([]rune{','}, t...)
			}
			runes = append(runes, t...)
		}
	}
	return runes
}

func (c *GopdfUsecase) CreateContractAMPage7Street(pdf *gopdf.GoPdf, Street string) {
	//Street = "321, Kensington Dr 321, Kensington Dr 321, Kensington Dr"

	runes := StreetClean(Street)
	for k, v := range runes {
		if k >= 30 {
			break
		}
		temp := string(v)
		pdf.SetXY(69+(float64(k)*17.25), 342.5)
		pdf.Cell(nil, temp)
	}

}
func (c *GopdfUsecase) CreateContractAMPage7Branch(pdf *gopdf.GoPdf, branch string) {

	// Branch of Service:
	//Public Health Service in Base - USPHS,
	//National Oceanic and Atmospheric Administration in Base - NOAA,
	//Air National Guard in base - Air Force,
	//Army National Guard in Base - Army
	ttfName := "Arial"
	//branch = "Army National Guard"
	letter := "X" //✗
	pdf.SetFont(ttfName, "", 9)
	if branch == "Army National Guard" { // Army
		pdf.SetXY(206.5, 293)
		pdf.Cell(nil, letter)
	}
	if branch == "Navy" { // Navy
		pdf.SetXY(257, 293)
		pdf.Cell(nil, letter)
	}
	if branch == "Air National Guard" { // Air Force
		pdf.SetXY(302, 293)
		pdf.Cell(nil, letter)
	}
	if branch == "Marine Corps" { // Marine Corps
		pdf.SetXY(361, 293)
		pdf.Cell(nil, letter)
	}
	if branch == "Coast Guard" { // Coast Guard
		pdf.SetXY(434.5, 293)
		pdf.Cell(nil, letter)
	}
	if branch == "Space Force" { // Space Force
		pdf.SetXY(206.5, 310)
		pdf.Cell(nil, letter)
	}
	if branch == "National Oceanic and Atmospheric Administration" { // NOAA
		pdf.SetXY(286, 310)
		pdf.Cell(nil, letter)
	}
	if branch == "Public Health Service" { // USPHS
		pdf.SetXY(333, 310)
		pdf.Cell(nil, letter)
	}
	pdf.SetFont(ttfName, "", 10)
}

func (c *GopdfUsecase) CreateContractAMPage7Dob(pdf *gopdf.GoPdf, dob string) {
	if dob == "" {
		return
	}
	dobs := strings.Split(dob, "-")
	if len(dobs) != 3 {
		return
	}
	mm := strings.Split(dobs[1], "")
	dd := strings.Split(dobs[2], "")
	year := strings.Split(dobs[0], "")
	for k, v := range mm {
		pdf.SetXY(float64(412+(k*17)), 267)
		pdf.Cell(nil, v)
	}
	for k, v := range dd {
		pdf.SetXY(float64(464+(k*17)), 267)
		pdf.Cell(nil, v)
	}
	for k, v := range year {
		pdf.SetXY(float64(518+(k*17)), 267)
		pdf.Cell(nil, v)
	}

}

func (c *GopdfUsecase) CreateContractAMPage7Ssn(pdf *gopdf.GoPdf, ssnForContract SsnForContract) {

	for k, v := range ssnForContract.First {
		pdf.SetXY(float64(44+(k*17)), 267.5)
		pdf.Cell(nil, v)
	}
	for k, v := range ssnForContract.Middle {
		pdf.SetXY(float64(117+(k*17)), 267.5)
		pdf.Cell(nil, v)
	}
	for k, v := range ssnForContract.Last {
		pdf.SetXY(float64(172+(k*17)), 267.5)
		pdf.Cell(nil, v)
	}

}
func (c *GopdfUsecase) CreateContractAMPage7LastNames(pdf *gopdf.GoPdf, lastNames []string) {
	for k, v := range lastNames {
		if k > 17 {
			break
		}
		pdf.SetXY(float64(278+(k*17)), 236)
		pdf.Cell(nil, v)
	}
}

func (c *GopdfUsecase) CreateContractAMPage7FirstNames(pdf *gopdf.GoPdf, firstNames []string) {
	for k, v := range firstNames {
		if k > 11 {
			break
		}
		if k <= 1 {
			pdf.SetXY(float64(42+(k*19)), 236)
		} else if k >= 6 {
			pdf.SetXY(float64(42+(k*17)), 236)
		} else {
			pdf.SetXY(float64(43+(k*17)), 236)
		}
		pdf.Cell(nil, v)
	}
}

func (c *GopdfUsecase) CreateContractAMPage1(pdf *gopdf.GoPdf,

	contractTime time.Time,
	contractVetVo ContractVetVo,
	contractAttorneyVo ContractAttorneyVo,
) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	pdf.SetXY(311, 117)
	pdf.Cell(nil, lib.DayWithSuffix(contractTime))

	pdf.SetXY(373, 117)
	pdf.Cell(nil, contractTime.Format("January"))

	pdf.SetXY(443, 117)
	pdf.Cell(nil, InterfaceToString(contractTime.Year()%100))

	pdf.SetXY(82, 131)
	pdf.Cell(nil, contractVetVo.GetFullName())

	pdf.SetXY(315, 131)
	pdf.Cell(nil, contractAttorneyVo.FullName())

	pdf.SetXY(204, 145)
	pdf.Cell(nil, contractAttorneyVo.GetAccreditationDate())

	pdf.SetXY(389, 145)
	pdf.Cell(nil, contractAttorneyVo.AccreditationNumber)

	return nil
}

func (c *GopdfUsecase) CreateContractAmPage6(pdf *gopdf.GoPdf) error {
	ttfName := "Arial"
	pdf.SetTextColor(255, 255, 255)
	err := pdf.SetFont(ttfName, "", 36)
	if err != nil {
		return err
	}
	pdf.SetXY(171, 550)
	pdf.Cell(nil, "[[s|1]]")

	err = pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return err
	}
	pdf.SetXY(115, 598)
	pdf.Cell(nil, "[[d|1              ]]")

	err = pdf.SetFont(ttfName, "", 36)
	if err != nil {
		return err
	}
	pdf.SetXY(198, 634)
	pdf.Cell(nil, "[[s|2]]")

	err = pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return err
	}
	pdf.SetXY(115, 682)
	pdf.Cell(nil, "[[d|2              ]]")
	pdf.SetTextColor(0, 0, 0)
	return nil
}

func (c *GopdfUsecase) CreateContract(createContractVo CreateContractVo, contractIndex int) (pdfFile []byte, err error) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})

	for i := 1; i <= 8; i++ {
		pdf.AddPage()
		ttfName := "Arial"
		err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/Arial.ttf")
		if err != nil {
			return nil, err
		}
		err = pdf.SetFont(ttfName, "", 10)
		if err != nil {
			return nil, err
		}
		pageJpgPath := c.ResourceUsecase.ResPath() + "/Contracts/AgreementForConsultingServices-000" + InterfaceToString(i) + ".jpg"
		if i == 2 {
			pageJpgPath = c.ResourceUsecase.ResPath() + "/Contracts/" + InterfaceToString(contractIndex) + "/AgreementForConsultingServices-000" + InterfaceToString(i) + ".jpg"
		}
		//image bytes
		b, err := os.ReadFile(pageJpgPath)
		if err != nil {
			return nil, err
		}
		imgH1, err := gopdf.ImageHolderByBytes(b)
		if err != nil {
			return nil, err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 612, H: 792}); err != nil {
			return nil, err
		}
		//// Color the page
		//pdf.SetLineWidth(0.1)
		//pdf.SetFillColor(124, 252, 0) //setup fill color
		//pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
		//pdf.SetFillColor(0, 0, 0)

		if i == 1 {
			err = c.CreateContractPage1(&pdf, createContractVo.ClientName)
			if err != nil {
				return nil, err
			}
		} else if i == 2 {
			err = c.CreateContractPage2(&pdf)
			if err != nil {
				return nil, err
			}
		} else if i == 6 {
			err = c.CreateContractPage6(&pdf, createContractVo.ClientName, createContractVo.ClientEmail, createContractVo.VsName, createContractVo.VsEmail)
			if err != nil {
				return nil, err
			}
		} else if i == 8 {
			err = c.CreateContractPage8(&pdf)
			if err != nil {
				return nil, err
			}
		}
	}

	return pdf.GetBytesPdf(), nil
}

func (c *GopdfUsecase) CreateContractPage1(pdf *gopdf.GoPdf, clientName string) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	pdf.SetXY(296, 170)
	pdf.Cell(nil, clientName)

	err = pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return err
	}
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(76, 168)
	pdf.Cell(nil, "[[d|1              ]]")
	return nil
}

func (c *GopdfUsecase) CreateContractPage2(pdf *gopdf.GoPdf) error {
	ttfName := "Arial"
	//pdf.SetTextColor(0, 0, 0)
	pdf.SetTextColor(255, 255, 255)
	err := pdf.SetFont(ttfName, "", 36)
	if err != nil {
		return err
	}
	pdf.SetXY(398, 266)
	pdf.Cell(nil, "[[i|1]]")
	return nil
}

func (c *GopdfUsecase) CreateContractPage6(pdf *gopdf.GoPdf,
	clientName string, clientEmail string,
	vsName string, vsEmail string) error {
	ttfName := "Arial"
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)

	pdf.SetXY(141, 133)
	pdf.Cell(nil, clientName)

	pdf.SetXY(362, 180)
	pdf.Cell(nil, clientEmail)

	pdf.SetXY(207, 228)
	pdf.Cell(nil, vsName)

	pdf.SetXY(362, 275)
	pdf.Cell(nil, vsEmail)

	err = pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return err
	}
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(358, 131)
	pdf.Cell(nil, "[[d|1              ]]")
	pdf.SetXY(358, 226)
	pdf.Cell(nil, "[[d|2              ]]")

	err = pdf.SetFont(ttfName, "", 47)
	if err != nil {
		return err
	}
	pdf.SetXY(151+5, 147)
	pdf.Cell(nil, "[[s|1]]")

	pdf.SetXY(148+5, 241)
	pdf.Cell(nil, "[[s|2]]")

	return nil
}

func (c *GopdfUsecase) CreateContractPage8(pdf *gopdf.GoPdf) error {
	ttfName := "Arial"

	err := pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return err
	}
	//pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(253, 114)
	pdf.Cell(nil, "[[d|1              ]]")
	err = pdf.SetFont(ttfName, "", 40)
	if err != nil {
		return err
	}
	pdf.SetXY(85, 88)
	pdf.Cell(nil, "[[s|1]]")

	return nil
}
func (c *GopdfUsecase) CreateMedicalTeamForm(createMedicalTeamFormVo CreateMedicalTeamFormVo) (pdfFile []byte, err error) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for i := 1; i <= 4; i++ {
		pdf.AddPage()
		ttfName := "Arial"
		err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/Arial.ttf")
		if err != nil {
			return nil, err
		}
		err = pdf.SetFont(ttfName, "", 10)
		if err != nil {
			return nil, err
		}
		//image bytes
		b, err := os.ReadFile(c.ResourceUsecase.ResPath() + "/MT/MedicalTeamFormsPage-000" + InterfaceToString(i) + ".jpg")
		if err != nil {
			return nil, err
		}
		imgH1, err := gopdf.ImageHolderByBytes(b)
		if err != nil {
			return nil, err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 595, H: 842}); err != nil {
			return nil, err
		}
		//// Color the page
		//pdf.SetLineWidth(0.1)
		//pdf.SetFillColor(124, 252, 0) //setup fill color
		//pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
		//pdf.SetFillColor(0, 0, 0)

		if i == 1 {
			err = c.MedicalTeamFormsPage1(&pdf, createMedicalTeamFormVo.ClientName)
			if err != nil {
				return nil, err
			}
		} else if i == 2 {

			err = c.MedicalTeamFormsPage2(&pdf, createMedicalTeamFormVo.Address, createMedicalTeamFormVo.Location)
			if err != nil {
				return nil, err
			}
		} else if i == 3 {
			err = c.MedicalTeamFormsPage3(&pdf, createMedicalTeamFormVo.ClientName)
			if err != nil {
				return nil, err
			}
		} else if i == 4 {
			err = c.MedicalTeamFormsPage4(&pdf,
				createMedicalTeamFormVo.ClientName,
				createMedicalTeamFormVo.Dob,
				createMedicalTeamFormVo.Ssn,
				createMedicalTeamFormVo.Phone,
				createMedicalTeamFormVo.Email,
				createMedicalTeamFormVo.Address,
				createMedicalTeamFormVo.Location,
				createMedicalTeamFormVo.Itf,
				createMedicalTeamFormVo.PrivateExamsNeededS)
			if err != nil {
				return nil, err
			}
		}
	}

	//pdf.WritePdf("example.pdf")
	return pdf.GetBytesPdf(), nil
}

func (c *GopdfUsecase) Test() error {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for i := 1; i <= 4; i++ {
		pdf.AddPage()
		ttfName := "Arial"
		err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/Arial.ttf")
		if err != nil {
			return err
		}
		err = pdf.SetFont(ttfName, "", 10)
		if err != nil {
			return err
		}
		//image bytes
		b, err := os.ReadFile(c.ResourceUsecase.ResPath() + "/MT/MedicalTeamFormsPage-000" + InterfaceToString(i) + ".jpg")
		if err != nil {
			return err
		}
		imgH1, err := gopdf.ImageHolderByBytes(b)
		if err != nil {
			return err
		}
		if err := pdf.ImageByHolder(imgH1, 0, 0, &gopdf.Rect{W: 595, H: 842}); err != nil {
			return err
		}
		//// Color the page
		//pdf.SetLineWidth(0.1)
		//pdf.SetFillColor(124, 252, 0) //setup fill color
		//pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
		//pdf.SetFillColor(0, 0, 0)

		if i == 1 {
			err = c.MedicalTeamFormsPage1(&pdf, "Gary Liao")
			if err != nil {
				return err
			}
		} else if i == 2 {

			err = c.MedicalTeamFormsPage2(&pdf, "3837 Battle Creek Road", "Chula Vista, California 910910")
			if err != nil {
				return err
			}
		} else if i == 3 {
			err = c.MedicalTeamFormsPage3(&pdf, "Gary Liao")
			if err != nil {
				return err
			}
		} else if i == 4 {
			err = c.MedicalTeamFormsPage4(&pdf,
				"Gary Liao",
				"2023-12-11",
				"3344-444-444",
				"8837-333-33-33",
				"aaaoh@qq.com",
				"3837 Battle Creek Road",
				"Chula Vista, California 910910",
				"2020-11-11",
				[]string{
					"1fdsafas",
					"2sss",
					"2sss",
					"2sss",
					"2sss",
					"2sss",
					"2sss",
					"8sss",
					"9sss",
					"10sss",
					"10sss",
					//"10sss",
					//"10sss",
					//"10sss",
					//"15sss",
				})
			if err != nil {
				return err
			}
		}
	}

	pdf.WritePdf("example.pdf")
	return nil
}

func (c *GopdfUsecase) MedicalTeamFormsPage1(pdf *gopdf.GoPdf, clientName string) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}

	pdf.SetXY(320, 140)
	pdf.Cell(nil, clientName)
	return nil
}

func (c *GopdfUsecase) MedicalTeamFormsPage2(pdf *gopdf.GoPdf, address string, location string) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}

	pdf.SetXY(90, 252)
	pdf.Cell(nil, address)

	pdf.SetXY(90, 265)
	pdf.Cell(nil, location)
	return nil
}

func (c *GopdfUsecase) MedicalTeamFormsPage3(pdf *gopdf.GoPdf, fullName string) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	pdf.SetXY(289, 282)
	pdf.Cell(nil, fullName)

	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(337, 294)
	pdf.Cell(nil, "[[d|1          ]]")

	err = pdf.SetFont(ttfName, "", 48)
	if err != nil {
		return err
	}
	pdf.SetXY(340, 233)

	pdf.Cell(nil, "[[s|1]]")

	return nil
}

func (c *GopdfUsecase) MedicalTeamFormsPage4(pdf *gopdf.GoPdf,
	clientName string,
	dob string,
	ssn string,
	phone string,
	email string,
	address string,
	location string,
	itf string,
	privateExamsNeededS []string) error {

	ttfName := "Arial"
	pdf.SetTextColor(0, 0, 0)
	err := pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	var gap, initY float64
	gap = 12
	initY = 82
	pdf.SetXY(170, initY)
	pdf.Cell(nil, clientName)
	pdf.SetXY(170, initY+gap)
	pdf.Cell(nil, dob)
	pdf.SetXY(170, initY+2*gap)
	pdf.Cell(nil, ssn)
	pdf.SetXY(170, initY+3*gap)
	pdf.Cell(nil, phone)
	pdf.SetXY(170, initY+4*gap)
	pdf.Cell(nil, email)
	pdf.SetXY(170, initY+5*gap)
	pdf.Cell(nil, address)
	pdf.SetXY(170, initY+6*gap)
	pdf.Cell(nil, location)
	pdf.SetXY(170, initY+7*gap)
	pdf.Cell(nil, itf)

	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(314, 622)
	pdf.Cell(nil, "[[d|1          ]]")

	err = pdf.SetFont(ttfName, "", 48)
	if err != nil {
		return err
	}
	pdf.SetXY(90, 588)
	pdf.Cell(nil, "[[s|1]]")

	pdf.SetTextColor(0, 0, 0)
	err = pdf.SetFont(ttfName, "", 10)
	if err != nil {
		return err
	}
	var privateExamsNeededY, privateExamsNeededGap float64
	privateExamsNeededY = 699
	privateExamsNeededGap = 10
	pdf.SetXY(30, privateExamsNeededY)

	for _, v := range privateExamsNeededS {

		texts, err := pdf.SplitTextWithWordWrap(v, 535)
		if err != nil {
			return err
		}
		//privateExamsNeededY = 5 + pdf.GetY()
		pdf.SetY(pdf.GetY() + 3)
		for _, text := range texts {
			//pdf.GetY()
			//lib.DPrintln("pdf.GetY():", pdf.GetY())
			privateExamsNeededY = privateExamsNeededGap + pdf.GetY()
			if privateExamsNeededY >= 820 {
				privateExamsNeededY = 20
				pdf.AddPage()
			}
			pdf.SetXY(30, privateExamsNeededY)
			pdf.Cell(nil, text)

			//pdf.MultiCell(&gopdf.Rect{
			//	H: 120,
			//	W: 300,
			//}, v)
		}
	}

	return nil
}

// CreatePersonalStatementsPDFForAiV1 statementConditionId=0时，全部； 非0时，指定
func (c *GopdfUsecase) CreatePersonalStatementsPDFForAiV1(dealName string, statementDetail StatementDetail, statementConditionId int32) (io.Reader, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})
	var lineWidth float64 = 2

	// Add first page
	pdf.AddPage()
	ttfName := "Arial"
	err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/Arial.ttf")
	// ttfName := "Times New Roman"
	// err := pdf.AddTTFFont(ttfName, c.ResourceUsecase.ResPath()+"/ttf/TimesNewRomanPSMT.ttf")
	if err != nil {
		return nil, err
	}

	err = pdf.SetFont(ttfName, "", 12)
	if err != nil {
		return nil, err
	}

	// Define page boundaries with consistent margins
	leftMargin := 50.0
	rightMargin := 562.0 // 612 - 50 (page width minus margins)
	topMargin := 50.0
	bottomMargin := 742.0 // 792 - 50 (page height minus margins)
	lineHeight := 15.0

	// Helper function to check if we need a new page
	checkNewPage := func(currentY float64, requiredSpace float64) float64 {
		if currentY+requiredSpace > bottomMargin {
			pdf.AddPage()
			return topMargin
		}
		return currentY
	}

	// Helper function to wrap text and write to PDF
	wrapAndWriteText := func(text string, x float64, y float64, maxWidth float64) (float64, error) {
		if text == "" {
			return y, nil
		}

		// Split text into words
		words := strings.Fields(text)
		if len(words) == 0 {
			return y, nil
		}

		currentLine := ""
		currentY := y

		for _, word := range words {
			testLine := currentLine
			if testLine == "" {
				testLine = word
			} else {
				testLine += " " + word
			}

			// Get text width to check if it fits
			textWidth, err := pdf.MeasureTextWidth(testLine)
			if err != nil {
				return y, err
			}

			if textWidth > maxWidth {
				// Current line doesn't fit, write what we have
				if currentLine != "" {
					currentY = checkNewPage(currentY, lineHeight)
					pdf.SetXY(x, currentY)
					pdf.Cell(nil, currentLine)
					currentY += lineHeight
				}
				currentLine = word
			} else {
				currentLine = testLine
			}
		}

		// Write the last line
		if currentLine != "" {
			currentY = checkNewPage(currentY, lineHeight)
			pdf.SetXY(x, currentY)
			pdf.Cell(nil, currentLine)
			currentY += lineHeight
		}

		return currentY, nil
	}

	// Add deal name as title
	var yPosition float64 = topMargin
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(leftMargin, yPosition)
	pdf.Cell(nil, dealName)
	yPosition += lineHeight * 2

	// Add separator line
	yPosition += 5
	pdf.SetXY(leftMargin, yPosition)
	// Draw line using Line method instead of underscores for better compatibility
	lineEndX := rightMargin - 20
	// Set line width to make it thicker
	pdf.SetLineWidth(lineWidth)
	pdf.Line(leftMargin, yPosition+6, lineEndX, yPosition+6)
	// Reset line width to default
	pdf.SetLineWidth(0.1)
	yPosition += 10
	yPosition = checkNewPage(yPosition, lineHeight*2)
	yPosition += lineHeight

	// Add base information
	statementBaseInfoList := []StatementBaseInfo{
		{Label: "Full Name", Value: statementDetail.BaseInfo.FullName},
		{Label: "Unique ID", Value: InterfaceToString(statementDetail.CaseId)},
		{Label: "Branch of Service", Value: statementDetail.BaseInfo.BranchOfService},
		{Label: "Years of Service", Value: statementDetail.BaseInfo.YearsOfService},
		{Label: "Retired from service", Value: statementDetail.BaseInfo.RetiredFromService},
		{Label: "Marital Status", Value: statementDetail.BaseInfo.MaritalStatus},
		{Label: "Children", Value: statementDetail.BaseInfo.Children},
		{Label: "Occupation in service", Value: statementDetail.BaseInfo.OccupationInService},
	}

	for _, v := range statementBaseInfoList {
		pdf.SetCharSpacing(0)
		text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		var err error
		yPosition, err = wrapAndWriteText(text, leftMargin, yPosition, rightMargin-leftMargin)
		if err != nil {
			return nil, err
		}
		// 恢复正常字符间距
		pdf.SetCharSpacing(0)
	}

	// Add separator line
	yPosition += 10
	yPosition = checkNewPage(yPosition, lineHeight*2)
	pdf.SetXY(leftMargin, yPosition)
	// Draw line using Line method instead of underscores for better compatibility
	lineEndX = rightMargin - 20
	// Set line width to make it thicker
	pdf.SetLineWidth(lineWidth)
	pdf.Line(leftMargin, yPosition+6, lineEndX, yPosition+6)
	// Reset line width to default
	pdf.SetLineWidth(0.1)
	yPosition += lineHeight * 2

	// Add statements
	for k, v := range statementDetail.Statements {
		if v.IsEmptyResult() {
			continue
		}

		if statementConditionId != 0 {
			if v.StatementCondition.StatementConditionId != statementConditionId {
				continue
			}
		}

		// Add disability/condition name with proper wrapping
		conditionText := "Name of Disability/Condition: " + v.StatementCondition.ConditionValue
		var err error
		yPosition, err = wrapAndWriteText(conditionText, leftMargin, yPosition, rightMargin-leftMargin)
		if err != nil {
			return nil, err
		}

		// Add current treatment facility and medication first
		for _, v1 := range v.Rows {
			if v1.SectionType == Statemt_Section_CurrentTreatmentFacility ||
				v1.SectionType == Statemt_Section_CurrentMedication {
				sectionText := GetSectionTitleFromSectionType(v1.SectionType) + ": " + v1.Body

				yPosition, err = wrapAndWriteText(sectionText, leftMargin, yPosition, rightMargin-leftMargin)
				if err != nil {
					return nil, err
				}
			}
		}

		yPosition += lineHeight

		// Add other sections
		for _, v1 := range v.Rows {
			if v1.SectionType != Statemt_Section_CurrentTreatmentFacility &&
				v1.SectionType != Statemt_Section_CurrentMedication {

				if v1.SectionType == Statemt_Section_SpecialNotes || v1.SectionType == Statemt_Section_IntroductionParagraph {
					if v1.Body != "" {
						lines := strings.Split(v1.Body, "\n")
						for _, line := range lines {
							line = strings.TrimSpace(line)
							if line == "" {
								continue
							}
							// Check if we need a new page for this line plus spacing
							yPosition = checkNewPage(yPosition, lineHeight*2)
							yPosition, err = wrapAndWriteText(line, leftMargin, yPosition, rightMargin-leftMargin)
							if err != nil {
								return nil, err
							}
							yPosition += lineHeight // Add spacing between lines
						}
					}
				} else {
					if v1.Body == "" {
						continue
					}

					// Add section title
					yPosition += lineHeight
					yPosition, err = wrapAndWriteText(v1.Title+":", leftMargin, yPosition, rightMargin-leftMargin)
					yPosition += lineHeight
					if err != nil {
						return nil, err
					}

					// Add section body with line-by-line processing
					if v1.Body != "" {
						lines := strings.Split(v1.Body, "\n")
						for lineKey, line := range lines {
							line = strings.TrimSpace(line)
							if line == "" {
								continue
							}
							yPosition, err = wrapAndWriteText(line, leftMargin, yPosition, rightMargin-leftMargin)
							if err != nil {
								return nil, err
							}

							// Add extra line break for non-request sections or non-last lines
							if v1.SectionType != Statemt_Section_Request || lineKey != len(lines)-1 {
								yPosition += lineHeight / 2
							}
						}
					}
				}
			}
		}

		// Add separator between statements (except for the last one)
		if k != len(statementDetail.Statements)-1 && statementConditionId == 0 {
			yPosition += lineHeight
			yPosition = checkNewPage(yPosition, lineHeight*3)
			pdf.SetXY(leftMargin, yPosition)
			// Draw line using Line method instead of underscores for better compatibility
			lineEndX := rightMargin - 20
			// Set line width to make it thicker
			pdf.SetLineWidth(lineWidth)
			pdf.Line(leftMargin, yPosition+6, lineEndX, yPosition+6)
			// Reset line width to default
			pdf.SetLineWidth(0.1)
			yPosition += lineHeight * 2
		}
	}

	// Return PDF as bytes reader
	pdfBytes := pdf.GetBytesPdf()

	return bytes.NewReader(pdfBytes), nil
}
