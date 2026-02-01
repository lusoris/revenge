# Jinja2 Documentation

> Source: https://jinja.palletsprojects.com/en/stable/
> Fetched: 2026-02-01T11:55:20.643549+00:00
> Content-Hash: 13a18fc38d745c42
> Type: html

---

### Navigation

- [index](genindex/ "General Index")
- [modules](py-modindex/ "Python Module Index") |
- [next](intro/ "Introduction") |
- Jinja Documentation (3.1.x) »
- [Jinja]()

# Jinja¶

[](_images/jinja-name.svg)

Jinja is a fast, expressive, extensible templating engine. Special placeholders in the template allow writing code similar to Python syntax. Then the template is passed data to render the final document.

Contents:

- [Introduction](intro/)
  - [Installation](intro/#installation)
- [API](api/)
  - [Basics](api/#basics)
  - [High Level API](api/#high-level-api)
  - [Autoescaping](api/#autoescaping)
  - [Notes on Identifiers](api/#notes-on-identifiers)
  - [Undefined Types](api/#undefined-types)
  - [The Context](api/#the-context)
  - [Loaders](api/#loaders)
  - [Bytecode Cache](api/#bytecode-cache)
  - [Async Support](api/#async-support)
  - [Policies](api/#policies)
  - [Utilities](api/#utilities)
  - [Exceptions](api/#exceptions)
  - [Custom Filters](api/#custom-filters)
  - [Custom Tests](api/#custom-tests)
  - [Evaluation Context](api/#evaluation-context)
  - [The Global Namespace](api/#the-global-namespace)
  - [Low Level API](api/#low-level-api)
  - [The Meta API](api/#the-meta-api)
- [Sandbox](sandbox/)
  - [Security Considerations](sandbox/#security-considerations)
  - [API](sandbox/#module-jinja2.sandbox)
  - [Operator Intercepting](sandbox/#operator-intercepting)
- [Native Python Types](nativetypes/)
  - [Examples](nativetypes/#examples)
  - [Sandboxed Native Environment](nativetypes/#sandboxed-native-environment)
  - [API](nativetypes/#api)
- [Template Designer Documentation](templates/)
  - [Synopsis](templates/#synopsis)
  - [Variables](templates/#variables)
  - [Filters](templates/#filters)
  - [Tests](templates/#tests)
  - [Comments](templates/#comments)
  - [Whitespace Control](templates/#whitespace-control)
  - [Escaping](templates/#escaping)
  - [Line Statements](templates/#line-statements)
  - [Template Inheritance](templates/#template-inheritance)
  - [HTML Escaping](templates/#html-escaping)
  - [List of Control Structures](templates/#list-of-control-structures)
  - [Import Context Behavior](templates/#import-context-behavior)
  - [Expressions](templates/#expressions)
  - [List of Builtin Filters](templates/#list-of-builtin-filters)
  - [List of Builtin Tests](templates/#list-of-builtin-tests)
  - [List of Global Functions](templates/#list-of-global-functions)
  - [Extensions](templates/#extensions)
  - [Autoescape Overrides](templates/#autoescape-overrides)
- [Extensions](extensions/)
  - [Adding Extensions](extensions/#adding-extensions)
  - [i18n Extension](extensions/#i18n-extension)
  - [Expression Statement](extensions/#expression-statement)
  - [Loop Controls](extensions/#loop-controls)
  - [With Statement](extensions/#with-statement)
  - [Autoescape Extension](extensions/#autoescape-extension)
  - [Debug Extension](extensions/#debug-extension)
  - [Writing Extensions](extensions/#module-jinja2.ext)
  - [Example Extensions](extensions/#example-extensions)
  - [Extension API](extensions/#extension-api)
- [Integration](integration/)
  - [Flask](integration/#flask)
  - [Django](integration/#django)
  - [Babel](integration/#babel)
  - [Pylons](integration/#pylons)
- [Switching From Other Template Engines](switching/)
  - [Django](switching/#django)
  - [Mako](switching/#mako)
- [Tips and Tricks](tricks/)
  - [Null-Default Fallback](tricks/#null-default-fallback)
  - [Alternating Rows](tricks/#alternating-rows)
  - [Highlighting Active Menu Items](tricks/#highlighting-active-menu-items)
  - [Accessing the parent Loop](tricks/#accessing-the-parent-loop)
- [Frequently Asked Questions](faq/)
  - [Why is it called Jinja?](faq/#why-is-it-called-jinja)
  - [How fast is Jinja?](faq/#how-fast-is-jinja)
  - [Isn’t it a bad idea to put logic in templates?](faq/#isn-t-it-a-bad-idea-to-put-logic-in-templates)
  - [Why is HTML escaping not the default?](faq/#why-is-html-escaping-not-the-default)
- [BSD-3-Clause License](license/)
- [Changes](changes/)
  - [Version 3.1.6](changes/#version-3-1-6)
  - [Version 3.1.5](changes/#version-3-1-5)
  - [Version 3.1.4](changes/#version-3-1-4)
  - [Version 3.1.3](changes/#version-3-1-3)
  - [Version 3.1.2](changes/#version-3-1-2)
  - [Version 3.1.1](changes/#version-3-1-1)
  - [Version 3.1.0](changes/#version-3-1-0)
  - [Version 3.0.3](changes/#version-3-0-3)
  - [Version 3.0.2](changes/#version-3-0-2)
  - [Version 3.0.1](changes/#version-3-0-1)
  - [Version 3.0.0](changes/#version-3-0-0)
  - [Version 2.11.3](changes/#version-2-11-3)
  - [Version 2.11.2](changes/#version-2-11-2)
  - [Version 2.11.1](changes/#version-2-11-1)
  - [Version 2.11.0](changes/#version-2-11-0)
  - [Version 2.10.3](changes/#version-2-10-3)
  - [Version 2.10.2](changes/#version-2-10-2)
  - [Version 2.10.1](changes/#version-2-10-1)
  - [Version 2.10](changes/#version-2-10)
  - [Version 2.9.6](changes/#version-2-9-6)
  - [Version 2.9.5](changes/#version-2-9-5)
  - [Version 2.9.4](changes/#version-2-9-4)
  - [Version 2.9.3](changes/#version-2-9-3)
  - [Version 2.9.2](changes/#version-2-9-2)
  - [Version 2.9.1](changes/#version-2-9-1)
  - [Version 2.9](changes/#version-2-9)
  - [Version 2.8.1](changes/#version-2-8-1)
  - [Version 2.8](changes/#version-2-8)
  - [Version 2.7.3](changes/#version-2-7-3)
  - [Version 2.7.2](changes/#version-2-7-2)
  - [Version 2.7.1](changes/#version-2-7-1)
  - [Version 2.7](changes/#version-2-7)
  - [Version 2.6](changes/#version-2-6)
  - [Version 2.5.5](changes/#version-2-5-5)
  - [Version 2.5.4](changes/#version-2-5-4)
  - [Version 2.5.3](changes/#version-2-5-3)
  - [Version 2.5.2](changes/#version-2-5-2)
  - [Version 2.5.1](changes/#version-2-5-1)
  - [Version 2.5](changes/#version-2-5)
  - [Version 2.4.1](changes/#version-2-4-1)
  - [Version 2.4](changes/#version-2-4)
  - [Version 2.3.1](changes/#version-2-3-1)
  - [Version 2.3](changes/#version-2-3)
  - [Version 2.2.1](changes/#version-2-2-1)
  - [Version 2.2](changes/#version-2-2)
  - [Version 2.1.1](changes/#version-2-1-1)
  - [Version 2.1](changes/#version-2-1)
  - [Version 2.0](changes/#version-2-0)
  - [Version 2.0rc1](changes/#version-2-0rc1)

### Project Links

- [Donate](https://palletsprojects.com/donate)
- [PyPI Releases](https://pypi.org/project/Jinja2/)
- [Source Code](https://github.com/pallets/jinja/)
- [Issue Tracker](https://github.com/pallets/jinja/issues/)
- [Chat](https://discord.gg/pallets)

### Quick search

© Copyright 2007 Pallets. Created using [Sphinx](https://www.sphinx-doc.org/) 8.1.3.
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
