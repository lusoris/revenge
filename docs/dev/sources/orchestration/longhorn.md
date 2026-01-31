# Longhorn Storage

> Source: https://longhorn.io/docs/
> Fetched: 2026-01-31T16:03:30.131713+00:00
> Content-Hash: c72cc7629c6f7649
> Type: html

---

[The Longhorn Documentation](/docs/1.11.0/)

  * [What is Longhorn?](/docs/1.11.0/what-is-longhorn/)
  * [Architecture and Concepts](/docs/1.11.0/concepts/)
  * [Terminology](/docs/1.11.0/terminology/)
  * [Best Practices](/docs/1.11.0/best-practices/)
  * [Document Conventions](/docs/1.11.0/document-conventions/)
  * [Contributing](/docs/1.11.0/contributing/)



[Important Notes](/docs/1.11.0/important-notes/)

[Installation and Setup](/docs/1.11.0/deploy/)

[Quick Installation](/docs/1.11.0/deploy/install/)

  * [Install as a Rancher Apps & Marketplace](/docs/1.11.0/deploy/install/install-with-rancher/)
  * [Install with Kubectl](/docs/1.11.0/deploy/install/install-with-kubectl/)
  * [Install with Helm](/docs/1.11.0/deploy/install/install-with-helm/)
  * [Install with Helm Controller](/docs/1.11.0/deploy/install/install-with-helm-controller/)
  * [Install with Fleet](/docs/1.11.0/deploy/install/install-with-fleet/)
  * [Install with Flux](/docs/1.11.0/deploy/install/install-with-flux/)
  * [Install with ArgoCD](/docs/1.11.0/deploy/install/install-with-argocd/)
  * [Air Gap Installation](/docs/1.11.0/deploy/install/airgap/)

[Accessing the UI](/docs/1.11.0/deploy/accessing-the-ui/)

  * [Create an Ingress with Basic Authentication (nginx)](/docs/1.11.0/deploy/accessing-the-ui/longhorn-ingress/)
  * [Create an HTTPRoute with Gateway API](/docs/1.11.0/deploy/accessing-the-ui/longhorn-httproute/)

[Upgrade](/docs/1.11.0/deploy/upgrade/)

  * [Upgrading Longhorn Manager](/docs/1.11.0/deploy/upgrade/longhorn-manager/)
  * [Manually Upgrading Longhorn Engine](/docs/1.11.0/deploy/upgrade/upgrade-engine/)
  * [Automatically Upgrading Longhorn Engine](/docs/1.11.0/deploy/upgrade/auto-upgrade-engine/)
  * [Instance Manager Pods During Upgrade](/docs/1.11.0/deploy/upgrade/instance-manager-pods-during-upgrade/)

[Uninstall Longhorn](/docs/1.11.0/deploy/uninstall/)

[Maintenance and Upgrade](/docs/1.11.0/maintenance/)

  * [Node Maintenance and Kubernetes Upgrade Guide](/docs/1.11.0/maintenance/maintenance/)



[Nodes and Volumes](/docs/1.11.0/nodes-and-volumes/)

[Nodes](/docs/1.11.0/nodes-and-volumes/nodes/)

  * [Node Space Usage](/docs/1.11.0/nodes-and-volumes/nodes/node-space-usage/)
  * [Storage Tags](/docs/1.11.0/nodes-and-volumes/nodes/storage-tags/)
  * [Configuring Defaults for Nodes and Disks](/docs/1.11.0/nodes-and-volumes/nodes/default-disk-and-node-config/)
  * [Multiple Disk Support](/docs/1.11.0/nodes-and-volumes/nodes/multidisk/)
  * [Scheduling](/docs/1.11.0/nodes-and-volumes/nodes/scheduling/)
  * [Evicting Replicas on Disabled Disks or Nodes](/docs/1.11.0/nodes-and-volumes/nodes/disks-or-nodes-eviction/)
  * [Node Conditions](/docs/1.11.0/nodes-and-volumes/nodes/node-conditions/)

