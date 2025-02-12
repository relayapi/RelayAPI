package utils

// ProviderURLs 定义一个映射，存储提供商和对应的基础 URL
var ProviderURLs = map[string]string{
	"openai":                "https://api.openai.com/v1",
	"dashscope":             "https://dashscope.aliyuncs.com/compatible-mode/v1",
	"googleai":              "https://generativelanguage.googleapis.com/v1beta",
	"anthropic":             "https://api.anthropic.com/v1",
	"cohere":                "https://api.cohere.ai/v1",
	"huggingface":           "https://api-inference.huggingface.co/models",
	"replicate":             "https://api.replicate.com/v1",
	"ai21":                  "https://api.ai21.com/v1",
	"stabilityai":           "https://api.stability.ai/v1",
	"deepl":                 "https://api.deepl.com/v2",
	"mistralai":             "https://api.mistral.ai/v1",
	"perplexityai":          "https://api.perplexity.ai/v1",
	"baiduai":               "https://aip.baidubce.com",
	"tencentai":             "https://ai.tencentcloudapi.com",
	"googlecloudai":         "https://aiplatform.googleapis.com",
	"googlecloud":           "https://language.googleapis.com/v1",
	"aws":                   "https://comprehend.us-east-1.amazonaws.com/",
	"ibmwatson":             "https://api.us-south.language-translator.watson.cloud.ibm.com/v3",
	"deepai":                "https://api.deepai.org/api",
	"clarifai":              "https://api.clarifai.com/v2",
	"opencognitive":         "https://api.opencognitive.com/v1",
	"assemblyai":            "https://api.assemblyai.com/v2",
	"azure":                 "https://api.cognitive.microsoft.com/sts/v1.0",
	"google":                "https://dialogflow.googleapis.com/v2",
	"amazon":                "https://runtime.sagemaker.amazonaws.com/v1",
	"ibm":                   "https://api.us-south.assistant.watson.cloud.ibm.com/instances",
	"paddle":                "https://aip.baidubce.com/rpc/2.0/ai_custom",
	"tencent":               "https://api.qcloud.com/v2/index.php",
	"googleaibeta":          "https://generativelanguage.googleapis.com/v1beta",
	"eleutheraibeta":        "https://api.eleuther.ai/v1beta",
	"openai-chat":           "https://api.openai.com/v1/chat/completions",
	"deepmind":              "https://api.deepmind.com/v1",
	"faceplusplus":          "https://api-us.faceplusplus.com/v3",
	"witai":                 "https://api.wit.ai/v1",
	"dialogflow":            "https://dialogflow.googleapis.com/v2",
	"speechmatics":          "https://asr.api.speechmatics.com/v2",
	"revai":                 "https://api.rev.ai/speechtotext/v1",
	"otterai":               "https://api.otter.ai/v1",
	"voiceflow":             "https://api.voiceflow.com/v2",
	"lobeai":                "https://api.lobe.ai/v1",
	"runwayml":              "https://api.runwayml.com/v1",
	"algorithmia":           "https://api.algorithmia.com/v1",
	"bigml":                 "https://bigml.io/andromeda",
	"h2oai":                 "https://api.h2o.ai/v1",
	"datarobot":             "https://api.datarobot.com/v2",
	"rapidapi":              "https://api.rapidapi.com/v1",
	"kite":                  "https://api.kite.com/v1",
	"wolframalpha":          "https://api.wolframalpha.com/v1",
	"yandexai":              "https://api.yandex.com/v1",
	"naverclova":            "https://clova.ai/v1",
	"samsungbixby":          "https://api.bixby.com/v1",
	"microsoftaibeta":       "https://api.cognitive.microsoft.com/v1",
	"salesforceeinstein":    "https://api.einstein.ai/v2",
	"oracleaibeta":          "https://api.oracle.com/v1",
	"sapleonardo":           "https://api.sap.com/v1",
	"accentureaibeta":       "https://api.accenture.com/v1",
	"infosysnia":            "https://api.infosys.com/v1",
	"tcsai":                 "https://api.tcs.com/v1",
	"cognizantai":           "https://api.cognizant.com/v1",
	"wiproholmes":           "https://api.wipro.com/v1",
	"capgeminiaibeta":       "https://api.capgemini.com/v1",
	"atosaibeta":            "https://api.atos.com/v1",
	"deloitteaibeta":        "https://api.deloitte.com/v1",
	"eyaibeta":              "https://api.ey.com/v1",
	"pwcaibeta":             "https://api.pwc.com/v1",
	"kpmgaibeta":            "https://api.kpmg.com/v1",
	"bcgaibeta":             "https://api.bcg.com/v1",
	"mckinseyaibeta":        "https://api.mckinsey.com/v1",
	"bainaibeta":            "https://api.bain.com/v1",
	"boozallenaibeta":       "https://api.boozallen.com/v1",
	"northropgrummanaibeta": "https://api.northropgrumman.com/v1",
	"lockheedmartinaibeta":  "https://api.lockheedmartin.com/v1",
	"raytheonaibeta":        "https://api.raytheon.com/v1",
	"generaldynamicsaibeta": "https://api.generaldynamics.com/v1",
	"boeingaibeta":          "https://api.boeing.com/v1",
	"airbusaibeta":          "https://api.airbus.com/v1",
	"spacexaibeta":          "https://api.spacex.com/v1",
	"blueoriginaibeta":      "https://api.blueorigin.com/v1",
	"virgingalacticaibeta":  "https://api.virgingalactic.com/v1",
	"nasajplaibeta":         "https://api.jpl.nasa.gov/v1",
	"esaaibeta":             "https://api.esa.int/v1",
	"isroaibeta":            "https://api.isro.gov.in/v1",
	"cnsaibeta":             "https://api.cnsa.gov.cn/v1",
	"roscosmosaibeta":       "https://api.roscosmos.ru/v1",
	"jaxaaibeta":            "https://api.jaxa.jp/v1",
	"cnesaibeta":            "https://api.cnes.fr/v1",
	"dlraibeta":             "https://api.dlr.de/v1",
	"googleaiplatform":      "https://aiplatform.googleapis.com/v1",
	"amazonai":              "https://apigateway.ap-northeast-1.amazonaws.com/ai",
	"microsoftml":           "https://api.ml.azure.com/v1.0",
	"ibmai":                 "https://api.ai.ibm.com/v1",
	"baiduairest":           "https://aip.baidubce.com/rest/2.0",
	"tencentaiapi":          "https://api.ai.qq.com/",
	"aliyunai":              "https://ai.aliyun.com/api/v1",
	"huaweiai":              "https://api.hicloud.com/ai/v1",
}

// GetProviderBaseURL 根据提供者获取基础 URL
func GetProviderBaseURL(provider string) string {
	if baseURL, exists := ProviderURLs[provider]; exists {
		return baseURL
	}
	return provider // 如果找不到对应的 URL，返回原始 provider 值作为 URL
}
