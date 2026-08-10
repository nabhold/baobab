# BAOBAB INTELLIGENCE ENGINE
## Business Case for Board Approval

---

**Document Control**

| Field | Value |
|---|---|
| Document Title | Baobab Intelligence Engine — Business Case |
| Document Type | Board Investment Paper |
| Document Status | Draft for Review |
| Prepared For | NABHOLD Group Africa Board of Directors |
| Prepared By | Office of the Enterprise Software Architect, in consultation with Executive Sponsor |
| Version | 0.1 |
| Date | [To be finalised on submission] |
| Classification | Confidential — Board Use Only |
| Master Format | Markdown (version-controlled). DOCX/PDF renderings are generated derivatives of this source for circulation and are not independently maintained. |
| Related Documents | *Baobab Intelligence Engine* (Product & Programme Proposal, Master Reference); *Baobab Platform — AI-Powered Market Research* (Companion Reference); Baobab Intelligence Engine — Product Requirements Document *(to follow)*; Baobab Intelligence Engine — Solution Architecture *(to follow)* |

---

## Executive Summary

NABHOLD Group Africa operates in a continent-wide information environment that is fragmented, inconsistent, and poorly suited to the pace at which commercial and investment decisions must be made. Market information exists — in government datasets, statistical agencies, financial filings, trade records, and research literature — but it is scattered, unequally reliable, and rarely organised around the questions that investors and businesses actually ask.

This paper proposes that NABHOLD approve, in principle, the development of the **Baobab Intelligence Engine**: an evidence-grounded, AI-assisted research and opportunity-discovery capability that will become part of the broader Baobab Platform. Its first commercial expression will be the production of high-quality African market and investment research reports. Its enduring value, however, is not the reports themselves — it is the accumulated, structured, evidence-backed intelligence asset that underlies them, which strengthens with every research engagement and becomes progressively more valuable to NABHOLD and to paying clients over time.

This is **not** a request to fund a seven-phase build. It is a request to:

1. Approve the strategic direction and long-term roadmap in principle;
2. Release an initial funding envelope for **Phase 0 (Foundation) and Phase 1 (Research MVP)** only;
3. Approve a **stage-gated funding model**, under which every subsequent phase must earn its own approval by demonstrating defined, measurable outcomes — including the Board's right to terminate the programme at any gate; and
4. Approve the governance structure through which this programme will be sponsored, reviewed, and escalated.

The financial case at this stage is presented as a **structured model with placeholder assumptions**, not as a costed budget. NABHOLD Finance and Executive Management are asked to validate the numerical assumptions before Gate 0 approval is finalised; this paper deliberately avoids presenting false precision on cost or revenue before the underlying research proposition has been proven.

**The decision requested of the Board** is set out in full in Section 12, and is, in summary: approve the strategic direction; approve the Gate 0–1 funding envelope; approve the governance and stage-gate model; and mandate NABHOLD Finance to validate the financial assumptions ahead of full Gate 1 commitment.

---

## 1. Purpose of This Document

This Business Case exists to secure a Board decision, not to describe a product vision. The strategic vision, functional scope, and long-term architecture of the Baobab Intelligence Engine are already documented in full in the master reference document, *Baobab Intelligence Engine* (and its companion, *Baobab Platform — AI-Powered Market Research*). Those documents remain the authoritative source for what the system is intended to become.

This paper translates that vision into the form the Board must act on: a problem statement, a strategic rationale, options considered, a phased investment ask, a financial structure, a risk position, and a governance model with explicit decision gates.

Three further documents will follow this Business Case, each derived from it and from the master reference documents, and each answering a distinct question:

| Document | Question It Answers |
|---|---|
| Business Case *(this document)* | Should NABHOLD build this? |
| Product Requirements Document | What must the product do? |
| Solution Architecture | How should we build it? |
| Programme / Project Plan | How will we deliver it? |

This separation is deliberate. Collapsing these into a single document tends to produce four versions of the same conversation rather than four documents that each do their job.

---

## 2. Strategic Context and Problem Statement

Businesses, investors, and institutions operating in Africa face a structural information problem. Economic information relevant to commercial and investment decisions exists — but it is:

