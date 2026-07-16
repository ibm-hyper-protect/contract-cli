<!--
Copyright (c) 2026 IBM Corp.
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Contributing to contract-cli

Thank you for considering contributing to `contract-cli`! We appreciate your time and effort in helping improve this project. This guide will help you understand our development process and how to contribute effectively.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Questions](#questions)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the maintainers listed in [MAINTAINERS.md](MAINTAINERS.md).

## How Can I Contribute?

### Reporting Bugs

We use GitHub issue templates to ensure we collect all necessary information. To report a bug:

1. **Check for existing issues**: Search [existing issues](https://github.com/ibm-hyper-protect/contract-cli/issues) to avoid duplicates
2. **Use the bug report template**: Click [here](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=bug_report.yml) or select "Bug Report" when creating a new issue
3. **Fill out all required fields**: The template will guide you through providing:
   - Bug description and impact
   - Steps to reproduce
   - Expected vs actual behavior
   - CLI command and flags used
   - Environment details (CLI version, OS, OpenSSL version, etc.)
   - Target platform (ccrt, ccrv, ccco, or hpvs)

**Important**: For security vulnerabilities, **do NOT create a public issue**. Instead, report them via [GitHub Security Advisories](https://github.com/ibm-hyper-protect/contract-cli/security/advisories/new) or follow our [Security Policy](SECURITY.md).

### Suggesting Enhancements

To suggest a new feature or enhancement:

1. **Check for existing requests**: Search [existing issues](https://github.com/ibm-hyper-protect/contract-cli/issues) to see if it's already been suggested
2. **Use the feature request template**: Click [here](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=feature_request.yml) or select "Feature Request" when creating a new issue
3. **Provide detailed information**: The template will guide you through:
   - Problem statement and motivation
   - Proposed solution
   - Alternatives considered
   - Use cases and examples
   - Expected CLI usage
   - Priority and impact assessment

### Asking Questions

If you have questions about using the CLI:

1. **Check the documentation first**: Review our [User Documentation](docs/README.md) and [README](README.md)
2. **Search existing Q&A**: Look through [closed issues](https://github.com/ibm-hyper-protect/contract-cli/issues?q=is%3Aissue+label%3Aquestion) with the "question" label
3. **Use GitHub Discussions**: For general questions, use [GitHub Discussions](https://github.com/ibm-hyper-protect/contract-cli/discussions)
4. **Create a question issue**: If needed, use our [question template](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=question.yml)

### Code Contributions

We actively welcome your pull requests! However, please follow this process:

1. **Open an issue first** - Before submitting a pull request, open an issue describing:
   - What bug you're fixing or feature you're adding
   - Why it should be fixed or added
   - How you plan to implement it

   This helps us discuss the approach early and avoid duplicated or unnecessary work.

   **Pull requests without a linked issue may be closed.**

2. **Get feedback** - Wait for maintainer feedback on your issue before starting work.

3. **Fork and create a branch** - Once approved, fork the repo and create a feature branch.

4. **Implement your changes** - Follow our coding standards and best practices.

5. **Test thoroughly** - Add tests for your changes and ensure all tests pass.

6. **Submit a pull request** - Reference the original issue in your PR description.

## Getting Started

### Prerequisites

- **Go 1.25.6 or later**
- **OpenSSL** - Required for encryption operations
- **Make** - For running build commands
- **Git** - For version control

### Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR-USERNAME/contract-cli.git
   cd contract-cli
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/ibm-hyper-protect/contract-cli.git
   ```

4. **Install dependencies**:
   ```bash
   make install-deps
   ```

5. **Verify your setup**:
   ```bash
   make test
   ```

## Development Workflow

1. **Create a feature branch** from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our [coding standards](#coding-standards)

3. **View available make targets**:
   ```bash
   make help
   ```

4. **Run tests** frequently during development:
   ```bash
   make test
   # Or simply run the default target
   make
   ```

5. **Build the CLI** to test locally:
   ```bash
   make build
   ```

6. **Tidy dependencies**:
   ```bash
   make tidy
   ```

7. **Commit your changes** with [proper commit messages](#commit-messages)

8. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

9. **Open a Pull Request** from your fork to the main repository

## Coding Standards

### Go Style Guide

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` (run `make fmt`) to format your code (run automatically with most editors)
- Follow [Effective Go](https://golang.org/doc/effective_go) principles
- Write idiomatic Go code

### Best Practices

- **Keep functions small and focused** - Each function should do one thing well
- **Write self-documenting code** - Use clear variable and function names
- **Add comments for complex logic** - Explain the "why", not the "what"
- **Handle errors explicitly** - Never ignore errors
- **Use Cobra best practices** - Follow the patterns established in the codebase
- **Validate user input** - Always validate flags and arguments

### CLI-Specific Guidelines

- **Consistent flag names** - Follow existing naming conventions (e.g., `--in`, `--key`, `--output`)
- **Clear error messages** - Provide helpful, actionable error messages to users
- **Help text** - Write clear descriptions for commands and flags
- **Exit codes** - Use appropriate exit codes (0 for success, non-zero for errors)

### Documentation

- **Add command documentation** - Update [docs/README.md](docs/README.md) for new commands
- **Update main README** - Add usage examples for significant new features
- **Include examples** - Provide practical examples in documentation

## Commit Signing

**All commits to this repository MUST be signed with GPG or SSH keys.** This ensures the authenticity and integrity of code contributions.

### Why Commit Signing?

- **Authentication**: Proves commits are from legitimate team members
- **Integrity**: Detects if commits have been tampered with
- **Non-repudiation**: Authors cannot deny making signed commits
- **Compliance**: Meets security requirements for regulated industries

### Quick Setup

#### Option 1: GPG Signing (Recommended)

**1. Generate GPG Key**
```bash
gpg --full-generate-key
```
- Select: `(1) RSA and RSA`
- Key size: `4096` bits
- Expiration: `1y` (1 year recommended)
- Enter your name and **GitHub email address**

**2. Get Your Key ID**
```bash
gpg --list-secret-keys --keyid-format=long
```
Your key ID is after `rsa4096/` (e.g., `3AA5C34371567BD2`)

**3. Export Public Key**
```bash
gpg --armor --export YOUR_KEY_ID
```
Copy the entire output (including BEGIN and END lines)

**4. Add to GitHub**
- Go to: https://github.com/settings/keys
- Click "New GPG key"
- Paste your public key
- Click "Add GPG key"

**5. Configure Git (Local - This Repo Only)**
```bash
cd contract-go
git config --local commit.gpgsign true
git config --local user.signingkey YOUR_KEY_ID
git config --local user.email "your-github-email@example.com"
```

**6. Fix GPG TTY Issue (macOS/Linux)**
```bash
echo 'export GPG_TTY=$(tty)' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc
```

#### Option 2: SSH Signing (Git 2.34+)

**1. Generate SSH Key** (if you don't have one)
```bash
ssh-keygen -t ed25519 -C "your-github-email@example.com"
```

**2. Add to GitHub**
- Copy your public key: `cat ~/.ssh/id_ed25519.pub`
- Go to: https://github.com/settings/keys
- Click "New SSH key"
- Select "Signing Key" as type
- Paste your public key

**3. Configure Git (Local - This Repo Only)**
```bash
cd contract-go
git config --local gpg.format ssh
git config --local user.signingkey ~/.ssh/id_ed25519.pub
git config --local commit.gpgsign true
```

### Making Signed Commits

Once configured, commits are signed automatically:
```bash
git commit -m "feat: add new feature"
```

Or explicitly sign:
```bash
git commit -S -m "feat: add new feature"
```

### Verify Your Commit is Signed

```bash
# Check last commit
git log --show-signature -1

# Should show:
# gpg: Good signature from "Your Name <your-email@example.com>"
```

### Troubleshooting

#### "gpg failed to sign the data"
```bash
export GPG_TTY=$(tty)
git commit -m "your message"
```

#### Email Mismatch
Ensure your Git email matches your GPG key email:
```bash
git config --local user.email "your-github-email@example.com"
```

#### Sign Previous Commit
```bash
git commit --amend --no-edit -S
git push --force-with-lease
```

### Important Notes

- **Local configuration only**: Signing is enabled only for this repository
- **Other repos unaffected**: Your other repositories won't require signing
- **Pipeline enforcement**: PRs with unsigned commits will fail CI/CD checks
- **Squash and merge**: Only the latest commit needs to be signed

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat:` - A new feature or command
- `fix:` - A bug fix
- `docs:` - Documentation only changes
- `refactor:` - Code changes that neither fix a bug nor add a feature
- `perf:` - Performance improvements
- `test:` - Adding or updating tests
- `chore:` - Changes to build process, dependencies, etc.
- `ci:` - CI/CD configuration changes

### Examples

```
feat(encrypt): add support for IBM Confidential Computing peer pod contracts

fix(base64-tgz): handle empty docker-compose files correctly

docs: update installation instructions for Windows

test(decrypt-attestation): add tests for malformed input
```

### Guidelines

- **Use imperative mood** - "add feature" not "added feature"
- **Keep subject line under 50 characters**
- **Capitalize the subject line**
- **Don't end subject line with a period**
- **Separate subject from body with a blank line**
- **Use body to explain what and why, not how**
- **Reference issues** (e.g., "Fixes #123")

## CI Requirements

Every pull request must pass the following checks before it can be merged:

### 1. Branch Name
Branch names must follow the format `<Type>/<name>`:

| Type | Purpose |
|---|---|
| `Feature` | A new feature |
| `Fix` | A bug fix |
| `Docs` | Documentation only changes |
| `Refactor` | Code changes that neither fix a bug nor add a feature |
| `Performance` | Performance improvements |
| `Test` | Adding or updating tests |
| `Chore` | Changes to build process, dependencies, etc. |
| `CI` | CI/CD configuration changes |

Valid: `Feature/add-encryption`, `Fix/null-pointer`, `Chore/upgrade-deps`
Invalid: `my-branch`, `fix-bug`, `HPVM-123`

### 2. Commit Message
Commits must follow [Conventional Commits](https://www.conventionalcommits.org/) format: `type(optional-scope): description`

Valid: `feat: add contract expiry`, `fix(api): handle null values`
Invalid: `updated stuff`, `WIP`, `bug fix`

### 3. No Force Pushes
Force pushing to a PR branch is not allowed. If you need to update your branch, rebase and push normally:
```bash
git pull --rebase upstream main
git push origin your-branch
```

## Pull Request Process

### Before Submitting

- [ ] Link to the related issue in your PR description
- [ ] Ensure all tests pass (`make test`)
- [ ] Run `make tidy` to clean up dependencies
- [ ] Build the CLI to verify compilation (`make build`)
- [ ] Update documentation if needed
- [ ] Add or update tests for your changes
- [ ] Follow the commit message conventions
- [ ] Rebase your branch on the latest `main` if needed

### PR Template

When opening a PR, include:

```markdown
## Description
Brief description of the changes

## Related Issue
Fixes #issue_number

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] Code refactor

## Testing
Describe the tests you ran and how to reproduce them

## Checklist
- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] I have commented my code where necessary
- [ ] I have updated the documentation
- [ ] I have added tests that prove my fix/feature works
- [ ] All new and existing tests pass
```

### Review Process

1. **Automated checks** - CI must pass before review
2. **Maintainer review** - Tag @Sashwat-K and @vikas-sharma24 as reviewers
3. **Address feedback** - Make requested changes promptly
4. **Approval** - At least one maintainer must approve
5. **Merge** - Maintainers will merge your PR

### After Your PR is Merged

- Delete your feature branch
- Update your local repository:
  ```bash
  git checkout main
  git pull upstream main
  ```

## Testing

### Running Tests

```bash
# Run all tests (default target)
make test
# Or simply
make

# Run tests with coverage
make test-cover

# Run tests with custom coverage
go test -v -race -coverprofile=coverage.out ./...

# Run specific package tests
go test ./cmd -v

# Run specific test
go test ./cmd -v -run TestEncryptCommand
```

### Writing Tests

- **Write tests for all new code** - Aim for good coverage
- **Use table-driven tests** where appropriate
- **Test edge cases** and error conditions
- **Use meaningful test names** - `TestCommandName_Scenario_ExpectedBehavior`
- **Keep tests independent** - Tests should not depend on each other
- **Test CLI output** - Verify command output and exit codes

### Example Test

```go
func TestEncryptCommand(t *testing.T) {
    tests := []struct {
        name      string
        args      []string
        wantErr   bool
        errMsg    string
    }{
        {
            name:    "valid contract",
            args:    []string{"--in", "test.yaml", "--key", "key.pem"},
            wantErr: false,
        },
        {
            name:    "missing contract flag",
            args:    []string{"--key", "key.pem"},
            wantErr: true,
            errMsg:  "required flag \"contract\" not set",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := NewEncryptCommand()
            cmd.SetArgs(tt.args)
            err := cmd.Execute()

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## Questions?

If you have questions about contributing:

1. Check the [documentation](docs/README.md)
2. Search [existing issues](https://github.com/ibm-hyper-protect/contract-cli/issues)
3. Open a new issue with the `question` label
4. Reach out to the maintainers listed in [MAINTAINERS.md](MAINTAINERS.md)

## License

By contributing to `contract-cli`, you agree that your contributions will be licensed under the Apache License 2.0.

---

Thank you for contributing to contract-cli! Your efforts help make this project better for everyone.
