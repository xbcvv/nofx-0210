<#
.SYNOPSIS
    任务文档归档脚本 (Archive Task Documentation Script)

.DESCRIPTION
    此脚本用于将当前的实施计划 (Plan) 和演练文档 (Walkthrough) 归档到 task/task_archives/ 目录。
    文件名格式为: YYYYMMDD_TaskName_Type.md
    
.PARAMETER TaskName
    任务名称 (必填)，用于生成文件名的一部分。推荐使用 PascalCase 或 Snake_Case。
    
.PARAMETER PlanPath
    实施计划文件的路径。默认为当前目录下的 "implementation_plan.md"。
    
.PARAMETER WalkthroughPath
    演练文档文件的路径。默认为当前目录下的 "walkthrough.md"。

.EXAMPLE
    .\scripts\archive_task.ps1 -TaskName "GlobalCommandConfig" -PlanPath "path/to/plan.md" -WalkthroughPath "path/to/walkthrough.md"
#>

param (
    [Parameter(Mandatory=$true)]
    [string]$TaskName,

    [string]$PlanPath,
    
    [string]$WalkthroughPath
)

$ErrorActionPreference = "Stop"

# 获取项目根目录 (假设脚本在 scripts/ 目录下)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir = Split-Path -Parent $ScriptDir
$ArchiveDir = Join-Path $RootDir "task\task_archives"

# 确保归档目录存在
if (-not (Test-Path $ArchiveDir)) {
    New-Item -ItemType Directory -Path $ArchiveDir | Out-Null
    Write-Host "已创建归档目录: $ArchiveDir" -ForegroundColor Cyan
}

# 生成日期前缀
$DatePrefix = Get-Date -Format "yyyyMMdd"

# 归档实施计划
if (-not [string]::IsNullOrWhiteSpace($PlanPath) -and (Test-Path $PlanPath)) {
    $TargetPlanName = "${DatePrefix}_${TaskName}_Plan.md"
    $TargetPlanPath = Join-Path $ArchiveDir $TargetPlanName
    Copy-Item -Path $PlanPath -Destination $TargetPlanPath -Force
    Write-Host "已归档实施计划: $TargetPlanName" -ForegroundColor Green
} elseif (-not [string]::IsNullOrWhiteSpace($PlanPath)) {
    Write-Warning "未找到实施计划文件: $PlanPath"
}

# 归档演练文档
if (-not [string]::IsNullOrWhiteSpace($WalkthroughPath) -and (Test-Path $WalkthroughPath)) {
    $TargetWalkthroughName = "${DatePrefix}_${TaskName}_Walkthrough.md"
    $TargetWalkthroughPath = Join-Path $ArchiveDir $TargetWalkthroughName
    Copy-Item -Path $WalkthroughPath -Destination $TargetWalkthroughPath -Force
    Write-Host "已归档演练文档: $TargetWalkthroughName" -ForegroundColor Green
} elseif (-not [string]::IsNullOrWhiteSpace($WalkthroughPath)) {
    Write-Warning "未找到演练文档文件: $WalkthroughPath"
}

Write-Host "任务 '$TaskName' 归档完成!" -ForegroundColor Green
