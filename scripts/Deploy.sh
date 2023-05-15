#!/usr/bin/env bash

CLOUD_USER=root
CLOUD_ADDR=skfw.net

cat << "EOF" | ssh -T $CLOUD_USER@$CLOUD_ADDR
systemctl stop app
EOF

scp \
cloud/init/base.psql \
cloud/app.service \
bin/start \
$CLOUD_USER@$CLOUD_ADDR:/app

cat << "EOF" | ssh -T $CLOUD_USER@$CLOUD_ADDR
cd /app
cat base.psql | sudo -u postgres psql
ln -sfn /app/app.service /lib/systemd/system/app.service
systemctl start app
systemctl daemon-reload
EOF

