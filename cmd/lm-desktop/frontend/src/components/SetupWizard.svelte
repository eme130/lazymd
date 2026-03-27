<script>
  import { SelectVaultDir, SaveVault } from '../../wailsjs/go/wailsplugin/App';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let vaultPath = '';
  let error = '';
  let saving = false;

  async function browse() {
    try {
      const path = await SelectVaultDir();
      if (path) {
        vaultPath = path;
        error = '';
      }
    } catch (e) {
      error = 'Failed to open directory picker';
    }
  }

  async function confirm() {
    if (!vaultPath) {
      error = 'Please select a directory';
      return;
    }
    saving = true;
    try {
      await SaveVault(vaultPath);
      dispatch('complete');
    } catch (e) {
      error = `Failed to save: ${e}`;
      saving = false;
    }
  }
</script>

<div class="setup">
  <h1>Welcome to LazyMD</h1>
  <p class="subtitle">Select your vault directory — where your notes live.</p>

  <div class="input-row">
    <input type="text" bind:value={vaultPath} placeholder="~/notes" />
    <button on:click={browse}>Browse</button>
  </div>

  {#if error}
    <p class="error">{error}</p>
  {/if}

  <button class="confirm" on:click={confirm} disabled={saving}>
    {saving ? 'Setting up...' : 'Get Started'}
  </button>
</div>

<style>
  .setup {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    background: var(--lm-bg, #1a1b26);
    color: var(--lm-fg, #a9b1d6);
  }
  h1 { color: var(--lm-link, #7aa2f7); margin-bottom: 0.5rem; }
  .subtitle { color: var(--lm-muted, #565f89); margin-bottom: 2rem; }
  .input-row {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  input {
    background: var(--lm-editor-bg, #24283b);
    border: 1px solid var(--lm-border, #3b4261);
    color: var(--lm-fg, #a9b1d6);
    padding: 0.5rem 1rem;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 1rem;
    width: 400px;
    border-radius: 4px;
  }
  input:focus { border-color: var(--lm-link, #7aa2f7); outline: none; }
  button {
    background: var(--lm-border, #3b4261);
    color: var(--lm-fg, #a9b1d6);
    border: none;
    padding: 0.5rem 1rem;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    cursor: pointer;
    border-radius: 4px;
  }
  button:hover { background: #414868; }
  .confirm {
    background: var(--lm-link, #7aa2f7);
    color: var(--lm-bg, #1a1b26);
    font-weight: bold;
    padding: 0.75rem 2rem;
    margin-top: 1rem;
  }
  .confirm:hover { background: #89b4fa; }
  .confirm:disabled { opacity: 0.5; cursor: default; }
  .error { color: #f7768e; margin: 0; }
</style>
