package exam_api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestGen(t *testing.T) {
	// 初始化数据库连接
	dsn := "root:Eas123!@tcp(115.190.122.151:3306)/cp_test?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	// 创建生成器配置
	g := gen.NewGenerator(gen.Config{
		// 关键：设置实体包路径
		ModelPkgPath: "./internal/data/entity",
		OutPath:      "",                 // 关键：禁用查询文件生成
		Mode:         gen.WithoutContext, // 可选：不生成带context的方法
	})

	// 使用数据库连接
	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable(
		gen.FieldModify(func(field gen.Field) gen.Field {

			// 统一处理时间类型字段
			timeTypes := map[string]bool{
				"time.Time": true,
			}

			if timeTypes[field.Type] {
				// 版本兼容的 NULL 检查
				isNullable := false

				if _, ok := field.GORMTag["not null"]; !ok {
					isNullable = true
				}

				if isNullable {
					field.Type = "*time.Time"
				}
			}
			return field
		}),
	)...)

	//g.GenerateAllTable()
	// 执行生成
	g.Execute()
}
