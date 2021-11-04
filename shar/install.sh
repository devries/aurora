#!/bin/sh

# Try some echo setup:
if (echo "testing\c"; echo 1,2,3) | grep c >/dev/null; then
  # Stardent Vistra SVR4 grep lacks -e, says ghazi@caip.rutgers.edu.
  if (echo -n testing; echo 1,2,3) | sed s/-n/xn/ | grep xn >/dev/null; then
    ac_n= ac_c='
' ac_t='	'
  else
    ac_n=-n ac_c= ac_t=
  fi
else
  ac_n= ac_c='\c' ac_t=
fi

less -eXF README-shar
echo $ac_n "Do you wish to continue? (yes/no): $ac_c"
read accept
if [ "$accept" != "y" ] && [ "$accept" != "Y" ] && [ "$accept" != "yes" ] && [ "$accept" != "YES" ]; then
    exit 1
fi

defaultprefix="/usr/local"

echo $ac_n "Where would you like to install aurora? [$defaultprefix]: $ac_c"
read prefix
if [ "$prefix" = "" ]; then
    prefix="$defaultprefix"
fi

# untar unix package to destination
binpath=${prefix}/bin
mkdir -p ${binpath}

# Determine Operating System
unameOut="$(uname -s)"
case "${unameOut}" in
  Linux*)
    arch="$(uname -m)"
    case $arch in
      x86_64*) machine=linux;;
      aarch64*) machine=linuxarm64;;
      arm*) machine=linuxarmhf;;
      *) machine=unknown;;
    esac
    ;;
  *)        machine=unknown;;
esac

if [ $machine != "unknown" ]; then
  echo "Detected: ${machine}"

  # Copy the binary
  cp ${machine}/aurora ${binpath}/aurora

  # make binary executable
  chmod -R ugo+x ${binpath}/aurora

  # Write some documentation
  echo "bouncerounds is now installed in ${binpath}"

  # Set up systemd script
  echo $ac_n "Would you like to create a systemd unit file? (yes/no): $ac_c"
  read create_systemd
  if [ "$create_systemd" != "y" ] && [ "$create_systemd" != "Y" ] && [ "$create_systemd" != "yes" ] && [ "$create_systemd" != "YES" ]; then
    exit 0
  fi

  echo $ac_n "Metrics address (ip:port): $ac_c"
  read metrics

  servicename="aurora"

  # Write out file
  cat > /etc/systemd/system/$servicename.service <<EOF
[Unit]
Description=Aurora Prometheus Service
Wants=network-online.target
After=network-online.target

[Service]
ExecStart=${binpath}/aurora -m ${metrics}

[Install]
WantedBy=multi-user.target
EOF

  echo "Created systemd service ${servicename}"

  # Load into systemd
  systemctl daemon-reload
  systemctl start $servicename.service
  systemctl enable $servicename.service

else
  arch="$(uname -m)"
  echo "Machine ${unameOut}-${arch} not recognized."
fi
