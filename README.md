# 🎯 River Queue Demo

> A comprehensive demo project showcasing task queue management with PostgreSQL and River Queue

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![River](https://img.shields.io/badge/River-v0.11.4-blue)](https://riverqueue.com/)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?style=flat&logo=python)](https://python.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[中文文档](README_zh.md) | English

## 📖 Introduction

This is a complete River Queue learning project demonstrating how to build a reliable task queue system using River. The project includes:

- **Go Workers**: Worker modules for task processing
- **Python Client**: Task insertion using official riverqueue-python client
- **Multiple Insertion Methods**: Support for Go, Python, and direct SQL insertion
- **Task Chaining**: Demonstration of task dependencies

**Key Features:**
- ✅ Modular Worker implementation (AddOne, MultiplyTwo, AddThree)
- ✅ Continuous running mode with graceful shutdown
- ✅ Task chaining support (MultiplyTwo → AddThree)
- ✅ Multiple task insertion methods (Go/Python/SQL)
- ✅ Official Python client integration
- ✅ Complete error handling and logging

## 🚀 Quick Start

### Prerequisites

**Go Environment:**
```bash
# Install Go 1.21+
go version
```

**PostgreSQL:**
```bash
# Option 1: Using Docker (Recommended)
chmod +x pg.sh
./pg.sh

# Option 2: Using local PostgreSQL
psql --version
createdb river_demo
```

**Python Environment (Optional, for task insertion):**
```bash
# Option 1: Using conda (Recommended for riverqueue-python)
conda create -n riverqueue python=3.12
conda activate riverqueue

# Clone and install riverqueue-python from source
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python
pip install .

# Option 2: Using venv
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# or venv\Scripts\activate  # Windows

# Install from PyPI (if available)
pip install riverqueue

# Or install dependencies for direct SQL insertion
pip install psycopg2-binary sqlalchemy
```

### Configuration

1. Copy and edit the configuration file:

```bash
cp setting/config_DEV.jsonc.example setting/config_DEV.jsonc
```

2. Update `setting/config_DEV.jsonc` with your database credentials:

```json
{
    "river_database_url": "postgres://riverqueue:riverqueue_password@localhost:5432/river_demo?search_path=riverqueue",
    "river_max_workers": 5,
    "river_test_only": false
}
```

**Configuration Parameters:**
- `river_database_url`: PostgreSQL connection string (if you used `pg.sh`, use the URL printed by the script)
- `river_max_workers`: Maximum number of concurrent workers
- `river_test_only`: Set to `true` for testing (tasks execute immediately)

### Running the Application

**Terminal 1 - Start River Queue Worker:**
```bash
# Make script executable
chmod +x run.sh

# Start the worker
./run.sh
```

**Terminal 2 - Insert Tasks:**

Using Python (Recommended):
```bash
chmod +x run_python.sh
./run_python.sh
```

Using Go:
```bash
chmod +x run_go.sh
./run_go.sh
```

Using direct SQL:
```bash
bash examples/insert_tasks_sql.sh
```

## 📁 Project Structure

```
riverqueue_demo/
├── main.go                 # Main application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── run.sh                  # Main runner script
├── run_go.sh              # Go task inserter script
├── run_python.sh          # Python task inserter script
├── setting/
│   └── config_DEV.jsonc   # Configuration file
├── worker/                # Worker implementations
│   ├── addone/           # AddOne worker (+1)
│   ├── multiplytwo/      # MultiplyTwo worker (×2, chains to AddThree)
│   └── addthree/         # AddThree worker (+3)
└── examples/             # Task insertion examples
    ├── insert_tasks.go   # Go client example
    ├── insert_tasks.py   # Python client example
    └── insert_tasks_sql.sh # Direct SQL insertion
```

## 🔧 Workers

### 1. AddOne Worker
- **Kind**: `add_one`
- **Function**: Adds 1 to the input number
- **Args**: `{"number": N}`

### 2. MultiplyTwo Worker
- **Kind**: `multiply_two`
- **Function**: Multiplies input by 2, then chains to AddThree
- **Args**: `{"number": N}`
- **Chaining**: Automatically creates an AddThree task with the result

### 3. AddThree Worker
- **Kind**: `add_three`
- **Function**: Adds 3 to the input number
- **Args**: `{"number": N}`

## 📝 Usage Examples

### Example 1: Simple Task
```bash
# Insert AddOne task with number=5
# Expected output: 5 + 1 = 6
```

### Example 2: Chained Task
```bash
# Insert MultiplyTwo task with number=5
# Step 1: 5 × 2 = 10
# Step 2: 10 + 3 = 13 (auto-chained AddThree)
```

## 🐍 Python Environment Setup

### Installing riverqueue-python

The official `riverqueue-python` package is currently under development. To use it:

**Option 1: Install from source (Recommended)**
```bash
# Create conda environment
conda create -n riverqueue python=3.12
conda activate riverqueue

# Clone the repository
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python

# Install in development mode
pip install -e .

# Or install directly
pip install .
```

**Option 2: Using venv**
```bash
# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Clone and install
git clone https://github.com/riverqueue/riverqueue-python.git
cd riverqueue-python
pip install .
```

**Dependencies:**
- `riverqueue`: Official River Queue Python client
- `sqlalchemy`: Database ORM
- `psycopg2-binary`: PostgreSQL adapter

## 🔍 Monitoring

View tasks in PostgreSQL:
```sql
-- View all tasks
SELECT * FROM riverqueue.river_job ORDER BY created_at DESC LIMIT 10;

-- View tasks by status
SELECT state, kind, COUNT(*) 
FROM riverqueue.river_job 
GROUP BY state, kind;

-- View recent completed tasks
SELECT id, kind, args, state, created_at, finalized_at 
FROM riverqueue.river_job 
WHERE state = 'completed' 
ORDER BY finalized_at DESC 
LIMIT 10;
```

## 🛠️ Development

### Building
```bash
go build -o river-queue-demo main.go
./river-queue-demo
```

### Testing
```bash
go test -v ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Update dependencies
go mod tidy
```

## � Docker PostgreSQL Management

The `pg.sh` script provides easy PostgreSQL management:

```bash
# Start PostgreSQL
./pg.sh

# Check status
docker ps

# Stop container
docker stop riverqueue-postgres

# Start existing container
docker start riverqueue-postgres

# Remove container (⚠️ deletes data)
docker rm -f riverqueue-postgres
```

## �📚 Learn More

- [River Queue Official Documentation](https://riverqueue.com/)
- [River Go Client](https://github.com/riverqueue/river)
- [River Python Client](https://github.com/riverqueue/riverqueue-python)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [River Queue](https://riverqueue.com/) - The awesome job queue system
- [PostgreSQL](https://www.postgresql.org/) - Reliable database system

---

**Happy Queueing! 🎉**
