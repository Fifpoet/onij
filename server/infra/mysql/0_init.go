package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func NewMysqlCli() *gorm.DB {
	ip := os.Getenv("ip")
	mysqlPwd := os.Getenv("mypwd")
	dsn := "root:%s@tcp(%s:3306)/footprint?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf(dsn, mysqlPwd, ip), // DSN data source name
		DefaultStringSize:         256,                            // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                           // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                           // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                           // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                          // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "fp_", // 表名前缀，`User`表为`dp_user`
			SingularTable: true,  // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 库表自动迁移
	err = Db.AutoMigrate(&Tag{}, &Tod{}, &Relay{}, &File{}, &Music{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return Db
}
