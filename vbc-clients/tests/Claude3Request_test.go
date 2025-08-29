package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Gen_API_AiResultString(t *testing.T) {
	aa := biz.Claude3Request{
		Messages: []biz.Claude3Message{
			{
				Role: "user",
				Content: []biz.Claude3Content{
					{
						Type: "text",
						Text: "From the list below, identify the entries that are directly or indirectly related to \"Tinnitus\". A connection can be:\n- the condition is caused by Tinnitus (e.g., \"secondary to tinnitus\"),\n- Tinnitus is a mentioned symptom or diagnosis,\n- or Tinnitus is part of the main condition description.\n- Only return the lines that are relevant to Tinnitus.\n\nData list:\n5004-Back-New-Back pain with radiculopathy secondary to Right knee strain with tibial stress fractures -5974846674216669565.pdf  \n5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf  \n5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf  \n5004-Knee-Increase-Right knee strain with tibial stress fractures with limitation of flexion and extension (increase)Bilateral left knee pain secondary to Right knee strain with tibial stress fractures with limitation of flexion and extension (opinion)-5974845864217999778.pdf  \n5004-Mental Disorders Secondaries-New-Major Depressive Disorder secondary to Right knee strain with tibial stress fractures-5974857144211924747.pdf",
					},
				},
			},
		},
	}
	c := biz.InterfaceToString(aa)
	lib.DPrintln(c)

	//aca := "{\"anthropic_version\":\"\",\"max_tokens\":0,\"messages\":[{\"role\":\"user\",\"content\":[{\"type\":\"text\",\"text\":\"From the list below, identify the entries that are directly or indirectly related to \\\"Tinnitus\\\". A connection can be:\\n- the condition is caused by Tinnitus (e.g., \\\"secondary to tinnitus\\\"),\\n- Tinnitus is a mentioned symptom or diagnosis,\\n- or Tinnitus is part of the main condition description.\\n- Only return the lines that are relevant to Tinnitus.\\n\\nData list:\\n5004-Back-New-Back pain with radiculopathy secondary to Right knee strain with tibial stress fractures -5974846674216669565.pdf  \\n5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf  \\n5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf  \\n5004-Knee-Increase-Right knee strain with tibial stress fractures with limitation of flexion and extension (increase)Bilateral left knee pain secondary to Right knee strain with tibial stress fractures with limitation of flexion and extension (opinion)-5974845864217999778.pdf  \\n5004-Mental Disorders Secondaries-New-Major Depressive Disorder secondary to Right knee strain with tibial stress fractures-5974857144211924747.pdf\"}]}]}"
}
