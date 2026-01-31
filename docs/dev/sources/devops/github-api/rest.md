# GitHub REST API Reference

> Source: https://docs.github.com/en/rest
> Fetched: 2026-01-31T16:06:04.316035+00:00
> Content-Hash: 9ce3210808f7b8a4
> Type: html

---

The REST API is now versioned. For more information, see "[About API versioning](/rest/overview/api-versions)."

## Start here

[View all ](/en/rest/guides)

  * ### [About the REST APIGet oriented to the REST API documentation.](/en/rest/about-the-rest-api/about-the-rest-api)
  * ### [Getting started with the REST APILearn how to use the GitHub REST API.](/en/rest/using-the-rest-api/getting-started-with-the-rest-api)
  * ### [Authenticating to the REST APIYou can authenticate to the REST API to access more endpoints and have a higher rate limit.](/en/rest/authentication/authenticating-to-the-rest-api)
  * ### [Best practices for using the REST APIFollow these best practices when using GitHub's API.](/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api)



## Popular

  * ### [Rate limits for the REST APILearn about REST API rate limits, how to avoid exceeding them, and what to do if you do exceed them.](/en/rest/using-the-rest-api/rate-limits-for-the-rest-api)
  * ### [Troubleshooting the REST APILearn how to diagnose and resolve common problems for the REST API.](/en/rest/using-the-rest-api/troubleshooting-the-rest-api)
  * ### [Scripting with the REST API and JavaScriptWrite a script using the Octokit.js SDK to interact with the REST API.](/en/rest/guides/scripting-with-the-rest-api-and-javascript)
  * ### [Keeping your API credentials secureFollow these best practices to keep your API credentials and tokens secure.](/en/rest/authentication/keeping-your-api-credentials-secure)



## Guides

  * ### [Delivering deploymentsUsing the Deployments REST API, you can build custom tooling that interacts with your server and a third-party app.](/en/rest/guides/delivering-deployments)
  * ### [Using the REST API to interact with checksYou can use the REST API to build GitHub Apps that run powerful checks against code changes in a repository. You can create apps that perform continuous integration, code linting, or code scanning services and provide detailed feedback on commits.](/en/rest/guides/using-the-rest-api-to-interact-with-checks)
  * ### [Using pagination in the REST APILearn how to navigate through paginated responses from the REST API.](/en/rest/using-the-rest-api/using-pagination-in-the-rest-api)



[Explore guides ](/en/rest/guides)

## All REST API docs

