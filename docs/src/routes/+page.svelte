<script lang="ts">
  import {
    ArrowRight,
    BookOpenText,
    Check,
    Code,
    FileCode,
    Github,
    Layers,
    Play,
    Puzzle,
    Rocket,
    Shield,
    Sparkles,
    Users,
    Workflow,
    Zap,
  } from "@lucide/svelte";

  import GithubStarsImg from "$lib/components/GithubStarsImg.svelte";
  import H1 from "$lib/components/H1.svelte";
  import H2 from "$lib/components/H2.svelte";
  import UfoAbduction from "$lib/components/UfoAbduction.svelte";
</script>

<!-- Hero Section -->
<section class="hero min-h-screen bg-base-100">
  <div class="hero-content text-center">
    <div class="max-w-6xl">
      <!-- Main Header -->
      <div class="mb-12">
        <H1 class="mb-6 text-5xl font-bold desk:text-7xl">UFO RPC</H1>
        <H2
          class="mx-auto mb-8 max-w-4xl text-xl text-base-content/80 desk:text-3xl"
        >
          The RPC framework that puts developer experience first
        </H2>

        <!-- UFO Animation -->
        <div class="mb-8 flex justify-center">
          <UfoAbduction width={250} />
        </div>

        <!-- Problem Statement -->
        <div
          class="mx-auto mb-8 alert flex items-center justify-center alert-warning shadow-lg"
        >
          <Sparkles class="size-5" />
          <span class="text-sm desk:text-base">
            <strong>Tired of API's complexity?</strong>
            Get type safety without the headaches.
          </span>
        </div>
      </div>

      <!-- CTA Buttons -->
      <div
        class="mb-16 flex flex-col items-center justify-center gap-4 desk:flex-row"
      >
        <a
          href="https://github.com/uforg/uforpc"
          target="_blank"
          class="btn w-full gap-2 shadow-lg transition-all btn-lg btn-primary hover:shadow-xl desk:w-auto"
        >
          <Github class="size-5" />
          <span>View on GitHub</span>
          <GithubStarsImg />
        </a>
        <a
          href="/docs"
          class="btn w-full gap-2 transition-all btn-outline btn-lg hover:btn-primary desk:w-auto"
        >
          <BookOpenText class="size-5" />
          Documentation
        </a>
      </div>
    </div>
  </div>
</section>

<!-- Pain Points vs Solution -->
<section class="bg-base-200 py-24">
  <div class="container mx-auto px-4">
    <div class="mb-16 text-center">
      <h2 class="mb-4 text-4xl font-bold">Why Developers Choose UFO RPC</h2>
      <p class="mx-auto max-w-3xl text-xl text-base-content/70">
        Stop fighting with complex RPC frameworks. UFO RPC gives you all the
        benefits with none of the pain.
      </p>
    </div>

    <div class="mx-auto grid max-w-6xl grid-cols-1 gap-8 desk:grid-cols-2">
      <!-- Before: Pain Points -->
      <div class="card border border-error/20 bg-error/5">
        <div class="card-body">
          <h3 class="mb-6 card-title text-error">
            <span class="text-2xl">❌</span>
            Traditional RPC Problems
          </h3>
          <ul class="space-y-3">
            <li class="flex items-start gap-2">
              <span class="mt-1 text-error">•</span>
              <span
                >Complex protobuf schemas that are hard to read and write</span
              >
            </li>
            <li class="flex items-start gap-2">
              <span class="mt-1 text-error">•</span>
              <span>Binary protocols that are impossible to debug</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="mt-1 text-error">•</span>
              <span>Manual client code that's error-prone and tedious</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="mt-1 text-error">•</span>
              <span>Poor tooling and development experience</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="mt-1 text-error">•</span>
              <span>Documentation that goes stale immediately</span>
            </li>
          </ul>
        </div>
      </div>

      <!-- After: UFO RPC Solution -->
      <div class="card border border-success/20 bg-success/5">
        <div class="card-body">
          <h3 class="mb-6 card-title text-success">
            <span class="text-2xl">✅</span>
            UFO RPC Solution
          </h3>
          <ul class="space-y-3">
            <li class="flex items-start gap-2">
              <Check class="mt-0.5 size-5 flex-shrink-0 text-success" />
              <span
                >Human-readable <code>.urpc</code> schemas anyone can understand</span
              >
            </li>
            <li class="flex items-start gap-2">
              <Check class="mt-0.5 size-5 flex-shrink-0 text-success" />
              <span>Standard HTTP/JSON that you can debug with curl</span>
            </li>
            <li class="flex items-start gap-2">
              <Check class="mt-0.5 size-5 flex-shrink-0 text-success" />
              <span>Fully-typed clients generated automatically</span>
            </li>
            <li class="flex items-start gap-2">
              <Check class="mt-0.5 size-5 flex-shrink-0 text-success" />
              <span>VS Code extension, LSP, and interactive playground</span>
            </li>
            <li class="flex items-start gap-2">
              <Check class="mt-0.5 size-5 flex-shrink-0 text-success" />
              <span>Documentation lives in your schema, always up-to-date</span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Features Deep Dive -->