[Volumes](/docs/1.11.0/nodes-and-volumes/volumes/)

  * [Create Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/create-volumes/)
  * [Longhorn PVC Ownership and Permission](/docs/1.11.0/nodes-and-volumes/volumes/pvc-ownership-and-permission/)
  * [Delete Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/delete-volumes/)
  * [Detach Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/detaching-volumes/)
  * [ReadWriteMany (RWX) Volume](/docs/1.11.0/nodes-and-volumes/volumes/rwx-volumes/)
  * [Use Longhorn Volume as an iSCSI Target](/docs/1.11.0/nodes-and-volumes/volumes/iscsi/)
  * [Volume Size](/docs/1.11.0/nodes-and-volumes/volumes/volume-size/)
  * [Viewing Workloads that Use a Volume](/docs/1.11.0/nodes-and-volumes/volumes/workload-identification/)
  * [Volume Expansion](/docs/1.11.0/nodes-and-volumes/volumes/expansion/)
  * [Trim Filesystem](/docs/1.11.0/nodes-and-volumes/volumes/trim-filesystem/)
  * [Volume Conditions](/docs/1.11.0/nodes-and-volumes/volumes/volume-conditions/)



[High Availability](/docs/1.11.0/high-availability/)

  * [Auto Balance Replicas](/docs/1.11.0/high-availability/auto-balance-replicas/)
  * [Data Locality](/docs/1.11.0/high-availability/data-locality/)
  * [Kubernetes Cluster Autoscaler Support (Experimental)](/docs/1.11.0/high-availability/k8s-cluster-autoscaler/)
  * [RWX Volume Fast Failover (Experimental)](/docs/1.11.0/high-availability/rwx-volume-fast-failover/)
  * [Volume Recovery](/docs/1.11.0/high-availability/recover-volume/)
  * [Node Failure Handling with Longhorn](/docs/1.11.0/high-availability/node-failure/)



[Backup and Restore](/docs/1.11.0/snapshots-and-backups/)

  * [Create a Snapshot](/docs/1.11.0/snapshots-and-backups/setup-a-snapshot/)
  * [Snapshot Space Management](/docs/1.11.0/snapshots-and-backups/snapshot-space-management/)
  * [Recurring Snapshots and Backups](/docs/1.11.0/snapshots-and-backups/scheduling-backups-and-snapshots/)
  * [Volume Clone Support](/docs/1.11.0/snapshots-and-backups/csi-volume-clone/)
  * [Disaster Recovery (DR) Volumes](/docs/1.11.0/snapshots-and-backups/setup-disaster-recovery-volumes/)



[Backup and Restore](/docs/1.11.0/snapshots-and-backups/backup-and-restore/)

  * [Setting a Backup Target](/docs/1.11.0/snapshots-and-backups/backup-and-restore/set-backup-target/)
  * [Configure The Block Size Of Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/configure-backup-block-size/)
  * [Create a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/create-a-backup/)
  * [Restore from a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-from-a-backup/)
  * [Restoring Volumes for Kubernetes StatefulSets](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-statefulset/)
  * [Restore Volume Recurring Jobs from a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-recurring-jobs-from-a-backup/)
  * [Synchronize Backup Volumes Manually](/docs/1.11.0/snapshots-and-backups/backup-and-restore/synchronize_backup_volumes_manually/)

[CSI Snapshot Support](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/)

  * [Enable CSI Snapshot Support on a Cluster](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/enable-csi-snapshot-support/)
  * [CSI VolumeSnapshot Associated with Longhorn BackingImage](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backing-image/)
  * [CSI VolumeSnapshot Associated with Longhorn Snapshot](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-snapshot/)
  * [CSI VolumeSnapshot Associated with Longhorn Backup](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backup/)



[Monitoring](/docs/1.11.0/monitoring/)

  * [Setting up Prometheus and Grafana to monitor Longhorn](/docs/1.11.0/monitoring/prometheus-and-grafana-setup/)
  * [Integrating Longhorn metrics into the Rancher monitoring system](/docs/1.11.0/monitoring/integrating-with-rancher-monitoring/)
  * [Longhorn Metrics for Monitoring](/docs/1.11.0/monitoring/metrics/)
  * [Disk Health Monitoring](/docs/1.11.0/monitoring/disk-heath/)
  * [Kubelet Volume Metrics Support](/docs/1.11.0/monitoring/kubelet-volume-metrics/)
  * [Longhorn Alert Rule Examples](/docs/1.11.0/monitoring/alert-rules-example/)



[Advanced Resources](/docs/1.11.0/advanced-resources/)

  * [Longhorn VolumeAttachment](/docs/1.11.0/advanced-resources/volumeattachment/)



