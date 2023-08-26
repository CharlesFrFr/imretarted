@echo off
title Zombie Server
echo [ Zombie Server ] Listing users...

cd ./../
server.exe -get_users -return

pause