package tests

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"testing"
	"vbc/lib"
	"vbc/lib/to"
)

func Test_AzopenaiUsecase_AzTextEbdSmall(t *testing.T) {
	a, err := UT.AzopenaiUsecase.AzTextEbdSmall(context.TODO(), azopenai.EmbeddingsOptions{
		Input:      []string{"hello"},
		Dimensions: to.Ptr(int32(512)),
	})
	lib.DPrintln(a)
	lib.DPrintln(err)
}

func Test_AzopenaiUsecase_AskGpt4o(t *testing.T) {
	result, err := UT.AzopenaiUsecase.AskGpt4o(context.TODO(), "", "hello")
	lib.DPrintln(err)
	lib.DPrintln(result)
}

func Test_AzopenaiUsecase_AskGpt4o_OutJSON(t *testing.T) {

	systemConfig := `# Role
学生信息提取器

## Skills
- 精通中文
- 能够理解文本
- 精通JSON数据格式

## Action
- 根据提学生的描述，提取出对应的 姓名、年龄、父母工作、家庭住址、爱好 等信息，并以JSON格式输出

## Constrains
- 忽略无关内容
- 针对学生的自我介绍你可以采用更合理的表述方式
- 针对父母工作你可以根据描述，转换为更合理的职业名称
- 必须保证你的结果只包含一个合法的JSON格式

## Format
- 对应JSON的key为：name, age, father_job, mother_job, address, hobby`

	content := "我叫张三，今年9岁了，我的家住在幸福小区，我的妈妈是老师，我的爸爸是医生。我喜欢画画、唱歌。"
	res, err := UT.AzopenaiUsecase.AskGpt4o(context.TODO(), systemConfig, content)
	lib.DPrintln(res, err)

}

func Test_AzopenaiUsecase_AzGpt4o(t *testing.T) {

	response, err := UT.AzopenaiUsecase.AzGpt4o(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			//&azopenai.ChatRequestSystemMessage{
			//	Content: to.Ptr("You're an intelligent assistant"),
			//},
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("Hello"),
			},
			//&azopenai.ChatRequestAssistantMessage{
			//	Content: to.Ptr(""),
			//},
			//&azopenai.ChatRequestUserMessage{
			//	Content: azopenai.NewChatRequestUserMessageContent(""),
			//},
		},
		MaxTokens:   to.Ptr(int32(2048 - 127)), // 最大响应数
		Temperature: to.Ptr(float32(0.0)),      // 温度 0.7
		//TopP:        to.Ptr(float32(0.95)),     // 默认 0。95 暂不定义
	})
	lib.DPrintln(err)
	lib.DPrintln(response)
}