[Deploy](/docs/1.11.0/advanced-resources/deploy/)

  * [Customizing Default Settings](/docs/1.11.0/advanced-resources/deploy/customizing-default-settings/)
  * [Taints and Tolerations](/docs/1.11.0/advanced-resources/deploy/taint-toleration/)
  * [Node Selector](/docs/1.11.0/advanced-resources/deploy/node-selector/)
  * [CSI Component Pod Anti-Affinity](/docs/1.11.0/advanced-resources/deploy/csi-pod-antiaffinity-preset/)
  * [Rancher Windows Cluster](/docs/1.11.0/advanced-resources/deploy/rancher_windows_cluster/)
  * [Priority Class](/docs/1.11.0/advanced-resources/deploy/priority-class/)
  * [Revision Counter](/docs/1.11.0/advanced-resources/deploy/revision_counter/)
  * [Storage Network](/docs/1.11.0/advanced-resources/deploy/storage-network/)

[OS/Distro Specific Configuration](/docs/1.11.0/advanced-resources/os-distro-specific/)

  * [Longhorn CSI on K3s](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-k3s/)
  * [Longhorn CSI on RKE and CoreOS](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-rke-and-coreos/)
  * [Longhorn CSI on GKE](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-gke/)
  * [OCP/OKD Support](/docs/1.11.0/advanced-resources/os-distro-specific/okd-support/)
  * [Container-Optimized OS (COS) Support](/docs/1.11.0/advanced-resources/os-distro-specific/container-optimized-os-support/)
  * [Talos Linux Support](/docs/1.11.0/advanced-resources/os-distro-specific/talos-linux-support/)

[Support Managed Kubernetes Service](/docs/1.11.0/advanced-resources/support-managed-k8s-service/)

  * [Manage Node-Group on AWS EKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-eks/)
  * [Manage Node-Group on Azure AKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-aks/)
  * [Manage Node-Group on GCP GKE](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-gke/)
  * [Upgrade Kubernetes on AWS EKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-eks/)
  * [Upgrade Kubernetes on Azure AKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-aks/)
  * [Upgrade Kubernetes on GCP GKE](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-gke/)

[Data Integrity](/docs/1.11.0/advanced-resources/data-integrity/)

  * [Snapshot Data Integrity Check](/docs/1.11.0/advanced-resources/data-integrity/snapshot-data-integrity-check/)

[Data Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/)

  * [Orphaned Data Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/orphaned-data-cleanup/)
  * [Orphaned Instance Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/orphaned-instance-cleanup/)

[Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/)

  * [Fast Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/fast-replica-rebuilding/)
  * [Offline Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/offline-replica-rebuilding/)
  * [Scale Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/scale-replica-rebuilding/)

[Backing Image](/docs/1.11.0/advanced-resources/backing-image/)

  * [Backing Image](/docs/1.11.0/advanced-resources/backing-image/backing-image/)
  * [Backing Image Backup](/docs/1.11.0/advanced-resources/backing-image/backing-image-backup/)
  * [Backing Image Encryption](/docs/1.11.0/advanced-resources/backing-image/backing-image-encryption/)

[Security](/docs/1.11.0/advanced-resources/security/)

  * [Volume Encryption](/docs/1.11.0/advanced-resources/security/volume-encryption/)
  * [MTLS Support](/docs/1.11.0/advanced-resources/security/mtls-support/)

[Command Line Tool (longhornctl)](/docs/1.11.0/advanced-resources/longhornctl/)

  * [Install longhornctl](/docs/1.11.0/advanced-resources/longhornctl/install-longhornctl/)

[Containerized Data Importer (CDI)](/docs/1.11.0/advanced-resources/containerized-data-importer/)

  * [Longhorn with CDI Imports](/docs/1.11.0/advanced-resources/containerized-data-importer/containerized-data-importer/)

[Longhorn System Backup And Restore](/docs/1.11.0/advanced-resources/system-backup-restore/)

  * [Backup Longhorn System](/docs/1.11.0/advanced-resources/system-backup-restore/backup-longhorn-system/)
  * [Restore Longhorn System](/docs/1.11.0/advanced-resources/system-backup-restore/restore-longhorn-system/)
  * [Restore to a cluster contains data using Rancher snapshot](/docs/1.11.0/advanced-resources/system-backup-restore/restore-to-a-cluster-contains-data-using-rancher-snapshot/)
  * [Restore to a new cluster using Velero](/docs/1.11.0/advanced-resources/system-backup-restore/restore-to-a-new-cluster-using-velero/)

