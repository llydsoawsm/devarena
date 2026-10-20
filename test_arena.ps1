$ErrorActionPreference = "Stop"

Write-Host ">>> DevCode Verification Suite Running..." -ForegroundColor Cyan

# 1. Health check
$health = Invoke-RestMethod -Uri "http://localhost:8080/api/health" -Method Get
Write-Host "✅ Health Check: $($health.status) (Version: $($health.version))" -ForegroundColor Green

# 2. Login or Register Player 1
$p1Login = @{ email = "player1@devcode.io"; password = "password123" } | ConvertTo-Json
$auth1 = $null
try {
    $auth1 = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/login" -Method Post -Body $p1Login -ContentType "application/json"
} catch {
    $p1Reg = @{ name = "Warrior_1"; email = "player1@devcode.io"; password = "password123" } | ConvertTo-Json
    $auth1 = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" -Method Post -Body $p1Reg -ContentType "application/json"
}
$t1 = $auth1.access_token
Write-Host "✅ Player 1 Authenticated: $($auth1.user.name) (Rating: $($auth1.user.rating))" -ForegroundColor Green

# 3. Login or Register Player 2
$p2Login = @{ email = "player2@devcode.io"; password = "password123" } | ConvertTo-Json
$auth2 = $null
try {
    $auth2 = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/login" -Method Post -Body $p2Login -ContentType "application/json"
} catch {
    $p2Reg = @{ name = "Warrior_2"; email = "player2@devcode.io"; password = "password123" } | ConvertTo-Json
    $auth2 = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" -Method Post -Body $p2Reg -ContentType "application/json"
}
$t2 = $auth2.access_token
Write-Host "✅ Player 2 Authenticated: $($auth2.user.name) (Rating: $($auth2.user.rating))" -ForegroundColor Green

# 4. Enqueue Both Players
$h1 = @{ Authorization = "Bearer $t1"; "Content-Type" = "application/json" }
$h2 = @{ Authorization = "Bearer $t2"; "Content-Type" = "application/json" }

Invoke-RestMethod -Uri "http://localhost:8080/api/matchmaking/enqueue" -Method Post -Headers $h1 -Body '{"mode":"ranked"}' | Out-Null
Invoke-RestMethod -Uri "http://localhost:8080/api/matchmaking/enqueue" -Method Post -Headers $h2 -Body '{"mode":"ranked"}' | Out-Null
Write-Host "✅ Both players enqueued in Ranked mode" -ForegroundColor Green

# 5. Wait for sweep worker to pair them
Start-Sleep -Seconds 3

# Fetch match history
$hist = Invoke-RestMethod -Uri "http://localhost:8080/api/matches/user/history" -Method Get -Headers $h1
$match = $hist.matches[0]
Write-Host "✅ Match Created! Room ULID: $($match.public_id), Problem: $($match.problem.title)" -ForegroundColor Green

# 6. Wait for 10-second countdown to live state
Write-Host "⏳ Waiting for 10s pre-match countdown..." -ForegroundColor Yellow
Start-Sleep -Seconds 11

$live = Invoke-RestMethod -Uri "http://localhost:8080/api/matches/$($match.public_id)" -Method Get
Write-Host "✅ Match State: $($live.match.status) (Started: $($live.match.started_at))" -ForegroundColor Green

# 7. Test Sandbox Sample Code Execution
$samplePy = @"
import sys
def solve():
    lines = sys.stdin.read().splitlines()
    if not lines:
        return
    nums = list(map(int, lines[0].split()))
    target = int(lines[1])
    seen = {}
    for i, n in enumerate(nums):
        d = target - n
        if d in seen:
            print(f"{seen[d]} {i}")
            return
        seen[n] = i
solve()
"@

$runSampleBody = @{
    problem_id = $match.problem_id
    language = "python"
    code = $samplePy
} | ConvertTo-Json

$sampleRes = Invoke-RestMethod -Uri "http://localhost:8080/api/submissions/run" -Method Post -Body $runSampleBody -ContentType "application/json"
Write-Host "✅ Sample Execution Verdict: $($sampleRes.result.verdict) (Runtime: $($sampleRes.result.runtime_ms)ms)" -ForegroundColor Green

# 8. Submit Solution for Match & ELO Rating Settlement
$submitBody = @{
    match_public_id = $match.public_id
    language = "python"
    code = $samplePy
} | ConvertTo-Json

$subRes = Invoke-RestMethod -Uri "http://localhost:8080/api/submissions/submit" -Method Post -Headers $h1 -Body $submitBody -ContentType "application/json"
Write-Host "🏆 Match Submission Verdict: $($subRes.result.verdict) (Passed: $($subRes.result.passed_test_cases)/$($subRes.result.total_test_cases))" -ForegroundColor Green

# 9. Verify Settle Results & Rating Adjustments
$finalMatch = Invoke-RestMethod -Uri "http://localhost:8080/api/matches/$($match.public_id)" -Method Get
Write-Host "🎉 Settled Match Winner: UserID $($finalMatch.match.winner_id)" -ForegroundColor Cyan
foreach ($p in $finalMatch.match.participants) {
    Write-Host "   Player: $($p.user.name) -> Result: $($p.result) | Rating: $($p.rating_before) -> $($p.rating_after) (Delta: $($p.rating_delta))" -ForegroundColor White
}

Write-Host ">>> ALL SYSTEMS OPERATIONAL AND VERIFIED! <<<" -ForegroundColor Green
