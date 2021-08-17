<div class="grid grid-cols-1 lg:grid-cols-2 gap-y-8 justify-items-center pt-4">
    <div class="col-span-2 lg:col-span-1" bind:this={pickerEl}></div>
    <div class="flex items-center">
        <button on:click={fillB} class="bg-gradient-to-r from-muse to-treelar via-treelar p-4 text-4xl rounded">Fill Board</button>
    </div>
</div>

<script lang="ts">
import iro from "@jaames/iro";
import { onMount } from "svelte";
import { fillBoard } from "./hat_api";
let colorPicker: iro.ColorPicker
let pickerEl: HTMLElement
const fillB = () => {
    fillBoard(colorPicker.color.rgb)
}
onMount(() => {
    colorPicker = iro.ColorPicker(pickerEl, {})
})
</script>