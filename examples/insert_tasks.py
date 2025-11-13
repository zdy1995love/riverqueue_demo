#!/usr/bin/env python3
"""
River Queue Task Insertion Example (Python)

This script demonstrates how to insert River tasks into PostgreSQL database using Python.
Supports two methods:
1. Official riverqueue-python client (Recommended)
2. Direct SQL insertion (Raw method)
"""

import json
from dataclasses import dataclass
from datetime import datetime, timezone

# Optional imports - depends on which method you use
try:
    import riverqueue
    import sqlalchemy
    from riverqueue.driver import riversqlalchemy
    HAS_RIVERQUEUE = True
except ImportError:
    HAS_RIVERQUEUE = False

try:
    import psycopg2
    HAS_PSYCOPG2 = True
except ImportError:
    HAS_PSYCOPG2 = False

# ==================== Job Args Definitions ====================
# Define Job Args according to riverqueue-python requirements

@dataclass
class AddOneArgs:
    """Arguments for AddOne Worker"""
    number: int
    
    kind: str = "add_one"
    
    def to_json(self) -> str:
        return json.dumps({"number": self.number})


@dataclass
class MultiplyTwoArgs:
    """Arguments for MultiplyTwo Worker"""
    number: int
    
    kind: str = "multiply_two"
    
    def to_json(self) -> str:
        return json.dumps({"number": self.number})


@dataclass
class AddThreeArgs:
    """Arguments for AddThree Worker"""
    number: int
    
    kind: str = "add_three"
    
    def to_json(self) -> str:
        return json.dumps({"number": self.number})


# ==================== Configuration and Utility Functions ====================

def load_config():
    """Load configuration file"""
    with open('setting/config_DEV.jsonc', 'r') as f:
        config = json.load(f)
    return config

def parse_database_url(url):
    """Parse database URL (for psycopg2)"""
    from urllib.parse import urlparse, unquote
    
    # Parse URL
    parsed = urlparse(url)
    
    # Extract dbname without parameters
    dbname = parsed.path.lstrip('/')
    if '?' in dbname:
        dbname = dbname.split('?')[0]
    
    return {
        'host': parsed.hostname,
        'port': parsed.port,
        'user': parsed.username,
        'password': unquote(parsed.password) if parsed.password else None,
        'dbname': dbname
    }


def convert_database_url_for_sqlalchemy(url):
    """Convert database URL to SQLAlchemy format and extract search_path"""
    # SQLAlchemy needs postgresql:// instead of postgres://
    url = url.replace('postgres://', 'postgresql://')
    
    # Extract and remove search_path parameter (psycopg2 doesn't support this parameter)
    search_path = None
    if '?' in url and 'search_path=' in url:
        base_url, params = url.split('?', 1)
        param_list = params.split('&')
        new_params = []
        
        for param in param_list:
            if param.startswith('search_path='):
                search_path = param.split('=', 1)[1]
            else:
                new_params.append(param)
        
        if new_params:
            url = base_url + '?' + '&'.join(new_params)
        else:
            url = base_url
    
    return url, search_path


# ==================== Official Client Insertion Method ====================

def insert_tasks_with_official_client(database_url):
    """Insert tasks using official riverqueue-python client (Recommended)"""
    print("\n🎯 Method 1: Using Official riverqueue-python Client")
    print("=" * 50)
    
    if not HAS_RIVERQUEUE:
        print("❌ Error: riverqueue package not installed")
        print("Please run: pip install riverqueue sqlalchemy")
        return
    
    try:
        # Convert URL format and extract search_path
        sqlalchemy_url, search_path = convert_database_url_for_sqlalchemy(database_url)
        
        # Create SQLAlchemy engine
        # If search_path exists, set it using connect_args
        connect_args = {}
        if search_path:
            connect_args['options'] = f'-csearch_path={search_path}'
            print(f"📌 Using schema: {search_path}")
        
        engine = sqlalchemy.create_engine(sqlalchemy_url, connect_args=connect_args)
        
        # Create River client
        client = riverqueue.Client(riversqlalchemy.Driver(engine))
        
        print("✅ River client created successfully")
        print("\n📝 Inserting test tasks...")
        
        # Task 1: AddOne(5)
        insert_res = client.insert(AddOneArgs(number=5))
        print(f"✅ Task inserted [ID: {insert_res.job.id}] AddOne(5)")
        
        # Task 2: MultiplyTwo(10)
        insert_res = client.insert(MultiplyTwoArgs(number=10))
        print(f"✅ Task inserted [ID: {insert_res.job.id}] MultiplyTwo(10) → will trigger AddThree(20)")
        
        # Task 3: AddOne(100)
        insert_res = client.insert(AddOneArgs(number=100))
        print(f"✅ Task inserted [ID: {insert_res.job.id}] AddOne(100)")
        
        # Task 4: MultiplyTwo(7)
        insert_res = client.insert(MultiplyTwoArgs(number=7))
        print(f"✅ Task inserted [ID: {insert_res.job.id}] MultiplyTwo(7) → will trigger AddThree(14)")
        
        # Task 5: AddThree(50)
        insert_res = client.insert(AddThreeArgs(number=50))
        print(f"✅ Task inserted [ID: {insert_res.job.id}] AddThree(50)")
        
        # Batch insert example
        print("\n📦 Batch inserting tasks...")
        results = client.insert_many([
            AddOneArgs(number=200),
            AddOneArgs(number=300),
        ])
        print(f"✅ Batch insert successful ({len(results)} tasks)")
        for result in results:
            print(f"   - [ID: {result.job.id}] {result.job.kind}")
        
        print("\n✅ Official client insertion complete")
        
        # Cleanup
        engine.dispose()
        
    except Exception as e:
        print(f"❌ Official client insertion failed: {e}")
        import traceback
        traceback.print_exc()


