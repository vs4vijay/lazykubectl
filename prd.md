# Product Requirements Document

**Product Name:** lazykubectl (Rust Edition)
**Version:** 1.0 (MVP)
**Date:** 2026-05-07
**Author:** Product Team
**Status:** Draft

---

## 1. Executive Summary

lazykubectl is a Terminal User Interface (TUI) client for Kubernetes that provides an intuitive, keyboard-driven alternative to the traditional `kubectl` command-line tool. This document outlines requirements for a complete rewrite of the existing Go-based lazykubectl in Rust, focusing on improved reliability, type safety, and maintainability while preserving the core user experience.

**Strategic Rationale:**
- **Performance**: Rust's zero-cost abstractions and memory efficiency for resource-constrained environments
- **Reliability**: Strong type system eliminates entire classes of runtime errors present in the Go version
- **Modern Foundation**: Leverage Rust's growing ecosystem (ratatui, kube-rs) for long-term sustainability
- **Developer Experience**: Better error messages, compile-time guarantees, and cleaner architecture

---

## 2. Problem Statement

### Current Challenges

**For DevOps Engineers and Platform Teams:**
- `kubectl` requires memorizing complex commands and flags
- Navigating between namespaces, pods, and logs requires multiple command invocations
- No visual overview of cluster state without third-party tools
- Switching context between resources is cumbersome
- Log viewing requires piping to external tools (less, grep)

**Limitations of Existing Go Implementation:**
- Silent error failures in rendering code
- Runtime panics due to nil pointer dereferences
- String-based state management prone to inconsistencies
- Monolithic TUI file difficult to maintain (357 lines)
- Limited test coverage

### Opportunity

A Rust rewrite addresses these technical limitations while maintaining the intuitive UX that makes lazykubectl valuable for daily Kubernetes operations.

---

## 3. Goals and Objectives

### Primary Goals

1. **Feature Parity (MVP)**: Deliver core functionality matching the Go version
   - Namespace browsing and deletion
   - Pod and container exploration
   - Log viewing (last 100 lines)
   - Real-time event monitoring

2. **Improved Reliability**: Zero runtime panics, comprehensive error handling

3. **Cross-Platform Support**: Windows, Linux, macOS compatibility

4. **Developer Experience**: Clean, testable codebase with clear module boundaries

### Success Metrics

- **Reliability**: Zero crashes during 100-hour continuous operation test
- **Performance**: Startup time < 500ms, UI responsiveness < 100ms latency
- **Adoption**: 100+ GitHub stars in first 3 months post-release
- **Quality**: 70%+ code coverage, all integration tests passing
- **Compatibility**: Works with Kubernetes 1.24-1.30

---

## 4. Target Users

### Primary Personas

**1. DevOps Engineer (Sarah)**
- Manages 5-10 Kubernetes clusters daily
- Needs quick cluster health checks and log troubleshooting
- Prefers keyboard-driven workflows over GUI tools
- Values speed and efficiency

**2. Platform Engineer (Marcus)**
- Develops and maintains internal Kubernetes platforms
- Frequently debugs deployment issues across namespaces
- Requires reliable tooling for production incident response
- Needs cross-platform support (Linux workstation, macOS laptop)

**3. SRE Team Member (Priya)**
- On-call rotation requires rapid cluster diagnostics
- Needs to quickly navigate from namespace → pod → container → logs
- Values terminal-based tools that work over SSH
- Requires stability and predictable behavior

### Secondary Personas

- Junior developers learning Kubernetes
- CI/CD pipeline operators monitoring deployments
- Security teams auditing cluster configurations

---

## 5. Use Cases

### UC-1: Troubleshoot Application Logs
**Actor:** DevOps Engineer
**Precondition:** User has kubectl configured with cluster access
**Flow:**
1. Launch `lazykubectl`
2. Navigate to problematic namespace (arrow keys + Enter)
3. Select failing pod from list
4. Select container if multi-container pod
5. View last 100 lines of logs to identify error

**Postcondition:** User identifies root cause from log output
**Priority:** P0 (Critical)

### UC-2: Monitor Cluster Events
**Actor:** SRE Team Member
**Precondition:** Cluster experiencing issues
**Flow:**
1. Launch `lazykubectl`
2. Event panel (bottom) automatically streams cluster events
3. Observe warnings, errors, and state changes in real-time

**Postcondition:** User detects event anomalies (e.g., CrashLoopBackOff, ImagePullBackOff)
**Priority:** P0 (Critical)

