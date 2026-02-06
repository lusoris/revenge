# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-02-06)


### ⚠ BREAKING CHANGES

* **library:** Shared libraries table no longer exists. Each content module now manages its own library table.

### Features

* add Helm chart, develop auto-builds, and production compose ([fa7389c](https://github.com/lusoris/revenge/commit/fa7389c22196cef09e1fce65f24ef4d7f7b6268b))
* add minimal Go skeleton for CI/CD ([3280748](https://github.com/lusoris/revenge/commit/3280748575a5b2eff396554b3d0b695190e686dd))
* Add Movie Module backend foundation ([59fb5d1](https://github.com/lusoris/revenge/commit/59fb5d13501e9b7011a8bc3bc78d54ede6de8565))
* add strict mode to sync-versions.py for CI validation ([f2c2aeb](https://github.com/lusoris/revenge/commit/f2c2aeb717cf21b6155e5507263bb4c43b9edbe4))
* add sync-sot-status.py script to sync YAML status to SOT tables ([bff0fc6](https://github.com/lusoris/revenge/commit/bff0fc648190e1644811aa5c2165ff16af6db864))
* **adult:** add QAR obfuscation + trans performer fields ([030a634](https://github.com/lusoris/revenge/commit/030a634c04658f571797fe8d03064019dd66e349))
* **api:** add adult content access control infrastructure ([a59678b](https://github.com/lusoris/revenge/commit/a59678b60053a7702efd2738733b70cb56abc6ba))
* **api:** add full Sonarr API handlers matching Radarr pattern ([a2508c9](https://github.com/lusoris/revenge/commit/a2508c9cb5eb84aa31b24bf1cf59a6ae49564d59))
* **api:** add metadata and image proxy endpoints ([31b596a](https://github.com/lusoris/revenge/commit/31b596a9cf76fb9030c7c5b229c2ab4e8a8489a4))
* **api:** add multi-language localization support ([b1dd14e](https://github.com/lusoris/revenge/commit/b1dd14e880cbfc49589a0a58a496bd09dc3d7e28))
* **api:** add QAR OpenAPI spec for adult content endpoints ([0e238da](https://github.com/lusoris/revenge/commit/0e238daedf307466d636a22423849d858d099be4))
* **api:** add rate limiting middleware and nice-to-have endpoints ([9414a61](https://github.com/lusoris/revenge/commit/9414a61a1c1e0c0f4e01ce42df2875f0af3eac0d))
* **api:** add settings API endpoints and handlers ([b1098a0](https://github.com/lusoris/revenge/commit/b1098a06b8886422cd1186ae3ca332f7b10e2353))
* **api:** add TV metadata endpoints to complete metadata API ([2353e22](https://github.com/lusoris/revenge/commit/2353e229b6f725972b72cbf43b5317f715361f42))
* **api:** add TV show API handlers and endpoints (A11.8) ([0008cdb](https://github.com/lusoris/revenge/commit/0008cdbf0c8d8438dfb295ac4ca055b0584cc3cc))
* **api:** complete Phase A12 shared metadata integration ([34ddad9](https://github.com/lusoris/revenge/commit/34ddad9328bf979d7289908214a778eba6d05612))
* **api:** Day 4-5 - OpenAPI foundation with ogen handlers ([02c15ca](https://github.com/lusoris/revenge/commit/02c15ca170c4f8c025f3198c81cee1650083fcd6))
* **api:** implement Session and RBAC API endpoints (Steps 7+8) ([69597be](https://github.com/lusoris/revenge/commit/69597bed29af1af2909b5293720524928990da34))
* **api:** implement TV Shows API handlers ([54938f0](https://github.com/lusoris/revenge/commit/54938f0d52f272a0d6dcc1b39025849bb4425966))
* **api:** implement User Service API handlers ([715045d](https://github.com/lusoris/revenge/commit/715045de6bf234d5bbe35aa2dcabd37637c6dfd0))
* **apikeys:** implement API Keys service (Step 9) ([ae0ad75](https://github.com/lusoris/revenge/commit/ae0ad75a2be302ff58cc050b9495a75713bf3bc4))
* **api:** wire TVShow module and fix graceful shutdown ([13f0984](https://github.com/lusoris/revenge/commit/13f0984a2b6d6670ce5607ea34db64cf5e596d43))
* **auth:** add auth token database tables ([7aab4eb](https://github.com/lusoris/revenge/commit/7aab4eb01d2c69b6f15b1770999e0d831589d584))
* **auth:** Add repository layer for auth service ([9dd5e0c](https://github.com/lusoris/revenge/commit/9dd5e0cf6451a85891f53afe0749b1866f5e31af))
* **auth:** implement account lockout / rate limiting (A7.5) ([9f57ecf](https://github.com/lusoris/revenge/commit/9f57ecf4d5ea7196f6ebf6dffcc6b9c6fd401f96))
* **auth:** implement auth API endpoints (Step 6.6) ([506fe1a](https://github.com/lusoris/revenge/commit/506fe1af9eee9084e14bbbc217304f875173f813))
* **auth:** Implement auth service layer with full auth flows ([eda25d7](https://github.com/lusoris/revenge/commit/eda25d72f88e0994b340d2e59289c8196a489b48))
* **auth:** implement JWT authentication middleware (Step 6.5) ([09e88f6](https://github.com/lusoris/revenge/commit/09e88f601fc63f7787d5c2f8c3833af69874228b))
* **auth:** Implement JWT token manager with stdlib crypto ([77bdbc9](https://github.com/lusoris/revenge/commit/77bdbc9937455fdd529ff43c0fe2b84f2210c338))
* **auth:** implement OIDC new user creation [A0.7] ([143bb85](https://github.com/lusoris/revenge/commit/143bb850f0d50169581101becca928ca2caad2ae))
* **avatar:** implement avatar upload with storage abstraction ([580670a](https://github.com/lusoris/revenge/commit/580670af1c8d8ab0a81415a50244bdc91bcc865f))
* **cache:** add otter L1 in-memory cache ([ca2752d](https://github.com/lusoris/revenge/commit/ca2752dcf000eb34ffc1a2239e45c9a3e13c123f))
* **cache:** add rueidis client integration ([a3bf39d](https://github.com/lusoris/revenge/commit/a3bf39de4870c8fee5ae796225a7b5bc3f1ba67a))
* **cache:** add unified cache operations with L1+L2 integration ([b13f5b1](https://github.com/lusoris/revenge/commit/b13f5b1ac04a2899b46b14f30a1d847af66cf560))
* **ci:** add full Kubernetes deployment testing workflow ([73cfedc](https://github.com/lusoris/revenge/commit/73cfedc1b344fbeb44d157935df5d6f7b7207bb7))
* **ci:** add k3s and Docker Swarm deployment testing workflows ([e973732](https://github.com/lusoris/revenge/commit/e973732720f2661a01fc12293f0ad31a91382005))
* **cluster:** add NFS volume support for shared media storage (A8.1.1) ([ab09fbb](https://github.com/lusoris/revenge/commit/ab09fbb583904711c92e83b494f207854bfd30f1))
* **cluster:** integrate Raft with cleanup jobs for leader-only execution ([1cb0247](https://github.com/lusoris/revenge/commit/1cb0247feb842df570637229d573885b14869caf))
* complete automation foundation (Phases 13-16) ([1f30838](https://github.com/lusoris/revenge/commit/1f30838dd2b410fb668f88b72d312167296c684b))
* complete test pipeline for local and CI ([f1de450](https://github.com/lusoris/revenge/commit/f1de4509ec95e3de8f9705ea7ea16229e7f0a8f2))
* **config:** add activity log retention configuration [A6.4] ([bb95b9b](https://github.com/lusoris/revenge/commit/bb95b9b6b313d3b767b4bca9ee3de788e5e0ac20))
* **content:** add movie module + shared interfaces ([ca3e62f](https://github.com/lusoris/revenge/commit/ca3e62f94a3921cf2c6ff213149a864a3cd78b11))
* **content:** add shared background jobs framework (A10.5) ([88f2f0f](https://github.com/lusoris/revenge/commit/88f2f0fd38ae4f67141d4904c5f6c0170fa2a250))
* **content:** add shared library service framework (A10.4) ([fd8ae30](https://github.com/lusoris/revenge/commit/fd8ae3032d27dcaba2b279b529ce43d603fd15b8))
* **content:** add shared matcher framework with fuzzy matching (A10.2) ([42dd1f0](https://github.com/lusoris/revenge/commit/42dd1f0fb3d3d0d9fae337fdef7881abd91c598b))
* **content:** add shared metadata provider framework (A10.3) ([7178bec](https://github.com/lusoris/revenge/commit/7178bec5e15a0984eaec1476b93d7e0c95455917))
* **content:** add shared scanner framework with movie adapter (A10.1) ([866ffe9](https://github.com/lusoris/revenge/commit/866ffe94ad264ec8bf9efa3d3fc9fcf731904d48))
* create HTTP_CLIENT.yaml pattern for proxy/VPN documentation ([52e2ddd](https://github.com/lusoris/revenge/commit/52e2ddd9ee429355be471ca98e21a21f1c50d080))
* **crypto:** add shared AES-256-GCM encryption service ([782a470](https://github.com/lusoris/revenge/commit/782a470b0d890646cd4f57e613821fa11923f96a))
* **database:** add new migration files ([3174228](https://github.com/lusoris/revenge/commit/3174228b81f7682b9639f68cffae38c032555e01))
* **db:** add movie module database schema with 6 migrations ([a72c8c8](https://github.com/lusoris/revenge/commit/a72c8c877cf6d857cd7626fcd45ee1c88324014e))
* **db:** add Prometheus metrics for PostgreSQL pool ([02278a9](https://github.com/lusoris/revenge/commit/02278a9bdc5bce5d2f93de16c2479b739bbc71e4))
* **db:** add query logging with slow query detection ([ad0cc3f](https://github.com/lusoris/revenge/commit/ad0cc3ff7068a7368a45e878c69c52ca67e494ce))
* **db:** add server settings table and sqlc code generation ([ca45b67](https://github.com/lusoris/revenge/commit/ca45b67abc2ffc512a58f31877e49e181ac54afd))
* **db:** add user_preferences and user_avatars tables (migrations 000006-000007) ([ce4fd85](https://github.com/lusoris/revenge/commit/ce4fd85917956c2444729c7e750c7ded71a1e41b))
* **db:** add user_settings table migration ([d926d04](https://github.com/lusoris/revenge/commit/d926d04f977cc457ed40af2a02b3a7923f904244))
* **deploy:** add Helm chart and Docker deployment infrastructure ([5fd024b](https://github.com/lusoris/revenge/commit/5fd024b3aac0ad68a39014bd69de6fef2968f150))
* **deps:** add media processing packages, link RBAC docs ([a15386c](https://github.com/lusoris/revenge/commit/a15386c87a08e0bd743c3fe70c819c8bc9ecb380))
* **deps:** add v0.1.0 TODO dependencies ([f5b9a33](https://github.com/lusoris/revenge/commit/f5b9a334d03bb128781167c906c6f26654cbb642))
* **diagrams:** convert all ASCII diagrams to Mermaid with LR subgraphs ([d8c276a](https://github.com/lusoris/revenge/commit/d8c276afaeec28d5a32ad9a61d9a4caf01adefdb))
* **docker:** add Typesense Dockerfile and entrypoint script ([c74abcf](https://github.com/lusoris/revenge/commit/c74abcff128ba9d8b6b974e8a59a861098d729a4))
* **docs:** Add horizontal layout support for Mermaid diagrams ([25176c9](https://github.com/lusoris/revenge/commit/25176c959e7793b6574de15c7d3c4a295c415509))
* **docs:** convert ASCII diagrams to Mermaid ([ad1820a](https://github.com/lusoris/revenge/commit/ad1820aa22fef0ca93b893e5b5a35c16b34d6292))
* **docs:** convert ASCII diagrams to Mermaid, fix template whitespace ([bfb206c](https://github.com/lusoris/revenge/commit/bfb206ca0814974555b2a235cd79a83d20849fe7))
* **docs:** Improve ASCII to Mermaid converter with proper box extraction ([7e28102](https://github.com/lusoris/revenge/commit/7e28102ed140bc57d315cd3a67f7443640dc6ead))
* **docs:** Phase 1 - Add basic architecture to all YAML files ([0a0387a](https://github.com/lusoris/revenge/commit/0a0387adac872bae74b30e556dcfc7df6b0c70ab))
* **email:** implement transactional email service ([108c7ff](https://github.com/lusoris/revenge/commit/108c7ff3d41e492b9fa20456a986ba35c7438e79))
* **foundation:** Day 1 complete - session service + cleanup ([7d3692b](https://github.com/lusoris/revenge/commit/7d3692b80b7952443487fe6c287a04056f125d2c))
* generate ogen server code from OpenAPI spec ([6477058](https://github.com/lusoris/revenge/commit/6477058dca83449c389d34d46bc57fb0856f5e2f))
* **health:** implement real health checks for cache, jobs, database ([2114fea](https://github.com/lusoris/revenge/commit/2114feab4b89a0b284a97707d7af2c30f62f4859))
* **i18n:** add multi-language database migration for movies ([4a9cbf5](https://github.com/lusoris/revenge/commit/4a9cbf5ffa166baad8fb54a468ede980c1f76692))
* **i18n:** add multi-language support to Movie domain model ([28a82f4](https://github.com/lusoris/revenge/commit/28a82f4dfdac6576f88f9a6a719d3a84f281466d))
* **image:** add image proxy service with caching ([b245b10](https://github.com/lusoris/revenge/commit/b245b10725341d6142639810ed9c0f7bde332c55))
* implement HTTP server with ogen handlers ([63a5202](https://github.com/lusoris/revenge/commit/63a5202131add198a560d11b402b2f5e8d8095fb))
* **infra:** add config and health checks ([7242d77](https://github.com/lusoris/revenge/commit/7242d77ad2549686ddb35aa46f3b32df51f4222f))
* **infra:** add GitHub Actions/Helm sources and fix deployment workflows ([b41eb3f](https://github.com/lusoris/revenge/commit/b41eb3faa108d81d9692f58fc3385790907be028))
* **infra:** Day 1 foundation - module integration, health checks, migrations ([955bf00](https://github.com/lusoris/revenge/commit/955bf00d82cd7790d0e424388d1354984cf85f70))
* integrate status sync into pipelines and CI ([bc59bd6](https://github.com/lusoris/revenge/commit/bc59bd613d0d8eeef4a260e7a4d0a9a6574104ed))
* **integration:** add Sonarr API v3 client for TV shows (A11.7) ([ed56cee](https://github.com/lusoris/revenge/commit/ed56ceebe04c16bf6abe000fe2bcd0922cf30489))
* **jobs:** add cleanup job worker implementation ([8122511](https://github.com/lusoris/revenge/commit/81225116530ba674b98ee97c3af77f22fe4e8bf9))
* **jobs:** add queue configuration and retry policies ([bcd708d](https://github.com/lusoris/revenge/commit/bcd708d3c78dca4ab91d5eef7691c59c1adff40f))
* **jobs:** add River job queue client integration ([2f6d021](https://github.com/lusoris/revenge/commit/2f6d021a1d919ede55ef7bf2c7d4ce5f205e8ce3))
* **jobs:** implement 5-level queue priority system [A6.3] ([d3262a7](https://github.com/lusoris/revenge/commit/d3262a7c3dbb4671d3a8da0fdacda69ac4ceca13))
* **jobs:** Week 2 Day 1 - River workers infrastructure ([8669bb7](https://github.com/lusoris/revenge/commit/8669bb72551c7e9680e0ce6efda0d63a3ce96ce1))
* **library:** add LibraryService for tvshow and qar/fleet modules ([3959778](https://github.com/lusoris/revenge/commit/39597781b9b8b2778d54f7f55c50995907a79d34))
* **library:** implement Library service (Step 12) ([45fe4a4](https://github.com/lusoris/revenge/commit/45fe4a4f7270836750bd2348c8d867f8d0349e3f))
* **main:** wire Day 2-3 services to fx dependency injection ([76d4336](https://github.com/lusoris/revenge/commit/76d4336888004f034b64302d8fcc728d97fa57bd))
* **metadata:** add movie adapter for shared metadata service (A12.6) ([b95f095](https://github.com/lusoris/revenge/commit/b95f09532adcbdfccbc9e60210baac0e3860edae))
* **metadata:** add provider interface and types (A12.1) ([ad9322f](https://github.com/lusoris/revenge/commit/ad9322ff10c1e76963b15e514322cdd87145cda7))
* **metadata:** add River jobs for metadata refresh (A12.5) ([7c65f37](https://github.com/lusoris/revenge/commit/7c65f37bfb32e6f9c3128071493e024a74cec491))
* **metadata:** add service and fx module (A12.4, A12.7) ([b368f75](https://github.com/lusoris/revenge/commit/b368f7567e3c51d586b8934268e05e5c53c2b776))
* **metadata:** add StashDB provider for QAR modules ([e8bcb53](https://github.com/lusoris/revenge/commit/e8bcb5388430a1ee4bc90e5d40fdd6803babfad9))
* **metadata:** add TMDb provider implementation (A12.2) ([42da49e](https://github.com/lusoris/revenge/commit/42da49e1336308a14e4d179b773bf718aa0ac98a))
* **metadata:** add TVDb provider implementation (A12.3) ([7e89b2e](https://github.com/lusoris/revenge/commit/7e89b2ed883fc2117e340e1c8ff81c9a557fc62a))
* **mfa:** add database migrations for MFA tables ([5e1913a](https://github.com/lusoris/revenge/commit/5e1913a5b3021198a25b6e7c06e5d59636bb779a))
* **mfa:** add MFA API handlers and integrate with server ([5cee136](https://github.com/lusoris/revenge/commit/5cee13616792da97a01c364fb84d18ba68a80616))
* **mfa:** add SQLC queries for MFA operations ([aa3c2b6](https://github.com/lusoris/revenge/commit/aa3c2b6b7d3ea3d10ac5337e36ce334cce1f343c))
* **mfa:** implement backup codes and MFA manager (Phase 4) ([e72d7f7](https://github.com/lusoris/revenge/commit/e72d7f7ff9a527a1349b2a4d276859add1fe89c1))
* **mfa:** implement remember device setting [A0.8] ([df9b5ff](https://github.com/lusoris/revenge/commit/df9b5ff38adf20cb65cf7a13b6122b3baafa5f67))
* **mfa:** implement TOTP service (Phase 2) ([3a7464f](https://github.com/lusoris/revenge/commit/3a7464f322741398e33bd5cd838d63a2367312a3))
* **mfa:** implement WebAuthn service (Phase 3) ([f0c3da6](https://github.com/lusoris/revenge/commit/f0c3da69cff1132f5617d4d1472434154ac3b745))
* **mfa:** implement WebAuthn session caching [A0.6] ([59b01a5](https://github.com/lusoris/revenge/commit/59b01a5e99f0db7e7248678454797579200c0707))
* **mfa:** integrate MFA with auth service and unify password hashing ([3a1ae62](https://github.com/lusoris/revenge/commit/3a1ae626acaf74395326b50f7d224d6ea181c7bd))
* **middleware:** add request metadata extraction ([5c48a59](https://github.com/lusoris/revenge/commit/5c48a59bfacaa6607fa0e98f09a49fa8c24637fa))
* **movie:** add HTTP handlers and integrate into app ([f18891b](https://github.com/lusoris/revenge/commit/f18891b880001f811db40dddf3cf05d3fe3786f8))
* **movie:** add Library Provider for file scanning and matching ([d8789fc](https://github.com/lusoris/revenge/commit/d8789fc4d32f808fd13b8873bcb2e8a01c674a49))
* **movie:** add multi-language metadata enrichment service ([0c5d05e](https://github.com/lusoris/revenge/commit/0c5d05ed8c7b06fc2feda5bd63bb8ebcb9989486))
* **movie:** add multi-language support to repository layer ([26c721c](https://github.com/lusoris/revenge/commit/26c721cd98a5fd48d3c1789da4d99cb4197c5809))
* **movie:** add multi-language support to TMDb mapper ([8da09cb](https://github.com/lusoris/revenge/commit/8da09cb6826f90fe4c19e67967abda3eb308a290))
* **movie:** add repository layer with PostgreSQL implementation ([4f6441a](https://github.com/lusoris/revenge/commit/4f6441a4fbb0db3aa4d1c5492c2342ea04fafbf8))
* **movie:** add River Jobs for background processing ([033accd](https://github.com/lusoris/revenge/commit/033accd17bad29a889a896ddf8e3b170e36a27ae))
* **movie:** add service layer and fx module ([5ac9fe3](https://github.com/lusoris/revenge/commit/5ac9fe313131c0313c634ddb96edfe2c1cbe0511))
* **movie:** add SQLC queries for movie module ([31d35ed](https://github.com/lusoris/revenge/commit/31d35ed99276dd17a1b5de680026c20a64805efe))
* **movie:** add TMDb metadata service ([a70c7b5](https://github.com/lusoris/revenge/commit/a70c7b57e23d28e0c8a1bf49270705ae8e8a907d))
* **movie:** implement all repository methods (A1) ([26d7bdf](https://github.com/lusoris/revenge/commit/26d7bdf982c4cc86ef46948f1354daaf2c173516))
* **movie:** implement collection repository methods and API handlers ([78ed59d](https://github.com/lusoris/revenge/commit/78ed59dbfee94aef978afc33f1bd19227d547c89))
* **movie:** implement file match and metadata refresh jobs (A2) ([b4efe69](https://github.com/lusoris/revenge/commit/b4efe691b20f345fbb9e98edcace4a79d5c21ef3))
* **movie:** implement library matcher with Levenshtein scoring (A5) ([94f75bb](https://github.com/lusoris/revenge/commit/94f75bbdefc4e2655ceb5791b0f86b86cfd76c7a))
* **movie:** implement TMDb multi-language client ([f8fee1b](https://github.com/lusoris/revenge/commit/f8fee1bb8e76d2eb4a7b7be07e1bf20b97cb546b))
* **movies:** wire movie module and enable API handlers ([9d11bf7](https://github.com/lusoris/revenge/commit/9d11bf715607798f7bd2d921ed92bb95369c081d))
* **observability:** add request correlation IDs with X-Request-ID header ([ec0a9f8](https://github.com/lusoris/revenge/commit/ec0a9f8e02208b5638a7a2107f8f2689f519d40f))
* **oidc:** implement OIDC service (Step 10) ([a6cb554](https://github.com/lusoris/revenge/commit/a6cb5544d48bebc1bcf63e846f1361ac6caaa703))
* **phase7-8:** observability, notifications, cache, tests ([d5747a9](https://github.com/lusoris/revenge/commit/d5747a90d167a01186054d6ede1fc9595f6567b7))
* **playback:** add playback service skeleton ([5602442](https://github.com/lusoris/revenge/commit/56024428b014065dae6ae8636a9323de109868ac))
* **qar:** add fingerprint service, whisparr client, and search collections ([e6f909f](https://github.com/lusoris/revenge/commit/e6f909f865691d680e96a183c0c2722fd9b92b80))
* **qar:** add obfuscated QAR module structure ([d43dd66](https://github.com/lusoris/revenge/commit/d43dd6670d1c4205e86266358d5d40effd587743))
* **qar:** add relationship handlers for performer/studio/tag movies ([1896a66](https://github.com/lusoris/revenge/commit/1896a66e84ba9cd4d4003e7e97fddbfbed0ee97b))
* **qar:** add request system foundation (provisions, quotas, rules) ([5d3ccce](https://github.com/lusoris/revenge/commit/5d3ccce551468c58c78fc739d1ce0e24a518c2f5))
* **qar:** complete API handlers and add 30-day continue watching filter ([fdcbbe4](https://github.com/lusoris/revenge/commit/fdcbbe44e3bfd545a819589ede8f755da5d4421c))
* **qar:** implement crew and flag repository methods ([afc70bd](https://github.com/lusoris/revenge/commit/afc70bdf0ccacfc3e6146caf871087baa36430cc))
* **qar:** implement repository methods for expedition, voyage, fleet, port ([9ac02e9](https://github.com/lusoris/revenge/commit/9ac02e9db16e4893f613986813a5e92af818e382))
* **qar:** implement request system API handlers + fix sqlc typo ([d4cd6c7](https://github.com/lusoris/revenge/commit/d4cd6c7a1c262ffbd182eda64ee8a9b1d7245772))
* **radarr:** add Radarr admin API handlers and webhook endpoint ([2e67151](https://github.com/lusoris/revenge/commit/2e67151e484e712268b3c71b0ec665d7a3a8791b))
* **radarr:** implement Radarr integration client and sync service ([6ad5379](https://github.com/lusoris/revenge/commit/6ad5379d83b48782af1e4bbe28049925fc0a9fde))
* **raft:** implement leader election for cluster deployments (A8.2.1) ([59211d3](https://github.com/lusoris/revenge/commit/59211d3f8bdd75503d43e72949bf7c5f8d9f4899))
* **rbac:** add request and metadata permissions ([88e69a9](https://github.com/lusoris/revenge/commit/88e69a9ef096742232b8459cd7d82af6df8173bf))
* **rbac:** add resource grants, audit logging, and user preferences ([4062695](https://github.com/lusoris/revenge/commit/40626958dda7cff88a6a80da17bc3a356e85b872))
* **rbac:** Day 2 complete - RBAC system with permissions ([ee1327d](https://github.com/lusoris/revenge/commit/ee1327d34cb4f4be731d6d9c5d745a620de6b1ca))
* **rbac:** implement dynamic RBAC with Casbin + shared video_people ([a231e97](https://github.com/lusoris/revenge/commit/a231e97103efaa68b200fdb67cee97e0ca051a8a))
* **rbac:** implement fine-grained permissions system [A6.1] ([68512bc](https://github.com/lusoris/revenge/commit/68512bc2c51bbfb0d482ba53d58a5894e48fcc5a))
* **rbac:** implement RBAC service with Casbin (Step 8) ([8e35103](https://github.com/lusoris/revenge/commit/8e351035bd244f5f0305d34f30b7a6cee182bab1))
* **requests:** add polls system design ([88e69a9](https://github.com/lusoris/revenge/commit/88e69a9ef096742232b8459cd7d82af6df8173bf))
* **search:** add River job for search index operations ([19b3f20](https://github.com/lusoris/revenge/commit/19b3f209e9c9b3aabb2e2f959346c70b67cb62cd))
* **search:** add Typesense movie search service and API endpoints ([6a8701c](https://github.com/lusoris/revenge/commit/6a8701c12f2d4ab766fdbbeaf7249606db0247ca))
* **services:** Day 3 - Global services (activity, settings, apikeys) ([d1abf69](https://github.com/lusoris/revenge/commit/d1abf693bdecfea76ef7a772de857da2943dd424))
* **session:** add write-through caching and configurable TTL [A6.2] ([a55307c](https://github.com/lusoris/revenge/commit/a55307cd982ffb8355ab7222423feeda5b35c223))
* **session:** implement session service (Step 7 - repository & service only) ([8dc4398](https://github.com/lusoris/revenge/commit/8dc4398333121be8c88b2a4ed1377319df2fa467))
* **settings:** add repository layer for server and user settings ([c280ffb](https://github.com/lusoris/revenge/commit/c280ffb83f1e31c21a0da4538e7acddf20751901))
* **settings:** add service layer with business logic ([932e138](https://github.com/lusoris/revenge/commit/932e138ea953d7c4a954276ec6fa93814f508de1))
* **skeleton:** complete v0.1.0 skeleton implementation ([f51e14e](https://github.com/lusoris/revenge/commit/f51e14e78f6702e4dd77df4ed41551599f958aff))
* **sonarr:** add full Sonarr integration matching Radarr ([34ad1ad](https://github.com/lusoris/revenge/commit/34ad1ada7856d17197c27ab2dbc9f9b1ba3f5f59))
* **storage:** add S3-compatible storage backend for avatars (A8.1.2) ([4c8bb7c](https://github.com/lusoris/revenge/commit/4c8bb7c77947921afe0d9627109d21b450208363))
* **testutil:** implement Dragonfly and Typesense testcontainers (A3) ([014c27e](https://github.com/lusoris/revenge/commit/014c27eae55010eaab593ad2a17c56b15bbe152c))
* **tools:** add tools.go to track development dependencies ([f472eae](https://github.com/lusoris/revenge/commit/f472eaeebb27366a37407c0936abcc14ea4db8ac))
* **tvshow:** add background job workers with River (A11.6) ([37c0e48](https://github.com/lusoris/revenge/commit/37c0e48247924b969790eded1dcf2b28546585cb))
* **tvshow:** add database schema and sqlc queries (A11.1) ([3a09be9](https://github.com/lusoris/revenge/commit/3a09be969081f3e303371dd01ee0bf88defe9a89))
* **tvshow:** add domain models and types (A11.2) ([96e285d](https://github.com/lusoris/revenge/commit/96e285de03a01dcec195b75cd3c011da8db55cdd))
* **tvshow:** add module.go, jobs.go, metadata_provider.go ([88e69a9](https://github.com/lusoris/revenge/commit/88e69a9ef096742232b8459cd7d82af6df8173bf))
* **tvshow:** add repository interface and postgres implementation (A11.4) ([ccdc35a](https://github.com/lusoris/revenge/commit/ccdc35abdf8c7ea966ae37d31fd3f746e3c99d15))
* **tvshow:** add scanner and metadata adapters (A11.3) ([be94d6a](https://github.com/lusoris/revenge/commit/be94d6affb5db72578243abf50dc5bfc560e1bc7))
* **tvshow:** add service layer with business logic (A11.5) ([573b0fc](https://github.com/lusoris/revenge/commit/573b0fcf0c9555ec9dd2b832d6cc916bd5923e5a))
* **tvshow:** add TV shows module database migrations ([71f75bb](https://github.com/lusoris/revenge/commit/71f75bb2e98081d9c81de5ac5ae7aa1b0e622e06))
* **user:** add repository layer with PostgreSQL implementation ([66f6722](https://github.com/lusoris/revenge/commit/66f6722665786dcbe7b03786e575759f996429f1))
* **user:** add service layer with business logic ([8efb4eb](https://github.com/lusoris/revenge/commit/8efb4eb6b739a5bf1b1e9c2118e4a2a2c687f471))
* **ux:** add avatar system design ([88e69a9](https://github.com/lusoris/revenge/commit/88e69a9ef096742232b8459cd7d82af6df8173bf))
* **v0.1.0:** complete all TODO items - migrations, error wrapping, testing infrastructure ([5eaa7c4](https://github.com/lusoris/revenge/commit/5eaa7c4295cbcb1b29f928bcaafc0f604de3c96e))
* **v0.1.0:** complete skeleton deliverables ([fbab82b](https://github.com/lusoris/revenge/commit/fbab82bd22abfa53c11e7652b374a286a87b47a8))
* **v0.1.0:** implement skeleton - project structure foundation ([aac743f](https://github.com/lusoris/revenge/commit/aac743f87af41a8284bc64092c90bd7a8a50a0cf))
* **validate:** add safe type conversion package ([48a23c8](https://github.com/lusoris/revenge/commit/48a23c85fe9e477ee8a9caee42d8db1d6db3e4e0))
* **webhook:** implement custom payload templates with Go templates (A4) ([c19de7c](https://github.com/lusoris/revenge/commit/c19de7c0be0af8df5de2d27a5c45d9248b2702a3))


### Bug Fixes

* add missing template fields to HTTP_CLIENT.yaml ([fe607ad](https://github.com/lusoris/revenge/commit/fe607ad2c2392d5145cf5ff210b6e55d8b91a13e))
* add missing template fields to HTTP_CLIENT.yaml ([80eeda0](https://github.com/lusoris/revenge/commit/80eeda05fa1e225ed0b8f4875e0a68304c995d84))
* add newline at end of METADATA.yaml ([87e3155](https://github.com/lusoris/revenge/commit/87e3155962e6fafcf632c3d092fa17238913b319))
* **api:** prevent integer overflow in activity handler pagination ([9fac0b8](https://github.com/lusoris/revenge/commit/9fac0b8403a4bb9b2d2f721cba47320d68f8a4a3))
* **api:** prevent integer overflow in library handler pagination ([9eb2b05](https://github.com/lusoris/revenge/commit/9eb2b050c1aa9201e4c1a7c452458c5f8a2f7356))
* **api:** replace placeholder UUIDs with auth context extraction ([35f3425](https://github.com/lusoris/revenge/commit/35f3425716c0261d4c062c1c4b5369b1541b58e4))
* **auth:** prevent email enumeration via password reset ([773ef2a](https://github.com/lusoris/revenge/commit/773ef2a8ba8d0c93926776057978449748b2ed9a))
* **auth:** prevent username enumeration via timing attack ([fe627d0](https://github.com/lusoris/revenge/commit/fe627d0c4c8d91e9c7209ef933d4ecc1e0fd664a))
* **auth:** wrap user registration in transaction for atomicity ([678102d](https://github.com/lusoris/revenge/commit/678102d7475dce33dd4a69ad7d32203f0d80cd69))
* **cache:** improve TTL handling and eviction logic ([1ab26fd](https://github.com/lusoris/revenge/commit/1ab26fd79730f2cfb7924371e5b8f694f33379c8))
* **ci:** add missing sources and fix kind cluster config ([c73458e](https://github.com/lusoris/revenge/commit/c73458ea6234576e95e6466d38a2cf0daa5c9bad))
* **ci:** add npm install to doc-validation + fix import mock test ([ee58824](https://github.com/lusoris/revenge/commit/ee5882409c0955a51b2a6ef6af96e6969a04492b))
* **ci:** add packages permission for GHCR login ([2ad342f](https://github.com/lusoris/revenge/commit/2ad342ff775abebaa6588f7604ea50e27a8f13a4))
* **ci:** add packages:write permission for Docker push ([c6cb48c](https://github.com/lusoris/revenge/commit/c6cb48c33be71190fc8504d2b2f0f0575df88e6c))
* **ci:** disable Go cache until go.sum exists ([ed539b6](https://github.com/lusoris/revenge/commit/ed539b661ecb596f961f99092d024f525460c214))
* **ci:** exclude .github/docs from version checks ([1e98907](https://github.com/lusoris/revenge/commit/1e98907cae71858d2bd39d124cd929826fe57ae8))
* **ci:** exclude docs/dev/sources from hardcoded version check ([416d7bc](https://github.com/lusoris/revenge/commit/416d7bc3b041756ad60b57bdc1e9948ff48efa26))
* **ci:** fix integration tests + add versions dependency ([2fd22f6](https://github.com/lusoris/revenge/commit/2fd22f6750133bea40aa107c0074f50e24461673))
* **ci:** fix release workflow - add ARM cross-compilers and remove Docker Hub ([428bf8f](https://github.com/lusoris/revenge/commit/428bf8fe213ed40171e2b54f1aa1cc34b541140b))
* **ci:** remove conflicting allow-licenses from dependency-review ([0364f00](https://github.com/lusoris/revenge/commit/0364f0092be0dc3a29b1ef5017e3b933adb15bbe))
* **ci:** remove Go 1.24 from test matrix ([c692629](https://github.com/lusoris/revenge/commit/c6926296317282aea218cb1c8e6c1ddc4425a1f7))
* **ci:** remove GPL/AGPL deny-licenses (project is AGPL-3.0) ([76ea746](https://github.com/lusoris/revenge/commit/76ea746a4ddf1df9cf4d63da0ebba315a8719c03))
* **ci:** remove unsupported --out-format flag from golangci-lint v2.8.0 ([845d085](https://github.com/lusoris/revenge/commit/845d085cee5578aacf144fd3c195a54c20810190))
* **ci:** replace golangci-lint-action with direct install for v2.x compatibility ([0d36cac](https://github.com/lusoris/revenge/commit/0d36cac970bf0fc2176b7cb7d00364b1be8c7abc))
* **ci:** resolve test port conflicts in CI pipelines ([df5378c](https://github.com/lusoris/revenge/commit/df5378c51e7e353866f833df5fa7259f3a6475f2))
* **ci:** resolve workflow failures for Go 1.25 compatibility ([d8c3b60](https://github.com/lusoris/revenge/commit/d8c3b60ede97e65f69b50bd85d548bc2894c3f7d))
* **ci:** skip integration tests on macOS/Windows ([fc9bc3a](https://github.com/lusoris/revenge/commit/fc9bc3ad87da9968e75f715bd93cb9eb115f4414))
* **ci:** update codecov-action to v5 and use files parameter ([c0a59bd](https://github.com/lusoris/revenge/commit/c0a59bd75b3ebef36df84385b85163d013bfbd3f))
* **ci:** upgrade codeql-action v3 to v4 and add security-events permission ([9a9cf96](https://github.com/lusoris/revenge/commit/9a9cf96563b404bba8184d70b41fa3c55d37fdc6))
* **ci:** use grep -v to exclude .github/docs from version checks ([e5838ae](https://github.com/lusoris/revenge/commit/e5838ae061fc8650c383134106b036e55f9bad9a))
* **config,database:** remove required validation on optional fields and update migration test ([7dd223e](https://github.com/lusoris/revenge/commit/7dd223eb19be2c40cc1bce91a94fff67067aa9f0))
* **config:** add default Database.URL placeholder value ([501a063](https://github.com/lusoris/revenge/commit/501a063020d60f4553cb97e1824ccede2f81318c))
* consolidate migrations to single embedded directory ([55d8863](https://github.com/lusoris/revenge/commit/55d88632225b307738ad7aee1a3c67cf8e60ae49))
* **context:** add timeouts to async goroutine contexts (A7.6) ([1599dcb](https://github.com/lusoris/revenge/commit/1599dcbf0dcc51ec1342cc76b03433000b8e5755))
* correct design_refs format in HTTP_CLIENT.yaml ([8c41e84](https://github.com/lusoris/revenge/commit/8c41e847e1274eee9a1f20e93083ab1b32b8437d))
* **crypto:** add bcrypt backward compatibility for password verification ([525f095](https://github.com/lusoris/revenge/commit/525f095f46a93188b3f6b8ca25132f424105f34e))
* **database:** prevent integer overflow in connection pool config ([a85ee94](https://github.com/lusoris/revenge/commit/a85ee94376264cb4ab181fa5bf791e624d53fc2e))
* **database:** remove context parameter from NewPool for fx compatibility ([06b44f7](https://github.com/lusoris/revenge/commit/06b44f7c511b69a198e3332ed46eb9fcbaa0af9b))
* **data:** fix yamllint indentation errors in shared-sot.yaml and 03_DESIGN_DOCS_STATUS.yaml ([96e530f](https://github.com/lusoris/revenge/commit/96e530f5d1e772a07a59493918696270020d9bed))
* **data:** replace invalid status emoji in TECH_STACK.yaml ([9ea53b8](https://github.com/lusoris/revenge/commit/9ea53b8a251bd95bb1f6129c1f32e48ed0320a03))
* **deploy:** add serviceAccount to values.yaml and fix helm lint ([a45f97a](https://github.com/lusoris/revenge/commit/a45f97a2bbde658f037f052374871fd96262d941))
* **deploy:** remove postgresql dependency from Helm chart ([bb48648](https://github.com/lusoris/revenge/commit/bb48648ef6e65c1a1cf02ed75c8cf6b7045c23b5))
* **deps:** correct Dependabot commit message format ([1c56569](https://github.com/lusoris/revenge/commit/1c565695ecd906a2cc069e026095820202c4ad61))
* **diagrams:** use flowchart LR with subgraphs for cleaner layout ([d783519](https://github.com/lusoris/revenge/commit/d7835198c5fe88f89bc4b750e0687fc171cda355))
* **diagrams:** use invisible links for horizontal layout in Mermaid ([084a6f6](https://github.com/lusoris/revenge/commit/084a6f68e6c503102ab82d59786f7225a81a0997))
* **docker:** add casbin model config to repository ([9d18bff](https://github.com/lusoris/revenge/commit/9d18bffd1fedded43a6159e7e914b9abb255cbe2))
* **docker:** improve Docker configuration ([7a1b245](https://github.com/lusoris/revenge/commit/7a1b2454d057afe8983cedf93b8e3e967b45d603))
* **docker:** remove non-existent Alpine FFmpeg dev packages ([1c1ae18](https://github.com/lusoris/revenge/commit/1c1ae180b2f9058bba6b6ad2e748b977363e5472))
* **docs:** add missing wiki generation step to doc pipeline ([fda21d8](https://github.com/lusoris/revenge/commit/fda21d846042c63465ffd0707437063025473cf1))
* **docs:** correct markdown link paths + enhance link checker for YAML sources ([c9fc037](https://github.com/lusoris/revenge/commit/c9fc03705760dc6238b5f1619e06e606ea30abc1))
* **docs:** correct relative paths to sources in subdirectory design docs ([ff61b46](https://github.com/lusoris/revenge/commit/ff61b465190b4a4ac905148b95d05250fd6c89f2))
* **docs:** correct source link depth calculation in doc generator ([bf3116d](https://github.com/lusoris/revenge/commit/bf3116d7a7459aac6c6f1cc263715a3d2fd8a6af))
* **docs:** Escape parentheses in Mermaid labels to prevent parse errors ([c4a4815](https://github.com/lusoris/revenge/commit/c4a481557918bd7b16d3591049422b73fd831953))
* **docs:** fix design_refs relative paths with correct depth calculation ([7c94ed9](https://github.com/lusoris/revenge/commit/7c94ed9f271de324a9af975ecc036628f2d560b1))
* **docs:** fix doc generation pipeline and TOC handling ([74cd3b2](https://github.com/lusoris/revenge/commit/74cd3b261251edf407f17045f6fa403b13f8ccb4))
* **docs:** fix final manual file references - 100% design doc links fixed! ([69d026d](https://github.com/lusoris/revenge/commit/69d026dbf120b68fb5e066828deb981ef1acc311))
* **docs:** fix template/schema type handling and regenerate all docs ([0708b1f](https://github.com/lusoris/revenge/commit/0708b1f5e75a4c5ac53d4ddfb5e66fdbedcc0b40))
* **docs:** remove self-referential links and fix root file refs ([e04e79e](https://github.com/lusoris/revenge/commit/e04e79ee8078d8f48729160c8a21eaad0a60372b))
* **docs:** use ✅ instead of 🟢 for overall_status in TECH_STACK.yaml ([6e12dca](https://github.com/lusoris/revenge/commit/6e12dca8c31cf43b71430b8db8852d628eec08b2))
* **docs:** Use quoted labels in Mermaid to preserve parentheses ([2de02f0](https://github.com/lusoris/revenge/commit/2de02f06f39e9d4280fba9ee9f02c95fa258f32a))
* **docs:** use realistic browser User-Agent for source fetching ([6dd6ad8](https://github.com/lusoris/revenge/commit/6dd6ad88eccec0c2b23c83890624d764849130cf))
* fix Ruff linting and validation errors ([79d1ccc](https://github.com/lusoris/revenge/commit/79d1cccb42247d7dfd7e884318690ff5c237ed7b))
* **helm:** add serviceAccountName template to helpers ([78b1bf9](https://github.com/lusoris/revenge/commit/78b1bf9612eaa1e78ba4c674b65315ab824e85ac))
* **helm:** remove optional templates for minimal chart ([5af4d94](https://github.com/lusoris/revenge/commit/5af4d9471e4e571934e74f0d8e1574eb27755196))
* **image:** replace x/image with govips for image processing ([be0a507](https://github.com/lusoris/revenge/commit/be0a5079502928cf55b495dbe006b2896b19831b))
* **jobs:** prevent integer overflow in backoff calculation ([d6be6a0](https://github.com/lusoris/revenge/commit/d6be6a0d37ef47ae5f4529a86b0b496718ab77b6))
* **lint:** handle error return values in test cleanup functions ([af805ad](https://github.com/lusoris/revenge/commit/af805ad9c05f158f726b4f528f8cc354fb79b98d))
* **lint:** handle remaining unchecked error returns ([d91a178](https://github.com/lusoris/revenge/commit/d91a17851d74d35075fc7ff1cfdfef826d9d79a8))
* **lint:** resolve golangci-lint errors (errcheck, unused, govet) ([fd5fb86](https://github.com/lusoris/revenge/commit/fd5fb8699caa16e75658f52eaad93801dc941079))
* **mfa:** add TOTP upsert logic and improve error handling ([6e79104](https://github.com/lusoris/revenge/commit/6e791049274013648184d7c10b2af032c9dbdf10))
* **mfa:** fix SQL bugs and extend MFA service tests ([f7e5111](https://github.com/lusoris/revenge/commit/f7e51112da46f6210d08e776feee7ff1988215bf))
* **mfa:** implement GetUserIDFromContext using existing context helper ([e8928e6](https://github.com/lusoris/revenge/commit/e8928e6f9d9dc1dac1ecb3d3d4b3709cec66e8ae))
* **mfa:** use shared util package for gosec safe conversion helpers ([8a7f418](https://github.com/lusoris/revenge/commit/8a7f41879070ea6f9b21919f342854f00ff5ab41))
* **movie:** align worker with shared metadata service job types ([a666308](https://github.com/lusoris/revenge/commit/a66630822a49b42941b17feab0180f0595dc994e))
* **notification:** prevent goroutine leaks in dispatcher ([c0db13f](https://github.com/lusoris/revenge/commit/c0db13f91f858e752888f7bb94846a27f5a7ea56))
* **release:** consolidate Release Please config and reset to v0.2.0 ([ac63577](https://github.com/lusoris/revenge/commit/ac63577b8cb74c284b2b10e386dab3ebec8f26b1))
* remove obsolete HTTP server test ([fe82010](https://github.com/lusoris/revenge/commit/fe82010a08eadb20621628c27f51f004e1c63b53))
* replace valkey-go with rueidis in architecture ([75ecc8f](https://github.com/lusoris/revenge/commit/75ecc8f16d3d71f27c36df0f5d1f5ecff0ca3567))
* resolve Ruff import sorting in doc_generator.py ([ad4c125](https://github.com/lusoris/revenge/commit/ad4c12592c011dfba2a093fee701c7171976b8b7))
* resolve YAML validation errors and add comprehensive template tests ([442ee94](https://github.com/lusoris/revenge/commit/442ee949703fd1892f7c7e94e75d4efe009f3fa2))
* **scripts:** fix ruff linting issues breaking tests ([b212f6e](https://github.com/lusoris/revenge/commit/b212f6eb2f456203d3391fa9117bc9d03acc0f78))
* **security:** add safe integer conversions to prevent G115 overflow ([2634302](https://github.com/lusoris/revenge/commit/2634302ed09f401b04a8d7c1ce6c6d2f7cd02e06))
* **security:** resolve gosec G104/G112/G301/G306 issues ([e22695f](https://github.com/lusoris/revenge/commit/e22695f606533f58d661fbb31e06c780d24a6f7f))
* **session:** reorder refresh operations for resilience ([c79a800](https://github.com/lusoris/revenge/commit/c79a800140f4254284671085c3008e6a8d86df0a))
* **session:** return actual count from CleanupExpiredSessions ([015c0be](https://github.com/lusoris/revenge/commit/015c0be998b10443657888752ad94b11eb5daed7))
* **session:** update UpdateActivity signature to match design ([f1df16a](https://github.com/lusoris/revenge/commit/f1df16a450093741668ad26d3495e4892211225a))
* skip version drift test and fix import sorting ([7d68246](https://github.com/lusoris/revenge/commit/7d68246a16193d72db8bfb0d7da1608b596c7ac8))
* sync YAML overall_status to SOURCE_OF_TRUTH tables ([cbf6e1b](https://github.com/lusoris/revenge/commit/cbf6e1bc4c49b9b47d303a0e937fb654fc752dda))
* **templates:** handle undefined variables with | default() ([4b29bbd](https://github.com/lusoris/revenge/commit/4b29bbd535ab923f115291880c1182e17f2739a5))
* **templates:** hide empty sections, fix API fields, revert mermaid ([9123f46](https://github.com/lusoris/revenge/commit/9123f4604b5f9de314af6d65f0125f0f15570fec))
* **templates:** make wiki documentation human-readable ([ad56bec](https://github.com/lusoris/revenge/commit/ad56bec16d6066fbeba9130c6f636438e0729b7a))
* **templates:** properly format complex data structures ([d8fb543](https://github.com/lusoris/revenge/commit/d8fb54368a0d5d96795b1f6b6db1310709aee901))
* **templates:** remove duplicate titles and empty frontmatter from wiki docs ([ce462bb](https://github.com/lusoris/revenge/commit/ce462bb55f7c79e04277b0e9f8b4f0fcbd2cc76d))
* **templates:** remove YAML frontmatter from all docs ([f869bbb](https://github.com/lusoris/revenge/commit/f869bbb330cd01eb9f0b124fa4da3c502bbfe2f3))
* **templates:** strip leading whitespace from generated docs ([7e65c36](https://github.com/lusoris/revenge/commit/7e65c362c49912fbe5163dc48a1ca6189f07cba3))
* **tests:** increase batch_regenerate timeout from 5s to 120s ([558f263](https://github.com/lusoris/revenge/commit/558f26379fc5d16fe58feb50a1b36df501fdff56))
* **tests:** Remove ALTER DATABASE from migrations ([990ab8e](https://github.com/lusoris/revenge/commit/990ab8e25de641fb04207383019f3d41f751bd8d))
* **tests:** update integration tests for ogen API changes ([ef6eabb](https://github.com/lusoris/revenge/commit/ef6eabbef6bfc19c03162e2a6a4b04ac5b1eaa05))
* **tests:** Use pool stats instead of config for assertions ([641e9b8](https://github.com/lusoris/revenge/commit/641e9b874a0de28c6c9e5e211658e6a5e637f03b))
* **testutil:** prevent integer overflow in test database port conversion ([d6e8d23](https://github.com/lusoris/revenge/commit/d6e8d23b115813073317b38613e717d933dda94a))
* **testutil:** use project root migrations instead of embedded copy ([94bfceb](https://github.com/lusoris/revenge/commit/94bfceb65d4e044728e8c076a7089aa1c491079b))
* **user:** update tests to expect Argon2id instead of bcrypt ([d04ed11](https://github.com/lusoris/revenge/commit/d04ed115c095c514c9a8ddc732fd9ae239103bf4))
* **user:** wrap avatar upload in transaction for atomicity ([5df5a27](https://github.com/lusoris/revenge/commit/5df5a27cc35100e0826c2648614405030ec9fd7d))
* **yamllint:** disable indentation rules and ignore Helm templates ([7cbb8d9](https://github.com/lusoris/revenge/commit/7cbb8d98ecfb1a868afa0643a548540f0932b3f4))


### Performance Improvements

* **decimal:** replace shopspring/decimal with govalues/decimal ([8adc754](https://github.com/lusoris/revenge/commit/8adc754b08d23609c09e04589bb4a33850b18ca7))
* **http:** replace go-resty with imroc/req v3 for HTTP clients ([4186ae1](https://github.com/lusoris/revenge/commit/4186ae1aab8e2fb9dc9f0ac61718be42ae8ed99a))
* **uuid:** switch UUID generation from v4 to v7 for sortable IDs ([5bdcbe5](https://github.com/lusoris/revenge/commit/5bdcbe5a4b984a67f5299c4fee9e183c98ecf862))


### Code Refactoring

* **library:** migrate to per-module library tables ([4a92009](https://github.com/lusoris/revenge/commit/4a92009e58a6b9039afb57eafd2841dc601122c5))

## [Unreleased]

### Added
- N/A

### Changed
- N/A

### Fixed
- N/A

## [0.1.2] - 2026-02-02

### Added
- **Errors Package Tests**: Complete test coverage for error handling utilities
  - `internal/errors`: 44% → 100% coverage

### Test Infrastructure
- **errors/wrap_test.go**: Comprehensive tests for wrap.go functions
  - `TestWrapf`: 6 subtests covering format args, nil handling, nesting
  - `TestWithStack`: 5 subtests covering stack trace verification
  - `TestWrapSentinel`: 8 subtests covering all sentinel errors
  - `TestFormatError`: 7 subtests covering formatting scenarios
  - `TestWrapChaining`: chaining Wrap, Wrapf, WithStack
  - `TestConcurrentErrorCreation`: concurrent safety (100 goroutines)
  - `TestErrorMessageFormat`: message format verification

## [0.1.1] - 2026-02-02

### Added
- **Comprehensive Unit Tests**: Major test coverage improvements across core packages
  - `internal/api`: 41% → 97.1% coverage with fx.Lifecycle integration tests
  - `internal/config`: 10% → 76.2% coverage with loader and validation tests
  - `internal/infra/health`: 55% → 68.1% coverage with health check tests
  - `internal/infra/database`: 20% → 22% coverage with PoolConfig tests

### Test Infrastructure
- **api/handler_test.go**: Tests for GetLiveness, GetStartup, GetReadiness, NewError
  - Embedded-postgres integration for realistic database testing
  - Concurrent request handling tests
  - Error response validation
- **api/server_test.go**: Full fx.Lifecycle tests using fxtest.New
  - Server startup/shutdown lifecycle
  - Graceful shutdown verification
  - Concurrent request handling (50 parallel requests)
  - Configuration application tests
  - Multiple port sequence tests
- **config/loader_test.go**: Comprehensive configuration tests
  - Default value loading
  - YAML file loading
  - Environment variable overrides (REVENGE_* prefix)
  - Validation failure scenarios
  - MustLoad panic handling
- **database/pool_test.go**: Connection pool configuration tests
  - MaxConns calculation (CPU * 2 + 1)
  - URL parsing and validation
  - Connection timeout settings
- **health/checks_test.go**: Stub health check tests
  - CheckCache, CheckSearch, CheckJobs
  - Status constants validation
  - Concurrent check execution

### Notes
- Tests use embedded-postgres for integration testing
- Run with `-p 1` flag to avoid port conflicts in parallel mode
- Full test suite: `go test ./internal/... -cover -count=1 -p 1`

## [0.1.0] - 2026-02-02

### Added
- **HTTP Server**: Fully functional HTTP server with ogen-generated code from OpenAPI spec
- **Health Endpoints**:
  - `GET /health/live` - Liveness probe (Kubernetes)
  - `GET /health/ready` - Readiness probe with dependency checks
  - `GET /health/startup` - Startup probe
- **Configuration System**: YAML-based configuration with environment variable support
- **Structured Logging**: Dual logging with Zap (JSON) and Slog (structured)
- **Database Support**: PostgreSQL connection pooling with pgxpool and migrations
- **Dependency Injection**: Uber fx for lifecycle management
- **Docker Support**: Multi-stage Dockerfile and docker-compose for development
- **CI/CD Pipeline**: GitHub Actions with build, test, lint, and security scanning
- **Integration Tests**: Comprehensive test suite with testcontainers
- **OpenAPI Spec**: Full API specification at `api/openapi/openapi.yaml`

### Infrastructure
- PostgreSQL 18.1 support
- Dragonfly (Redis-compatible) cache client stub
- Typesense search client stub
- River background jobs client stub

### Developer Experience
- Makefile with common targets (`make build`, `make test`, `make lint`)
- GoReleaser configuration for releases
- Renovate for dependency updates
- CodeQL security scanning
