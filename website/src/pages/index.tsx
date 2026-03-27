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
        <p className={s.heroDedication}>
          LazyMathDeath created LazyMathDragon. In loving memory of MBZ, for my
          Kleopetra &mdash; and in future, to my Mary.
        </p>
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

/* ── Vision: The Math Playground ──────────────────────────────────── */

const visionPillars = [
  {
    icon: '\u{1D70B}',
    title: 'LaTeX-Native Markdown',
    desc: 'Write theorems, proofs, and equations inline with your prose. Markdown extended with first-class LaTeX \u2014 no plugins, no hacks. The notation you think in is the notation you type.',
  },
  {
    icon: '\u{1F40D}',
    title: 'Embedded Python Interpreter',
    desc: 'Run computations right inside your document. Symbolic algebra, numerical proofs, visualizations \u2014 your markdown is alive. Write a conjecture, then prove it in the next cell.',
  },
  {
    icon: '\u{1F4DC}',
    title: 'Document Your Theorems',
    desc: 'Markdown is the connective tissue. Explain the intuition, state the theorem in LaTeX, verify it in Python, and publish it all as one artifact. The new generation of the research pad.',
  },
  {
    icon: '\u{1F52C}',
    title: 'Research Yourself',
    desc: 'No gatekeepers. No paywalls on your own tools. A super-fluid playground where students and researchers alike experiment with mathematics \u2014 for the living, and for those who came before us.',
  },
];

function Vision(): React.JSX.Element {
  return (
    <section className={s.visionSection} id="vision">
      <div className={s.visionGlow} aria-hidden="true" />
      <div className={s.container}>
        <div style={{textAlign: 'center'}}>
          <span className={s.visionLabel}>The Aim</span>
        </div>
        <h2 className={s.sectionTitle}>
          A super-fluid<br />math playground.
        </h2>
        <p className={s.sectionDesc}>
          We are building the instrument that Euler would have used if he had a
          laptop. Markdown for the narrative, LaTeX for the notation, Python for
          the proof. One document, one truth. Mathematics should be explored the
          way it was meant to be &mdash; fluidly, fearlessly, and without
          friction between thought and expression.
        </p>
        <div className={s.visionGrid}>
          {visionPillars.map(({icon, title, desc}) => (
            <article key={title} className={s.visionCard}>
              <span className={s.visionIcon} aria-hidden="true">{icon}</span>
              <h3>{title}</h3>
              <p>{desc}</p>
            </article>
          ))}
        </div>
        <p className={s.visionManifesto}>
          We believe the next breakthrough in mathematics won't come from a
          closed-source tool with a subscription fee. It will come from someone
          with a text file, an idea, and a compiler that sends it upward.
          LazyMD is not a runtime &mdash; runtimes are for the gods. This is a
          compiler. It takes the notes of a mathematician, compiles them through
          LaTeX and Python, and dispatches them to the gods to verify &mdash;
          first-class AI that proves, challenges, and extends your thinking.
          The Red Dragon rides again. Same energy as Aho, Sethi, and Ullman.
          Same lineage as the LLVM dragon. A new beast for a new era: the
          LazyMathDragon. For students scribbling their first proof. For
          researchers pushing the boundary. For the living and the dead &mdash;
          because great ideas deserve compilers that outlive their creators.
        </p>
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
    status: 'Available now',
    live: true,
    desc: 'The same editing experience in your browser. Real-time collaboration, cloud sync, and zero install. Open a link and start writing.',
    ascii: '  \u250C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510\n  \u2502 \u25CF\u25CF\u25CF  lazymd.com   \u2502\n  \u251C\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2524\n  \u2502 # Hello world.  \u2502\n  \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518',
  },
  {
    label: 'Native Apps',
    status: 'Available now',
    live: true,
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

/* ── Buy Me a Cigarette ───────────────────────────────────────────── */

const MONERO_ADDRESS =
  '84Jd3E44j54ZpmH2xAnJ6qJstnDtaJEmvK4pmvR78i5xcLGADnviwDpSa1uZMzCcrkVqH2u8E8hbBU4g4bn9sfB14t5Yjoi';

function BuyMeACigarette(): React.JSX.Element {
  const [copied, setCopied] = React.useState(false);

  const handleClick = () => {
    navigator.clipboard.writeText(MONERO_ADDRESS).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
    // Also try to open Monero client
    window.open(`monero:${MONERO_ADDRESS}`, '_self');
  };

  return (
    <section className={s.section} style={{textAlign: 'center'}}>
      <div className={s.container}>
        <span className={s.sectionLabel}>Support</span>
        <h2 className={s.sectionTitle}>Buy me a cigarette.</h2>
        <p className={s.sectionDesc}>
          If LazyMD saves you time, consider sending some Monero.
        </p>
        <button className={s.btnCigarette} onClick={handleClick}>
          {copied ? 'Address Copied!' : '\uD83D\uDEAC Buy Me a Cigarette'}
        </button>
        <p className={s.moneroAddr}>{MONERO_ADDRESS}</p>
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
        <Vision />
        <Interfaces />
        <Install />
        <BuyMeACigarette />
        <CTASection />
      </main>
    </Layout>
  );
}
