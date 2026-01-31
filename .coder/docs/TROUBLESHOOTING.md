# Coder Troubleshooting Guide

**Purpose**: Solutions for common issues with Coder workspaces and remote development

**Last Updated**: 2026-01-31

---

## Table of Contents

- [Workspace Issues](#workspace-issues)
- [Connection Issues](#connection-issues)
- [IDE Integration Issues](#ide-integration-issues)
- [Port Forwarding Issues](#port-forwarding-issues)
- [Performance Issues](#performance-issues)
- [Service Issues](#service-issues)
- [SSH Issues](#ssh-issues)
- [Resource Issues](#resource-issues)
- [Template Issues](#template-issues)

---

## Workspace Issues

### Workspace Won't Start

**Problem**: Workspace fails to start or hangs on "Starting..."

**Symptoms**:
- `coder list` shows status as "Starting" for >5 minutes
- "Failed to start workspace" error
- Workspace immediately stops after starting

**Solutions**:

1. **Check workspace status**:
   ```bash
   coder list
   # Look for error messages in Status column
   ```

2. **View workspace logs**:
   ```bash
   coder logs revenge-dev
   # Look for error messages during startup
   ```

3. **Stop and restart workspace**:
   ```bash
   coder stop revenge-dev
   # Wait 10 seconds
   coder start revenge-dev
   ```

4. **Delete and recreate workspace**:
   ```bash
   # CAUTION: This deletes all data in workspace
   coder delete revenge-dev --yes

   # Recreate from template
   coder create revenge-dev --template revenge
   ```

5. **Check Coder server status**:
   ```bash
   curl https://coder.ancilla.lol/healthz
   # Should return "OK"
   ```

### Workspace Keeps Stopping

**Problem**: Workspace stops unexpectedly or won't stay running

**Symptoms**:
- Workspace status changes to "Stopped" randomly
- Work is interrupted frequently

**Solutions**:

1. **Check auto-stop settings**:
   ```bash
   # View workspace details
   coder show revenge-dev
   # Look for "Auto-stop" duration
   ```

2. **Disable auto-stop** (if enabled):
   - Edit workspace template to increase/disable auto-stop
   - Or keep terminal connection active

3. **Check resource limits**:
   - Workspace might be stopped due to exceeding resource limits
   - View logs: `coder logs revenge-dev`

4. **Keep SSH connection alive**:
   ```bash
   # In ~/.ssh/config
   Host coder-*
       ServerAliveInterval 60
       ServerAliveCountMax 10
   ```

### Can't Delete Workspace

**Problem**: `coder delete` fails or hangs

**Solutions**:

1. **Force delete**:
   ```bash
   coder delete revenge-dev --yes --force
   ```

2. **Stop before deleting**:
   ```bash
   coder stop revenge-dev
   coder delete revenge-dev --yes
   ```

3. **Check for dependent resources**:
   - Some templates create external resources (volumes, networks)
   - These might prevent deletion
   - Contact Coder admin if stuck

---

## Connection Issues

### Can't Connect to Coder Server

**Problem**: `coder login` fails or can't reach Coder server

**Symptoms**:
- "Connection refused" error
- "Unable to connect to Coder server" error
- Timeout errors

**Solutions**:

1. **Verify server URL**:
   ```bash
   # Correct URL
   coder login https://coder.ancilla.lol

   # NOT http:// (use https://)
   # NOT trailing slash
   ```

2. **Check network connectivity**:
   ```bash
   ping coder.ancilla.lol
   # Should respond with low latency

   curl https://coder.ancilla.lol/healthz
   # Should return "OK"
   ```

3. **Check VPN/proxy settings**:
   - If behind corporate VPN/proxy, ensure Coder URL is accessible
   - Add to proxy exceptions if needed

4. **Clear Coder cache**:
   ```bash
   rm -rf ~/.config/coder/
   coder login https://coder.ancilla.lol
   ```

5. **Update Coder CLI**:
   ```bash
   # macOS
   brew upgrade coder

   # Linux
   curl -L https://coder.com/install.sh | sh

   # Verify version
   coder version
   ```

### Authentication Fails

**Problem**: "Authentication failed" or token errors

**Solutions**:

1. **Re-authenticate**:
   ```bash
   coder logout
   coder login https://coder.ancilla.lol
   # Browser opens for authentication
   ```

2. **Check token expiration**:
   ```bash
   coder tokens ls
   # Look for expired tokens

   # Remove expired tokens
   coder tokens rm <token-id>
   ```

3. **Clear session**:
   ```bash
   rm ~/.config/coder/session
   coder login https://coder.ancilla.lol
   ```

4. **Check browser cookies**:
   - Authentication happens via browser
   - Ensure cookies are enabled
   - Try different browser

---

## IDE Integration Issues

### VS Code Can't Connect

**Problem**: VS Code Remote-SSH can't connect to workspace

**Symptoms**:
- "Could not establish connection" error
- VS Code hangs on "Setting up SSH Host"
- "Permission denied" errors

**Solutions**:

1. **Verify SSH config**:
   ```bash
   coder config-ssh
   # Updates ~/.ssh/config with workspace entries
   ```

2. **Test SSH connection**:
   ```bash
   ssh coder-revenge-dev
   # Should connect without errors
   ```

3. **Regenerate SSH keys**:
   ```bash
   ssh-keygen -R coder-revenge-dev
   coder config-ssh
   ```

4. **Check VS Code Remote-SSH extension**:
   - Ensure extension is installed and updated
   - Check extension logs: Output → Remote-SSH

5. **Try VS Code browser** (workaround):
   ```bash
   coder open revenge-dev
   # Opens browser-based VS Code
   ```

### Zed SSH Connection Fails

**Problem**: Zed can't connect via SSH to workspace

**Solutions**:

1. **Update SSH config**:
   ```bash
   coder config-ssh
   ```

2. **Test connection**:
   ```bash
   ssh coder-revenge-dev
   # Should work before trying in Zed
   ```

3. **Check Zed SSH settings**:
   - Zed → Settings → Remote
   - Verify SSH command and options

4. **Use manual SSH config**:
   - If `coder config-ssh` doesn't work
   - Manually add to `~/.ssh/config`:
     ```
     Host coder-revenge-dev
         HostName coder.ancilla.lol
         User coder
         Port 22
         ProxyCommand coder ssh --stdio revenge-dev
     ```

### JetBrains Gateway Won't Connect

**Problem**: Gateway can't connect to Coder workspace

**Symptoms**:
- "Connection failed" in Gateway
- "Failed to start IDE backend" error
- Gateway hangs on connecting

**Solutions**:

1. **Re-login in Gateway**:
   - Open Gateway
   - Settings → Coder → Logout
   - Login again with https://coder.ancilla.lol

2. **Verify Coder plugin**:
   - Settings → Plugins
   - Ensure Coder plugin is installed and updated

3. **Check workspace status**:
   ```bash
   coder list
   # Workspace must be "Running"

   coder start revenge-dev
   ```

4. **Delete IDE backend cache** (in workspace):
   ```bash
   coder ssh revenge-dev
   rm -rf ~/.cache/JetBrains/
   exit
   ```

5. **Try different IDE**:
   - If GoLand fails, try IntelliJ IDEA
   - Backend download might be corrupted

See: [JETBRAINS_INTEGRATION.md](JETBRAINS_INTEGRATION.md#troubleshooting)

---

## Port Forwarding Issues

### Ports Not Accessible

**Problem**: Can't access forwarded ports (e.g., http://localhost:8096)

**Symptoms**:
- "Connection refused" on localhost:PORT
- "ERR_CONNECTION_REFUSED" in browser
- Service runs in workspace but not accessible locally

**Solutions**:

1. **Verify service is running** (in workspace):
   ```bash
   coder ssh revenge-dev
   curl http://localhost:8096/health
   # Should respond

   # Check what's listening
   netstat -tlnp | grep 8096
   ```

2. **Check port forwarding**:
   ```bash
   # List forwarded ports
   coder port-forward list revenge-dev
   ```

3. **Manually forward port**:
   ```bash
   coder port-forward revenge-dev --tcp 8096:8096
   # Keep terminal open
   # Access at http://localhost:8096
   ```

4. **Check firewall** (local machine):
   ```bash
   # macOS
   sudo /usr/libexec/ApplicationFirewall/socketfilterfw --list

   # Allow Coder CLI if blocked
   ```

5. **Try different local port**:
   ```bash
   # If 8096 is occupied locally
   coder port-forward revenge-dev --tcp 9000:8096
   # Access at http://localhost:9000
   ```

### Port Forward Keeps Dropping

**Problem**: Port forwarding connection drops frequently

**Solutions**:

1. **Check network stability**:
   ```bash
   ping -c 100 coder.ancilla.lol
   # Look for packet loss
   ```

2. **Use TCP instead of HTTP forwarding**:
   ```bash
   # TCP is more reliable
   coder port-forward revenge-dev --tcp 8096:8096
   ```

3. **Keep terminal active**:
   - Don't close terminal running port-forward
   - Port forward stops when terminal closes

4. **Use background process**:
   ```bash
   nohup coder port-forward revenge-dev --tcp 8096:8096 &
   ```

---

## Performance Issues

### Slow Workspace Performance

**Problem**: Workspace is slow, laggy, or unresponsive

**Symptoms**:
- Commands take long to execute
- IDE is slow to respond
- File operations are slow

**Solutions**:

1. **Check workspace resources**:
   ```bash
   coder stat revenge-dev
   # Look at CPU, memory, disk usage
   ```

2. **Increase workspace resources**:
   - Stop workspace
   - Recreate with more CPU/RAM:
     ```bash
     coder create revenge-dev --template revenge \
       --parameter cpu=8 \
       --parameter memory=16
     ```

3. **Check running processes** (in workspace):
   ```bash
   coder ssh revenge-dev
   top
   # Look for processes using high CPU/memory
   ```

4. **Clear caches**:
   ```bash
   coder ssh revenge-dev

   # Clear Go cache
   go clean -cache -modcache -testcache

   # Clear npm cache
   cd web && npm cache clean --force

   # Clear IDE caches
   rm -rf ~/.cache/
   ```

5. **Check disk space**:
   ```bash
   coder ssh revenge-dev
   df -h
   # Ensure sufficient free space (>5GB)
   ```

### High Latency

**Problem**: High latency between local machine and workspace

**Solutions**:

1. **Check network latency**:
   ```bash
   ping coder.ancilla.lol
   # Should be <100ms for good experience
   ```

2. **Use wired connection**:
   - WiFi can have higher latency
   - Use Ethernet if possible

3. **Choose closer region** (if multiple regions available):
   - Contact Coder admin about workspace regions

4. **Optimize IDE settings**:
   - Reduce real-time features (linting, hints)
   - See IDE-specific docs for performance tuning

---

## Service Issues

### PostgreSQL Won't Start

**Problem**: PostgreSQL service fails to start in workspace

**Solutions**:

1. **Check service status**:
   ```bash
   coder ssh revenge-dev
   docker ps -a | grep postgres
   # Or if using systemd
   systemctl status postgresql
   ```

2. **View service logs**:
   ```bash
   # Docker
   docker logs <postgres-container-id>

   # Or systemd
   journalctl -u postgresql -n 50
   ```

3. **Restart service**:
   ```bash
   # Docker
   docker restart <postgres-container-id>

   # Or systemd
   sudo systemctl restart postgresql
   ```

4. **Check ports**:
   ```bash
   netstat -tlnp | grep 5432
   # Ensure port 5432 is not already in use
   ```

5. **Recreate database** (if corrupted):
   ```bash
   # Stop service
   # Remove data directory
   # Restart service (will recreate)
   ```

### Dragonfly/Redis Connection Fails

**Problem**: Can't connect to Dragonfly cache

**Solutions**:

1. **Check if running**:
   ```bash
   coder ssh revenge-dev
   docker ps | grep dragonfly
   ```

2. **Test connection**:
   ```bash
   redis-cli -h localhost -p 6379 ping
   # Should return PONG
   ```

3. **Restart Dragonfly**:
   ```bash
   docker restart <dragonfly-container-id>
   ```

4. **Check logs**:
   ```bash
   docker logs <dragonfly-container-id>
   ```

---

## SSH Issues

### SSH Connection Refused

**Problem**: Can't SSH into workspace

**Solutions**:

1. **Verify workspace is running**:
   ```bash
   coder list
   # Status should be "Running"
   ```

2. **Update SSH config**:
   ```bash
   coder config-ssh
   ```

3. **Test with coder ssh**:
   ```bash
   coder ssh revenge-dev
   # Should connect without errors
   ```

4. **Check SSH keys**:
   ```bash
   ls -la ~/.ssh/
   # Should have id_rsa or id_ed25519

   # Generate if missing
   ssh-keygen -t ed25519
   ```

5. **Try verbose SSH**:
   ```bash
   ssh -v coder-revenge-dev
   # Shows detailed connection info
   ```

### Permission Denied

**Problem**: "Permission denied (publickey)" error

**Solutions**:

1. **Add key to SSH agent**:
   ```bash
   eval "$(ssh-agent -s)"
   ssh-add ~/.ssh/id_ed25519
   ```

2. **Check key permissions**:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   chmod 644 ~/.ssh/id_ed25519.pub
   ```

3. **Regenerate SSH config**:
   ```bash
   coder config-ssh
   ```

4. **Use password authentication** (if enabled):
   ```bash
   ssh -o PreferredAuthentications=password coder-revenge-dev
   ```

---

## Resource Issues

### Out of Memory

**Problem**: Workspace runs out of memory

**Symptoms**:
- "Cannot allocate memory" errors
- OOMKilled processes
- Workspace becomes unresponsive

**Solutions**:

1. **Check memory usage**:
   ```bash
   coder ssh revenge-dev
   free -h
   # Look at available memory
   ```

2. **Identify memory hogs**:
   ```bash
   ps aux --sort=-%mem | head -20
   ```

3. **Increase workspace memory**:
   ```bash
   coder stop revenge-dev
   # Recreate with more memory
   coder create revenge-dev --template revenge --parameter memory=16
   ```

4. **Close unused processes**:
   ```bash
   # Stop IDE backends
   pkill -f "jetbrains\|vscode"

   # Stop unnecessary services
   docker stop <unused-container>
   ```

### Out of Disk Space

**Problem**: Workspace runs out of disk space

**Solutions**:

1. **Check disk usage**:
   ```bash
   coder ssh revenge-dev
   df -h
   du -sh ~/* | sort -h
   ```

2. **Clean up caches**:
   ```bash
   # Go cache
   go clean -cache -modcache -testcache

   # Docker
   docker system prune -a

   # npm
   cd web && npm cache clean --force

   # IDE caches
   rm -rf ~/.cache/
   ```

3. **Remove old files**:
   ```bash
   # Old logs
   rm -rf ~/logs/*.log

   # Build artifacts
   rm -rf bin/
   ```

4. **Increase disk size**:
   - Contact Coder admin to increase workspace disk

---

## Template Issues

### Template Build Fails

**Problem**: Can't create workspace - template build fails

**Solutions**:

1. **Check template status**:
   ```bash
   coder templates list
   # Look for errors
   ```

2. **View template build logs**:
   ```bash
   coder templates versions list revenge
   coder templates versions <version-id> logs
   ```

3. **Contact Coder admin**:
   - Template issues usually require admin access
   - Provide error messages from logs

### Template Updates Not Applied

**Problem**: Workspace doesn't reflect template changes

**Solutions**:

1. **Check template version**:
   ```bash
   coder show revenge-dev
   # Look at "Template version"
   ```

2. **Update workspace**:
   ```bash
   coder update revenge-dev
   # Applies latest template version
   ```

3. **Recreate workspace** (if update fails):
   ```bash
   # Backup important data first
   coder delete revenge-dev --yes
   coder create revenge-dev --template revenge
   ```

---

## General Troubleshooting

### Check Coder CLI Version

```bash
coder version
# Ensure you're on latest version

# Update if needed
brew upgrade coder  # macOS
# Or
curl -L https://coder.com/install.sh | sh  # Linux
```

### Enable Debug Logging

```bash
export CODER_VERBOSE=true
coder <command>
# Shows detailed debug output
```

### Check Coder Status Page

- Visit Coder status page (if available)
- Check for known issues or maintenance

### Collect Logs for Support

```bash
# Workspace logs
coder logs revenge-dev > workspace.log

# CLI debug logs
CODER_VERBOSE=true coder list > cli-debug.log 2>&1

# SSH debug logs
ssh -vvv coder-revenge-dev > ssh-debug.log 2>&1
```

---

## Getting Help

### Internal Resources

- **Shared Troubleshooting**: [../../.shared/docs/TROUBLESHOOTING.md](../../.shared/docs/TROUBLESHOOTING.md)
- **IDE-specific issues**:
  - [../../.vscode/docs/TROUBLESHOOTING.md](../../.vscode/docs/TROUBLESHOOTING.md) (if exists)
  - [../../.zed/docs/TROUBLESHOOTING.md](../../.zed/docs/TROUBLESHOOTING.md)
  - [../../.jetbrains/docs/TROUBLESHOOTING.md](../../.jetbrains/docs/TROUBLESHOOTING.md)

### External Resources

- **Coder Docs**: https://coder.com/docs
- **Coder GitHub**: https://github.com/coder/coder/issues
- **Coder Discord**: https://discord.gg/coder (if available)

### Contact Admin

For infrastructure issues:
- Contact Coder server administrator
- Provide: workspace name, error messages, logs
- Include: `coder version`, `coder stat revenge-dev`

---

## Quick Reference

### Common Commands

```bash
# Workspace lifecycle
coder list                          # List workspaces
coder start revenge-dev             # Start workspace
coder stop revenge-dev              # Stop workspace
coder restart revenge-dev           # Restart workspace
coder delete revenge-dev --yes      # Delete workspace

# Connection
coder ssh revenge-dev               # SSH into workspace
coder config-ssh                    # Update SSH config

# Port forwarding
coder port-forward revenge-dev --tcp 8096:8096

# Logs and status
coder logs revenge-dev              # View logs
coder stat revenge-dev              # Resource usage
coder show revenge-dev              # Workspace details

# Authentication
coder logout                        # Logout
coder login https://coder.ancilla.lol  # Login
```

### Emergency Reset

If everything is broken:

```bash
# 1. Logout
coder logout

# 2. Clear config
rm -rf ~/.config/coder/

# 3. Clear SSH config
coder config-ssh --remove

# 4. Re-login
coder login https://coder.ancilla.lol

# 5. Regenerate SSH config
coder config-ssh

# 6. Restart workspace
coder restart revenge-dev
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
