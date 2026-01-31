# Longhorn Storage

> Source: https://longhorn.io/docs/
> Fetched: 2026-01-30T23:56:16.885701+00:00
> Content-Hash: e099644f20c4f814
> Type: html

---

The Longhorn Documentation

What is Longhorn?

Architecture and Concepts

Terminology

Best Practices

Document Conventions

Contributing

Important Notes

Installation and Setup

Quick Installation

Install as a Rancher Apps & Marketplace

Install with Kubectl

Install with Helm

Install with Helm Controller

Install with Fleet

Install with Flux

Install with ArgoCD

Air Gap Installation

Accessing the UI

Create an Ingress with Basic Authentication (nginx)

Create an HTTPRoute with Gateway API

Upgrade

Upgrading Longhorn Manager

Manually Upgrading Longhorn Engine

Automatically Upgrading Longhorn Engine

Instance Manager Pods During Upgrade

Uninstall Longhorn

Maintenance and Upgrade

Node Maintenance and Kubernetes Upgrade Guide

Nodes and Volumes

Nodes

Node Space Usage

Storage Tags

Configuring Defaults for Nodes and Disks

Multiple Disk Support

Scheduling

Evicting Replicas on Disabled Disks or Nodes

Node Conditions

Volumes

Create Longhorn Volumes

Longhorn PVC Ownership and Permission

Delete Longhorn Volumes

Detach Longhorn Volumes

ReadWriteMany (RWX) Volume

Use Longhorn Volume as an iSCSI Target

Volume Size

Viewing Workloads that Use a Volume

Volume Expansion

Trim Filesystem

Volume Conditions

High Availability

Auto Balance Replicas

Data Locality

Kubernetes Cluster Autoscaler Support (Experimental)

RWX Volume Fast Failover (Experimental)

Volume Recovery

Node Failure Handling with Longhorn

Backup and Restore

Create a Snapshot

Snapshot Space Management

Recurring Snapshots and Backups

Volume Clone Support

Disaster Recovery (DR) Volumes

Backup and Restore

Setting a Backup Target

Configure The Block Size Of Backup

Create a Backup

Restore from a Backup

Restoring Volumes for Kubernetes StatefulSets

Restore Volume Recurring Jobs from a Backup

Synchronize Backup Volumes Manually

CSI Snapshot Support

Enable CSI Snapshot Support on a Cluster

CSI VolumeSnapshot Associated with Longhorn BackingImage

CSI VolumeSnapshot Associated with Longhorn Snapshot

CSI VolumeSnapshot Associated with Longhorn Backup

Monitoring

Setting up Prometheus and Grafana to monitor Longhorn

Integrating Longhorn metrics into the Rancher monitoring system

Longhorn Metrics for Monitoring

Disk Health Monitoring

Kubelet Volume Metrics Support

Longhorn Alert Rule Examples

Advanced Resources

Longhorn VolumeAttachment

Deploy

Customizing Default Settings

Taints and Tolerations

Node Selector

CSI Component Pod Anti-Affinity

Rancher Windows Cluster

Priority Class

Revision Counter

Storage Network

OS/Distro Specific Configuration

Longhorn CSI on K3s

Longhorn CSI on RKE and CoreOS

Longhorn CSI on GKE

OCP/OKD Support

Container-Optimized OS (COS) Support

Talos Linux Support

Support Managed Kubernetes Service

Manage Node-Group on AWS EKS

Manage Node-Group on Azure AKS

Manage Node-Group on GCP GKE

Upgrade Kubernetes on AWS EKS

Upgrade Kubernetes on Azure AKS

Upgrade Kubernetes on GCP GKE

Data Integrity

Snapshot Data Integrity Check

Data Cleanup

Orphaned Data Cleanup

Orphaned Instance Cleanup

Replica Rebuilding

Fast Replica Rebuilding

Offline Replica Rebuilding

Scale Replica Rebuilding

Backing Image

Backing Image

Backing Image Backup

Backing Image Encryption

Security

Volume Encryption

MTLS Support

Command Line Tool (longhornctl)

Install longhornctl

Containerized Data Importer (CDI)

Longhorn with CDI Imports

Longhorn System Backup And Restore

Backup Longhorn System

Restore Longhorn System

Restore to a cluster contains data using Rancher snapshot

Restore to a new cluster using Velero

Cluster Restore

Restore cluster with a Rancher snapshot

Data Recovery

Identifying and Recovering from Data Errors

Exporting a Volume from a Single Replica

Identifying Corrupted Replicas

Recovering from a Full Disk

Recovering from a Longhorn Backup without System Installed

CSI Driver Migration

Migrating from the Flexvolume Driver to CSI

References

Settings

Storage Class Parameters

Python Client

Reference Setup, Performance, Scalability, and Sizing Guidelines

Longhorn Networking

Examples

Helm Values

Troubleshoot

Troubleshooting Problems

Support Bundle

V2 Data Engine (Technical Preview)

Prerequisites

Quick Start

Performance

Troubleshooting

Features

V2 Volume Clone Support

Replica Rebuild QoS

Configurable CPU Cores

Interrupt Mode Support

Node Disk Support

Selective V2 Data Engine Activation

UBLK Frontend Support

V2 Volume Expansion

The Longhorn Documentation

Cloud native distributed block storage for Kubernetes

Edit this page

Longhorn

is a lightweight, reliable, and powerful distributed

block storage

system for Kubernetes.

Longhorn implements distributed block storage using containers and microservices. Longhorn creates a dedicated storage controller for each block device volume and synchronously replicates the volume across multiple replicas stored on multiple nodes. The storage controller and replicas are themselves orchestrated using Kubernetes.