### [About the REST API](/en/rest/about-the-rest-api)

  * [About the REST API](/en/rest/about-the-rest-api/about-the-rest-api)
  * [Comparing GitHub's REST API and GraphQL API](/en/rest/about-the-rest-api/comparing-githubs-rest-api-and-graphql-api)
  * [API Versions](/en/rest/about-the-rest-api/api-versions)
  * [Breaking changes](/en/rest/about-the-rest-api/breaking-changes)
  * [About the OpenAPI description for the REST API](/en/rest/about-the-rest-api/about-the-openapi-description-for-the-rest-api)



### [Using the REST API](/en/rest/using-the-rest-api)

  * [Getting started with the REST API](/en/rest/using-the-rest-api/getting-started-with-the-rest-api)
  * [Rate limits for the REST API](/en/rest/using-the-rest-api/rate-limits-for-the-rest-api)
  * [Using pagination in the REST API](/en/rest/using-the-rest-api/using-pagination-in-the-rest-api)
  * [Libraries for the REST API](/en/rest/using-the-rest-api/libraries-for-the-rest-api)
  * [Best practices for using the REST API](/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api)
  * [Troubleshooting the REST API](/en/rest/using-the-rest-api/troubleshooting-the-rest-api)
  * [Timezones and the REST API](/en/rest/using-the-rest-api/timezones-and-the-rest-api)
  * [Using CORS and JSONP to make cross-origin requests](/en/rest/using-the-rest-api/using-cors-and-jsonp-to-make-cross-origin-requests)
  * [Issue event types](/en/rest/using-the-rest-api/issue-event-types)
  * [GitHub event types](/en/rest/using-the-rest-api/github-event-types)



### [Authenticating to the REST API](/en/rest/authentication)

  * [Authenticating to the REST API](/en/rest/authentication/authenticating-to-the-rest-api)
  * [Keeping your API credentials secure](/en/rest/authentication/keeping-your-api-credentials-secure)
  * [Endpoints available for GitHub App installation access tokens](/en/rest/authentication/endpoints-available-for-github-app-installation-access-tokens)
  * [Endpoints available for GitHub App user access tokens](/en/rest/authentication/endpoints-available-for-github-app-user-access-tokens)
  * [Endpoints available for fine-grained personal access tokens](/en/rest/authentication/endpoints-available-for-fine-grained-personal-access-tokens)
  * [Permissions required for GitHub Apps](/en/rest/authentication/permissions-required-for-github-apps)
  * [Permissions required for fine-grained personal access tokens](/en/rest/authentication/permissions-required-for-fine-grained-personal-access-tokens)



### [Guides](/en/rest/guides)

  * [Scripting with the REST API and JavaScript](/en/rest/guides/scripting-with-the-rest-api-and-javascript)
  * [Scripting with the REST API and Ruby](/en/rest/guides/scripting-with-the-rest-api-and-ruby)
  * [Discovering resources for a user](/en/rest/guides/discovering-resources-for-a-user)
  * [Delivering deployments](/en/rest/guides/delivering-deployments)
  * [Rendering data as graphs](/en/rest/guides/rendering-data-as-graphs)
  * [Working with comments](/en/rest/guides/working-with-comments)
  * [Building a CI server](/en/rest/guides/building-a-ci-server)
  * [Using the REST API to interact with your Git database](/en/rest/guides/using-the-rest-api-to-interact-with-your-git-database)
  * [Using the REST API to interact with checks](/en/rest/guides/using-the-rest-api-to-interact-with-checks)
  * [Encrypting secrets for the REST API](/en/rest/guides/encrypting-secrets-for-the-rest-api)



### [REST API endpoints for GitHub Actions](/en/rest/actions)

  * [REST API endpoints for GitHub Actions artifacts](/en/rest/actions/artifacts)
  * [REST API endpoints for GitHub Actions cache](/en/rest/actions/cache)
  * [GitHub-hosted runners](/en/rest/actions/hosted-runners)
  * [REST API endpoints for GitHub Actions OIDC](/en/rest/actions/oidc)
  * [REST API endpoints for GitHub Actions permissions](/en/rest/actions/permissions)
  * [REST API endpoints for GitHub Actions Secrets](/en/rest/actions/secrets)
  * [REST API endpoints for self-hosted runner groups](/en/rest/actions/self-hosted-runner-groups)
  * [REST API endpoints for self-hosted runners](/en/rest/actions/self-hosted-runners)
  * [REST API endpoints for GitHub Actions variables](/en/rest/actions/variables)
  * [REST API endpoints for workflow jobs](/en/rest/actions/workflow-jobs)
  * [REST API endpoints for workflow runs](/en/rest/actions/workflow-runs)
  * [REST API endpoints for workflows](/en/rest/actions/workflows)



### [REST API endpoints for activity](/en/rest/activity)

  * [REST API endpoints for events](/en/rest/activity/events)
  * [REST API endpoints for feeds](/en/rest/activity/feeds)
  * [REST API endpoints for notifications](/en/rest/activity/notifications)
  * [REST API endpoints for starring](/en/rest/activity/starring)
  * [REST API endpoints for watching](/en/rest/activity/watching)



### [REST API endpoints for apps](/en/rest/apps)

  * [REST API endpoints for GitHub Apps](/en/rest/apps/apps)
  * [REST API endpoints for GitHub App installations](/en/rest/apps/installations)
  * [REST API endpoints for GitHub Marketplace](/en/rest/apps/marketplace)
  * [REST API endpoints for OAuth authorizations](/en/rest/apps/oauth-applications)
  * [REST API endpoints for GitHub App webhooks](/en/rest/apps/webhooks)



### [REST API endpoints for billing](/en/rest/billing)

  * [Budgets](/en/rest/billing/budgets)
  * [Billing usage](/en/rest/billing/usage)



### [REST API endpoints for branches and their settings](/en/rest/branches)

  * [REST API endpoints for branches](/en/rest/branches/branches)
  * [REST API endpoints for protected branches](/en/rest/branches/branch-protection)



### [REST API endpoints for security campaigns](/en/rest/campaigns)

  * [REST API endpoints for security campaigns](/en/rest/campaigns/campaigns)



### [REST API endpoints for checks](/en/rest/checks)

  * [REST API endpoints for check runs](/en/rest/checks/runs)
  * [REST API endpoints for check suites](/en/rest/checks/suites)



### [REST API endpoints for GitHub Classroom](/en/rest/classroom)

  * [REST API endpoints for GitHub Classroom](/en/rest/classroom/classroom)



### [REST API endpoints for code scanning](/en/rest/code-scanning)

  * [REST API endpoints for code scanning](/en/rest/code-scanning/code-scanning)



### [REST API endpoints for code security settings](/en/rest/code-security)

  * [Configurations](/en/rest/code-security/configurations)



### [REST API endpoints for codes of conduct](/en/rest/codes-of-conduct)

  * [REST API endpoints for codes of conduct](/en/rest/codes-of-conduct/codes-of-conduct)



### [REST API endpoints for Codespaces](/en/rest/codespaces)

  * [REST API endpoints for Codespaces](/en/rest/codespaces/codespaces)
  * [REST API endpoints for Codespaces organizations](/en/rest/codespaces/organizations)
  * [REST API endpoints for Codespaces organization secrets](/en/rest/codespaces/organization-secrets)
  * [REST API endpoints for Codespaces machines](/en/rest/codespaces/machines)
  * [REST API endpoints for Codespaces repository secrets](/en/rest/codespaces/repository-secrets)
  * [REST API endpoints for Codespaces user secrets](/en/rest/codespaces/secrets)



### [REST API endpoints for collaborators](/en/rest/collaborators)

  * [REST API endpoints for collaborators](/en/rest/collaborators/collaborators)
  * [REST API endpoints for repository invitations](/en/rest/collaborators/invitations)



### [REST API endpoints for commits](/en/rest/commits)

  * [REST API endpoints for commits](/en/rest/commits/commits)
  * [REST API endpoints for commit comments](/en/rest/commits/comments)
  * [REST API endpoints for commit statuses](/en/rest/commits/statuses)



### [REST API endpoints for Copilot](/en/rest/copilot)

  * [REST API endpoints for Copilot metrics](/en/rest/copilot/copilot-metrics)
  * [REST API endpoints for Copilot user management](/en/rest/copilot/copilot-user-management)



### [Credentials](/en/rest/credentials)

  * [Revocation](/en/rest/credentials/revoke)



### [REST API endpoints for Dependabot](/en/rest/dependabot)

  * [REST API endpoints for Dependabot alerts](/en/rest/dependabot/alerts)
  * [REST API endpoints for Dependabot repository access](/en/rest/dependabot/repository-access)
  * [REST API endpoints for Dependabot secrets](/en/rest/dependabot/secrets)



### [REST API endpoints for the dependency graph](/en/rest/dependency-graph)

  * [REST API endpoints for dependency review](/en/rest/dependency-graph/dependency-review)
  * [REST API endpoints for dependency submission](/en/rest/dependency-graph/dependency-submission)
  * [REST API endpoints for software bill of materials (SBOM)](/en/rest/dependency-graph/sboms)



### [REST API endpoints for deploy keys](/en/rest/deploy-keys)

  * [REST API endpoints for deploy keys](/en/rest/deploy-keys/deploy-keys)



### [REST API endpoints for deployments](/en/rest/deployments)

  * [REST API endpoints for deployment branch policies](/en/rest/deployments/branch-policies)
  * [REST API endpoints for deployments](/en/rest/deployments/deployments)
  * [REST API endpoints for deployment environments](/en/rest/deployments/environments)
  * [REST API endpoints for protection rules](/en/rest/deployments/protection-rules)
  * [REST API endpoints for deployment statuses](/en/rest/deployments/statuses)



### [REST API endpoints for emojis](/en/rest/emojis)

  * [REST API endpoints for emojis](/en/rest/emojis/emojis)



### [Enterprise teams](/en/rest/enterprise-teams)

  * [REST API endpoints for enterprise team memberships](/en/rest/enterprise-teams/enterprise-team-members)
  * [REST API endpoints for enterprise team organizations](/en/rest/enterprise-teams/enterprise-team-organizations)
  * [REST API endpoints for enterprise teams](/en/rest/enterprise-teams/enterprise-teams)



### [REST API endpoints for gists and gist comments](/en/rest/gists)

  * [REST API endpoints for gists](/en/rest/gists/gists)
  * [REST API endpoints for gist comments](/en/rest/gists/comments)



### [REST API endpoints for Git database](/en/rest/git)

  * [REST API endpoints for Git blobs](/en/rest/git/blobs)
  * [REST API endpoints for Git commits](/en/rest/git/commits)
  * [REST API endpoints for Git references](/en/rest/git/refs)
  * [REST API endpoints for Git tags](/en/rest/git/tags)
  * [REST API endpoints for Git trees](/en/rest/git/trees)



### [REST API endpoints for gitignore](/en/rest/gitignore)

  * [REST API endpoints for gitignore](/en/rest/gitignore/gitignore)



### [REST API endpoints for interactions](/en/rest/interactions)

  * [REST API endpoints for organization interactions](/en/rest/interactions/orgs)
  * [REST API endpoints for repository interactions](/en/rest/interactions/repos)
  * [REST API endpoints for user interactions](/en/rest/interactions/user)



### [REST API endpoints for issues](/en/rest/issues)

  * [REST API endpoints for issue assignees](/en/rest/issues/assignees)
  * [REST API endpoints for issue comments](/en/rest/issues/comments)
  * [REST API endpoints for issue events](/en/rest/issues/events)
  * [REST API endpoints for issues](/en/rest/issues/issues)
  * [REST API endpoints for issue dependencies](/en/rest/issues/issue-dependencies)
  * [REST API endpoints for labels](/en/rest/issues/labels)
  * [REST API endpoints for milestones](/en/rest/issues/milestones)
  * [REST API endpoints for sub-issues](/en/rest/issues/sub-issues)
  * [REST API endpoints for timeline events](/en/rest/issues/timeline)



### [REST API endpoints for licenses](/en/rest/licenses)

  * [REST API endpoints for licenses](/en/rest/licenses/licenses)



### [REST API endpoints for Markdown](/en/rest/markdown)

  * [REST API endpoints for Markdown](/en/rest/markdown/markdown)



### [REST API endpoints for meta data](/en/rest/meta)

  * [REST API endpoints for meta data](/en/rest/meta/meta)



### [REST API endpoints for metrics](/en/rest/metrics)

  * [REST API endpoints for community metrics](/en/rest/metrics/community)
  * [REST API endpoints for repository statistics](/en/rest/metrics/statistics)
  * [REST API endpoints for repository traffic](/en/rest/metrics/traffic)



### [REST API endpoints for migrations](/en/rest/migrations)

  * [REST API endpoints for organization migrations](/en/rest/migrations/orgs)
  * [REST API endpoints for source imports](/en/rest/migrations/source-imports)
  * [REST API endpoints for user migrations](/en/rest/migrations/users)



### [Models](/en/rest/models)

  * [REST API endpoints for models catalog](/en/rest/models/catalog)
  * [REST API endpoints for model embeddings](/en/rest/models/embeddings)
  * [REST API endpoints for models inference](/en/rest/models/inference)



### [REST API endpoints for organizations](/en/rest/orgs)

  * [REST API endpoints for API Insights](/en/rest/orgs/api-insights)
  * [REST API endpoints for artifact metadata](/en/rest/orgs/artifact-metadata)
  * [REST API endpoints for artifact attestations](/en/rest/orgs/attestations)
  * [REST API endpoints for blocking users](/en/rest/orgs/blocking)
  * [REST API endpoints for custom properties](/en/rest/orgs/custom-properties)
  * [REST API endpoints for issue types](/en/rest/orgs/issue-types)
  * [REST API endpoints for organization members](/en/rest/orgs/members)
  * [REST API endpoints for network configurations](/en/rest/orgs/network-configurations)
  * [REST API endpoints for organization roles](/en/rest/orgs/organization-roles)
  * [REST API endpoints for organizations](/en/rest/orgs/orgs)
  * [REST API endpoints for outside collaborators](/en/rest/orgs/outside-collaborators)
  * [REST API endpoints for personal access tokens](/en/rest/orgs/personal-access-tokens)
  * [REST API endpoints for rule suites](/en/rest/orgs/rule-suites)
  * [REST API endpoints for rules](/en/rest/orgs/rules)
  * [REST API endpoints for security managers](/en/rest/orgs/security-managers)
  * [REST API endpoints for organization webhooks](/en/rest/orgs/webhooks)



### [REST API endpoints for packages](/en/rest/packages)

  * [REST API endpoints for packages](/en/rest/packages/packages)



### [REST API endpoints for GitHub Pages](/en/rest/pages)

  * [REST API endpoints for GitHub Pages](/en/rest/pages/pages)



### [Private registries](/en/rest/private-registries)

  * [Organization configurations](/en/rest/private-registries/organization-configurations)



### [Projects](/en/rest/projects)

  * [REST API endpoints for draft Project items](/en/rest/projects/drafts)
  * [REST API endpoints for Project fields](/en/rest/projects/fields)
  * [REST API endpoints for Project items](/en/rest/projects/items)
  * [REST API endpoints for Projects](/en/rest/projects/projects)
  * [REST API endpoints for Project views](/en/rest/projects/views)



### [REST API endpoints for pull requests](/en/rest/pulls)

  * [REST API endpoints for pull requests](/en/rest/pulls/pulls)
  * [REST API endpoints for pull request review comments](/en/rest/pulls/comments)
  * [REST API endpoints for review requests](/en/rest/pulls/review-requests)
  * [REST API endpoints for pull request reviews](/en/rest/pulls/reviews)



### [REST API endpoints for rate limits](/en/rest/rate-limit)

  * [REST API endpoints for rate limits](/en/rest/rate-limit/rate-limit)



### [REST API endpoints for reactions](/en/rest/reactions)

  * [REST API endpoints for reactions](/en/rest/reactions/reactions)



### [REST API endpoints for releases and release assets](/en/rest/releases)

  * [REST API endpoints for releases](/en/rest/releases/releases)
  * [REST API endpoints for release assets](/en/rest/releases/assets)



### [REST API endpoints for repositories](/en/rest/repos)

  * [REST API endpoints for repository attestations](/en/rest/repos/attestations)
  * [REST API endpoints for repository autolinks](/en/rest/repos/autolinks)
  * [REST API endpoints for repository contents](/en/rest/repos/contents)
  * [REST API endpoints for custom properties](/en/rest/repos/custom-properties)
  * [REST API endpoints for forks](/en/rest/repos/forks)
  * [REST API endpoints for repositories](/en/rest/repos/repos)
  * [REST API endpoints for rule suites](/en/rest/repos/rule-suites)
  * [REST API endpoints for rules](/en/rest/repos/rules)
  * [REST API endpoints for repository tags](/en/rest/repos/tags)
  * [REST API endpoints for repository webhooks](/en/rest/repos/webhooks)



### [REST API endpoints for search](/en/rest/search)

  * [REST API endpoints for search](/en/rest/search/search)



### [REST API endpoints for secret scanning](/en/rest/secret-scanning)

  * [REST API endpoints for secret scanning push protection](/en/rest/secret-scanning/push-protection)
  * [REST API endpoints for secret scanning](/en/rest/secret-scanning/secret-scanning)



### [REST API endpoints for security advisories](/en/rest/security-advisories)

  * [REST API endpoints for global security advisories](/en/rest/security-advisories/global-advisories)
  * [REST API endpoints for repository security advisories](/en/rest/security-advisories/repository-advisories)



### [REST API endpoints for teams](/en/rest/teams)

  * [REST API endpoints for team members](/en/rest/teams/members)
  * [REST API endpoints for teams](/en/rest/teams/teams)



### [REST API endpoints for users](/en/rest/users)

  * [REST API endpoints for artifact attestations](/en/rest/users/attestations)
  * [REST API endpoints for blocking users](/en/rest/users/blocking)
  * [REST API endpoints for emails](/en/rest/users/emails)
  * [REST API endpoints for followers](/en/rest/users/followers)
  * [REST API endpoints for GPG keys](/en/rest/users/gpg-keys)
  * [REST API endpoints for Git SSH keys](/en/rest/users/keys)
  * [REST API endpoints for social accounts](/en/rest/users/social-accounts)
  * [REST API endpoints for SSH signing keys](/en/rest/users/ssh-signing-keys)
  * [REST API endpoints for users](/en/rest/users/users)


  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
