import {writable, get} from "svelte/store"
import { emitter } from "./notifs";
import {chunk, reverse, map, flatten} from "lodash-es"

export const hostname = writable("192.168.1.42");
if(window.location.hostname !== "localhost") {
    hostname.set(window.location.hostname)
}

const displayHatLocked = () => {
    alert("The hat is currently locked, please ask someone who knows the password to unlock it")
}

export const setMode = async (m) => {
    let resp = await fetch(`http://${get(hostname)}/api?mode=`+m)
    if(await resp.text() === "LOCKED") {
        displayHatLocked()
    }
}

export const showArt = async (a) => {
    let resp = await fetch(`http://${get(hostname)}/api?art=`+a)
    if(await resp.text() === "LOCKED") {
        displayHatLocked()
    }
}

export const showGif = async (a) => {
    let resp = await fetch(`http://${get(hostname)}/api?gif=`+a)
    if(await resp.text() === "LOCKED") {
        displayHatLocked()
    }
}

export const fillBoard = async ({r, g, b}) => {
    let color = (r << 16) + (g << 8) + b
    let resp = await fetch(`http://${get(hostname)}/api/fill?color=`+color)
    if(await resp.text() === "LOCKED") {
        displayHatLocked()
    }
}

export const getArt = async () => {
    let a = await fetch(`http://${get(hostname)}/api/art`)
    let {data} = await a.json()
    return data
}

export const getGifs = async () => {
    let a = await fetch(`http://${get(hostname)}/api/gifs`)
    let {data} = await a.json()
    return data
}

export const toggleLock = async (pw: string) => {
    let formData = new FormData()
    formData.append("pass", pw)
    let resp = await fetch(`http://${get(hostname)}/lock`, {
        method: "POST",
        mode: "cors",
        body: formData
    })
    emitter.emit("LOCK", await resp.text())
}

//pixelArr in this case takes the left to right, top to bottom type of array
export const sendToBoard = async (pixelArr: number[]) => {
    let boardDTO = []
    boardDTO = chunk(pixelArr, 14)
    boardDTO = reverse(boardDTO)
    boardDTO = map(boardDTO, (v, idx) => idx % 2 !== 0 ? reverse(v) : v)
    boardDTO = flatten(boardDTO)
    const formData = new FormData()
    formData.append("pixels", JSON.stringify({data: boardDTO}))
    let resp = await fetch(`http://${get(hostname)}/api/draw`, {
        method: "POST",
        mode: "cors",
        body: formData
    })
    if(await resp.text() === "LOCKED") {
        displayHatLocked()
    }
}