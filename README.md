# okd4-utils
A few useful setup and management utilities for OKD4 clusters

This tool can be used to automatically do the configuration required to set up an OKD4 cluster according to [Craig Robinson](https://itnext.io/guide-installing-an-okd-4-5-cluster-508a2631cbee?gi=be44dbb2f87f). This tool is intended to be used for educational purposes and the configuration of dev clusters, and will most likely end up being too inflexible for use with a real cluster. But, we'll see.

The ultimate goal is to automate absolutely everything in the linked guide, starting with config file generation, then file placement, then software installation, firewall configuration, and if I can find some magical way to configure the coreos nodes duirng this, then I also want to handle that.

This is designed to be run on the service machine that you are using for HAProxy and bind.
