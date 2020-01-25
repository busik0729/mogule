#!/bin/bash

go build &&
scp -r config busik@45.11.26.20:/home/busik/web/crm &&
scp -r migrations busik@45.11.26.20:/home/busik/web/crm &&
scp -r schemas busik@45.11.26.20:/home/busik/web/crm &&
scp package.txt busik@45.11.26.20:/home/busik/web/crm &&
scp update.sh busik@45.11.26.20:/home/busik/web/crm &&
scp runGetDependencies.sh busik@45.11.26.20:/home/busik/web/crm &&
scp crm busik@45.11.26.20:/home/busik/web/crm &&
ssh busik@45.11.26.20 '~/web/crm/update.sh'