- Fragmented across jurisdictions, agencies, and formats;
- Inconsistent in quality and update frequency;
- Difficult to compare across countries and sectors;
- Frequently disconnected from the specific commercial or investment decision a user needs to make.

The practical consequence is that market research in Africa is slow, expensive, difficult to reproduce, and difficult to compare consistently across opportunities. Investors and businesses either commission bespoke research at high cost and long lead times, or make decisions on incomplete information.

This is the problem the Baobab Intelligence Engine is designed to close — not by generating more content faster, but by building a structured intelligence capability that converts fragmented information into a disciplined chain:

> **Evidence → Knowledge → Analysis → Opportunity Intelligence → Decision Support.**

The strategic premise, carried directly from the master reference document, is explicit and should be treated as a design constraint on everything that follows:

> **Baobab should not be built as an AI report generator. It should be built as an intelligence system from which reports are one output.**

---

## 3. Strategic Objective and Fit Within NABHOLD

The Baobab Intelligence Engine is proposed as a **strategic capability of NABHOLD Group Africa**, not a standalone software product. It sits within the broader Baobab Platform and is intended, over time, to become the intelligence foundation on which both NABHOLD's own investment decisions and NABHOLD's commercial research products are built.

Two outcome streams justify the investment:

**Internal outcomes** — NABHOLD gains a systematic capability to research markets faster, identify and compare investment opportunities consistently, support due diligence, monitor markets continuously, and build proprietary institutional knowledge that compounds over time rather than being re-derived on every engagement.

**Commercial outcomes** — the same underlying capability produces a family of saleable products: market intelligence briefs, sector and country reports, bespoke investment research, opportunity dashboards, and — over the longer term — enterprise intelligence subscriptions and machine-readable data/API products.

These two streams are not competing uses of the same investment; they are the same asset expressed twice. This is the strategic rationale for building Baobab as a platform capability rather than a single-purpose reporting tool, and it is the basis of the long-term compounding effect — described in the master reference document as the **Commercial Intelligence Flywheel** — in which NABHOLD's own investment activity, informed by Baobab, generates operational data and outcomes that in turn improve the intelligence engine itself.

---

## 4. Proposed Solution — Overview

At full maturity, the Baobab Intelligence Engine is conceived as a layered system: a Data layer, a Research layer, and a Knowledge layer feeding an Analytics layer (market, risk, and opportunity analysis), which in turn drives an Intelligence layer producing reports, dashboards, and — eventually — APIs. Underpinning all of this is a discipline in which AI is used as a reasoning and automation layer, while evidence and structured data remain the authoritative source of truth. This principle governs the entire system and should be treated as non-negotiable regardless of which AI models or techniques are used at any point in the build.

This Business Case does **not** ask the Board to approve the construction of that full system in one step. It asks the Board to approve:

- The **strategic direction and seven-phase roadmap in principle** (summarised in Section 6 and detailed in Appendix A), as the long-term context within which near-term investment decisions are made; and
- **Funding for Phase 0 (Foundation) and Phase 1 (Research MVP) only** — a single, complete, narrow vertical slice proving the entire intelligence chain end-to-end: a research question is received, sources are discovered, evidence is retrieved and extracted, analysis is performed, findings are synthesised with citations, and a research output is produced.

This vertical-slice approach is deliberate. It proves the hardest and most commercially material assumption — that Baobab can produce evidence-backed, traceable, decision-useful African market research — before capital is committed to the opportunity engine, monitoring, enterprise platform, or API layers described in later phases.

A guiding architectural principle has been agreed for how Phase 0–1 is to be built, and is recorded here because it directly affects cost and risk: **build narrowly, architect intelligently.** The Solution Architecture (to follow this Business Case) will implement only what Phase 0–1 requires, but will define, as explicit and documented extension points, the seams through which the Opportunity Engine, Knowledge Graph, Scoring, Monitoring, Reporting, Client Portal, and API layers will later attach. The Architecture document will state plainly, for every major component, whether it is *implemented*, *stubbed with a defined interface*, or *deferred design* — so that no future team mistakes a placeholder for a commitment, and no Phase 1 decision silently forecloses a Phase 3 capability.

---

## 5. Options Considered

