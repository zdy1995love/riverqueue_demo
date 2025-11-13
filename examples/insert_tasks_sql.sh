#!/bin/bash

# River Queue Task Insertion Example (using psql)

echo "🚀 River Queue Task Insertion Example (SQL)"
echo "================================"

# Read database connection info from config file
# Note: You need to set this manually or parse config_DEV.jsonc using jq

# Example configuration (modify according to your setup)
# If you used pg.sh, use these values:
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="riverqueue"
DB_PASS="riverqueue_password"
DB_NAME="river_demo"

echo "💡 Tip: Please ensure the following information is correct:"
echo "   Host: $DB_HOST:$DB_PORT"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo ""

# Check if running in Docker
if command -v docker &> /dev/null; then
    echo "🐳 Docker detected, searching for PostgreSQL container..."
    PG_CONTAINER=$(docker ps --format '{{.Names}}' | grep -E "riverqueue-postgres|postgres" | head -1)
    
    if [ -n "$PG_CONTAINER" ]; then
        echo "✅ Found container: $PG_CONTAINER"
        echo ""
        echo "📝 Inserting test tasks..."
        
        # Execute SQL using Docker exec
        docker exec -i "$PG_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Task 1: AddOne(5)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_one', '{"number": 5}', 'default', 'available', 1, 3, NOW(), NOW())
RETURNING id, kind, args;

-- Task 2: MultiplyTwo(10) - will auto-trigger AddThree(20)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('multiply_two', '{"number": 10}', 'default', 'available', 1, 3, NOW(), NOW())
RETURNING id, kind, args;

-- Task 3: AddOne(100)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_one', '{"number": 100}', 'default', 'available', 1, 3, NOW(), NOW())
RETURNING id, kind, args;

-- Task 4: MultiplyTwo(7) - will auto-trigger AddThree(14)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('multiply_two', '{"number": 7}', 'default', 'available', 1, 3, NOW(), NOW())
RETURNING id, kind, args;

-- Task 5: AddThree(50)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_three', '{"number": 50}', 'default', 'available', 1, 3, NOW(), NOW())
RETURNING id, kind, args;
EOF
        
        echo ""
        echo "✅ All tasks inserted successfully"
        echo "💡 Check the River Queue main program output to see task execution"
        
    else
        echo "❌ PostgreSQL container not found"
        echo "💡 If PostgreSQL is not in Docker, use the following command to insert manually:"
        echo ""
        echo "psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f examples/insert_tasks.sql"
    fi
else
    echo "📝 Inserting tasks using local psql..."
    PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Task 1: AddOne(5)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_one', '{"number": 5}', 'default', 'available', 1, 3, NOW(), NOW());

-- Task 2: MultiplyTwo(10)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('multiply_two', '{"number": 10}', 'default', 'available', 1, 3, NOW(), NOW());

-- Task 3: AddOne(100)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_one', '{"number": 100}', 'default', 'available', 1, 3, NOW(), NOW());

-- Task 4: MultiplyTwo(7)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('multiply_two', '{"number": 7}', 'default', 'available', 1, 3, NOW(), NOW());

-- Task 5: AddThree(50)
INSERT INTO river_job (kind, args, queue, state, priority, max_attempts, scheduled_at, created_at)
VALUES ('add_three', '{"number": 50}', 'default', 'available', 1, 3, NOW(), NOW());
EOF
    
    echo ""
    echo "✅ All tasks inserted successfully"
fi
