#!/bin/sh

cat > .netrc <<EOF
machine github.com
  login $NET_RC_LOGIN
  password $NET_RC_PASSWORD
EOF

