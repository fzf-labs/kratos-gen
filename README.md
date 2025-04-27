# Kratos-Gen

## 项目介绍

Kratos-Gen 是一个用于 [Kratos](https://github.com/go-kratos/kratos) 微服务框架的代码生成工具，可以根据 protobuf 文件或数据库表结构自动生成 service 和 data 层的代码，大大提高开发效率。

## 特性

- 根据 protobuf 文件生成 service 层代码
- 根据数据库表结构生成 data 层代码和 biz 层接口
- 支持 MySQL 和 PostgreSQL 数据库
- 自动处理依赖注入和 wire 配置
- 支持自定义输出路径

## 安装

```bash
go install github.com/fzf-labs/kratos-gen@latest
```

或者从源码编译：

```bash
git clone https://github.com/fzf-labs/kratos-gen.git
cd kratos-gen
go build -o kratos-gen .
```

## 使用方法

### 生成 service 层代码

根据 protobuf 文件生成 service 层代码：

```bash
kratos-gen service [flags]
```

参数说明：

```
Flags:
  -h, --help                       帮助信息
  -i, --inPutPbPath string         protobuf 文件输入路径 (默认 "./api")
  -o, --outPutServicePath string   service 层代码输出路径 (默认 "./internal/service")
```

### 生成 data 层代码

根据数据库表结构生成 data 层代码：

```bash
kratos-gen data [flags]
```

参数说明：

```
Flags:
      --db string               数据库类型 (mysql 或 postgres)
      --dsn string              数据库连接字符串
  -h, --help                    帮助信息
      --outPutBizPath string    biz 层接口代码输出路径 (默认 "./internal/biz")
      --outPutDataPath string   data 层代码输出路径 (默认 "./internal/data")
      --tables string           指定要生成的表名，多个表用逗号分隔，不指定则生成所有表
```

## 示例

### 生成 service 层代码

```bash
# 使用默认路径
kratos-gen service

# 自定义路径
kratos-gen service -i ./proto -o ./internal/service
```

### 生成 data 层代码

```bash
# MySQL 示例
kratos-gen data --db mysql --dsn "user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"

# PostgreSQL 示例
kratos-gen data --db postgres --dsn "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"

# 指定表
kratos-gen data --db mysql --dsn "user:password@tcp(127.0.0.1:3306)/database_name" --tables "user,order,product"
```

## 项目结构

```
├── data        # 数据层代码生成相关
│   ├── cmd.go  # data 命令定义
│   ├── data.go # data 代码生成逻辑
│   └── tpl     # data 层模板
├── service     # 服务层代码生成相关
│   ├── cmd.go  # service 命令定义
│   ├── service.go # service 代码生成逻辑
│   └── tpl     # service 层模板
├── proto       # protobuf 解析
├── utils       # 工具函数
├── go.mod      # Go 模块定义
├── go.sum      # Go 依赖版本
└── main.go     # 主程序入口
```

## 依赖

- Go 1.22+
- github.com/dave/dst
- github.com/emicklei/proto
- github.com/samber/lo
- github.com/spf13/cobra
- gorm.io/gorm
- gorm.io/driver/mysql
- gorm.io/driver/postgres

## 贡献

欢迎提交 issue 和 pull request。

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。