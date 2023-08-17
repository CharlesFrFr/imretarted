@echo off
title Zombie Server
echo [ Zombie Server ] Starting Server

cd ./../

:a
server.exe
echo [ Zombie Server ] Server Crashed
goto a


pause