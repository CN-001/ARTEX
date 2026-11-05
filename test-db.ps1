<#
.SYNOPSIS
  用一次性 Postgres 跑 ARTEX 里"依赖数据库"的测试，跑完自动清理。

.DESCRIPTION
  ARTEX 的 server 测试在检测不到数据库配置时会整体跳过（t.Skipf "postgres unavailable"），
  所以 `go test ./...` 在没有库的机器上只看得到"跳过"。本脚本起一个**临时的、只绑回环的**
  Postgres 容器，把 ARTEX_PG_DSN 指过去，让这些用例真跑一遍，然后无论成败都删掉容器。

  安全边界：
    - 容器端口只绑定 127.0.0.1（局域网不可达）；
    - 库名/用户/密码都是一次性的 artex_test/artex/test，与生产库完全无关；
    - 脚本结束（含失败/中断）一定执行 docker rm -f，不留残留容器。

.PARAMETER Run
  传给 `go test -run` 的过滤表达式。默认跑三类依赖库的用例。

.PARAMETER Port
  临时库在宿主机回环上的端口，默认 5433（不占用生产库的 5432）。

.PARAMETER Image
  Postgres 镜像，默认 postgres:16-alpine（与 docker-compose.yml 一致）。

.PARAMETER KeepContainer
  保留容器（排查用）。默认跑完即删。

.EXAMPLE
  pwsh -File .\test-db.ps1
  跑默认的 DB 用例集。

.EXAMPLE
  pwsh -File .\test-db.ps1 -Run 'TestNotify' -KeepContainer
  只跑通知相关用例，并保留容器供排查。
#>
[CmdletBinding()]
param(
  [string]$Run = 'TestFindingTraffic|TestFindingWorkflow|TestNotify',
  [int]$Port = 5433,
  [string]$Image = 'postgres:16-alpine',
  [string]$Container = 'artex-testpg',
  [switch]$KeepContainer
)

$ErrorActionPreference = 'Stop'
$repo = $PSScriptRoot
if (-not (Test-Path (Join-Path $repo 'go.mod'))) {
  throw "找不到 go.mod：请把本脚本放在 ARTEX 仓库根目录（当前 $repo）"
}

function Info($m) { Write-Host "[test-db] $m" }
function Fail($m) { Write-Host "[test-db] $m" -ForegroundColor Red }

# 1) Docker 可用性
docker info *> $null
if ($LASTEXITCODE -ne 0) {
  Fail 'Docker 不可用（守护进程没起来）。先启动 Docker Desktop 再跑；'
  Fail '或者：设置 ARTEX_PG_DSN 指向一个已存在的测试库，然后直接 go test。'
  exit 2
}

# 2) 清掉同名残留容器（上次异常退出的）
docker rm -f $Container *> $null
$started = $false
$logFile = Join-Path ([IO.Path]::GetTempPath()) ("artex-test-db-{0:yyyyMMdd-HHmmss}.log" -f (Get-Date))

try {
  # 3) 起临时库（只绑回环）
  Info "启动临时库容器 $Container（$Image，127.0.0.1:$Port）"
  $null = docker run -d --name $Container `
    -e POSTGRES_PASSWORD=test -e POSTGRES_USER=artex -e POSTGRES_DB=artex_test `
    -p "127.0.0.1:${Port}:5432" $Image
  if ($LASTEXITCODE -ne 0) { throw "docker run 失败（端口 $Port 被占用？）" }
  $started = $true

  # 4) 等库就绪
  $ready = $false
  foreach ($i in 1..30) {
    $ping = docker exec $Container pg_isready -U artex -d artex_test 2>$null
    if ($ping -like '*accepting connections*') { $ready = $true; break }
    Start-Sleep -Seconds 2
  }
  if (-not $ready) { throw '临时库 60 秒内没有就绪' }
  Info '临时库已就绪'

  # 5) 指向临时库并跑测试
  $env:ARTEX_PG_DSN = "postgres://artex:test@127.0.0.1:$Port/artex_test?sslmode=disable"
  Info "ARTEX_PG_DSN -> 127.0.0.1:$Port/artex_test（一次性库，与生产无关）"
  Info "go test ./server/ -run '$Run' -count=1 -v"

  Push-Location $repo
  try {
    go test ./server/ -run $Run -count=1 -v *>&1 | Tee-Object -FilePath $logFile
    $code = $LASTEXITCODE
  } finally {
    Pop-Location
    Remove-Item Env:\ARTEX_PG_DSN -ErrorAction SilentlyContinue
  }

  # 6) 汇总
  $lines   = Get-Content $logFile
  $passed  = ($lines | Select-String '^--- PASS').Count
  $failed  = ($lines | Select-String '^--- FAIL').Count
  $skipped = ($lines | Select-String '^--- SKIP').Count
  Write-Host ''
  Info "结果：PASS=$passed FAIL=$failed SKIP=$skipped（日志：$logFile）"
  if ($skipped -gt 0) {
    Info '注意：仍有 SKIP。若库是好的，说明有用例走了别的跳过分支（例如平台相关）。'
  }
  if ($failed -gt 0) {
    $lines | Select-String '^--- FAIL' | ForEach-Object { Fail $_.Line }
  }
  exit $code
} finally {
  # 7) 无论如何都清理
  Remove-Item Env:\ARTEX_PG_DSN -ErrorAction SilentlyContinue
  if ($started -and -not $KeepContainer) {
    docker rm -f $Container *> $null
    Info '临时库容器已删除'
  } elseif ($started) {
    Info "保留容器 $Container（-KeepContainer）：docker rm -f $Container 手动清理"
  }
}
