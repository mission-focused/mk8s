# mk8s
mk8s is a lightweight and portable tool for installing Kubernetes across many nodes with a focus on air-gap support. It uses a declarative manifest to define cluster nodes and their roles. It uses local or SSH access to perform the bootstrap and execution of Kubernetes processes with the intent of being extensible and flexible. 

## ⚠️ Notice

There is a lot of work still to be done on this project and we need your help. There will be many permutations across many systems across many Kubernetes distributions and if you find a scenario that is not supported - please leave an issue with as much information as you can.

## Why?
Orchestrating the installation of Kubernetes clusters - whether connected or disconnected - across many nodes is often left to the operator (hence why this project exists). Cloud and Hypervisors have understood workflows for deploying clusters - but bare-metal and disconnected environments are often more complicated. What if we could reduce the dependency down to networked compute (linux). 

What if we could do this without the complexities of Ansible and have more granular capabilities and error handling to respond to events accordingly. Enter mk8s. 

## Distro Support
| Distribution | Supported |
|--------------|-----------|
| RKE2 | True |
| K3S | False |

## Inspiration
- [RKE2](https://github.com/rancher/rke2)
- [k3sup](https://github.com/alexellis/k3sup)

## Future Objectives
- Distro Support
  - Start with RKE2
  - Support K3s/Kubeadm in future builds

## Areas of Interest

Concurrency
- Distributing artifacts (when required)
- Bootstrapping nodes and awaiting signal to start primary process - should be much faster than Ansible
  - Begin artifact transfer and bootstrap on all nodes simultaneously
  - Continue Initial server node install - pause all other nodes
  - On signal from primary server ready - initiate all server nodes k8s processes to join
  - On signal from server nodes healthy - initiate all agent nodes 


## Potential Capabilities
- Automated Artifact download
  - Ability to pass in alternative download for primary artifacts
    - Add `artifacts` key to the manifest
  - Ability to check for the existence of the artifacts
    - During download to use local instead of re-downloading
    - `mk8s tools check-artifacts`
  - SHA validation of artifacts
    - RKE2/K3S
      - Do something with the shasum file
- Extensibility
  - Handle race conditions for primary node when installing fresh cluster
  - Handle joining standard server nodes to existing primary
  - Handle joining agents to an established control-plane

- Error Handling
  - Retry logic where applicable
