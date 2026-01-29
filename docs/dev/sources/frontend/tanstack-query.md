# TanStack Query

> Auto-fetched from [https://tanstack.com/query/latest/docs/framework/svelte/overview](https://tanstack.com/query/latest/docs/framework/svelte/overview)
> Last Updated: 2026-01-29T20:14:14.436498+00:00

---

TanStack
Query
v5
v5
Alpha
Try TanStack CLI
Search...
K
Auto
Log In
Start
RC
Start
RC
Router
Router
Query
Query
Table
Table
DB
beta
DB
beta
AI
alpha
AI
alpha
Form
new
Form
new
Virtual
Virtual
Pacer
beta
Pacer
beta
Store
alpha
Store
alpha
Devtools
alpha
Devtools
alpha
CLI
alpha
CLI
alpha
More Libraries
More Libraries
Builder
Alpha
Builder
Alpha
Feed
Beta
Feed
Beta
Maintainers
Maintainers
Partners
Partners
Showcase
Showcase
Blog
Blog
Learn
NEW
Learn
NEW
Support
Support
Stats
Stats
Discord
Discord
Merch
Merch
GitHub
GitHub
Ethos
Ethos
Tenets
Tenets
Brand Guide
Brand Guide
Docs
Partners
Svelte
Latest
Search...
K
Menu
Home
Frameworks
Contributors
NPM Stats
Community Resources
GitHub
Discord
Getting Started
Overview
Installation
Devtools
SSR & SvelteKit
Migrate from v5 to v6
API Reference
QueryClient
QueryCache
MutationCache
QueryObserver
InfiniteQueryObserver
QueriesObserver
streamedQuery
focusManager
onlineManager
notifyManager
timeoutManager
Svelte Reference
Functions / createQuery
Functions / createQueries
Functions / createInfiniteQuery
Functions / createMutation
Functions / useIsFetching
Functions / useIsMutating
Functions / useMutationState
Functions / queryOptions
Functions / infiniteQueryOptions
ESLint
ESLint Plugin Query
Exhaustive Deps
Stable Query Client
No Rest Destructuring
No Unstable Deps
Infinite Query Property Order
No void Query Functions
Mutation Property Order
Examples
Simple
Basic
Auto Refetching / Polling / Realtime
SSR
Optimistic Updates
Playground
Star Wars
Infinite Queries
latest
Svelte
Latest
Menu
Home
Frameworks
Contributors
NPM Stats
Community Resources
GitHub
Discord
Getting Started
Overview
Installation
Devtools
SSR & SvelteKit
Migrate from v5 to v6
API Reference
QueryClient
QueryCache
MutationCache
QueryObserver
InfiniteQueryObserver
QueriesObserver
streamedQuery
focusManager
onlineManager
notifyManager
timeoutManager
Svelte Reference
Functions / createQuery
Functions / createQueries
Functions / createInfiniteQuery
Functions / createMutation
Functions / useIsFetching
Functions / useIsMutating
Functions / useMutationState
Functions / queryOptions
Functions / infiniteQueryOptions
ESLint
ESLint Plugin Query
Exhaustive Deps
Stable Query Client
No Rest Destructuring
No Unstable Deps
Infinite Query Property Order
No void Query Functions
Mutation Property Order
Examples
Simple
Basic
Auto Refetching / Polling / Realtime
SSR
Optimistic Updates
Playground
Star Wars
Infinite Queries
AI/LLM: This documentation page is available in plain markdown format at
/query/latest/docs/framework/svelte/overview
.md
Learn about TanStack Ads
Hide Ads
Getting Started
On this page
Overview
Copy page
The
@tanstack/svelte-query
package offers a 1st-class API for using TanStack Query via Svelte.
Migrating from stores to the runes syntax? See the
migration guide
.
Example
Include the QueryClientProvider near the root of your project:
svelte
<script lang="ts">
import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query'
import Example from './lib/Example.svelte'

const queryClient = new QueryClient()
</script>

<QueryClientProvider client={queryClient}>
<Example />
</QueryClientProvider>
Then call any function (e.g. createQuery) from any component:
svelte
<script lang="ts">
import { createQuery } from '@tanstack/svelte-query'

const query = createQuery(() => ({
queryKey: ['todos'],
queryFn: () => fetchTodos(),
}))
</script>

<div>
{#if query.isLoading}
<p>Loading...</p>
{:else if query.isError}
<p>Error: {query.error.message}</p>
{:else if query.isSuccess}
{#each query.data as todo}
<p>{todo.title}</p>
{/each}
{/if}
</div>
SvelteKit
If you are using SvelteKit, please have a look at
SSR & SvelteKit
.
Available Functions
Svelte Query offers useful functions and components that will make managing server state in Svelte apps easier.
createQuery
createQueries
createInfiniteQuery
createMutation
useQueryClient
useIsFetching
useIsMutating
useMutationState
useIsRestoring
useHydrate
<QueryClientProvider>
<HydrationBoundary>
Important Differences between Svelte Query & React Query
Svelte Query offers an API similar to React Query, but there are some key differences to be mindful of.
The arguments to the
create*
functions must be wrapped in a function to preserve reactivity.
Edit on GitHub
Previous
Community Resources
Next
Installation
On this page
Example
SvelteKit
Available Functions
Important Differences between Svelte Query & React Query
Learn about TanStack Ads
Hide Ads
Partners
Become a Partner
Learn about TanStack Ads
Hide Ads
Want to Skip the Docs?
Query.gg - The Official React Query Course
“If you’re serious about *really* understanding React Query, there’s no better way than with query.gg”
—Tanner Linsley
Learn More