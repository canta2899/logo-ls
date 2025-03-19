$arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq 'Arm64') {
    "arm64"
} else {
    "amd64"
}

$os = "windows"

$latestVersion = (Invoke-RestMethod -Uri "https://api.github.com/repos/canta2899/logo-ls/releases/latest").tag_name

# Download URL
$downloadUrl = "https://github.com/canta2899/logo-ls/releases/download/$latestVersion/logo-ls-$latestVersion-$os-$arch.zip"

$installDir = "$env:USERPROFILE\.local\bin"

if (-Not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

if (-Not ($env:PATH -split ";" | Where-Object { $_ -eq $installDir })) {
    Write-Warning "$installDir is not in your PATH."
    Write-Host "You should either add it to your PATH or move logo-ls to a directory that is in your PATH."
}

$tempZip = [System.IO.Path]::Combine($env:TEMP, "logo-ls.zip")

Write-Host "Downloading logo-ls $latestVersion for $os-$arch..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip

Write-Host "Extracting to $installDir..."
Expand-Archive -Path $tempZip -DestinationPath $installDir -Force

Remove-Item $tempZip

Write-Host "logo-ls installed successfully to $installDir"
