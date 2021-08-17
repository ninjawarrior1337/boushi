<template>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-y-8 justify-items-center pt-4">
        <div>
            <div class="grid w-full grid-cols-14 grid-rows-12 gap-2 select-none">
                {#each pixel as p, idx}
                <div 
                class="w-4 h-3 p-3" 
                on:mousedown={() => drawPixel(idx, true)}
                on:mouseover={() => drawPixel(idx, false)}
                style="background-color: {rgb32AsHex(p)}">
                </div>
                {/each}
            </div>
        </div>
        <div class="grid lg:grid-cols-2 gap-6 items-center">
            <div bind:this={pickerEl}></div>
            <div class="flex flex-col space-y-4">
                <button on:click={() => sendToBoard(pixel)} class="bg-gradient-to-r from-muse to-treelar via-treelar p-4 text-4xl rounded">Send to Board</button>
                <button on:click={clearBoard} class="bg-gradient-to-r from-muse to-red-500 p-4 text-4xl rounded">Clear Board</button>
            </div>
        </div>
    </div>
</template>

<svelte:window on:mousedown={() => isMouseDown = true} on:mouseup={() => isMouseDown = false}></svelte:window>

<script lang="ts">
import iro from "@jaames/iro";
import { onMount } from "svelte";
import { sendToBoard } from "./hat_api";
let pixel = new Array(168)
let isMouseDown = false
let pickerEl: HTMLElement;
let colorPicker;

const clearBoard = () => {
    pixel = pixel.fill(0, 0, 168)
}

clearBoard()

onMount(() => {
    colorPicker = iro.ColorPicker(pickerEl, {})
})

const rgb32AsHex = (pixel: number) => {
    return "#"+pixel.toString(16).padStart(6, "0")
}
const drawPixel = (i: number, bypass: boolean) => {
    let {r, g, b} = colorPicker.color.rgb
    let currentColor = (r << 16) + (g << 8) + b
    if(isMouseDown || bypass) {
        pixel = pixel.map((v, idx) => {
            if(v != currentColor) {
                if (idx === i) return currentColor; else return v
            } else {
                if (idx === i) return 0; else return v
            }
        })
    }
}
</script>

<style>
</style>