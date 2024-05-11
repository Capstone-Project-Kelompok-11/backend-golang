#!powershell

$env:DOCKER_BUILDKIT=1

$docker = "C:\Program Files\Docker\Docker\resources\bin\docker.exe"

& $docker compose -f .\docker-compose.yaml -p academy build
& $docker compose up
