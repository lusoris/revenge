# Mermaid.js Documentation

> Source: https://mermaid.js.org/intro/
> Fetched: 2026-02-01T11:48:06.358768+00:00
> Content-Hash: 0ee2e261e41ab729
> Type: html

---

# About Mermaid ​

**Mermaid lets you create diagrams and visualizations using text and code.**

It is a JavaScript based diagramming and charting tool that renders Markdown-inspired text definitions to create and modify diagrams dynamically.

> If you are familiar with Markdown you should have no problem learning [Mermaid's Syntax](./syntax-reference.html).

[](https://github.com/mermaid-js/mermaid/actions/workflows/build.yml)[](https://www.npmjs.com/package/mermaid)[](https://bundlephobia.com/package/mermaid)[](https://coveralls.io/github/mermaid-js/mermaid?branch=master)[](https://www.jsdelivr.com/package/npm/mermaid)[](https://www.npmjs.com/package/mermaid)[](https://discord.gg/sKeNQX4Wtj)[](https://twitter.com/mermaidjs_)

[](https://mermaid-js.github.io/mermaid/landing/)

Mermaid is a JavaScript based diagramming and charting tool that uses Markdown-inspired text definitions and a renderer to create and modify complex diagrams. The main purpose of Mermaid is to help documentation catch up with development.

> Doc-Rot is a Catch-22 that Mermaid helps to solve.

Diagramming and documentation costs precious developer time and gets outdated quickly. But not having diagrams or docs ruins productivity and hurts organizational learning.  
Mermaid addresses this problem by enabling users to create easily modifiable diagrams, it can also be made part of production scripts (and other pieces of code).  
  
Mermaid allows even non-programmers to easily create detailed and diagrams through the [Mermaid Live Editor](https://mermaid.live/).  
[Tutorials](./../ecosystem/tutorials.html) has video tutorials.

Use Mermaid with your favorite applications, check out the list of [Community Integrations](./../ecosystem/integrations-community.html).

For a more detailed introduction to Mermaid and some of its more basic uses, look to the [Beginner's Guide](./../intro/getting-started.html) and [Usage](./../config/usage.html).

🌐 [CDN](https://www.jsdelivr.com/package/npm/mermaid) | 📖 [Documentation](https://mermaidjs.github.io) | 🙌 [Contribution](./../community/contributing.html) | 🔌 [Plug-Ins](./../ecosystem/integrations-community.html)

> 🖖 Keep a steady pulse: [mermaid needs more Collaborators](https://github.com/mermaid-js/mermaid/issues/866).

🏆 **Mermaid was nominated and won the[JS Open Source Awards (2019)](https://osawards.com/javascript/#nominees) in the category "The most exciting use of technology"!!!**

**Thanks to all involved, people committing pull requests, people answering questions and special thanks to Tyler Long who is helping me maintain the project 🙏**

Our PR Visual Regression Testing is powered by [Argos](https://argos-ci.com/?utm_source=mermaid&utm_campaign=oss) with their generous Open Source plan. It makes the process of reviewing PRs with visual changes a breeze.

[](https://argos-ci.com?utm_source=mermaid&utm_campaign=oss)

In our release process we rely heavily on visual regression tests using [applitools](https://applitools.com/). Applitools is a great service which has been easy to use and integrate with our tests.

[](https://applitools.com/)

## Diagram Types ​

### [Flowchart](./../syntax/flowchart.html?id=flowcharts-basic-syntax) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Sequence diagram](./../syntax/sequenceDiagram.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Gantt diagram](./../syntax/gantt.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Class diagram](./../syntax/classDiagram.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Git graph](./../syntax/gitgraph.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Entity Relationship Diagram - ❗ experimental](./../syntax/entityRelationshipDiagram.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [User Journey Diagram](./../syntax/userJourney.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [Quadrant Chart](./../syntax/quadrantChart.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

### [XY Chart](./../syntax/xyChart.html) ​

##### Code:

mermaid


Ctrl + Enter|Run ▶

## Installation ​

**In depth guides and examples can be found at[Getting Started](./getting-started.html) and [Usage](./../config/usage.html).**

**It would also be helpful to learn more about mermaid's[Syntax](./syntax-reference.html).**

### CDN ​
    
    
    https://cdn.jsdelivr.net/npm/mermaid@<version>/dist/

To select a version:

Replace `<version>` with the desired version number.

Latest Version: <https://cdn.jsdelivr.net/npm/mermaid@11>

## Deploying Mermaid ​

To Deploy Mermaid:

  1. You will need to install node v16, which would have npm
  2. Install mermaid 
     * NPM: `npm i mermaid`
     * Yarn: `yarn add mermaid`
     * Pnpm: `pnpm add mermaid`



### [Mermaid API](./../config/setup/README.html): ​

**To deploy mermaid without a bundler, insert a`script` tag with an absolute address and a `mermaid.initialize` call into the HTML using the following example:**

html
    
    
    <script type="module">
      import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs';
      mermaid.initialize({ startOnLoad: true });
    </script>

**Doing so commands the mermaid parser to look for the`<div>` or `<pre>` tags with `class="mermaid"`. From these tags, mermaid tries to read the diagram/chart definitions and render them into SVG charts.**

**Examples can be found in** [Other examples](./../syntax/examples.html)

## Sibling projects ​

  * [Mermaid Live Editor](https://github.com/mermaid-js/mermaid-live-editor)
  * [Mermaid CLI](https://github.com/mermaid-js/mermaid-cli)
  * [Mermaid Tiny](https://github.com/mermaid-js/mermaid/tree/develop/packages/tiny)
  * [Mermaid Webpack Demo](https://github.com/mermaidjs/mermaid-webpack-demo)
  * [Mermaid Parcel Demo](https://github.com/mermaidjs/mermaid-parcel-demo)



## Request for Assistance ​

Things are piling up and I have a hard time keeping up. It would be great if we could form a core team of developers to cooperate with the future development of mermaid.

As part of this team you would get write access to the repository and would represent the project when answering questions and issues.

Together we could continue the work with things like:

  * Adding more types of diagrams like mindmaps, ert diagrams, etc.
  * Improving existing diagrams



Don't hesitate to contact me if you want to get involved!

## Contributors ​

[](https://github.com/mermaid-js/mermaid/issues?q=is%3Aissue+is%3Aopen+label%3A%22Good+first+issue%21%22)[](https://github.com/mermaid-js/mermaid/graphs/contributors)[](https://github.com/mermaid-js/mermaid/graphs/contributors)

Mermaid is a growing community and is always accepting new contributors. There's a lot of different ways to help out and we're always looking for extra hands! Look at [this issue](https://github.com/mermaid-js/mermaid/issues/866) if you want to know where to start helping out.

Detailed information about how to contribute can be found in the [contribution guideline](./../community/contributing.html).

### Requirements ​

  * [volta](https://volta.sh/) to manage node versions.
  * [Node.js](https://nodejs.org/en/). `volta install node`
  * [pnpm](https://pnpm.io/) package manager. `volta install pnpm`



### Development Installation ​

bash
    
    
    git clone git@github.com:mermaid-js/mermaid.git
    cd mermaid
    # npx is required for first install as volta support for pnpm is not added yet.
    npx pnpm install
    pnpm test

### Lint ​

sh
    
    
    pnpm lint

We use [eslint](https://eslint.org/). We recommend you to install [editor plugins](https://eslint.org/docs/user-guide/integrations) to get real time lint result.

### Test ​

sh
    
    
    pnpm test

Manual test in browser: open `dist/index.html`

### Release ​

For those who have the permission to do so:

Update version number in `package.json`.

sh
    
    
    npm publish

The above command generates files into the `dist` folder and publishes them to [npmjs.com](https://www.npmjs.com/).

## Security and safe diagrams ​

For public sites, it can be precarious to retrieve text from users on the internet, storing that content for presentation in a browser at a later stage. The reason is that the user content can contain embedded malicious scripts that will run when the data is presented. For Mermaid this is a risk, specially as mermaid diagrams contain many characters that are used in html which makes the standard sanitation unusable as it also breaks the diagrams. We still make an effort to sanitize the incoming code and keep refining the process but it is hard to guarantee that there are no loop holes.

As an extra level of security for sites with external users we are happy to introduce a new security level in which the diagram is rendered in a sandboxed iframe preventing JavaScript in the code from being executed. This is a great step forward for better security.

_Unfortunately you cannot have a cake and eat it at the same time which in this case means that some of the interactive functionality gets blocked along with the possible malicious code._

## Reporting vulnerabilities ​

To report a vulnerability, please e-mail [security@mermaid.live](mailto:security@mermaid.live) with a description of the issue, the steps you took to create the issue, affected versions, and if known, mitigations for the issue.

## Appreciation ​

A quick note from Knut Sveidqvist:

> _Many thanks to the[d3](https://d3js.org/) and [dagre-d3](https://github.com/cpettitt/dagre-d3) projects for providing the graphical layout and drawing libraries!_
> 
> _Thanks also to the[js-sequence-diagram](https://bramp.github.io/js-sequence-diagrams) project for usage of the grammar for the sequence diagrams. Thanks to Jessica Peter for inspiration and starting point for gantt rendering._
> 
> _Thank you to[Tyler Long](https://github.com/tylerlong) who has been a collaborator since April 2017._
> 
> _Thank you to the ever-growing list of[contributors](https://github.com/mermaid-js/mermaid/graphs/contributors) that brought the project this far!_

* * *

_Mermaid was created by Knut Sveidqvist for easier documentation._
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
