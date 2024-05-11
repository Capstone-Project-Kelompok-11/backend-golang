#!pwsh

$Rules = @(
    @{ DisplayName = "Docker Client"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\com.docker.cli.exe" },
    @{ DisplayName = "Docker Desktop Backend"; ProgramPath = "C:\program files\docker\docker\resources\com.docker.backend.exe" },
    @{ DisplayName = "Docker Daemon"; ProgramPath = "C:\Program Files\Docker\Docker\resources\dockerd.exe" },
    @{ DisplayName = "Docker CLI"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\docker.exe" },
    @{ DisplayName = "Docker Compose CLI"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\docker-compose.exe" },
    @{ DisplayName = "Docker Credential Desktop"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\docker-credential-desktop.exe" },
    @{ DisplayName = "Docker Hub Tool CLI"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\hub-tool.exe" },
    @{ DisplayName = "Docker Kubernetes CLI"; ProgramPath = "C:\Program Files\Docker\Docker\resources\bin\kubectl.exe" }
)

foreach ($Rule in $Rules) {
    $RuleDisplayName = $Rule.DisplayName
    $RuleProgramPath = $Rule.ProgramPath
    $RuleIsExists = Get-NetFirewallRule -DisplayName $RuleDisplayName -ErrorAction SilentlyContinue

    if (-not $RuleIsExists) {
        Write-Host "Firewall rule '$RuleDisplayName' does not exist. Creating..."
        New-NetFirewallRule -DisplayName $RuleDisplayName -Enabled true -Profile Private,Public -Direction Inbound -Action Allow -Program $RuleProgramPath
        Write-Host "Firewall rule '$RuleDisplayName' created."
        
    } else {
        Write-Host "Firewall rule '$RuleDisplayName' already exists. Updating..."
        Set-NetFirewallRule -DisplayName $RuleDisplayName -Enabled true -Profile Private,Public -Direction Inbound -Action Allow -Program $RuleProgramPath
        Write-Host "Firewall rule '$RuleDisplayName' updated."
        
    }
}