Features

Enterprise-grade distributed block storage with no single point of failure

Incremental snapshot of block storage

Backup to secondary storage (

NFS

or

S3

-compatible object storage) built on efficient change block detection

Recurring snapshots and backups

Automated, non-disruptive upgrades. You can upgrade the entire Longhorn software stack without disrupting running storage volumes.

An intuitive GUI dashboard

© 2019-2026 Longhorn Authors | Documentation Distributed under CC-BY-4.0

© 2026 The Linux Foundation. All rights reserved. The Linux Foundation has registered trademarks and uses trademarks. For a list of trademarks of The Linux Foundation,
please see our

Trademark Usage

page.

The Longhorn Documentation

What is Longhorn?

Architecture and Concepts

Terminology

Best Practices

Document Conventions

Contributing

Important Notes

Installation and Setup

Quick Installation

Install as a Rancher Apps & Marketplace

Install with Kubectl

Install with Helm

Install with Helm Controller

Install with Fleet

Install with Flux

Install with ArgoCD

Air Gap Installation

Accessing the UI

Create an Ingress with Basic Authentication (nginx)

Create an HTTPRoute with Gateway API

Upgrade

Upgrading Longhorn Manager

Manually Upgrading Longhorn Engine

Automatically Upgrading Longhorn Engine

Instance Manager Pods During Upgrade

Uninstall Longhorn

Maintenance and Upgrade

Node Maintenance and Kubernetes Upgrade Guide

Nodes and Volumes

Nodes

Node Space Usage

Storage Tags

Configuring Defaults for Nodes and Disks

Multiple Disk Support

Scheduling

Evicting Replicas on Disabled Disks or Nodes

Node Conditions

Volumes

Create Longhorn Volumes

Longhorn PVC Ownership and Permission

Delete Longhorn Volumes

Detach Longhorn Volumes

ReadWriteMany (RWX) Volume

Use Longhorn Volume as an iSCSI Target

Volume Size

Viewing Workloads that Use a Volume

Volume Expansion

Trim Filesystem

Volume Conditions

High Availability

Auto Balance Replicas

Data Locality

Kubernetes Cluster Autoscaler Support (Experimental)

RWX Volume Fast Failover (Experimental)

Volume Recovery

Node Failure Handling with Longhorn

Backup and Restore

Create a Snapshot

Snapshot Space Management

Recurring Snapshots and Backups

Volume Clone Support

Disaster Recovery (DR) Volumes

Backup and Restore

Setting a Backup Target

Configure The Block Size Of Backup

Create a Backup

Restore from a Backup

Restoring Volumes for Kubernetes StatefulSets

Restore Volume Recurring Jobs from a Backup

Synchronize Backup Volumes Manually

CSI Snapshot Support

Enable CSI Snapshot Support on a Cluster

CSI VolumeSnapshot Associated with Longhorn BackingImage

CSI VolumeSnapshot Associated with Longhorn Snapshot

CSI VolumeSnapshot Associated with Longhorn Backup

Monitoring

Setting up Prometheus and Grafana to monitor Longhorn

Integrating Longhorn metrics into the Rancher monitoring system

Longhorn Metrics for Monitoring

Disk Health Monitoring

Kubelet Volume Metrics Support

Longhorn Alert Rule Examples

Advanced Resources

Longhorn VolumeAttachment

Deploy

Customizing Default Settings

Taints and Tolerations

Node Selector

CSI Component Pod Anti-Affinity

Rancher Windows Cluster

Priority Class

Revision Counter

Storage Network

OS/Distro Specific Configuration

Longhorn CSI on K3s

Longhorn CSI on RKE and CoreOS

Longhorn CSI on GKE

OCP/OKD Support

Container-Optimized OS (COS) Support

Talos Linux Support

Support Managed Kubernetes Service

Manage Node-Group on AWS EKS

Manage Node-Group on Azure AKS

Manage Node-Group on GCP GKE

Upgrade Kubernetes on AWS EKS

Upgrade Kubernetes on Azure AKS

Upgrade Kubernetes on GCP GKE

Data Integrity

Snapshot Data Integrity Check

Data Cleanup

Orphaned Data Cleanup

Orphaned Instance Cleanup

Replica Rebuilding

Fast Replica Rebuilding

Offline Replica Rebuilding

Scale Replica Rebuilding

Backing Image

Backing Image

Backing Image Backup

Backing Image Encryption

Security

Volume Encryption

MTLS Support

Command Line Tool (longhornctl)

Install longhornctl

Containerized Data Importer (CDI)

Longhorn with CDI Imports

Longhorn System Backup And Restore

Backup Longhorn System

Restore Longhorn System

Restore to a cluster contains data using Rancher snapshot

Restore to a new cluster using Velero

Cluster Restore

Restore cluster with a Rancher snapshot

Data Recovery

Identifying and Recovering from Data Errors

Exporting a Volume from a Single Replica

Identifying Corrupted Replicas

Recovering from a Full Disk

Recovering from a Longhorn Backup without System Installed

CSI Driver Migration

Migrating from the Flexvolume Driver to CSI

References

Settings

Storage Class Parameters

Python Client

Reference Setup, Performance, Scalability, and Sizing Guidelines

Longhorn Networking

Examples

Helm Values

Troubleshoot

Troubleshooting Problems

Support Bundle

V2 Data Engine (Technical Preview)

Prerequisites

Quick Start

Performance

Troubleshooting

Features

V2 Volume Clone Support

Replica Rebuild QoS

Configurable CPU Cores

Interrupt Mode Support

Node Disk Support

Selective V2 Data Engine Activation

UBLK Frontend Support

V2 Volume Expansion