# Svelte 5 Runes

> Source: https://svelte.dev/docs/svelte/$state
> Fetched: 2026-01-31T15:59:58.267776+00:00
> Content-Hash: fcd0763ab3339330
> Type: html

---

### On this page

The `$state` rune allows you to create _reactive state_ , which means that your UI _reacts_ when it changes.
    
    
    <script>
    	let count = $state(0);
    </script>
    
    <button onclick={() => count++}>
    	clicks: {count}
    </button>

Unlike other frameworks you may have encountered, there is no API for interacting with state — `count` is just a number, rather than an object or a function, and you can update it like you would update any other variable.

### Deep state

If `$state` is used with an array or a simple object, the result is a deeply reactive _state proxy_. [Proxies](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy) allow Svelte to run code when you read or write properties, including via methods like `array.push(...)`, triggering granular updates.

State is proxified recursively until Svelte finds something other than an array or simple object (like a class or an object created with `Object.create`). In a case like this...
    
    
    let 
    
    
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos = `
    
    
    function $state<{
        done: boolean;
        text: string;
    }[]>(initial: {
        done: boolean;
        text: string;
    }[]): {
        done: boolean;
        text: string;
    }[] (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state([ { `done: boolean`done: false, `text: string`text: 'add more todos' } ]);`

...modifying an individual todo’s property will trigger updates to anything in your UI that depends on that specific property:
    
    
    module todos
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos[0].`done: boolean`done = !`
    
    
    module todos
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos[0].`done: boolean`done;`

If you push a new object to the array, it will also be proxified:
    
    
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos.`
    
    
    Array<{ done: boolean; text: string; }>.push(...items: {
        done: boolean;
        text: string;
    }[]): number

`

Appends new elements to the end of an array, and returns the new length of the array.

@paramitems New elements to add to the array.

push({ `done: boolean`done: false, `text: string`text: 'eat lunch' });`

> When you update properties of proxies, the original object is _not_ mutated. If you need to use your own proxy handlers in a state proxy, [you should wrap the object _after_ wrapping it in `$state`](https://svelte.dev/playground/hello-world?version=latest#H4sIAAAAAAAACpWR3WoDIRCFX2UqhWyIJL3erAulL9C7XnQLMe5ksbUqOpsfln33YuyGFNJC8UKdc2bOhw7Myk9kJXsJ0nttO9jcR5KEG9AWJDwHdzwxznbaYGTl68Do5JM_FRifuh-9X8Y9Gkq1rYx4q66cJbQUWcmqqIL2VDe2IYMEbvuOikBADi-GJDSkXG-phId0G-frye2DO2psQYDFQ0Ys8gQO350dUkEydEg82T0GOs0nsSG9g2IqgxACZueo2ZUlpdvoDC6N64qsg1QKY8T2bpZp8gpIfbCQ85Zn50Ud82HkeY83uDjspenxv3jXcSDyjPWf9L1vJf0GH666J-jLu1ery4dV257IWXBWGa0-xFDMQdTTn2ScxWKsn86ROsLwQxqrVR5QM84Ij8TKFD2-cUZSm4O2LSt30kQcvwCgCmfZnAIAAA==).

Note that if you destructure a reactive value, the references are not reactive — as in normal JavaScript, they are evaluated at the point of destructuring:
    
    
    let { let done: booleandone, let text: stringtext } = 
    
    
    module todos
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos[0]; // this will not affect the value of `done` `
    
    
    module todos
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos[0].`done: boolean`done = !`
    
    
    module todos
    let todos: {
        done: boolean;
        text: string;
    }[]

`todos[0].`done: boolean`done;`

### Classes

Class instances are not proxied. Instead, you can use `$state` in class fields (whether public or private), or as the first assignment to a property immediately inside the `constructor`:
    
    
    class class TodoTodo {
    	Todo.done: booleandone = 
    
    
    function $state<false>(initial: false): false (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(false); constructor(text) { this.`Todo.text: any`text = `
    
    
    function $state<any>(initial: any): any (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(`text: any`text); } `Todo.reset(): void`reset() { this.`Todo.text: any`text = ''; this.`Todo.done: boolean`done = false; } }`

> The compiler transforms `done` and `text` into `get` / `set` methods on the class prototype referencing private fields. This means the properties are not enumerable.

When calling methods in JavaScript, the value of [`this`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/this) matters. This won’t work, because `this` inside the `reset` method will be the `<button>` rather than the `Todo`:
    
    
    <button onclick={todo.reset}>
    	reset
    </button>

You can either use an inline function...
    
    
    <button onclick={() => todo.reset()}>
    	reset
    </button>

...or use an arrow function in the class definition:
    
    
    class class TodoTodo {
    	Todo.done: booleandone = 
    
    
    function $state<false>(initial: false): false (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(false); constructor(text) { this.`Todo.text: any`text = `
    
    
    function $state<any>(initial: any): any (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(`text: any`text); } `Todo.reset: () => void`reset = () => { this.`Todo.text: any`text = ''; this.`Todo.done: boolean`done = false; } }`

### Built-in classes

Svelte provides reactive implementations of built-in classes like `Set`, `Map`, `Date` and `URL` that can be imported from [`svelte/reactivity`](svelte-reactivity).

## $state.raw

In cases where you don’t want objects and arrays to be deeply reactive you can use `$state.raw`.

State declared with `$state.raw` cannot be mutated; it can only be _reassigned_. In other words, rather than assigning to a property of an object, or using an array method like `push`, replace the object or array altogether if you’d like to update it:
    
    
    let 
    
    
    let person: {
        name: string;
        age: number;
    }

`person = `
    
    
    namespace $state
    function $state<T>(initial: T): T (+1 overload)

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state.`
    
    
    function $state.raw<{
        name: string;
        age: number;
    }>(initial: {
        name: string;
        age: number;
    }): {
        name: string;
        age: number;
    } (+1 overload)

