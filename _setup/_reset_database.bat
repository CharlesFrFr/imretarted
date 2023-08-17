@echo off
title Zombie Server
echo [ Zombie Server ] Resetting Database

cd ./../
server.exe -reset_database

pause