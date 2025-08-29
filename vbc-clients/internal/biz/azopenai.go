package biz

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/to"
)

type AzopenaiUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewAzopenaiUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *AzopenaiUsecase {
	uc := &AzopenaiUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	return uc
}

func AzureGptKey() string {
	return configs.EnvAzureGptKey()
}

func (c *AzopenaiUsecase) AzTextEbdSmall(ctx context.Context, embeddingsOptions azopenai.EmbeddingsOptions) (res azopenai.GetEmbeddingsResponse, err error) {

	deployConfig := lib.ToTypeMapByString(`{"Key":"` + AzureGptKey() + `","DeploymentName":"textebd3small"}`)
	keyCredential := azcore.NewKeyCredential(deployConfig.GetString("Key"))
	InstanceName := "openaieu2"
	client, err := azopenai.NewClientWithKeyCredential("https://"+InstanceName+".openai.azure.com",
		keyCredential, nil)
	if err != nil {
		c.log.Error(err)
		return res, err
	}
	deploymentName := deployConfig.GetString("DeploymentName")
	/*
		azopenai.EmbeddingsOptions{
				Input:          []string{"hello"},
				Dimensions:     to.Ptr(int32(512)),
				DeploymentName: &deploymentName,
			}
	*/
	embeddingsOptions.DeploymentName = &deploymentName
	return client.GetEmbeddings(ctx, embeddingsOptions, nil)
	/*
		{
			"data": [{
				"embedding": [0.023896476, -0.07919617, 0.008054113, 0.09416796, 0.012709338, -0.06722708, -0.039952572, 0.08691144, -0.0031408435, -0.06218088, 0.0139917405, -0.047334205, -0.017380202, -0.041036878, 0.020758238, 0.08274102, -0.09950609, 0.0533396, 0.014429634, 0.058010466, 0.09441818, 0.014398356, -0.007934214, 0.0337178, 0.030819364, 0.016869327, -0.006808202, 0.027441328, 0.034656145, -0.09224957, 0.047334205, -0.0633486, 0.028963529, -0.0069802315, 0.015409682, -0.021060593, -0.037366915, 0.049336005, 0.0033415447, -0.047334205, -0.031444926, -0.012146332, 0.07677733, 0.03459359, -0.00005188582, 0.022165753, -0.029818464, -0.0061357226, 0.04587456, 0.042830158, 0.009732705, -0.02425096, 0.06109657, 0.1821637, 0.06280644, -0.025105895, 0.050962467, 0.033238202, 0.014513043, 0.010530297, -0.034447625, 0.008106243, -0.015503516, 0.013220214, 0.01483625, -0.016775493, -0.0033441512, 0.014023019, -0.017536594, 0.025189305, -0.0012980415, 0.07394145, 0.013981314, -0.034197398, 0.0110828765, -0.017140403, -0.06585085, -0.032570936, -0.014033445, 0.010227942, -0.0874953, 0.0027993908, -0.018850274, -0.05321449, -0.055966962, 0.00887777, -0.100173354, 0.004365903, -0.01003506, -0.027962629, -0.0206227, 0.04837681, -0.0038185357, -0.019663505, 0.029943576, 0.005012317, -0.01449219, -0.00430074, 0.09316706, 0.030798512, 0.079821736, -0.06530869, 0.02675321, -0.01559735, 0.028108595, 0.07052171, 0.037554584, 0.00882564, -0.041870963, 0.041516475, -0.11702183, -0.06464142, -0.040661544, -0.033822063, 0.028796712, 0.08757871, 0.08411726, -0.12744787, 0.010264433, -0.07173113, -0.007783036, -0.038784854, 0.022019789, -0.04454003, -0.11460299, 0.010790948, -0.010274859, 0.039181046, -0.060762938, 0.00046037466, 0.000027083259, 0.03177856, -0.010822225, -0.0019327265, -0.069854446, 0.02325006, -0.08119797, -0.0066987285, -0.03742947, 0.0061357226, 0.064099275, -0.029755907, 0.019653078, -0.04708398, -0.01328277, -0.047709543, -0.026106795, 0.07673563, -0.04520729, -0.01625419, 0.022270015, 0.05221359, -0.08274102, 0.010029847, -0.027191103, -0.03559449, 0.08432578, -0.045415815, -0.022874724, 0.02433437, 0.057927057, 0.018746013, -0.02506419, -0.046249896, -0.06272303, 0.020257788, 0.020643553, 0.06543381, -0.032341566, 0.028796712, -0.026231907, 0.03432251, -0.073899746, -0.009649296, -0.003943648, 0.015712038, 0.042308856, -0.024000736, -0.1106828, -0.0017789424, 0.022499386, 0.12928285, 0.0030835003, -0.039222747, -0.043497425, -0.011583326, -0.00005860573, -0.034051433, 0.0073555685, 0.046333306, -0.050962467, 0.07727778, 0.06288985, 0.025877422, 0.010728392, -0.06147191, 0.12252678, -0.032279007, -0.004444098, 0.0014518256, 0.023875624, -0.038993377, 0.051754843, -0.015722463, -0.035990678, -0.0013814499, 0.004394574, 0.06376564, 0.050170086, -0.037783954, 0.08824597, -0.02354199, -0.0015847575, 0.06288985, -0.039994277, -0.0834917, 0.08908006, 0.0039592874, 0.049878158, -0.017557446, -0.06097146, 0.0533396, 0.10692943, 0.029338866, 0.03250838, -0.02354199, 0.0058594323, -0.028046038, 0.02304154, -0.022103198, 0.040598985, 0.059011366, -0.018193433, -0.050962467, 0.02254109, 0.059345, 0.030444026, -0.04253823, 0.00047731696, 0.0057655983, -0.09516886, -0.0019613982, 0.0030261572, 0.02616935, 0.035073187, 0.0064485036, 0.008580628, -0.07836209, 0.005463243, 0.003589163, 0.019872025, 0.053381305, 0.030214654, -0.07181454, 0.023854772, 0.003265956, -0.0025621983, 0.03058999, -0.012709338, -0.0785289, -0.0336761, -0.019183908, 0.0066570244, -0.041912667, -0.04658353, -0.072273284, 0.0014140311, 0.01207335, -0.020945907, 0.023479434, 0.0140542975, -0.022499386, 0.061388504, -0.01036348, -0.061221685, 0.007725693, 0.024104996, -0.00727216, 0.060095675, 0.043872762, -0.08983073, 0.005145249, 0.009670149, -0.003865453, -0.074066564, 0.034843814, 0.03770055, -0.068019465, 0.040890913, 0.041912667, 0.03559449, 0.007386847, 0.022290865, -0.0286716, -0.025668902, 0.039848313, 0.052088477, -0.012928285, 0.030464878, -0.06326519, 0.08707826, 0.013376605, 0.00844509, -0.030131245, -0.03469785, 0.064140975, 0.0022155328, -0.017317647, 0.033071388, -0.023500286, -0.05221359, 0.03936871, 0.073774636, -0.0017971881, -0.059219886, -0.044915363, 0.03451018, 0.0052364767, 0.014648581, -0.038096737, 0.019851172, -0.04687546, -0.007960279, -0.013022119, -0.029463978, 0.0008073662, 0.046625234, -0.005160888, -0.05037861, -0.042788453, -0.004326805, -0.0914989, 0.013950037, -0.003941042, 0.07531769, 0.005755172, -0.054840952, 0.028317114, -0.046792053, 0.037095837, -0.012803173, -0.07656881, -0.029318014, 0.00766835, -0.019527966, -0.010134107, 0.04162074, -0.022916429, -0.01031135, 0.0058385804, -0.0012400467, 0.04745932, -0.016222913, -0.018339397, -0.0108118, -0.06831139, 0.084367484, -0.032279007, -0.032842014, 0.07669392, -0.051337805, -0.023959031, 0.017056996, -0.053381305, 0.07319078, 0.08332488, -0.074316785, -0.01870431, 0.030569138, 0.04666694, -0.0038419943, -0.014502617, 0.05121269, 0.032049637, -0.08474282, 0.034843814, 0.011124581, -0.0061148703, 0.050962467, 0.02322921, 0.079821736, 0.034572735, -0.0050905123, -0.045957968, 0.04155818, -0.022353422, 0.012208888, -0.025898274, -0.0336761, -0.02483482, 0.11702183, -0.013241067, 0.055675033, -0.061388504, 0.012323575, -0.020883352, 0.04247567, -0.015430534, -0.037054133, -0.030840216, 0.02685747, -0.07786164, 0.027545588, -0.007381634, -0.026711505, -0.030193802, -0.031820264, 0.0020096186, 0.046625234, -0.039118487, 0.034676995, -0.055174585, 0.07056341, 0.051546324, -0.008033261, -0.023312617, -0.014586025, -0.042100336, 0.014460912, -0.02022651, 0.035281707, -0.011635456, -0.013449587, 0.06472483, -0.042600784, 0.06543381, 0.07869572, -0.050545424, 0.017880652, -0.102091745, -0.073983155, -0.008705741, -0.0125738, 0.01575374, -0.023083245, 0.0010830045, 0.002657336, 0.03974405, -0.0004822042, -0.008387746, -0.015743315, -0.018735588, -0.02254109, -0.07360782, 0.00053400855, -0.01575374, -0.015357551, -0.0013241066, 0.0012354853, -0.029109493, -0.007558876, -0.032216452, 0.02646128, 0.0052104117, -0.0034223464, 0.028087743, -0.019653078, -0.003477083, 0.01212548, 0.0023862591, 0.007423338, -0.030965328, 0.017567871, 0.06981274, 0.024313517, -0.060596123, -0.030861069, -0.027962629, -0.036783054, -0.05021179, 0.044122986, 0.05613378, 0.0814899, -0.016921457, -0.027983481, -0.0021034528, -0.00005229309, -0.005927202, -0.0022415977, -0.029901871, -0.003901944, 0.004050515, -0.028296262, 0.072314985, -0.019788617, 0.105261266, 0.0138666285, 0.0065214857, -0.068353094, -0.0023927754, 0.019851172, 0.0693957, 0.07456701, -0.026315317, -0.057176385, -0.024522038, -0.0012934802, -0.0015378403, -0.054882657, -0.07298225, -0.03478126, -0.02625276, -0.021561043, 0.0053381305, 0.08319977, -0.023792215, -0.02817115, -0.0025361334, 0.020549718, -0.029964428, 0.013188936, 0.03209134],
				"index": 0
			}],
			"usage": {
				"prompt_tokens": 1,
				"total_tokens": 1
			}
		}
	*/
}

