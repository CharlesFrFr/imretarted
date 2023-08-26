@echo off
title Zombie Server
echo [ Zombie Server ] Resetting Database

cd ./../
server.exe -get_users

pause