The Board should be aware of the alternatives considered and why the proposed approach was preferred.

| Option | Description | Assessment |
|---|---|---|
| **A — Do nothing / continue with ad hoc research** | Continue commissioning bespoke research per engagement, as at present. | Does not build a proprietary asset; cost and lead time do not improve over time; no commercial product emerges. Rejected as a strategic dead end. |
| **B — Buy / license an existing market-intelligence platform** | Procure a third-party research or business-intelligence platform and adapt it to African markets. | No existing commercial platform is purpose-built for African evidence sourcing, source-authority tiering, or the opportunity-scoring model NABHOLD requires. Would create vendor dependency on NABHOLD's core intelligence asset — the one component the master reference document identifies as the long-term strategic value driver. Rejected as strategically unsound, though selected third-party components (data feeds, retrieval infrastructure, LLM providers) remain appropriate as building blocks within Option D. |
| **C — Build the full seven-phase system immediately ("big bang")** | Commit capital and resourcing to the entire scope described in the master reference document from the outset. | Highest cost and risk exposure; commits capital before the core research proposition is validated; directly contravenes the programme's own guiding principle of building the intelligence foundation before the commercial surface. Rejected. |
| **D — Phased, stage-gated build, starting with a narrow vertical slice (Phase 0–1)** | Prove the research-to-evidence-to-output chain first; expand only on demonstrated success; architect for future phases without building them prematurely. | **Recommended.** Lowest initial capital exposure; earliest possible validation of the core value proposition; preserves full strategic optionality; aligns with the programme's own stated development philosophy. |

---

## 6. Scope and Phasing

### 6.1 Two Horizons

This programme deliberately operates on two distinct horizons, and the Board's approval decision differs for each:

```
                   BAOBAB VISION
                         │
                         ▼
                 SEVEN-PHASE ROADMAP  ← approved in principle, as strategic context
                         │
              ┌──────────┴──────────┐
              │                     │
         CURRENT RELEASE       FUTURE CAPABILITIES
              │
              ▼
           PHASE 0 – 1  ← funding requested now
              │
              ▼
       INITIAL VERTICAL SLICE
```

**Strategic roadmap (Phases 0–7):** approved in principle only, as the context within which near-term decisions are made and against which the Architecture's extension points are designed. Full detail is at Appendix A.

**Funded scope (Phases 0–1):** the only scope for which capital release is requested at this time.

### 6.2 What Phase 0–1 Delivers

- **Phase 0 — Foundation:** product requirements, domain model, research methodology, intelligence taxonomy, evidence model, initial source strategy, technical architecture, and evaluation framework.
- **Phase 1 — Research MVP:** a working, end-to-end demonstration of research question → research plan → source discovery → evidence retrieval and extraction → analysis → citation → research output.

### 6.3 What Phase 0–1 Explicitly Does Not Deliver

To keep the funding ask honest, the Board should note that Phase 0–1 does **not** include: the Opportunity Discovery Engine and Opportunity Score, the Knowledge Graph, risk and scenario modelling, continuous monitoring and alerts, the client portal, multi-tenant enterprise features, or any API/data-product layer. These remain Phase 2 and beyond, each subject to its own gate.

---

## 7. Commercial Model and Benefits Case

### 7.1 Internal Benefits (NABHOLD)

- Faster, more systematic market research to support NABHOLD's own investment activity.
- A consistent, comparable basis for evaluating opportunities across geographies and sectors.
- Reduced repetitive research cost over time as institutional knowledge accumulates.
- A defensible, evidence-traceable basis for investment committee decisions.

### 7.2 Commercial Benefits (External)

- A saleable family of research products — market intelligence briefs, sector and country reports, bespoke investment research — with a credible path to recurring subscription and enterprise revenue as the platform matures.
- A defensible commercial position built on proprietary evidence and methodology rather than repackaged third-party content.

### 7.3 The Nature of the Asset Being Built

The Board should approve this investment on the understanding that **the reports are the first commercial expression of the asset, not the asset itself.** The underlying value being built is the accumulated evidence, structured knowledge, relationships, analytical models, and institutional experience that make each subsequent research engagement faster, cheaper, and more reliable than the last. This distinction should inform how the Board evaluates progress at each gate: report output volume alone is an insufficient measure of programme success (see Section 11).

