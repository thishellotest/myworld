package configs

import (
	"os"
)

func EnvMySQLDSN() string {
	return os.Getenv("MYSQL_DSN")
}

func EnvRedisDSN() string {
	return os.Getenv("REDIS_DSN")
}

func EnvMailGlliaoPWD() string {
	return os.Getenv("MailGlliao_PWD")
}

func EnvMailYwangPWD() string {
	return os.Getenv("MailYwang_PWD")
}

func EnvMailTeamAgsPWD() string {
	return os.Getenv("MailTeamAgs_PWD")
}

func EnvJotformKey() string {
	return os.Getenv("JOTFORM_KEY")
}

func EnvAzureStorageAccountName() string {
	return os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
}

func EnvAzureStorageAccountKey() string {
	return os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
}

func EnvAzureCognitiveKey() string {
	return os.Getenv("AZURE_COGNITIVE_KEY")
}

func EnvAzureGptKey() string {
	return os.Getenv("AZURE_GPT_KEY")
}

func EnvDialpadKey() string {
	return os.Getenv("DIALPAD_KEY")
}

func EnvDialpadWebhookSecret() string {
	return os.Getenv("DIALPAD_WEBHOOK_SECRET")
}

func EnvZoomAccountIdKey() string {
	return os.Getenv("ZOOM_ACCOUNT_ID")
}

func EnvZoomClientIdKey() string {
	return os.Getenv("ZOOM_CLIENT_ID")
}

func EnvZoomSecretKey() string {
	return os.Getenv("ZOOM_SECRET")
}

func EnvSensitiveDataKey() string {
	return os.Getenv("SENSITIVE_DATA_KEY")
}