<section class="py-24">
  <div class="container mx-auto px-4">
    <div class="mb-16 text-center">
      <h2 class="mb-4 text-4xl font-bold">Built for Developer Happiness</h2>
      <p class="mx-auto max-w-3xl text-xl text-base-content/70">
        Every feature is designed to make your development workflow faster,
        safer, and more enjoyable.
      </p>
    </div>

    <!-- Feature Cards Grid -->
    <div class="mb-16 grid grid-cols-1 gap-8 desk:grid-cols-3">
      <!-- Type Safety -->
      <div
        class="card border border-primary/20 bg-base-100 shadow-xl transition-all duration-300 hover:shadow-2xl"
      >
        <div class="card-body">
          <div class="mb-4 w-fit rounded-full bg-primary/10 p-3">
            <Code class="size-8 text-primary" />
          </div>
          <h3 class="mb-2 card-title text-xl">End-to-End Type Safety</h3>
          <p class="mb-4 text-base-content/70">
            One schema, multiple languages. Changes propagate everywhere
            automatically.
          </p>
          <div class="flex flex-wrap gap-1">
            <div class="badge badge-sm badge-primary">Go</div>
            <div class="badge badge-sm badge-secondary">TypeScript</div>
            <div class="badge badge-sm badge-accent">Dart</div>
          </div>
        </div>
      </div>

      <!-- Developer Experience -->
      <div
        class="card border border-secondary/20 bg-base-100 shadow-xl transition-all duration-300 hover:shadow-2xl"
      >
        <div class="card-body">
          <div class="mb-4 w-fit rounded-full bg-secondary/10 p-3">
            <Zap class="size-8 text-secondary" />
          </div>
          <h3 class="mb-2 card-title text-xl">Premium DX</h3>
          <p class="mb-4 text-base-content/70">
            LSP autocompletion, syntax highlighting, real-time validation, and
            more.
          </p>
          <div class="flex flex-wrap gap-1">
            <div class="badge badge-outline badge-sm">VS Code</div>
            <div class="badge badge-outline badge-sm">LSP</div>
            <div class="badge badge-outline badge-sm">Formatter</div>
          </div>
        </div>
      </div>

      <!-- Production Ready -->
      <div
        class="card border border-accent/20 bg-base-100 shadow-xl transition-all duration-300 hover:shadow-2xl"
      >
        <div class="card-body">
          <div class="mb-4 w-fit rounded-full bg-accent/10 p-3">
            <Shield class="size-8 text-accent" />
          </div>
          <h3 class="mb-2 card-title text-xl">Production Hardened</h3>
          <p class="mb-4 text-base-content/70">
            Built-in retries, timeouts, circuit breakers, and observability
            hooks.
          </p>
          <div class="flex flex-wrap gap-1">
            <div class="badge badge-sm badge-success">Retries</div>
            <div class="badge badge-sm badge-warning">Timeouts</div>
            <div class="badge badge-sm badge-info">Metrics</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Detailed Features -->
    <div class="grid grid-cols-1 gap-8 desk:grid-cols-2">
      <!-- Interactive Playground -->
      <div class="card border border-base-300 bg-base-100 shadow-lg">
        <div class="card-body">
          <div class="mb-4 flex items-center gap-3">
            <div class="rounded-lg bg-primary/10 p-2">
              <Workflow class="size-6 text-primary" />
            </div>
            <h3 class="card-title text-2xl">Interactive Playground</h3>
          </div>
          <p class="mb-6 text-base-content/70">
            Every UFO RPC project generates a beautiful, static web portal that
            serves as your API documentation, testing interface, and developer
            onboarding tool.
          </p>

          <div class="mockup-browser mb-4 border bg-base-300">
            <div class="mockup-browser-toolbar">
              <div class="input">https://api.yourproject.com</div>
            </div>
            <div class="bg-base-200 px-4 py-6 text-sm">
              <div class="mb-2 flex items-center gap-2">
                <span class="badge badge-sm badge-success">POST</span>
                <span>/CreateUser</span>
              </div>
              <div class="text-base-content/70">
                Auto-generated forms • Live curl commands • Code snippets
              </div>
            </div>
          </div>

          <ul class="space-y-2 text-sm">
            <li class="flex items-center gap-2">
              <Check class="size-4 text-success" />
              <span>Auto-generated forms for testing endpoints</span>
            </li>
            <li class="flex items-center gap-2">
              <Check class="size-4 text-success" />
              <span>Copy-paste ready curl commands</span>
            </li>
            <li class="flex items-center gap-2">
              <Check class="size-4 text-success" />
              <span>Client SDK code snippets</span>
            </li>
            <li class="flex items-center gap-2">
              <Check class="size-4 text-success" />
              <span>Live documentation that never goes stale</span>
            </li>
          </ul>
        </div>
      </div>

      <!-- Schema-First Development -->
      <div class="card border border-base-300 bg-base-100 shadow-lg">
        <div class="card-body">
          <div class="mb-4 flex items-center gap-3">
            <div class="rounded-lg bg-secondary/10 p-2">
              <Layers class="size-6 text-secondary" />
            </div>
            <h3 class="card-title text-2xl">Schema as Source of Truth</h3>
          </div>
          <p class="mb-6 text-base-content/70">
            Define your API once in human-readable <code>.urpc</code> files. Everything
            else is generated automatically, keeping your entire stack in perfect
            sync.
          </p>

          <div class="mockup-code mb-4 text-sm">
            <pre data-prefix="// "><code>schema.urpc</code></pre>
            <pre data-prefix=""><code>type User &#123;</code></pre>
            <pre data-prefix=""><code>  id: string</code></pre>
            <pre data-prefix=""><code>  name: string</code></pre>
            <pre data-prefix=""><code>  email: string</code></pre>
            <pre data-prefix=""><code>&#125;</code></pre>
          </div>

          <div class="alert alert-info">
            <Puzzle class="size-5" />
            <span class="text-sm">
              <strong>One schema, infinite possibilities:</strong>
              Servers, clients, docs, and playground all stay perfectly synchronized.
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Getting Started -->
<section class="bg-base-200 py-24">
  <div class="container mx-auto px-4">
    <div class="mb-16 text-center">
      <h2 class="mb-4 text-4xl font-bold">Start Building in 60 Seconds</h2>
      <p class="mx-auto max-w-2xl text-xl text-base-content/70">
        No complex setup, no configuration hell. Just pure productivity.
      </p>
    </div>

    <div class="mx-auto max-w-5xl">
      <!-- Installation Steps -->
      <div class="steps mb-12 w-full desk:steps-horizontal">
        <div class="step step-primary">
          <div class="max-w-xs text-left">
            <h3 class="mb-2 text-lg font-bold">Install CLI</h3>
            <div class="mockup-code text-xs">
              <pre>
                <code>
                  // TODO
                </code>
              </pre>
            </div>
          </div>
        </div>
        <div class="step step-primary">
          <div class="max-w-xs text-left">
            <h3 class="mb-2 text-lg font-bold">Initialize Project</h3>
            <div class="mockup-code text-xs">
              <pre><code>urpc init .</code></pre>
            </div>
          </div>
        </div>
        <div class="step step-primary">
          <div class="max-w-xs text-left">
            <h3 class="mb-2 text-lg font-bold">Generate & Use</h3>
            <div class="mockup-code text-xs">
              <pre><code>urpc generate</code></pre>
            </div>
          </div>
        </div>
      </div>

      <!-- What You Get -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h3 class="mb-6 card-title justify-center text-2xl">
            <Rocket class="size-6" />
            What You Get Out of the Box
          </h3>

          <div class="grid grid-cols-1 gap-6 desk:grid-cols-2">
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <div class="badge badge-primary">Server</div>
                <span>Production-ready Go & TypeScript HTTP server</span>
              </div>
              <div class="flex items-center gap-3">
                <div class="badge badge-secondary">Client</div>
                <span>Type-safe Go, TypeScript & Dart clients</span>
              </div>
              <div class="flex items-center gap-3">
                <div class="badge badge-accent">Docs</div>
                <span>Interactive API documentation</span>
              </div>
            </div>
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <div class="badge badge-success">Playground</div>
                <span>Web-based testing interface</span>
              </div>
              <div class="flex items-center gap-3">
                <div class="badge badge-warning">OpenAPI</div>
                <span>Compatible with existing tools</span>
              </div>
              <div class="flex items-center gap-3">
                <div class="badge badge-info">Resilience</div>
                <span>Built-in retries & error handling</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- CTA -->
      <div class="mt-12 text-center">
        <a
          href="/docs"
          class="btn gap-2 shadow-lg btn-lg btn-primary hover:shadow-xl"
        >
          <BookOpenText class="size-5" />
          View Complete Guide
          <ArrowRight class="size-4" />
        </a>
      </div>
    </div>
  </div>
