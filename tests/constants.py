"""Shared constants for test configuration.

Centralized thresholds and configuration values used across test modules.
Update these values as documentation improves.
"""

# =============================================================================
# PLACEHOLDER Thresholds
# =============================================================================
# These thresholds allow gradual PLACEHOLDER replacement while preventing
# new PLACEHOLDERs from being added. Decrease as docs are completed.
# Set to 0 for strict mode when all docs should be complete.

# Current counts - decrease as docs complete
MAX_WIKI_OVERVIEW_PLACEHOLDERS = 113
MAX_WIKI_TAGLINE_PLACEHOLDERS = 11

# =============================================================================
# Required Fields
# =============================================================================

BASE_REQUIRED_FIELDS = [
    "doc_title",
    "doc_category",
    "created_date",
    "overall_status",
    "status_design",
    "technical_summary",
]

FEATURE_REQUIRED_FIELDS = [
    "feature_name",
    "schema_name",
]

SERVICE_REQUIRED_FIELDS = [
    "service_name",
    "package_path",
    "fx_module",
]

INTEGRATION_REQUIRED_FIELDS = [
    "integration_name",
    "external_service",
    "integration_id",
]

# =============================================================================
# Generation Thresholds
# =============================================================================

# Minimum expected successful doc generations
MIN_SUCCESSFUL_GENERATIONS = 140

# Maximum allowed generation failures
MAX_GENERATION_FAILURES = 5

# Minimum expected YAML files in data directory
MIN_YAML_FILES = 100
