---
title: "Upgrade Guide"
linkTitle: "Upgrade Guide"
weight: 10
description: >
  Learn about upgrading Karpenter
---

Karpenter is a controller that runs in your cluster, but it is not tied to a specific Kubernetes version, as the Cluster Autoscaler is.
Use your existing upgrade mechanisms to upgrade your core add-ons in Kubernetes and keep Karpenter up to date on bug fixes and new features.

To make upgrading easier we aim to minimize introduction of breaking changes with the followings:

# Compatibility issues

To make upgrading easier, we aim to minimize the introduction of breaking changes with the followings components:

* Provisioner API
* Helm Chart

We try to maintain compatibility with:

* The application itself
* The documentation of the application

When we introduce a breaking change, we do so only as described in this document.

Karpenter follows [Semantic Versioning 2.0.0](https://semver.org/) in its release versions, while in
major version zero (v0.y.z) [anything may change at any time](https://semver.org/#spec-item-4).
However, to further protect users during this phase we will only introduce breaking changes in minor releases (releases that increment y in x.y.z).
Note this does not mean every minor upgrade has a breaking change as we will also increment the
minor version when we release a new feature.

Users should therefore check to see if there is a breaking change every time they are upgrading to a new minor version.

## How Do We Break Incompatibility?

When there is a breaking change we will:

* Increment the minor version when in major version 0
* Add a permanent separate section named `upgrading to vx.y.z+` under [released upgrade notes](#released-upgrade-notes)
  clearly explaining the breaking change and what needs to be done on the user side to ensure a safe upgrade
* Add the sentence “This is a breaking change, please refer to the above link for upgrade instructions” to the top of the release notes and in all our announcements

## How Do We Find Incompatibilities

Besides the peer review process for all changes to the code base we also do the followings in order to find
incompatibilities:
* (To be implemented) To check the compatibility of the application, we will automate tests for installing, uninstalling, upgrading from an older version, and downgrading to an older version
* (To be implemented) To check the compatibility of the documentation with the application, we will turn the commands in our documentation into scripts that we can automatically run

## Nightly Builds

(To be implemented) Every night we will build and release everything that has been checked into the source code.
This enables us to detect problems including breaking changes and potential drifts in our external dependencies sooner than we otherwise would.
It also allows some advanced Karpenter users who have their own nightly builds to test the upcoming changes before they are released.

## Release Candidates

We consider having release candidates for major and important minor versions. Our release candidates are tagged like `vx.y.z-rc.0`, `vx.y.z-rc.1`. The release candidate will then graduate to `vx.y.z` as a normal stable release.
By adopting this practice we allow our users who are early adopters to test out new releases before they are available to the wider community, thereby providing us with early feedback resulting in more stable releases.

## Security Patches

While we are in major version 0 we will not release security patches to older versions.
Rather we provide the patches in the latest versions.
When at major version 1 we will have an EOL (end of life) policy where we provide security patches
for a subset of older versions and deprecate the others.

# Released Upgrade Notes

## Upgrading to v0.13.0+
* v0.13.0 introduces a new CRD named `AWSNodeTemplate` which can be used to specify AWS Cloud Provider parameters. Everything that was previously specified under `spec.provider` in the Provisioner resource, can now be specified in the spec of the new resource. The use of `spec.provider` is deprecated but will continue to function to maintain backwards compatibility for the current API version (v1alpha5) of the Provisioner resource. v0.13.0 also introduces support for custom user data that doesn't require the use of a custom launch template. The user data can be specified in-line in the AWSNodeTemplate resource. Read the [UserData documentation here](../aws/user-data) to get started.

  If you are upgrading from v0.10.1 - v0.11.1, a new CRD `awsnodetemplate` was added. In v0.12.0, this crd was renamed to `awsnodetemplates`. Since helm does not manage the lifecycle of CRDs, you will need to perform a few manual steps for this CRD upgrade:
  1. Make sure any `awsnodetemplate` manifests are saved somewhere so that they can be reapplied to the cluster.
  2. `kubectl delete crd awsnodetemplate`
  3. `kubectl apply -f https://raw.githubusercontent.com/aws/karpenter/v0.13.0/charts/karpenter/crds/karpenter.k8s.aws_awsnodetemplates.yaml`
  4. Perform the Karpenter upgrade to v0.13.x, which will install the new `awsnodetemplates` CRD.
  5. Reapply the `awsnodetemplate` manifests you saved from step 1, if applicable. 
* v0.13.0 also adds EC2/spot price fetching to Karpenter to allow making more accurate decisions regarding node deployments.  Our getting started guide documents this, but if you are upgrading Karpenter you will need to modify your Karpenter controller policy to add the `pricing:GetProducts` and `ec2:DescribeSpotPriceHistory` permissions.

## Upgrading to v0.12.0+
* v0.12.0 adds an OwnerReference to each Node created by a provisioner. Previously, deleting a provisioner would orphan nodes. Now, deleting a provisioner will cause Kubernetes [cascading delete](https://kubernetes.io/docs/concepts/architecture/garbage-collection/#cascading-deletion) logic to gracefully terminate the nodes using the Karpenter node finalizer. You may still orphan nodes by removing the owner reference.
* If you are upgrading from v0.10.1 - v0.11.1, a new CRD `awsnodetemplate` was added. In v0.12.0, this crd was renamed to `awsnodetemplates`. Since helm does not manage the lifecycle of CRDs, you will need to perform a few manual steps for this CRD upgrade:
  1. Make sure any `awsnodetemplate` manifests are saved somewhere so that they can be reapplied to the cluster.
  2. `kubectl delete crd awsnodetemplate`
  3. `kubectl apply -f https://raw.githubusercontent.com/aws/karpenter/v0.12.1/charts/karpenter/crds/karpenter.k8s.aws_awsnodetemplates.yaml`
  4. Perform the Karpenter upgrade to v0.12.x, which will install the new `awsnodetemplates` CRD.
  5. Reapply the `awsnodetemplate` manifests you saved from step 1, if applicable. 


## Upgrading to v0.11.0+

v0.11.0 changes the way that the `vpc.amazonaws.com/pod-eni` resource is reported.  Instead of being reported for all nodes that could support the resources regardless of if the cluster is configured to support it, it is now controlled by a command line flag or environment variable. The parameter defaults to false and must be set if your cluster uses [security groups for pods](https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html).  This can be enabled by setting the environment variable `AWS_ENABLE_POD_ENI` to true via the helm value `controller.env`.

Other extended resources must be registered on nodes by their respective device plugins which are typically installed as DaemonSets (e.g. the `nvidia.com/gpu` resource will be registered by the [NVIDIA device plugin](https://github.com/NVIDIA/k8s-device-plugin). Previously, Karpenter would register these resources on nodes at creation and they would be zeroed out by `kubelet` at startup.  By allowing the device plugins to register the resources, pods will not bind to the nodes before any device plugin initialization has occurred.

## Upgrading to v0.10.0+

v0.10.0 adds a new field, `startupTaints` to the provisioner spec.  Standard Helm upgrades [do not upgrade CRDs](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/#some-caveats-and-explanations) so the  field will not be available unless the CRD is manually updated.  This can be performed prior to the standard upgrade by applying the new CRD manually:

```shell
kubectl replace -f https://raw.githubusercontent.com/aws/karpenter/v0.10.0/charts/karpenter/crds/karpenter.sh_provisioners.yaml
```

📝 If you don't perform this manual CRD update, Karpenter will work correctly except for rejecting the creation/update of provisioners that use `startupTaints`.

## Upgrading to v0.6.2+

If using Helm, the variable names have changed for the cluster's name and endpoint. You may need to update any configuration
that sets the old variable names.

- `controller.clusterName` is now `clusterName`
- `controller.clusterEndpoint` is now `clusterEndpoint`