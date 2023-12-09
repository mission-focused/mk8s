# mk8s
mk8s is a lightweight and portable tool for installing Kubernetes across many nodes with a focus on air-gap support. It uses a declarative manifest to define cluster nodes and their roles. It uses local or SSH access to perform the bootstrap and execution of Kubernetes processes with the intent of being extensible and flexible. 

**Note:** This is merely an experiment at this point in time. All of this work is solely done in the open. Please leave comments or issues if you have any. 

## Why?
Orchestrating the installation of Kubernetes clusters - whether connected or disconnected - across many nodes is often left to the operator (hence why this project exists). Cloud and Hypervisors have understood workflows for deploying clusters - but bare-metal and disconnected environments are often more complicated. What if we could reduce the dependency down to networked compute (linux). 

What if we could do this without the complexities of Ansible and have more granular capabilities and error handling to respond to events accordingly. Enter mk8s. 

## Future Objectives
- Concurrency
  - Drive bootstrap to installation for many Kubernetes nodes simultaneously
  - Reduce total time to High-availability 
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

## Initial Build Phases

Phase 1: manifest mockup - Done
Phase 2: Pick a single distro and build a download command - Done
Phase 3: Write functionality to copy artifacts to target
Phase 4: Write functionality to bootstrap nodes
Phase 5: Write functionality to install Kubernetes and check for health

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

## Where to start next
- RKE2 download logic
  - How to download the list of artifacts... 