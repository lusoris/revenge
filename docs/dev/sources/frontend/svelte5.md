# Svelte 5 Documentation

> Auto-fetched from [https://svelte.dev/docs/svelte/overview](https://svelte.dev/docs/svelte/overview)
> Last Updated: 2026-01-28T21:46:01.454242+00:00

---

Introduction
Overview
Getting started
.svelte files
.svelte.js and .svelte.ts files
Runes
What are runes?
$state
$derived
$effect
$props
$bindable
$inspect
$host
Template syntax
Basic markup
{#if ...}
{#each ...}
{#key ...}
{#await ...}
{#snippet ...}
{@render ...}
{@html ...}
{@attach ...}
{@const ...}
{@debug ...}
bind:
use:
transition:
in: and out:
animate:
style:
class
await
Styling
Scoped styles
Global styles
Custom properties
Nested <style> elements
Special elements
<svelte:boundary>
<svelte:window>
<svelte:document>
<svelte:body>
<svelte:head>
<svelte:element>
<svelte:options>
Runtime
Stores
Context
Lifecycle hooks
Imperative component API
Hydratable data
Misc
Testing
TypeScript
Custom elements
Svelte 4 migration guide
Svelte 5 migration guide
Frequently asked questions
Reference
svelte
svelte/action
svelte/animate
svelte/attachments
svelte/compiler
svelte/easing
svelte/events
svelte/legacy
svelte/motion
svelte/reactivity/window
svelte/reactivity
svelte/server
svelte/store
svelte/transition
Compiler errors
Compiler warnings
Runtime errors
Runtime warnings
Legacy APIs
Overview
Reactive let/var declarations
Reactive $: statements
export let
$$props and $$restProps
on:
<slot>
$$slots
<svelte:fragment>
<svelte:component>
<svelte:self>
Imperative component API
Svelte
Introduction
Overview
On this page
Overview
Svelte is a framework for building user interfaces on the web. It uses a compiler to turn declarative components written in HTML, CSS and JavaScript...
App
<
script
>
function
greet
() {
alert
(
'Welcome to Svelte!'
);
}
</
script
>
<
button
onclick
=
{greet}>click me</
button
>
<
style
>
button
{
font-size
:
2
em
;
}
</
style
>
<
script
lang
=
"ts"
>
function
greet
() {
alert
(
'Welcome to Svelte!'
);
}
</
script
>
<
button
onclick
=
{greet}>click me</
button
>
<
style
>
button
{
font-size
:
2
em
;
}
</
style
>
...into lean, tightly optimized JavaScript.
You can use it to build anything on the web, from standalone components to ambitious full stack apps (using Svelte’s companion application framework,
SvelteKit
) and everything in between.
These pages serve as reference documentation. If you’re new to Svelte, we recommend starting with the
interactive tutorial
and coming back here when you have questions.
You can also try Svelte online in the
playground
or, if you need a more fully-featured environment, on
StackBlitz
.
Edit this page on GitHub
llms.txt
previous
next
Getting started