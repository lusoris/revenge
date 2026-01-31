# Pro Git: Git Hooks

> Source: https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks
> Fetched: 2026-01-31T12:46:43.196879+00:00
> Content-Hash: a211a3fd16f55222
> Type: html

---

Chapters ▾

  1. ## 1\. [Getting Started](/book/en/v2/Getting-Started-About-Version-Control)

     1. 1.1 [About Version Control](/book/en/v2/Getting-Started-About-Version-Control)
     2. 1.2 [A Short History of Git](/book/en/v2/Getting-Started-A-Short-History-of-Git)
     3. 1.3 [What is Git?](/book/en/v2/Getting-Started-What-is-Git%3F)
     4. 1.4 [The Command Line](/book/en/v2/Getting-Started-The-Command-Line)
     5. 1.5 [Installing Git](/book/en/v2/Getting-Started-Installing-Git)
     6. 1.6 [First-Time Git Setup](/book/en/v2/Getting-Started-First-Time-Git-Setup)
     7. 1.7 [Getting Help](/book/en/v2/Getting-Started-Getting-Help)
     8. 1.8 [Summary](/book/en/v2/Getting-Started-Summary)
  2. ## 2\. [Git Basics](/book/en/v2/Git-Basics-Getting-a-Git-Repository)

     1. 2.1 [Getting a Git Repository](/book/en/v2/Git-Basics-Getting-a-Git-Repository)
     2. 2.2 [Recording Changes to the Repository](/book/en/v2/Git-Basics-Recording-Changes-to-the-Repository)
     3. 2.3 [Viewing the Commit History](/book/en/v2/Git-Basics-Viewing-the-Commit-History)
     4. 2.4 [Undoing Things](/book/en/v2/Git-Basics-Undoing-Things)
     5. 2.5 [Working with Remotes](/book/en/v2/Git-Basics-Working-with-Remotes)
     6. 2.6 [Tagging](/book/en/v2/Git-Basics-Tagging)
     7. 2.7 [Git Aliases](/book/en/v2/Git-Basics-Git-Aliases)
     8. 2.8 [Summary](/book/en/v2/Git-Basics-Summary)
  3. ## 3\. [Git Branching](/book/en/v2/Git-Branching-Branches-in-a-Nutshell)

     1. 3.1 [Branches in a Nutshell](/book/en/v2/Git-Branching-Branches-in-a-Nutshell)
     2. 3.2 [Basic Branching and Merging](/book/en/v2/Git-Branching-Basic-Branching-and-Merging)
     3. 3.3 [Branch Management](/book/en/v2/Git-Branching-Branch-Management)
     4. 3.4 [Branching Workflows](/book/en/v2/Git-Branching-Branching-Workflows)
     5. 3.5 [Remote Branches](/book/en/v2/Git-Branching-Remote-Branches)
     6. 3.6 [Rebasing](/book/en/v2/Git-Branching-Rebasing)
     7. 3.7 [Summary](/book/en/v2/Git-Branching-Summary)
  4. ## 4\. [Git on the Server](/book/en/v2/Git-on-the-Server-The-Protocols)

     1. 4.1 [The Protocols](/book/en/v2/Git-on-the-Server-The-Protocols)
     2. 4.2 [Getting Git on a Server](/book/en/v2/Git-on-the-Server-Getting-Git-on-a-Server)
     3. 4.3 [Generating Your SSH Public Key](/book/en/v2/Git-on-the-Server-Generating-Your-SSH-Public-Key)
     4. 4.4 [Setting Up the Server](/book/en/v2/Git-on-the-Server-Setting-Up-the-Server)
     5. 4.5 [Git Daemon](/book/en/v2/Git-on-the-Server-Git-Daemon)
     6. 4.6 [Smart HTTP](/book/en/v2/Git-on-the-Server-Smart-HTTP)
     7. 4.7 [GitWeb](/book/en/v2/Git-on-the-Server-GitWeb)
     8. 4.8 [GitLab](/book/en/v2/Git-on-the-Server-GitLab)
     9. 4.9 [Third Party Hosted Options](/book/en/v2/Git-on-the-Server-Third-Party-Hosted-Options)
     10. 4.10 [Summary](/book/en/v2/Git-on-the-Server-Summary)
  5. ## 5\. [Distributed Git](/book/en/v2/Distributed-Git-Distributed-Workflows)

     1. 5.1 [Distributed Workflows](/book/en/v2/Distributed-Git-Distributed-Workflows)
     2. 5.2 [Contributing to a Project](/book/en/v2/Distributed-Git-Contributing-to-a-Project)
     3. 5.3 [Maintaining a Project](/book/en/v2/Distributed-Git-Maintaining-a-Project)
     4. 5.4 [Summary](/book/en/v2/Distributed-Git-Summary)



  1. ## 6\. [GitHub](/book/en/v2/GitHub-Account-Setup-and-Configuration)

     1. 6.1 [Account Setup and Configuration](/book/en/v2/GitHub-Account-Setup-and-Configuration)
     2. 6.2 [Contributing to a Project](/book/en/v2/GitHub-Contributing-to-a-Project)
     3. 6.3 [Maintaining a Project](/book/en/v2/GitHub-Maintaining-a-Project)
     4. 6.4 [Managing an organization](/book/en/v2/GitHub-Managing-an-organization)
     5. 6.5 [Scripting GitHub](/book/en/v2/GitHub-Scripting-GitHub)
     6. 6.6 [Summary](/book/en/v2/GitHub-Summary)
  2. ## 7\. [Git Tools](/book/en/v2/Git-Tools-Revision-Selection)

     1. 7.1 [Revision Selection](/book/en/v2/Git-Tools-Revision-Selection)
     2. 7.2 [Interactive Staging](/book/en/v2/Git-Tools-Interactive-Staging)
     3. 7.3 [Stashing and Cleaning](/book/en/v2/Git-Tools-Stashing-and-Cleaning)
     4. 7.4 [Signing Your Work](/book/en/v2/Git-Tools-Signing-Your-Work)
     5. 7.5 [Searching](/book/en/v2/Git-Tools-Searching)
     6. 7.6 [Rewriting History](/book/en/v2/Git-Tools-Rewriting-History)
     7. 7.7 [Reset Demystified](/book/en/v2/Git-Tools-Reset-Demystified)
     8. 7.8 [Advanced Merging](/book/en/v2/Git-Tools-Advanced-Merging)
     9. 7.9 [Rerere](/book/en/v2/Git-Tools-Rerere)
     10. 7.10 [Debugging with Git](/book/en/v2/Git-Tools-Debugging-with-Git)
     11. 7.11 [Submodules](/book/en/v2/Git-Tools-Submodules)
     12. 7.12 [Bundling](/book/en/v2/Git-Tools-Bundling)
     13. 7.13 [Replace](/book/en/v2/Git-Tools-Replace)
     14. 7.14 [Credential Storage](/book/en/v2/Git-Tools-Credential-Storage)
     15. 7.15 [Summary](/book/en/v2/Git-Tools-Summary)
  3. ## 8\. [Customizing Git](/book/en/v2/Customizing-Git-Git-Configuration)

     1. 8.1 [Git Configuration](/book/en/v2/Customizing-Git-Git-Configuration)
     2. 8.2 [Git Attributes](/book/en/v2/Customizing-Git-Git-Attributes)
     3. 8.3 [Git Hooks](/book/en/v2/Customizing-Git-Git-Hooks)
     4. 8.4 [An Example Git-Enforced Policy](/book/en/v2/Customizing-Git-An-Example-Git-Enforced-Policy)
     5. 8.5 [Summary](/book/en/v2/Customizing-Git-Summary)
  4. ## 9\. [Git and Other Systems](/book/en/v2/Git-and-Other-Systems-Git-as-a-Client)

     1. 9.1 [Git as a Client](/book/en/v2/Git-and-Other-Systems-Git-as-a-Client)
     2. 9.2 [Migrating to Git](/book/en/v2/Git-and-Other-Systems-Migrating-to-Git)
     3. 9.3 [Summary](/book/en/v2/Git-and-Other-Systems-Summary)
  5. ## 10\. [Git Internals](/book/en/v2/Git-Internals-Plumbing-and-Porcelain)

     1. 10.1 [Plumbing and Porcelain](/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
     2. 10.2 [Git Objects](/book/en/v2/Git-Internals-Git-Objects)
     3. 10.3 [Git References](/book/en/v2/Git-Internals-Git-References)
     4. 10.4 [Packfiles](/book/en/v2/Git-Internals-Packfiles)
     5. 10.5 [The Refspec](/book/en/v2/Git-Internals-The-Refspec)
     6. 10.6 [Transfer Protocols](/book/en/v2/Git-Internals-Transfer-Protocols)
     7. 10.7 [Maintenance and Data Recovery](/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery)
     8. 10.8 [Environment Variables](/book/en/v2/Git-Internals-Environment-Variables)
     9. 10.9 [Summary](/book/en/v2/Git-Internals-Summary)



  1. ## A1. [Appendix A: Git in Other Environments](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Graphical-Interfaces)

     1. A1.1 [Graphical Interfaces](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Graphical-Interfaces)
     2. A1.2 [Git in Visual Studio](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-Visual-Studio)
     3. A1.3 [Git in Visual Studio Code](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-Visual-Studio-Code)
     4. A1.4 [Git in IntelliJ / PyCharm / WebStorm / PhpStorm / RubyMine](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-IntelliJ-/-PyCharm-/-WebStorm-/-PhpStorm-/-RubyMine)
     5. A1.5 [Git in Sublime Text](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-Sublime-Text)
     6. A1.6 [Git in Bash](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-Bash)
     7. A1.7 [Git in Zsh](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-Zsh)
     8. A1.8 [Git in PowerShell](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Git-in-PowerShell)
     9. A1.9 [Summary](/book/en/v2/Appendix-A:-Git-in-Other-Environments-Summary)
  2. ## A2. [Appendix B: Embedding Git in your Applications](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-Command-line-Git)

     1. A2.1 [Command-line Git](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-Command-line-Git)
     2. A2.2 [Libgit2](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-Libgit2)
     3. A2.3 [JGit](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-JGit)
     4. A2.4 [go-git](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-go-git)
     5. A2.5 [Dulwich](/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-Dulwich)
  3. ## A3. [Appendix C: Git Commands](/book/en/v2/Appendix-C:-Git-Commands-Setup-and-Config)

     1. A3.1 [Setup and Config](/book/en/v2/Appendix-C:-Git-Commands-Setup-and-Config)
     2. A3.2 [Getting and Creating Projects](/book/en/v2/Appendix-C:-Git-Commands-Getting-and-Creating-Projects)
     3. A3.3 [Basic Snapshotting](/book/en/v2/Appendix-C:-Git-Commands-Basic-Snapshotting)
     4. A3.4 [Branching and Merging](/book/en/v2/Appendix-C:-Git-Commands-Branching-and-Merging)
     5. A3.5 [Sharing and Updating Projects](/book/en/v2/Appendix-C:-Git-Commands-Sharing-and-Updating-Projects)
     6. A3.6 [Inspection and Comparison](/book/en/v2/Appendix-C:-Git-Commands-Inspection-and-Comparison)
     7. A3.7 [Debugging](/book/en/v2/Appendix-C:-Git-Commands-Debugging)
     8. A3.8 [Patching](/book/en/v2/Appendix-C:-Git-Commands-Patching)
     9. A3.9 [Email](/book/en/v2/Appendix-C:-Git-Commands-Email)
     10. A3.10 [External Systems](/book/en/v2/Appendix-C:-Git-Commands-External-Systems)
     11. A3.11 [Administration](/book/en/v2/Appendix-C:-Git-Commands-Administration)
     12. A3.12 [Plumbing Commands](/book/en/v2/Appendix-C:-Git-Commands-Plumbing-Commands)



2nd Edition 

# 8.3 Customizing Git - Git Hooks

## Git Hooks

Like many other Version Control Systems, Git has a way to fire off custom scripts when certain important actions occur. There are two groups of these hooks: client-side and server-side. Client-side hooks are triggered by operations such as committing and merging, while server-side hooks run on network operations such as receiving pushed commits. You can use these hooks for all sorts of reasons.

### Installing a Hook

The hooks are all stored in the `hooks` subdirectory of the Git directory. In most projects, that’s `.git/hooks`. When you initialize a new repository with `git init`, Git populates the hooks directory with a bunch of example scripts, many of which are useful by themselves; but they also document the input values of each script. All the examples are written as shell scripts, with some Perl thrown in, but any properly named executable scripts will work fine – you can write them in Ruby or Python or whatever language you are familiar with. If you want to use the bundled hook scripts, you’ll have to rename them; their file names all end with `.sample`.

To enable a hook script, put a file in the `hooks` subdirectory of your `.git` directory that is named appropriately (without any extension) and is executable. From that point forward, it should be called. We’ll cover most of the major hook filenames here.

### Client-Side Hooks

There are a lot of client-side hooks. This section splits them into committing-workflow hooks, email-workflow scripts, and everything else.

Note |  It’s important to note that client-side hooks are **not** copied when you clone a repository. If your intent with these scripts is to enforce a policy, you’ll probably want to do that on the server side; see the example in [An Example Git-Enforced Policy](/book/en/v2/ch00/_an_example_git_enforced_policy).  
---|---  
  
#### Committing-Workflow Hooks

The first four hooks have to do with the committing process.

The `pre-commit` hook is run first, before you even type in a commit message. It’s used to inspect the snapshot that’s about to be committed, to see if you’ve forgotten something, to make sure tests run, or to examine whatever you need to inspect in the code. Exiting non-zero from this hook aborts the commit, although you can bypass it with `git commit --no-verify`. You can do things like check for code style (run `lint` or something equivalent), check for trailing whitespace (the default hook does exactly this), or check for appropriate documentation on new methods.

The `prepare-commit-msg` hook is run before the commit message editor is fired up but after the default message is created. It lets you edit the default message before the commit author sees it. This hook takes a few parameters: the path to the file that holds the commit message so far, the type of commit, and the commit SHA-1 if this is an amended commit. This hook generally isn’t useful for normal commits; rather, it’s good for commits where the default message is auto-generated, such as templated commit messages, merge commits, squashed commits, and amended commits. You may use it in conjunction with a commit template to programmatically insert information.

The `commit-msg` hook takes one parameter, which again is the path to a temporary file that contains the commit message written by the developer. If this script exits non-zero, Git aborts the commit process, so you can use it to validate your project state or commit message before allowing a commit to go through. In the last section of this chapter, we’ll demonstrate using this hook to check that your commit message is conformant to a required pattern.

After the entire commit process is completed, the `post-commit` hook runs. It doesn’t take any parameters, but you can easily get the last commit by running `git log -1 HEAD`. Generally, this script is used for notification or something similar.

#### Email Workflow Hooks

You can set up three client-side hooks for an email-based workflow. They’re all invoked by the `git am` command, so if you aren’t using that command in your workflow, you can safely skip to the next section. If you’re taking patches over email prepared by `git format-patch`, then some of these may be helpful to you.

The first hook that is run is `applypatch-msg`. It takes a single argument: the name of the temporary file that contains the proposed commit message. Git aborts the patch if this script exits non-zero. You can use this to make sure a commit message is properly formatted, or to normalize the message by having the script edit it in place.

The next hook to run when applying patches via `git am` is `pre-applypatch`. Somewhat confusingly, it is run _after_ the patch is applied but before a commit is made, so you can use it to inspect the snapshot before making the commit. You can run tests or otherwise inspect the working tree with this script. If something is missing or the tests don’t pass, exiting non-zero aborts the `git am` script without committing the patch.

The last hook to run during a `git am` operation is `post-applypatch`, which runs after the commit is made. You can use it to notify a group or the author of the patch you pulled in that you’ve done so. You can’t stop the patching process with this script.

#### Other Client Hooks

The `pre-rebase` hook runs before you rebase anything and can halt the process by exiting non-zero. You can use this hook to disallow rebasing any commits that have already been pushed. The example `pre-rebase` hook that Git installs does this, although it makes some assumptions that may not match with your workflow.

The `post-rewrite` hook is run by commands that replace commits, such as `git commit --amend` and `git rebase` (though not by `git filter-branch`). Its single argument is which command triggered the rewrite, and it receives a list of rewrites on `stdin`. This hook has many of the same uses as the `post-checkout` and `post-merge` hooks.

After you run a successful `git checkout`, the `post-checkout` hook runs; you can use it to set up your working directory properly for your project environment. This may mean moving in large binary files that you don’t want source controlled, auto-generating documentation, or something along those lines.

The `post-merge` hook runs after a successful `merge` command. You can use it to restore data in the working tree that Git can’t track, such as permissions data. This hook can likewise validate the presence of files external to Git control that you may want copied in when the working tree changes.

The `pre-push` hook runs during `git push`, after the remote refs have been updated but before any objects have been transferred. It receives the name and location of the remote as parameters, and a list of to-be-updated refs through `stdin`. You can use it to validate a set of ref updates before a push occurs (a non-zero exit code will abort the push).

Git occasionally does garbage collection as part of its normal operation, by invoking `git gc --auto`. The `pre-auto-gc` hook is invoked just before the garbage collection takes place, and can be used to notify you that this is happening, or to abort the collection if now isn’t a good time.

### Server-Side Hooks

In addition to the client-side hooks, you can use a couple of important server-side hooks as a system administrator to enforce nearly any kind of policy for your project. These scripts run before and after pushes to the server. The pre hooks can exit non-zero at any time to reject the push as well as print an error message back to the client; you can set up a push policy that’s as complex as you wish.

#### `pre-receive`

The first script to run when handling a push from a client is `pre-receive`. It takes a list of references that are being pushed from stdin; if it exits non-zero, none of them are accepted. You can use this hook to do things like make sure none of the updated references are non-fast-forwards, or to do access control for all the refs and files they’re modifying with the push.

#### `update`

The `update` script is very similar to the `pre-receive` script, except that it’s run once for each branch the pusher is trying to update. If the pusher is trying to push to multiple branches, `pre-receive` runs only once, whereas `update` runs once per branch they’re pushing to. Instead of reading from stdin, this script takes three arguments: the name of the reference (branch), the SHA-1 that reference pointed to before the push, and the SHA-1 the user is trying to push. If the `update` script exits non-zero, only that reference is rejected; other references can still be updated.

#### `post-receive`

The `post-receive` hook runs after the entire process is completed and can be used to update other services or notify users. It takes the same stdin data as the `pre-receive` hook. Examples include emailing a list, notifying a continuous integration server, or updating a ticket-tracking system – you can even parse the commit messages to see if any tickets need to be opened, modified, or closed. This script can’t stop the push process, but the client doesn’t disconnect until it has completed, so be careful if you try to do anything that may take a long time.

Tip |  If you’re writing a script/hook that others will need to read, prefer the long versions of command-line flags; six months from now you’ll thank us.  
---|---  
  
[prev](/book/en/v2/Customizing-Git-Git-Attributes) | [next](/book/en/v2/Customizing-Git-An-Example-Git-Enforced-Policy)