### UC-3: Clean Up Test Namespaces
**Actor:** Platform Engineer
**Precondition:** Multiple test namespaces exist
**Flow:**
1. Launch `lazykubectl`
2. Navigate to test namespace
3. Press 'd' to delete namespace
4. Confirm deletion

**Postcondition:** Namespace and all resources deleted
**Priority:** P1 (High)

### UC-4: Explore Cluster Resources
**Actor:** Junior Developer
**Precondition:** New to Kubernetes cluster
**Flow:**
1. Launch `lazykubectl`
2. Browse available namespaces
3. Select namespace to view pods
4. Navigate through pods to understand application structure

**Postcondition:** User gains understanding of cluster organization
**Priority:** P2 (Medium)

### UC-5: Quick Cluster Health Check
**Actor:** DevOps Engineer
**Precondition:** Morning routine check
**Flow:**
1. Launch `lazykubectl`
2. Scan namespace list for unexpected entries
3. View event log for errors/warnings
4. Verify critical pods are running

**Postcondition:** User confirms cluster health or identifies issues
**Priority:** P1 (High)

---

## 6. Requirements

### 6.1 Functional Requirements (MVP)

#### FR-1: Kubernetes Connectivity
- **FR-1.1**: Load kubeconfig from default location (`~/.kube/config`)
- **FR-1.2**: Support custom kubeconfig path via `--kubeconfig` flag
- **FR-1.3**: Use current context from kubeconfig
- **FR-1.4**: Validate cluster connectivity on startup
- **FR-1.5**: Display clear error message if connection fails

#### FR-2: Namespace Management
- **FR-2.1**: List all namespaces in cluster
- **FR-2.2**: Display namespace metadata (name, status, age)
- **FR-2.3**: Select namespace via keyboard navigation (up/down arrows + Enter)
- **FR-2.4**: Delete namespace via 'd' key
- **FR-2.5**: Refresh namespace list via 'r' key

#### FR-3: Pod Browsing
- **FR-3.1**: List pods in selected namespace
- **FR-3.2**: Display pod status (Running, Pending, Failed, etc.)
- **FR-3.3**: Select pod via keyboard navigation
- **FR-3.4**: Automatically load pods when namespace is selected
- **FR-3.5**: Handle empty namespaces gracefully

#### FR-4: Container Navigation
- **FR-4.1**: Extract containers from selected pod
- **FR-4.2**: Display container names in list
- **FR-4.3**: Select container via keyboard navigation
- **FR-4.4**: Support multi-container pods
- **FR-4.5**: Handle init containers

#### FR-5: Log Viewing
- **FR-5.1**: Display last 100 lines of container logs
- **FR-5.2**: Support scrolling through logs (up/down arrows)
- **FR-5.3**: Automatically load logs when container is selected
- **FR-5.4**: Handle containers with no logs
- **FR-5.5**: Display timestamps in log output

#### FR-6: Event Monitoring
- **FR-6.1**: Stream cluster events in real-time
- **FR-6.2**: Display events in dedicated panel (bottom)
- **FR-6.3**: Show event type, reason, and message
- **FR-6.4**: Auto-scroll to newest events
- **FR-6.5**: Buffer last 100 events

#### FR-7: User Interface
- **FR-7.1**: Display cluster info (context, cluster name, user) in top panel
- **FR-7.2**: Highlight selected item in lists
- **FR-7.3**: Show current view context (e.g., "Pods in namespace: default")
- **FR-7.4**: Support mouse click selection (optional enhancement)
- **FR-7.5**: Render in terminal with proper color support

#### FR-8: Keyboard Shortcuts
- **FR-8.1**: `↑/↓` - Navigate items in current list
- **FR-8.2**: `Enter` - Select item and drill down
- **FR-8.3**: `Tab` - Cycle between panels (future)
- **FR-8.4**: `d` - Delete namespace (only in namespace view)
- **FR-8.5**: `r` - Refresh current view
- **FR-8.6**: `Ctrl+C` - Quit application
- **FR-8.7**: `Esc` - Navigate back to previous view (future)

### 6.2 Non-Functional Requirements

#### NFR-1: Performance
- **NFR-1.1**: Application startup < 500ms on modern hardware
- **NFR-1.2**: UI response to keyboard input < 100ms
- **NFR-1.3**: Kubernetes API calls timeout after 30 seconds
- **NFR-1.4**: Support clusters with 100+ namespaces without lag
- **NFR-1.5**: Memory usage < 50MB for typical operations

