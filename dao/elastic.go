package dao

import (
	"github.com/olivere/elastic/v7"
)

func ConnectES() *elastic.Client {
	// // 创建Client, 连接ES
	// client, err := elastic.NewClient(
	// 	// elasticsearch 服务地址，多个服务地址使用逗号分隔
	// 	elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"),
	// 	// 基于http base auth验证机制的账号和密码
	// 	elastic.SetBasicAuth("user", "secret"),
	// 	// 启用gzip压缩
	// 	elastic.SetGzip(true),
	// 	// 设置监控检查时间间隔
	// 	elastic.SetHealthcheckInterval(10*time.Second),
	// 	// 设置错误日志输出
	// 	elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	// 	// 设置info日志输出
	// 	elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	// if err != nil {
	// 	// Handle error
	// 	log.Fatalf("连接失败: %v\n", err)
	// }
	// return client
	return nil
}
