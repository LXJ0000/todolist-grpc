package bootstrap

import (
	"log"
	"time"

	"github.com/LXJ0000/todolist-grpc/app/user/domain"
	"github.com/LXJ0000/todolist-grpc/pkg/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewOrmDatabase(env *Env) orm.Database {
	//db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", env.DBName)), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	// In WSL how to connect sqlite ?
	// move go-backend.db to /mnt/c/Users/JANNAN/Desktop/go-backend.db then
	// ln -s /mnt/c/Users/JANNAN/Desktop/go-backend.db ./go-backend.db
	//dsn := "root:root@tcp(127.0.0.1:3306)/go-backend?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       env.MySQLAddress,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	if err = db.AutoMigrate(
		&domain.User{},
	); err != nil {
		log.Fatal(err)
	}
	database := orm.NewDatabase(db)

	return database
}
