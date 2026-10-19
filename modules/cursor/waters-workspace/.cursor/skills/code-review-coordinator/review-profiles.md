# Review profiles

Model, effort and focus per review type. The coordinator selects one or more profiles in
Step 2 and passes the matching **Focus** list verbatim into each subagent prompt.

## Model availability

Cursor subagents accept only the slugs listed in `<available_subagent_models>`. Map the
intended model to a slug before calling the Task tool:

| Intended model | Slug to pass | Note |
|---|---|---|
| Claude Sonnet 5 | `claude-sonnet-5-5-high` | Nearest available Sonnet |
| Claude Opus 5 | `claude-opus-5-5-medium` | Nearest available Opus |
| GPT-5.6 Sol | `gpt-5.6-sol-medium` | Only effort tier available |
| GPT-5.6 Terra | **unavailable** | No Terra slug exists in this environment |
| Fable | **unavailable** | Also excluded from security-sensitive reviews |

Effort is encoded in the slug, so the requested effort levels below cannot always be set
independently. Treat them as intent, and state in the final review which model actually ran.

**When a profile routes to Terra** (Large refactoring, Infrastructure), stop and ask the user
whether to run it on `gpt-5.6-sol-medium` — the escalation target those profiles already name —
or to wait for Terra. Do not substitute silently.

---

## Routine pull requests

Model: Claude Sonnet 5
Effort: Medium, Thinking on

Focus on:

- Correctness
- Maintainability and readability
- Test coverage
- Error handling
- Logging and observability
- Backward compatibility

## Large refactoring reviews

Model: GPT-5.6 Terra
Effort: High

Escalate to Sol when the change crosses service or security boundaries.

Focus on:

- Behavioural regressions
- Broken dependency chains
- Missing updates
- Runtime and build risks
- Removed safeguards
- Validation needed to prove behaviour is unchanged

Trace affected paths from entry points to side effects.

## Architecture-sensitive reviews

Model: GPT-5.6 Sol
Effort: High

Use for authentication, authorisation, tenancy, security, reliability, scalability, and
shared contracts.

Evaluate:

- Trust and authorisation boundaries
- Tenant isolation
- Data ownership
- Concurrency and consistency
- Scalability and reliability
- Backward compatibility
- Operational impact

For each risk, explain why it matters, its likely production impact, and mitigation options.

## Critical production fixes

Model: GPT-5.6 Sol
Effort: High or Max

Challenge:

- Original incident recurrence
- Silent failure
- Data corruption
- Race conditions
- Partial failure
- Rollback failure
- Monitoring gaps

Conclude with: Safe to deploy / Deploy with mitigations / Do not deploy.

## Security reviews

Primary model: GPT-5.6 Sol
Independent challenger: Claude Opus 5
Effort: High or Max

Review for:

- Privilege escalation
- Authentication and authorisation bypass
- Secret disclosure
- Injection
- Insecure defaults
- Tenant isolation failure
- Unsafe deserialisation
- Supply-chain and dependency risk

For each issue, provide the attack scenario, impacted assets, severity, and mitigation.

Avoid Fable for security-sensitive code unless Anthropic retention has been explicitly
approved for the material.

## Infrastructure reviews

Model: GPT-5.6 Terra
Effort: High

Escalate to Sol for production-wide, identity, network, or data-loss risk.

Focus on:

- Deployment and rollback safety
- Destructive resource changes
- IAM permissions
- Networking
- State migration
- Monitoring
- Cost impact
- Production availability

---

## Challenger review pattern

Use for platform-wide changes, authentication, tenant isolation, shared SDKs, and critical
production systems.

Workflow:

1. Routine review with Sonnet 5
2. Risk review with GPT-5.6 Sol
3. Independent challenge review with Claude Opus 5

Run these **in order**, not in parallel — the challenger receives the PR plus the first two
reviews so it can attack their conclusions.

Challenger prompt:

```
Assume this pull request contains at least one serious flaw.
Find the strongest evidence for and against merging it.

Challenge:
- assumptions
- design decisions
- scalability expectations
- operational behaviour
- security boundaries
- rollback and recovery

Identify the most likely source of a future production incident.
Do not invent findings: distinguish confirmed defects from unresolved risks.
```