`

Declares state that is _not_ made deeply reactive — instead of mutating it, you must reassign it.

Example:
    
    
    &#x3C;script>
      let items = $state.raw([0]);
    
      const addItem = () => {
    	items = [...items, items.length];
      };
    &#x3C;/script>
    
    &#x3C;button onclick={addItem}>
      {items.join(', ')}
    &#x3C;/button>

<https://svelte.dev/docs/svelte/$state#$state.raw>

@paraminitial The initial value

raw({ `name: string`name: 'Heraclitus', `age: number`age: 49 }); // this will have no effect `
    
    
    let person: {
        name: string;
        age: number;
    }

`person.`age: number`age += 1; // this will work, because we're creating a new person `
    
    
    let person: {
        name: string;
        age: number;
    }

`person = { `name: string`name: 'Heraclitus', `age: number`age: 50 };`

This can improve performance with large arrays and objects that you weren’t planning to mutate anyway, since it avoids the cost of making them reactive. Note that raw state can _contain_ reactive state (for example, a raw array of reactive objects).

As with `$state`, you can declare class fields using `$state.raw`.

## $state.snapshot

To take a static snapshot of a deeply reactive `$state` proxy, use `$state.snapshot`:
    
    
    <script>
    	let counter = $state({ count: 0 });
    
    	function onclick() {
    		// Will log `{ count: ... }` rather than `Proxy { ... }`
    		console.log($state.snapshot(counter));
    	}
    </script>

This is handy when you want to pass some state to an external library or API that doesn’t expect a proxy, such as `structuredClone`.

## $state.eager

