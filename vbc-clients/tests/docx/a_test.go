// Copyright 2017 Baliance. All rights reserved.

package docx

import (
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"fmt"
	"log"
	"testing"
)

var lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin lobortis, lectus dictum feugiat tempus, sem neque finibus enim, sed eleifend sem nunc ac diam. Vestibulum tempus sagittis elementum`

func TestA(t *testing.T) {
	// When Word saves a document, it removes all unused styles.  This means to
	// copy the styles from an existing document, you must first create a
	// document that contains text in each style of interest.  As an example,
	// see the template.docx in this directory.  It contains a paragraph set in
	// each style that Word supports by default.
	doc, err := document.OpenTemplate("/Users/garyliao/code/vbc-clients/resource/Personal_Statements_tpl.docx")
	if err != nil {
		log.Fatalf("error opening Windows Word 2016 document: %s", err)
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	for _, s := range doc.Styles.Styles() {
		fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	}
	return

	Normal := "1"
	SeptalLine := "33"
	ListParagraph := "34"

	para := doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("Test Test-70#5120")

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("1.50 - Headaches secondary to hypertension (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("20* - Cervical pain secondary to degenerative disc disease with lumbosacral strain (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("20* - Bilateral knee pain with limitation of flexion and extension secondary to degenerative disc disease with lumbosacral strain and chronic coccyx contusion (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("4.20* - Right hip pain with limitation of flexion and extension secondary to degenerative disc disease with lumbosacral strain and chronic coccyx contusion (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("20* - Left leg radiculopathy secondary to degenerative disc disease with lumbosacral strain and chronic coccyx contusion (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("123 20* - Cervical pain secondary to degenerative disc disease with lumbosacral strain (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(ListParagraph)
	para.AddRun().AddText("+++ 20* - Cervical pain secondary to degenerative disc disease with lumbosacral strain (opinion)")

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	//para = doc.AddParagraph()
	//para.SetStyle("Subtitle")
	//para.AddRun().AddText("Document Subtitle")
	//
	//para = doc.AddParagraph()
	//para.SetStyle("Heading1")
	//para.AddRun().AddText("Major Section")
	//para = doc.AddParagraph()
	//para = doc.AddParagraph()
	//for i := 0; i < 4; i++ {
	//	para.AddRun().AddText(lorem)
	//}
	//
	//para = doc.AddParagraph()
	//para.SetStyle("Heading2")
	//para.AddRun().AddText("Minor Section")
	//para = doc.AddParagraph()
	//for i := 0; i < 4; i++ {
	//	para.AddRun().AddText(lorem)
	//}
	//
	//// using a pre-defined table style
	//table := doc.AddTable()
	//table.Properties().SetWidthPercent(90)
	//table.Properties().SetStyle("GridTable4-Accent1")
	//look := table.Properties().TableLook()
	//// these have default values in the style, so we manually turn some of them off
	//look.SetFirstColumn(false)
	//look.SetFirstRow(true)
	//look.SetLastColumn(false)
	//look.SetLastRow(true)
	//look.SetHorizontalBanding(true)
	//
	//for r := 0; r < 5; r++ {
	//	row := table.AddRow()
	//	for c := 0; c < 5; c++ {
	//		cell := row.AddCell()
	//		cell.AddParagraph().AddRun().AddText(fmt.Sprintf("row %d col %d", r+1, c+1))
	//	}
	//}

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)
	//section.SetPageSizeAndOrientation(measurement.Inch*8.3, measurement.Inch*11.7, wml.ST_PageOrientationLandscape)

	doc.SaveToFile("use-template2.docx")
}
