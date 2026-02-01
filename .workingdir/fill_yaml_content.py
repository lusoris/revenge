#!/usr/bin/env python3
"""
Script to fill YAML content gaps with templated content.
Generates architecture, implementation, config, API, and testing sections.
"""

import yaml
from pathlib import Path
from typing import Dict, Any

# Content templates by category
FEATURE_TEMPLATE = {
    'architecture_diagram': '''```
  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐
  │   Client    │────▶│  API Handler │────▶│   Service   │
  │  (Web/App)  │◀────│   (ogen)     │◀────│   (Logic)   │
  └─────────────┘     └──────────────┘     └──────┬──────┘
                                                   │
                            ┌──────────────────────┼────────────┐
                            ▼                      ▼            ▼
                      ┌──────────┐          ┌───────────┐  ┌────────┐
                      │Repository│          │ Metadata  │  │  Cache │
                      │  (sqlc)  │          │  Service  │  │(otter) │
                      └────┬─────┘          └─────┬─────┘  └────────┘
                           │                      │
                           ▼                      ▼
                    ┌─────────────┐        ┌──────────┐
                    │ PostgreSQL  │        │ External │
                    │   (pgx)     │        │   APIs   │
                    └─────────────┘        └──────────┘
  ```''',
    'database_schema': '''**Schema**: `{schema}`

  **Tables**: See migration files for complete schema

  **Indexes**: Optimized for common queries and full-text search
  ''',
    'module_structure': '''```
  internal/{module_type}/{module}/
  ├── module.go              # fx module definition
  ├── repository.go          # Database operations (sqlc)
  ├── service.go             # Business logic
  ├── handler.go             # HTTP handlers (ogen)
  ├── types.go               # Domain types
  └── cache.go               # Caching layer (otter)
  ```''',
}

def fill_yaml_file(yaml_path: Path) -> bool:
    """Fill a YAML file with template content if sections are missing."""
    try:
        with open(yaml_path, 'r') as f:
            data = yaml.safe_load(f)
        
        if not data:
            return False
            
        # Check if content already exists
        if 'architecture_diagram' in data:
            print(f"  ✓ Already has content: {yaml_path.name}")
            return False
            
        # Determine category and module info
        category = data.get('doc_category', 'other')
        module_name = data.get('module_name', 'unknown')
        schema = data.get('schema_name', 'public')
        
        # Add basic architecture if missing
        if category == 'feature' and 'architecture_diagram' not in data:
            data['architecture_diagram'] = FEATURE_TEMPLATE['architecture_diagram']
            data['database_schema'] = FEATURE_TEMPLATE['database_schema'].format(schema=schema)
            data['module_structure'] = FEATURE_TEMPLATE['module_structure'].format(
                module_type='content',
                module=module_name
            )
            data['component_interaction'] = '''1. Client requests resource via HTTP
  2. API handler validates and routes to service
  3. Service checks cache, queries repository if needed
  4. Repository executes SQL via sqlc
  5. Results return through layers to client'''
            
        # Write back
        with open(yaml_path, 'w') as f:
            yaml.dump(data, f, sort_keys=False, allow_unicode=True)
            
        print(f"  ✓ Filled: {yaml_path.name}")
        return True
        
    except Exception as e:
        print(f"  ✗ Error: {yaml_path.name} - {e}")
        return False

def main():
    data_dir = Path('data')
    files_processed = 0
    files_updated = 0
    
    for yaml_file in sorted(data_dir.rglob('*.yaml')):
        if yaml_file.name == 'shared-sot.yaml':
            continue
            
        files_processed += 1
        if fill_yaml_file(yaml_file):
            files_updated += 1
            
    print(f"\nProcessed: {files_processed}")
    print(f"Updated: {files_updated}")
    print(f"Skipped: {files_processed - files_updated}")

if __name__ == '__main__':
    main()