---

## 8. Financial Case

### 8.1 Approach

Consistent with the decision frozen ahead of this paper, **no budget figures are invented at this stage.** Presenting numerical precision before Phase 0 has scoped the actual cost drivers — particularly AI/compute consumption, which is inherently usage-dependent — would give the Board a false basis for decision-making. Instead, this section presents the **financial structure** the Business Case will be costed against, with all numerical assumptions explicitly marked for validation by NABHOLD Finance and Executive Management ahead of Gate 0/Gate 1 sign-off. Full detail is at Appendix C.

### 8.2 Cost Structure

```
CAPEX                              OPEX
├── Product / Software Development ├── Cloud infrastructure
├── AI / Compute (setup)           ├── AI model consumption
├── Data acquisition               ├── Data subscriptions
├── Infrastructure                 ├── Personnel
├── Security                       ├── Maintenance
├── Research                       ├── Security
└── Professional services          └── Support
```

### 8.3 Revenue Structure

```
COMMERCIAL REVENUE
├── Research reports
├── Bespoke research
├── Subscriptions
├── Enterprise licences
├── Intelligence services
└── API / data products
```

### 8.4 Scenarios

Three scenarios will be modelled once Finance validates cost and adoption assumptions:

- **Conservative** — low adoption, slower commercialisation.
- **Base** — expected commercial adoption.
- **Upside** — strong adoption and enterprise expansion.

Each scenario will ultimately produce: initial investment, operating cost, revenue, gross margin, break-even point, payback period, ROI, and (where appropriate) NPV/IRR, with sensitivity analysis on the key adoption and AI-consumption-cost assumptions.

### 8.5 Financial KPI to Track From Day One

One operating metric is specified now, ahead of full financial modelling, because it is the single most Board-relevant cost signal in an AI-driven research system:

> **Cost per research request / research unit** — the fully loaded cost (AI/compute, data, and marginal infrastructure) of producing one unit of research output at defined quality, tracked from the first Phase 1 evaluation runs onward.

This KPI should be reported at every stage gate from Gate 1 forward, as it is the leading indicator of whether the commercial model in Section 7.2 is economically viable at scale.

---

## 9. Governance, Stage-Gating, and Decision Rights

### 9.1 Governance Chain

```
Business Analyst
      │
      ▼
Enterprise Software Architect
      │
      ▼
Project Manager / Programme Manager
      │
      ▼
Executive Sponsor
      │
      ▼
Investment / Management Committee
      │
      ▼
NABHOLD Group Africa Board of Directors
      │
      ▼
Funding Gate
```

The **NABHOLD Group Africa Board of Directors** holds final investment authority. The Executive Sponsor and the Investment/Management Committee provide recommendation, technical and commercial governance, and gate-review functions, but do not themselves authorise capital release beyond the envelope the Board has approved.

### 9.2 Stage-Gated Funding

Funding is **not** released as a single envelope. The Board is asked to approve an initial envelope covering Gate 0 and Gate 1 only. Each subsequent gate must be earned on the basis of defined, demonstrated outcomes before further capital is released.

| Gate | Decision |
|---|---|
| Gate 0 | Approve concept and discovery |
| Gate 1 | Approve MVP development |
| Gate 2 | Approve commercial pilot |
| Gate 3 | Approve production scale-up |
| Gate 4 | Approve enterprise expansion |

**Gate 1 is the critical gate for this programme.** No major scale-up of investment should occur until the Phase 0–1 vertical slice has demonstrated, against the quantified criteria in Section 11, that the core research proposition actually works.

### 9.3 Gates Can Terminate the Programme

This must be stated explicitly, and is recorded here as a formal governance principle rather than an implied option: **a stage gate is a decision point with three possible outcomes — proceed, proceed with conditions, or stop.** The Board and the Investment/Management Committee retain the right, at any gate, to discontinue the programme if the evidence does not support continuation. This is the specific mechanism by which the programme avoids the common technology-investment failure mode of continuing to fund a system because money has already been spent on it. Each phase earns the right to proceed to the next; none is entitled to it.

---

## 10. Architecture and Technology Approach

