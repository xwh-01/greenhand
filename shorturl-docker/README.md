# Day 8: 短链接引擎 V2.0 (MySQL版)

## 功能特性
- 基于 MySQL 数据持久化
- 支持创建短链接
- 支持短链接跳转  
- 支持查看统计数据
- 服务重启数据不丢失

## 技术栈
- Go 语言
- GORM ORM 框架
- MySQL 数据库
- 分层架构设计

## 启动前准备
1. 确保 MySQL 服务运行：`net start MySQL`
2. 创建数据库：`short_url`

## 启动方式
```bash
go run main.go