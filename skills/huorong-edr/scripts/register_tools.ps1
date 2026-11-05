# register_tools.ps1 —— 把火绒查询注册为 ARTEX agent 自定义工具（kind=http）
# 前置: 1) hr_api.py 已在主机运行 (http://127.0.0.1:18888)
#       2) ARTEX 已启动，且你有 API Key（ARTEX 设置页获取）
# 用法:
#   .\register_tools.ps1 -ApiKey "your-arte-key"
#   .\register_tools.ps1 -ApiKey "..." -ArtexBase "http://127.0.0.1:8000" -HrBase "http://host.docker.internal:18888"
param(
    [Parameter(Mandatory=$true)][string]$ApiKey,
    [string]$ArtexBase = "http://127.0.0.1:18787",
    [string]$HrBase = "http://host.docker.internal:18888",
    [string]$HrToken = ""
)

if (-not $HrToken) {
    $tokFile = 'E:\ARTEX-main\data\huorong\api_token.txt'
    if (Test-Path $tokFile) { $HrToken = (Get-Content $tokFile -Raw).Trim() }
}
if (-not $HrToken) { Write-Host "缺少火绒 API token（先启动 hr_api.py）" -ForegroundColor Red; exit 1 }

$tools = @(
    [ordered]@{
        key         = "huorong_query"
        description = "查询本机火绒(Huorong)安全软件数据。endpoint: status=防护状态/病毒库版本/事件统计; events=安全事件(需since如60m/24h/3d, types可逗号分隔: filemon,sysprot,quaran,scan,c_malinfo,instmon,update); quarantine=隔离区; trust=信任区。红蓝对抗中用于对照：哪些攻击动作被火绒检出。"
        schema      = [ordered]@{
            type       = "object"
            properties = [ordered]@{
                endpoint = [ordered]@{ type = "string"; enum = @("status","events","quarantine","trust"); description = "查询类型" }
                since    = [ordered]@{ type = "string"; description = "相对时间窗(仅events): 60m/24h/3d" }
                until    = [ordered]@{ type = "string"; description = "截止相对时长(仅events,可选): 10m" }
                types    = [ordered]@{ type = "string"; description = "事件类型,逗号分隔(仅events,可选)" }
                limit    = [ordered]@{ type = "integer"; description = "条数上限(默认100)" }
                format   = [ordered]@{ type = "string"; enum = @("markdown","json"); description = "输出格式(默认markdown)" }
            }
        }
        agents      = @()
        enabled     = $true
        kind        = "http"
        exec        = [ordered]@{
            method     = "GET"
            url        = "$HrBase/api/{endpoint}?since={since}&until={until}&types={types}&limit={limit}&format={format}"
            headers    = [ordered]@{ "X-HR-Token" = $HrToken }
            timeout_ms = 30000
        }
        deferred    = $false
    }
)

$headers = @{ Authorization = "Bearer $ApiKey" }
foreach ($t in $tools) {
    $body = $t | ConvertTo-Json -Depth 12
    try {
        $resp = Invoke-RestMethod -Method POST -Uri "$ArtexBase/api/tools/custom" `
            -Headers $headers -ContentType 'application/json' -Body $body -TimeoutSec 30
        Write-Host "✅ 已注册工具: $($t.key) -> $($resp | ConvertTo-Json -Compress)" -ForegroundColor Green
    } catch {
        Write-Host "❌ 注册失败 $($t.key): $($_.Exception.Message)" -ForegroundColor Red
        if ($_.ErrorDetails.Message) { Write-Host "   $($_.ErrorDetails.Message)" -ForegroundColor DarkYellow }
    }
}

Write-Host ""
Write-Host "注册后: ARTEX 页面刷新，agent 工具列表应出现 huorong_query。"
Write-Host "验证: 新建任务说「查一下火绒最近的拦截事件」，观察 agent 是否调用 huorong_query。"