This section summarises, for the Board's awareness, the technical governance approach that will structure the Solution Architecture document to follow. It is included here because it directly affects the risk and cost profile of the investment being requested.

**Classification.** The Baobab Intelligence Engine is a **greenfield capability within the broader Baobab Platform.** It is neither an isolated new NABHOLD system nor a bolt-on module to an unrelated existing application; it shares a platform context with other Baobab enterprise capabilities (planned or existing) and will be architected accordingly.

**Technology direction.** The intended Baobab technology ecosystem — Python, Django, PostgreSQL, Redis, Celery, FastAPI for specialised AI services, containerised deployment, GitHub-based CI/CD, and AWS as the eventual primary cloud environment — is adopted as the **architectural context**, not as an immutable constraint. The Solution Architecture will explicitly distinguish confirmed architectural decisions from proposed technologies subject to architecture review, and every material technology choice or deviation (vector store, graph store, agent-orchestration framework, and similar) will be recorded as a numbered Architecture Decision Record from the first day of Phase 0. This is a deliberate governance discipline to prevent an early convenience choice from hardening into unreviewed architectural dogma.

**AI governance principle.** AI is treated throughout as a reasoning and automation layer. It is not the source of truth. The source of truth is the evidence and structured data layer beneath it. This principle is carried directly from the master reference document and will be enforced as a design constraint in the Solution Architecture and, subsequently, in engineering practice.