func (c *AzopenaiUsecase) AskGpt4o(ctx context.Context, systemConfig string, prompt string) (result string, err error) {

	var messages []azopenai.ChatRequestMessageClassification
	if systemConfig != "" {
		messages = append(messages, &azopenai.ChatRequestSystemMessage{
			Content: to.Ptr(systemConfig),
		})
	}
	messages = append(messages, &azopenai.ChatRequestUserMessage{
		Content: azopenai.NewChatRequestUserMessageContent(prompt),
	})

	response, err := c.AzGpt4o(ctx, azopenai.ChatCompletionsOptions{
		Messages:    messages,
		MaxTokens:   to.Ptr(int32(2048 - 127)), // 最大响应数
		Temperature: to.Ptr(float32(0.0)),      // 温度 0.7
		//TopP:        to.Ptr(float32(0.95)),     // 默认 0。95 暂不定义
	})
	if err != nil {
		c.log.Error("AskGpt4o: ", InterfaceToString(messages), " : ", InterfaceToString(response), " : ", err.Error())
		return "", err
	}
	if len(response.Choices) <= 0 {
		c.log.Error("AskGpt4o1: ", InterfaceToString(messages), " : ", InterfaceToString(response))
		return "", errors.New("response.Choices is wrong")
	}
	if response.Choices[0].Message.Content == nil || *response.Choices[0].Message.Content == "" {
		c.log.Error("AskGpt4o2: ", InterfaceToString(messages), " : ", InterfaceToString(response))
		return "", errors.New("response.Choices[0].Message.Content is wrong")
	}
	return *response.Choices[0].Message.Content, nil
}

