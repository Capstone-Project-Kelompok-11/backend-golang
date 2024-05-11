#!powershell

$env:DOCKER_BUILDKIT=1

$docker = "C:\Program Files\Docker\Docker\resources\bin\docker.exe"

& $docker system prune