#### NFR-2: Reliability
- **NFR-2.1**: Zero panics during normal operation
- **NFR-2.2**: Gracefully handle Kubernetes API errors
- **NFR-2.3**: Recover from transient network failures
- **NFR-2.4**: Display user-friendly error messages for all error conditions
- **NFR-2.5**: Never crash due to malformed Kubernetes responses

#### NFR-3: Compatibility
- **NFR-3.1**: Support Kubernetes 1.24+ (current and previous 6 minor versions)
- **NFR-3.2**: Work with all authentication methods (token, cert, OIDC, GCP, AWS, Azure)
- **NFR-3.3**: Compatible with Windows 10+, Linux (kernel 4.0+), macOS 11+
- **NFR-3.4**: Support common terminal emulators (Windows Terminal, iTerm2, Alacritty, etc.)
- **NFR-3.5**: Handle terminal resize gracefully

#### NFR-4: Maintainability
- **NFR-4.1**: Modular architecture with clear separation of concerns
- **NFR-4.2**: Comprehensive inline documentation
- **NFR-4.3**: 70%+ code coverage
- **NFR-4.4**: Pass Clippy lints (Rust static analysis)
- **NFR-4.5**: Follow Rust API guidelines

#### NFR-5: Security
- **NFR-5.1**: Never log sensitive credentials
- **NFR-5.2**: Respect kubeconfig file permissions
- **NFR-5.3**: Use official kube-rs library (no custom auth)
- **NFR-5.4**: No network calls except to Kubernetes API
- **NFR-5.5**: Sanitize log output to prevent terminal injection

#### NFR-6: Usability
- **NFR-6.1**: Clear visual hierarchy in UI layout
- **NFR-6.2**: Consistent color scheme across views
- **NFR-6.3**: Display help hints in status bar
- **NFR-6.4**: Provide `--help` flag with usage examples
- **NFR-6.5**: Support `--version` flag

---

## 7. User Experience

### 7.1 UI Layout

```
┌─────────────────────────────────────────────────────────────┐
│ Cluster Info Panel (20% height)                             │
│ Context: prod-us-west | Cluster: gke-prod | User: admin    │
│ Nodes: 12 | Namespaces: 45                                  │
├─────────────────────────────────────────────────────────────┤
│ Main View Panel (50% height)                                │
│                                                              │
│ [Dynamic content based on navigation state]                 │
│ - Namespaces list (initial view)                            │
│ - Pods list (after namespace selection)                     │
│ - Containers list (after pod selection)                     │
│ - Logs (after container selection)                          │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│ Events Panel (30% height)                                   │
│ [Real-time cluster events streaming]                        │
│ Normal  | ScalingReplicaSet | Scaled up replica set...     │
│ Warning | BackOff            | Back-off restarting...       │
└─────────────────────────────────────────────────────────────┘
```

### 7.2 Navigation Flow

```
Start
  │
  ├─> Namespaces View
  │     │
  │     ├─ Select namespace → Pods View
  │     │                        │
  │     │                        ├─ Select pod → Containers View
  │     │                        │                  │
  │     │                        │                  └─ Select container → Logs View
  │     │                        │
  │     │                        └─ Press 'r' → Refresh pods
  │     │
  │     └─ Press 'd' → Delete namespace
  │
  └─> Events continuously update in bottom panel
```

### 7.3 Interaction Patterns

**Selection Pattern:**
1. Highlight moves with ↑/↓ keys
2. Visual indicator (green background) shows current selection
3. Enter key confirms selection and transitions to next view

**Error Display Pattern:**
1. Error appears in status bar at top of main panel
2. Red color indicates error state
3. Error message provides actionable guidance
4. User can dismiss by pressing any key

**Loading Pattern:**
1. "Loading..." indicator shows during API calls
2. Spinner animation (optional)
3. Timeout after 30 seconds with retry option

---

## 8. Technical Architecture

### 8.1 Technology Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|---------|-----------|
| Language | Rust | 1.75+ | Memory safety, performance, strong typing |
| TUI Framework | ratatui | 0.26+ | Most active Rust TUI library, immediate mode rendering |
| Terminal Backend | crossterm | 0.27+ | Cross-platform terminal manipulation |
| Kubernetes Client | kube-rs | 0.88+ | Official Rust Kubernetes client |
| Async Runtime | tokio | 1.x | De facto standard for async Rust |
| CLI Parser | clap | 4.x | Popular, ergonomic CLI framework |
| Error Handling | thiserror + anyhow | 1.x | Standard error handling pattern |

