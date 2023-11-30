# mk8s
Declarative Multi-Node Kubernetes Installer

**Note:** This is merely an experiment at this point in time. All of my work is solely done in the open. Please leave comments or issues if you have any. 

## Purpose
Declarative multi-node kubernetes installer. Given a declarative manifest - install kubernetes across many nodes (to include the node running the tool). Make this process idempotent and air gap friendly. Allow for abstractions of various Kubernetes distributions. 

## Thoughts
Why hasn't anyone built an orchestrator of this kind before? Something that can install kubernetes across many nodes from a single point without the mess of dependencies that is Ansible. 