When state changes, it may not be reflected in the UI immediately if it is used by an `await` expression, because [updates are synchronized](await-expressions#Synchronized-updates).

In some cases, you may want to update the UI as soon as the state changes. For example, you might want to update a navigation bar when the user clicks on a link, so that they get visual feedback while waiting for the new page to load. To do this, use `$state.eager(value)`:
    
    
    <nav>
    	<a href="/" aria-current={$state.eager(pathname) === '/' ? 'page' : null}>home</a>
    	<a href="/about" aria-current={$state.eager(pathname) === '/about' ? 'page' : null}>about</a>
    </nav>

Use this feature sparingly, and only to provide feedback in response to user action — in general, allowing Svelte to coordinate updates will provide a better user experience.

## Passing state into functions

JavaScript is a _pass-by-value_ language — when you call a function, the arguments are the _values_ rather than the _variables_. In other words:

index
    
    
    /**
     * @param {number} a
     * @param {number} b
     */
    function function add(a: number, b: number): number
    
    @parama 
    
    @paramb 
    
    add(a: number
    
    @parama 
    
    a, b: number
    
    @paramb 
    
    b) {
    	return a: number
    
    @parama 
    
    a + b: number
    
    @paramb 
    
    b;
    }
    
    let let a: numbera = 1;
    let let b: numberb = 2;
    let let total: numbertotal = function add(a: number, b: number): number
    
    @parama 
    
    @paramb 
    
    add(let a: numbera, let b: numberb);
    var console: Console
    
    The console module provides a simple debugging console that is similar to the
    JavaScript console mechanism provided by web browsers.
    
    
    
    
    The module exports two specific components:
    
    
    
    
    
    
        * A Console class with methods such as console.log(), console.error() and console.warn() that can be used to write to any Node.js stream.
    
    
        * A global console instance configured to write to [process.stdout](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and
        [process.stderr](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global console can be used without importing the node:console module.
    
    
    
    
    
    _**Warning**_ : The global console object’s methods are neither consistently
    synchronous like the browser APIs they resemble, nor are they consistently
    asynchronous like all other Node.js streams. See the [note on process I/O](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for
    more information.
    
    
    
    
    Example using the global console:
    
    
    
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: number`total); // 3 `let a: number`a = 3; `let b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: number`total); // still 3!`
    
    
    function function add(a: number, b: number): numberadd(a: numbera: number, b: numberb: number) {
    	return a: numbera + b: numberb;
    }
    
    let let a: numbera = 1;
    let let b: numberb = 2;
    let let total: numbertotal = function add(a: number, b: number): numberadd(let a: numbera, let b: numberb);
    var console: Console
    
    The console module provides a simple debugging console that is similar to the
    JavaScript console mechanism provided by web browsers.
    
    
    
    
    The module exports two specific components:
    
    
    
    
    
    
        * A Console class with methods such as console.log(), console.error() and console.warn() that can be used to write to any Node.js stream.
    
    
        * A global console instance configured to write to [process.stdout](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and
        [process.stderr](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global console can be used without importing the node:console module.
    
    
    
    
    
    _**Warning**_ : The global console object’s methods are neither consistently
    synchronous like the browser APIs they resemble, nor are they consistently
    asynchronous like all other Node.js streams. See the [note on process I/O](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for
    more information.
    
    
    
    
    Example using the global console:
    
    
    
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: number`total); // 3 `let a: number`a = 3; `let b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: number`total); // still 3!`

If `add` wanted to have access to the _current_ values of `a` and `b`, and to return the current `total` value, you would need to use functions instead:

index
    
    
    /**
     * @param {() => number} getA
     * @param {() => number} getB
     */
    function function add(getA: () => number, getB: () => number): () => number
    
    @paramgetA 
    
    @paramgetB 
    
    add(getA: () => number
    
    @paramgetA 
    
    getA, getB: () => number
    
    @paramgetB 
    
    getB) {
    	return () => getA: () => number
    
    @paramgetA 
    
    getA() + getB: () => number
    
    @paramgetB 
    
    getB();
    }
    
    let let a: numbera = 1;
    let let b: numberb = 2;
    let let total: () => numbertotal = function add(getA: () => number, getB: () => number): () => number
    
    @paramgetA 
    
    @paramgetB 
    
    add(() => let a: numbera, () => let b: numberb);
    var console: Console
    
    The console module provides a simple debugging console that is similar to the
    JavaScript console mechanism provided by web browsers.
    
    
    
    
    The module exports two specific components:
    
    
    
    
    
    
        * A Console class with methods such as console.log(), console.error() and console.warn() that can be used to write to any Node.js stream.
    
    
        * A global console instance configured to write to [process.stdout](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and
        [process.stderr](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global console can be used without importing the node:console module.
    
    
    
    
    
    _**Warning**_ : The global console object’s methods are neither consistently
    synchronous like the browser APIs they resemble, nor are they consistently
    asynchronous like all other Node.js streams. See the [note on process I/O](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for
    more information.
    
    
    
    
    Example using the global console:
    
    
    
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: () => number`total()); // 3 `let a: number`a = 3; `let b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: () => number`total()); // 7`
    
    
    function function add(getA: () => number, getB: () => number): () => numberadd(getA: () => numbergetA: () => number, getB: () => numbergetB: () => number) {
    	return () => getA: () => numbergetA() + getB: () => numbergetB();
    }
    
    let let a: numbera = 1;
    let let b: numberb = 2;
    let let total: () => numbertotal = function add(getA: () => number, getB: () => number): () => numberadd(() => let a: numbera, () => let b: numberb);
    var console: Console
    
    The console module provides a simple debugging console that is similar to the
    JavaScript console mechanism provided by web browsers.
    
    
    
    
    The module exports two specific components:
    
    
    
    
    
    
        * A Console class with methods such as console.log(), console.error() and console.warn() that can be used to write to any Node.js stream.
    
    
        * A global console instance configured to write to [process.stdout](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and
        [process.stderr](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global console can be used without importing the node:console module.
    
    
    
    
    
    _**Warning**_ : The global console object’s methods are neither consistently
    synchronous like the browser APIs they resemble, nor are they consistently
    asynchronous like all other Node.js streams. See the [note on process I/O](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for
    more information.
    
    
    
    
    Example using the global console:
    
    
    
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: () => number`total()); // 3 `let a: number`a = 3; `let b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`let total: () => number`total()); // 7`

State in Svelte is no different — when you reference something declared with the `$state` rune...
    
    
    let let a: numbera = 
    
    
    function $state<1>(initial: 1): 1 (+1 overload)

` namespace $state`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(1); let `let b: number`b = `
    
    
    function $state<2>(initial: 2): 2 (+1 overload)

` namespace $state

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(2);

...you’re accessing its _current value_.

Note that ‘functions’ is broad — it encompasses properties of proxies and [`get`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Functions/get)/[`set`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Functions/set) properties...

index
    
    
    /**
     * @param {{ a: number, b: number }} input
     */
    function 
    
    
    function add(input: {
        a: number;
        b: number;
    }): {
        readonly value: number;
    }

`

@paraminput

add(`
    
    
    input: {
        a: number;
        b: number;
    }

`

@paraminput

input) { return { get `value: number`value() { return `
    
    
    input: {
        a: number;
        b: number;
    }

`

@paraminput

input.`a: number`a + `
    
    
    input: {
        a: number;
        b: number;
    }

`

@paraminput

input.`b: number`b; } }; } let `
    
    
    module input
    let input: {
        a: number;
        b: number;
    }

`input = `
    
    
    function $state<{
        a: number;
        b: number;
    }>(initial: {
        a: number;
        b: number;
    }): {
        a: number;
        b: number;
    } (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state({ `a: number`a: 1, `b: number`b: 2 }); let `
    
    
    let total: {
        readonly value: number;
    }

`total = `
    
    
    function add(input: {
        a: number;
        b: number;
    }): {
        readonly value: number;
    }

`

@paraminput

add(`
    
    
    module input
    let input: {
        a: number;
        b: number;
    }

`input); `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`
    
    
    let total: {
        readonly value: number;
    }

`total.`value: number`value); // 3 `
    
    
    module input
    let input: {
        a: number;
        b: number;
    }

`input.`a: number`a = 3; `
    
    
    module input
    let input: {
        a: number;
        b: number;
    }

`input.`b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`
    
    
    let total: {
        readonly value: number;
    }

`total.`value: number`value); // 7`
    
    
    function 
    
    
    function add(input: {
        a: number;
        b: number;
    }): {
        readonly value: number;
    }

`add(`
    
    
    input: {
        a: number;
        b: number;
    }

`input: { `a: number`a: number, `b: number`b: number }) { return { get `value: number`value() { return `
    
    
    input: {
        a: number;
        b: number;
    }

`input.`a: number`a + `
    
    
    input: {
        a: number;
        b: number;
    }

`input.`b: number`b; } }; } let `
    
    
    let input: {
        a: number;
        b: number;
    }

`input = `
    
    
    function $state<{
        a: number;
        b: number;
    }>(initial: {
        a: number;
        b: number;
    }): {
        a: number;
        b: number;
    } (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state({ `a: number`a: 1, `b: number`b: 2 }); let `
    
    
    let total: {
        readonly value: number;
    }

`total = `
    
    
    function add(input: {
        a: number;
        b: number;
    }): {
        readonly value: number;
    }

`add(`
    
    
    let input: {
        a: number;
        b: number;
    }

`input); `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`
    
    
    let total: {
        readonly value: number;
    }

`total.`value: number`value); // 3 `
    
    
    let input: {
        a: number;
        b: number;
    }

`input.`a: number`a = 3; `
    
    
    let input: {
        a: number;
        b: number;
    }

`input.`b: number`b = 4; `var console: Console`

The `console` module provides a simple debugging console that is similar to the JavaScript console mechanism provided by web browsers.

The module exports two specific components:

  * A `Console` class with methods such as `console.log()`, `console.error()` and `console.warn()` that can be used to write to any Node.js stream.
  * A global `console` instance configured to write to [`process.stdout`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and [`process.stderr`](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global `console` can be used without importing the `node:console` module.



_**Warning**_ : The global console object’s methods are neither consistently synchronous like the browser APIs they resemble, nor are they consistently asynchronous like all other Node.js streams. See the [`note on process I/O`](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for more information.

Example using the global `console`:
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(`
    
    
    let total: {
        readonly value: number;
    }

`total.`value: number`value); // 7`

...though if you find yourself writing code like that, consider using classes instead.

## Passing state across modules

You can declare state in `.svelte.js` and `.svelte.ts` files, but you can only _export_ that state if it’s not directly reassigned. In other words you can’t do this:

state.svelte
    
    
    export let let count: numbercount = 
    
    
    function $state<0>(initial: 0): 0 (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(0); export function `function increment(): void`increment() { `let count: number`count += 1; }`

That’s because every reference to `count` is transformed by the Svelte compiler — the code above is roughly equivalent to this:

state.svelte
    
    
    export let let count: Signal<number>count = const $: Svelte$.Svelte.state<number>(value?: number | undefined): Signal<number>state(0);
    
    export function function increment(): voidincrement() {
    	const $: Svelte$.Svelte.set<number>(source: Signal<number>, value: number): voidset(let count: Signal<number>count, const $: Svelte$.Svelte.get<number>(source: Signal<number>): numberget(let count: Signal<number>count) + 1);
    }

> You can see the code Svelte generates by clicking the ‘JS Output’ tab in the [playground](/playground).

Since the compiler only operates on one file at a time, if another file imports `count` Svelte doesn’t know that it needs to wrap each reference in `$.get` and `$.set`:
    
    
    import { let count: numbercount } from './state.svelte.js';
    
    var console: Console
    
    The console module provides a simple debugging console that is similar to the
    JavaScript console mechanism provided by web browsers.
    
    
    
    
    The module exports two specific components:
    
    
    
    
    
    
        * A Console class with methods such as console.log(), console.error() and console.warn() that can be used to write to any Node.js stream.
    
    
        * A global console instance configured to write to [process.stdout](https://nodejs.org/docs/latest-v20.x/api/process.html#processstdout) and
        [process.stderr](https://nodejs.org/docs/latest-v20.x/api/process.html#processstderr). The global console can be used without importing the node:console module.
    
    
    
    
    
    _**Warning**_ : The global console object’s methods are neither consistently
    synchronous like the browser APIs they resemble, nor are they consistently
    asynchronous like all other Node.js streams. See the [note on process I/O](https://nodejs.org/docs/latest-v20.x/api/process.html#a-note-on-process-io) for
    more information.
    
    
    
    
    Example using the global console:
    
    
    
    
    
    console.log('hello world');
    // Prints: hello world, to stdout
    console.log('hello %s', 'world');
    // Prints: hello world, to stdout
    console.error(new Error('Whoops, something bad happened'));
    // Prints error message and stack trace to stderr:
    //   Error: Whoops, something bad happened
    //     at [eval]:5:15
    //     at Script.runInThisContext (node:vm:132:18)
    //     at Object.runInThisContext (node:vm:309:38)
    //     at node:internal/process/execution:77:19
    //     at [eval]-wrapper:6:22
    //     at evalScript (node:internal/process/execution:76:60)
    //     at node:internal/main/eval_string:23:3
    
    const name = 'Will Robinson';
    console.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to stderr

Example using the `Console` class:
    
    
    const out = getStreamSomehow();
    const err = getStreamSomehow();
    const myConsole = new console.Console(out, err);
    
    myConsole.log('hello world');
    // Prints: hello world, to out
    myConsole.log('hello %s', 'world');
    // Prints: hello world, to out
    myConsole.error(new Error('Whoops, something bad happened'));
    // Prints: [Error: Whoops, something bad happened], to err
    
    const name = 'Will Robinson';
    myConsole.warn(`Danger ${name}! Danger!`);
    // Prints: Danger Will Robinson! Danger!, to err

@see[source](https://github.com/nodejs/node/blob/v20.11.1/lib/console.js)

console.`Console.log(message?: any, ...optionalParams: any[]): void (+1 overload)`

Prints to `stdout` with newline. Multiple arguments can be passed, with the first used as the primary message and all additional used as substitution values similar to [`printf(3)`](http://man7.org/linux/man-pages/man3/printf.3.html) (the arguments are all passed to [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args)).
    
    
    const count = 5;
    console.log('count: %d', count);
    // Prints: count: 5, to stdout
    console.log('count:', count);
    // Prints: count: 5, to stdout

See [`util.format()`](https://nodejs.org/docs/latest-v20.x/api/util.html#utilformatformat-args) for more information.

@sincev0.1.100

log(typeof `let count: number`count); // 'object', not 'number'`

This leaves you with two options for sharing state between modules — either don’t reassign it...
    
    
    // This is allowed — since we're updating
    // `counter.count` rather than `counter`,
    // Svelte doesn't wrap it in `$.state`
    export const 
    
    
    const counter: {
        count: number;
    }

`counter = `
    
    
    function $state<{
        count: number;
    }>(initial: {
        count: number;
    }): {
        count: number;
    } (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state({ `count: number`count: 0 }); export function `function increment(): void`increment() { `
    
    
    const counter: {
        count: number;
    }

`counter.`count: number`count += 1; }`

...or don’t directly export it:
    
    
    let let count: numbercount = 
    
    
    function $state<0>(initial: 0): 0 (+1 overload)
    namespace $state

`

Declares reactive state.

Example:
    
    
    let count = $state(0);

<https://svelte.dev/docs/svelte/$state>

@paraminitial The initial value

$state(0); export function `function getCount(): number`getCount() { return `let count: number`count; } export function `function increment(): void`increment() { `let count: number`count += 1; }`

[ Edit this page on GitHub](https://github.com/sveltejs/svelte/edit/main/documentation/docs/02-runes/02-$state.md) [ llms.txt](/docs/svelte/$state/llms.txt)

previous next

[What are runes?](/docs/svelte/what-are-runes) [$derived](/docs/svelte/$derived)
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