[Cluster Restore](/docs/1.11.0/advanced-resources/cluster-restore/)

  * [Restore cluster with a Rancher snapshot](/docs/1.11.0/advanced-resources/cluster-restore/rancher-cluster-restore/)

[Data Recovery](/docs/1.11.0/advanced-resources/data-recovery/)

  * [Identifying and Recovering from Data Errors](/docs/1.11.0/advanced-resources/data-recovery/data-error/)
  * [Exporting a Volume from a Single Replica](/docs/1.11.0/advanced-resources/data-recovery/export-from-replica/)
  * [Identifying Corrupted Replicas](/docs/1.11.0/advanced-resources/data-recovery/corrupted-replica/)
  * [Recovering from a Full Disk](/docs/1.11.0/advanced-resources/data-recovery/full-disk/)
  * [Recovering from a Longhorn Backup without System Installed](/docs/1.11.0/advanced-resources/data-recovery/recover-without-system/)

[CSI Driver Migration](/docs/1.11.0/advanced-resources/driver-migration/)

  * [Migrating from the Flexvolume Driver to CSI](/docs/1.11.0/advanced-resources/driver-migration/migrating-flexvolume/)



[References](/docs/1.11.0/references/)

  * [Settings](/docs/1.11.0/references/settings/)
  * [Storage Class Parameters](/docs/1.11.0/references/storage-class-parameters/)
  * [Python Client](/docs/1.11.0/references/longhorn-client-python/)
  * [Reference Setup, Performance, Scalability, and Sizing Guidelines](/docs/1.11.0/references/reference-setup-performance-scalability-and-sizing-guidelines/)
  * [Longhorn Networking](/docs/1.11.0/references/networking/)
  * [Examples](/docs/1.11.0/references/examples/)
  * [Helm Values](/docs/1.11.0/references/helm-values/)



[Troubleshoot](/docs/1.11.0/troubleshoot/)

  * [Troubleshooting Problems](/docs/1.11.0/troubleshoot/troubleshooting/)
  * [Support Bundle](/docs/1.11.0/troubleshoot/support-bundle/)



[V2 Data Engine (Technical Preview)](/docs/1.11.0/v2-data-engine/)

  * [Prerequisites](/docs/1.11.0/v2-data-engine/prerequisites/)
  * [Quick Start](/docs/1.11.0/v2-data-engine/quick-start/)
  * [Performance](/docs/1.11.0/v2-data-engine/performance/)
  * [Troubleshooting](/docs/1.11.0/v2-data-engine/troubleshooting/)



[Features](/docs/1.11.0/v2-data-engine/features/)

  * [V2 Volume Clone Support](/docs/1.11.0/v2-data-engine/features/volume-clone/)
  * [Replica Rebuild QoS](/docs/1.11.0/v2-data-engine/features/replica-rebuild-qos/)
  * [Configurable CPU Cores](/docs/1.11.0/v2-data-engine/features/configurable-cpu-cores/)
  * [Interrupt Mode Support](/docs/1.11.0/v2-data-engine/features/interrupt-mode/)
  * [Node Disk Support](/docs/1.11.0/v2-data-engine/features/node-disk-support/)
  * [Selective V2 Data Engine Activation](/docs/1.11.0/v2-data-engine/features/selective-v2-data-engine-activation/)
  * [UBLK Frontend Support](/docs/1.11.0/v2-data-engine/features/ublk-frontend-support/)
  * [V2 Volume Expansion](/docs/1.11.0/v2-data-engine/features/volume-expansion/)



The Longhorn Documentation

Cloud native distributed block storage for Kubernetes

