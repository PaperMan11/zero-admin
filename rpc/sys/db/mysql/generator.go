// configuration.go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./rpc/sys/db/mysql/query", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	var dsn = "root:123456@tcp(192.168.241.128:3306)/zeroadmin?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	// Initialize a *gorm.DB instance
	db, _ := gorm.Open(mysql.Open(dsn))

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("sys_user"),
		g.GenerateModel("sys_role"),
		g.GenerateModel("sys_user_role"),
		g.GenerateModel("sys_scope"),
		g.GenerateModel("sys_role_scope"),
		g.GenerateModel("sys_menu"),
		g.GenerateModel("sys_login_log"),
		g.GenerateModel("sys_operate_log"),
	)

	// Execute the generator
	g.Execute()
}
