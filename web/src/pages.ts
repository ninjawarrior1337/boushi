import {derived, writable} from "svelte/store"

export const pages = ["Presets", "GIFs","Fill Board", "Draw"]
export const selectedPageIndex = writable(0)
export const activePage = derived(selectedPageIndex, (i) => pages[i])