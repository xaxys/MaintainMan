@echo off
cd %~dp0..
swag init
redoc-cli bundle %~dp0..\docs\swagger.json -o %~dp0..\docs\index.html