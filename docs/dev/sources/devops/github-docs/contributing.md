# GitHub CONTRIBUTING Guide

> Source: https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/setting-guidelines-for-repository-contributors
> Fetched: 2026-02-01T11:53:21.880425+00:00
> Content-Hash: 39ae807d0d2de3af
> Type: html

---

# Setting guidelines for repository contributors

You can create guidelines to communicate how people should contribute to your project.

View page as Markdown

## In this article

## About contributing guidelines

To help your project contributors do good work, you can add a file with contribution guidelines to your project repository's root, `docs`, or `.github` folder. When someone opens a pull request or creates an issue, they will see a link to that file. The link to the contributing guidelines also appears on your repository's `contribute` page. For an example of a `contribute` page, see [github/docs/contribute](https://github.com/github/docs/contribute).

If your repository includes a `CONTRIBUTING.md` file, GitHub also surfaces it in two other places to make it easier for contributors to discover:

- A " Contributing" tab in the repository overview (next to the " README" and " Code of conduct")
- A "Contributing" link in the repository sidebar

For the repository owner, contribution guidelines are a way to communicate how people should contribute.

For contributors, the guidelines help them verify that they're submitting well-formed pull requests and opening useful issues.

For both owners and contributors, contribution guidelines save time and hassle caused by improperly created pull requests or issues that have to be rejected and re-submitted.

You can create default contribution guidelines for your organization or personal account. For more information, see [Creating a default community health file](/en/communities/setting-up-your-project-for-healthy-contributions/creating-a-default-community-health-file).

Tip

Repository maintainers can set specific guidelines for issues by creating an issue or pull request template for the repository. For more information, see [About issue and pull request templates](/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates).

## Adding a `CONTRIBUTING.md` file

  1. On GitHub, navigate to the main page of the repository.

  2. Above the list of files, select the **Add file** dropdown menu, then click **Create new file**.

Alternatively, you can click  in the file tree view on the left.

  1. Decide whether to store your contributing guidelines in your repository's root, `docs`, or `.github` directory. Then, in the filename field, type the name and extension for the file. Contributing guidelines filenames are not case sensitive. Files are rendered in rich text format if the file extension is in a supported format. For more information, see [Working with non-code files](/en/repositories/working-with-files/using-files/working-with-non-code-files#rendering-differences-in-prose-documents).

     - To make your contributing guidelines visible in the repository's root directory, type _CONTRIBUTING_.
     - To make your contributing guidelines visible in the repository's `docs` directory, type _docs/_ to create the new directory, then _CONTRIBUTING_.
     - If a repository contains more than one _CONTRIBUTING_ file, then the file shown in links is chosen from locations in the following order: the `.github` directory, then the repository's root directory, and finally the `docs` directory.
  2. In the new file, add contribution guidelines. These could include:

     - Steps for creating good issues or pull requests.
     - Links to external documentation, mailing lists, or a code of conduct.
     - Community and behavioral expectations.
  3. Click **Commit changes...**

  4. In the "Commit message" field, type a short, meaningful commit message that describes the change you made to the file. You can attribute the commit to more than one author in the commit message. For more information, see [Creating a commit with multiple authors](/en/pull-requests/committing-changes-to-your-project/creating-and-editing-commits/creating-a-commit-with-multiple-authors).

  5. Below the commit message fields, decide whether to add your commit to the current branch or to a new branch. If your current branch is the default branch, you should choose to create a new branch for your commit and then create a pull request. For more information, see [Creating a pull request](/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request).

  6. Click **Commit changes** or **Propose changes**.

## Examples of contribution guidelines

If you're stumped, here are some good examples of contribution guidelines:

- The GitHub Docs [contribution guidelines](/en/contributing).
- The Ruby on Rails [contribution guidelines](https://github.com/rails/rails/blob/main/CONTRIBUTING.md).
- The Open Government [contribution guidelines](https://github.com/opengovernment/opengovernment/blob/master/CONTRIBUTING.md).

## Further reading

- The Open Source Guides' section [Starting an Open Source Project](https://opensource.guide/starting-a-project/)
- [GitHub Skills](https://skills.github.com/)
- [Adding a license to a repository](/en/communities/setting-up-your-project-for-healthy-contributions/adding-a-license-to-a-repository)

  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
