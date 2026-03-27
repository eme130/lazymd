import { GetThemeColors } from '../../wailsjs/go/wailsplugin/App';

export async function applyTheme() {
  const colors = await GetThemeColors();
  const root = document.documentElement;
  for (const [key, value] of Object.entries(colors)) {
    if (key.startsWith('--')) {
      root.style.setProperty(key, value);
    }
  }
}
