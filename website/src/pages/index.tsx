import React from 'react';
import Layout from '@theme/Layout';
import Head from '@docusaurus/Head';
import Link from '@docusaurus/Link';
import TerminalDemo from '@site/src/components/TerminalDemo';
import s from './index.module.css';

/* ── Hero ──────────────────────────────────────────────────────────── */

function Hero(): React.JSX.Element {
  return (
    <header className={s.hero}>
      <div className={s.heroGlow} aria-hidden="true" />
      <div className={s.heroGrain} aria-hidden="true" />
      <div className={s.heroInner}>
        <div className={s.heroTag}>
          <span className={s.heroTagDot} aria-hidden="true" />
          Open Source &mdash; Written in Go
        </div>
        <h1 className={s.heroTitle}>
          The compiler<br />
          <em className={s.heroTitleItalic}>for mortals.</em>
        </h1>
        <p className={s.heroSubtitle}>
          Not another Obsidian alternative. LazyMD is an LLVM-inspired
          markdown compiler for philosophers, product engineers, founders,
          lawyers, researchers &mdash; anyone turning raw thought into
          structured output. Plain text in, anything out.
        </p>
        <div className={s.heroActions}>
          <Link className={s.btnPrimary} to="/docs/getting-started/installation">
            Get Started
          </Link>
          <Link className={s.btnSecondary} to="/docs/getting-started/quick-start">
            Documentation
          </Link>
        </div>
      </div>
    </header>
  );
}

/* ── Social Proof ─────────────────────────────────────────────────── */

const badges: Array<[string, string]> = [
  ['\u26A1', 'Written in Go'],
  ['\u{1F517}', 'MCP Protocol'],
  ['\u2328', 'Vim Keybindings'],
  ['\u{1F4E6}', 'Single Binary'],
  ['\u{1F9E0}', 'Knowledge Graph'],
];

function SocialProof(): React.JSX.Element {
  return (
    <section className={s.socialProof}>
      <div className={s.socialProofInner}>
        {badges.map(([icon, label], i) => (
          <React.Fragment key={label}>
            {i > 0 && <span className={s.divider} aria-hidden="true" />}
            <span className={s.badge}>
              <span className={s.badgeIcon} aria-hidden="true">{icon}</span>
              {label}
            </span>
          </React.Fragment>
        ))}
      </div>
    </section>
  );
}

/* ── Philosophy: Why Plain Text ───────────────────────────────────── */

const philosophy = [
  {
    icon: '\u2726',
    title: 'LLVM for prose',
    desc: 'LazyMD follows the LLVM architecture: a universal frontend parses markdown into an intermediate representation. Backends compile it to any target \u2014 HTML, PDF, slides, structured data, or AI-consumable context.',
  },
  {
    icon: '\u2696',
    title: 'Git-native by default',
    desc: 'Every line is diffable. Every change is mergeable. No binary blobs, no sync conflicts. Your notes and code live in the same version control as everything else.',
  },
  {
    icon: '\u221E',
    title: 'Zero lock-in, infinite portability',
    desc: '.md files work in any editor, on any OS, in any decade. Your content outlives every app. Markdown written in 2004 reads perfectly today \u2014 and will in 2044.',
  },
  {
    icon: '\u2318',
    title: 'Human-readable, machine-parseable',
    desc: 'Plain text is the only format that both humans and machines read natively. No schema, no runtime. Open a file and it just makes sense \u2014 to you and to every AI agent.',
  },
];

