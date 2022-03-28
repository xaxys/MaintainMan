@echo off
swag init -d %~dp0\.. -o %~dp0\..\docs\
redoc-cli bundle %~dp0\..\docs\swagger.json -o %~dp0\..\docs\index.html