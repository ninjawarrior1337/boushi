<template>
    <div transition:fly class="fixed inset-0 bg-black bg-opacity-60 grid place-items-center">
        <div class="flex flex-col items-center justify-center space-y-8 shadow-2xl rounded-xl bg-dark-300 text-white p-12 mx-5 opacity-100">
            <h2 class="text-5xl text-center">Enter Password</h2>
            <input type="password" class="bg-dark-100 p-4 shadow-xl" placeholder="**********" bind:value={password}>
            <div>
                {state}
            </div>
            <div class="grid lg:grid-cols-2 gap-4 text-xl">
                <button on:click={() => toggleLock(password)} class="p-4 bg-gradient-to-r from-muse via-treelar to-treelar rounded shadow-xl font-semibold">Toggle Lock</button>
                <button class="p-4 bg-red-500 rounded shadow-xl font-semibold"on:click={dismiss}>Dismiss</button>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { onDestroy } from "svelte";

import {fly} from "svelte/transition"
import { toggleLock } from "./hat_api";
import { emitter } from "./notifs";
let password = ""
export let dismiss: () => void;
let state = ""
const onLockStatus = (v) => {
    console.log(v)
    state = v
}
emitter.on("LOCK", onLockStatus)
onDestroy(() => {
    emitter.off("LOCK", onLockStatus)
})
</script>
