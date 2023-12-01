# mk8s
Declarative Multi-Node Kubernetes Installer

**Note:** This is merely an experiment at this point in time. All of my work is solely done in the open. Please leave comments or issues if you have any. 

## Purpose
Declarative multi-node kubernetes installer. Given a declarative manifest - install kubernetes across many nodes (to include the node running the tool). Make this process idempotent and air gap friendly. Allow for abstractions of various Kubernetes distributions. 

## Thoughts
Why hasn't anyone built an orchestrator of this kind before? Something that can install kubernetes across many nodes from a single point without the mess of dependencies that is Ansible. 

## Areas of Interest

Concurrency
- Distributing artifacts (when required)
- Bootstrapping nodes and awaiting signal to start primary process - should be much faster than Ansible
  - Begin artifact transfer and bootstrap on all nodes simultaneously
  - Continue Initial server node install - pause all other nodes
  - On signal from primary server ready - initiate all server nodes k8s processes to join
  - On signal from server nodes healthy - initiate all agent nodes 

## Build Phases

Phase 1: manifest mockup
Phase 2: Pick a single distro and build a download command / implementation
Phase 3: Write functionality to copy artifacts to target
Phase 4: Write functionality to bootstrap nodes
Phase 5: Write functionality to install Kubernetes and check for health