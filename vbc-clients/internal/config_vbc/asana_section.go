package config_vbc

import (
	"vbc/configs"
	"vbc/lib"
)

type AsanaSections lib.TypeMap

const (
	AsanaSection_GETTING_STARTED_EMAIL = "GETTING STARTED EMAIL"
)

func (c AsanaSections) GetSectionGidByName(name string) string {
	r := lib.ToTypeList(lib.TypeMap(c).Get("data"))
	for _, v := range r {
		if v.GetString("name") == name {
			return v.GetString("gid")
		}
	}
	return ""
}

func GetAsanaSections() AsanaSections {
	if configs.IsProd() {
		return AsanaSections(lib.ToTypeMapByString(ProdSections))
	} else {
		return AsanaSections(lib.ToTypeMapByString(TestSections))
	}
}

var ProdSections = `{
  "data": [
    {
      "gid": "1206472580135543",
      "name": "INCOMING REQUEST",
      "resource_type": "section"
    },
    {
      "gid": "1206472580135544",
      "name": "FEE SCHEDULE AND CONTRACT",
      "resource_type": "section"
    },
    {
      "gid": "1206472580135555",
      "name": "GETTING STARTED EMAIL",
      "resource_type": "section"
    },
    {
      "gid": "1206472580135566",
      "name": "UPDATE CLIENT INTAKE INFO",
      "resource_type": "section"
    },
    {
      "gid": "1206472580135567",
      "name": "AWAITING C-FILE",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162482",
      "name": "RECORD REVIEW",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162484",
      "name": "SCHEDULE CALL",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162486",
      "name": "STATEMENT DRAFTS",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162488",
      "name": "STATEMENTS FINALIZED",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162490",
      "name": "CURRENT TREATMENT",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162492",
      "name": "MINI-DBQS",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162494",
      "name": "NEXUS LETTERS",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162496",
      "name": "MEDICAL TEAM",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162498",
      "name": "FILE CLAIMS (NEW OR SUPPLEMENTAL)",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162500",
      "name": "VERIFY EVIDENCE RECEIVED",
      "resource_type": "section"
    },
    {
      "gid": "1206479748162502",
      "name": "AWAITING DECISION",
      "resource_type": "section"
    },
    {
      "gid": "1206479748212992",
      "name": "AWAITING PAYMENT (UPDATE RATING ON BOARD, SEND INVOICE AND CLOTHING EMAIL)",
      "resource_type": "section"
    },
    {
      "gid": "1206479748212994",
      "name": "COMPLETED",
      "resource_type": "section"
    }
  ]
}`

var TestSections = `{
  "data": [
    {
      "gid": "1205963735008551",
      "name": "INCOMING REQUEST",
      "resource_type": "section"
    },
    {
      "gid": "1206176340563332",
      "name": "FEE SCHEDULE AND CONTRACT",
      "resource_type": "section"
    },
    {
      "gid": "1205962480830599",
      "name": "GETTING STARTED EMAIL",
      "resource_type": "section"
    },
    {
      "gid": "1205962480830600",
      "name": "UPDATE CLIENT INTAKE INFO",
      "resource_type": "section"
    },
    {
      "gid": "1206176340563337",
      "name": "AWAITING C-FILE",
      "resource_type": "section"
    },
    {
      "gid": "1205962480830611",
      "name": "RECORD REVIEW",
      "resource_type": "section"
    },
    {
      "gid": "1205962481167456",
      "name": "SCHEDULE CALL",
      "resource_type": "section"
    },
    {
      "gid": "1205964025637422",
      "name": "STATEMENT DRAFTS",
      "resource_type": "section"
    },
    {
      "gid": "1205964025637424",
      "name": "STATEMENTS FINALIZED",
      "resource_type": "section"
    },
    {
      "gid": "1205964025637426",
      "name": "CURRENT TREATMENT",
      "resource_type": "section"
    },
    {
      "gid": "1205964025637428",
      "name": "MINI-DBQS",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667212",
      "name": "NEXUS LETTERS",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667214",
      "name": "MEDICAL TEAM",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667216",
      "name": "FILE CLAIMS (NEW OR SUPPLEMENTAL)",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667218",
      "name": "VERIFY EVIDENCE RECEIVED",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667220",
      "name": "AWAITING DECISION",
      "resource_type": "section"
    },
    {
      "gid": "1206241124546068",
      "name": "AWAITING PAYMENT (UPDATE RATING ON BOARD, SEND INVOICE AND CLOTHING EMAIL)",
      "resource_type": "section"
    },
    {
      "gid": "1205964025667224",
      "name": "COMPLETED",
      "resource_type": "section"
    }
  ]
}`
