#!/bin/bash

source /pkg/bin/ztp_helper.sh

function configure_xr()
{
   ## Apply a blind config
   xrapply $1
   if [ $? -ne 0 ]; then
       echo "xrapply failed to run"
   fi
   xrcmd "show config failed" > /home/vagrant/config_failed_check
}

config_file=$1
configure_xr $config_file

cat /home/vagrant/config_failed_check
grep -q "ERROR" /home/vagrant/config_failed_check

if [ $? -ne 0 ]; then
    echo "Configuration was successful!"
    echo "Last applied configuration was:"
    xrcmd "show configuration commit changes last 1"
else
    echo "Configuration Failed. Check /home/vagrant/config_failed on the router for logs"
    xrcmd "show configuration failed" > /home/vagrant/config_failed
    exit 1
fi