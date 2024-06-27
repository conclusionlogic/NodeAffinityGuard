# Chart README

## Overview

This Helm Chart deploys the node-affinity-guard application on Kubernetes. It includes configurations for deployment, DaemonSet, service accounts, RBAC rules, and more.

## Chart Details

- **Chart Name:** node-affinity-guard
- **Version:** 0.0.1
- **App Version:** 0.0.1
- **Maintainer:**
  - Name: conclusionlogic

## Components

1. **_helpers.tpl**
   - Contains template functions for generating names, labels, and images used in the chart.
   - Relevant functions: `node-affinity-guard.name`, `node-affinity-guard.fullname`, `node-affinity-guard.chart`, `node-affinity-guard.image`, `node-affinity-guard.labels`, `node-affinity-guard.selectorLabels`, `node-affinity-guard.serviceAccountName`.

2. **values.yaml**
   - Defines configuration values for deployment, image settings, service accounts, RBAC, resources, node selectors, tolerations, and affinity.

3. **deployment.yaml**
   - Configures a Deployment resource with specified replicas, labels, annotations, and container settings.

4. **daemonset.yaml**
   - Configures a DaemonSet resource with labels, annotations, container settings, and pod security context.

5. **clusterrole.yaml**
   - Defines a ClusterRole with rules for accessing resources like nodes and deployments.

6. **clusterrolebinding.yaml**
   - Sets up a ClusterRoleBinding linking the ClusterRole to a ServiceAccount.

7. **serviceaccount.yaml**
   - Creates a ServiceAccount with specified labels and annotations.

## Usage

To deploy the node-affinity-guard application using this Helm Chart, follow the standard Helm deployment process.

For detailed configurations and customisation options, refer to the specific template files in the `templates` directory.
