# Using the Garden Linux extension with Gardener as end-user

The [`core.gardener.cloud/v1beta1.Shoot` resource](https://github.com/gardener/gardener/blob/master/example/90-shoot.yaml) declares a few fields that should be considered when this OS extension is used. It essentially allows you to configure [Garden Linux](https://github.com/gardenlinux/gardenlinux) specific settings from the `Shoot` manifest.

In this document we describe how this configuration looks like and under which circumstances your attention may be required.

## Declaring Garden Linux specific configuration

To configure Garden Linux specific settings, you can declare a `OperatingSystemConfiguration` in the `Shoot` manifest for each worker pool at `.spec.provider.workers[].machine.image.providerConfig`. 

An example `OperatingSystemConfiguration` would look like this:

```yaml
providerConfig:
  apiVersion: gardenlinux.os.extensions.gardener.cloud/v1alpha1
  kind: OperatingSystemConfiguration
  cgroupVersion: v2
  netfilterBackend: iptables
  linuxSecurityModule: SELinux
```

You might also want to have a look at [how Garden Linux configures these settings](https://github.com/gardenlinux/gardenlinux/blob/main/docs/configuration/gardener-kernel-restart.md).

## Setting cgroup version of Garden Linux

Kubernetes version `>= v1.19` support the unified cgroup hierarchy (a.k.a. cgroup v2) on the worker nodes' operating system.

To configure cgroup v2, the following line can be included into the `OperatingSystemConfiguration`:

```yaml
  cgroupVersion: v2
```

If not specified, this setting will default to cgroup `v1`. Also, for Shoot clusters with K8S `< v1.19`, cgroup `v1` will be enforced.

### Possible values for `cgroupVersion` (case matters):

| value | result |
|---|---|
| `v1` | Garden Linux will be configured to use the classic cgroup hierarchy (cgroup v1) |
| `v2` | Garden Linux will be configured to use the unified cgroup hierarchy (cgroup v2) |

Configuration of the setting in Garden Linux takes place by [this script](https://github.com/gardenlinux/gardenlinux/blob/main/features/gardener/file.include/var/lib/gardener-gardenlinux/01_cgroup_configure.sh).


## Setting the Linux Security Module

This setting allows you to configure the Linux Security Module (lsm) to be `SELinux` or `AppArmor`. Certain Kubernetes workloads might require either lsm to be loaded at boot of the worker node and will fail to run if it is not active.

To configure SELinux, the following line can be included into the `OperatingSystemConfiguration`:

```yaml
  linuxSecurityModule: SELinux
```

If not specifief, this setting will default to `AppArmor`.

### Possible values for `linuxSecurityModule` (case matters):

| value | result |
|---|---|
| `AppArmor` | Garden Linux will be configured with _AppArmor_ as lsm |
| `SELinux` | Garden Linux will be configured with _SELinux_ as lsm |

Configuration of the setting in Garden Linux takes place by [this script](https://github.com/gardenlinux/gardenlinux/blob/main/features/gardener/file.include/var/lib/gardener-gardenlinux/02_configure_lsm.sh).


## Setting the netfilter backend:

This setting allows you to configure wether `iptables` will use the new `nf_tables` or the old `iptables` netfilter backend.

To configure `nf_tables`, the following line can be included into the `OperatingSystemConfiguration`:

```yaml
  netfilterBackend: iptables
```

If not specifief, this setting will default to `nftables`.

### Possible values for `netfilterBackend` (case matters):

| value | result |
|---|---|
| `nftables` | Garden Linux will be configured to use `nftables` as netfilter backend |
| `iptables` | Garden Linux will be configured to use `iptables` as netfilter backend |

Configuration of the setting in Garden Linux takes place by [this script](https://github.com/gardenlinux/gardenlinux/blob/main/features/gardener/file.include/var/lib/gardener-gardenlinux/03_iptables_backend.sh).
