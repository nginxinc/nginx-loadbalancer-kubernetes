# Act events

This project uses [act](https://github.com/nektos/act) to help test GitHub Actions locally. This directory contains act event definitions.

## Requirements

- [act](https://github.com/nektos/act), this guide assumes you have installed act as a [GitHub CLI extension](https://github.com/nektos/act#installation-as-github-cli-extension).
- [Docker](https://www.docker.com/)

## Running act

The following commands can be run to exercise the GitHub Actions locally.

Will pass the `main` conditional and run the `build-and-sign-image` job.

```bash
gh act  -W .github/workflows/build-and-sign-image.yml -e .github/workflows/.act/push-event-valid.json
```

Will fail the `main` conditional and not run the `build-and-sign-image` job.

```bash
gh act  -W .github/workflows/build-and-sign-image.yml -e .github/workflows/.act/push-event-invalid.json
```
