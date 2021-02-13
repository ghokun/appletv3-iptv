#!/bin/bash
openssl req -new -nodes -newkey rsa:2048 -out redbulltv.pem -keyout redbulltv.key -x509 -days 7300 -subj "/C=US/CN=appletv.redbull.tv"
openssl x509 -in redbulltv.pem -outform der -out redbulltv.cer && cat redbulltv.key >> redbulltv.pem