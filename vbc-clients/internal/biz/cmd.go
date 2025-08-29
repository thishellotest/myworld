package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
)

type CmdUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	BoxUsecase    *BoxUsecase
}

func NewCmdUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxUsecase *BoxUsecase,
) *CmdUsecase {
	uc := &CmdUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		BoxUsecase:    BoxUsecase,
	}

	return uc
}

func (c *CmdUsecase) RunBackup() error {
	path, fileName, err := c.ExeBackupMySQL()
	if err != nil {
		return err
	}
	destPath, compressFileName, err := c.CompressBackup(path, fileName)
	if err != nil {
		return err
	}
	cmd := exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", path, fileName))
	er := cmd.Run()
	if er != nil {
		c.log.Error(er)
	}
	c.log.Info("destPath: ", destPath, " compressFileName: ", compressFileName)
	err = c.UploadToBox(destPath, compressFileName)
	if err != nil {
		return err
	}
	cmd = exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", destPath, compressFileName))
	er = cmd.Run()
	if er != nil {
		c.log.Error(er)
	}

	return nil
}

func (c *CmdUsecase) UploadToBox(path string, fileName string) error {

	file := fmt.Sprintf("%s/%s", path, fileName)
	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = c.BoxUsecase.UploadFile("314164957958", fp, fileName)

	return err
}

func (c *CmdUsecase) CompressBackup(path string, fileName string) (destPath, compressFileName string, err error) {
	destPath = path
	compressFileName = fileName + ".tar.gz"
	cmd := exec.Command("tar", "-czvf", fmt.Sprintf("%s/%s", destPath, compressFileName), "-C", path, fileName)
	// 执行命令并捕获标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}

	return destPath, compressFileName, nil
}

func (c *CmdUsecase) ExeBackupMySQL() (path string, fileName string, err error) {
	/*
		mysqldump -h snzor5.stackhero-network.com -P 3306 -u vbc_readonly -p'LYXXmyz8z1WUP8Je' --single-transaction  --no-tablespaces  vbcdbzoho > vbcdbzoho_4.sql
		mysqldump -h snzor5.stackhero-network.com -P 3306 -u vbc_readonly -p'LYXXmyz8z1WUP8Je' --single-transaction  --no-tablespaces  vbcdbzoho --result-file=/home/azureuser/vbcdbzoho_5.sql
	*/

	cfg, err := mysql.ParseDSN(configs.EnvMySQLDSN())

	if err != nil {
		return "", "", err
	}
	addrs := strings.Split(cfg.Addr, ":")
	port := "3306"
	addr := ""
	if len(addrs) == 2 {
		addr = addrs[0]
		port = addrs[1]
	}
	fileName = fmt.Sprintf("vbcdb_%s.sql", time.Now().In(configs.GetVBCDefaultLocation()).Format("2006-01-02_15_04"))
	path = "/home/azureuser"
	// 设置 mysqldump 命令和参数 密码在一起
	cmd := exec.Command("mysqldump", "-h", addr, "-P", port, "-u", cfg.User, "-p"+cfg.Passwd, "--single-transaction", "--no-tablespaces", cfg.DBName,
		fmt.Sprintf("--result-file=%s", fmt.Sprintf("%s/%s", path, fileName)))

	// 执行命令并捕获标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}
	c.log.Info("Backup completed successfully")
	return path, fileName, nil
}
