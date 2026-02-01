# Vidstack Player

> Source: https://www.vidstack.io/docs/player/getting-started/installation
> Fetched: 2026-02-01T11:47:07.367893+00:00
> Content-Hash: b71783470c904c9c
> Type: html

---

Getting Started

# Installation

Instructions to get your player installed and on-screen.

## 1. Select Framework

Section titled 1. Select Framework

[](/docs/player/getting-started/installation/javascript)

JavaScript  More information

The JavaScript option works with any framework and allows constructing the player and layouts imperatively like so:

    VidstackPlayer.create({
      target: '#player',
      src: '...',
      layout: new VidstackPlayerLayout(),
    });
    

You can target an existing audio, video, or iframe element to progressively enhance it.

[](/docs/player/getting-started/installation/angular)

Angular  More information

We don’t have a specific integration for [Angular](https://angularjs.org/) at the moment. However, you can enjoy all Vidstack Player features through our [Web Component](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) library as [Angular has perfect support](https://custom-elements-everywhere.com/#angular).

[](/docs/player/getting-started/installation/react)

React  More information

The React integration will provide a seamless experience when working with [React](https://reactjs.org) and frameworks like [Next.js](https://nextjs.org). Our library is React Sever Component and Next.js `app/` directory ready.

[](/docs/player/getting-started/installation/svelte)

Svelte  More information

We don’t have a specific integration for [Svelte](https://svelte.dev) at the moment. However, you can enjoy all Vidstack Player features through our [Web Component](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) library as [Svelte has perfect support](https://custom-elements-everywhere.com/#svelte).

We also ship JSX types to make sure you have complete TypeScript support for props and events! In addition, we don’t use Shadow DOM so it will work with SSR-frameworks like SvelteKit.

[](/docs/player/getting-started/installation/vue)

Vue  More information

We don’t have a specific integration for [Vue](https://vuejs.org) at the moment. However, you can enjoy all Vidstack Player features through our [Web Component](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) library as [Vue has perfect support](https://custom-elements-everywhere.com/#vue).

We also ship JSX types to make sure you have complete TypeScript support for props and events! In addition, we don’t use Shadow DOM so it will work with SSR-frameworks like Nuxt.

[](/docs/player/getting-started/installation/solid)

Solid  More information

We don’t have a specific integration for [Solid](https://www.solidjs.com) at the moment. However, you can enjoy all Vidstack Player features through our [Web Component](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) library as [Solid has perfect support](https://custom-elements-everywhere.com/#solid).

We also ship JSX types to make sure you have complete TypeScript support for props and events! In addition, we don’t use Shadow DOM so it will work with SSR-frameworks like SolidStart.

[](/docs/player/getting-started/installation/web-components)

Web Components  More information

Our [Web Component](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) library can be used anywhere with the simple drop of an import or CDN link as they’re natively supported by all browsers.

This option is best when writing plain HTML or using a JS library that supports them. Custom Elements have [perfect support in most libraries](https://custom-elements-everywhere.com).

[](/docs/player/getting-started/installation/cdn)

CDN  More information

Using the [JSDelivr CDN](https://www.jsdelivr.com) is the simplest and fastest way to start using the player library via [Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_components).

We provide a CDN bundle that includes all package dependencies, and it’s specially minified to get the bundle size as small as possible. Add a few `<script>` and `<style>` tags to your `<head>` element and you’re ready to start building!

## 2. Select Provider

Section titled 2. Select Provider

Audio  More information

You can use this provider to build an audio player. This provider uses the native `<audio>` element to support [audio codecs](https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Audio_codecs) such as AAC, MP3, FLAC, etc.

Do note, the type of audio codec supported is completely dependent on the browser so use the prior link as reference to see what’s supported.

Video  More information

You can use this provider to build a video player. This provider uses the native `<video>` element to support [video codecs](https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Video_codecs) such as AV1, AVC (H.264), VP9, etc.

Do note, the type of video codec supported is completely dependent on the browser so use the prior link as reference to see what’s supported.

HLS  More information

Embed video content into documents via the native video element. This provider enables streaming video using the HTTP Live Streaming (HLS) protocol.

[HLS isn’t widely supported yet](https://caniuse.com/?search=hls), but we use the popular [hls.js](https://github.com/video-dev/hls.js) library to ensure it works anywhere [Media Source Extensions](https://caniuse.com/mediasource) (MSE) or [Managed Media Source](https://caniuse.com/mdn-api_managedmediasource) (MMS) is supported.

DASH  More information

Embed video content into documents via the native video element. This provider enables streaming video using the Dynamic Adaptive Streaming over HTTP (DASH) protocol.

[DASH isn’t supported in any browser](https://caniuse.com/?search=dash), but we use the popular [dash.js](https://github.com/Dash-Industry-Forum/dash.js) library to ensure it works anywhere [Media Source Extensions](https://caniuse.com/mediasource) (MSE) or [Managed Media Source](https://caniuse.com/mdn-api_managedmediasource) (MMS) is supported.

YouTube  More information

Embed YouTube content using [iframes](https://developers.google.com/youtube/iframe_api_reference).

Vimeo  More information

Embed YouTube content using [iframes](https://developers.google.com/youtube/iframe_api_reference).

Remotion  More information

[Remotion](https://www.remotion.dev) enables creating complex animations and videos programatically using React. You can use this provider to embed, preview, and play dynamic React components that are using Remotion. You can also provide the rendered MP4 directly to the player in production (see the video provider).

## 3. Select Styling

Section titled 3. Select Styling

A layout refers to the arrangement and presentation of various player components. The CSS and Tailwind CSS options below are if you want to style components from scratch and build your own layout. The Default Theme option is if you want to build your own layout on top of our component styles. Finally, the Default and Plyr layouts are our production-ready templates to get you up and running quickly.
  
CSS  More information

The CSS option provides you with a minimal starting point and completely unstyled components. All components including the player itself provide styling hooks via data attributes and support animations.

This option is best when you want to build your player yourself from scratch using vanilla CSS.

Default Theme  More information

The Default Theme is best when you want to build your own player but you don’t want to design each component from zero. Our default styles have been built to be extremely easy to customize! We provide the shell and base styles and you provide the content and customization.

Default Layout  More information

The [Default Layout](/docs/player/components/layouts/default-layout) is our production-ready UI. If you’re looking for something pre-designed and ready out of the box then this is the best option. You can easily customize the icons, branding, colors, and components to your liking.

Plyr Layout  More information

Based on the beautiful 2015 design by [Sam Potts](https://twitter.com/sam_potts), the [Plyr Layout](https://plyr.io) is a simple and elegant option. This layout includes the same features, styles, CSS variables, icons, and options as the original Plyr player to make migrating over simple.

Tailwind CSS  More information

The [Tailwind](https://tailwindcss.com) option provides you with a minimal starting point and completely unstyled components. All components including the player itself provide styling hooks via data attributes and support animations.

Our [optional plugin](/docs/player/styling/tailwind) can help speed you up even more by providing you with easy to use media variants such as `media-paused:opacity-0`.

* * *

Previous

[Introduction](/docs/player)

Next

[Architecture](/docs/player/getting-started/architecture)
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
