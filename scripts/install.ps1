# Requires PowerShell 3.0 or higher
$ErrorActionPreference = "Stop"

function Get-LatestRelease {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/NethermindEth/sedge/releases/latest"
    return $response.tag_name
}

# Detect architecture
$ARCH = ""
if ([Environment]::Is64BitOperatingSystem) {
    $ARCH = "amd64"
} else {
    Write-Host "Unsupported architecture."
    Exit 1
}

# Get the latest version
$VERSION = Get-LatestRelease
Write-Host "Latest version is $VERSION"

# Construct the download URL
$BINARY_NAME = "sedge-${VERSION}-windows-${ARCH}.exe"
$DOWNLOAD_URL = "https://github.com/NethermindEth/sedge/releases/download/${VERSION}/${BINARY_NAME}"

# Download the binary
Write-Host "Downloading ${BINARY_NAME} from ${DOWNLOAD_URL}..."
$Output = "sedge.exe"
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $Output

# Move to a directory in PATH (e.g., C:\Program Files\sedge)
$Destination = "C:\Program Files\sedge"
if (-Not (Test-Path -Path $Destination)) {
    New-Item -ItemType Directory -Path $Destination
}

Write-Host "Installing sedge to $Destination..."
Move-Item -Path $Output -Destination $Destination -Force

# Add to PATH if not already present
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
if ($CurrentPath -notlike "*$Destination*") {
    Write-Host "Adding $Destination to system PATH..."
    $NewPath = "$CurrentPath;$Destination"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, [System.EnvironmentVariableTarget]::Machine)
}

Write-Host "Installation complete. You may need to restart your terminal for changes to take effect."