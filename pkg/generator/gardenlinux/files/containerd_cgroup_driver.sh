#!/bin/bash

set -Eeuo pipefail

source "$(dirname $0)/g_functions.sh"

# reconfigures containerd to use systemd as a cgroup driver on cgroup v2 enabled systems
function configure_containerd {
    desired_cgroup=$1
    CONTAINERD_CONFIG="/etc/containerd/config.toml"

    if [ ! -s "$CONTAINERD_CONFIG" ]; then
        echo "$CONTAINERD_CONFIG does not exist" >&2
        return
    fi

    if [[ "$desired_cgroup" == "v2" ]]; then
        echo "Configuring containerd cgroup driver to systemd"
        sed -i "s/SystemdCgroup *= *false/SystemdCgroup = true/" "$CONTAINERD_CONFIG"
    else
        echo "Configuring containerd cgroup driver to cgroupfs"
        sed -i "s/SystemdCgroup *= *true/SystemdCgroup = false/" "$CONTAINERD_CONFIG"
    fi
}

configure_containerd "$(check_current_cgroup)"
