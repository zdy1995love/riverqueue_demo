# 🎯 River Queue 演示项目

> 基于 PostgreSQL 和 River Queue 的任务队列完整演示项目

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![River](https://img.shields.io/badge/River-v0.11.4-blue)](https://riverqueue.com/)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?style=flat&logo=python)](https://python.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

中文文档 | [English](README.md)

## 📖 简介

这是一个完整的 River Queue 学习项目，展示了如何使用 River 构建可靠的任务队列系统。该项目包含：

- **Go Workers**: 实现任务处理的 Worker 模块
- **Python Client**: 使用官方 riverqueue-python 客户端插入任务
- **多种插入方式**: 支持 Go、Python 和 SQL 直接插入
- **链式任务**: 演示任务间的依赖关系

**核心特性:**
- ✅ 模块化的 Worker 实现 (AddOne, MultiplyTwo, AddThree)
- ✅ 持续运行模式与优雅关闭
- ✅ 链式任务支持 (MultiplyTwo → AddThree)
- ✅ 多种任务插入方式 (Go/Python/SQL)
- ✅ 官方 Python 客户端集成
- ✅ 完整的错误处理和日志记录

## 🚀 快速开始

### 环境准备

**Go 环境:**
```bash
# 安装 Go 1.21+
go version
```

**PostgreSQL:**
```bash
# 方式1: 使用 Docker (推荐)
chmod +x pg.sh
./pg.sh

# 方式2: 使用本地 PostgreSQL
psql --version
createdb river_demo
```

**Python 环境 (可选，用于任务插入):**
```bash
# 方式1: 使用 conda (推荐用于 riverqueue-python)
conda create -n riverqueue python=3.12
conda activate riverqueue

# 克隆并从源码安装 riverqueue-python
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python
pip install .

# 方式2: 使用 venv
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\Scripts\activate  # Windows

# 从 PyPI 安装 (如果可用)
pip install riverqueue

# 或安装直接 SQL 插入所需的依赖
pip install psycopg2-binary sqlalchemy
```

### 配置

1. 复制并编辑配置文件:

```bash
cp setting/config_DEV.jsonc.example setting/config_DEV.jsonc
```

2. 更新 `setting/config_DEV.jsonc` 中的数据库凭证:

```json
{
    "river_database_url": "postgres://riverqueue:riverqueue_password@localhost:5432/river_demo?search_path=riverqueue",
    "river_max_workers": 5,
    "river_test_only": false
}
```

**配置参数说明:**
- `river_database_url`: PostgreSQL 连接字符串（如果使用 `pg.sh`，请使用脚本输出的 URL）
- `river_max_workers`: 最大并发 worker 数量
- `river_test_only`: 设为 `true` 用于测试（任务立即执行）

### 运行应用

**终端 1 - 启动 River Queue Worker:**
```bash
# 添加执行权限
chmod +x run.sh

# 启动 worker
./run.sh
```

**终端 2 - 插入任务:**

使用 Python (推荐):
```bash
chmod +x run_python.sh
./run_python.sh
```

使用 Go:
```bash
chmod +x run_go.sh
./run_go.sh
```

使用直接 SQL:
```bash
bash examples/insert_tasks_sql.sh
```

## 📁 项目结构

```
riverqueue_demo/
├── main.go                 # 主程序入口
├── go.mod                  # Go 模块依赖
├── go.sum                  # 依赖校验和
├── run.sh                  # 主运行脚本
├── run_go.sh              # Go 任务插入脚本
├── run_python.sh          # Python 任务插入脚本
├── setting/
│   └── config_DEV.jsonc   # 配置文件
├── worker/                # Worker 实现
│   ├── addone/           # AddOne worker (+1)
│   ├── multiplytwo/      # MultiplyTwo worker (×2, 链接到 AddThree)
│   └── addthree/         # AddThree worker (+3)
└── examples/             # 任务插入示例
    ├── insert_tasks.go   # Go 客户端示例
    ├── insert_tasks.py   # Python 客户端示例
    └── insert_tasks_sql.sh # 直接 SQL 插入
```

## 🔧 Workers 说明

### 1. AddOne Worker
- **类型**: `add_one`
- **功能**: 给输入数字加 1
- **参数**: `{"number": N}`

### 2. MultiplyTwo Worker
- **类型**: `multiply_two`
- **功能**: 将输入数字乘以 2，然后链接到 AddThree
- **参数**: `{"number": N}`
- **链式**: 自动创建 AddThree 任务处理结果

### 3. AddThree Worker
- **类型**: `add_three`
- **功能**: 给输入数字加 3
- **参数**: `{"number": N}`

## 📝 使用示例

### 示例 1: 简单任务
```bash
# 插入 AddOne 任务，数字为 5
# 预期输出: 5 + 1 = 6
```

### 示例 2: 链式任务
```bash
# 插入 MultiplyTwo 任务，数字为 5
# 步骤 1: 5 × 2 = 10
# 步骤 2: 10 + 3 = 13 (自动链接的 AddThree)
```

## 🐍 Python 环境设置

### 安装 riverqueue-python

官方 `riverqueue-python` 包目前正在开发中。使用方法:

**方式 1: 从源码安装 (推荐)**
```bash
# 创建 conda 环境
conda create -n riverqueue python=3.12
conda activate riverqueue

# 克隆仓库
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python

# 以开发模式安装
pip install -e .

# 或直接安装
pip install .
```

**方式 2: 使用 venv**
```bash
# 创建虚拟环境
python3 -m venv venv
source venv/bin/activate

# 克隆并安装
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python
pip install .
```

**依赖包:**
- `riverqueue`: River Queue 官方 Python 客户端
- `sqlalchemy`: 数据库 ORM
- `psycopg2-binary`: PostgreSQL 适配器

## 🔍 监控

在 PostgreSQL 中查看任务:
```sql
-- 查看所有任务
SELECT * FROM riverqueue.river_job ORDER BY created_at DESC LIMIT 10;

-- 按状态查看任务
SELECT state, kind, COUNT(*) 
FROM riverqueue.river_job 
GROUP BY state, kind;

-- 查看最近完成的任务
SELECT id, kind, args, state, created_at, finalized_at 
FROM riverqueue.river_job 
WHERE state = 'completed' 
ORDER BY finalized_at DESC 
LIMIT 10;
```

## 🛠️ 开发

### 构建
```bash
go build -o river-queue-demo main.go
./river-queue-demo
```

### 测试
```bash
go test -v ./...
```

### 依赖管理
```bash
# 下载依赖
go mod download

# 更新依赖
go mod tidy
```

## � Docker PostgreSQL 管理

`pg.sh` 脚本提供了简便的 PostgreSQL 管理:

```bash
# 启动 PostgreSQL
./pg.sh

# 检查状态
docker ps

# 停止容器
docker stop riverqueue-postgres

# 启动已存在的容器
docker start riverqueue-postgres

# 删除容器 (⚠️ 会删除数据)
docker rm -f riverqueue-postgres
```

## �📚 了解更多

- [River Queue 官方文档](https://riverqueue.com/)
- [River Go 客户端](https://github.com/riverqueue/river)
- [River Python 客户端](https://github.com/riverqueue/riverqueue-python)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [River Queue](https://riverqueue.com/) - 优秀的任务队列系统
- [PostgreSQL](https://www.postgresql.org/) - 可靠的数据库系统

---

**祝您使用愉快! 🎉**
