import axios from 'axios'
const basePath = "/api/canvas"

let initialized = false

export default {
    init() {
        initialized = true
    },
    setColor(color) {
        if (initialized) {
            axios.post(`${basePath}/color`, {color: color})
        }
    },
    undo() {
        if (initialized) {
            axios.post(`${basePath}/undo`, {undo: true})
        }
    },
    setTool(toolName) {
        if (initialized) {
            axios.post(`${basePath}/tool`, {toolName: toolName})
        }
    },
    reset() {
        if (initialized) {
            axios.post(`${basePath}/reset`, {reset: true})
        }
    },
    download(info) {
        if (initialized) {
            const url = `${basePath}/download?pixelSize=${info.pixelSize}&fileName=${info.fileName}`

            axios
                .request({
                    url: url,
                    method: 'GET',
                    responseType: "blob"
                })
                .then(response => {
                    const fileURL = window.URL.createObjectURL(
                        new Blob([response.data]),
                        {
                            type: response.headers["content-type"]
                        }
                    );
                    const fileLink = document.createElement("a");
                    fileLink.href = fileURL;
                    const fileNames = response.headers["content-disposition"].match(
                        /filename="([^"]+)"/
                    )

                    const fileName = fileNames.length > 1 ? fileNames[1] : "untitled.png"
                    fileLink.setAttribute("download", fileName);
                    document.body.appendChild(fileLink);

                    fileLink.click();
                    fileLink.remove();
               })
                .catch(() => {
                });
        }
    },
}
