import {createStore} from 'vuex'

export const store = createStore({
    state: {initializing: true},
    mutations: {
        replace(state, data) {
            let newState = Object.assign({}, state)

            if (data.canvas) {
                newState.canvas = data.canvas
            } else if(data.pixels) {
                for (const pixel of data.pixels) {
                    newState.canvas[pixel.y][pixel.x] = pixel.color
                }
            }

            if (data.window) {
                const win = data.window
                newState.window = {
                    top: win.y,
                    left: win.x,
                    bottom: win.y + 7,
                    right: win.x + 7,
                }
            }

            if (data.cursor) {
                newState.cursor = Object.assign({}, data.cursor)
            }
            if (data.color) {
                newState.color = data.color
            }

            if (data.toolName) {
                newState.tool = data.toolName
                let toolChar
                switch (data.toolName){
                    case "pen": toolChar = "+"; break;
                    case "eraser": toolChar = "x"; break;
                    case "bucket": toolChar = "o"; break;
                    default: toolChar = "?"; break;
                }
                newState.toolChar = toolChar
            }

            newState.initializing = false

            this.replaceState(newState)
        }
    },
    actions: {},
    modules: {
    }
})