func (c *AzopenaiUsecase) AzGpt4o(ctx context.Context, ChatCompletionsOptions azopenai.ChatCompletionsOptions) (*azopenai.GetChatCompletionsResponse, error) {

	deployConfig := lib.ToTypeMapByString(`{"Key":"` + AzureGptKey() + `","DeploymentName":"gpt4o"}`)
	keyCredential := azcore.NewKeyCredential(deployConfig.GetString("Key"))
	InstanceName := "openaieu2"

	client, err := azopenai.NewClientWithKeyCredential("https://"+InstanceName+".openai.azure.com",
		keyCredential, nil)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	//lib.DPrintln("client:", client)
	//lib.DPrintln("err:", err)
	//_ = client

	deploymentName := deployConfig.GetString("DeploymentName")
	ChatCompletionsOptions.DeploymentName = &deploymentName
	resp, err := client.GetChatCompletions(ctx, ChatCompletionsOptions, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *AzopenaiUsecase) AzGpt35(ctx context.Context, ChatCompletionsOptions azopenai.ChatCompletionsOptions) (*azopenai.GetChatCompletionsResponse, error) {

	deployConfig := lib.ToTypeMapByString(`{"Key":"","DeploymentName":""}`)
	keyCredential := azcore.NewKeyCredential(deployConfig.GetString("Key"))
	InstanceName := ""
	client, err := azopenai.NewClientWithKeyCredential("https://"+InstanceName+".openai.azure.com",
		keyCredential, nil)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	deploymentName := deployConfig.GetString("DeploymentName")
	ChatCompletionsOptions.DeploymentName = &deploymentName
	resp, err := client.GetChatCompletions(ctx, ChatCompletionsOptions, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
