<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EditorView, basicSetup } from 'codemirror';
  import { markdown } from '@codemirror/lang-markdown';
  import { oneDark } from '@codemirror/theme-one-dark';
  import { EditorState } from '@codemirror/state';
  import { GetContent, InsertText, DeleteRange, SetCursor, SaveFile } from '../../wailsjs/go/wailsplugin/App';
  import { onBufferChanged, onFileOpened } from '../lib/events';

  let editorContainer: HTMLDivElement;
  let view: EditorView;
  let ignoreNextUpdate = false;

  onMount(async () => {
    const content = await GetContent();

    const updateListener = EditorView.updateListener.of((update) => {
      if (!update.docChanged) return;
      if (ignoreNextUpdate) {
        ignoreNextUpdate = false;
        return;
      }

      update.changes.iterChanges((fromA, toA, fromB, toB, inserted) => {
        const fromPos = view.state.doc.lineAt(fromA);
        const toPos = view.state.doc.lineAt(toA);

        if (toA > fromA) {
          DeleteRange(
            fromPos.number - 1, fromA - fromPos.from,
            toPos.number - 1, toA - toPos.from,
          );
        }

        const text = inserted.toString();
        if (text.length > 0) {
          InsertText(fromPos.number - 1, fromA - fromPos.from, text);
        }
      });
    });

    const state = EditorState.create({
      doc: content,
      extensions: [
        basicSetup,
        markdown(),
        oneDark,
        updateListener,
        EditorView.theme({
          '&': { height: '100%' },
          '.cm-scroller': { overflow: 'auto' },
        }),
      ],
    });

    view = new EditorView({
      state,
      parent: editorContainer,
    });

    onBufferChanged((data) => {
      if (data?.origin === 'wails-gui') return;
      const content = data?.content;
      if (content != null && view) {
        ignoreNextUpdate = true;
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: content },
        });
      }
    });

    onFileOpened(async () => {
      const content = await GetContent();
      if (view) {
        ignoreNextUpdate = true;
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: content },
        });
      }
    });

    editorContainer.addEventListener('keydown', (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 's') {
        e.preventDefault();
        SaveFile();
      }
    });
  });

  onDestroy(() => {
    if (view) view.destroy();
  });
</script>

<div class="editor-container" bind:this={editorContainer}></div>

<style>
  .editor-container {
    height: 100%;
    overflow: hidden;
  }
  .editor-container :global(.cm-editor) {
    height: 100%;
  }
</style>
