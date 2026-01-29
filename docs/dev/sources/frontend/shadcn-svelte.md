# shadcn-svelte

> Auto-fetched from [https://www.shadcn-svelte.com/docs](https://www.shadcn-svelte.com/docs)
> Last Updated: 2026-01-29T20:14:08.641135+00:00

---

Sections
Get Started
Components
Changelog
Get Started
Installation
components.json
Theming
Dark Mode
CLI
JavaScript
Figma
llms.txt
Legacy Docs
Migration
Svelte 5
Tailwind v4
Components
Accordion
Alert Dialog
Alert
Aspect Ratio
Avatar
Badge
Breadcrumb
Button Group
Button
Calendar
Card
Carousel
Chart
Checkbox
Collapsible
Combobox
Command
Context Menu
Data Table
Date Picker
Dialog
Drawer
Dropdown Menu
Empty
Field
Formsnap
Hover Card
Input Group
Input OTP
Input
Item
Kbd
Label
Menubar
Native Select
Navigation Menu
Pagination
Popover
Progress
Radio Group
Range Calendar
Resizable
Scroll Area
Select
Separator
Sheet
Sidebar
Skeleton
Slider
Sonner
Spinner
Switch
Table
Tabs
Textarea
Toggle Group
Toggle
Tooltip
Typography
Installation
SvelteKit
Vite
Astro
Manual Installation
Dark Mode
Svelte
Astro
Registry
Registry
Getting Started
FAQ
Examples
registry.json
registry-item.json
On This Page
Open Code
Composition
Distribution
Beautiful Defaults
AI-Ready
Special sponsor
We're looking for one partner to be featured here.
Support the project and reach thousands of developers.
Reach out
Introduction
Copy Page
Next
Re-usable components built with Bits UI and Tailwind CSS.
An unofficial, community-led
Svelte
port of
shadcn/ui
. We are not affiliated with
shadcn
, but we did get his blessing before creating a Svelte version of his work. This project was born out of the need for a similar project for the Svelte ecosystem.
This is not a component library. It is how you build your component library.
You know how most traditional component libraries work: you install a package from NPM, import the components, and use them in your app.
This approach works well until you need to customize a component to fit your design system or require one that isnât included in the library.
Often, you end up wrapping library components, writing workarounds to override styles, or mixing components from different libraries with incompatible APIs.
This is what shadcn-svelte aims to solve. It is built around the following principles:
Open Code:
The top layer of your component code is open for modification.
Composition:
Every component uses a common, composable interface, making them predictable.
Distribution:
A flat-file schema and command-line tool make it easy to distribute components.
Beautiful Defaults:
Carefully chosen default styles, so you get great design out-of-the-box.
AI-Ready:
Open code for LLMs to read, understand, and improve.
Open Code
shadcn-svelte hands you the actual component code. You have full control to customize and extend the components to your needs. This means:
Full Transparency:
You see exactly how each component is built.
Easy Customization:
Modify any part of a component to fit your design and functionality requirements.
AI Integration:
Access to the code makes it straightforward for LLMs to read, understand, and even improve your components.
In a typical library, if you need to change a buttonâs behavior, you have to override styles or wrap the component. With shadcn-svelte, you simply edit the button code directly.
How do I pull upstream updates in an Open Code approach?
shadcn-svelte follows a headless component architecture. This means the core of your app can receive fixes by updating your dependencies, for instance, bits-ui or paneforge.
The topmost layer, i.e., the one closest to your design system, is not
coupled with the implementation of the library. It stays open for
modification.
Composition
Every component in shadcn-svelte shares a common, composable interface.
If a component does not exist, we bring it in, make it composable, and adjust its style to match and work with the rest of the design system.
A shared, composable interface means it's predictable for both your team and LLMs. You are not learning different APIs for every new component. Even for third-party ones.
Distribution
shadcn-svelte is also a code distribution system. It defines a schema for components and a CLI to distribute them.
Schema:
A flat-file structure that defines the components, their dependencies, and properties.
CLI:
A command-line tool to distribute and install components across projects with cross-framework support.
You can use the schema to distribute your components to other projects or have AI generate completely new components based on existing schema.
Beautiful Defaults
shadcn-svelte comes with a large collection of components that have carefully chosen default styles. They are designed to look good on their own and to work well together as a consistent system:
Good Out-of-the-Box:
Your UI has a clean and minimal look without extra work.
Unified Design:
Components naturally fit with one another. Each component is built to match the others, keeping your UI consistent.
Easily Customizable:
If you want to change something, it's simple to override and extend the defaults.
AI-Ready
The design of shadcn-svelte makes it easy for AI tools to work with your code. Its open code and consistent API allow AI models to read, understand, and even generate new components.
An AI model can learn how your components work and suggest improvements or even create new components that integrate with your existing design.
Installation