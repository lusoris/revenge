# Adult Content Documentation

> Adult module with isolated `c` schema (Queen Anne's Revenge)

---

## Documents

| Document | Description |
|----------|-------------|
| [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md) | Complete adult module isolation design |
| [ADULT_METADATA.md](ADULT_METADATA.md) | Whisparr/StashDB metadata integration |
| [DATA_RECONCILIATION.md](DATA_RECONCILIATION.md) | Data matching and deduplication |
| [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) | Database schema for adult content |

---

## Key Concepts

- **Schema Isolation**: All adult content in `qar` schema (Queen Anne's Revenge obfuscation) - see [SOURCE_OF_TRUTH.md](../../SOURCE_OF_TRUTH.md#qar-obfuscation-terminology)
- **Metadata**: Whisparr as primary, StashDB as fallback
- **Access Control**: Requires explicit RBAC permission

---

## Related

- [Integrations: Whisparr](../../integrations/servarr/WHISPARR.md)
- [Integrations: StashDB](../../integrations/metadata/adult/STASHDB.md)
- [Shared: NSFW Toggle](../shared/NSFW_TOGGLE.md)
