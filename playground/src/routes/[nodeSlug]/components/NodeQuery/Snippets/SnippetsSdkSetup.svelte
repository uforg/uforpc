<script lang="ts">
  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Code from "$lib/components/Code.svelte";

  const dartPackageName = $derived.by(
    () => uiStore.codeSnippetsSdkDartPackageName.trim() || "your_dart_package",
  );
  const golangPackageName = $derived.by(
    () => uiStore.codeSnippetsSdkGolangPackageName.trim() || "yourpkg",
  );

  const tsSetup = $derived.by(
    () => `// 1) Place the generated file in your project, e.g. src/lib/uforpc-client-sdk.ts
import { NewClient } from "./path/to/uforpc-client-sdk.ts";

// 2) Build the client
const client = NewClient("${store.baseUrl}").build();`,
  );

  const goSetup = $derived.by(
    () => `// Put the generated file into your Go module, inside any package you want

package main

import "yourmodule/${golangPackageName}"

func main() {
	client := ${golangPackageName}.NewClient("${store.baseUrl}").Build()
	_ = client // ready to use
}`,
  );

  const dartYaml = $derived.by(
    () => `# Make sure to use the package name you chose when you downloaded the SDK
# Currently, the package name is ${dartPackageName}
dependencies:
  ${dartPackageName}:
    path: ./path/to/${dartPackageName}`,
  );

  const dartSetup = $derived.by(
    () => `import "package:${dartPackageName}/client.dart" as urpc;

final client = urpc.NewClient("${store.baseUrl}").build();`,
  );
</script>

<div class="prose prose-sm text-base-content max-w-none space-y-4">
  {#if uiStore.codeSnippetsSdkLang === "typescript-client"}
    <h3>TypeScript setup</h3>
    <p>
      The SDK is a single <code>.ts</code> file with no external dependencies.
      Move it to your project and import <code>NewClient</code> from it.
    </p>
    <ol class="list-decimal pl-5">
      <li>
        Move the generated file (for example, <code>uforpc-client-sdk.ts</code>)
        to a folder in your project.
      </li>
      <li>Import and build the client:</li>
    </ol>
    <div class="not-prose">
      <Code code={tsSetup} lang="ts" />
    </div>
    <p>No additional configuration or dependencies are required.</p>
  {:else if uiStore.codeSnippetsSdkLang === "golang-client"}
    <h3>Go setup</h3>
    <p>
      The SDK is a single <code>.go</code> file with no external dependencies. Place
      it inside the Go package where you will use it.
    </p>
    <p>
      The package name you selected in the previous step is
      <code>{golangPackageName}</code>.
    </p>
    <ol class="list-decimal pl-5">
      <li>Move the file into any go package inside your module.</li>
      <li>Import the package and build the client where you want to use it:</li>
    </ol>
    <div class="not-prose">
      <Code code={goSetup} lang="go" />
    </div>
    <p>
      If you keep the generated client in a different package, import that
      package and call <code>NewClient</code> through it.
    </p>
  {:else if uiStore.codeSnippetsSdkLang === "dart-client"}
    <h3>Dart setup</h3>
    <p>
      The download is a zip containing a full Dart package. Unzip it and add it
      to your project as a local dependency using the package name you chose
      when you downloaded the SDK.
    </p>
    <ol class="list-decimal pl-5">
      <li>Unzip the archive to a local folder.</li>
      <li>
        In your application's <code>pubspec.yaml</code>, add a local dependency
        pointing to the unzipped folder:
      </li>
    </ol>
    <div class="not-prose">
      <Code code={dartYaml} lang="yaml" />
    </div>
    <ol class="list-decimal pl-5" start="3">
      <li>Run <code>dart pub get</code> or <code>flutter pub get</code>.</li>
      <li>Import the client package and build the client:</li>
    </ol>
    <div class="not-prose">
      <Code code={dartSetup} lang="dart" />
    </div>
  {/if}
</div>
