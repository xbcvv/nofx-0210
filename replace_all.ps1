
$files = @(
    "docker-compose.stable.yml",
    "web\src\constants\branding.ts",
    "web\src\components\landing\FooterSection.tsx",
    "web\src\components\faq\FAQLayout.tsx",
    "web\src\components\faq\FAQContent.tsx",
    "scripts\pr-fix.sh",
    "scripts\pr-check.sh",
    "install-stable.sh"
)

$replacements = @{
    "NoFxAiOS/nofx" = "xbcvv/nofx-0210"
    "ghcr.io/nofxaios/nofx/nofx-backend:stable" = "xbcvv/nofx-backend:latest"
    "ghcr.io/nofxaios/nofx/nofx-frontend:stable" = "xbcvv/nofx-frontend:latest"
    "https://raw.githubusercontent.com/NoFxAiOS/nofx/release/stable" = "https://raw.githubusercontent.com/xbcvv/nofx-0210/main"
}

foreach ($file in $files) {
    if (Test-Path $file) {
        $content = Get-Content $file -Raw -Encoding UTF8
        $newContent = $content
        
        foreach ($key in $replacements.Keys) {
            $newContent = $newContent.Replace($key, $replacements[$key])
        }
        
        if ($content -ne $newContent) {
            $newContent | Set-Content $file -Encoding UTF8
            Write-Host "Updated $file"
        } else {
            Write-Host "No changes needed for $file"
        }
    } else {
        Write-Host "Skipping $file (not found)"
    }
}
