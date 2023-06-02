#!/usr/bin/env bash

CLOUD_USER=root
CLOUD_ADDR=skfw.net

cat << "EOF" | ssh -T $CLOUD_USER@$CLOUD_ADDR
mkdir -p /app/assets/fonts /app/assets/public /app/assets/public/caches /app/assets/public/documents /app/assets/public/images /app/assets/public/videos /app/templates/documents/certificates
ln -sfn /app/assets /assets
ln -sfn /app/templates /templates
systemctl stop app
EOF

scp \
cloud/init/base.psql \
cloud/app.service \
bin/start \
$CLOUD_USER@$CLOUD_ADDR:/app

scp \
assets/fonts/arial.ttf \
$CLOUD_USER@$CLOUD_ADDR:/app/assets/fonts/arial.ttf

scp \
templates/documents/certificates/cert.pdf \
$CLOUD_USER@$CLOUD_ADDR:/app/templates/documents/certificates

cat << "EOF" | ssh -T $CLOUD_USER@$CLOUD_ADDR
cd /app
cat base.psql | sudo -u postgres psql
ln -sfn /app/app.service /lib/systemd/system/app.service
systemctl start app
systemctl daemon-reload
EOF

