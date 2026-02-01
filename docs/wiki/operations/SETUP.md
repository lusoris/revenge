

---
---

## Table of Contents

- [Production Deployment Setup](#production-deployment-setup)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Troubleshooting](#troubleshooting)
    - [Container won't start](#container-wont-start)
    - [Cannot connect to database](#cannot-connect-to-database)
    - [Reverse proxy returns 502 Bad Gateway](#reverse-proxy-returns-502-bad-gateway)
    - [Slow performance / high CPU](#slow-performance-high-cpu)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)


# Production Deployment Setup




> Deploy Revenge to your homelab or server with Docker, Kubernetes, or bare metal


This guide covers everything you need to deploy Revenge in production. The recommended approach uses Docker Compose for simple homelab setups - just create a docker-compose.yml, set a few environment variables, and run docker compose up. For larger deployments, Kubernetes (K3s) provides high availability and auto-scaling. Includes reverse proxy configuration (Traefik, Caddy, nginx), SSL setup with LetsEncrypt, backup strategies, and security hardening tips.


---





---


## How It Works

<!-- User-friendly explanation -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->








## Troubleshooting
### Container won't start

**Cause**: Environment misconfiguration or dependency issues

**Solution**: 1. Check logs: `docker compose logs revenge`
2. Verify environment variables in .env
3. Ensure PostgreSQL is ready: `docker compose ps postgres`

### Cannot connect to database

**Cause**: PostgreSQL not ready or network issue

**Solution**: 1. Check PostgreSQL health: `docker exec revenge-postgres pg_isready`
2. Verify DATABASE_URL in environment
3. Check network: `docker compose exec revenge ping postgres`

### Reverse proxy returns 502 Bad Gateway

**Cause**: Backend not running or unhealthy

**Solution**: 1. Verify Revenge is running: `docker compose ps`
2. Check backend is healthy: `curl http://localhost:8080/health/live`
3. Review proxy logs (Traefik/nginx/Caddy)

### Slow performance / high CPU

**Cause**: Resource constraints or heavy transcoding load

**Solution**: 1. Check Docker resource limits
2. Review active transcoding jobs: `docker compose logs revenge | grep transcode`
3. Increase dragonfly cache size
4. Consider offloading transcoding to Blackbeard


## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md)
- [FFmpeg Documentation](../../sources/media/ffmpeg.md)
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md)
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md)
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md)
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md)
- [Go io](../../sources/go/stdlib/io.md)
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)