### 8.2 System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         main.rs                              │
│ - Parse CLI arguments (clap)                                 │
│ - Initialize logging                                         │
│ - Bootstrap application                                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                         app.rs                               │
│ - Create KubeClient                                          │
│ - Create TUI App                                             │
│ - Spawn async tasks (event watcher, log streamer)          │
│ - Main event loop coordination                              │
└──────────┬──────────────────────────┬───────────────────────┘
           │                          │
           ▼                          ▼
┌────────────────────────┐   ┌───────────────────────────────┐
│   k8s/client.rs        │   │   tui/app.rs                  │
│ - KubeClient struct    │   │ - AppState management         │
│ - get_namespaces()     │   │ - Event handling              │
│ - get_pods()           │   │ - View rendering              │
│ - get_logs()           │   │ - Navigation logic            │
│ - delete_namespace()   │   │ - Keyboard input processing   │
│ - watch_events()       │   │                               │
└────────────────────────┘   └───────────────────────────────┘
```

### 8.3 State Management

**Type-Safe State Machine:**

```rust
pub enum ViewMode {
    Namespaces,
    Pods { namespace: String },
    Containers { namespace: String, pod: String },
    Logs { namespace: String, pod: String, container: String },
}

pub struct AppState {
    pub view_mode: ViewMode,
    pub namespaces: Vec<Namespace>,
    pub pods: Vec<Pod>,
    pub containers: Vec<Container>,
    pub logs: Vec<String>,
    pub events: Vec<String>,
    pub selected_index: usize,
    pub info: ClusterInfo,
}
```

**Benefits:**
- Impossible to have invalid state (e.g., viewing logs without selecting container)
- Compiler enforces valid state transitions
- Clear data ownership and lifecycle

### 8.4 Async Communication

**Message-Passing Pattern:**

```rust
pub enum AppEvent {
    NamespaceSelected(String),
    PodSelected(String),
    ContainerSelected(String),
    LogsReceived(Vec<String>),
    EventReceived(String),
    Refresh,
    Error(String),
    Quit,
}
```

Background tasks (event watcher) → `tokio::mpsc::channel` → Main TUI loop

### 8.5 Error Handling Strategy

**Custom Error Types:**

```rust
#[derive(Error, Debug)]
pub enum LazyKubectlError {
    #[error("Kubernetes API error: {0}")]
    KubeError(#[from] kube::Error),

    #[error("Failed to connect to cluster")]
    NoClusterConnectivity,

    #[error("Configuration error: {0}")]
    ConfigError(String),

    #[error("TUI error: {0}")]
    TuiError(#[from] std::io::Error),
}
```

**Philosophy:**
- Fail fast with clear messages
- Propagate errors with `?` operator
- Display errors in TUI without crashing
- Log errors for debugging

### 8.6 Module Structure

```
src/
├── main.rs              # CLI entry point
├── app.rs               # Application orchestration
├── error.rs             # Error types
├── config.rs            # Configuration handling
├── k8s/
│   ├── mod.rs
│   ├── client.rs        # Kubernetes API client
│   ├── types.rs         # Type aliases
│   └── watch.rs         # Watch/stream utilities
├── tui/
│   ├── mod.rs
│   ├── app.rs           # TUI state machine
│   ├── ui.rs            # Rendering logic
│   ├── event.rs         # Input handling
│   ├── state.rs         # State types
│   └── widgets/
│       ├── mod.rs
│       ├── info.rs      # Info panel
│       ├── namespaces.rs
│       ├── pods.rs
│       └── logs.rs
└── utils/
    └── helpers.rs       # Utility functions
```

---

## 9. Success Metrics

### 9.1 Launch Criteria (MVP)

- ✅ All FR (Functional Requirements) implemented
- ✅ All P0 use cases working
- ✅ Zero known panics or crashes
- ✅ Works on Windows, Linux, macOS
- ✅ Integration tests pass against Kubernetes 1.28+
- ✅ Documentation complete (README, --help)
- ✅ Binary builds available for all platforms

### 9.2 Quality Metrics

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Code Coverage | 70%+ | cargo-tllvm-cov |
| Build Success | 100% | GitHub Actions CI |
| Clippy Warnings | 0 | cargo clippy --all-targets |
| Performance (Startup) | < 500ms | Benchmark suite |
| Performance (UI Response) | < 100ms | Manual testing |
| Memory Usage | < 50MB | valgrind / activity monitor |

### 9.3 User Adoption Metrics (Post-Launch)

- **Week 1**: 50+ downloads
- **Month 1**: 100+ GitHub stars
- **Month 3**: 5+ community contributions
- **Month 6**: 1000+ weekly downloads

---

## 10. Out of Scope (MVP)

The following features are **explicitly excluded** from the MVP and deferred to future releases:

### Excluded Features

- **Services View**: Service listing and details (present in Go version)
- **Deployments View**: Deployment status and management
- **Node Details**: Detailed node information beyond count
- **ConfigMaps/Secrets**: Resource viewing and editing
- **Shell Execution**: `kubectl exec` equivalent
- **Port Forwarding**: Interactive port-forward management
- **Resource Editing**: Inline YAML editing
- **Multi-Cluster**: Switching between multiple clusters in-app
- **Search/Filter**: Resource filtering by name or label
- **Log Streaming**: Real-time log following (vs. tail -100)
- **Resource Metrics**: CPU/Memory usage graphs
- **Custom Themes**: User-configurable color schemes
- **Configuration File**: `~/.lazykubectl.yaml` support
- **Plugin System**: Extensibility via plugins
- **CRD Support**: Custom Resource Definition browsing

### Rationale

These features add significant complexity and are not essential for the core value proposition (quick namespace/pod/log navigation). They can be prioritized based on user feedback post-MVP.

---

## 11. Release Plan

### Phase 1: Foundation (Weeks 1-2)
**Deliverable:** CLI tool with K8s connectivity

- Initialize Cargo project
- Implement error handling (`error.rs`)
- Implement Kubernetes client (`k8s/client.rs`)
  - Connection establishment
  - `get_namespaces()`
  - `get_pods()`
- CLI argument parsing (`config.rs`)
- Basic integration test

### Phase 2: Basic TUI (Weeks 3-4)
**Deliverable:** Interactive namespace and pod browser

- Implement state management (`tui/state.rs`)
- Implement UI rendering (`tui/ui.rs`)
- Build widgets (info, namespaces, pods)
- Keyboard event handling (`tui/event.rs`)
- Navigation: namespace → pod

### Phase 3: Complete Navigation (Week 5)
**Deliverable:** Full drill-down to logs

- Container extraction from pod spec
- Log fetching (tail 100 lines)
- Logs widget
- Navigation: pod → container → logs
- Arrow key navigation

### Phase 4: Actions and Events (Week 6)
**Deliverable:** Full MVP feature set

- Namespace deletion
- Event watching background task
- Event streaming to UI
- Keyboard shortcuts ('d', 'r')
- Error display in UI

### Phase 5: Polish and Release (Weeks 7-8)
**Deliverable:** Production-ready v0.1.0

- Error message improvements
- Cross-platform testing
- Documentation (README, help text)
- Integration test suite
- CI/CD pipeline (GitHub Actions)
- Binary releases
- Performance optimization

### Post-MVP Roadmap

**v0.2.0** (Q3 2026)
- Service viewing
- Deployment status
- Search/filter functionality
- Log follow mode

**v0.3.0** (Q4 2026)
- Resource metrics (CPU/Memory)
- ConfigMap/Secret viewing
- Multi-cluster support
- Custom themes

**v1.0.0** (Q1 2027)
- Full feature parity with Go version
- Plugin system
- Configuration file support
- Comprehensive testing and documentation

---

## 12. Dependencies and Risks

### 12.1 External Dependencies

| Dependency | Risk Level | Mitigation |
|------------|-----------|------------|
| kube-rs API changes | Medium | Pin to specific version, monitor changelog |
| ratatui breaking changes | Low | Stable API, large community |
| Kubernetes API deprecations | Low | Support wide version range (1.24+) |
| Terminal compatibility | Medium | Extensive cross-platform testing |
| Crossterm platform bugs | Low | Well-tested library, fallback options |

### 12.2 Technical Risks

**Risk: Async complexity in TUI**
- **Impact:** High - Could block entire development
- **Likelihood:** Medium
- **Mitigation:** Use tokio::select!, prototype early, reference existing projects (gitui)

**Risk: Performance with large clusters**
- **Impact:** Medium - Poor UX for enterprise users
- **Likelihood:** Medium
- **Mitigation:** Implement pagination, lazy loading, caching (post-MVP)

**Risk: Cross-platform terminal rendering**
- **Impact:** Medium - Windows users affected
- **Likelihood:** Low
- **Mitigation:** Crossterm handles this, test on Windows early

**Risk: kube-rs authentication edge cases**
- **Impact:** High - Users can't connect to clusters
- **Likelihood:** Low
- **Mitigation:** Use kube-rs defaults, extensive integration testing

### 12.3 Resource Constraints

- **Development Time**: 6-8 weeks (1 developer)
- **Testing Infrastructure**: Requires local Kubernetes cluster (minikube/kind)
- **CI/CD**: GitHub Actions free tier sufficient for builds

---

## 13. Acceptance Criteria

### MVP Acceptance Checklist

**Functionality:**
- [ ] Connect to Kubernetes cluster via kubeconfig
- [ ] List all namespaces
- [ ] Navigate to namespace and view pods
- [ ] Navigate to pod and view containers
- [ ] View container logs (last 100 lines)
- [ ] Delete namespace with 'd' key
- [ ] Refresh view with 'r' key
- [ ] Stream cluster events in bottom panel
- [ ] Quit with Ctrl+C

**Quality:**
- [ ] Zero panics during 1-hour continuous use test
- [ ] All unit tests passing
- [ ] Integration tests against K8s 1.28, 1.29, 1.30
- [ ] Clippy passes with zero warnings
- [ ] Code coverage ≥70%
- [ ] README with installation and usage instructions
- [ ] --help flag provides clear usage guide

**Platforms:**
- [ ] Builds successfully on Windows (MSVC)
- [ ] Builds successfully on Linux (GNU)
- [ ] Builds successfully on macOS (Apple Silicon and Intel)
- [ ] Manual testing completed on all three platforms
- [ ] Terminal resize handled gracefully on all platforms

**Performance:**
- [ ] Startup time < 500ms (measured on reference hardware)
- [ ] UI responds to input < 100ms
- [ ] No memory leaks during 12-hour stress test

**Documentation:**
- [ ] README.md with quick start guide
- [ ] Architecture documentation (this PRD + code comments)
- [ ] Inline documentation for all public APIs
- [ ] Release notes for v0.1.0

---

## 14. Appendix

### A. Comparison with Go Version

| Aspect | Go Version | Rust Version |
|--------|-----------|--------------|
| **Language** | Go 1.14+ | Rust 1.75+ |
| **TUI Library** | gocui (retained mode) | ratatui (immediate mode) |
| **K8s Client** | client-go | kube-rs |
| **CLI Framework** | Cobra | clap |
| **Concurrency** | Goroutines | Tokio async/await |
| **State Management** | String map | Type-safe enum |
| **Error Handling** | Often ignored | Enforced with Result |
| **Type Safety** | Runtime checks | Compile-time guarantees |
| **Memory Safety** | Garbage collected | Borrow checker |
| **Binary Size** | ~20MB | ~15MB (estimated) |

### B. Key Technical Decisions

**Decision 1: Single-threaded async runtime**
- **Choice:** `tokio::main(flavor = "current_thread")`
- **Rationale:** I/O-bound workload, simpler state management

**Decision 2: Enum-based state machine**
- **Choice:** `ViewMode` enum with associated data
- **Rationale:** Impossible to represent invalid states, better than string maps

**Decision 3: Immediate mode UI (ratatui)**
- **Choice:** Redraw entire UI each frame
- **Rationale:** Simpler than retained mode, no state synchronization issues

**Decision 4: Message-passing concurrency**
- **Choice:** Tokio channels for async task communication
- **Rationale:** Clear separation, prevents data races, idiomatic Rust

**Decision 5: Fixed 3-panel layout**
- **Choice:** Info (top), Main (middle), Events (bottom)
- **Rationale:** Matches 90% use case, simpler than complex grid system

### C. References

- **Original Go Implementation:** `S:/GitHub/lazykubectl`
- **kube-rs Documentation:** https://docs.rs/kube
- **ratatui Documentation:** https://ratatui.rs
- **Rust API Guidelines:** https://rust-lang.github.io/api-guidelines/

### D. Glossary

- **TUI**: Terminal User Interface - text-based graphical interface in terminal
- **kubeconfig**: Kubernetes configuration file containing cluster credentials
- **Namespace**: Kubernetes resource for logical cluster partitioning
- **Pod**: Smallest deployable unit in Kubernetes (one or more containers)
- **Container**: Isolated runtime environment for application
- **Watch API**: Kubernetes API for streaming resource changes
- **Immediate Mode**: UI rendering pattern where entire UI is redrawn each frame
- **Message-Passing**: Concurrency pattern using channels to communicate between tasks

---

**Document Version:** 1.0
**Last Updated:** 2026-05-07
**Next Review:** Upon MVP completion
