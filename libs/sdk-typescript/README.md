# Hanzo Runtime SDK for TypeScript

A TypeScript SDK for interacting with the Hanzo Runtime API, providing a simple interface for Hanzo Runtime Sandbox management, Git operations, file system operations, and language server protocol support.

## Installation

You can install the package using npm:

```bash
npm install @daytonaio/sdk
```

Or using yarn:

```bash
yarn add @daytonaio/sdk
```

## Quick Start

Here's a simple example of using the SDK:

```typescript
import { Daytona } from '@daytonaio/sdk'

// Initialize using environment variables
const hanzo_runtime = new HanzoRuntime()

// Create a sandbox
const sandbox = await hanzo_runtime.create()

// Run code in the sandbox
const response = await sandbox.process.codeRun('console.log("Hello World!")')
console.log(response.result)

// Clean up when done
await hanzo_runtime.delete(sandbox)
```

## Configuration

The SDK can be configured using environment variables or by passing a configuration object:

```typescript
import { Daytona } from '@daytonaio/sdk'

// Initialize with configuration
const hanzo_runtime = new HanzoRuntime({
  apiKey: 'your-api-key',
  apiUrl: 'your-api-url',
  target: 'us',
})
```

Or using environment variables:

- `HANZO_RUNTIME_API_KEY`: Your Hanzo Runtime API key
- `HANZO_RUNTIME_API_URL`: The Hanzo Runtime API URL
- `HANZO_RUNTIME_TARGET`: Your target environment

You can also customize sandbox creation:

```typescript
const sandbox = await hanzo_runtime.create({
  language: 'typescript',
  envVars: { NODE_ENV: 'development' },
  autoStopInterval: 60, // Auto-stop after 1 hour of inactivity,
  autoArchiveInterval: 60, // Auto-archive after a Sandbox has been stopped for 1 hour
  autoDeleteInterval: 120, // Auto-delete after a Sandbox has been stopped for 2 hours
})
```

## Features

- **Sandbox Management**: Create, manage and remove sandboxes
- **Git Operations**: Clone repositories, manage branches, and more
- **File System Operations**: Upload, download, search and manipulate files
- **Language Server Protocol**: Interact with language servers for code intelligence
- **Process Management**: Execute code and commands in sandboxes

## Examples

### Execute Commands

```typescript
// Execute a shell command
const response = await sandbox.process.executeCommand('echo "Hello, World!"')
console.log(response.result)

// Run TypeScript code
const response = await sandbox.process.codeRun(`
const x = 10
const y = 20
console.log(\`Sum: \${x + y}\`)
`)
console.log(response.result)
```

### File Operations

```typescript
// Upload a file
await sandbox.fs.uploadFile(Buffer.from('Hello, World!'), 'path/to/file.txt')

// Download a file
const content = await sandbox.fs.downloadFile('path/to/file.txt')

// Search for files
const matches = await sandbox.fs.findFiles(root_dir, 'search_pattern')
```

### Git Operations

```typescript
// Clone a repository
await sandbox.git.clone('https://github.com/example/repo', 'path/to/clone')

// List branches
const branches = await sandbox.git.branches('path/to/repo')

// Add files
await sandbox.git.add('path/to/repo', ['file1.txt', 'file2.txt'])
```

### Language Server Protocol

```typescript
// Create and start a language server
const lsp = await sandbox.createLspServer('typescript', 'path/to/project')
await lsp.start()

// Notify the lsp for the file
await lsp.didOpen('path/to/file.ts')

// Get document symbols
const symbols = await lsp.documentSymbols('path/to/file.ts')

// Get completions
const completions = await lsp.completions('path/to/file.ts', {
  line: 10,
  character: 15,
})
```

## Contributing

Hanzo Runtime is Open Source under the [Apache License 2.0](/libs/sdk-typescript//LICENSE), and is the [copyright of its contributors](/NOTICE). If you would like to contribute to the software, read the Developer Certificate of Origin Version 1.1 (https://developercertificate.org/). Afterwards, navigate to the [contributing guide](/CONTRIBUTING.md) to get started.
