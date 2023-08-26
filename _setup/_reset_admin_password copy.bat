@echo off
title Zombie Server
echo [ Zombie Server ] Resetting Admin Password

cd ./../
server.exe -reset_admin_password -return

pause