[Edit this page](https://github.com/longhorn/website/edit/master/content/docs/1.11.0/_index.md)

**Longhorn** is a lightweight, reliable, and powerful distributed [block storage](https://cloudacademy.com/blog/object-storage-block-storage/) system for Kubernetes.

Longhorn implements distributed block storage using containers and microservices. Longhorn creates a dedicated storage controller for each block device volume and synchronously replicates the volume across multiple replicas stored on multiple nodes. The storage controller and replicas are themselves orchestrated using Kubernetes.

## Features

  * Enterprise-grade distributed block storage with no single point of failure
  * Incremental snapshot of block storage
  * Backup to secondary storage ([NFS](https://www.extrahop.com/resources/protocols/nfs/) or [S3](https://aws.amazon.com/s3/)-compatible object storage) built on efficient change block detection
  * Recurring snapshots and backups
  * Automated, non-disruptive upgrades. You can upgrade the entire Longhorn software stack without disrupting running storage volumes.
  * An intuitive GUI dashboard



* * *

© 2019-2026 Longhorn Authors | Documentation Distributed under CC-BY-4.0

  


© 2026 The Linux Foundation. All rights reserved. The Linux Foundation has registered trademarks and uses trademarks. For a list of trademarks of The Linux Foundation, please see our [Trademark Usage](https://www.linuxfoundation.org/trademark-usage/) page.

  
[](https://k3s.io/) [](https://harvesterhci.io/)

[The Longhorn Documentation](/docs/1.11.0/)

  * [What is Longhorn?](/docs/1.11.0/what-is-longhorn/)
  * [Architecture and Concepts](/docs/1.11.0/concepts/)
  * [Terminology](/docs/1.11.0/terminology/)
  * [Best Practices](/docs/1.11.0/best-practices/)
  * [Document Conventions](/docs/1.11.0/document-conventions/)
  * [Contributing](/docs/1.11.0/contributing/)



[Important Notes](/docs/1.11.0/important-notes/)

[Installation and Setup](/docs/1.11.0/deploy/)

[Quick Installation](/docs/1.11.0/deploy/install/)

  * [Install as a Rancher Apps & Marketplace](/docs/1.11.0/deploy/install/install-with-rancher/)
  * [Install with Kubectl](/docs/1.11.0/deploy/install/install-with-kubectl/)
  * [Install with Helm](/docs/1.11.0/deploy/install/install-with-helm/)
  * [Install with Helm Controller](/docs/1.11.0/deploy/install/install-with-helm-controller/)
  * [Install with Fleet](/docs/1.11.0/deploy/install/install-with-fleet/)
  * [Install with Flux](/docs/1.11.0/deploy/install/install-with-flux/)
  * [Install with ArgoCD](/docs/1.11.0/deploy/install/install-with-argocd/)
  * [Air Gap Installation](/docs/1.11.0/deploy/install/airgap/)

[Accessing the UI](/docs/1.11.0/deploy/accessing-the-ui/)

  * [Create an Ingress with Basic Authentication (nginx)](/docs/1.11.0/deploy/accessing-the-ui/longhorn-ingress/)
  * [Create an HTTPRoute with Gateway API](/docs/1.11.0/deploy/accessing-the-ui/longhorn-httproute/)

[Upgrade](/docs/1.11.0/deploy/upgrade/)

  * [Upgrading Longhorn Manager](/docs/1.11.0/deploy/upgrade/longhorn-manager/)
  * [Manually Upgrading Longhorn Engine](/docs/1.11.0/deploy/upgrade/upgrade-engine/)
  * [Automatically Upgrading Longhorn Engine](/docs/1.11.0/deploy/upgrade/auto-upgrade-engine/)
  * [Instance Manager Pods During Upgrade](/docs/1.11.0/deploy/upgrade/instance-manager-pods-during-upgrade/)

[Uninstall Longhorn](/docs/1.11.0/deploy/uninstall/)

[Maintenance and Upgrade](/docs/1.11.0/maintenance/)

  * [Node Maintenance and Kubernetes Upgrade Guide](/docs/1.11.0/maintenance/maintenance/)



[Nodes and Volumes](/docs/1.11.0/nodes-and-volumes/)

[Nodes](/docs/1.11.0/nodes-and-volumes/nodes/)

  * [Node Space Usage](/docs/1.11.0/nodes-and-volumes/nodes/node-space-usage/)
  * [Storage Tags](/docs/1.11.0/nodes-and-volumes/nodes/storage-tags/)
  * [Configuring Defaults for Nodes and Disks](/docs/1.11.0/nodes-and-volumes/nodes/default-disk-and-node-config/)
  * [Multiple Disk Support](/docs/1.11.0/nodes-and-volumes/nodes/multidisk/)
  * [Scheduling](/docs/1.11.0/nodes-and-volumes/nodes/scheduling/)
  * [Evicting Replicas on Disabled Disks or Nodes](/docs/1.11.0/nodes-and-volumes/nodes/disks-or-nodes-eviction/)
  * [Node Conditions](/docs/1.11.0/nodes-and-volumes/nodes/node-conditions/)

[Volumes](/docs/1.11.0/nodes-and-volumes/volumes/)

  * [Create Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/create-volumes/)
  * [Longhorn PVC Ownership and Permission](/docs/1.11.0/nodes-and-volumes/volumes/pvc-ownership-and-permission/)
  * [Delete Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/delete-volumes/)
  * [Detach Longhorn Volumes](/docs/1.11.0/nodes-and-volumes/volumes/detaching-volumes/)
  * [ReadWriteMany (RWX) Volume](/docs/1.11.0/nodes-and-volumes/volumes/rwx-volumes/)
  * [Use Longhorn Volume as an iSCSI Target](/docs/1.11.0/nodes-and-volumes/volumes/iscsi/)
  * [Volume Size](/docs/1.11.0/nodes-and-volumes/volumes/volume-size/)
  * [Viewing Workloads that Use a Volume](/docs/1.11.0/nodes-and-volumes/volumes/workload-identification/)
  * [Volume Expansion](/docs/1.11.0/nodes-and-volumes/volumes/expansion/)
  * [Trim Filesystem](/docs/1.11.0/nodes-and-volumes/volumes/trim-filesystem/)
  * [Volume Conditions](/docs/1.11.0/nodes-and-volumes/volumes/volume-conditions/)



[High Availability](/docs/1.11.0/high-availability/)

  * [Auto Balance Replicas](/docs/1.11.0/high-availability/auto-balance-replicas/)
  * [Data Locality](/docs/1.11.0/high-availability/data-locality/)
  * [Kubernetes Cluster Autoscaler Support (Experimental)](/docs/1.11.0/high-availability/k8s-cluster-autoscaler/)
  * [RWX Volume Fast Failover (Experimental)](/docs/1.11.0/high-availability/rwx-volume-fast-failover/)
  * [Volume Recovery](/docs/1.11.0/high-availability/recover-volume/)
  * [Node Failure Handling with Longhorn](/docs/1.11.0/high-availability/node-failure/)



[Backup and Restore](/docs/1.11.0/snapshots-and-backups/)

  * [Create a Snapshot](/docs/1.11.0/snapshots-and-backups/setup-a-snapshot/)
  * [Snapshot Space Management](/docs/1.11.0/snapshots-and-backups/snapshot-space-management/)
  * [Recurring Snapshots and Backups](/docs/1.11.0/snapshots-and-backups/scheduling-backups-and-snapshots/)
  * [Volume Clone Support](/docs/1.11.0/snapshots-and-backups/csi-volume-clone/)
  * [Disaster Recovery (DR) Volumes](/docs/1.11.0/snapshots-and-backups/setup-disaster-recovery-volumes/)



[Backup and Restore](/docs/1.11.0/snapshots-and-backups/backup-and-restore/)

  * [Setting a Backup Target](/docs/1.11.0/snapshots-and-backups/backup-and-restore/set-backup-target/)
  * [Configure The Block Size Of Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/configure-backup-block-size/)
  * [Create a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/create-a-backup/)
  * [Restore from a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-from-a-backup/)
  * [Restoring Volumes for Kubernetes StatefulSets](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-statefulset/)
  * [Restore Volume Recurring Jobs from a Backup](/docs/1.11.0/snapshots-and-backups/backup-and-restore/restore-recurring-jobs-from-a-backup/)
  * [Synchronize Backup Volumes Manually](/docs/1.11.0/snapshots-and-backups/backup-and-restore/synchronize_backup_volumes_manually/)

[CSI Snapshot Support](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/)

  * [Enable CSI Snapshot Support on a Cluster](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/enable-csi-snapshot-support/)
  * [CSI VolumeSnapshot Associated with Longhorn BackingImage](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backing-image/)
  * [CSI VolumeSnapshot Associated with Longhorn Snapshot](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-snapshot/)
  * [CSI VolumeSnapshot Associated with Longhorn Backup](/docs/1.11.0/snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backup/)



[Monitoring](/docs/1.11.0/monitoring/)

  * [Setting up Prometheus and Grafana to monitor Longhorn](/docs/1.11.0/monitoring/prometheus-and-grafana-setup/)
  * [Integrating Longhorn metrics into the Rancher monitoring system](/docs/1.11.0/monitoring/integrating-with-rancher-monitoring/)
  * [Longhorn Metrics for Monitoring](/docs/1.11.0/monitoring/metrics/)
  * [Disk Health Monitoring](/docs/1.11.0/monitoring/disk-heath/)
  * [Kubelet Volume Metrics Support](/docs/1.11.0/monitoring/kubelet-volume-metrics/)
  * [Longhorn Alert Rule Examples](/docs/1.11.0/monitoring/alert-rules-example/)



[Advanced Resources](/docs/1.11.0/advanced-resources/)

  * [Longhorn VolumeAttachment](/docs/1.11.0/advanced-resources/volumeattachment/)



[Deploy](/docs/1.11.0/advanced-resources/deploy/)

  * [Customizing Default Settings](/docs/1.11.0/advanced-resources/deploy/customizing-default-settings/)
  * [Taints and Tolerations](/docs/1.11.0/advanced-resources/deploy/taint-toleration/)
  * [Node Selector](/docs/1.11.0/advanced-resources/deploy/node-selector/)
  * [CSI Component Pod Anti-Affinity](/docs/1.11.0/advanced-resources/deploy/csi-pod-antiaffinity-preset/)
  * [Rancher Windows Cluster](/docs/1.11.0/advanced-resources/deploy/rancher_windows_cluster/)
  * [Priority Class](/docs/1.11.0/advanced-resources/deploy/priority-class/)
  * [Revision Counter](/docs/1.11.0/advanced-resources/deploy/revision_counter/)
  * [Storage Network](/docs/1.11.0/advanced-resources/deploy/storage-network/)

[OS/Distro Specific Configuration](/docs/1.11.0/advanced-resources/os-distro-specific/)

  * [Longhorn CSI on K3s](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-k3s/)
  * [Longhorn CSI on RKE and CoreOS](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-rke-and-coreos/)
  * [Longhorn CSI on GKE](/docs/1.11.0/advanced-resources/os-distro-specific/csi-on-gke/)
  * [OCP/OKD Support](/docs/1.11.0/advanced-resources/os-distro-specific/okd-support/)
  * [Container-Optimized OS (COS) Support](/docs/1.11.0/advanced-resources/os-distro-specific/container-optimized-os-support/)
  * [Talos Linux Support](/docs/1.11.0/advanced-resources/os-distro-specific/talos-linux-support/)

[Support Managed Kubernetes Service](/docs/1.11.0/advanced-resources/support-managed-k8s-service/)

  * [Manage Node-Group on AWS EKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-eks/)
  * [Manage Node-Group on Azure AKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-aks/)
  * [Manage Node-Group on GCP GKE](/docs/1.11.0/advanced-resources/support-managed-k8s-service/manage-node-group-on-gke/)
  * [Upgrade Kubernetes on AWS EKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-eks/)
  * [Upgrade Kubernetes on Azure AKS](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-aks/)
  * [Upgrade Kubernetes on GCP GKE](/docs/1.11.0/advanced-resources/support-managed-k8s-service/upgrade-k8s-on-gke/)

[Data Integrity](/docs/1.11.0/advanced-resources/data-integrity/)

  * [Snapshot Data Integrity Check](/docs/1.11.0/advanced-resources/data-integrity/snapshot-data-integrity-check/)

[Data Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/)

  * [Orphaned Data Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/orphaned-data-cleanup/)
  * [Orphaned Instance Cleanup](/docs/1.11.0/advanced-resources/data-cleanup/orphaned-instance-cleanup/)

[Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/)

  * [Fast Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/fast-replica-rebuilding/)
  * [Offline Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/offline-replica-rebuilding/)
  * [Scale Replica Rebuilding](/docs/1.11.0/advanced-resources/rebuilding/scale-replica-rebuilding/)

[Backing Image](/docs/1.11.0/advanced-resources/backing-image/)

  * [Backing Image](/docs/1.11.0/advanced-resources/backing-image/backing-image/)
  * [Backing Image Backup](/docs/1.11.0/advanced-resources/backing-image/backing-image-backup/)
  * [Backing Image Encryption](/docs/1.11.0/advanced-resources/backing-image/backing-image-encryption/)

[Security](/docs/1.11.0/advanced-resources/security/)

  * [Volume Encryption](/docs/1.11.0/advanced-resources/security/volume-encryption/)
  * [MTLS Support](/docs/1.11.0/advanced-resources/security/mtls-support/)

[Command Line Tool (longhornctl)](/docs/1.11.0/advanced-resources/longhornctl/)

  * [Install longhornctl](/docs/1.11.0/advanced-resources/longhornctl/install-longhornctl/)

[Containerized Data Importer (CDI)](/docs/1.11.0/advanced-resources/containerized-data-importer/)

  * [Longhorn with CDI Imports](/docs/1.11.0/advanced-resources/containerized-data-importer/containerized-data-importer/)

[Longhorn System Backup And Restore](/docs/1.11.0/advanced-resources/system-backup-restore/)

  * [Backup Longhorn System](/docs/1.11.0/advanced-resources/system-backup-restore/backup-longhorn-system/)
  * [Restore Longhorn System](/docs/1.11.0/advanced-resources/system-backup-restore/restore-longhorn-system/)
  * [Restore to a cluster contains data using Rancher snapshot](/docs/1.11.0/advanced-resources/system-backup-restore/restore-to-a-cluster-contains-data-using-rancher-snapshot/)
  * [Restore to a new cluster using Velero](/docs/1.11.0/advanced-resources/system-backup-restore/restore-to-a-new-cluster-using-velero/)

[Cluster Restore](/docs/1.11.0/advanced-resources/cluster-restore/)

  * [Restore cluster with a Rancher snapshot](/docs/1.11.0/advanced-resources/cluster-restore/rancher-cluster-restore/)

[Data Recovery](/docs/1.11.0/advanced-resources/data-recovery/)

  * [Identifying and Recovering from Data Errors](/docs/1.11.0/advanced-resources/data-recovery/data-error/)
  * [Exporting a Volume from a Single Replica](/docs/1.11.0/advanced-resources/data-recovery/export-from-replica/)
  * [Identifying Corrupted Replicas](/docs/1.11.0/advanced-resources/data-recovery/corrupted-replica/)
  * [Recovering from a Full Disk](/docs/1.11.0/advanced-resources/data-recovery/full-disk/)
  * [Recovering from a Longhorn Backup without System Installed](/docs/1.11.0/advanced-resources/data-recovery/recover-without-system/)

[CSI Driver Migration](/docs/1.11.0/advanced-resources/driver-migration/)

  * [Migrating from the Flexvolume Driver to CSI](/docs/1.11.0/advanced-resources/driver-migration/migrating-flexvolume/)



[References](/docs/1.11.0/references/)

  * [Settings](/docs/1.11.0/references/settings/)
  * [Storage Class Parameters](/docs/1.11.0/references/storage-class-parameters/)
  * [Python Client](/docs/1.11.0/references/longhorn-client-python/)
  * [Reference Setup, Performance, Scalability, and Sizing Guidelines](/docs/1.11.0/references/reference-setup-performance-scalability-and-sizing-guidelines/)
  * [Longhorn Networking](/docs/1.11.0/references/networking/)
  * [Examples](/docs/1.11.0/references/examples/)
  * [Helm Values](/docs/1.11.0/references/helm-values/)



[Troubleshoot](/docs/1.11.0/troubleshoot/)

  * [Troubleshooting Problems](/docs/1.11.0/troubleshoot/troubleshooting/)
  * [Support Bundle](/docs/1.11.0/troubleshoot/support-bundle/)



[V2 Data Engine (Technical Preview)](/docs/1.11.0/v2-data-engine/)

  * [Prerequisites](/docs/1.11.0/v2-data-engine/prerequisites/)
  * [Quick Start](/docs/1.11.0/v2-data-engine/quick-start/)
  * [Performance](/docs/1.11.0/v2-data-engine/performance/)
  * [Troubleshooting](/docs/1.11.0/v2-data-engine/troubleshooting/)



[Features](/docs/1.11.0/v2-data-engine/features/)

  * [V2 Volume Clone Support](/docs/1.11.0/v2-data-engine/features/volume-clone/)
  * [Replica Rebuild QoS](/docs/1.11.0/v2-data-engine/features/replica-rebuild-qos/)
  * [Configurable CPU Cores](/docs/1.11.0/v2-data-engine/features/configurable-cpu-cores/)
  * [Interrupt Mode Support](/docs/1.11.0/v2-data-engine/features/interrupt-mode/)
  * [Node Disk Support](/docs/1.11.0/v2-data-engine/features/node-disk-support/)
  * [Selective V2 Data Engine Activation](/docs/1.11.0/v2-data-engine/features/selective-v2-data-engine-activation/)
  * [UBLK Frontend Support](/docs/1.11.0/v2-data-engine/features/ublk-frontend-support/)
  * [V2 Volume Expansion](/docs/1.11.0/v2-data-engine/features/volume-expansion/)


  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
