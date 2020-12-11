#!/usr/bin/env bash
sudo yum -y install java-1.8.0
sudo mkdir /minecraft
sudo chown -R ec2-user:ec2-user /minecraft
cd /minecraft
aws s3 cp s3://merp-mc/setup/forge-1.16.4-35.1.0-installer.jar /minecraft/forge-1.16.4-35.1.0-installer.jar
java -jar forge-1.16.4-35.1.0-installer.jar --installServer
echo '#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).
#Mon Aug 06 18:11:14 UTC 2018
eula=true' > eula.txt
aws s3 cp s3://merp-mc/setup/mods.zip /minecraft/mods.zip
unzip mods.zip
sudo aws s3 cp s3://merp-mc/setup/minecraft.service /etc/systemd/system/minecraft.service
sudo systemctl daemon-reload
sudo service minecraft start