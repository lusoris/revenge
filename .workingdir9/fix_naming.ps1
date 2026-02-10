#!/usr/bin/env pwsh
# Rename camelCase properties to snake_case in OpenAPI spec

$file = "c:\Users\ms\dev\revenge\api\openapi\openapi.yaml"
$content = Get-Content $file -Raw

$replacements = [ordered]@{
    "pageSize"              = "page_size"
    "providerDisplayName"   = "provider_display_name"
    "authorizationEndpoint" = "authorization_endpoint"
    "endSessionEndpoint"    = "end_session_endpoint"
    "userInfoEndpoint"      = "user_info_endpoint"
    "tokenEndpoint"         = "token_endpoint"
    "autoCreateUsers"       = "auto_create_users"
    "updateUserInfo"        = "update_user_info"
    "masterPlaylistUrl"     = "master_playlist_url"
    "displayName"           = "display_name"
    "providerName"          = "provider_name"
    "providerType"          = "provider_type"
    "clientSecret"          = "client_secret"
    "allowLinking"          = "allow_linking"
    "claimMappings"         = "claim_mappings"
    "roleMappings"          = "role_mappings"
    "accessToken"           = "access_token"
    "refreshToken"          = "refresh_token"
    "lastLoginAt"           = "last_login_at"
    "tokenType"             = "token_type"
    "expiresIn"             = "expires_in"
    "issuerUrl"             = "issuer_url"
    "isDefault"             = "is_default"
    "isEnabled"             = "is_enabled"
    "clientId"              = "client_id"
    "linkedAt"              = "linked_at"
    "jwksUri"               = "jwks_uri"
    "authUrl"               = "auth_url"
    "resourceType"          = "resource_type"
    "errorMessage"          = "error_message"
    "totalCount"            = "total_count"
    "successCount"          = "success_count"
    "failedCount"           = "failed_count"
    "oldestEntry"           = "oldest_entry"
    "newestEntry"           = "newest_entry"
    "resourceId"            = "resource_id"
    "ipAddress"             = "ip_address"
    "userAgent"             = "user_agent"
    "createdAt"             = "created_at"
    "updatedAt"             = "updated_at"
    "userId"                = "user_id"
    "realtimeMonitoring"    = "realtime_monitoring"
    "metadataProvider"      = "metadata_provider"
    "preferredLanguage"     = "preferred_language"
    "scanOnStartup"         = "scan_on_startup"
    "scannerConfig"         = "scanner_config"
    "itemsScanned"          = "items_scanned"
    "itemsUpdated"          = "items_updated"
    "itemsRemoved"          = "items_removed"
    "completedAt"           = "completed_at"
    "itemsAdded"            = "items_added"
    "latestScan"            = "latest_scan"
    "errorCount"            = "error_count"
    "startedAt"             = "started_at"
    "libraryId"             = "library_id"
    "scanType"              = "scan_type"
    "durationSeconds"       = "duration_seconds"
    "progressSeconds"       = "progress_seconds"
    "subtitleTracks"        = "subtitle_tracks"
    "subtitleTrack"         = "subtitle_track"
    "startPosition"         = "start_position"
    "audioTracks"           = "audio_tracks"
    "audioTrack"            = "audio_track"
    "isOriginal"            = "is_original"
    "mediaType"             = "media_type"
    "sessionId"             = "session_id"
    "expiresAt"             = "expires_at"
    "isForced"              = "is_forced"
    "mediaId"               = "media_id"
    "fileId"                = "file_id"
    "accountName"           = "account_name"
}

$totalChanges = 0
foreach ($old in $replacements.Keys) {
    $new = $replacements[$old]
    $keyPattern = "(?m)^(\s+)${old}:"
    $reqPattern = "(?m)^(\s+- )${old}\s*$"
    $keyMatches = [regex]::Matches($content, $keyPattern)
    $reqMatches = [regex]::Matches($content, $reqPattern)
    $count = $keyMatches.Count + $reqMatches.Count
    if ($count -gt 0) {
        $content = [regex]::Replace($content, $keyPattern, "`${1}${new}:")
        $content = [regex]::Replace($content, $reqPattern, "`${1}${new}")
        Write-Host "  $old -> $new  ($count)" -ForegroundColor Green
        $totalChanges += $count
    }
}

[System.IO.File]::WriteAllText($file, $content)
Write-Host "Total: $totalChanges changes" -ForegroundColor Cyan
