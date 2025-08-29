package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

// 有
func Test_HaiUsecase_GetDiseaseNamesByMedicalTextWithAI_1(t *testing.T) {
	text := `.\n' NAVMED 6150 3 (REV 7.72) FRONT (MODIFIED)\nHEALTH RECORD\nSICK CALL TREATMENT RECORD\nDATE\nNAME OF TREATING FACILITY. COMPLAINT, TREATMENT ADMINISTERED. SIGNATURE AND GRADE/RATE OF PERSON ADMINISTERING TREATMENT\n4CCTTL\nMISS DETROIT (ACC. 4)\nHEADACHE,RUNNY NOSE SORE THROAT REYAM TeP. 99-2; THROAT CLEAR -\nRX; ACTIFED 2 T TIO - TyLENOL #3X TI GID .\n54 11 acum 5N\n.\n-.\n22 oct 74\nUSS Detail AOE 4\n/0 Nausea, vanutri\n.. 5. Vanuitil x q today has felt Nauseated for the past week aldo loss of Appetite for\n. : the past week. It was Recruit at San Diego\nrate it Callers # 5 + 8 dienas\nAugust -left there on Aug 12 1974 + arrived\nNon t outbreak of barnabitis. Has not been -\ntreated for marasite infection since lite Jal\nNouses queaus to coincide with rough weather\n0- Tem. 97: No salerad una interes\n. Abd - No larato galenomical BSN. No mais tardemen. V/A - WN2\n. (A) Prote motion unkness. v. wird gastroenteritis.\n(p) Compagne tale now + 6h. #3\nDramamine # 15 min pausea .\nE-a Castle Stue\nSEX\nMALE Malayer\nRACE\nGRADE. RATING, OR POSITION\nORGANIZATION UNIT\nCOMPONENT OR BRANCH USN\nSERVICE. DEPT. OR AGENCY\nDOD\nSR\nNTC SDIEGO\nPATIENT 'S LAST NAME . FIRST NAME . MIDDLE NAME\nFANQUILUT REMIGIO DIMALANTA 75 C168\nDATE OF BIRTH ( DAY-MONTH . YEAR)\nIDENTIFICATION NO\n01-21-49\n20/ 571299855/\n21Jan44\nSICK CALL TREATMENT RECORD NAVMED 6150/3\n3\nIME\n(3F)`
	response, err := UT.HaiUsecase.GetDiseaseNamesByMedicalTextWithAI(context.TODO(), text, "test")

	lib.DPrintln(response, err)
	/*

			```json
		{
		  "diseaseNames": [
		    "HEADACHE",
		    "RUNNY NOSE",
		    "SORE THROAT",
		    "Nausea",
		    "loss of Appetite",
		    "barnabitis",
		    "marasite infection",
		    "gastroenteritis"
		  ]
		}
		```


	*/
}

// 无
func Test_HaiUsecase_GetDiseaseNamesByMedicalTextWithAI_2(t *testing.T) {
	text := `DATE\nSYMPTOMS, DIAGNOSIS, TREATMENT, TREATING ORGANIZATION (Sign each entry)\n16OCT 96\nCont 1 Volar Split / No un of (R) hand at work.\nHunAls LEON\n:\n.\n.\n.\n*U.S. Government Printing Office: 1996 - 404-763/20098\nSTANDARD FORM 600 BACK (REV. 5-84)\n-\n-`
	response, err := UT.HaiUsecase.GetDiseaseNamesByMedicalTextWithAI(context.TODO(), text, "test")

	lib.DPrintln(response, err)
	/*
		```json
		{
		  "diseaseNames": []
		}
		```
	*/
}

// 无
func Test_HaiUsecase_GetDiseaseNamesByMedicalTextWithAI_3(t *testing.T) {
	text := `NSN 7540-00-634-4178\nHEALTH RECORD\nDATE\n16 OCT 96\nS/47 40 Q> RHD/RA Symptoms. C/0 9 mos, pan\n600108\nCHRONOLOGICAL RECORD OF MEDICAL CARE\nSYMPTOMS. DIAGNOSIS, TREATMENT TREATING ORGANIZATION (Sign each entry) ORTHOPEDIC HAND SURGERY\nNAVAL MEDICAL CENTER SAN DIEGO, CA\nAppointments: 532-8429\nat RF MC-P joint with gripping\nCystic man at this area slowly\ncow ley Li 40\nSize\nAlso do painful mas on dorsum g (R) writ\n& poor aspirations\n0 1 RF MC-P jout Volar Saface à tender Moveable module 3-4 mm.\nXray - wik FDS/FOP/ EDC/ GPL/FPL /punch/IO 5/5\n2 Dorsal wrist & Kan man TIP Quentantad & Fluxuri\nCRT\nL25\n2pt\n25ML\nRefinacular apt 4/ 1 Dorsal wrist ganglion\nDWG-> P(2) Reformular cyst desires\nNo fluid Bolle Symptomatic Ptaspiration\nretrieved\nand\nwould leta Resultan, scheduled\n26 Nov 96\nAspirato Both Cysts toly\nPATIENT'S IDENTIFICATION (Use this space for Mechanical Imprint)\n5\n2130719 BENICIO\nRECORDS MAINTAINED AT:\nSwals Place Celestine -> Volar Spunt removablex\nPATIENT'S NAME (Lot, First, Midi(Initial) STATUS\n1\nRELATIONSHIP TO SPONSOR\nRANK/GRADE\n611 679-04\nOUTPATIENT RECAKUSS\nSPONSOR'S NAME\nORGANIZATION\nDEPART./SERVICE |SSN/IDENTIFICATION NO.\nDATE OF BIRTH\nCHRONOLOGICAL RECORD OF MEDICAL CARE\nSTANDARD FORM 600 (REV. 5-84) Prescribed by GSA and ICMR FIRMR (41 CFR) 201-45:505`
	response, err := UT.HaiUsecase.GetDiseaseNamesByMedicalTextWithAI(context.TODO(), text, "test")

	lib.DPrintln(response, err)
	/*
		```json
		{
		  "diseaseNames": []
		}
		```
	*/
}
