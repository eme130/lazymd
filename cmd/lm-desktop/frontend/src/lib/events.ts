import { EventsOn } from '../../wailsjs/runtime/runtime';

export function onBufferChanged(callback: (data: Record<string, any>) => void) {
  return EventsOn('buffer:changed', callback);
}

export function onFileOpened(callback: (data: Record<string, any>) => void) {
  return EventsOn('file:opened', callback);
}

export function onFileSaved(callback: (data: Record<string, any>) => void) {
  return EventsOn('file:saved', callback);
}

export function onCursorMoved(callback: (data: Record<string, any>) => void) {
  return EventsOn('cursor:moved', callback);
}

export function onGraphUpdated(callback: (data: Record<string, any>) => void) {
  return EventsOn('graph:updated', callback);
}
