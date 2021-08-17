<div class="grid grid-cols-1 lg:grid-cols-3 gap-8 justify-items-center pt-4">
    {#each art as a}
        <button on:click={() => showGif(a)} class="bg-blue-600 w-full p-4 rounded text-xl">{a}</button>
    {/each}
</div>

<script lang="ts">
import { onMount } from "svelte";

import { getGifs, showGif } from "./hat_api";
let art = []
onMount(async () => {
    art = await getGifs()
})
</script>