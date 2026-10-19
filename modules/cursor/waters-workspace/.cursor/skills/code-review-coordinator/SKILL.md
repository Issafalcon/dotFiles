---
name: code-review-coordinator
description: Coordinates targeted pull request reviews across WCC repositories — loads the PR from GitHub, traces cross-repo impact through graphify graphs, gathers Jira issue and design documentation context via the Atlassian API, then dispatches focused review subagents and consolidates their findings into one review. Use when reviewing a GitHub pull request in the wcc workspace.
disable-model-invocation: true
---

# Code review coordinator

Coordinate a review of the pull request in $ARGUMENTS.

You are the **coordinator**. You gather context, dispatch focused review subagents, and
consolidate their output. You do **not** write the review findings yourself and you never
modify code.

## Step 1 — Load the pull request

```bash
gh pr view "$PR_URL" --json title,body,author,url,number,headRefName,baseRefName,files,additions,deletions,labels,statusCheckRollup,reviewDecision
gh pr diff "$PR_URL"
```

Record the repository as `<owner>/<repo>` and map it to the local checkout at `<repo>/` in
the workspace root. If that directory is missing, say so — the cross-repo step and any local
verification will be limited to what `gh` returns.

All paths below are relative to the workspace root (the directory holding the repos and this
`.cursor/` directory), so this skill works on any machine. Run the shell steps from there.

Note the scope as `N files changed (+X/-Y)`; the output contract requires it.

## Step 2 — Classify the change

Classification decides which profiles in [review-profiles.md](review-profiles.md) apply.
More than one may apply. Judge from the changed paths and the diff, not the PR title alone.

| Signal in the diff | Profile |
|---|---|
| Auth, authorisation, tenancy, shared contracts, SDK or OpenAPI surface | Architecture-sensitive |
| Secrets, crypto, deserialisation, input validation, dependency bumps | Security |
| Hotfix, revert, or incident fix against a release branch | Critical production fix |
| `infrastructure/`, `deployment/`, `monitoring/`, `.github/workflows/`, `Dockerfile` | Infrastructure |
| Wide rename, move or restructure; high file count with near-neutral net lines | Large refactoring |
| None of the above | Routine |

Announce the classification, the profiles selected, and the model each will run on **before**
dispatching. If a selected profile routes to an unavailable model, follow the availability
rules in [review-profiles.md](review-profiles.md) and ask the user before proceeding.

## Step 3 — Trace cross-repo impact

A WCC change rarely stops at its own repo. For every changed symbol or contract that is
visible outside the service — public types, SDK surface, message schemas, OpenAPI operations,
NuGet package surface, shared Terraform modules — look for consumers in the **other** repos:

```bash
for g in */graphify-out/graph.json; do
  echo "--- ${g%%/*}"
  graphify explain "<symbol>" --graph "$g"
done
```

Also available, all accepting `--graph`:

- `graphify query "<question>" --graph <g>` — BFS traversal for a broader question
- `graphify path "<A>" "<B>" --graph <g>` — dependency chain between two symbols
- `graphify merge-graphs <g1> <g2> --out /tmp/wcc-merged-graph.json` — one cross-repo graph to
  query when a chain spans repos

When `explain` reports an ambiguous match, retry with the repo-relative path or the full node
id it prints.

**Report graph coverage honestly.** Not every repo has a graph. List the repos that had no
`graphify-out/graph.json`, and state that they were *not analysed* rather than implying no
impact. Offer `graphify update <repo>` to build one.

Carry forward an impact list: consuming repo, symbol, how it is consumed, and the risk if the
PR's behaviour changes.

## Step 4 — Gather Jira and design context

Find the issue key in the PR title, body, or branch name (e.g. `INFSAAS-20417`). If no key is
obvious, **ask the user** for it — do not guess or proceed without it.

Load the issue by following `~/.agents/skills/jira-context/SKILL.md`, which handles
connectivity and falls back to asking the user when Atlassian credentials are absent.

Then walk outward for design intent and epic-level context:

```bash
AUTH="$ATLASSIAN_EMAIL:$ATLASSIAN_API_TOKEN"
BASE="https://waterscorporation.atlassian.net"

# Parent / epic and linked issues
curl -sf -u "$AUTH" "$BASE/rest/api/3/issue/$ISSUE_KEY?fields=parent,issuelinks,summary,description,labels" | jq .

# Remote links — Confluence pages, design docs, external specs
curl -sf -u "$AUTH" "$BASE/rest/api/3/issue/$ISSUE_KEY/remotelink" | jq -r '.[].object | "\(.title) \(.url)"'

# Epic body, once the parent key is known
curl -sf -u "$AUTH" "$BASE/rest/api/3/issue/$PARENT_KEY?expand=renderedFields" | jq .

# Confluence page content, for a page id taken from a remote link
curl -sf -u "$AUTH" "$BASE/wiki/api/v2/pages/$PAGE_ID?body-format=storage" | jq -r '.body.storage.value'
```

Also check the repo's own design specs (`docs/design/`, `docs/architecture/`, `Docs/design/`)
for a spec whose Upstream metadata matches the issue key.

Summarise the acceptance criteria and design constraints the PR must satisfy. A review that
does not know the intended behaviour cannot judge correctness.

## Step 5 — Dispatch review subagents

Split the work by **area of concern**, not by file. Typical areas: domain and use-case
correctness; API contract and backward compatibility; tenancy and authorisation; tests;
infrastructure and deployment; observability. Only dispatch areas the diff actually touches.

Launch them in **parallel** — one message containing multiple Task calls — using
`subagent_type: generalPurpose`, with the model from the selected profile. Use the
`security-review` subagent only when the user explicitly asks for a security review.

Every subagent prompt must contain all six of these:

1. The PR URL, repo, and the `gh pr view` / `gh pr diff` commands to fetch the diff itself
2. The **Focus** list for its profile, verbatim from [review-profiles.md](review-profiles.md)
3. The cross-repo impact list from Step 3
4. The Jira acceptance criteria and design constraints from Step 4
5. `Read .cursor/rules/code-review-output.mdc from the workspace root and follow it exactly.`
   Pass the absolute path you resolved for it, so the subagent can read it from any directory.
6. Its area boundary, and: report findings only, change no code, open no PR comments

## Step 6 — Challenger pattern

For platform-wide changes, authentication, tenant isolation, shared SDKs, and critical
production systems, run the challenger pattern in [review-profiles.md](review-profiles.md)
**sequentially** instead of the single parallel pass. The challenger receives the PR plus the
two prior reviews so it can attack their conclusions.

## Step 7 — Consolidate

Merge every subagent response into **one** review that follows
`.cursor/rules/code-review-output.mdc`.

- Each issue appears exactly once, in Findings, at its highest reported severity.
- Where subagents disagree, resolve it and say which view won and why. Do not list both.
- Keep confirmed defects separate from unresolved risks.
- Attribute cross-repo findings to the consuming repo, and name the graph evidence.
- State which models actually ran, and name any repo that had no graph coverage.

Post to GitHub only if the user asks. When asked, write the body to a temp file **outside**
any repo so nothing lands in git:

```bash
tmp="$(mktemp /tmp/pr-review.XXXXXX.md)"
printf '%s\n' "$REVIEW_MD" > "$tmp"
gh pr comment "$PR_URL" --body-file "$tmp"
rm -f "$tmp"
```
