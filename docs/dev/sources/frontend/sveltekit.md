# SvelteKit Documentation

> Auto-fetched from [https://svelte.dev/docs/kit/introduction](https://svelte.dev/docs/kit/introduction)
> Last Updated: 2026-01-28T21:46:06.828975+00:00

---

Getting started
Introduction
Creating a project
Project types
Project structure
Web standards
Core concepts
Routing
Loading data
Form actions
Page options
State management
Remote functions
Build and deploy
Building your app
Adapters
Zero-config deployments
Node servers
Static site generation
Single-page apps
Cloudflare
Cloudflare Workers
Netlify
Vercel
Writing adapters
Advanced
Advanced routing
Hooks
Errors
Link options
Service workers
Server-only modules
Snapshots
Shallow routing
Observability
Packaging
Best practices
Auth
Performance
Icons
Images
Accessibility
SEO
Appendix
Frequently asked questions
Integrations
Breakpoint Debugging
Migrating to SvelteKit v2
Migrating from Sapper
Additional resources
Glossary
Reference
@sveltejs/kit
@sveltejs/kit/hooks
@sveltejs/kit/node/polyfills
@sveltejs/kit/node
@sveltejs/kit/vite
$app/environment
$app/forms
$app/navigation
$app/paths
$app/server
$app/state
$app/stores
$app/types
$env/dynamic/private
$env/dynamic/public
$env/static/private
$env/static/public
$lib
$service-worker
Configuration
Command Line Interface
Types
SvelteKit
Getting started
Introduction
On this page
Introduction
Before we begin
What is SvelteKit?
What is Svelte?
SvelteKit vs Svelte
Before we begin
If you’re new to Svelte or SvelteKit we recommend checking out the
interactive tutorial
.
If you get stuck, reach out for help in the
Discord chatroom
.
What is SvelteKit?
SvelteKit is a framework for rapidly developing robust, performant web applications using
Svelte
. If you’re coming from React, SvelteKit is similar to Next. If you’re coming from Vue, SvelteKit is similar to Nuxt.
To learn more about the kinds of applications you can build with SvelteKit, see the
documentation regarding project types
.
What is Svelte?
In short, Svelte is a way of writing user interface components — like a navigation bar, comment section, or contact form — that users see and interact with in their browsers. The Svelte compiler converts your components to JavaScript that can be run to render the HTML for the page and to CSS that styles the page. You don’t need to know Svelte to understand the rest of this guide, but it will help. If you’d like to learn more, check out
the Svelte tutorial
.
SvelteKit vs Svelte
Svelte renders UI components. You can compose these components and render an entire page with just Svelte, but you need more than just Svelte to write an entire app.
SvelteKit helps you build web apps while following modern best practices and providing solutions to common development challenges. It offers everything from basic functionalities — like a
router
that updates your UI when a link is clicked — to more advanced capabilities. Its extensive list of features includes
build optimizations
to load only the minimal required code;
offline support
;
preloading
pages before user navigation;
configurable rendering
to handle different parts of your app on the server via
SSR
, in the browser through
client-side rendering
, or at build-time with
prerendering
;
image optimization
; and much more. Building an app with all the modern best practices is fiendishly complicated, but SvelteKit does all the boring stuff for you so that you can get on with the creative part.
It reflects changes to your code in the browser instantly to provide a lightning-fast and feature-rich development experience by leveraging
Vite
with a
Svelte plugin
to do
Hot Module Replacement (HMR)
.
Edit this page on GitHub
llms.txt
previous
next
Creating a project