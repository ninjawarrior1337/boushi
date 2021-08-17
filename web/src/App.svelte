<template>
    <div class="w-full min-h-screen flex flex-col items-center bg-dark-800 text-white">
        <div class="w-full flex items-center py-6 px-12">
            <img src={logo} alt="treelar logo" class="w-12">
            <h1 class="text-transparent font-semibold bg-clip-text bg-gradient-to-r from-muse via-treelar to-treelar text-5xl mx-auto">
                Boushi
            </h1>
            <button on:click={() => lockModalVisible = true}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
            </button>
        </div>
        <div class="container px-4 lg:px-12 pb-12">
            <div class="flex space-x-4 text-xl border-b-2">
                {#each pages as page, index}
                    <button on:click={() => selectedPageIndex.set(index)} class="bg-dark-200 p-2 rounded-x rounded-t">{page}</button>
                {/each}
            </div>
            {#if $activePage === "Presets"}
            {#await import("./Presets.svelte") then P}
                <svelte:component this={P.default}></svelte:component>
            {/await}
            {/if}
            {#if $activePage === "Fill Board"}
            {#await import("./FillBoard.svelte") then F}
                <svelte:component this={F.default}></svelte:component>
            {/await}
            {/if}
            {#if $activePage === "Draw"}
            {#await import("./Canvas.svelte") then C}
                <svelte:component this={C.default}></svelte:component>
            {/await}
            {/if}
            {#if $activePage === "GIFs"}
            {#await import("./Gifs.svelte") then G}
                <svelte:component this={G.default}></svelte:component>
            {/await}
            {/if}
        </div>
    </div>
    {#if lockModalVisible}
        <LockModal dismiss={() => lockModalVisible = false}></LockModal>
    {/if}
</template>

<script lang="ts">
import { onMount } from "svelte";
import iro from "@jaames/iro"
import { writable } from "svelte/store";
import Canvas from "./Canvas.svelte";
import logo from "../assets/logo2020.svg"
import { getArt, showArt, fillBoard } from "./hat_api";
import LockModal from "./LockModal.svelte";
import { selectedPageIndex, pages, activePage } from "./pages";

let lockModalVisible = false;
</script>

<style>
.hidden {
    display: none;
}
</style>