**Documentation standard.** The Solution Architecture will be documented using **arc42** for structure and **C4** (Context → Container → Component) for visualisation — a combination chosen because it separates enterprise-level architecture governance (which remains within NABHOLD's existing TOGAF-aligned governance where applicable) from the product-level engineering documentation this programme actually needs day to day.

---

## 11. Risks and Mitigations

| Risk | Description | Mitigation |
|---|---|---|
| **AI hallucination** | AI-generated analysis may present unsupported claims convincingly. | Evidence-grounded generation; mandatory citation; independent evidence validation; confidence scoring separate from opportunity/quality scoring. |
| **Poor source quality** | Research conclusions may rest on low-authority sources. | Source authority tiering (Tier 1–4) enforced at ingestion; provenance tracking. |
| **Outdated information** | Market conditions in Africa can shift faster than research is refreshed. | Freshness metadata on all evidence; defined update schedules; monitoring introduced from Phase 5 onward. |
| **False precision** | A numeric opportunity or quality score can imply certainty the evidence does not support. | Confidence is always reported separately from any score; assumptions and evidence are shown alongside conclusions, never in place of them. |
| **Over-automation** | AI may produce fluent but poorly reasoned conclusions without adequate challenge. | Structured research workflow; mandatory human review states for commercial output; contrarian/counter-thesis analysis built into the method from Phase 3 onward. |
| **Data licensing exposure** | Commercial use of third-party information may carry licensing restrictions. | Source licensing controls; provenance tracking; data-rights review built into the source-authority framework. |
| **Excessive technical complexity** | The platform could become overbuilt before commercial value is proven. | Stage-gated, vertical-slice delivery (Section 6); "build narrowly, architect intelligently" as a binding architectural principle; no Phase 2+ component is built ahead of its gate. |
| **Cost overrun on AI/compute** | LLM and retrieval consumption costs are usage-dependent and can scale unpredictably with research volume. | Cost-per-research-unit tracked as a financial KPI from Gate 1 onward (Section 8.5); scenario modelling with sensitivity analysis on this specific variable. |
| **Programme continuation bias** | Sunk-cost pressure to continue funding despite weak Gate 1 evidence. | Explicit Board-retained right to terminate at any gate (Section 9.3); Gate 1 exit criteria are quantified and pre-agreed (Section 11 below / Appendix E), not assessed retrospectively. |

---

## 12. Success Criteria and Phase 1 Exit Criteria

### 12.1 Defining Success Correctly

Phase 1 success is **not** defined as "an AI research application was built." That is a technology deliverable, not a business outcome. The exit criterion for Gate 1 is defined as:

> **Can Baobab reliably conduct a defined class of African market research and produce evidence-backed, traceable, decision-useful intelligence at commercially credible quality?**

### 12.2 Quantified Phase 1 Exit Criteria

The evaluation dimensions identified for Phase 1 are converted below into measurable, testable thresholds, so that Gate 1 is assessed against pre-agreed evidence rather than subjective impression. These thresholds are proposed as a starting baseline for the Investment/Management Committee to ratify; the full evaluation methodology will be detailed in the Product Requirements Document and executed as part of Phase 0's evaluation framework.

| Dimension | Question | Proposed Quantified Threshold |
|---|---|---|
| Research quality | Did it answer the research question? | ≥85% of test research questions (from an agreed evaluation set) receive a "materially answered" rating from independent analyst review. |
| Evidence | Are claims supported? | ≥90% of factual claims in sampled outputs resolve to at least one cited evidence record. |
| Source quality | Are sources credible? | ≥80% of cited evidence in sampled outputs derives from Tier 1–3 sources. |
| Citation | Can claims be traced? | 100% of factual claims presented in a research output carry a resolvable citation to source, publisher, and date. |
| Accuracy | Are factual statements correct? | ≥95% of sampled factual statements verified as accurate on independent spot-check. |
| Completeness | Did it cover material dimensions? | ≥80% of the pre-defined material dimensions for a given research question class are addressed in the output. |
| Freshness | Is the information sufficiently current? | ≥90% of cited evidence meets the freshness threshold defined for its evidence type. |
| Reasoning | Are conclusions logically supported? | ≥85% of sampled conclusions assessed by independent analyst review as logically derived from the cited evidence. |
| Reproducibility | Can the research be repeated? | Re-running an identical research question produces materially consistent findings in ≥90% of test cases. |
| Cost | Is the research economically viable? | Cost per research unit (Section 8.5) is measured and reported; no fixed threshold is set until Gate 0 baseline data exists, but the trend must be flat or improving across the Phase 1 evaluation period. |
| Time | Is it materially faster than conventional research? | Demonstrates at least a 5x reduction in elapsed time versus a comparable conventional (manual analyst) research baseline. |
| Human review | Can analysts validate the output efficiently? | Median analyst review-and-validation time per output does not exceed [to be baselined in Phase 0] hours. |

Full evaluation methodology, sampling approach, and the analyst review protocol are detailed in Appendix E, and will be finalised during Phase 0 as part of the evaluation-framework deliverable.

---

## 13. Recommendation and Decision Requested

It is recommended that the NABHOLD Group Africa Board of Directors:

1. **Approve, in principle**, the strategic direction and seven-phase roadmap of the Baobab Intelligence Engine as set out in this paper and in the master reference documents (Appendix A), as the long-term context for the programme — without committing capital beyond Phase 0–1 at this time.
2. **Approve the release of an initial funding envelope** covering **Gate 0 (concept and discovery) and Gate 1 (MVP development)** only, on the financial structure set out in Section 8 and Appendix C, subject to numerical validation by NABHOLD Finance and Executive Management.
3. **Approve the stage-gated governance and funding model** set out in Section 9, including the explicit right of the Board and the Investment/Management Committee to discontinue the programme at any gate.
4. **Approve the governance chain** set out in Section 9.1, with the Executive Sponsor and Investment/Management Committee holding recommendation and gate-review authority, and the Board retaining final investment authority.
5. **Ratify the proposed Phase 1 exit criteria** in Section 12.2 (or refer them to the Investment/Management Committee for refinement) as the pre-agreed basis on which Gate 1 will be assessed.
6. **Mandate NABHOLD Finance** to validate the cost, revenue, and adoption assumptions underlying the financial model in Appendix C ahead of final Gate 1 sign-off.

Subject to Board approval of the above, the next documents in the chain — the Product Requirements Document and the Solution Architecture — will be developed against the frozen scope, governance, and success criteria set out in this paper, and will in turn form the basis of the Programme/Project Plan.

---

## Appendix A — Seven-Phase Strategic Roadmap (Context Only — Not Funded by This Paper)

The following phases are drawn directly from the master reference document and are presented here as strategic context. Only Phase 0–1 is the subject of the current funding request; Phases 2–7 will each require their own gate approval as set out in Section 9.

| Phase | Name | Establishes |
|---|---|---|
| 0 | Foundation | Product requirements, domain model, research methodology, intelligence taxonomy, evidence model, initial source strategy, technical architecture, evaluation framework |
| 1 | Research MVP | Research questions, source discovery, retrieval, evidence extraction, basic synthesis, citation management |
| 2 | Intelligence Foundation | Country model, sector model, market model, company model, evidence graph, historical data, knowledge relationships |
| 3 | Opportunity Engine | Opportunity discovery, opportunity objects, opportunity scoring, confidence scoring, risk assessment, investment thesis, counter-thesis |
| 4 | Commercial Research Products | Market reports, opportunity reports, country reports, executive briefs, report generation, client delivery |
| 5 | Continuous Intelligence | Monitoring, alerts, opportunity score changes, market updates, continuous research |
| 6 | Enterprise Intelligence | Enterprise dashboards, multi-tenant intelligence, advanced permissions, custom research environments, enterprise reporting |
| 7 | Intelligence Infrastructure | APIs, data products, external integrations, machine-readable intelligence, partner ecosystem |

---

## Appendix B — Stage-Gate Definitions (Detail)

*[To be expanded with entry criteria, exit criteria, evidence required, and reviewing body for each gate, as ratified by the Investment/Management Committee. Gate 0 and Gate 1 definitions should be finalised before Board submission; Gates 2–4 may be defined in outline only at this stage.]*

---

## Appendix C — Financial Model (Detail)

*[To be completed in consultation with NABHOLD Finance. Structure: CAPEX/OPEX/Revenue as per Section 8.2–8.3; Conservative/Base/Upside scenarios as per Section 8.4; cost-per-research-unit baseline methodology as per Section 8.5. No figures are to be entered without Finance sign-off.]*

---

## Appendix D — Risk Register (Detail)

*[To be expanded from Section 11 into a full register with likelihood, impact, owner, and review cadence for each risk, maintained as a living document from Phase 0 onward.]*

---

## Appendix E — Phase 1 Evaluation Framework (Detail)

*[To be developed in Phase 0 as a formal evaluation-framework deliverable, expanding Section 12.2 into a full methodology: evaluation question set, sampling approach, analyst review protocol, and scoring rubric. This appendix will become a direct input to the Product Requirements Document's acceptance criteria.]*

---

## Appendix F — Traceability Matrix

| Business Case Section | Source Reference |
|---|---|
| §2 Strategic Context | *Baobab Intelligence Engine* §1–2 |
| §3 Strategic Objective & Fit | *Baobab Intelligence Engine* §4, §42, §60 |
| §4 Proposed Solution | *Baobab Intelligence Engine* §7, §43, §52, §55 |
| §6 Scope & Phasing | *Baobab Intelligence Engine* §54–55 |
| §7 Commercial Model | *Baobab Intelligence Engine* §5, §39; *Baobab Platform* Commercial Products / UX sections |
| §9 Governance & Stage-Gating | Foundational Decisions (frozen), this document chain |
| §10 Architecture Approach | *Baobab Intelligence Engine* §43–48, §52 |
| §11 Risks | *Baobab Intelligence Engine* §58 |
| §12 Success Criteria | *Baobab Intelligence Engine* §50, §56–57; Foundational Decisions §6 |

---

## Appendix G — Glossary

| Term | Meaning |
|---|---|
| Evidence | A discrete, sourced, dated, and provenance-tracked information unit supporting a research claim. |
| Source Authority Tier | Classification (Tier 1–4) of a source's reliability, from primary official data to secondary/open-web material. |
| Opportunity Score | A weighted, configurable scoring model for comparing commercial/investment opportunities (Phase 3+; not in Phase 0–1 scope). |
| Confidence Score | A measure of evidential reliability, always reported separately from any opportunity or quality score. |
| Vertical Slice | A narrow but complete implementation of the full research chain, end to end, used to validate the core proposition before broader build-out. |
| Stage Gate | A formal governance decision point at which the programme may proceed, proceed with conditions, or be discontinued. |
| ADR | Architecture Decision Record — a numbered, dated record of a material architectural or technology decision and its rationale. |

---

*End of Business Case. This document is the frozen basis for the Product Requirements Document and Solution Architecture to follow.*