</section>

<!-- Community & Open Source -->
<section class="py-24">
  <div class="container mx-auto px-4 text-center">
    <div class="mx-auto max-w-4xl">
      <h2 class="mb-8 text-4xl font-bold">Join the Revolution</h2>
      <p class="mb-12 text-xl text-base-content/70">
        UFO RPC is 100% open source with a thriving community of developers who
        believe APIs should be simple, not complex.
      </p>

      <!-- GitHub Stats -->
      <div
        class="mb-12 flex flex-col items-center justify-center gap-4 desk:flex-row"
      >
        <a
          href="https://github.com/uforg/uforpc"
          target="_blank"
          class="btn gap-2 btn-ghost btn-lg"
        >
          <Github class="size-5" />
          <span>Star on GitHub</span>
          <GithubStarsImg />
        </a>
        <a
          href="https://github.com/uforg/uforpc/issues"
          target="_blank"
          class="btn gap-2 btn-ghost btn-lg"
        >
          Report Issues
        </a>
        <a
          href="https://github.com/uforg/uforpc/discussions"
          target="_blank"
          class="btn gap-2 btn-ghost btn-lg"
        >
          Discussions
        </a>
      </div>

      <!-- License & Philosophy -->
      <div class="mb-12 grid grid-cols-1 gap-6 desk:grid-cols-3">
        <div class="stats bg-base-100 shadow">
          <div class="stat">
            <div class="stat-figure text-primary">
              <Shield class="size-8" />
            </div>
            <div class="stat-title">License</div>
            <div class="stat-value text-primary">AGPL-3.0</div>
            <div class="stat-desc">Generated code: MIT</div>
          </div>
        </div>

        <div class="stats bg-base-100 shadow">
          <div class="stat">
            <div class="stat-figure text-secondary">
              <Users class="size-8" />
            </div>
            <div class="stat-title">Philosophy</div>
            <div class="stat-value text-secondary">DX First</div>
            <div class="stat-desc">Developer Experience</div>
          </div>
        </div>

        <div class="stats bg-base-100 shadow">
          <div class="stat">
            <div class="stat-figure text-accent">
              <FileCode class="size-8" />
            </div>
            <div class="stat-title">Languages</div>
            <div class="stat-value text-accent">3+</div>
            <div class="stat-desc">Go, TypeScript, Dart</div>
          </div>
        </div>
      </div>

      <!-- Final CTA -->
      <div class="card bg-primary text-primary-content shadow-2xl">
        <div class="card-body text-center">
          <h3 class="mb-4 card-title justify-center text-2xl">
            Ready to Build Something Amazing?
          </h3>
          <p class="mb-6 text-primary-content/90">
            Join the developers who've already discovered the joy of simple,
            powerful APIs.
          </p>
          <div
            class="flex flex-col items-center gap-3 desk:flex-row desk:justify-center"
          >
            <a href="/docs" class="btn gap-2 btn-lg btn-secondary">
              <Rocket class="size-5" />
              Get Started Now
            </a>
            <a
              href="https://playground.uforpc.uforg.dev"
              target="_blank"
              class="btn gap-2 border-primary-content text-primary-content btn-outline btn-lg hover:bg-primary-content hover:text-primary"
            >
              <Play class="size-5" />
              Try the Playground
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
