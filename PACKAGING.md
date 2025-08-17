# Packaging Guidelines for Runtime

The Runtime team appreciates any efforts to make our software more accessible to users on various platforms.
While we encourage packaging and distribution of our open-source project, we have some important guidelines, particularly regarding naming.

## Critical Naming Guideline

**Important**: While you are free to package and distribute our software, you **MUST NOT** name your package `runtime` or, in any way, suggest that, the package you distribute, is an official distribution of `runtime`. This restriction is to prevent confusion and maintain the integrity of our project identity.

- Acceptable: "unofficial-runtime-package", "unofficial-runtime-distribution", etc.
- Not Acceptable: "runtime", "official-runtime", etc.

## General Guidelines

1. **License Compliance**: Ensure that the AGPL 3.0/Apache 2.0 license is included with the package and that all copyright notices are preserved.

2. **Version Accuracy**: Use the exact version number of Runtime that you are packaging. Do not modify the version number or add custom suffixes without explicit permission.

3. **Dependencies**: Include all necessary dependencies as specified in our project documentation. Do not add extra dependencies without consulting the project maintainers.

4. **Modifications**: If you need to make any modifications to the source code for packaging purposes, please document these changes clearly and consider submitting them as pull requests to the main project.

5. **Standard Note**: Please include the following standard note in your package description or metadata:

   ```
   This package contains an unofficial distribution of Runtime, an open source project
   developed by Hanzo Industries Inc. This package is not officially supported or endorsed
   by the Runtime project. For the official version, please visit https://github.com/hanzoai/runtime.
   ```

## Feedback and Questions

If you have any questions about packaging Runtime or need clarification on these guidelines, especially regarding naming conventions, please open an issue in our GitHub repository.

We appreciate your contribution to making Runtime more accessible to users across different platforms, while respecting our project's identity!
