#!/usr/bin/env bash

crontab -r -u busik
crontab -u busik /home/busik/web/crm/cron_config
