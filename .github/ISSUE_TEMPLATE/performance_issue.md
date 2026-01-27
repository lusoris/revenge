---
name: Performance Issue
about: Report a performance problem or regression
title: '[PERF] '
labels: 'type/performance, status/needs-triage'
assignees: ''
---

## Performance Issue Description
<!-- A clear description of the performance problem -->

## Environment
- **OS**: [e.g. Ubuntu 22.04, Windows 11, macOS 14]
- **revenge Version**: [e.g. v0.1.0 or commit hash]
- **Hardware**:
  - CPU: [e.g. Intel i7-10700K, AMD Ryzen 5600X]
  - RAM: [e.g. 16GB]
  - Storage: [e.g. SSD, HDD, NVMe]
  - GPU (if applicable): [e.g. NVIDIA RTX 3070]
- **Database**: [SQLite or PostgreSQL version]
- **Library Size**: [e.g. 5000 movies, 100 TV shows]
- **Concurrent Users**: [estimated]

## Metrics

### Current Performance
<!-- Provide specific metrics -->
- Response time: [e.g. 5000ms]
- Memory usage: [e.g. 2GB]
- CPU usage: [e.g. 80%]
- Throughput: [e.g. 10 req/s]

### Expected Performance
<!-- What you expect based on documentation or previous versions -->
- Response time: [e.g. < 100ms]
- Memory usage: [e.g. < 512MB]

## Steps to Reproduce
1. Start server with configuration [...]
2. Perform action [...]
3. Observe performance metrics

## Profiling Data
<!-- If you have profiling data, please attach or link it -->

### pprof profiles
- CPU profile: [link or attachment]
- Memory profile: [link or attachment]
- Trace: [link or attachment]

### How to collect profiles
```bash
# CPU profile
curl http://localhost:8096/debug/pprof/profile?seconds=30 > cpu.pprof

# Memory profile
curl http://localhost:8096/debug/pprof/heap > mem.pprof

# Analyze
go tool pprof -http=:8080 cpu.pprof
```

## Logs
<!-- Relevant log output with timestamps -->
```
Paste logs here
```

## Additional Context
<!-- Add any other context about the problem -->

## Checklist
- [ ] I have searched existing issues
- [ ] I have provided specific metrics
- [ ] I have included environment details
- [ ] I have collected profiling data (if possible)
