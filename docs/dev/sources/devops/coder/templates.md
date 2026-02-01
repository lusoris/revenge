# Coder Templates

> Source: https://coder.com/docs/templates
> Fetched: 2026-02-01T11:52:46.194133+00:00
> Content-Hash: 36ac6c178a15b080
> Type: html

---

[Home](/docs "Home")[Administration](/docs/admin "Administration")Templates

Templates are written in [Terraform](https://developer.hashicorp.com/terraform/intro) and define the underlying infrastructure that all Coder workspaces run on.

The "Starter Templates" page within the Coder dashboard.

## Learn the concepts

While templates are written in standard Terraform, it's important to learn the Coder-specific concepts behind templates. The best way to learn the concepts is by [creating a basic template from scratch](/docs/tutorials/template-from-scratch). If you are unfamiliar with Terraform, see [Hashicorp's Tutorials](https://developer.hashicorp.com/terraform/tutorials) for common cloud providers.

## Starter templates

After learning the basics, use starter templates to import a template with sensible defaults for popular platforms (e.g. AWS, Kubernetes, Docker, etc). Docs: [Create a template from a starter template](/docs/admin/templates/creating-templates#from-a-starter-template).

## Extending templates

It's often necessary to extend the template to make it generally useful to end users. Common modifications are:

- Your image(s) (e.g. a Docker image with languages and tools installed). Docs: [Image management](/docs/admin/templates/managing-templates/image-management).
- Additional parameters (e.g. disk size, instance type, or region). Docs: [Template parameters](/docs/admin/templates/extending-templates/parameters).
- Additional IDEs (e.g. JetBrains) or features (e.g. dotfiles, RDP). Docs: [Adding IDEs and features](/docs/admin/templates/extending-templates).

Learn more about the various ways you can [extend your templates](/docs/admin/templates/extending-templates).

## Best Practices

We recommend starting with a universal template that can be used for basic tasks. As your Coder deployment grows, you can create more templates to meet the needs of different teams.

- [Image management](/docs/admin/templates/managing-templates/image-management): Learn how to create and publish images for use within Coder workspaces & templates.
- [Dev Containers integration](/docs/admin/integrations/devcontainers/integration): Enable native dev containers support using `@devcontainers/cli` and Docker.
- [Envbuilder](/docs/admin/integrations/devcontainers/envbuilder): Alternative approach for environments without Docker access.
- [Template hardening](/docs/admin/templates/extending-templates/resource-persistence#-bulletproofing): Configure your template to prevent certain resources from being destroyed (e.g. user disks).
- [Manage templates with Ci/Cd pipelines](/docs/admin/templates/managing-templates/change-management): Learn how to source control your templates and use GitOps to ensure template changes are reviewed and tested.
- [Permissions and Policies](/docs/admin/templates/template-permissions): Control who may access and modify your template.
- [External Workspaces](/docs/admin/templates/managing-templates/external-workspaces): Learn how to connect your existing infrastructure to Coder workspaces.

##### [Creating TemplatesLearn how to create templates with Terraform](/docs/admin/templates/creating-templates)

##### [Managing TemplatesLearn how to manage templates and best practices](/docs/admin/templates/managing-templates)

##### [Extending TemplatesLearn best practices in extending templates](/docs/admin/templates/extending-templates)

##### [Open in CoderOpen workspaces in Coder](/docs/admin/templates/open-in-coder)

##### [Permissions & PoliciesLearn how to create templates with Terraform](/docs/admin/templates/template-permissions)

##### [Troubleshooting TemplatesLearn how to troubleshoot template issues](/docs/admin/templates/troubleshooting)
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