# ==================== Direct SQL Insertion Method ====================

def insert_task(cursor, kind, args, queue='default'):
    """Insert River task into database (Direct SQL method)"""
    sql = """
    INSERT INTO river_job (
        kind,
        args,
        queue,
        priority,
        max_attempts,
        state,
        scheduled_at,
        created_at
    ) VALUES (
        %s,
        %s,
        %s,
        1,
        3,
        'available',
        NOW(),
        NOW()
    ) RETURNING id
    """
    
    cursor.execute(sql, (kind, json.dumps(args), queue))
    job_id = cursor.fetchone()[0]
    return job_id


def insert_tasks_with_sql(database_url):
    """Insert tasks using raw SQL method"""
    print("\n🎯 Method 2: Using Raw SQL Direct Insertion")
    print("=" * 50)
    
    if not HAS_PSYCOPG2:
        print("❌ Error: psycopg2 package not installed")
        print("Please run: pip install psycopg2-binary")
        return
    
    # Parse database configuration
    db_config = parse_database_url(database_url)
    
    # Extract search_path from URL
    from urllib.parse import urlparse, parse_qs
    parsed = urlparse(database_url)
    search_path = None
    if parsed.query:
        params = parse_qs(parsed.query)
        if 'search_path' in params:
            search_path = params['search_path'][0]
    
    # Connect to database
    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()
        
        # Set search_path if specified
        if search_path:
            cursor.execute(f"SET search_path TO {search_path}")
            print(f"✅ Database connected successfully (schema: {search_path})")
        else:
            print("✅ Database connected successfully")
    except Exception as e:
        print(f"❌ Failed to connect to database: {e}")
        return
    
    print("\n📝 Inserting test tasks...")
    
    try:
        # Task 1: AddOne(5)
        job_id = insert_task(cursor, 'add_one', {'number': 5})
        print(f"✅ Task inserted [ID: {job_id}] AddOne(5)")
        
        # Task 2: MultiplyTwo(10)
        job_id = insert_task(cursor, 'multiply_two', {'number': 10})
        print(f"✅ Task inserted [ID: {job_id}] MultiplyTwo(10) → will trigger AddThree(20)")
        
        # Task 3: AddOne(100)
        job_id = insert_task(cursor, 'add_one', {'number': 100})
        print(f"✅ Task inserted [ID: {job_id}] AddOne(100)")
        
        # Task 4: MultiplyTwo(7)
        job_id = insert_task(cursor, 'multiply_two', {'number': 7})
        print(f"✅ Task inserted [ID: {job_id}] MultiplyTwo(7) → will trigger AddThree(14)")
        
        # Task 5: AddThree(50)
        job_id = insert_task(cursor, 'add_three', {'number': 50})
        print(f"✅ Task inserted [ID: {job_id}] AddThree(50)")
        
        # Commit transaction
        conn.commit()
        
        print("\n✅ SQL direct insertion complete")
        
    except Exception as e:
        conn.rollback()
        print(f"❌ Failed to insert tasks: {e}")
    finally:
        cursor.close()
        conn.close()


def main():
    print("🚀 River Queue Task Insertion Example (Python)")
    print("=" * 60)
    print("📚 This example supports two insertion methods:")
    print("   1. Official riverqueue-python client (Recommended)")
    print("   2. Raw SQL direct insertion")
    print("=" * 60)
    
    # Load configuration
    config = load_config()
    database_url = config['river_database_url']
    
    # Method selection
    print("\nPlease select insertion method:")
    print("  [1] Use official riverqueue-python client (Recommended)")
    print("  [2] Use raw SQL direct insertion")
    print("  [3] Execute both methods")
    
    try:
        choice = input("\nEnter your choice (1/2/3, default is 1): ").strip() or "1"
    except (EOFError, KeyboardInterrupt):
        choice = "1"
        print("1")
    
    if choice == "1":
        insert_tasks_with_official_client(database_url)
    elif choice == "2":
        insert_tasks_with_sql(database_url)
    elif choice == "3":
        insert_tasks_with_official_client(database_url)
        insert_tasks_with_sql(database_url)
    else:
        print("❌ Invalid choice, using default method (official client)")
        insert_tasks_with_official_client(database_url)
    
    print("\n" + "=" * 60)
    print("✅ Task insertion complete")
    print("💡 Check the River Queue main program output to see task execution")
    print("=" * 60)

if __name__ == '__main__':
    main()
