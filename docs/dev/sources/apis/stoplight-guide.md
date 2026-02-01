# Stoplight API Design Guide

> Source: https://stoplight.io/api-design-guide
> Fetched: 2026-02-01T11:46:10.114343+00:00
> Content-Hash: e5b9688fcc377bdf
> Type: html

---

[](/)

# APIÂ Design Guide

Learn about API Design through our comprehensive guide covering Web API Design principles and best practices.

[Learn about Stoplight APIÂ Design](/api-design)

## What is APIÂ design?

In this guide, we will answer the key question of âwhat is API design,â as well as cover the basics of API design.

API design is the collection of planning and architectural decisions you make when building an API. Your basic API design influences how well developers are able to consume it and even how they use it.

Just like website design or product design, API design informs the user experience. Good API design principles meet initial expectations and continue to behave consistently and predictably.

Organizations with [high design maturity](https://www.invisionapp.com/design-better/design-maturity-model/) experience better quality outcomes for their end users, faster time to market, and better innovation. Thatâs why this API design guide assists in supporting good design throughout your API creation process; good API design leads to better overall APIs.

Better quality outcomes

Faster time to market

Enhanced innovation

There is not a single approach on how to design an API or even how to design good APIs âthe right way.â Instead, we need to lean on good industry basic API design guidelines, best practices and patterns where relevant, then take cues from those who will use our APIs.

Looking to get started with API design? Try [Stoplight Platform](/welcome) to bring a design-first approach to your API workflow.

## How to choose your APIÂ specification

Before you can communicate your API design, you need an artifact that someone else can use to understand your design & API guidelines. Historically, this might have been called documentation.

While itâs still important to have human-facing documentation that is easy to use, more is required of modern APIs for technology to connect with other technology. In recent years the industry has rallied around the OpenAPI Specification.

OpenAPI allows you to define how your REST API works, in a way that can be easily consumed by both humans and machines. It serves as a contract that specifies how a consumer can use the API and what responses you can expect.

## OpenAPI documentation and specification

The industry has selected OpenAPI as the way forward, so letâs understand it and explore what OpenAPI includes in our OpenAPI design guide. From a technical standpoint, it is a YAML or JSON file that follows a specific document structure. You should be able to describe any REST API using a document that adheres to the OpenAPI v3 schema.

### OpenAPI versions v2.0 vs. v3.0

While OpenAPI v3 is the most recent version of OpenAPI, it replaced OpenAPI v2 - previously known as Swagger. The newer version provides a simpler way to describe APIs, while also offering more flexibility. Because there were a lot of legacy Swagger documents, itâs important to have a compatible community-owned version. But API practitioners wanted to move the OpenAPI specification forward with OpenAPI v3, with its latest release being 3.1 in 2021.

### Differences between OpenAPI 2.0 and 3.0

The two major versions of OpenAPI have the most significant differences, which come from their history. OpenAPI 2.0 was previously known as Swagger and is intended to replace it with backward compatibility. Once adopted as an open format, the community began working on OpenAPI 3.0, released in 2017. Let's explore some of the significant changes made to OpenAPI below.

### Differences between OpenAPI 3.0 and 3.1

Moving from the legacy Swagger description format of OpenAPI 2.0 to 3.0 brought many changes. While OpenAPI 3.1 may be a minor release, there are some significant differences between version 3.1 and 3.0.

The changes were notable enough that the community wondered whether the latest release qualified OpenAPI 4.0. OpenAPI Version 3.1 is fully compatible with the latest draft of [JSONÂ Schema, version 2019-09](https://json-schema.org/specification.html).

Now, OpenAPI 3.1 supports all JSON Schema Keywords, so if the keyword exists in the JSON Schema vocabulary, then you can use it with OpenAPI 3.1.

Stoplight Platform supports OpenAPI 2.0, 3.0, and 3.1. [Learn more here.](https://blog.stoplight.io/stoplight-now-supports-openapi-3-1-documents)

Source: [What's the Difference Between OpenAPI 2.0, 3.0, and 3.1?](https://blog.stoplight.io/difference-between-open-v2-v3-v31)

## Why APIÂ design-first matters

Now that youâve chosen OpenAPI v3, you may be tempted to set that aside until after you build your API, and âwrite the docs later.â While itâs useful to describe existing APIs, you can also create API descriptions before and during writing the code.

When you design your API alongside a description, you always have the artifact to communicate whatâs possible with your API. The design-first approach offers a single source of truth, readable by collaborators and machines alike.

### Code-first vs. design-first APIÂ development

Design-first becomes clearer when you consider the alternative. If you go straight into building your API, thereâs no returning to design. Thatâs like constructing a house and then going to an architect to draw up plans. It just makes no sense.

Yet, software teams frequently make similar choices. They may output an API spec from code, which sounds efficient. Unfortunately, by the time youâve built an API in code, youâve lost out on a lot of the advantages of a design-first approach. When your API design exists before the implementation, you can get early feedback, connect your API to tools from the start, and collaborate across departments and functions.

Do you know who will use your API? Even for an internal project, youâre likely to have multiple consumers. An API spec allows you to share details about how the API will work. You can send the spec document itself, or use tools to prototype your API or documentation. You could generate mock servers based on your spec, as described in another section, and have your consumers make live calls.

Your collaboration can go beyond technical teams, as well. You could get great insights from product, marketing, partnerships, and many other areas of your organization.

### The importance of knowing use cases

When you understand how your software will be used you can design it better. The biggest mistake in API design is to make decisions based on how your system works, rather than what your consumers need to support. In order to design around use cases, youâll need to talk to the consumers, or at least include those who know them better.

Software is rarely built entirely by engineers. There are stakeholders throughout the organization. And while many engineers can be very product-minded, they donât always have the visibility of the full picture. If your organization has a product group, thatâs often where the voice of the customer is most heard. Involve anyone who understands how an API will be used in discussions as you design the API.

For example, letâs say you want to design a contact API. Naturally, you would expect to be able to create, list, update, and delete contacts. However, if you donât dig deeper, you are designing an API based on your system. Instead, find out how contacts are created. Do the details come from a user in the field, or are they passed through an online form? Ask the same questions about the other potential endpoints.

When you involve others in API design, you build something better. The API spec becomes an artifact upon which they can comment. You still need ways to coordinate the cross-department conversation, but design-first makes it possible in the first place.

## APIÂ design tools

When you use OpenAPI to design your API, it becomes part of your workflow. That means as soon as you have even a single potential endpoint of your API described, you can begin to gather feedback and piece together how your API will be used. Rather than toiling away in an API silo, your API description allows for collaboration with colleagues and across departments. You can work the API description into your approval processes, so everyone is on the same page with its progress.

Tooling built around the OpenAPI specification can help in the very early stages of design, throughout the life of an API, and even as you consider versioning and deprecation.

[Learn more about Stoplight Platform, a popular and robust API design tool.](/api-design)

### APIÂ documentation tools

Well-documented APIs are more likely to have higher adoption and better user experience. API documentation is one critical component to good design.

There are many tools on the market to help generate quality, up-to-date documentation from your API descriptions.

Before developers and architects used a description document to help them design APIs, documentation was the biggest use case. While OpenAPI allows for much more than generated documentation, that remains a huge advantage to having your API described in OpenAPI.

There are different types of documentation, but OpenAPI-generated docs thrive for API references and interactive documentation. As you add and update your API endpoints, you can automatically keep your documentation updated. You may even be able to connect these tools to your CI/CD workflow, so that as your new API hits production, so does your new API documentation.

Reading documentation is one way to determine how an API works. Live calls add another dimension to that understanding. Interactive documentation means that consumers can test requests against your API, supply their own inputs and see the response inline.

Youâll want to add other types of documentation, too, such as tutorials. Look for a tool that allows you to have customized documentation alongside your generated docs. You also may want to match your siteâs style and navigation.

For an example of fully customizable, generated docs, see Stoplightâs hosted API documentation.

### API mock servers

Just as interactive documentation adds another dimension beyond simple reference, you can benefit from making calls against your API while you design. Your OpenAPI description can be used to create mock servers that use responses youâve included in your design. You can collaborate with others around real data and seek early feedback from API consumers.

Much as documentation is built and rebuilt as you update your API description, mock servers can also automatically have your latest changes. Integrate with your own API as you build it by including mock server endpoints in your code, or coordinate with API consumers and collaborators to write tests or sample code. Code you write against a mock server isnât wasted, because only the server root will change when you move to production.

[Accelerate your API development and collaboration with open source mock servers powered by Prism, a Stoplight open source project.](/open-source/prism)

### Automatically test your APIs

Mocking API calls before theyâre in production is a good idea. Once your API is live, youâll also want to make sure itâs built the way youâve described. Thatâs where API testing comes in.

Your OpenAPI definition describes exactly how your API can be used and what response to expect. During testing, you create scenarios for how your API is used, then run them to make sure you get the correct HTTP status code for the method used. If your OpenAPI document is a contact, testing makes sure youâve built it true to the terms.

Testing can be built into your CI/CD pipeline, so you always know that your tests are passing. Like other software testing, you can track coverage, ensuring that errors are unlikely to slip through. You can build fully customizable tests with built-in coverage reporting with Stoplight OpenAPI testing.

### Use linting to spot errors

As you design your APIs using OpenAPI, youâll need to help your entire team and program conform to your chosen specification. First, create a set of API style guidelines. Then, you can use automated linting tools to validate your JSON or YAML as you write so that it adheres to your style guidelines. An accurate API description is important so that you can feel confident that other tools will interpret your API the way you expect.

Linting tools come in command line, editor plugin, and built-in varieties. It helps you spot errors before you commit them to your repository. Since the OpenAPI spec becomes your source of truth, you want it to be right!

More advanced linting tools can also help you design consistent APIs. For example, have you decided to use plural terms for your resources? If you have an API style guide, you may be able to use a linter to catch that singular endpoint before it goes live. Consistency leads to a better developer experience and a greater likelihood that your API wonât need major changes.

[Improve the quality of your APIÂ descriptions with Stoplight Spectral.](/open-source/spectral)

Regardless of how much your tools help you, itâs a good idea to become familiar with the structure and elements of your OpenAPI documents.

## API design best practices

Armed with an understanding of your use cases, youâre ready to begin your API design. Each project is different, so API design best practices may not always fit your situation. However, these are guidelines to keep in mind as you design your API.

**These are the high level tenets of good API design:**

- Approach design collaboratively
- Maintain internal consistency
- When possible, use an established convention

Youâll want to keep your entire team updated as you make design decisions together. Your OpenAPI spec is your single source of truth, so make sure it is available in a place where everyone can see revisions and discuss changes. A GitHub repository or Stoplightâs Visual OpenAPI Designer can help keep everyone on the same page.

Contents

Example H2

Example H3

Design quality APIs with Stoplight.

[Get Started for Free](/welcome)[See a Demo](/demo)

[](/)

Products

- [Stoplight Solutions](/solutions)
- [Enterprise Sales](/enterprise)
- [OpenÂ Source](/open-source)
- [Pricing](/pricing)

Resources

- [Stoplight Docs](https://docs.stoplight.io/)
- [Blog](https://blog.stoplight.io/)
- [Podcast](/podcast)
- [Guides](/guides)
- [Webinars](/webinars)

Help

- [See a Demo](/demo)
- [Get Support](https://support.stoplight.io/)
- [Contact Us](/contact)
- [Stoplight Community](https://community.stoplight.io/)
- [Status Page](https://status.stoplight.io/)

About

- [About Us](/about)
- [Press](/press)
- [Case Studies](/case-studies)
- [Roadmap](https://roadmap.stoplight.io/)
- [Careers](/careers-old)

Â© 2024 SmartBear Software. All Rights Reserved.

[Website Terms of Use](/website-terms)[Subscription Agreement](/terms)[Privacy Policy](/privacy)[Support Policy](/policy)[Security](/security-practices)

[](https://twitter.com/stoplightio)[](https://www.linkedin.com/company/stoplight)[](https://www.github.com/stoplightio)[](https://www.facebook.com/stoplightio/)
  *[↑]: Back to Top
