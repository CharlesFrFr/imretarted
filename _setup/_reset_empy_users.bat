@echo off
title Zombie Server
echo [ Zombie Server ] Resetting empty users...

cd ./../
server.exe -remove_empty_users -return

pause