function Philosophy(): React.JSX.Element {
  return (
    <section className={s.section} id="philosophy">
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.sectionLabel}>Philosophy</span>
        </div>
        <h2 className={s.sectionTitle}>
          Parse once,<br />compile anywhere.
        </h2>
        <p className={s.sectionDesc}>
          Like LLVM separated frontends from backends, LazyMD separates
          writing from output. One IR, infinite targets. Plain text is the
          universal source code of human thought.
        </p>
        <div className={s.philosophyGrid}>
          {philosophy.map(({icon, title, desc}) => (
            <article key={title} className={s.philosophyCard}>
              <span className={s.philosophyIcon} aria-hidden="true">{icon}</span>
              <h3>{title}</h3>
              <p>{desc}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ── Interfaces: Runs Everywhere ──────────────────────────────────── */

const interfaces = [
  {
    label: 'Terminal (TUI)',
    status: 'Available now',
    live: true,
    desc: 'The original. SSH into a server, open a tmux pane, or work locally. Vim keybindings, live preview, zero dependencies. Feels like home.',
    ascii: '  \u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510\n  \u2502 $ lm note.md \u2502\n  \u2502 \u2588\u2588\u2588\u2588\u2588 \u2502 # Hello   \u2502\n  \u2502 docs  \u2502 world.    \u2502\n  \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518',
  },
  {
    label: 'Web Editor',
    status: 'Coming soon',
    live: false,
    desc: 'The same editing experience in your browser. Real-time collaboration, cloud sync, and zero install. Open a link and start writing.',
    ascii: '  \u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510\n  \u2502 \u25CF\u25CF\u25CF  lazymd.com   \u2502\n  \u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524\n  \u2502 # Hello world.  \u2502\n  \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518',
  },
  {
    label: 'Native Apps',
    status: 'Coming soon',
    live: false,
    desc: 'macOS, Windows, Linux, iOS, Android. Native performance with platform-native feel. One markdown vault, every device in your life.',
    ascii: '  \u250C\u2500\u2500\u2500\u2500\u2510  \u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2510\n  \u2502 \u2588\u2588 \u2502  \u2502  \u2588\u2588  \u2502\n  \u2502 \u2588\u2588 \u2502  \u2502  \u2588\u2588  \u2502\n  \u2514\u2500\u2500\u2500\u2500\u2518  \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2518\n   iOS      macOS',
  },
  {
    label: 'Any Device',
    status: 'The vision',
    live: false,
    desc: 'Rabbit R1, smart displays, e-ink tablets, embedded systems. If it has a screen and a network connection, LazyMD can run there.',
    ascii: '  \u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510\n  \u2502 \u25B7 R1   \u2502\n  \u2502 ~~~~   \u2502\n  \u2502 ~~~~   \u2502\n  \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518\n   Any surface.',
  },
];

function Interfaces(): React.JSX.Element {
  return (
    <section className={s.sectionAlt} id="interfaces">
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.sectionLabel}>Interfaces</span>
        </div>
        <h2 className={s.sectionTitle}>
          One compiler.<br />Every surface.
        </h2>
        <p className={s.sectionDesc}>
          Your terminal is just the first backend. LazyMD is designed as a
          protocol-first compiler that targets any device with a screen.
        </p>
        <div className={s.interfaceGrid}>
          {interfaces.map(({label, status, live, desc, ascii}) => (
            <article key={label} className={s.interfaceCard}>
              <pre className={s.interfaceAscii} aria-hidden="true">{ascii}</pre>
              <div className={s.interfaceMeta}>
                <h3>{label}</h3>
                <span className={live ? s.statusLive : s.statusSoon}>{status}</span>
              </div>
              <p>{desc}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ── Features ─────────────────────────────────────────────────────── */

const features = [
  {icon: '\u2328', title: 'Vim-Native Editing', desc: 'Full modal editing with Normal, Insert, and Command modes. Navigate with hjkl, motions with w/b, delete with dd, undo with u \u2014 muscle-memory compatible.', wide: false},
  {icon: '\u25CE', title: 'Live Preview', desc: 'Rendered markdown in a side panel. Headers, bold, italic, code blocks with syntax highlighting \u2014 all updating as you type.', wide: false},
  {icon: '\u2588', title: 'Multi-Panel Layout', desc: 'Inspired by lazygit \u2014 file tree, editor, preview, and brain graph side by side. Toggle panels with Alt+1/2/3.', wide: false},
  {icon: '\u2726', title: 'Syntax Highlighting', desc: 'Built-in highlighting for Go, Python, JavaScript, TypeScript, Rust, Zig, C, C++, Java, and 16+ languages. Theme-aware colors with a pluggable highlighter backend.', wide: true},
  {icon: '\u2699', title: 'Plugin System', desc: 'Register custom commands, hook into editor events, add panels. Build and share community plugins.', wide: false},
  {icon: '\u2192', title: 'Zero Dependencies', desc: 'Pure Go with Bubble Tea and Lip Gloss. No runtime dependencies. Fast startup, tiny single binary.', wide: false},
  {icon: '\u2387', title: 'Mouse Support', desc: 'Click to position cursor, scroll with mouse wheel, click panels to focus. Works in iTerm2, Alacritty, kitty, and more.', wide: false},
  {icon: '\u2B21', title: 'MCP Server', desc: 'Built-in Model Context Protocol server with 22 tools. AI agents like Claude Code and Gemini CLI connect via JSON-RPC 2.0 over stdio.', wide: false},
  {icon: '\u{1F9E0}', title: 'Brain: Knowledge Graph', desc: 'Obsidian-style graph view for [[wiki-links]]. Visualize connections between notes with a force-directed ASCII layout. Navigate, explore backlinks, and find orphan notes.', wide: true},
  {icon: '\u26A1', title: 'Instant Startup', desc: 'Compiles to a single ~2MB binary. Launches in milliseconds, even on large files. No JVM, no Electron, no wait.', wide: false},
];

function Features(): React.JSX.Element {
  return (
    <section className={s.section} id="features">
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.sectionLabel}>Features</span>
        </div>
        <h2 className={s.sectionTitle}>Everything you need,<br />nothing you don't.</h2>
        <p className={s.sectionDesc}>
          A compiler built for mortals. Whether you're drafting briefs, shipping products, or writing research &mdash; LazyMD stays out of your way.
        </p>
        <div className={s.featureGrid}>
          {features.map(({icon, title, desc, wide}) => (
            <article key={title} className={`${s.featureCard} ${wide ? s.featureCardWide : ''}`}>
              <div className={s.featureIcon} aria-hidden="true">{icon}</div>
              <h3>{title}</h3>
              <p>{desc}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ── MCP / AI Integration ─────────────────────────────────────────── */

const mcpTools = [
  'open_file', 'read_document', 'write_document', 'list_headings',
  'edit_section', 'read_section', 'list_tasks', 'update_task',
  'get_breadcrumb', 'move_section', 'search_content', 'get_structure',
  'list_links', 'get_backlinks', 'get_graph',
  'get_neighbors', 'find_path', 'get_orphans', 'get_hub_notes',
];

function MCPSection(): React.JSX.Element {
  return (
    <section className={s.sectionAlt} id="mcp-server">
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.sectionLabel}>AI Integration</span>
        </div>
        <h2 className={s.sectionTitle}>An MCP server, built in.</h2>
        <p className={s.sectionDesc}>
          AI agents connect via stdio to read, navigate, and edit markdown documents programmatically.
        </p>
        <div className={s.mcpGrid}>
          <div className={s.mcpContent}>
            <h3>22 tools over JSON-RPC 2.0</h3>
            <p>
              Document tools for reading, writing, and searching. Navigation tools for
              heading-based traversal, task management, and section manipulation.
              Brain tools for wiki-link analysis, backlink discovery, and graph queries.
            </p>
            <div className={s.mcpTools}>
              {mcpTools.map(t => <span key={t} className={s.mcpTool}>{t}</span>)}
            </div>
            <Link className={s.btnSecondary} to="/docs/mcp-server/overview">
              MCP Documentation
            </Link>
          </div>
          <div className={s.mcpCode}>
            <div className={s.mcpCodeHeader} aria-label="Code examples for Claude Code and Gemini CLI">
              <span className={`${s.mcpCodeTab} ${s.mcpCodeTabActive}`} aria-current="true">Claude Code</span>
              <span className={s.mcpCodeTab}>Gemini CLI</span>
            </div>
            <div className={s.mcpCodeBody}>
              <code>{`# Add LazyMD as an MCP server
claude mcp add LazyMD -- \\
  /path/to/lm --mcp-server

# Start with a file preloaded
lm --mcp-server myfile.md`}</code>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── Teams (Coming Soon) ──────────────────────────────────────────── */

const teamFeatures = [
  {
    icon: '\u{1F91D}',
    title: 'Real-Time Collaboration',
    desc: 'Real-time multiplayer editing. See your team\'s cursors, edits, and selections live. Work together across continents.',
  },
  {
    icon: '\u{1F3C6}',
    title: 'Rankings & Leaderboards',
    desc: 'Track contributions, streaks, and impact across your organization. Gamify the writing experience.',
  },
  {
    icon: '\u23F1',
    title: 'Time Tracking',
    desc: 'Built-in time tracking per file, project, and team. Know where your hours actually go \u2014 no third-party tools.',
  },
  {
    icon: '\u{1F3E2}',
    title: 'Multi-Tenancy',
    desc: 'Isolated workspaces for teams and organizations. Role-based access, audit logs, and compliance-ready infrastructure.',
  },
];

function Teams(): React.JSX.Element {
  return (
    <section className={s.teamsSection} id="teams">
      <div className={s.teamsGlow} aria-hidden="true" />
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.teamsLabel}>Coming Soon</span>
        </div>
        <h2 className={s.sectionTitle}>
          Built for teams.
        </h2>
        <p className={s.sectionDesc}>
          Everything below ships with <strong>LazyMD Cloud</strong> &mdash; the managed
          version for organizations that want collaboration, analytics, and
          enterprise-grade infrastructure out of the box.
        </p>
        <div className={s.teamsGrid}>
          {teamFeatures.map(({icon, title, desc}) => (
            <article key={title} className={s.teamsCard}>
              <span className={s.teamsIcon} aria-hidden="true">{icon}</span>
              <h3>{title}</h3>
              <p>{desc}</p>
            </article>
          ))}
        </div>
        <div className={s.teamsNotify}>
          <p>Interested in early access?</p>
          <Link className={s.btnWarm} to="https://github.com/EME130/lazymd/discussions">
            Join the waitlist
          </Link>
        </div>
      </div>
    </section>
  );
}

/* ── Install ──────────────────────────────────────────────────────── */

function Install(): React.JSX.Element {
  return (
    <section className={s.section} id="installation">
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.sectionLabel}>Get Started</span>
        </div>
        <h2 className={s.sectionTitle}>Up and running in seconds.</h2>
        <p className={s.sectionDesc}>
          Compiles to a single binary with zero runtime dependencies.
        </p>
        <div className={s.installGrid}>
          <div className={s.installCard}>
            <div className={s.installStep} aria-hidden="true">1</div>
            <h3>Prerequisites</h3>
            <p>Install <a href="https://go.dev/dl/">Go</a> 1.24.2 or later from the official site.</p>
          </div>
          <div className={s.installCard}>
            <div className={s.installStep} aria-hidden="true">2</div>
            <h3>Build</h3>
            <div className={s.codeBlock}>
              <code>{`git clone https://github.com/\nEME130/lazymd.git\ncd lazymd && go build ./cmd/lm`}</code>
            </div>
          </div>
          <div className={s.installCard}>
            <div className={s.installStep} aria-hidden="true">3</div>
            <h3>Run</h3>
            <div className={s.codeBlock}>
              <code>{`./lm myfile.md`}</code>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── CTA ──────────────────────────────────────────────────────────── */

function CTASection(): React.JSX.Element {
  return (
    <section className={s.ctaSection}>
      <div className={s.ctaGlow} aria-hidden="true" />
      <h2 className={s.ctaTitle}>The future compiles from plain text.</h2>
      <p className={s.ctaDesc}>
        Open source, MIT licensed, and built to outlast every proprietary tool
        that came before it. LLVM for thought.
      </p>
      <div className={s.heroActions}>
        <Link className={s.btnPrimary} to="/docs/getting-started/installation">
          Get Started
        </Link>
        <Link className={s.btnSecondary} to="https://github.com/EME130/lazymd">
          Star on GitHub
        </Link>
      </div>
    </section>
  );
}

/* ── Page ──────────────────────────────────────────────────────────── */

export default function Home(): React.JSX.Element {
  return (
    <Layout
      title="The Compiler for Mortals"
      description="LazyMD is an LLVM-inspired markdown compiler for mortals. Plain text in, anything out. Runs everywhere — terminal, web, native apps, any device. Written in Go.">
      <Head>
        <html lang="en" />
      </Head>
      <main>
        <Hero />
        <div className={s.terminalDemo}>
          <TerminalDemo />
        </div>
        <SocialProof />
        <Philosophy />
        <Interfaces />
        <Features />
        <MCPSection />
        <Teams />
        <Install />
        <BuyMeACigarette />
        <CTASection />
      </main>
    </Layout>
  );
}
